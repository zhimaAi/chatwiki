// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package business

import (
	"chatwiki/internal/app/plugin/define"
	"chatwiki/internal/app/plugin/php"
	"chatwiki/internal/pkg/lib_web"
	"errors"
	"net/http"
	"strings"

	"github.com/zhimaAi/go_tools/logs"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetLocalPluginList(c *gin.Context) {
	adminUserId := c.GetHeader(`admin_user_id`)
	if adminUserId == "" {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(`缺少admin_user_id`)))
		return
	}

	reqData := map[string]any{
		"filter_type": c.Query("filter_type"),
		"title":       c.Query("title"),
	}

	// 获取远程插件列表
	resp, err := requestXiaokefu(`kf/ChatWiki/CommonGetPluginList`, reqData)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	rawRemoteList, ok := resp.Data.([]interface{})
	if !ok {
		err = errors.New(`invalid remote data format`)
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	var (
		remoteList    []map[string]any
		remoteNameSet = make(map[string]map[string]any)
	)
	for _, item := range rawRemoteList {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
		}
		name := cast.ToString(itemMap[`name`])
		if name == "" {
			continue
		}
		if !cast.ToBool(itemMap[`enabled`]) {
			continue
		}
		remoteNameSet[name] = itemMap
		remoteList = append(remoteList, itemMap)
	}

	dbPluginList, err := msql.Model(`plugin_config`, define.Postgres).
		Where(`admin_user_id`, adminUserId).
		Select()
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	type PluginDetail struct {
		Local  *php.PluginManifest    `json:"local"`
		Remote map[string]interface{} `json:"remote"`
	}
	var result []PluginDetail

	for _, dbPlugin := range dbPluginList {
		manifest, err := php.GetPluginManifest(dbPlugin[`name`])
		if err != nil {
			logs.Error(err.Error())
			continue
		}
		remote, ok := remoteNameSet[manifest.Name]
		if !ok {
			continue
		}
		manifest.HasLoaded = cast.ToBool(dbPlugin[`has_loaded`])
		result = append(result, PluginDetail{
			Local:  manifest,
			Remote: remote,
		})
	}
	c.String(http.StatusOK, lib_web.FmtJson(result, nil))
}

func GetLocalPluginDetail(c *gin.Context) {
	adminUserId := c.GetHeader(`admin_user_id`)
	if adminUserId == "" {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(`缺少admin_user_id`)))
		return
	}
	name := strings.TrimSpace(c.Query("name"))
	if len(name) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(`缺少插件名称参数`)))
		return
	}

	// 获取远程插件详情
	resp, err := requestXiaokefu(`kf/ChatWiki/CommonGetPluginDetail`, map[string]any{
		`name`: name,
	})
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	remotePlugin, ok := resp.Data.(map[string]any)
	if !ok {
		err = errors.New(`invalid remote data format`)
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	type PluginDetail struct {
		Local  *php.PluginManifest    `json:"local"`
		Remote map[string]interface{} `json:"remote"`
	}
	result := PluginDetail{
		Remote: remotePlugin,
	}

	// 获取插件manifest
	manifest, err := php.GetPluginManifest(name)
	if err == nil && manifest != nil {
		dbPluginDetail, err := msql.Model(`plugin_config`, define.Postgres).
			Where(`admin_user_id`, adminUserId).
			Where(`name`, name).
			Find()
		if err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, err))
			return
		}
		if len(dbPluginDetail) > 0 {
			result.Local = manifest
			result.Local.HasLoaded = cast.ToBool(dbPluginDetail[`has_loaded`])
		}
	}

	c.String(http.StatusOK, lib_web.FmtJson(result, nil))
}

