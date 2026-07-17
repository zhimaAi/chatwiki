// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package business

import (
	"chatwiki/internal/app/message_service/common"
	"chatwiki/internal/app/message_service/define"
	"chatwiki/internal/pkg/lib_define"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/curl"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

// TelegramUpdate represents the Update object from Telegram Bot API.
// Reference: https://core.telegram.org/bots/api#update
type TelegramUpdate struct {
	UpdateId int64             `json:"update_id"`
	Message  *TelegramMessage  `json:"message"`
	Callback *TelegramCallback `json:"callback_query"`
}

// TelegramMessage represents the Message object from Telegram Bot API.
type TelegramMessage struct {
	MessageId int64             `json:"message_id"`
	From      *TelegramUser     `json:"from"`
	Chat      *TelegramChat     `json:"chat"`
	Date      int64             `json:"date"`
	Text      string            `json:"text"`
	Photo     []TelegramPhoto   `json:"photo"`
	Document  *TelegramDocument `json:"document"`
	Voice     *TelegramVoice    `json:"voice"`
	Video     *TelegramVideo    `json:"video"`
	Audio     *TelegramAudio    `json:"audio"`
	Caption   string            `json:"caption"`
}

// TelegramUser represents the User object from Telegram Bot API.
type TelegramUser struct {
	Id        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

// TelegramChat represents the Chat object from Telegram Bot API.
type TelegramChat struct {
	Id        int64  `json:"id"`
	Type      string `json:"type"` // "private", "group", "supergroup", "channel"
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Title     string `json:"title"`
}

// TelegramPhoto represents a PhotoSize object from Telegram Bot API.
type TelegramPhoto struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int    `json:"file_size"`
}

// TelegramDocument represents the Document object from Telegram Bot API.
type TelegramDocument struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	FileName     string `json:"file_name"`
	MimeType     string `json:"mime_type"`
	FileSize     int64  `json:"file_size"`
}

// TelegramVoice represents the Voice object from Telegram Bot API.
type TelegramVoice struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	Duration     int    `json:"duration"`
	MimeType     string `json:"mime_type"`
	FileSize     int64  `json:"file_size"`
}

// TelegramVideo represents the Video object from Telegram Bot API.
type TelegramVideo struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Duration     int    `json:"duration"`
	FileName     string `json:"file_name"`
	MimeType     string `json:"mime_type"`
	FileSize     int64  `json:"file_size"`
}

// TelegramAudio represents the Audio object from Telegram Bot API.
type TelegramAudio struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	Duration     int    `json:"duration"`
	FileName     string `json:"file_name"`
	MimeType     string `json:"mime_type"`
	FileSize     int64  `json:"file_size"`
}

// TelegramCallback represents the CallbackQuery object from Telegram Bot API.
type TelegramCallback struct {
	Id      string           `json:"id"`
	From    *TelegramUser    `json:"from"`
	Message *TelegramMessage `json:"message"`
	Data    string           `json:"data"`
}

// TelegramRobot handles incoming Telegram webhook updates.
func TelegramRobot(c *gin.Context) {
	accessKey := strings.TrimSpace(c.Param(`access_key`))
	if len(accessKey) == 0 {
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}

	// Look up app info by access_key
	appInfo, err := common.GetWechatAppInfo(`access_key`, accessKey)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}
	if len(appInfo) == 0 || appInfo[`app_type`] != lib_define.TelegramRobot {
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}

	adminUserId := cast.ToInt(appInfo[`admin_user_id`])

	// Verify secret_token in X-Telegram-Bot-Api-Secret-Token header
	secretToken := c.GetHeader(`X-Telegram-Bot-Api-Secret-Token`)
	expectedSecret := cast.ToString(appInfo[`app_secret`])
	if len(expectedSecret) > 0 && secretToken != expectedSecret {
		logs.Error(`telegram secret_token mismatch for access_key:%s`, accessKey)
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}

	// Read and parse the request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}
	if len(body) == 0 || len(strings.TrimSpace(string(body))) == 0 {
		c.String(http.StatusOK, lib_define.SUCCESS)
		logs.Notice(strings.Join([]string{c.ClientIP(), c.Request.Method, c.Request.Host, c.Request.RequestURI}, ` | `))
		return
	}

	update := TelegramUpdate{}
	if err = tool.JsonDecodeUseNumber(string(body), &update); err != nil {
		logs.Error(`telegram parse update error:%s, body:%s`, err.Error(), string(body))
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}

	// Handle callback_query
	if update.Callback != nil {
		handleTelegramCallback(appInfo, update.Callback)
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}

	// Handle message
	if update.Message == nil {
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}

	go handleTelegramMessage(appInfo, update.Message, adminUserId)
	c.String(http.StatusOK, lib_define.SUCCESS)
}

