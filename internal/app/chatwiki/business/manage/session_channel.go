// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_define"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

type sessionChannelCount struct {
	common.SessionChannel
	Count int `json:"count"`
}

func sessionChannelsWithCounts(channels []common.SessionChannel, counts []msql.Params, lang string) []sessionChannelCount {
	list := make([]sessionChannelCount, 0, len(channels)+len(counts))
	indexes := make(map[string]int)
	for _, channel := range channels {
		key := channel.AppType + `::` + channel.AppId
		if _, exists := indexes[key]; exists {
			continue
		}
		indexes[key] = len(list)
		list = append(list, sessionChannelCount{SessionChannel: channel})
	}
	for _, count := range counts {
		key := count[`app_type`] + `::` + count[`app_id`]
		index, exists := indexes[key]
		if !exists {
			name := count[`app_type`]
			if tool.InArray(name, lib_define.AppTypeList) {
				name = i18n.Show(lang, `channel_`+name)
			}
			index = len(list)
			indexes[key] = index
			list = append(list, sessionChannelCount{SessionChannel: common.SessionChannel{
				AppType: count[`app_type`], AppId: count[`app_id`], AppName: name,
			}})
		}
		list[index].Count = cast.ToInt(count[`total`])
	}
	return list
}

func getSessionChannels(ctx *gin.Context, userId int, robotId uint) ([]sessionChannelCount, int, error) {
	query := func() *msql.Builder {
		model := msql.Model(`chat_ai_session`, define.Postgres).Alias(`s`).
			Join(`chat_ai_dialogue d`, `s.dialogue_id=d.id`, `left`).
			Where(`s.admin_user_id`, cast.ToString(userId)).
			Where(`s.robot_id`, cast.ToString(robotId)).Where(`d.is_background`, `0`)
		startTime, endTime := cast.ToInt(ctx.Query(`start_time`)), cast.ToInt(ctx.Query(`end_time`))
		if startTime > 0 && endTime >= startTime {
			model.Where(`s.last_chat_time`, `between`, fmt.Sprintf(`%d,%d`, startTime, endTime))
		}
		if name := strings.TrimSpace(ctx.Query(`name`)); name != `` {
			model.Join(`chat_ai_customer customer`, fmt.Sprintf(`customer.admin_user_id=%d AND customer.openid=s.openid`, userId), `left`).
				Where(`customer.name`, `like`, name)
		}
		return model
	}
	counts, err := query().Field(`s.app_type,s.app_id,count(distinct s.openid) total`).Group(`s.app_type,s.app_id`).Select()
	if err != nil {
		return nil, 0, err
	}
	list := sessionChannelsWithCounts(common.GetChannelList(userId, robotId), counts, common.GetLang(ctx))
	if ctx.Query(`with_stats`) != `1` {
		return list, 0, nil
	}
	total, err := query().Count(`distinct s.openid`)
	return list, total, err
}
