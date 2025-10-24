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

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
)

func GetExportTaskList(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	robotId := cast.ToUint(c.Query(`robot_id`))
	libraryId := cast.ToUint(c.Query(`library_id`))
	startTime := cast.ToInt(c.Query(`start_time`))
	endTime := cast.ToInt(c.Query(`end_time`))
	page := max(1, cast.ToInt(c.Query(`page`)))
	size := max(1, cast.ToInt(c.Query(`size`)))
	source := cast.ToInt(c.DefaultQuery(`source`, cast.ToString(define.ExportSourceSession)))
	m := msql.Model(`chat_ai_export_task`, define.Postgres).
		Field(`id,create_time,file_name,source,status,file_url,err_msg`).
		Where(`admin_user_id`, cast.ToString(userId)).Where(`robot_id`, cast.ToString(robotId)).
		Where(`source`, cast.ToString(source))
	if libraryId > 0 {
		m.Where(`library_id`, cast.ToString(libraryId))
	}
	if startTime > 0 && endTime > 0 && endTime >= startTime {
		m.Where(`create_time`, `between`, fmt.Sprintf(`%d,%d`, startTime, endTime))
	}
	list, total, err := m.Order(`id desc`).Paginate(page, size)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	data := map[string]any{`map`: define.ExportSourceList, `list`: list, `total`: total, `page`: page, `size`: size}
	c.String(http.StatusOK, lib_web.FmtJson(data, nil))
}

func DownloadExportFile(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToUint(c.Query(`id`))
	robotId := cast.ToUint(c.Query(`robot_id`))
	info, err := msql.Model(`chat_ai_export_task`, define.Postgres).Where(`id`, cast.ToString(id)).
		Where(`admin_user_id`, cast.ToString(userId)).Where(`robot_id`, cast.ToString(robotId)).
		Where(`status`, cast.ToString(define.ExportStatusSucceed)).Field(`file_name,file_url`).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(info) == 0 || !common.LinkExists(info[`file_url`]) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	c.FileAttachment(common.GetFileByLink(info[`file_url`]), info[`file_name`])
}
