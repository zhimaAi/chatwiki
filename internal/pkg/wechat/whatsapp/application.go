// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package whatsapp

import (
	"chatwiki/internal/pkg/lib_define"
	"errors"
	"fmt"
	"net/http"

	cams "github.com/alibabacloud-go/cams-20200606/v5/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/response"
	openresponse "github.com/ArtisanCloud/PowerWeChat/v3/src/openPlatform/authorizer/miniProgram/account/response"
)

type Application struct {
	AppID       string // WhatsApp phone number (sender From)
	Secret      string // Alibaba Cloud AccessKeySecret
	AccessKeyId string // Alibaba Cloud AccessKeyId
	CustSpaceId string // Alibaba Cloud CustSpaceId / channel ID
}

func (a *Application) sendMessage(to, messageType, content string) (errcode int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
		}
	}()

	credential, err := credentials.NewCredential(&credentials.Config{
		Type:            tea.String(`access_key`),
		AccessKeyId:     tea.String(a.AccessKeyId),
		AccessKeySecret: tea.String(a.Secret),
	})
	if err != nil {
		return 0, err
	}

	client, err := cams.NewClient(&openapi.Config{
		Credential: credential,
		Endpoint:   tea.String(`cams.ap-southeast-1.aliyuncs.com`),
	})
	if err != nil {
		return 0, err
	}

	req := &cams.SendChatappMessageRequest{
		ChannelType: tea.String(`whatsapp`),
		Type:        tea.String(`message`),
		MessageType: tea.String(messageType),
		From:        tea.String(a.AppID),
		To:          tea.String(to),
		Content:     tea.String(content),
		CustSpaceId: tea.String(a.CustSpaceId),
	}

	resp, err := client.SendChatappMessageWithOptions(req, &util.RuntimeOptions{})
	if err != nil {
		var sdkErr *tea.SDKError
		if errors.As(err, &sdkErr) {
			err = errors.New(*sdkErr.Message)
		}
		return 0, err
	}

	if resp.Body == nil || resp.Body.Code == nil || *resp.Body.Code != `OK` {
		code := ``
		msg := ``
		requestId := ``
		if resp.Body != nil {
			if resp.Body.Code != nil {
				code = *resp.Body.Code
			}
			if resp.Body.Message != nil {
				msg = *resp.Body.Message
			}
			if resp.Body.RequestId != nil {
				requestId = *resp.Body.RequestId
			}
		}
		return 0, fmt.Errorf("whatsapp send failed: code=%s message=%s requestId=%s", code, msg, requestId)
	}

	return 0, nil
}

func (a *Application) SetTyping(customer, command string) (int, error) {
	return 0, errors.New(`not supported`)
}

func (a *Application) SendText(customer, content string, push *lib_define.PushMessage) (int, error) {
	jsonContent := fmt.Sprintf(`{"text":%q}`, content)
	return a.sendMessage(customer, `text`, jsonContent)
}

func (a *Application) SendImage(customer, link string, push *lib_define.PushMessage) (int, error) {
	jsonContent := fmt.Sprintf(`{"link":%q}`, link)
	return a.sendMessage(customer, `image`, jsonContent)
}

func (a *Application) SendVideo(customer, link string, push *lib_define.PushMessage) (int, error) {
	jsonContent := fmt.Sprintf(`{"link":%q}`, link)
	return a.sendMessage(customer, `video`, jsonContent)
}

func (a *Application) SendVoice(customer, link string, push *lib_define.PushMessage) (int, error) {
	jsonContent := fmt.Sprintf(`{"link":%q}`, link)
	return a.sendMessage(customer, `audio`, jsonContent)
}

func (a *Application) SendUrl(customer, url, title string, push *lib_define.PushMessage) (int, error) {
	return 0, errors.New(`not supported`)
}

func (a *Application) SendMiniProgramPage(customer, appid, title, pagePath, localThumbURL string, push *lib_define.PushMessage) (int, error) {
	return 0, errors.New(`not supported`)
}

func (a *Application) SendImageTextLink(customer, url, title, description, localThumbURL, picurl string, push *lib_define.PushMessage) (int, error) {
	return 0, errors.New(`not supported`)
}

func (a *Application) SendSmartMenu(customer string, smartMenu lib_define.SmartMenu, push *lib_define.PushMessage) (int, error) {
	return 0, errors.New(`not supported`)
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

func (a *Application) GetFileByMedia(mediaId string, push *lib_define.PushMessage) ([]byte, http.Header, int, error) {
	return nil, nil, 0, errors.New(`not supported`)
}

func (a *Application) GetAccountBasicInfo() (*openresponse.ResponseGetBasicInfo, int, error) {
	return nil, 0, errors.New(`not supported`)
}
