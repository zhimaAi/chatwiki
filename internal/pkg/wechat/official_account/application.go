// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package official_account

import (
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/wechat/common"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/messages"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/power"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/response"
	response2 "github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/response"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount"
	menuRequest "github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount/menu/request"
	publishRequest "github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount/publish/request"
	publishresponse "github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount/publish/response"
	openresponse "github.com/ArtisanCloud/PowerWeChat/v3/src/openPlatform/authorizer/miniProgram/account/response"
)

type Application struct {
	AppID  string
	Secret string
}

func (a *Application) SendSmartMenu(customer string, smartMenu lib_define.SmartMenu, push *lib_define.PushMessage) (int, error) {
	content := common.WechatFormatSmartMenu2C(smartMenu)
	return a.SendText(customer, content, push)
}

func (a *Application) SendImageTextLink(customer, url, title, description, localThumbURL, picurl string, push *lib_define.PushMessage) (int, error) {
	jsonStr, err := messages.NewLink(&power.HashMap{
		`url`:         url,
		`title`:       title,
		`description`: description,
	}).TransformForJsonRequest(gzhBuildSendMsgAppends(customer, push), true)
	if err != nil {
		return 0, err
	}
	app, err := a.GetApp()
	if err != nil {
		return 0, err
	}
	resp, err := app.CustomerService.Send(context.Background(), jsonStr)
	if err != nil {
		return 0, err
	}
	if resp.ErrCode != 0 {
		return resp.ErrCode, errors.New(resp.ErrMsg)
	}
	return 0, nil
}

func (a *Application) SendMiniProgramPage(customer, appid, title, pagePath, localThumbURL string, push *lib_define.PushMessage) (int, error) {
	app, err := a.GetApp()
	if err != nil {
		return 0, err
	}
	mediaId, errCode, err := a.UploadTempImage(localThumbURL)
	if err != nil {
		return errCode, err
	}
	jsonStr, err := messages.NewMiniProgramPage(&power.HashMap{
		`appid`:          appid,
		`title`:          title,
		`pagepath`:       pagePath,
		"thumb_media_id": mediaId,
	}).TransformForJsonRequest(gzhBuildSendMsgAppends(customer, push), true)
	if err != nil {
		return 0, err
	}

	resp, err := app.CustomerService.Send(context.Background(), jsonStr)
	if err != nil {
		return 0, err
	}
	if resp.ErrCode != 0 {
		return resp.ErrCode, errors.New(resp.ErrMsg)
	}
	return 0, nil
}

func (a *Application) SendUrl(customer, url, title string, push *lib_define.PushMessage) (int, error) {
	content := "<a href='" + url + "'>" + title + "</a>"
	jsonStr, err := messages.NewText(content).
		TransformForJsonRequest(gzhBuildSendMsgAppends(customer, push), true)
	if err != nil {
		return 0, err
	}
	app, err := a.GetApp()
	if err != nil {
		return 0, err
	}
	resp, err := app.CustomerService.Send(context.Background(), jsonStr)
	if err != nil {
		return 0, err
	}
	if resp.ErrCode != 0 {
		return resp.ErrCode, errors.New(resp.ErrMsg)
	}
	return 0, nil
}

func (a *Application) GetApp() (*officialAccount.OfficialAccount, error) {
	config := &officialAccount.UserConfig{
		AppID: a.AppID, Secret: a.Secret,
		HttpDebug: false, Debug: lib_define.IsDev,
		Cache: common.GetWechatCache(),
	}
	return officialAccount.NewOfficialAccount(config)
}

func (a *Application) SetTyping(customer, command string) (int, error) {
	app, err := a.GetApp()
	if err != nil {
		return 0, err
	}
	var resp *response2.ResponseOfficialAccount
	switch command {
	case lib_define.CommandCancelTyping:
		resp, err = app.CustomerService.HideTypingStatusToUser(context.Background(), customer)
	default:
		resp, err = app.CustomerService.ShowTypingStatusToUser(context.Background(), customer)
	}
	if err != nil {
		return 0, err
	}
	if resp.ErrCode != 0 {
		return resp.ErrCode, errors.New(resp.ErrMsg)
	}
	return 0, nil
}

