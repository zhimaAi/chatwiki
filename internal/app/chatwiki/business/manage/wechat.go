// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_redis"
	"chatwiki/internal/pkg/lib_web"
	"chatwiki/internal/pkg/wechat"
	"chatwiki/internal/pkg/wechat/messenger"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	cams "github.com/alibabacloud-go/cams-20200606/v5/client"
	openapiutil "github.com/alibabacloud-go/darabonba-openapi/v2/utils"
	"github.com/alibabacloud-go/tea/dara"
	"github.com/aliyun/credentials-go/credentials"
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
	m := msql.Model(define.TableChatAiWechatApp, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId))
	robotId := cast.ToInt(c.Query(`robot_id`))
	if robotId > 0 {
		m.Where(`robot_id`, cast.ToString(robotId))
	}
	appType := strings.TrimSpace(c.Query(`app_type`))
	if len(appType) > 0 {
		m.Where(`app_type`, appType)
		//For official accounts, sort by the sort field
		if appType == lib_define.AppOfficeAccount {
			m.Order(`sort asc`)
		}
	}
	appName := strings.TrimSpace(c.Query(`app_name`))
	if len(appName) > 0 {
		m.Where(`app_name`, `like`, appName)
	}
	queryAll := cast.ToBool(c.Query(`is_all`))

	list, err := m.Order(`id desc`).Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	lang := common.GetLang(c)
	for i, appInfo := range list {

		if appInfo[`app_type`] == lib_define.AppWhatsapp {
			// WhatsApp: expose stored values under the same field names the save form uses, so the
			// frontend can prefill the edit form straight from the list without a detail request.
			// Sensitive values (access_key_secret) are returned in plaintext for now (no encryption yet).
			list[i][`access_key_id`] = appInfo[`cams_access_key_id`]
			list[i][`cust_space_id`] = appInfo[`cams_cust_space_id`]
			list[i][`access_key_secret`] = appInfo[`app_secret`]
		}

		accountIsVerify := lib_define.WechatAccountIsVerify(appInfo[`account_customer_type`])
		list[i][`account_is_verify`] = cast.ToString(accountIsVerify)
		list[i][`account_corp_verify`] = cast.ToString(lib_define.WechatAccountIsCorpVerify(appInfo[`account_customer_type`]))

		if appInfo[`app_type`] == lib_define.AppOfficeAccount && !accountIsVerify {
			list[i][`wechat_reply_type`] = i18n.Show(lang, `wechat_no_answer_manual_get`, define.Config.WebService[`wechat_wait`])
		} else {
			list[i][`wechat_reply_type`] = i18n.Show(lang, `system_auto_reply`)
		}
		//account_type_name
		if appInfo[`account_type`] != `-1` && appInfo[`account_type`] != `0` {
			list[i][`account_type_name`] = i18n.Show(lang, fmt.Sprintf(`account_type_%s`, appInfo[`account_type`]))
		} else {
			list[i][`account_type_name`] = ``
		}
		//account_customer_type_name
		if appInfo[`account_customer_type`] != `-1` && appInfo[`account_customer_type`] != `0` {
			list[i][`account_customer_type_name`] = i18n.Show(lang, fmt.Sprintf(`account_customer_type_%s`, appInfo[`account_customer_type`]))
		} else {
			list[i][`account_customer_type_name`] = ``
		}
		if (appInfo[`account_customer_type`] == `-1` || appInfo[`account_type`] == `-1`) && tool.InArrayString(appInfo[`app_type`], []string{lib_define.AppOfficeAccount, lib_define.AppMini}) {
			go func() {
				_ = common.RefreshAccountVerify(appInfo) //refresh defaults asynchronously
			}()
		}
		if appInfo[`app_type`] == lib_define.AppOfficeAccount && !queryAll && !accountIsVerify {
			list[i] = nil
			continue //filter out unverified wechat official account
		}
	}
	list = filterNilByMSQLParams(list)
	c.String(http.StatusOK, lib_web.FmtJson(list, nil))
}

func filterNilByMSQLParams(list []msql.Params) []msql.Params {
	var newList = make([]msql.Params, 0)
	for _, item := range list {
		if item != nil {
			newList = append(newList, item)
		}
	}
	return newList
}

