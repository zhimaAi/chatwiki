// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package manage

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/lib_web"
	"chatwiki/internal/pkg/wechat"
	"chatwiki/internal/pkg/wechat/feishu_robot"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

// FeishuUserAuthLoginRedirect 跳转到飞书授权地址
func FeishuUserAuthLoginRedirect(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	feishuAppId := strings.TrimSpace(c.Query(`feishu_app_id`))
	feishuAppSecret := strings.TrimSpace(c.Query(`feishu_app_secret`))
	feishuFrontendRedirectUrl := strings.TrimSpace(c.Query(`feishu_frontend_auth_redirect_url`))
	if len(feishuAppId) == 0 || len(feishuAppSecret) == 0 || len(feishuFrontendRedirectUrl) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	app, err := wechat.GetApplication(msql.Params{`app_type`: lib_define.FeiShuRobot, `app_id`: feishuAppId, `app_secret`: feishuAppSecret})
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	feishuApp, ok := app.(wechat.FeishuInterface)
	if !ok {
		logs.Error(`飞书app初始化失败`)
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	backendRedirectUrl := define.Config.WebService[`admin_domain`] + "/manage/feishuUserAuthLogin/callback"
	callbackUrl, err := feishuApp.BuildUserAuthLoginUrl(backendRedirectUrl, feishuFrontendRedirectUrl)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.Redirect(http.StatusFound, callbackUrl)
}

// FeishuUserAuthLoginCallback 飞书授权回调
func FeishuUserAuthLoginCallback(c *gin.Context) {
	stateStr := strings.TrimSpace(c.Query("state"))
	code := strings.TrimSpace(c.Query(`code`))
	stateStateJson, err := tool.Base64Decode(stateStr)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	var state feishu_robot.FeishuUserAuthLoginState
	if err := json.Unmarshal([]byte(stateStateJson), &state); err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	if len(state.FrontRedirectUrl) == 0 || len(state.AppId) == 0 || len(state.AppSecret) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}
	if len(code) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}

	app, err := wechat.GetApplication(msql.Params{`app_type`: lib_define.FeiShuRobot, `app_id`: state.AppId, `app_secret`: state.AppSecret})
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	feishuApp, ok := app.(wechat.FeishuInterface)
	if !ok {
		logs.Error(`飞书app初始化失败`)
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	backendRedirectUrl := define.Config.WebService[`admin_domain`] + "/manage/feishuUserAuthLogin/callback"
	resp, err := feishuApp.GetUserAccessToken(code, backendRedirectUrl)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	params := url.Values{}
	params.Add("user_access_token", resp.AccessToken)
	c.Redirect(http.StatusFound, tool.Base64DecodeNoError(state.FrontRedirectUrl)+"&"+params.Encode())
}

// GetFeishuDocFileList 获取飞书知识库文档列表
func GetFeishuDocFileList(c *gin.Context) {
	var adminUserId int
	if adminUserId = GetAdminUserId(c); adminUserId == 0 {
		return
	}
	feishuAppId := strings.TrimSpace(c.PostForm(`feishu_app_id`))
	feishuAppSecret := strings.TrimSpace(c.PostForm(`feishu_app_secret`))
	userAccessToken := strings.TrimSpace(c.PostForm(`user_access_token`))
	if len(feishuAppId) == 0 || len(feishuAppSecret) == 0 || len(userAccessToken) == 0 {
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `param_lack`))))
		return
	}

	app, err := wechat.GetApplication(msql.Params{`app_type`: lib_define.FeiShuRobot, `app_id`: feishuAppId, `app_secret`: feishuAppSecret})
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}

	feishuApp, ok := app.(wechat.FeishuInterface)
	if !ok {
		logs.Error(`飞书app初始化失败`)
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	result, err := feishuApp.GetDocFileTree(userAccessToken, "")
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_web.FmtJson(nil, errors.New(i18n.Show(common.GetLang(c), `sys_err`))))
		return
	}
	c.String(http.StatusOK, lib_web.FmtJson(result, nil))
}