func (a *Application) SendText(customer, content string, push *lib_define.PushMessage) (int, error) {
	content = common.ReplaceDate(content)
	jsonStr, err := messages.NewText(content).
		TransformForJsonRequest(gzhBuildSendMsgAppends(customer, push), true)
	if err != nil {
		return 0, err
	}
	app, err := a.GetApp()
	if err != nil {
		return 0, err
	}
	resp, err := app.CustomerService.Send(context.Background(), jsonStr)
	if err != nil {
		return 0, err
	}
	if resp.ErrCode != 0 {
		return resp.ErrCode, errors.New(resp.ErrMsg)
	}
	return 0, nil
}

func (a *Application) GetToken(refresh bool) (*response.ResponseGetToken, int, error) {
	app, err := a.GetApp()
	if err != nil {
		return nil, 0, err
	}
	resp, err := app.AccessToken.GetToken(refresh)
	if err != nil {
		return nil, 0, err
	}
	if resp.ErrCode != 0 {
		return nil, resp.ErrCode, errors.New(resp.ErrMsg)
	}
	return resp, 0, nil
}

func (a *Application) SendMsgOnEvent(_, _ string) (int, error) {
	return 0, errors.New(`not supported`)
}

func (a *Application) GetCustomerInfo(_ string) (map[string]any, int, error) {
	return nil, 0, errors.New(`not supported`)
}

func (a *Application) UploadTempImage(filePath string) (string, int, error) {
	app, err := a.GetApp()
	if err != nil {
		return ``, 0, err
	}
	resp, err := app.Media.UploadImage(context.Background(), filePath)
	if err != nil {
		return ``, 0, err
	}
	if resp.ErrCode != 0 {
		return ``, resp.ErrCode, errors.New(resp.ErrMsg)
	}
	return resp.MediaID, 0, nil
}

func (a *Application) SendImage(customer, filePath string, push *lib_define.PushMessage) (int, error) {
	app, err := a.GetApp()
	if err != nil {
		return 0, err
	}
	mediaId, errCode, err := a.UploadTempImage(filePath)
	if err != nil {
		return errCode, err
	}
	jsonStr, err := messages.NewImage(mediaId, nil).
		TransformForJsonRequest(gzhBuildSendMsgAppends(customer, push), true)
	if err != nil {
		return 0, err
	}
	resp, err := app.CustomerService.Send(context.Background(), jsonStr)
	if err != nil {
		return 0, err
	}
	if resp.ErrCode != 0 {
		return resp.ErrCode, errors.New(resp.ErrMsg)
	}
	return 0, nil
}

func (a *Application) GetFileByMedia(mediaId string, push *lib_define.PushMessage) ([]byte, http.Header, int, error) {
	app, err := a.GetApp()
	if err != nil {
		return nil, nil, 0, err
	}
	resp, err := app.Media.Get(context.Background(), mediaId)
	if err != nil {
		return nil, nil, 0, err
	}
	bytes, err := common.HttpRead(resp)
	temp := response.ResponseWork{}
	if err := json.Unmarshal(bytes, &temp); err == nil {
		return nil, nil, temp.ErrCode, errors.New(temp.ErrMsg)
	}
	return bytes, resp.Header, 0, nil
}

// GetMenu 获取菜单
func (a *Application) GetMenu() (*common.ResponseMenuGet, error) {
	app, err := a.GetApp()
	if err != nil {
		return nil, err
	}
	//resp, err := app.Menu.Get(context.Background())
	resp := &common.ResponseMenuGet{}
	_, err = app.Base.BaseClient.HttpGet(context.Background(), "cgi-bin/menu/get", nil, nil, resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, errors.New(resp.ErrMsg)
	}
	return resp, nil
}