// RobotRelateOfficialAccount binds a robot to an official account
func RobotRelateOfficialAccount(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	//Get parameters
	robotId := cast.ToInt(c.PostForm(`robot_id`))
	appIdList := strings.TrimSpace(c.PostForm(`app_id_list`))

	//Get robot info
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

	//Get configs already linked to this robot
	oldAppInfoList, err := msql.Model(`chat_ai_wechat_app`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`robot_id`, cast.ToString(robotId)).
		Where(`app_type`, lib_define.AppOfficeAccount).
		Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	var shouldDelIdList []string
	for _, appInfo := range oldAppInfoList {
		shouldDelIdList = append(shouldDelIdList, appInfo[`id`])
	}

	//Get configs to add
	shouldAddAppInfoList, err := msql.Model(`chat_ai_wechat_app`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(userId)).
		Where(`app_id`, `in`, appIdList).
		Select()
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	var shouldAddIdList []string
	for _, appInfo := range shouldAddAppInfoList {
		if cast.ToInt(appInfo[`robot_id`]) != robotId {
			old, err := msql.Model(`chat_ai_robot`, define.Postgres).
				Where(`id`, cast.ToString(appInfo[`robot_id`])).
				Where(`admin_user_id`, cast.ToString(userId)).
				Field(`id,robot_name`).
				Find()
			if err != nil {
				logs.Error(err.Error())
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
				return
			}
			if len(old) > 0 {
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `official_has_bind`, appInfo[`app_name`], old[`robot_name`]))))
				return
			}
		}
		shouldAddIdList = append(shouldAddIdList, appInfo[`id`])
	}

	//Clear old records first
	if len(shouldDelIdList) > 0 {
		_, err = msql.Model(`chat_ai_wechat_app`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(userId)).
			Where(`id`, `in`, strings.Join(shouldDelIdList, ",")).
			Where(`app_type`, lib_define.AppOfficeAccount).
			Update(msql.Datas{`robot_id`: 0, `robot_key`: ``})
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
	}

	//Then update records
	if len(shouldAddIdList) > 0 {
		_, err = msql.Model(`chat_ai_wechat_app`, define.Postgres).
			Where(`id`, `in`, strings.Join(shouldAddIdList, `,`)).
			Update(msql.Datas{`robot_id`: robotId, `robot_key`: robot[`robot_key`]})
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
	}

	//Clear cache
	for _, appInfo := range append(oldAppInfoList, shouldAddAppInfoList...) {
		lib_redis.DelCacheData(define.Redis, &common.WechatAppCacheBuildHandler{Field: `id`, Value: appInfo[`id`]})
		lib_redis.DelCacheData(define.Redis, &common.WechatAppCacheBuildHandler{Field: `app_id`, Value: appInfo[`app_id`]})
		lib_redis.DelCacheData(define.Redis, &common.WechatAppCacheBuildHandler{Field: `access_key`, Value: appInfo[`access_key`]})
	}

	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
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
	appSecret := strings.TrimSpace(c.DefaultPostForm(`app_secret`, `no-need`))
	appAvatar := ``
	appType := strings.TrimSpace(c.PostForm(`app_type`))
	if appType == lib_define.AppWhatsapp {
		appId = normalizeWhatsappAppID(appId)
	}
	//unchangeable
	var accessKey string
	var oldAppId string
	var oldAccessKey string
	if id > 0 {
		requestAppId := appId
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
		oldAccessKey = appInfo[`access_key`]
		robotId = cast.ToInt(appInfo[`robot_id`])
		appType = appInfo[`app_type`]
		oldAppId = appInfo[`app_id`]
		appId = resolveWechatAppSaveAppID(appType, oldAppId, requestAppId)
		if appType == lib_define.AppWhatsapp {
			appId = normalizeWhatsappAppID(appId)
		}
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
	// WhatsApp reuses app_id (= phone number) but supplies its credentials via
	// access_key_secret/access_key_id/cust_space_id, validated inside the AppWhatsapp branch.
	// app_secret is never sent for WhatsApp (it defaults to "no-need"), so skip the generic
	// app_secret requirement here; app_id is still required and checked in the branch.
	if appType == lib_define.AppWhatsapp {
		if id < 0 || len(appName) == 0 || len(appType) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
			return
		}
	} else if id < 0 || len(appName) == 0 || len(appId) == 0 || len(appSecret) == 0 || len(appType) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if !tool.InArrayString(appType, lib_define.AppTypeList) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_invalid`, `app_type`))))
		return
	}
	//official account access quantity limit
	if appType == lib_define.AppOfficeAccount && id == 0 {
		count, err := msql.Model(`chat_ai_wechat_app`, define.Postgres).Where(`admin_user_id`, cast.ToString(userId)).Where(`app_type`, appType).Count()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if count >= 10 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `max_official_accounts_limit`))))
			return
		}
	}

	//WhatsApp channel: read its own fields, run connectivity check, skip generic token verification
	if appType == lib_define.AppWhatsapp {
		accessKeySecret := strings.TrimSpace(c.PostForm(`access_key_secret`))
		accessKeyId := strings.TrimSpace(c.PostForm(`access_key_id`))
		custSpaceId := strings.TrimSpace(c.PostForm(`cust_space_id`))
		// Same edit model as other channels: the frontend prefills the form from
		// getWechatAppList (incl. the plaintext secret) and re-submits every field, so we
		// simply overwrite with what was posted — no blank-keep special-casing.
		if len(appId) == 0 || len(accessKeySecret) == 0 || len(accessKeyId) == 0 || len(custSpaceId) == 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
			return
		}
		// Connectivity verification (appId is the phone number)
		if verifyErr := verifyWhatsappConnectivity(common.GetLang(c), accessKeyId, accessKeySecret, custSpaceId, appId); verifyErr != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, verifyErr))
			return
		}
		// Duplicate phone check for edit scenario (new records are checked above)
		if id > 0 && appId != oldAppId {
			if existing, err := common.GetWechatAppInfo(`app_id`, appId); err != nil {
				logs.Error(err.Error())
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
				return
			} else if len(existing) > 0 && cast.ToInt(existing[`id`]) != id {
				c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `app_exist`))))
				return
			}
		}
		// app_secret column reuses the WhatsApp access_key_secret (shared schema)
		appSecret = accessKeySecret

		//app_avatar upload
		fileHeader, _ := c.FormFile(`app_avatar`)
		uploadInfo, err := common.SaveUploadedFile(fileHeader, define.ImageLimitSize, userId, `app_avatar`, define.ImageAllowExt)
		if err == nil && uploadInfo != nil {
			appAvatar = uploadInfo.Link
		}
		// Derive a channel-level callback key from CustSpaceId. All phone numbers in
		// the same channel share one push_url, then inbound callbacks are routed by To.
		accessKey = deriveWhatsappAccessKey(custSpaceId)
		//database dispose
		data := msql.Datas{
			`app_name`:           appName,
			`app_secret`:         appSecret,
			`access_key`:         accessKey,
			`cams_access_key_id`: accessKeyId,
			`cams_cust_space_id`: custSpaceId,
			`update_time`:        tool.Time2Int(),
		}
		//WhatsApp must be bound to a robot
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
		data[`robot_key`] = robot[`robot_key`]
		if len(appAvatar) > 0 {
			data[`app_avatar`] = appAvatar
		}
		m := msql.Model(`chat_ai_wechat_app`, define.Postgres)
		if id > 0 {
			if appId != oldAppId {
				data[`app_id`] = appId
			}
			// access_key follows cust_space_id, so changing the channel updates push_url.
			_, err = m.Where(`id`, cast.ToString(id)).Update(data)
		} else {
			data[`admin_user_id`] = userId
			data[`robot_id`] = robotId
			data[`app_id`] = appId
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
		if len(oldAppId) > 0 && oldAppId != appId {
			lib_redis.DelCacheData(define.Redis, &common.WechatAppCacheBuildHandler{Field: `app_id`, Value: oldAppId})
		}
		lib_redis.DelCacheData(define.Redis, &common.WechatAppCacheBuildHandler{Field: `access_key`, Value: accessKey})
		if len(oldAccessKey) > 0 && oldAccessKey != accessKey {
			lib_redis.DelCacheData(define.Redis, &common.WechatAppCacheBuildHandler{Field: `access_key`, Value: oldAccessKey})
		}
		//configure external service parameters
		// Query by app_id here because access_key is shared by all phone numbers in a channel.
		appInfo, err := common.GetWechatAppInfo(`app_id`, appId)
		if err == nil {
			// WhatsApp callback authenticates purely via the access_key in push_url; it does
			// not use the WeChat-style push_token/push_aeskey, so we don't return them here.
			appInfo[`push_url`] = fmt.Sprintf(`%s/push_pwd/%s/%s`, define.Config.WebService[`push_domain`], appInfo[`app_type`], appInfo[`access_key`])
		}
		c.String(http.StatusOK, lib_web.FmtJson(appInfo, err))
		return
	}

	//get token verification
	app, err := wechat.GetApplication(msql.Params{`app_type`: appType, `app_id`: appId, `app_secret`: appSecret})
	if err != nil {
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	if _, _, err := app.GetToken(false); err != nil {
		if appType == lib_define.AppMessenger {
			if errors.Is(err, messenger.ErrSecretInvalid) {
				err = errors.New(i18n.Show(common.GetLang(c), `messenger_secret_error`))
			} else if errors.Is(err, messenger.ErrPageIDMismatch) {
				err = errors.New(i18n.Show(common.GetLang(c), `messenger_page_id_mismatch`))
			}
		}
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

	//Official accounts do not require strong binding to a robot; other types must be linked directly
	if appType != lib_define.AppOfficeAccount {
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
		data[`robot_key`] = robot[`robot_key`]
	}
	//For official accounts, app_id should be globally unique
	if appType == lib_define.AppOfficeAccount && id == 0 {
		count, err := msql.Model(`chat_ai_wechat_app`, define.Postgres).Where(`app_id`, appId).Where(`app_type`, appType).Count()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if count > 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `app_id_exist`))))
			return
		}
	}
	//For Messenger, app_id should be globally unique
	if appType == lib_define.AppMessenger && id == 0 {
		count, err := msql.Model(`chat_ai_wechat_app`, define.Postgres).Where(`app_id`, appId).Where(`app_type`, appType).Count()
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		if count > 0 {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `app_id_exist`))))
			return
		}
	}
	if appType == lib_define.FeiShuRobot {
		data[`encrypt_key`] = strings.TrimSpace(c.PostForm(`encrypt_key`))
		data[`verification_token`] = strings.TrimSpace(c.PostForm(`verification_token`))
	}
	if id > 0 && (appType == lib_define.AppWecomRobot || appType == lib_define.AppMessenger) && appId != oldAppId {
		if appInfo, err := common.GetWechatAppInfo(`app_id`, appId); err != nil {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		} else if len(appInfo) > 0 && cast.ToInt(appInfo[`id`]) != id {
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `app_exist`))))
			return
		}
		data[`app_id`] = appId
	}

	if len(appAvatar) > 0 {
		data[`app_avatar`] = appAvatar
	}
	//WeChat app verification type
	if basic, _, err := app.GetAccountBasicInfo(); err == nil {
		data[`account_customer_type`] = basic.CustomerType
		data[`account_type`] = basic.AccountType
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
	if len(oldAppId) > 0 && oldAppId != appId {
		lib_redis.DelCacheData(define.Redis, &common.WechatAppCacheBuildHandler{Field: `app_id`, Value: oldAppId})
	}
	lib_redis.DelCacheData(define.Redis, &common.WechatAppCacheBuildHandler{Field: `access_key`, Value: accessKey})
	//configure external service parameters
	appInfo, err := common.GetWechatAppInfo(`access_key`, accessKey)
	if err == nil {
		appInfo[`push_url`] = fmt.Sprintf(`%s/push_pwd/%s`, define.Config.WebService[`push_domain`], appInfo[`app_type`])
		if !tool.InArrayString(appInfo[`app_type`], []string{lib_define.AppWechatKefu, lib_define.AppWecomRobot}) {
			appInfo[`push_url`] += fmt.Sprintf(`/%s`, appInfo[`access_key`])
		}
		appInfo[`push_token`] = lib_define.SignToken
		if appInfo[`app_type`] == lib_define.AppMessenger {
			appInfo[`push_token`] = appInfo[`access_key`]
		}
		appInfo[`push_aeskey`] = lib_define.AesKey
		appInfo[`account_is_verify`] = cast.ToString(lib_define.WechatAccountIsVerify(appInfo[`account_customer_type`]))
	}
	c.String(http.StatusOK, lib_web.FmtJson(appInfo, err))
}

func resolveWechatAppSaveAppID(appType, oldAppId, requestAppId string) string {
	if appType == lib_define.AppWecomRobot || appType == lib_define.AppMessenger || appType == lib_define.AppWhatsapp {
		return requestAppId
	}
	return oldAppId
}

func normalizeWhatsappAppID(appId string) string {
	return strings.TrimPrefix(strings.TrimSpace(appId), `+`)
}

// deriveWhatsappAccessKey deterministically maps one CustSpaceId to one callback key.
// The prefix avoids collisions with legacy random keys from other channels.
func deriveWhatsappAccessKey(custSpaceId string) string {
	return `wa` + tool.MD5(`whatsapp:` + custSpaceId)[:18]
}

// verifyWhatsappConnectivity calls QueryChatappPhoneNumbers to confirm the AK is valid
// and that phoneNumber exists in the CustSpace's registered phone list.
func verifyWhatsappConnectivity(lang, accessKeyId, accessKeySecret, custSpaceId, phoneNumber string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
		}
	}()
	credential, credErr := credentials.NewCredential(&credentials.Config{
		Type:            dara.String(`access_key`),
		AccessKeyId:     dara.String(accessKeyId),
		AccessKeySecret: dara.String(accessKeySecret),
	})
	if credErr != nil {
		return credErr
	}
	client, clientErr := cams.NewClient(&openapiutil.Config{
		Credential: credential,
		Endpoint:   dara.String(`cams.ap-southeast-1.aliyuncs.com`),
	})
	if clientErr != nil {
		return clientErr
	}
	req := &cams.QueryChatappPhoneNumbersRequest{
		CustSpaceId: dara.String(custSpaceId),
	}
	resp, callErr := client.QueryChatappPhoneNumbersWithOptions(req, &dara.RuntimeOptions{})
	if callErr != nil {
		var sdkErr *dara.SDKError
		if errors.As(callErr, &sdkErr) {
			return errors.New(i18n.Show(lang, `whatsapp_verify_failed`, *sdkErr.Message))
		}
		return callErr
	}
	if resp == nil || resp.Body == nil {
		return errors.New(i18n.Show(lang, `whatsapp_verify_empty_resp`))
	}
	body := resp.Body
	code := ``
	if body.Code != nil {
		code = *body.Code
	}
	if code != `OK` {
		msg := ``
		if body.Message != nil {
			msg = *body.Message
		}
		return errors.New(i18n.Show(lang, `whatsapp_verify_failed`, fmt.Sprintf(`%s %s`, code, msg)))
	}
	// Normalize phone: strip leading '+' for comparison since Aliyun stores without '+'
	normalizedPhone := strings.TrimPrefix(phoneNumber, `+`)
	for _, p := range body.PhoneNumbers {
		if p != nil && p.PhoneNumber != nil && *p.PhoneNumber == normalizedPhone {
			return nil
		}
	}
	return errors.New(i18n.Show(lang, `whatsapp_phone_not_registered`))
}

func SortWechatApp(c *gin.Context) {
	var userId int
	if userId = GetAdminUserId(c); userId == 0 {
		return
	}
	idList := strings.Split(c.PostForm(`id_list`), ",")
	if len(idList) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	sort := 0
	for _, id := range idList {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		_, err := msql.Model(`chat_ai_wechat_app`, define.Postgres).Where(`admin_user_id`, cast.ToString(userId)).Where(`id`, id).Update(msql.Datas{`sort`: sort})
		if err != nil {
			logs.Error(err.Error())
			c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
			return
		}
		sort++
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
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
	if appInfo[`app_type`] == lib_define.AppWhatsapp {
		appInfo[`access_key_id`] = appInfo[`cams_access_key_id`]
		appInfo[`cust_space_id`] = appInfo[`cams_cust_space_id`]
		appInfo[`access_key_secret`] = appInfo[`app_secret`]
		appInfo[`push_url`] = fmt.Sprintf(`%s/push_pwd/%s/%s`, define.Config.WebService[`push_domain`], appInfo[`app_type`], appInfo[`access_key`])
		c.String(http.StatusOK, lib_web.FmtJson(appInfo, nil))
		return
	}
	//configure external service parameters
	appInfo[`push_url`] = fmt.Sprintf(`%s/push_pwd/%s`, define.Config.WebService[`push_domain`], appInfo[`app_type`])
	if !tool.InArrayString(appInfo[`app_type`], []string{lib_define.AppWechatKefu, lib_define.AppWecomRobot}) {
		appInfo[`push_url`] += fmt.Sprintf(`/%s`, appInfo[`access_key`])
	}
	appInfo[`push_token`] = lib_define.SignToken
	if appInfo[`app_type`] == lib_define.AppMessenger {
		appInfo[`push_token`] = appInfo[`access_key`]
	}
	appInfo[`push_aeskey`] = lib_define.AesKey
	appInfo[`account_is_verify`] = cast.ToString(lib_define.WechatAccountIsVerify(appInfo[`account_customer_type`]))
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

func RefreshAccountVerify(c *gin.Context) {
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
	err = common.RefreshAccountVerify(appInfo)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, err))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(nil, nil))
}

func SetWechatNotVerifyConfig(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	//get params
	id := cast.ToInt64(c.PostForm(`id`))
	wechatNotVerifyHandGetReply := strings.TrimSpace(c.PostForm(`wechat_not_verify_hand_get_reply`))
	wechatNotVerifyHandGetWord := strings.TrimSpace(c.PostForm(`wechat_not_verify_hand_get_word`))
	wechatNotVerifyHandGetNext := strings.TrimSpace(c.PostForm(`wechat_not_verify_hand_get_next`))
	//check required
	if id <= 0 || len(wechatNotVerifyHandGetReply) == 0 || len(wechatNotVerifyHandGetWord) == 0 || len(wechatNotVerifyHandGetNext) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	//data check
	m := msql.Model(`chat_ai_robot`, define.Postgres)
	robotKey, err := m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(adminUserId)).Value(`robot_key`)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(robotKey) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	//database dispose
	data := msql.Datas{
		`wechat_not_verify_hand_get_reply`: common.MbSubstr(wechatNotVerifyHandGetReply, 0, 100),
		`wechat_not_verify_hand_get_word`:  common.MbSubstr(wechatNotVerifyHandGetWord, 0, 100),
		`wechat_not_verify_hand_get_next`:  common.MbSubstr(wechatNotVerifyHandGetNext, 0, 100),
		`update_time`:                      tool.Time2Int(),
	}
	if _, err = m.Where(`id`, cast.ToString(id)).Update(data); err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.RobotCacheBuildHandler{RobotKey: robotKey})
	c.String(http.StatusOK, lib_web.FmtJson(common.GetRobotInfo(robotKey)))
}

func SetWechatConfigSwitch(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	//get params
	id := cast.ToInt64(c.PostForm(`id`))
	key := strings.TrimSpace(c.PostForm(`key`))
	val := cast.ToInt(cast.ToBool(c.PostForm(`val`)))
	//check required
	if id <= 0 || !tool.InArrayString(key, []string{`show_ai_msg_gzh`, `show_typing_gzh`, `show_ai_msg_mini`, `show_typing_mini`}) {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	//data check
	m := msql.Model(`chat_ai_robot`, define.Postgres)
	robotKey, err := m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(adminUserId)).Value(`robot_key`)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(robotKey) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `no_data`))))
		return
	}
	//database dispose
	data := msql.Datas{key: val, `update_time`: tool.Time2Int()}
	if _, err = m.Where(`id`, cast.ToString(id)).Update(data); err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	//clear cached data
	lib_redis.DelCacheData(define.Redis, &common.RobotCacheBuildHandler{RobotKey: robotKey})
	c.String(http.StatusOK, lib_web.FmtJson(common.GetRobotInfo(robotKey)))
}
