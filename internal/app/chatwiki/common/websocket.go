// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_redis"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

type UsersCacheBuildHandler struct{ ParentId int }

func (h *UsersCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.users.%d`, h.ParentId)
}
func (h *UsersCacheBuildHandler) GetCacheData() (any, error) {
	return msql.Model(define.TableUser, define.Postgres).Where(`parent_id`, cast.ToString(h.ParentId)).ColumnArr(`id`)
}

func GetUsersInfo(parentId int) ([]string, error) {
	result := make([]string, 0)
	err := lib_redis.GetCacheWithBuild(define.Redis, &UsersCacheBuildHandler{ParentId: parentId}, &result, time.Minute)
	return result, err
}

func ReceiverChangeNotify(adminUserId int, change string, data any) {
	users, err := GetUsersInfo(adminUserId)
	if err != nil {
		logs.Error(err.Error())
	}
	openid := strings.Join(append([]string{cast.ToString(adminUserId)}, users...), `,`)
	message := msql.Datas{`msg_type`: `receiver_notify`, `change`: change, `data`: data}
	messageStr, _ := tool.JsonEncode(map[string]any{`openid`: openid, `message`: message})
	_ = AddJobs(lib_define.WsMessagePushTopic, messageStr)
	if define.IsDev {
		logs.Other(`websocket`, messageStr)
	}
}