// SetMenu 设置菜单
func (a *Application) SetMenu(menu menuRequest.RequestMenuCreate) (int, error) {
	app, err := a.GetApp()
	if err != nil {
		return 0, err
	}
	resp, err := app.Menu.Create(context.Background(), menu.Buttons)
	if err != nil {
		return 0, err
	}
	if resp.ErrCode != 0 {
		return resp.ErrCode, errors.New(resp.ErrMsg)
	}
	return 0, nil
}

// DeleteMenu 删除菜单
func (a *Application) DeleteMenu() (int, error) {
	app, err := a.GetApp()
	if err != nil {
		return 0, err
	}
	resp, err := app.Menu.Delete(context.Background())
	if err != nil {
		return 0, err
	}
	if resp.ErrCode != 0 {
		return resp.ErrCode, errors.New(resp.ErrMsg)
	}
	return 0, nil
}

func (a *Application) GetAccountBasicInfo() (*openresponse.ResponseGetBasicInfo, int, error) {
	app, err := a.GetApp()
	if err != nil {
		return nil, 0, err
	}
	resp := &openresponse.ResponseGetBasicInfo{}
	_, err = app.Base.BaseClient.HttpPostJson(context.Background(),
		`cgi-bin/account/getaccountbasicinfo`, nil, nil, nil, resp)
	if err != nil {
		return nil, 0, err
	}
	if resp.ErrCode != 0 {
		return nil, resp.ErrCode, errors.New(resp.ErrMsg)
	}
	return resp, 0, nil
}

// GetSubscribeScene 获取用户关注场景
func (a *Application) GetSubscribeScene(openid string) (string, error) {
	app, err := a.GetApp()
	if err != nil {
		return ``, err
	}
	resp, err := app.User.Get(context.Background(), openid, lib_define.LangZhCn)
	if err != nil {
		return ``, err
	}
	return resp.SubscribeScene, nil
}

// GetPublishedMessageList 获取已发布的消息列表
func (a *Application) GetPublishedMessageList(offset, count, notContent int) (*publishresponse.ResponseBatchGet, error) {
	app, err := a.GetApp()
	if err != nil {
		return nil, err
	}

	resp, err := app.Publish.PublishBatchGet(context.Background(), &publishRequest.RequestBatchGet{
		Offset:    offset,
		Count:     count,
		NoContent: notContent,
	})
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, errors.New(resp.ErrMsg)
	}
	return resp, err
}

// GetPublishedArticle 获取已发布图文信息
func (a *Application) GetPublishedArticle(articleId string) (*publishresponse.ResponsePublishGetArticle, error) {
	app, err := a.GetApp()
	if err != nil {
		return nil, err
	}

	resp, err := app.Publish.PublishGetArticle(context.Background(), articleId)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, errors.New(resp.ErrMsg)
	}
	return resp, err
}

func (a *Application) GetAccountClient() (*officialAccount.OfficialAccount, error) {
	return a.GetApp()
}

func (a *Application) SendVoice(customer, filePath string, push *lib_define.PushMessage) (int, error) {
	app, err := a.GetApp()
	if err != nil {
		return 0, err
	}
	mediaId, errCode, err := a.UploadTempVoice(filePath)
	if err != nil {
		return errCode, err
	}
	jsonStr, err := messages.NewVoice(mediaId, nil).
		TransformForJsonRequest(gzhBuildSendMsgAppends(customer, push), true)
	if err != nil {
		return 0, err
	}
	resp, err := app.CustomerService.Send(context.Background(), jsonStr)
	if err != nil {
		return 0, err
	}
	if resp.ErrCode != 0 {
		return resp.ErrCode, errors.New(resp.ErrMsg)
	}
	return 0, nil
}

func (a *Application) UploadTempVoice(filePath string) (string, int, error) {
	app, err := a.GetApp()
	if err != nil {
		return ``, 0, err
	}
	resp, err := app.Media.UploadVoice(context.Background(), filePath)
	if err != nil {
		return ``, 0, err
	}
	if resp.ErrCode != 0 {
		return ``, resp.ErrCode, errors.New(resp.ErrMsg)
	}
	return resp.MediaID, 0, nil
}
