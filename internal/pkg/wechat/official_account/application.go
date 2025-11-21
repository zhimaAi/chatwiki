// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package official_account

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/wechat/common"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/messages"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/power"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/response"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount"
	openresponse "github.com/ArtisanCloud/PowerWeChat/v3/src/openPlatform/authorizer/miniProgram/account/response"
)

type Application struct {
	AppID  string
	Secret string
}

func (a *Application) SendImageTextLink(customer, url, title, description, localThumbURL, picurl string, push *define.PushMessage) (int, error) {
	jsonStr, err := messages.NewLink(&power.HashMap{
		`url`:         url,
		`title`:       title,
		`description`: description,
	}).TransformForJsonRequest(&object.HashMap{`touser`: customer}, true)
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

func (a *Application) SendMiniProgramPage(customer, appid, title, pagePath, localThumbURL string, push *define.PushMessage) (int, error) {
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
	}).TransformForJsonRequest(&object.HashMap{`touser`: customer}, true)
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

func (a *Application) SendUrl(customer, url, title string, push *define.PushMessage) (int, error) {
	content := "<a href='" + url + "'>" + title + "</a>"
	jsonStr, err := messages.NewText(content).
		TransformForJsonRequest(&object.HashMap{`touser`: customer}, true)
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

func (a *Application) SendText(customer, content string, push *define.PushMessage) (int, error) {
	jsonStr, err := messages.NewText(content).
		TransformForJsonRequest(&object.HashMap{`touser`: customer}, true)
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

func (a *Application) SendImage(customer, filePath string, push *define.PushMessage) (int, error) {
	app, err := a.GetApp()
	if err != nil {
		return 0, err
	}
	mediaId, errCode, err := a.UploadTempImage(filePath)
	if err != nil {
		return errCode, err
	}
	jsonStr, err := messages.NewImage(mediaId, nil).
		TransformForJsonRequest(&object.HashMap{`touser`: customer}, true)
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

func (a *Application) GetFileByMedia(mediaId string, push *define.PushMessage) ([]byte, http.Header, int, error) {
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

// GetSubscribeScene 获取用户关注场景
func (a *Application) GetSubscribeScene(openid string) (string, error) {
	app, err := a.GetApp()
	if err != nil {
		return ``, err
	}
	resp, err := app.User.Get(context.Background(), openid, define.LangZhCn)
	if err != nil {
		return ``, err
	}
	return resp.SubscribeScene, nil
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
