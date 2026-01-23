// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package feishu_robot

import (
	"bytes"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/wechat/common"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/response"
	openresponse "github.com/ArtisanCloud/PowerWeChat/v3/src/openPlatform/authorizer/miniProgram/account/response"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
	larkdrive "github.com/larksuite/oapi-sdk-go/v3/service/drive/v1"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

type Application struct {
	AppID  string
	Secret string
}

type FeishuUserAuthLoginState struct {
	AppId            string `json:"app_id"`
	AppSecret        string `json:"app_secret"`
	FrontRedirectUrl string `json:"front_redirect_url"`
}
type FeishuUserAccessToken struct {
	Code                  int    `json:"code"`
	Error                 string `json:"error"`
	ErrorDescription      string `json:"error_description"`
	AccessToken           string `json:"access_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	TokenType             string `json:"token_type"`
	Scope                 string `json:"scope"`
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

type FeishuDocFile struct {
	Name       string `json:"name"`
	DocumentId string `json:"document_id"`
}

// FeishuDocFileTree 飞书文档树形结构
type FeishuDocFileTree struct {
	Name         string               `json:"name"`
	Token        string               `json:"token"`
	Type         string               `json:"type"`
	ParentToken  string               `json:"parent_token"`
	CreatedTime  string               `json:"created_time"`
	ModifiedTime string               `json:"modified_time"`
	OwnerID      string               `json:"owner_id"`
	URL          string               `json:"url"`
	Children     []*FeishuDocFileTree `json:"children,omitempty"`
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

func (a *Application) SetTyping(customer, command string) (int, error) {
	return 0, errors.New(`not supported`)
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

func (a *Application) SendVoice(customer, filePath string, push *lib_define.PushMessage) (int, error) {
	return 0, errors.New(`miniprogram not supported voice `)
}

// GetDocFileList 获取飞书doc类型文档列表
func (a *Application) GetDocFileList(userAccessToken string) ([]FeishuDocFile, error) {
	// 获取应用实例
	app, err := a.GetApp()
	if err != nil {
		return nil, err
	}

	var allDocFiles []FeishuDocFile
	pageToken := ""
	pageSize := 200 // 最大值

	// 循环获取所有文件列表
	for {
		reqBuilder := larkdrive.NewListFileReqBuilder()
		if pageToken != "" {
			reqBuilder.PageToken(pageToken)
		}
		req := reqBuilder.PageSize(pageSize).OrderBy(`EditedTime`).Direction(`DESC`).Build()
		resp, err := app.Drive.V1.File.List(context.Background(), req, larkcore.WithUserAccessToken(userAccessToken))
		if err != nil {
			logs.Error(fmt.Sprintf("获取文件列表失败: %v", err))
			return nil, err
		}
		if !resp.Success() {
			errMsg := fmt.Sprintf("logId: %s, error response: %s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
			logs.Error(errMsg)
			return nil, errors.New(errMsg)
		}
		if resp.Data != nil && resp.Data.Files != nil {
			for _, file := range resp.Data.Files {
				// 只获取 doc 和 docx 类型的文件
				if file.Type != nil && (*file.Type == "doc" || *file.Type == "docx") {
					docFile := FeishuDocFile{
						DocumentId: *file.Token,
					}
					if file.Name != nil {
						docFile.Name = *file.Name
					}
					allDocFiles = append(allDocFiles, docFile)
				}
			}
		}

		// 检查是否还有更多数据
		if resp.Data.HasMore == nil || !*resp.Data.HasMore {
			break
		}

		// 获取下一页的 pageToken
		if resp.Data.NextPageToken != nil {
			pageToken = *resp.Data.NextPageToken
		} else {
			break
		}
	}

	return allDocFiles, nil
}

// GetDocFileTree 获取飞书文档树形结构（包含文件夹和文档）
func (a *Application) GetDocFileTree(userAccessToken string, folderToken string) ([]*FeishuDocFileTree, error) {
	app, err := a.GetApp()
	if err != nil {
		return nil, err
	}

	var result []*FeishuDocFileTree
	pageToken := ""
	pageSize := 200

	// 循环获取指定文件夹下的所有文件
	for {
		reqBuilder := larkdrive.NewListFileReqBuilder()
		if pageToken != "" {
			reqBuilder.PageToken(pageToken)
		}

		// 如果指定了文件夹token，则获取该文件夹下的文件
		if folderToken != "" {
			reqBuilder.FolderToken(folderToken)
		}

		req := reqBuilder.PageSize(pageSize).OrderBy(`EditedTime`).Direction(`DESC`).Build()
		resp, err := app.Drive.V1.File.List(context.Background(), req, larkcore.WithUserAccessToken(userAccessToken))

		if err != nil {
			logs.Error(fmt.Sprintf("获取文件列表失败: %v", err))
			return nil, err
		}

		if !resp.Success() {
			errMsg := fmt.Sprintf("logId: %s, error response: %s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
			logs.Error(errMsg)
			return nil, errors.New(errMsg)
		}

		if resp.Data != nil && resp.Data.Files != nil {
			for _, file := range resp.Data.Files {
				if file.Token == nil || file.Type == nil {
					continue
				}

				node := &FeishuDocFileTree{
					Token: *file.Token,
					Type:  *file.Type,
				}

				if file.Name != nil {
					node.Name = *file.Name
				}
				if file.ParentToken != nil {
					node.ParentToken = *file.ParentToken
				}
				if file.CreatedTime != nil {
					node.CreatedTime = *file.CreatedTime
				}
				if file.ModifiedTime != nil {
					node.ModifiedTime = *file.ModifiedTime
				}
				if file.OwnerId != nil {
					node.OwnerID = *file.OwnerId
				}
				if file.Url != nil {
					node.URL = *file.Url
				}

				// 如果是文件夹，递归获取子文件
				if *file.Type == "folder" {
					children, err := a.GetDocFileTree(userAccessToken, *file.Token)
					if err != nil {
						logs.Error(fmt.Sprintf("获取文件夹 %s 子文件失败: %v", *file.Token, err))
						// 继续处理其他文件，不中断
					} else {
						node.Children = children
					}
				}

				result = append(result, node)
			}
		}

		// 检查是否还有更多数据
		if resp.Data.HasMore == nil || !*resp.Data.HasMore {
			break
		}

		// 获取下一页的 pageToken
		if resp.Data.NextPageToken != nil {
			pageToken = *resp.Data.NextPageToken
		} else {
			break
		}
	}

	return result, nil
}

func (a *Application) BuildUserAuthLoginUrl(redirectUri, frontRedirectUrl string) (string, error) {
	baseUrl := "https://accounts.feishu.cn/open-apis/authen/v1/authorize"
	state := tool.Base64Encode(tool.JsonEncodeNoError(FeishuUserAuthLoginState{
		AppId:            a.AppID,
		AppSecret:        a.Secret,
		FrontRedirectUrl: frontRedirectUrl,
	}))
	params := url.Values{}
	params.Add("client_id", a.AppID)
	params.Add("response_type", "code")
	params.Add("redirect_uri", redirectUri)
	params.Add("scope", "space:document:retrieve")
	params.Add("state", state)
	return fmt.Sprintf("%s?%s", baseUrl, params.Encode()), nil
}

func (a *Application) GetUserAccessToken(code, redirectUri string) (*FeishuUserAccessToken, error) {
	reqData := map[string]string{
		"grant_type":    "authorization_code",
		"code":          code,
		"client_id":     a.AppID,
		"client_secret": a.Secret,
		"redirect_uri":  redirectUri,
	}
	bodyBytes, _ := json.Marshal(reqData)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://open.feishu.cn/open-apis/authen/v2/oauth/token", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logs.Error(fmt.Sprintf("close body failed: %v", err))
		}
	}(resp.Body)

	var result FeishuUserAccessToken
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode json failed: %w", err)
	}

	// 飞书 OAuth 接口成功时 code 为 0
	if result.Code != 0 {
		return nil, fmt.Errorf("feishu error: %s (desc: %s, code: %d)",
			result.Error, result.ErrorDescription, result.Code)
	}

	return &result, nil
}

// GetDocFileDetail 获取飞书doc文档详情
func (a *Application) GetDocFileDetail(documentId string) (string, string, error) {
	app, err := a.GetApp()
	if err != nil {
		return "", "", err
	}

	// 先获取文档基本信息
	req1 := larkdocx.NewGetDocumentReqBuilder().DocumentId(documentId).Build()
	resp1, err := app.Docx.V1.Document.Get(context.Background(), req1)
	if err != nil {
		logs.Error(fmt.Sprintf("获取文档基本信息失败: %v", err))
		return "", "", err
	}
	if !resp1.Success() {
		errMsg := fmt.Sprintf("logId: %s, error response: %s", resp1.RequestId(), larkcore.Prettify(resp1.CodeError))
		logs.Error(errMsg)
		return "", "", errors.New(errMsg)
	}

	// 再获取文档内容
	req2 := larkdocx.NewRawContentDocumentReqBuilder().DocumentId(documentId).Build()
	resp2, err := app.Docx.V1.Document.RawContent(context.Background(), req2)
	if err != nil {
		logs.Error(fmt.Sprintf("获取文档内容失败: %v", err))
		return "", "", err
	}
	if !resp2.Success() {
		errMsg := fmt.Sprintf("logId: %s, error response: %s", resp2.RequestId(), larkcore.Prettify(resp2.CodeError))
		logs.Error(errMsg)
		return "", "", errors.New(errMsg)
	}
	if resp2.Data != nil && resp2.Data.Content != nil {
		return *resp1.Data.Document.Title, *resp2.Data.Content, nil
	}

	return "", "", nil
}
