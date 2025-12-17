// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_web"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetMcpProviderList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	query := msql.Model(`mcp_provider`, define.Postgres).Where(`admin_user_id`, cast.ToString(adminUserId))
	hasAuth := cast.ToInt(c.Query("has_auth"))
	if hasAuth == 1 {
		query = query.Where(`has_auth`, cast.ToString(hasAuth))
	}

	list, err := query.Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
}

func GetMcpProviderDetail(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	providerId := cast.ToInt64(c.Query(`provider_id`))
	if providerId <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}

	list, err := msql.Model(`mcp_provider`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(providerId)).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
}

func SaveMcpProvider(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}

	providerId := cast.ToInt64(c.PostForm(`provider_id`))
	name := strings.TrimSpace(c.PostForm(`name`))
	description := strings.TrimSpace(c.PostForm(`description`))
	url := strings.TrimSpace(c.PostForm(`url`))
	requestTimeout := cast.ToInt(c.PostForm(`request_timeout`))
	headers := strings.TrimSpace(c.PostForm(`headers`))

	// check params
	if len(name) == 0 || len(description) == 0 || len(url) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if utf8.RuneCountInString(name) > 100 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `name`))))
		return
	}
	if utf8.RuneCountInString(description) > 500 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `description`))))
		return
	}
	if !common.IsUrl(url) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `url`))))
		return
	}
	clientType, err := common.DetectMCPTransportType(url)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `url`))))
		return
	}
	if requestTimeout <= 0 || requestTimeout > 500 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `request_timeout`))))
		return
	}
	if len(headers) > 0 {
		var js map[string]string
		if err := json.Unmarshal([]byte(headers), &js); err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `headers`))))
			return
		}
	} else {
		headers = "{}"
	}

	// upload file
	avatar := ``
	if providerId == 0 {
		avatar = define.LocalUploadPrefix + `default/mcp_avatar.svg`
	}
	fileHeader, _ := c.FormFile(`avatar`)
	uploadInfo, err := common.SaveUploadedFile(fileHeader, define.ImageLimitSize, adminUserId, `mcp_provider_avatar`, define.ImageAllowExt)
	if err == nil && uploadInfo != nil {
		avatar = uploadInfo.Link
	}
	if providerId > 0 { // edit
		mcpProvider, err := msql.Model(`mcp_provider`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`id`, cast.ToString(providerId)).
			Find()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if len(mcpProvider) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
			return
		}
		updateData := msql.Datas{
			`name`:            name,
			`description`:     description,
			`url`:             url,
			`request_timeout`: requestTimeout,
			`headers`:         headers,
			`client_type`:     clientType,
			`update_time`:     tool.Time2Int(),
		}
		if len(avatar) > 0 {
			updateData[`avatar`] = avatar
		}
		_, err = msql.Model(`mcp_provider`, define.Postgres).
			Where(`id`, cast.ToString(providerId)).
			Update(updateData)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
	} else {
		providerId, err = msql.Model(`mcp_provider`, define.Postgres).Insert(msql.Datas{
			`admin_user_id`:   adminUserId,
			`avatar`:          avatar,
			`name`:            name,
			`description`:     description,
			`url`:             url,
			`request_timeout`: requestTimeout,
			`headers`:         headers,
			`client_type`:     clientType,
			`create_time`:     tool.Time2Int(),
			`update_time`:     tool.Time2Int(),
		}, `id`)
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}

	}
	resp := map[string]int64{"provider_id": providerId}
	c.String(http.StatusOK, lib_web.FmtJson(resp, nil))
}
func AuthMcpProvider(c *gin.Context) {
	adminUserId := GetAdminUserId(c)
	if adminUserId == 0 {
		return
	}

	providerId := cast.ToInt(c.PostForm("provider_id"))
	provider, err := msql.Model("mcp_provider", define.Postgres).
		Where("admin_user_id", cast.ToString(adminUserId)).
		Where("id", cast.ToString(providerId)).
		Find()
	if err != nil {
		sendError(c, "sys_err", err)
		return
	}
	if len(provider) == 0 {
		sendError(c, "no_data", nil)
		return
	}

	ctx := context.Background()

	mcpClient, err := common.NewMcpClient(ctx,
		cast.ToInt(provider["client_type"]),
		provider["url"],
		provider["headers"],
	)
	if err != nil {
		_ = updateProviderAuth(providerId, false, "[]")
		sendError(c, "auth_fail", err)
		return
	}

	tools, err := common.ListTools(ctx, mcpClient)
	if err != nil {
		_ = updateProviderAuth(providerId, false, "[]")
		sendError(c, "auth_fail", err)
		return
	}

	toolsStr, err := json.Marshal(tools)
	if err != nil {
		_ = updateProviderAuth(providerId, false, "[]")
		sendError(c, "sys_err", err)
		return
	}

	if err := updateProviderAuth(providerId, true, string(toolsStr)); err != nil {
		sendError(c, "sys_err", err)
		return
	}

	provider["has_auth"] = "1"
	provider["tools"] = string(toolsStr)

	c.String(http.StatusOK, lib_web.FmtJson(provider, nil))
}

// updateProviderAuth 封装更新 provider 状态
func updateProviderAuth(providerId int, hasAuth bool, tools string) error {
	auth := 0
	if hasAuth {
		auth = 1
	}
	_, err := msql.Model("mcp_provider", define.Postgres).
		Where("id", cast.ToString(providerId)).
		Update(msql.Datas{
			"tools":       tools,
			"has_auth":    auth,
			"update_time": tool.Time2Int(),
		})
	if err != nil {
		logs.Error(err.Error())
	}
	return err
}

func sendError(c *gin.Context, code string, err error) {
	if err != nil {
		logs.Error(err.Error())
	}
	msg := i18n.Show(common.GetLang(c), code)
	c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(msg)))
}

func DeleteMcpProvider(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	providerId := cast.ToInt(c.PostForm(`provider_id`))
	provider, err := msql.Model(`mcp_provider`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`id`, cast.ToString(providerId)).
		Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(provider) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	_, err = msql.Model(`mcp_provider`, define.Postgres).
		Where(`id`, cast.ToString(providerId)).
		Delete()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}
