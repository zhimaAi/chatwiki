// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package wecom_robot

import (
	"chatwiki/internal/pkg/lib_define"
	"errors"
	"net/http"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/response"
	openresponse "github.com/ArtisanCloud/PowerWeChat/v3/src/openPlatform/authorizer/miniProgram/account/response"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/curl"
)

type Application struct {
	AppID string
}

func (a *Application) SendSmartMenu(customer string, smartMenu lib_define.SmartMenu, push *lib_define.PushMessage) (int, error) {
	return 0, errors.New(`not supported`)
}

func (a *Application) SendImageTextLink(customer, url, title, description, localThumbURL, picurl string, push *lib_define.PushMessage) (int, error) {
	return 0, errors.New(`not supported`)
}

func (a *Application) SendMiniProgramPage(customer, appid, title, pagePath, localThumbURL string, push *lib_define.PushMessage) (int, error) {
	return 0, errors.New(`not supported`)
}

func (a *Application) SendUrl(customer, url, title string, push *lib_define.PushMessage) (int, error) {
	return 0, errors.New(`not supported`)
}

func (a *Application) SetTyping(customer, command string) (int, error) {
	return 0, errors.New(`not supported`)
}

func (a *Application) SendText(customer, content string, push *lib_define.PushMessage) (int, error) {
	responseUrl := cast.ToString(push.Message[`response_url`])
	if len(responseUrl) == 0 {
		return 0, errors.New(`response_url is empty`)
	}
	body := map[string]any{
		`msgtype`:  `markdown`,
		`markdown`: map[string]any{`content`: content},
	}
	request, err := curl.Post(responseUrl).JSONBody(body)
	if err != nil {
		return 0, err
	}
	resp := struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}{}
	if err = request.ToJSON(&resp); err != nil {
		return 0, err
	}
	if resp.ErrCode != 0 {
		return resp.ErrCode, errors.New(resp.ErrMsg)
	}
	delete(push.Message, `response_url`) //release it it can only be used once
	return 0, nil
}

func (a *Application) GetToken(refresh bool) (*response.ResponseGetToken, int, error) {
	return &response.ResponseGetToken{}, 0, nil
}

func (a *Application) SendMsgOnEvent(code, content string) (int, error) {
	return 0, errors.New(`not supported`)
}

func (a *Application) GetCustomerInfo(customer string) (map[string]any, int, error) {
	return nil, 0, errors.New(`not supported`)
}

func (a *Application) UploadTempImage(filePath string) (string, int, error) {
	return ``, 0, errors.New(`not supported`)
}

func (a *Application) SendImage(customer, filePath string, push *lib_define.PushMessage) (int, error) {
	return 0, errors.New(`not supported`)
}

func (a *Application) GetFileByMedia(mediaId string, push *lib_define.PushMessage) ([]byte, http.Header, int, error) {
	return nil, nil, 0, errors.New(`not supported`)
}

func (a *Application) GetAccountBasicInfo() (*openresponse.ResponseGetBasicInfo, int, error) {
	return nil, 0, errors.New(`not supported`)
}

func (a *Application) SendVoice(customer, filePath string, push *lib_define.PushMessage) (int, error) {
	return 0, errors.New(`not supported`)
}

func (a *Application) SendVideo(customer, filePath string, push *lib_define.PushMessage) (int, error) {
	return 0, errors.New(`not supported`)
}