func DestroyLocalPlugin(c *gin.Context) {
	adminUserId := c.GetHeader(`admin_user_id`)
	if adminUserId == "" {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(`缺少admin_user_id`)))
		return
	}
	var req struct {
		Name string `form:"name" binding:"required"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	_, err := msql.Model(`plugin_config`, define.Postgres).
		Where(`admin_user_id`, adminUserId).
		Where(`name`, req.Name).
		Delete()
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	define.PhpPlugin.UnloadPhpPlugin(req.Name)
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func LoadLocalPlugin(c *gin.Context) {
	adminUserId := c.GetHeader(`admin_user_id`)
	if adminUserId == "" {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(`缺少admin_user_id`)))
		return
	}

	var req struct {
		Name string `form:"name" binding:"required"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	Name := strings.Split(req.Name, ",")

	for _, name := range Name {
		info, err := msql.Model(`plugin_config`, define.Postgres).
			Where(`admin_user_id`, adminUserId).
			Where(`name`, name).
			Find()
		if err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, err))
			return
		}
		if len(info) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New("no data")))
			return
		}

		// 加载插件
		err = define.PhpPlugin.LoadPhpPlugin(name, define.Version, define.DefaultPhpEnv)
		if err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, err))
			return
		}

		// 更新数据库记录
		_, err = msql.Model(`plugin_config`, define.Postgres).
			Where(`admin_user_id`, adminUserId).
			Where(`name`, name).
			Update(msql.Datas{
				`update_time`: tool.Time2Int(),
				`has_loaded`:  true,
			})
		if err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, err))
			return
		}
	}

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func UnloadLocalPlugin(c *gin.Context) {
	adminUserId := c.GetHeader(`admin_user_id`)
	if adminUserId == "" {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(`缺少admin_user_id`)))
		return
	}

	var req struct {
		Name string `form:"name" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	_, err := msql.Model(`plugin_config`, define.Postgres).
		Where(`admin_user_id`, adminUserId).
		Where(`name`, req.Name).
		Update(msql.Datas{
			`update_time`: tool.Time2Int(),
			`has_loaded`:  false,
		})
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func GetLocalPluginConfig(c *gin.Context) {
	adminUserId := c.GetHeader(`admin_user_id`)
	if adminUserId == "" {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(`缺少admin_user_id`)))
		return
	}

	var req struct {
		Name string `form:"name" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	info, err := msql.Model(`plugin_config`, define.Postgres).
		Where(`admin_user_id`, adminUserId).
		Where(`name`, req.Name).
		Find()
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(info[`data`], nil))
}

func UpdateLocalPluginConfig(c *gin.Context) {
	adminUserId := c.GetHeader(`admin_user_id`)
	if adminUserId == "" {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(`缺少admin_user_id`)))
		return
	}

	var req struct {
		Name string `form:"name" binding:"required"`
		Data string `form:"data" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	_, err := msql.Model(`plugin_config`, define.Postgres).
		Where(`admin_user_id`, adminUserId).
		Where(`name`, req.Name).
		Update(msql.Datas{
			`update_time`: tool.Time2Int(),
			`data`:        req.Data,
		})
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func RunLocalPluginLambda(c *gin.Context) {
	adminUserId := c.GetHeader(`admin_user_id`)
	if adminUserId == "" {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(`缺少admin_user_id`)))
		return
	}

	var req struct {
		Name   string         `form:"name" binding:"required"`
		Action string         `form:"action" binding:"required"`
		Params map[string]any `form:"params"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	if req.Params == nil {
		req.Params = make(map[string]any)
	}

	manifest, err := php.GetPluginManifest(req.Name)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	info, err := msql.Model(`plugin_config`, define.Postgres).
		Where(`admin_user_id`, adminUserId).
		Where(`name`, req.Name).
		Where(`has_loaded`, `true`).
		Find()
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	if len(info) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(`插件未加载`)))
		return
	}

	resp, err := define.PhpPlugin.ExecLambdaPhpPlugin(req.Name, manifest.Title, req.Action, req.Params)
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}

	c.JSON(http.StatusOK, resp)
}
