// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package wechat_kefu

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/wechat/common"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/zhimaAi/go_tools/tool"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/response"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/accountService/message/request"
)

type Application struct {
	AppID  string
	Secret string
}

func (a *Application) GetApp() (*work.Work, error) {
	config := &work.UserConfig{
		CorpID: a.AppID, Secret: a.Secret,
		OAuth:     work.OAuth{Callback: `https://xxx.xxx`},
		HttpDebug: false, Debug: define.IsDev,
	}
	return work.NewWork(config)
}

func (a *Application) SendText(customer, content string) (int, error) {
	app, err := a.GetApp()
	if err != nil {
		return 0, err
	}
	externalUserid, openKfid := common.GetExternalUserInfo(customer)
	if len(externalUserid) == 0 || len(openKfid) == 0 {
		return 0, errors.New(`customer not exist`)
	}
	//replace the blue interactive content
	content = strings.ReplaceAll(content, `weixin://bizmsgmenu?msgmenucontent=`, `weixin://kefumenu?kefumenucontent=`)
	content = strings.ReplaceAll(content, `&msgmenuid=`, `&kefumenuid=`)
	options := &request.RequestAccountServiceSendMsg{
		ToUser: externalUserid, OpenKfid: openKfid, MsgType: `text`,
		Text: &request.RequestAccountServiceMsgText{Content: content},
	}
	resp, err := app.AccountServiceMessage.SendMsg(context.Background(), options)
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

func (a *Application) SendMsgOnEvent(code, content string) (int, error) {
	app, err := a.GetApp()
	if err != nil {
		return 0, err
	}
	//replace the blue interactive content
	content = strings.ReplaceAll(content, `weixin://bizmsgmenu?msgmenucontent=`, `weixin://kefumenu?kefumenucontent=`)
	content = strings.ReplaceAll(content, `&msgmenuid=`, `&kefumenuid=`)
	options := &request.RequestAccountServiceSendMsgOnEvent{
		Code: code, MsgType: `text`, MsgID: tool.Random(20),
		Text: request.RequestAccountServiceMsgText{Content: content},
	}
	resp, err := app.AccountServiceMessage.SendMsgOnEvent(context.Background(), options)
	if err != nil {
		return 0, err
	}
	if resp.ErrCode != 0 {
		return resp.ErrCode, errors.New(resp.ErrMsg)
	}
	return 0, nil
}

func (a *Application) GetCustomerInfo(customer string) (map[string]any, int, error) {
	app, err := a.GetApp()
	if err != nil {
		return nil, 0, err
	}
	externalUserid, openKfid := common.GetExternalUserInfo(customer)
	if len(externalUserid) == 0 || len(openKfid) == 0 {
		return nil, 0, errors.New(`customer not exist`)
	}
	resp, err := app.AccountServiceCustomer.BatchGet(context.Background(), []string{externalUserid})
	if err != nil {
		return nil, 0, err
	}
	if resp.ErrCode != 0 {
		return nil, resp.ErrCode, errors.New(resp.ErrMsg)
	}
	if len(resp.CustomerList) == 0 {
		return nil, 0, errors.New(`invalid external userid`)
	}
	return *(*map[string]any)(resp.CustomerList[0]), 0, nil
}

func (a *Application) UploadTempImage(filePath string) (string, int, error) {
	app, err := a.GetApp()
	if err != nil {
		return ``, 0, err
	}
	resp, err := app.Media.UploadTempImage(context.Background(), filePath, nil)
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
	externalUserid, openKfid := common.GetExternalUserInfo(customer)
	if len(externalUserid) == 0 || len(openKfid) == 0 {
		return 0, errors.New(`customer not exist`)
	}
	mediaId, errCode, err := a.UploadTempImage(filePath)
	if err != nil {
		return errCode, err
	}
	options := &request.RequestAccountServiceSendMsg{
		ToUser: externalUserid, OpenKfid: openKfid, MsgType: `image`,
		Image: &request.RequestAccountServiceMsgImage{MediaID: mediaId},
	}
	resp, err := app.AccountServiceMessage.SendMsg(context.Background(), options)
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