// handleTelegramMessage processes a Telegram Message and pushes it to NSQ.
func handleTelegramMessage(appInfo msql.Params, msg *TelegramMessage, adminUserId int) {
	if msg.From == nil || msg.Chat == nil {
		logs.Error(`telegram message missing from or chat`)
		return
	}

	// Build display nickname
	nickname := buildTelegramNickname(msg.From)

	// Determine message type and content
	msgType, content := extractTelegramContent(msg)
	botToken := cast.ToString(appInfo[`app_id`])

	nsqMsg := map[string]any{
		`ToUserName`:   cast.ToString(msg.Chat.Id),
		`FromUserName`: cast.ToString(msg.From.Id),
		`CreateTime`:   cast.ToString(msg.Date),
		`MsgType`:      msgType,
		`Content`:      content,
		`MsgId`:        cast.ToString(msg.MessageId),
		`appid`:        appInfo[`app_id`],
		`nickname`:     nickname,
		`chattype`:     msg.Chat.Type, // "private", "group", "supergroup", "channel"
	}

	// Download media files at webhook time (align with Messenger flow)
	if tool.InArray(msgType, []string{`image`, `voice`, `video`}) {
		if len(content) > 0 {
			ossUrl := downloadTelegramMedia(content, botToken, adminUserId)
			if len(ossUrl) > 0 {
				nsqMsg[`oss_url`] = ossUrl
			}
		}
		nsqMsg[`MediaId`] = content // keep file_id for fallback
	}

	// Set SessionType for group chat detection (used in reply flow)
	if msg.Chat.Type == `group` || msg.Chat.Type == `supergroup` {
		nsqMsg[`SessionType`] = `group`
	} else {
		nsqMsg[`SessionType`] = `single`
	}

	logs.Info(`telegram handleMessage appid:%s, msgType:%s, from:%s, chat:%s, content:%s`,
		appInfo[`app_id`], msgType, cast.ToString(msg.From.Id), cast.ToString(msg.Chat.Id), content)
	go common.PushNSQ(nsqMsg)
}

// handleTelegramCallback processes a Telegram CallbackQuery.
func handleTelegramCallback(appInfo msql.Params, callback *TelegramCallback) {
	if callback.From == nil {
		return
	}
	chatId := ``
	if callback.Message != nil && callback.Message.Chat != nil {
		chatId = cast.ToString(callback.Message.Chat.Id)
	}
	nsqMsg := map[string]any{
		`ToUserName`:   chatId,
		`FromUserName`: cast.ToString(callback.From.Id),
		`CreateTime`:   `0`,
		`MsgType`:      `callback_query`,
		`Content`:      callback.Data,
		`MsgId`:        callback.Id,
		`appid`:        appInfo[`app_id`],
		`nickname`:     buildTelegramNickname(callback.From),
		`SessionType`:  `single`,
	}
	logs.Info(`telegram handleCallback appid:%s, from:%s, chat:%s, data:%s`,
		appInfo[`app_id`], cast.ToString(callback.From.Id), chatId, callback.Data)
	go common.PushNSQ(nsqMsg)
}

// buildTelegramNickname constructs a display name from Telegram user info.
func buildTelegramNickname(user *TelegramUser) string {
	name := strings.TrimSpace(user.FirstName + ` ` + user.LastName)
	if len(name) == 0 {
		name = user.Username
	}
	if len(name) == 0 {
		name = cast.ToString(user.Id)
	}
	return name
}

