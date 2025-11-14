// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/lib_web"
	"chatwiki/internal/pkg/wechat"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func GetWechatAppList(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	m := msql.Model(`chat_ai_wechat_app`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId))
	robotId := cast.ToInt(c.Query(`robot_id`))
	if robotId > 0 {
		m.Where(`robot_id`, cast.ToString(robotId))
	}
	appType := strings.TrimSpace(c.Query(`app_type`))
	if len(appType) > 0 {
		m.Where(`app_type`, appType)
	}
	appName := strings.TrimSpace(c.Query(`app_name`))
	if len(appName) > 0 {
		m.Where(`app_name`, `like`, appName)
	}
	list, err := m.Order(`id desc`).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
}

func SaveWechatApp(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//get params
	id := cast.ToInt(c.PostForm(`id`))
	robotId := cast.ToInt(c.PostForm(`robot_id`))
	appName := strings.TrimSpace(c.PostForm(`app_name`))
	appId := strings.TrimSpace(c.PostForm(`app_id`))
	appSecret := strings.TrimSpace(c.PostForm(`app_secret`))
	appAvatar := ``
	appType := strings.TrimSpace(c.PostForm(`app_type`))
	//unchangeable
	var accessKey string
	if id > 0 {
		appInfo, err := common.GetWechatAppInfo(`id`, cast.ToString(id))
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if len(appInfo) == 0 || cast.ToInt(appInfo[`admin_user_id`]) != userId {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
			return
		}
		accessKey = appInfo[`access_key`]
		robotId = cast.ToInt(appInfo[`robot_id`])
		appId = appInfo[`app_id`]
		appType = appInfo[`app_type`]
	} else {
		if appInfo, err := common.GetWechatAppInfo(`app_id`, appId); err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		} else if len(appInfo) > 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `app_exist`))))
			return
		}
		appAvatar = define.LocalUploadPrefix + fmt.Sprintf(`default/%s_avatar.png`, appType)
	}
	//check required
	if id < 0 || robotId <= 0 || len(appName) == 0 || len(appId) == 0 || len(appSecret) == 0 || len(appType) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if !tool.InArrayString(appType, lib_define.AppTypeList) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `app_type`))))
		return
	}
	//data check
	robot, err := msql.Model(`chat_ai_robot`, define.Postgres).Where(`id`, cast.ToString(robotId)).
		Where(`admin_user_id`, cast.ToString(userId)).Field(`id,robot_key`).Find()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(robot) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	//get token verification
	app, err := wechat.GetApplication(msql.Params{`app_type`: appType, `app_id`: appId, `app_secret`: appSecret})
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	if _, _, err := app.GetToken(false); err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	//app_avatar upload
	fileHeader, _ := c.FormFile(`app_avatar`)
	uploadInfo, err := common.SaveUploadedFile(fileHeader, define.ImageLimitSize, userId, `app_avatar`, define.ImageAllowExt)
	if err == nil && uploadInfo != nil {
		appAvatar = uploadInfo.Link
	}
	//database dispose
	data := msql.Datas{
		`app_name`:    appName,
		`app_secret`:  appSecret,
		`update_time`: tool.Time2Int(),
	}
	if len(appAvatar) > 0 {
		data[`app_avatar`] = appAvatar
	}
	m := msql.Model(`chat_ai_wechat_app`, define.Postgres)
	if id > 0 {
		_, err = m.Where(`id`, cast.ToString(id)).Update(data)
	} else {
		for i := 0; i < 5; i++ {
			tempKey := tool.Random(10)
			if appInfo, e := common.GetWechatAppInfo(`access_key`, tempKey); e == nil && len(appInfo) == 0 {
				accessKey = tempKey
				break
			}
			time.Sleep(time.Nanosecond) //sleep 1 ns
		}
		if len(accessKey) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		data[`admin_user_id`] = userId
		data[`robot_id`] = robotId
		data[`robot_key`] = robot[`robot_key`]
		data[`app_id`] = appId
		data[`access_key`] = accessKey
		data[`app_type`] = appType
		data[`set_type`] = lib_define.PwdSetType
		data[`create_time`] = data[`update_time`]
		_, err = m.Insert(data)
	}
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.WechatAppCacheBuildHandler{Field: `id`, Value: cast.ToString(id)})
	lib_redis.DelCacheData(define.Redis, &common.WechatAppCacheBuildHandler{Field: `app_id`, Value: appId})
	lib_redis.DelCacheData(define.Redis, &common.WechatAppCacheBuildHandler{Field: `access_key`, Value: accessKey})
	//configure external service parameters
	appInfo, err := common.GetWechatAppInfo(`access_key`, accessKey)
	if err == nil {
		appInfo[`push_url`] = fmt.Sprintf(`%s/push_pwd/%s`, define.Config.WebService[`push_domain`], appInfo[`app_type`])
		if appInfo[`app_type`] != lib_define.AppWechatKefu {
			appInfo[`push_url`] += fmt.Sprintf(`/%s`, appInfo[`access_key`])
		}
		appInfo[`push_token`] = lib_define.SignToken
		appInfo[`push_aeskey`] = lib_define.AesKey
	}
	c.String(http.StatusOK, lib_web.FmtJson(appInfo, err))
}

func GetWechatAppInfo(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.Query(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	appInfo, err := common.GetWechatAppInfo(`id`, cast.ToString(id))
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(appInfo) == 0 || cast.ToInt(appInfo[`admin_user_id`]) != userId {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	//configure external service parameters
	appInfo[`push_url`] = fmt.Sprintf(`%s/push_pwd/%s`, define.Config.WebService[`push_domain`], appInfo[`app_type`])
	if appInfo[`app_type`] != lib_define.AppWechatKefu {
		appInfo[`push_url`] += fmt.Sprintf(`/%s`, appInfo[`access_key`])
	}
	appInfo[`push_token`] = lib_define.SignToken
	appInfo[`push_aeskey`] = lib_define.AesKey
	c.String(http.StatusOK, lib_web.FmtJson(appInfo, nil))
}

func DeleteWechatApp(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	id := cast.ToInt(c.PostForm(`id`))
	if id <= 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	appInfo, err := common.GetWechatAppInfo(`id`, cast.ToString(id))
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(appInfo) == 0 || cast.ToInt(appInfo[`admin_user_id`]) != userId {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	m := msql.Model(`chat_ai_wechat_app`, define.Postgres)
	_, err = m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(userId)).Delete()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.WechatAppCacheBuildHandler{Field: `id`, Value: appInfo[`id`]})
	lib_redis.DelCacheData(define.Redis, &common.WechatAppCacheBuildHandler{Field: `app_id`, Value: appInfo[`app_id`]})
	lib_redis.DelCacheData(define.Redis, &common.WechatAppCacheBuildHandler{Field: `access_key`, Value: appInfo[`access_key`]})
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}
