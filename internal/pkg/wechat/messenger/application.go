// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package messenger

import (
	"bytes"
	"chatwiki/internal/pkg/lib_define"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/response"
	openresponse "github.com/ArtisanCloud/PowerWeChat/v3/src/openPlatform/authorizer/miniProgram/account/response"
	"github.com/zhimaAi/go_tools/curl"
)

var (
	ErrSecretInvalid  = errors.New("messenger_secret_error")
	ErrPageIDMismatch = errors.New("messenger_page_id_mismatch")
)

type Application struct {
	AppID        string
	Secret       string
	GraphAPIBase string
}

// SendText sends text message
func (a *Application) SendText(customer, content string, push *lib_define.PushMessage) (int, error) {
	url := fmt.Sprintf("%s/%s/messages?access_token=%s", a.GraphAPIBase, a.AppID, a.Secret)
	body := map[string]any{
		"recipient": map[string]any{"id": customer},
		"message":   map[string]any{"text": content},
	}
	resp := struct {
		Error struct {
			Message   string `json:"message"`
			Type      string `json:"type"`
			Code      int    `json:"code"`
			FbtraceId string `json:"fbtrace_id"`
		} `json:"error"`
	}{}
	request, err := curl.Post(url).JSONBody(body)
	if err != nil {
		return 0, err
	}
	if err = request.ToJSON(&resp); err != nil {
		return 0, err
	}
	if resp.Error.Code != 0 {
		return resp.Error.Code, errors.New(resp.Error.Message)
	}
	return 0, nil
}

// SendImage sends image message (via Attachment Upload API)
func (a *Application) SendImage(customer, filePath string, push *lib_define.PushMessage) (int, error) {
	attachmentId, errcode, err := a.UploadAttachment(filePath, "image")
	if err != nil {
		return errcode, err
	}
	return a.sendAttachment(customer, attachmentId, "image")
}

// GetToken validates Page Access Token and Page ID
func (a *Application) GetToken(refresh bool) (*response.ResponseGetToken, int, error) {
	url := fmt.Sprintf("%s/me?access_token=%s", a.GraphAPIBase, a.Secret)
	resp := struct {
		Name string `json:"name"`
		Id   string `json:"id"`
	}{}
	request := curl.Get(url)
	var err error
	if err = request.ToJSON(&resp); err != nil {
		return nil, 0, errors.New("messenger get token err:" + err.Error())
	}
	if len(resp.Id) == 0 {
		return nil, -1, ErrSecretInvalid
	}
	if resp.Id != a.AppID {
		return nil, -1, ErrPageIDMismatch
	}
	return &response.ResponseGetToken{}, 0, nil
}

// SendVideo sends video message (via Attachment Upload API)
func (a *Application) SendVideo(customer, filePath string, push *lib_define.PushMessage) (int, error) {
	attachmentId, errcode, err := a.UploadAttachment(filePath, "video")
	if err != nil {
		return errcode, err
	}
	return a.sendAttachment(customer, attachmentId, "video")
}

func (a *Application) SetTyping(customer, command string) (int, error) {
	return 0, errors.New(`not supported`)
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

func (a *Application) GetAccountBasicInfo() (*openresponse.ResponseGetBasicInfo, int, error) {
	return nil, 0, errors.New(`not supported`)
}

// SendVoice sends voice/audio message (via Attachment Upload API)
func (a *Application) SendVoice(customer, filePath string, push *lib_define.PushMessage) (int, error) {
	attachmentId, errcode, err := a.UploadAttachment(filePath, "audio")
	if err != nil {
		return errcode, err
	}
	return a.sendAttachment(customer, attachmentId, "audio")
}

// UploadAttachment uploads a local file to Messenger Attachment Upload API, returns attachment_id
func (a *Application) UploadAttachment(filePath, mediaType string) (string, int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", 0, fmt.Errorf("open file error: %w", err)
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// message field: JSON with attachment type
	messageData := map[string]any{
		"attachment": map[string]any{
			"type":    mediaType,
			"payload": map[string]any{"is_reusable": true},
		},
	}
	messageJSON, _ := json.Marshal(messageData)
	_ = writer.WriteField("message", string(messageJSON))

	// filedata field: the actual file
	part, err := writer.CreateFormFile("filedata", filepath.Base(filePath))
	if err != nil {
		return "", 0, fmt.Errorf("create form file error: %w", err)
	}
	if _, err = io.Copy(part, file); err != nil {
		return "", 0, fmt.Errorf("copy file data error: %w", err)
	}
	if err = writer.Close(); err != nil {
		return "", 0, fmt.Errorf("close writer error: %w", err)
	}

	url := fmt.Sprintf("%s/%s/message_attachments?access_token=%s", a.GraphAPIBase, a.AppID, a.Secret)
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	result := struct {
		AttachmentId string `json:"attachment_id"`
		Error        struct {
			Message   string `json:"message"`
			Type      string `json:"type"`
			Code      int    `json:"code"`
			FbtraceId string `json:"fbtrace_id"`
		} `json:"error"`
	}{}
	if err = json.Unmarshal(respBody, &result); err != nil {
		return "", 0, err
	}
	if result.Error.Code != 0 {
		return "", result.Error.Code, errors.New(result.Error.Message)
	}
	return result.AttachmentId, 0, nil
}

// sendAttachment sends a message with attachment_id
func (a *Application) sendAttachment(customer, attachmentId, mediaType string) (int, error) {
	url := fmt.Sprintf("%s/%s/messages?access_token=%s", a.GraphAPIBase, a.AppID, a.Secret)
	body := map[string]any{
		"recipient": map[string]any{"id": customer},
		"message": map[string]any{
			"attachment": map[string]any{
				"type":    mediaType,
				"payload": map[string]any{"attachment_id": attachmentId},
			},
		},
	}
	resp := struct {
		Error struct {
			Message   string `json:"message"`
			Type      string `json:"type"`
			Code      int    `json:"code"`
			FbtraceId string `json:"fbtrace_id"`
		} `json:"error"`
	}{}
	request, err := curl.Post(url).JSONBody(body)
	if err != nil {
		return 0, err
	}
	if err = request.ToJSON(&resp); err != nil {
		return 0, err
	}
	if resp.Error.Code != 0 {
		return resp.Error.Code, errors.New(resp.Error.Message)
	}
	return 0, nil
}
