// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package feishu_robot

import (
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/wechat/common"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/response"
	openresponse "github.com/ArtisanCloud/PowerWeChat/v3/src/openPlatform/authorizer/miniProgram/account/response"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

type Application struct {
	AppID  string
	Secret string
}

type FeiShuCustomer struct {
	ReceiveIdType string `json:"receive_id_type"`
	ReceiveId     any    `json:"receive_id"`
	MsgType       string `json:"msg_type"`
}

type FeiShuTextMsgContent struct {
	Text string `json:"text"`
}
type FeiShuImgMsgContent struct {
	ImageKey string `json:"image_key"`
}

type FeiShuMsgContent struct {
	FeiShuTextMsgContent
	FeiShuImgMsgContent
	MsgType string `json:"msg_type"`
}

type FeiShuPost struct {
	ZhCn FeiShuPostMsg `json:"zh_cn"`
}

type FeiShuPostMsg struct {
	Title   string             `json:"title"`
	Content [][]map[string]any `json:"content"`
}

type FeiShuPostMsgContentText struct {
	Tag   string   `json:"tag"` //  text
	Text  string   `json:"text"`
	Style []string `json:"style"`
}
type FeiShuPostMsgContentUrl struct {
	Tag   string   `json:"tag"` // a
	Text  string   `json:"text"`
	Href  string   `json:"href"`
	Style []string `json:"style"`
}

type FeiShuPostMsgContentAt struct {
	Tag    string   `json:"tag"` //at
	UserId string   `json:"user_id"`
	Style  []string `json:"style"`
}

type FeiShuPostMsgContentImage struct {
	Tag      string `json:"tag"` // img
	ImageKey string `json:"image_key"`
}
type FeiShuPostMsgContentMedia struct {
	Tag      string `json:"tag"` // media
	FileKey  string `json:"file_key"`
	ImageKey string `json:"image_key"`
}

type FeiShuPostMsgContentEmotion struct {
	Tag       string `json:"tag"` // emotion
	EmojiType string `json:"emoji_type"`
}

type FeiShuPostMsgContentHr struct {
	Tag string `json:"tag"` // hr
}

type FeiShuPostMsgContentCodeBlock struct {
	Tag      string `json:"tag"`      // code_block
	Language string `json:"language"` // GO
	Text     string `json:"text"`
}

type FeiShuPostMsgContentMd struct {
	Tag  string `json:"tag"` // md
	Text string `json:"text"`
}

func structToMapAny(v interface{}) map[string]any {
	var m map[string]any
	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return m
	}
	return m
}

func (a *Application) GetApp() (*lark.Client, error) {

	return lark.NewClient(a.AppID, a.Secret), nil
}

func (a *Application) getReceiveInfo(push *lib_define.PushMessage) (receiveIdType string, receiveId string) {
	receiveIdType = "open_id"
	receiveId = cast.ToString(push.Message["FromUserName"])
	if push.Message["SessionType"] == "group" { //群聊消息
		receiveIdType = "chat_id"
		receiveId = cast.ToString(push.Message["ToUserName"])
	}
	return
}

func (a *Application) SendContent(msgType, content, receiveIdType, receiveId string) (int, error) {
	// 获取应用实例
	app, err := a.GetApp()
	if err != nil {
		return 0, err
	}

	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(receiveIdType).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(receiveId).
			MsgType(msgType).
			Content(content).
			Uuid(tool.Random(30)).
			Build()).
		Build()

	// 发起请求
	resp, err := app.Im.V1.Message.Create(context.Background(), req)
	fmt.Printf("content: \n%s", content)

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Printf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
		return 0, err
	}

	return 1, nil
}

func (a *Application) SendText(customer, content string, push *lib_define.PushMessage) (int, error) {
	// 获取接收者信息
	receiveIdType, receiveId := a.getReceiveInfo(push)
	// 获取基础消息信息
	replyContent := a.getTextBaseContent(push)
	//新发消息
	content = common.ReplaceDate(content)
	replyContent += content

	// 创建文本消息内容
	contentStr, _ := tool.JsonEncode(FeiShuTextMsgContent{Text: replyContent})

	// 调用通用发送方法
	return a.SendContent("text", contentStr, receiveIdType, receiveId)
}

func (a *Application) getTextBaseContent(push *lib_define.PushMessage) string {
	// 客户发送的数据
	replyContent := ""
	msgType := strings.ToLower(cast.ToString(push.Message[`MsgType`]))
	if push.Content != "" && msgType == lib_define.FeShuMsgTypeText {
		replyContent = push.Content + " \n ---------- \n"
	}

	// 添加@信息（如果是群聊）
	if push.Message["SessionType"] == "group" {
		replyContent += "<at user_id=\"" + cast.ToString(push.Message["FromUserName"]) + "\">Tom</at> " //at对应的人
	}
	return replyContent
}

func (a *Application) GetToken(refresh bool) (*response.ResponseGetToken, int, error) {
	return nil, 0, nil
}

