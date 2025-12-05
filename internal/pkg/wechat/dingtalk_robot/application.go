// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package dingtalk_robot

import (
	"bytes"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/wechat/common"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/response"
	openresponse "github.com/ArtisanCloud/PowerWeChat/v3/src/openPlatform/authorizer/miniProgram/account/response"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dingtalkoauth2_1_0 "github.com/alibabacloud-go/dingtalk/oauth2_1_0"
	dingtalkrobot_1_0 "github.com/alibabacloud-go/dingtalk/robot_1_0"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

type Application struct {
	AppID  string
	Secret string
}

type DingtalkUploadMedia struct {
	Errcode   int    `json:"errcode"`
	Errmsg    string `json:"errmsg"`
	MediaId   string `json:"media_id"`
	CreatedAt int64  `json:"created_at"`
	Type      string `json:"type"`
}

var DingTaklHost = "https://oapi.dingtalk.com"

var DingtalkMsgTypeMap = map[string]string{
	"sampleText":     "SampleText",
	"sampleMarkdown": "SampleMarkdown",
	"sampleImageMsg": "SampleImageMsg",
	"sampleLink":     "SampleLink",
}

type DingtalkSampleText struct {
	Content string `json:"content"`
}
type DingtalkSampleMarkdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}
type DingtalkSampleImageMsg struct {
	PhotoURL string `json:"photoURL"`
}
type DingtalkSampleLink struct {
	Text       string `json:"text"`
	Title      string `json:"title"`
	PicUrl     string `json:"picUrl"`
	MessageUrl string `json:"messageUrl"`
}

type DingTalkMsgType struct {
	SampleText     DingtalkSampleText     `json:"sampleText"`     //文本类型
	SampleMarkdown DingtalkSampleMarkdown `json:"sampleMarkdown"` //markdown
	SampleImageMsg DingtalkSampleImageMsg `json:"sampleImageMsg"` //图片
	SampleLink     DingtalkSampleLink     `json:"sampleLink"`     //链接
	MsgKey         string                 `json:"msg_key"`
}

func (a *Application) SendUrl(customer, url, title string, push *lib_define.PushMessage) (int, error) {

	markDownText := "[" + title + "](" + url + ")"

	//构建消息消息
	dingtalkMsg, _ := tool.JsonEncode(DingTalkMsgType{MsgKey: "sampleMarkdown", SampleMarkdown: DingtalkSampleMarkdown{
		Title: title,
		Text:  markDownText,
	}})

	return a.SendText("", dingtalkMsg, push)
}

func (a *Application) SendMiniProgramPage(customer, appid, title, pagePath, localThumbURL string, push *lib_define.PushMessage) (int, error) {
	mediaId, _, err := a.UploadTempImage(localThumbURL)
	if err != nil {
		return 0, nil
	}

	markDownText := "### " + title + "  \n  小程序APPID：" + appid + "  \n  小程序链接：" + pagePath + "![](" + mediaId + ")"

	//构建消息消息
	dingtalkMsg, _ := tool.JsonEncode(DingTalkMsgType{MsgKey: "sampleMarkdown", SampleMarkdown: DingtalkSampleMarkdown{
		Title: title,
		Text:  markDownText,
	}})

	return a.SendText("", dingtalkMsg, push)
}

func (a *Application) SendSmartMenu(customer string, smartMenu lib_define.SmartMenu, push *lib_define.PushMessage) (int, error) {
	description := common.ProcessEscapeSequences(smartMenu.MenuDescription)
	description = common.ReplaceDate(description)
	// 按\n切割smartMenu.MenuDescription，然后按行添加
	descriptionLines := strings.Split(description, "\n")

	markDownText := ``
	for _, line := range descriptionLines {
		markDownText += line + "  \n  "
	}

	// 遍历菜单内容
	if len(smartMenu.MenuContent) > 0 {
		for _, content := range smartMenu.MenuContent {
			var line = ``
			// 判断是普通文本还是链接
			if content.SerialNo != `` {
				line += content.SerialNo + ` `
			}
			//反向解析 a标签 是链接， 还是小程序
			linkInfo := common.ContentToALabel(content.Content)
			switch linkInfo.Type {
			case common.PlainText:
				line += linkInfo.Text
				break
			case common.NormalLink:
				line += "[" + linkInfo.Text + "](" + linkInfo.URL + ")"
				break
			case common.MiniProgramLink:
				//小程序
				line += linkInfo.Text
				break
			}
			markDownText += line + "  \n  "
		}
	}
	//构建消息消息
	dingtalkMsg, _ := tool.JsonEncode(DingTalkMsgType{MsgKey: "sampleMarkdown", SampleMarkdown: DingtalkSampleMarkdown{
		Title: `智能菜单`,
		Text:  markDownText,
	}})

	return a.SendText("", dingtalkMsg, push)
}