// extractTelegramContent determines the message type and extracts content from a Telegram Message.
func extractTelegramContent(msg *TelegramMessage) (msgType, content string) {
	if len(msg.Text) > 0 {
		return `text`, msg.Text
	}
	if len(msg.Caption) > 0 {
		content = msg.Caption
	}
	if len(msg.Photo) > 0 {
		// Get the largest photo (last in the array per Telegram API)
		largestPhoto := msg.Photo[len(msg.Photo)-1]
		return `image`, largestPhoto.FileId
	}
	if msg.Document != nil {
		if len(content) == 0 {
			content = msg.Document.FileName
		}
		return `document`, msg.Document.FileId
	}
	if msg.Voice != nil {
		return `voice`, msg.Voice.FileId
	}
	if msg.Video != nil {
		if len(content) == 0 {
			content = msg.Video.FileName
		}
		return `video`, msg.Video.FileId
	}
	if msg.Audio != nil {
		if len(content) == 0 {
			content = msg.Audio.FileName
		}
		return `audio`, msg.Audio.FileId
	}
	return `unknown`, ``
}

// downloadTelegramMedia downloads a Telegram media file (image/voice/video) by file_id
// and saves it locally, returning the oss_url. Aligns with Messenger's downloadMessengerMedia flow.
func downloadTelegramMedia(fileId, botToken string, adminUserId int) string {
	if len(fileId) == 0 || len(botToken) == 0 {
		return ``
	}

	apiBase := define.Config.Telegram[`api_base`]
	if len(apiBase) == 0 {
		apiBase = `https://api.telegram.org`
	}

	// Step 1: Resolve file_id → file_path via getFile API
	getFileUrl := apiBase + `/bot` + botToken + `/getFile?file_id=` + fileId
	gfResp, err := curl.Get(getFileUrl).Response()
	if err != nil {
		logs.Error(`telegram getFile request error:%s, file_id:%s`, err.Error(), fileId)
		return ``
	}
	defer gfResp.Body.Close()
	gfBody, err := io.ReadAll(gfResp.Body)
	if err != nil {
		logs.Error(`telegram getFile read error:%s, file_id:%s`, err.Error(), fileId)
		return ``
	}

	getFileResult := struct {
		Ok     bool `json:"ok"`
		Result struct {
			FilePath string `json:"file_path"`
		} `json:"result"`
	}{}
	if err = tool.JsonDecode(string(gfBody), &getFileResult); err != nil || !getFileResult.Ok || len(getFileResult.Result.FilePath) == 0 {
		logs.Error(`telegram getFile failed:%s, file_id:%s`, string(gfBody), fileId)
		return ``
	}

	filePath := getFileResult.Result.FilePath

	// Step 2: Determine extension and local save path (same convention as Messenger)
	ext := strings.ToLower(filepath.Ext(filePath))
	if len(ext) == 0 {
		ext = ".jpg"
	}
	objectKey := fmt.Sprintf("chat_ai/%d/received_message_images/%s/%s%s", adminUserId, tool.Date(`Ym`), tool.MD5(fileId), ext)
	localFile := `internal/app/chatwiki/upload/` + objectKey
	if tool.IsFile(localFile) {
		return `【image_domain】/upload/` + objectKey
	}

	if err := os.MkdirAll(filepath.Dir(localFile), 0755); err != nil {
		logs.Error("create telegram media dir err:%s, path:%s", err.Error(), filepath.Dir(localFile))
		return ``
	}

	// Step 3: Download the actual file from Telegram's file server
	downloadUrl := apiBase + `/file/bot` + botToken + `/` + filePath
	if err := downloadTelegramFile(downloadUrl, localFile, 40*time.Second); err != nil {
		logs.Error("download telegram media err:%s, file_id:%s", err.Error(), fileId)
		return ``
	}

	return `【image_domain】/upload/` + objectKey
}

// downloadTelegramFile downloads a file from URL with overall timeout (includes body read).
func downloadTelegramFile(fileUrl, temFile string, timeout time.Duration) error {
	request := curl.Get(fileUrl).SetTimeout(timeout, timeout)
	response, err := request.Response()
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf(`StatusCode:%d`, response.StatusCode))
	}

	// use goroutine + timer to enforce total timeout including response body read,
	// since http.Client.Timeout only covers headers, not body streaming
	errCh := make(chan error, 1)
	go func() {
		errCh <- request.ToFile(temFile)
	}()

	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case err = <-errCh:
	case <-timer.C:
		err = errors.New(fmt.Sprintf(`download timeout after %s`, timeout))
	}

	if err != nil && tool.IsFile(temFile) {
		_ = os.Remove(temFile)
	}
	return err
}
