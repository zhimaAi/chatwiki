// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package mini_program

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/wechat/common"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/response"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/miniProgram"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/miniProgram/customerServiceMessage/request"
)

type Application struct {
	AppID  string
	Secret string
}

func (a *Application) GetApp() (*miniProgram.MiniProgram, error) {
	config := &miniProgram.UserConfig{
		AppID: a.AppID, Secret: a.Secret,
		HttpDebug: false, Debug: define.IsDev,
	}
	return miniProgram.NewMiniProgram(config)
}

func (a *Application) SendText(customer, content string) (int, error) {
	app, err := a.GetApp()
	if err != nil {
		return 0, err
	}
	resp, err := app.CustomerServiceMessage.SendText(context.Background(),
		customer, &request.CustomerServiceMsgText{Content: content})
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
	resp, err := app.CustomerServiceMessage.UploadTempMedia(context.Background(), `image`, filePath, nil)
	if err != nil {
		return ``, 0, err
	}
	if resp.ErrCode != 0 {
		return ``, resp.ErrCode, errors.New(resp.ErrMsg)
	}
	return resp.MediaID, 0, nil
}

func (a *Application) SendImage(customer, filePath string) (int, error) {
	app, err := a.GetApp()
	if err != nil {
		return 0, err
	}
	mediaId, errCode, err := a.UploadTempImage(filePath)
	if err != nil {
		return errCode, err
	}
	resp, err := app.CustomerServiceMessage.SendImage(context.Background(),
		customer, &request.CustomerServiceMsgImage{MediaID: mediaId})
	if err != nil {
		return 0, err
	}
	if resp.ErrCode != 0 {
		return resp.ErrCode, errors.New(resp.ErrMsg)
	}
	return 0, nil
}

func (a *Application) GetFileByMedia(mediaId string) ([]byte, http.Header, int, error) {
	app, err := a.GetApp()
	if err != nil {
		return nil, nil, 0, err
	}
	resp, err := app.CustomerServiceMessage.GetTempMedia(context.Background(), mediaId)
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
