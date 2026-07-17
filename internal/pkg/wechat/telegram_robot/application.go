// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package telegram_robot

import (
	"bytes"
	"chatwiki/internal/pkg/lib_define"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/response"
	openresponse "github.com/ArtisanCloud/PowerWeChat/v3/src/openPlatform/authorizer/miniProgram/account/response"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/curl"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

// TelegramApiBase Telegram Bot API base URL, can be overridden via config
var TelegramApiBase = "https://api.telegram.org"

type Application struct {
	AppID  string // bot token
	Secret string // HMAC secret_token for webhook verification
}

// TelegramResponse is the generic response from Telegram Bot API
type TelegramResponse struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
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

// SendText sends a text message via Telegram Bot API.
// push.Message["ToUserName"] is expected to contain the chat_id.
func (a *Application) SendText(customer, content string, push *lib_define.PushMessage) (int, error) {
	chatId := cast.ToString(push.Message[`ToUserName`])
	if len(chatId) == 0 {
		return 0, errors.New(`chat_id is empty`)
	}
	if len(a.AppID) == 0 {
		return 0, errors.New(`bot token is empty`)
	}
	body := map[string]any{
		`chat_id`:    chatId,
		`text`:       content,
		`parse_mode`: `HTML`,
	}
	url := TelegramApiBase + `/bot` + a.AppID + `/sendMessage`
	request, err := curl.Post(url).JSONBody(body)
	if err != nil {
		return 0, err
	}
	resp := TelegramResponse{}
	if err = request.ToJSON(&resp); err != nil {
		return 0, err
	}
	if !resp.Ok {
		logs.Error(`telegram sendMessage failed: code=%d, desc=%s`, resp.ErrorCode, resp.Description)
		return resp.ErrorCode, errors.New(resp.Description)
	}
	return 0, nil
}

// sendTelegramFile sends a file (photo/video/voice) to a Telegram chat via multipart/form-data.
// apiMethod should be one of "sendPhoto", "sendVideo", "sendVoice".
func (a *Application) sendTelegramFile(chatId, filePath, apiMethod string) (int, error) {
	if len(chatId) == 0 {
		return 0, errors.New(`chat_id is empty`)
	}
	if len(a.AppID) == 0 {
		return 0, errors.New(`bot token is empty`)
	}

	file, err := os.Open(filePath)
	if err != nil {
		logs.Error(`telegram %s open file error:%s, path:%s`, apiMethod, err.Error(), filePath)
		return 0, err
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	_ = writer.WriteField(`chat_id`, chatId)

	fileField := `photo`
	if apiMethod == `sendVideo` {
		fileField = `video`
	} else if apiMethod == `sendVoice` {
		fileField = `voice`
	}

	part, err := writer.CreateFormFile(fileField, filepath.Base(filePath))
	if err != nil {
		logs.Error(`telegram %s create form file error:%s`, apiMethod, err.Error())
		return 0, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		logs.Error(`telegram %s copy file data error:%s`, apiMethod, err.Error())
		return 0, err
	}

	err = writer.Close()
	if err != nil {
		logs.Error(`telegram %s close writer error:%s`, apiMethod, err.Error())
		return 0, err
	}

	url := TelegramApiBase + `/bot` + a.AppID + `/` + apiMethod
	req, err := http.NewRequest(`POST`, url, &requestBody)
	if err != nil {
		logs.Error(`telegram %s create request error:%s`, apiMethod, err.Error())
		return 0, err
	}
	req.Header.Set(`Content-Type`, writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logs.Error(`telegram %s request error:%s`, apiMethod, err.Error())
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logs.Error(`telegram %s read response error:%s`, apiMethod, err.Error())
		return 0, err
	}

	tgResp := TelegramResponse{}
	if err = tool.JsonDecode(string(body), &tgResp); err != nil {
		logs.Error(`telegram %s parse response error:%s, body:%s`, apiMethod, err.Error(), string(body))
		return 0, err
	}
	if !tgResp.Ok {
		logs.Error(`telegram %s failed: code=%d, desc=%s`, apiMethod, tgResp.ErrorCode, tgResp.Description)
		return tgResp.ErrorCode, errors.New(tgResp.Description)
	}

	logs.Info(`telegram %s success, chatId:%s`, apiMethod, chatId)
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

// SendImage sends a photo via Telegram sendPhoto API.
func (a *Application) SendImage(customer, filePath string, push *lib_define.PushMessage) (int, error) {
	chatId := cast.ToString(push.Message[`ToUserName`])
	return a.sendTelegramFile(chatId, filePath, `sendPhoto`)
}

// GetFileByMedia downloads a file from Telegram by file_id.
// It first calls getFile API to resolve the file_id to a file_path,
// then downloads the actual file from Telegram's file server.
func (a *Application) GetFileByMedia(mediaId string, push *lib_define.PushMessage) ([]byte, http.Header, int, error) {
	if len(mediaId) == 0 {
		return nil, nil, 0, errors.New(`file_id is empty`)
	}
	if len(a.AppID) == 0 {
		return nil, nil, 0, errors.New(`bot token is empty`)
	}

	// Step 1: Call getFile API to resolve file_id → file_path
	getFileUrl := TelegramApiBase + `/bot` + a.AppID + `/getFile?file_id=` + mediaId
	getFileResp := struct {
		Ok     bool `json:"ok"`
		Result struct {
			FilePath string `json:"file_path"`
		} `json:"result"`
		ErrorCode   int    `json:"error_code"`
		Description string `json:"description"`
	}{}
	gfResp, err := http.Get(getFileUrl)
	if err != nil {
		logs.Error(`telegram getFile request error:%s, file_id:%s`, err.Error(), mediaId)
		return nil, nil, 0, err
	}
	defer gfResp.Body.Close()
	gfBody, err := io.ReadAll(gfResp.Body)
	if err != nil {
		logs.Error(`telegram getFile read error:%s, file_id:%s`, err.Error(), mediaId)
		return nil, nil, 0, err
	}
	if err = tool.JsonDecode(string(gfBody), &getFileResp); err != nil {
		logs.Error(`telegram getFile parse error:%s, file_id:%s`, err.Error(), mediaId)
		return nil, nil, 0, err
	}
	if !getFileResp.Ok {
		logs.Error(`telegram getFile failed: code=%d, desc=%s, file_id:%s`, getFileResp.ErrorCode, getFileResp.Description, mediaId)
		return nil, nil, getFileResp.ErrorCode, errors.New(getFileResp.Description)
	}
	if len(getFileResp.Result.FilePath) == 0 {
		return nil, nil, 0, errors.New(`file_path is empty`)
	}

	// Step 2: Download the actual file
	downloadUrl := TelegramApiBase + `/file/bot` + a.AppID + `/` + getFileResp.Result.FilePath
	resp, err := http.Get(downloadUrl)
	if err != nil {
		logs.Error(`telegram download file error:%s, url:%s`, err.Error(), downloadUrl)
		return nil, nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, resp.StatusCode, fmt.Errorf(`download failed with status %d`, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logs.Error(`telegram read file body error:%s, file_id:%s`, err.Error(), mediaId)
		return nil, nil, 0, err
	}

	// Fix Content-Type: Telegram file server may return application/octet-stream,
	// use file_path extension to infer the correct MIME type for downstream processing.
	contentType := resp.Header.Get(`Content-Type`)
	if contentType == `` || contentType == `application/octet-stream` {
		ext := filepath.Ext(getFileResp.Result.FilePath)
		if fixedType := extToContentType(ext); len(fixedType) > 0 {
			resp.Header.Set(`Content-Type`, fixedType)
			logs.Info(`telegram GetFileByMedia fixed Content-Type: %s → %s (ext:%s, file_path:%s)`,
				contentType, fixedType, ext, getFileResp.Result.FilePath)
		}
	}

	logs.Info(`telegram GetFileByMedia success, file_id:%s, size:%d`, mediaId, len(body))
	return body, resp.Header, 0, nil
}

// extToContentType maps file extension to MIME Content-Type.
func extToContentType(ext string) string {
	switch strings.ToLower(ext) {
	case `.jpg`, `.jpeg`:
		return `image/jpeg`
	case `.png`:
		return `image/png`
	case `.gif`:
		return `image/gif`
	case `.webp`:
		return `image/webp`
	case `.ogg`, `.oga`:
		return `audio/ogg`
	case `.mp3`:
		return `audio/mpeg`
	case `.mp4`, `.m4v`:
		return `video/mp4`
	case `.pdf`:
		return `application/pdf`
	}
	return ``
}

func (a *Application) GetAccountBasicInfo() (*openresponse.ResponseGetBasicInfo, int, error) {
	return nil, 0, errors.New(`not supported`)
}

// SendVoice sends a voice message via Telegram sendVoice API.
func (a *Application) SendVoice(customer, filePath string, push *lib_define.PushMessage) (int, error) {
	chatId := cast.ToString(push.Message[`ToUserName`])
	return a.sendTelegramFile(chatId, filePath, `sendVoice`)
}

// SendVideo sends a video via Telegram sendVideo API.
func (a *Application) SendVideo(customer, filePath string, push *lib_define.PushMessage) (int, error) {
	chatId := cast.ToString(push.Message[`ToUserName`])
	return a.sendTelegramFile(chatId, filePath, `sendVideo`)
}
