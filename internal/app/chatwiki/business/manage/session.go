// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetSessionChannelList(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	robotId := cast.ToUint(c.Query(`robot_id`))
	list := common.GetChannelList(userId, robotId)
	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
}

func GetSessionRecordList(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	robotId := cast.ToUint(c.Query(`robot_id`))
	appType := strings.TrimSpace(c.Query(`app_type`))
	appId := strings.TrimSpace(c.Query(`app_id`))
	startTime := cast.ToInt(c.Query(`start_time`))
	endTime := cast.ToInt(c.Query(`end_time`))
	name := strings.TrimSpace(c.Query(`name`))
	page := max(1, cast.ToInt(c.Query(`page`)))
	size := max(1, cast.ToInt(c.Query(`size`)))
	//get session_id list
	m := msql.Model(`chat_ai_session`, define.Postgres).Alias(`s`).
		Join(`chat_ai_dialogue d`, `s.dialogue_id=d.id`, `left`)
	m.Where(`s.admin_user_id`, cast.ToString(userId))
	m.Where(`s.robot_id`, cast.ToString(robotId))
	m.Where(`d.is_background`, `0`)
	if len(appType) > 0 {
		m.Where(`s.app_type`, appType)
	}
	if len(appId) > 0 {
		m.Where(`s.app_id`, appId)
	}
	if startTime > 0 && endTime > 0 && endTime >= startTime {
		m.Where(`s.last_chat_time`, `between`, fmt.Sprintf(`%d,%d`, startTime, endTime))
	}
	if len(name) > 0 {
		m.Join(`chat_ai_customer c`, fmt.Sprintf(`c.admin_user_id=%d AND c.openid=s.openid`, userId), `left`)
		m.Where(`c.name`, `like`, name)
	}
	m.Limit(size*(page-1), size).Group(`s.openid`).Order(`max(s.last_chat_time) DESC`)
	sessionIds, err := m.ColumnArr(`max(s.id) session_id`)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	//get record
	list := make([]msql.Params, 0)
	if len(sessionIds) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(list, nil))
		return
	}
	list, err = m.Where(`id`, `in`, strings.Join(sessionIds, `,`)).
		Field(`id session_id,dialogue_id,last_chat_time,last_chat_message,app_type,app_id,openid,rel_user_id`).
		Order(`last_chat_time DESC`).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(list) > 0 { //customer info
		relUserIds := make([]string, 0)
		for _, one := range list {
			if cast.ToInt(one[`rel_user_id`]) > 0 && !tool.InArray(one[`rel_user_id`], relUserIds) {
				relUserIds = append(relUserIds, one[`rel_user_id`])
			}
		}
		userInfoMaps := make(map[string]msql.Params)
		if len(relUserIds) > 0 {
			userInfoMaps, err = msql.Model(define.TableUser, define.Postgres).
				Where(`id`, `in`, strings.Join(relUserIds, `,`)).ColumnMap(`user_name,nick_name,avatar,id,parent_id`, `id`)
			if err != nil {
				logs.Error(err.Error())
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
				return
			}
			for userId, one := range userInfoMaps {
				if one[`parent_id`] == userId || one[`id`] == userId {
					continue
				}
				delete(userInfoMaps, userId)
			}
		}
		for i, one := range list {
			if cast.ToInt(one[`rel_user_id`]) > 0 {
				common.FillRelUserInfo2(list[i], cast.ToInt(one[`rel_user_id`]), userInfoMaps[one[`rel_user_id`]])
			} else if customer, _ := common.GetCustomerInfo(one[`openid`], userId); len(customer) > 0 {
				list[i][`name`] = customer[`name`]
				list[i][`avatar`] = customer[`avatar`]
			} else {
				list[i][`name`] = `访客XXXX`
				list[i][`avatar`] = define.DefaultCustomerAvatar
			}
		}
	}
	data := map[string]any{`list`: list, `page`: page, `size`: size, `has_more`: len(list) == size}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func CreateSessionExport(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	robotId := cast.ToUint(c.Query(`robot_id`))
	fileName := fmt.Sprintf(`会话记录%s.xlsx`, tool.Date(`YmdHis`))
	params := map[string]any{
		`admin_user_id`: userId,
		`robot_id`:      robotId,
		`app_type`:      strings.TrimSpace(c.Query(`app_type`)),
		`app_id`:        strings.TrimSpace(c.Query(`app_id`)),
		`start_time`:    cast.ToInt(c.Query(`start_time`)),
		`end_time`:      cast.ToInt(c.Query(`end_time`)),
		`name`:          strings.TrimSpace(c.Query(`name`)),
	}
	id, err := common.CreateExportTask(uint(userId), robotId, define.ExportSourceSession, fileName, params)
	c.String(http.StatusOK, lib_web.FmtJson(id, err))
}