func (a *Application) SendImageTextLink(customer, url, title, description, localThumbURL, picurl string, push *lib_define.PushMessage) (int, error) {
	mediaId, _, err := a.UploadTempImage(localThumbURL)
	if err != nil {

		return 0, nil
	}

	//构建图片消息
	dingtalkMsg, _ := tool.JsonEncode(DingTalkMsgType{MsgKey: "sampleLink", SampleLink: DingtalkSampleLink{
		Text:       description,
		Title:      title,
		PicUrl:     mediaId,
		MessageUrl: url,
	}})

	return a.SendText("", dingtalkMsg, push)

}

func (a *Application) GetAccountBasicInfo() (*openresponse.ResponseGetBasicInfo, int, error) {
	return nil, 0, errors.New(`not supported`)
}

func (a *Application) GetApp() (_result *dingtalkrobot_1_0.Client, _err error) {
	config := &openapi.Config{}
	config.Protocol = tea.String("https")
	config.RegionId = tea.String("central")
	_result = &dingtalkrobot_1_0.Client{}
	_result, _err = dingtalkrobot_1_0.NewClient(config)

	return _result, _err
}

func (a *Application) SendText(customer, content string, push *lib_define.PushMessage) (int, error) {
	app, err := a.GetApp()
	if err != nil {
		return 0, err
	}

	DingtalkMsg := ""

	dingTalkMsgType := DingTalkMsgType{}
	_ = tool.JsonDecode(content, &dingTalkMsgType)

	if dingTalkMsgType.MsgKey == "" { //默认空值为文本消息
		dingTalkMsgType.MsgKey = "sampleText"
		dingTalkMsgType.SampleText = DingtalkSampleText{Content: content}
	}

	//反射取一下对应的值
	ref := reflect.ValueOf(dingTalkMsgType)
	field := ref.FieldByName(DingtalkMsgTypeMap[dingTalkMsgType.MsgKey])
	if field.IsValid() {
		DingtalkMsg, _ = tool.JsonEncode(field.Interface())
	}

	responseToken, _, err := a.GetToken(false)

	if err != nil && responseToken.AccessToken == "" {
		return 0, err
	}

	respErr := *new(error)
	if cast.ToInt(push.Message["SessionType"]) == 2 { //群聊
		orgGroupSendHeaders := &dingtalkrobot_1_0.OrgGroupSendHeaders{}
		orgGroupSendHeaders.XAcsDingtalkAccessToken = tea.String(responseToken.AccessToken)
		orgGroupSendRequest := &dingtalkrobot_1_0.OrgGroupSendRequest{
			MsgParam:           tea.String(DingtalkMsg),
			MsgKey:             tea.String(dingTalkMsgType.MsgKey),
			OpenConversationId: tea.String(cast.ToString(push.Message["ToUserName"])),
			RobotCode:          tea.String(cast.ToString(push.Message["RobotCode"])),
		}

		_, respErr = app.OrgGroupSendWithOptions(orgGroupSendRequest, orgGroupSendHeaders, &util.RuntimeOptions{})

	} else { //单聊
		batchSendOTOHeaders := &dingtalkrobot_1_0.BatchSendOTOHeaders{}
		batchSendOTOHeaders.XAcsDingtalkAccessToken = tea.String(responseToken.AccessToken)
		batchSendOTORequest := &dingtalkrobot_1_0.BatchSendOTORequest{
			RobotCode: tea.String(cast.ToString(push.Message["RobotCode"])),
			UserIds:   []*string{tea.String(cast.ToString(push.Message["FromUserName"]))},
			MsgKey:    tea.String(dingTalkMsgType.MsgKey),
			MsgParam:  tea.String(DingtalkMsg),
		}
		_, respErr = app.BatchSendOTOWithOptions(batchSendOTORequest, batchSendOTOHeaders, &util.RuntimeOptions{})

	}

	if respErr != nil {
		var err = &tea.SDKError{}
		if _t, ok := respErr.(*tea.SDKError); ok {
			err = _t
		} else {
			err.Message = tea.String(respErr.Error())
		}
		if !tea.BoolValue(util.Empty(err.Code)) && !tea.BoolValue(util.Empty(err.Message)) {
			logs.Error("错误码：" + *err.Code + " 错误内容：" + *err.Message)
			return 0, errors.New(*err.Message)
		}

	}

	return 1, nil
}