func (a *Application) SendMsgOnEvent(code, content string) (int, error) {
	//app, err := a.GetApp()

	return 0, nil
}

func (a *Application) GetCustomerInfo(customer string) (map[string]any, int, error) {

	return map[string]any{}, 0, nil
}

func (a *Application) UploadTempImage(filePath string) (string, int, error) {
	app, err := a.GetApp()
	if err != nil {
		return "", 0, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		logs.Error("读取文件错误：" + err.Error())
		return "", 0, err
	}

	req := larkim.NewCreateImageReqBuilder().
		Body(larkim.NewCreateImageReqBodyBuilder().
			ImageType(`message`).
			Image(file).
			Build()).
		Build()

	// 发起请求
	resp, err := app.Im.V1.Image.Create(context.Background(), req)
	// 处理错误
	if err != nil {
		fmt.Println(err)
		return "", 0, err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Printf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
		return "", 0, err
	}

	return *resp.Data.ImageKey, 0, nil

}

func (a *Application) SendImage(customer, filePath string, push *lib_define.PushMessage) (int, error) {
	imageKey, _, err := a.UploadTempImage(filePath)

	if imageKey == "" {
		logs.Error("图片上传错误")
		return 0, err
	}

	// 获取接收者信息
	receiveIdType, receiveId := a.getReceiveInfo(push)

	// 创建图片消息内容
	contentStr, _ := tool.JsonEncode(FeiShuImgMsgContent{ImageKey: imageKey})

	// 调用通用发送方法
	return a.SendContent("image", contentStr, receiveIdType, receiveId)
}

func (a *Application) SendUrl(customer, url, title string, push *lib_define.PushMessage) (int, error) {
	// 获取接收者信息
	receiveIdType, receiveId := a.getReceiveInfo(push)
	// 获取基础消息信息
	replyContent := a.getTextBaseContent(push)
	// 需要发送的数据
	replyContent += "[" + title + "](" + url + ")"

	// 创建文本消息内容
	contentStr, _ := tool.JsonEncode(FeiShuTextMsgContent{Text: replyContent})

	// 调用通用发送方法
	return a.SendContent("text", contentStr, receiveIdType, receiveId)
}

func (a *Application) SendMiniProgramPage(customer, appid, title, pagePath, localThumbURL string, push *lib_define.PushMessage) (int, error) {

	imageKey, _, err := a.UploadTempImage(localThumbURL)

	if imageKey == "" {
		logs.Error("图片上传错误")
		return 0, err
	}

	// 获取接收者信息
	receiveIdType, receiveId := a.getReceiveInfo(push)
	// 获取基础消息信息
	var replyContent = make([][]map[string]any, 0)

	var firstLine = make([]map[string]any, 0)

	// 添加@信息（如果是群聊）
	if push.Message["SessionType"] == "group" {
		firstLine = append(firstLine, structToMapAny(FeiShuPostMsgContentAt{
			Tag:    "at",
			UserId: cast.ToString(push.Message["FromUserName"]),
		}))
	}

	firstLine = append(firstLine, structToMapAny(FeiShuPostMsgContentText{
		Tag:  "text",
		Text: "小程序Appid：" + appid,
	}))

	var secondLine = make([]map[string]any, 0)
	secondLine = append(secondLine, structToMapAny(FeiShuPostMsgContentText{
		Tag:  "text",
		Text: "小程序链接：" + pagePath,
	}))

	var thirdLine = make([]map[string]any, 0)
	thirdLine = append(thirdLine, structToMapAny(FeiShuPostMsgContentImage{
		Tag:      "img",
		ImageKey: imageKey,
	}))
	replyContent = append(replyContent, firstLine, secondLine, thirdLine)

	// 创建富文本消息内容
	var postMsg = FeiShuPostMsg{
		Title:   title,
		Content: replyContent,
	}
	// 创建文本消息内容
	contentStr, _ := tool.JsonEncode(FeiShuPost{ZhCn: postMsg})

	// 调用通用发送方法
	return a.SendContent("post", contentStr, receiveIdType, receiveId)
}

