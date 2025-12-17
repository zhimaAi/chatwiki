// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetReceiverList(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	robotId := cast.ToUint(c.Query(`robot_id`))
	page := max(1, cast.ToInt(c.Query(`page`)))
	size := max(1, cast.ToInt(c.Query(`size`)))
	m := msql.Model(`chat_ai_receiver`, define.Postgres).Where(`admin_user_id`, cast.ToString(userId))
	if robotId > 0 {
		m.Where(`robot_id`, cast.ToString(robotId))
	}
	keyword := strings.TrimSpace(c.Query(`keyword`))
	if len(keyword) > 0 {
		m.Where(`openid|nickname|name`, `like`, keyword)
	}
	list, err := m.Order(`id DESC`).Limit(size*(page-1), size).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	data := map[string]any{`list`: list, `page`: page, `size`: size, `has_more`: len(list) == size}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func SetReceiverRead(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToUint(c.PostForm(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	m := msql.Model(`chat_ai_receiver`, define.Postgres)
	info, err := m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(userId)).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(info, nil))
		return
	}
	upData := msql.Datas{`unread`: 0, `update_time`: tool.Time2Int()}
	if _, err = m.Where(`id`, info[`id`]).Update(upData); err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	for key, val := range upData {
		info[key] = cast.ToString(val)
	}
	//websocket notify
	common.ReceiverChangeNotify(userId, `update`, info)
	c.String(http.StatusOK, lib_web.FmtJson(info, nil))
}