func (a *Application) GetToken(refresh bool) (*response.ResponseGetToken, int, error) {
	config := &openapi.Config{}
	config.Protocol = tea.String("https")
	config.RegionId = tea.String("central")
	client, err := dingtalkoauth2_1_0.NewClient(config)
	if err != nil {
		return nil, 0, err
	}

	resp := response.ResponseGetToken{}

	redisKey := fmt.Sprintf(lib_define.RedisPrefixDingtalkAccessToken, a.AppID, a.Secret)
	cacheRes, err := lib_define.Redis.Get(context.Background(), redisKey).Result()

	if err != nil {
		logs.Error(err.Error())
	}

	//如果拿到了，退出
	if err == nil && cacheRes != "" {
		resp.AccessToken = cacheRes
		return &resp, 0, nil
	}

	getAccessTokenRequest := &dingtalkoauth2_1_0.GetAccessTokenRequest{
		AppKey:    tea.String(a.AppID),
		AppSecret: tea.String(a.Secret),
	}
	AccessToken := ""
	ExpireIn := 0
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		res, err := client.GetAccessToken(getAccessTokenRequest)
		if err != nil {
			return err
		}

		AccessToken = *res.Body.AccessToken
		ExpireIn = cast.ToInt(*res.Body.ExpireIn)
		return nil
	}()

	if tryErr != nil {
		var err = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			err = _t
		} else {
			err.Message = tea.String(tryErr.Error())
		}
		if !tea.BoolValue(util.Empty(err.Code)) && !tea.BoolValue(util.Empty(err.Message)) {
			// err 中含有 code 和 message 属性，可帮助开发定位问题
			logs.Error("错误码：" + *err.Code + "，错误内容：" + *err.Message)
		}

	}

	resp.AccessToken = AccessToken

	//设置一下缓存，缓存时间减半
	_, err = lib_define.Redis.Set(context.Background(), redisKey, AccessToken, time.Duration(ExpireIn/2)*time.Second).Result()

	if err != nil {
		logs.Error(err.Error())
	}

	return &resp, 0, nil
}

func (a *Application) SendMsgOnEvent(code, content string) (int, error) {
	//app, err := a.GetApp()

	return 0, nil
}

func (a *Application) GetCustomerInfo(customer string) (map[string]any, int, error) {

	return map[string]any{}, 0, nil
}

func (a *Application) UploadTempImage(filePath string) (string, int, error) {
	responseToken, _, err := a.GetToken(false)

	if responseToken.AccessToken == "" {
		return "", 0, errors.New("miss token")
	}
	// 请求URL
	url := DingTaklHost + "/media/upload?access_token=" + responseToken.AccessToken

	// 创建缓冲区用于存储multipart表单数据
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// 添加type字段
	_ = writer.WriteField("type", "image")

	// 添加媒体文件
	file, err := os.Open(filePath)
	if err != nil {

		return "", 0, errors.New("open file error")
	}
	defer file.Close()

	part, err := writer.CreateFormFile("media", filepath.Base(filePath))
	if err != nil {
		fmt.Printf("创建表单文件失败: %v\n", err)
		return "", 0, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Printf("写入文件数据失败: %v\n", err)
		return "", 0, err
	}

	// 关闭multipart writer
	err = writer.Close()
	if err != nil {
		fmt.Printf("关闭writer失败: %v\n", err)
		return "", 0, errors.New("open file error")
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		fmt.Printf("创建请求失败: %v\n", err)
		return "", 0, err
	}

	// 设置Content-Type
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("发送请求失败: %v\n", err)
		return "", 0, err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应失败: %v\n", err)
		return "", 0, err
	}

	media := DingtalkUploadMedia{}
	err = tool.JsonDecode(string(body), &media)

	if err != nil {
		return "", 0, err
	}

	return media.MediaId, 1, nil

}

func (a *Application) SendImage(customer, filePath string, push *lib_define.PushMessage) (int, error) {

	mediaId, _, err := a.UploadTempImage(filePath)

	if err != nil {

		return 0, nil
	}

	//构建图片消息
	dingtalkMsg, _ := tool.JsonEncode(DingTalkMsgType{MsgKey: "sampleImageMsg", SampleImageMsg: DingtalkSampleImageMsg{PhotoURL: mediaId}})

	_, _ = a.SendText("", dingtalkMsg, push)

	return 0, nil
}

func (a *Application) GetFileByMedia(mediaId string, push *lib_define.PushMessage) ([]byte, http.Header, int, error) {
	app, err := a.GetApp()
	if err != nil {
		return nil, nil, 0, err
	}

	responseToken, _, err := a.GetToken(false)

	if err != nil && responseToken.AccessToken == "" {
		return nil, nil, 0, err
	}

	robotMessageFileDownloadHeaders := &dingtalkrobot_1_0.RobotMessageFileDownloadHeaders{}
	robotMessageFileDownloadHeaders.XAcsDingtalkAccessToken = tea.String(responseToken.AccessToken)
	robotMessageFileDownloadRequest := &dingtalkrobot_1_0.RobotMessageFileDownloadRequest{
		DownloadCode: tea.String(mediaId),
		RobotCode:    tea.String(cast.ToString(push.Message["RobotCode"])),
	}

	downloadRed, err := app.RobotMessageFileDownloadWithOptions(robotMessageFileDownloadRequest, robotMessageFileDownloadHeaders, &util.RuntimeOptions{})

	if err != nil {
		return nil, nil, 0, err
	}

	resp, _ := http.Get(*downloadRed.Body.DownloadUrl)

	defer func() {
		_ = resp.Body.Close()
	}()

	bodyByte, _ := io.ReadAll(resp.Body)

	return bodyByte, resp.Header, 0, nil
}