func (a *Application) SendImageTextLink(customer, url, title, description, localThumbURL, picurl string, push *lib_define.PushMessage) (int, error) {
	imageKey, _, err := a.UploadTempImage(localThumbURL)

	if imageKey == "" {
		logs.Error("图片上传错误")
		return 0, err
	}
	// 获取接收者信息
	receiveIdType, receiveId := a.getReceiveInfo(push)
	// 获取基础消息信息
	var replyContent = make([][]map[string]any, 0)

	var firstLine = make([]map[string]any, 0)
	// 添加@信息（如果是群聊）
	if push.Message["SessionType"] == "group" {
		firstLine = append(firstLine, structToMapAny(FeiShuPostMsgContentAt{
			Tag:    "at",
			UserId: cast.ToString(push.Message["FromUserName"]),
		}))
	}
	firstLine = append(firstLine, structToMapAny(FeiShuPostMsgContentText{
		Tag:  "text",
		Text: "介绍：" + description,
	}))

	var secondLine = make([]map[string]any, 0)
	secondLine = append(secondLine, structToMapAny(FeiShuPostMsgContentUrl{
		Tag:  "a",
		Text: "跳转链接：" + url,
		Href: url,
	}))

	var thirdLine = make([]map[string]any, 0)
	thirdLine = append(thirdLine, structToMapAny(FeiShuPostMsgContentImage{
		Tag:      "img",
		ImageKey: imageKey,
	}))
	replyContent = append(replyContent, firstLine, secondLine, thirdLine)

	// 创建富文本消息内容
	var postMsg = FeiShuPostMsg{
		Title:   title,
		Content: replyContent,
	}
	// 创建文本消息内容
	contentStr, _ := tool.JsonEncode(FeiShuPost{ZhCn: postMsg})

	// 调用通用发送方法
	return a.SendContent("post", contentStr, receiveIdType, receiveId)
}

func (a *Application) SendSmartMenu(customer string, smartMenu lib_define.SmartMenu, push *lib_define.PushMessage) (int, error) {
	// 获取接收者信息
	receiveIdType, receiveId := a.getReceiveInfo(push)
	// 获取基础消息信息
	var replyContent = make([][]map[string]any, 0)
	description := common.ProcessEscapeSequences(smartMenu.MenuDescription)
	description = common.ReplaceDate(description)
	// 按\n切割smartMenu.MenuDescription，然后按行添加
	descriptionLines := strings.Split(description, "\n")

	// 添加@信息（如果是群聊）
	if push.Message["SessionType"] == "group" {
		var firstLine = make([]map[string]any, 0)
		firstLine = append(firstLine, structToMapAny(FeiShuPostMsgContentAt{
			Tag:    "at",
			UserId: cast.ToString(push.Message["FromUserName"]),
		}))
		replyContent = append(replyContent, firstLine)
	}

	// 按行添加菜单描述内容
	for _, line := range descriptionLines {
		var lineContent = make([]map[string]any, 0)
		lineContent = append(lineContent, structToMapAny(FeiShuPostMsgContentText{
			Tag:  "text",
			Text: line,
		}))
		replyContent = append(replyContent, lineContent)
	}

	// 遍历菜单内容
	if len(smartMenu.MenuContent) > 0 {
		for _, content := range smartMenu.MenuContent {
			var line = make([]map[string]any, 0)
			// 判断是普通文本还是链接
			if content.SerialNo != `` {
				line = append(line, structToMapAny(FeiShuPostMsgContentText{
					Tag:  "text",
					Text: content.SerialNo + ` `,
				}))
			}

			//反向解析 a标签 是链接， 还是小程序
			linkInfo := common.ContentToALabel(content.Content)
			switch linkInfo.Type {
			case common.PlainText:
				line = append(line, structToMapAny(FeiShuPostMsgContentText{
					Tag:  "text",
					Text: linkInfo.Text,
				}))
				break
			case common.NormalLink:
				line = append(line, structToMapAny(FeiShuPostMsgContentUrl{
					Tag:  "a",
					Text: linkInfo.Text,
					Href: linkInfo.URL,
				}))
				break
			case common.MiniProgramLink:
				line = append(line, structToMapAny(FeiShuPostMsgContentText{
					Tag:  "text",
					Text: linkInfo.Text,
				}))
				break

			}
			replyContent = append(replyContent, line)
		}
	}

	// 创建富文本消息内容
	var postMsg = FeiShuPostMsg{
		Title:   `智能菜单`,
		Content: replyContent,
	}
	// 创建文本消息内容
	contentStr, _ := tool.JsonEncode(FeiShuPost{ZhCn: postMsg})

	// 调用通用发送方法
	return a.SendContent("post", contentStr, receiveIdType, receiveId)
}

func (a *Application) GetAccountBasicInfo() (*openresponse.ResponseGetBasicInfo, int, error) {
	return nil, 0, errors.New(`not supported`)
}

func (a *Application) GetFileByMedia(mediaId string, push *lib_define.PushMessage) ([]byte, http.Header, int, error) {
	app, err := a.GetApp()
	if err != nil {
		return nil, nil, 0, err
	}

	req := larkim.NewGetMessageResourceReqBuilder().
		MessageId(cast.ToString(push.Message[`MsgId`])).
		FileKey(mediaId).
		Type(cast.ToString(push.Message[`MsgType`])).
		Build()

	// 发起请求
	resp, err := app.Im.V1.MessageResource.Get(context.Background(), req)
	// 发起请求
	if err != nil {
		return nil, nil, 0, err
	}
	// 服务端错误处理
	if !resp.Success() {
		return nil, nil, 0, errors.New(fmt.Sprintf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError)))
	}
	// 业务处理
	bytes := resp.RawBody
	return bytes, resp.Header, 0, nil
}
