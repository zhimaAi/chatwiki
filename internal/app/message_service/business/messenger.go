// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package business

import (
	"chatwiki/internal/app/message_service/common"
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
	"github.com/zhimaAi/go_tools/tool"
)

func MessengerPush(c *gin.Context) {
	accessKey := strings.TrimSpace(c.Param(`access_key`))
	appInfo, err := common.GetWechatAppInfo(`access_key`, accessKey)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}
	if len(appInfo) == 0 || appInfo[`app_type`] != lib_define.AppMessenger {
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}
	// GET: Webhook verification
	if c.Request.Method == http.MethodGet {
		mode := c.Query(`hub_mode`)
		//verifyToken := c.Query(`hub.verify_token`)
		challenge := c.Query(`hub_challenge`)
		//ignore verifyToken == accessKey
		if mode == "subscribe" {
			c.String(http.StatusOK, challenge)
			return
		}
		c.String(http.StatusOK, "verification failed")
		return
	}

	// POST: message receiving
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logs.Error(err.Error())
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}
	if len(body) == 0 || len(strings.TrimSpace(string(body))) == 0 {
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}
	var webhookData lib_define.MessengerWebhookEvent
	if err = tool.JsonDecodeUseNumber(string(body), &webhookData); err != nil {
		logs.Error("messenger webhook decode err:%s, body:%s", err.Error(), string(body))
		c.String(http.StatusOK, lib_define.SUCCESS)
		return
	}

	adminUserId := cast.ToInt(appInfo[`admin_user_id`])
	appId := appInfo[`app_id`]
	for _, entry := range webhookData.Entry {
		pageID := entry.ID
		for _, msg := range entry.Messaging {
			go func(msg lib_define.MessengerMessaging) {
				nsqMsg := buildMessengerNSQMsg(appId, pageID, msg, adminUserId)
				if nsqMsg == nil {
					return
				}
				common.PushNSQ(nsqMsg)
			}(msg)
		}
	}
	c.String(http.StatusOK, lib_define.SUCCESS)
}

// buildMessengerNSQMsg builds NSQ message
func buildMessengerNSQMsg(appId, pageID string, msg lib_define.MessengerMessaging, adminUserId int) map[string]any {
	// process message content
	msgType, content, ossUrl := parseMessengerMessage(msg, adminUserId)
	if len(msgType) == 0 {
		return nil
	}

	nsqMsg := map[string]any{
		"appid":        appId,
		"ToUserName":   pageID,
		"FromUserName": msg.Sender.ID,
		"CreateTime":   msg.Timestamp,
		"MsgId":        msg.Message.MID,
		"MsgType":      msgType,
		"Content":      content,
	}
	if len(ossUrl) > 0 {
		nsqMsg[`oss_url`] = ossUrl
	}

	return nsqMsg
}

// parseMessengerMessage parses Messenger message, returns MsgType, Content and oss_url
func parseMessengerMessage(msg lib_define.MessengerMessaging, adminUserId int) (msgType, content, ossUrl string) {
	if len(msg.Message.Text) > 0 {
		// text message
		msgType = lib_define.MessengerMsgTypeText
		content = msg.Message.Text
		return
	}

	if len(msg.Message.Attachments) > 0 {
		// prioritize first image attachment (stickers treated as images)
		for _, att := range msg.Message.Attachments {
			switch att.Type {
			case lib_define.MessengerMsgTypeImage:
				msgType = lib_define.MessengerMsgTypeImage
				ossUrl = downloadMessengerMedia(att.Payload.URL, adminUserId)
				return
			case lib_define.MessengerMsgTypeAudio:
				msgType = lib_define.MessengerMsgTypeAudio
				ossUrl = downloadMessengerMedia(att.Payload.URL, adminUserId)
				return
			case lib_define.MessengerMsgTypeVideo:
				msgType = lib_define.MessengerMsgTypeVideo
				ossUrl = downloadMessengerMedia(att.Payload.URL, adminUserId)
				return
			case lib_define.MessengerMsgTypeSticker:
				msgType = lib_define.MessengerMsgTypeImage
				ossUrl = downloadMessengerMedia(att.Payload.URL, adminUserId)
				return
			}
		}
	}

	return
}

// downloadMessengerMedia downloads Messenger media file and saves locally
// fallback to original CDN URL on download failure to avoid data loss
func downloadMessengerMedia(mediaUrl string, adminUserId int) string {
	if len(mediaUrl) == 0 || !common.IsUrl(mediaUrl) {
		return ``
	}
	// extract path from URL (remove query params), then get extension with filepath.Ext
	pathPart := mediaUrl
	if idx := strings.Index(mediaUrl, `?`); idx >= 0 {
		pathPart = mediaUrl[:idx]
	}
	ext := strings.ToLower(filepath.Ext(pathPart))
	if len(ext) == 0 {
		ext = ".jpg"
	}
	// consistent storage path with WeCom robot: upload/chat_ai/{userId}/received_message_images/{Ym}/{md5}{ext}
	// served directly by Go gin StaticFS("/upload", ...), bypassing nginx
	objectKey := fmt.Sprintf("chat_ai/%d/received_message_images/%s/%s%s", adminUserId, tool.Date(`Ym`), tool.MD5(mediaUrl), ext)
	localFile := `internal/app/chatwiki/upload/` + objectKey
	if tool.IsFile(localFile) {
		return `【image_domain】/upload/` + objectKey
	}
	// ensure directory exists
	if err := os.MkdirAll(filepath.Dir(localFile), 0755); err != nil {
		logs.Error("create messenger media dir err:%s, path:%s", err.Error(), filepath.Dir(localFile))
		return mediaUrl
	}
	if err := downloadMediaFile(mediaUrl, localFile, 40*time.Second); err != nil {
		logs.Error("download messenger media err:%s, url:%s, fallback to original CDN url", err.Error(), mediaUrl)
		return mediaUrl
	}
	return `【image_domain】/upload/` + objectKey
}

// downloadMediaFile downloads media file with custom timeout (includes body read)
func downloadMediaFile(fileUrl, temFile string, timeout time.Duration) error {
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
