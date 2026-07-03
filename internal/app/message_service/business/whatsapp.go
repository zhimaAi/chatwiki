// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package business

import (
	"chatwiki/internal/app/message_service/common"
	"chatwiki/internal/app/message_service/define"
	"chatwiki/internal/pkg/lib_define"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
)

// whatsappInboundMsg represents a single message in the Aliyun ChatApp inbound batch.
// Timestamp is declared as json.Number because Aliyun sends a numeric millisecond
// timestamp; decoding it into a Go string via tool.JsonDecodeUseNumber would fail
// ("cannot unmarshal number into ... of type string") and drop the whole batch.
type whatsappInboundMsg struct {
	MessageId   string      `json:"MessageId"`
	ChannelType string      `json:"ChannelType"`
	From        string      `json:"From"`
	To          string      `json:"To"`
	Type        string      `json:"Type"`
	Message     string      `json:"Message"`
	Timestamp   json.Number `json:"Timestamp"`
}

// whatsappMediaMessage is the embedded JSON in Message for non-text types.
type whatsappMediaMessage struct {
	Url      string `json:"url"`
	MimeType string `json:"mimeType"`
	Id       string `json:"id"`
	Filename string `json:"filename"`
}

// getWhatsappFileByUrl downloads a plain-HTTP(S) media URL and stores it under the
// shared upload directory (internal/app/chatwiki/upload/), returning the
// 【image_domain】-prefixed link. This mirrors the existing GetWecomFileByUrl flow in
// the same service: message_service has no OSS subsystem of its own, so it writes to
// the upload volume that the chatwiki core serves via image_domain. The chatwiki-core
// OSS-aware writer (WriteFileByString) lives in a different binary's package and is not
// reachable here without a cross-service import. No AES decrypt step is needed because
// Aliyun ChatApp media URLs are not encrypted (unlike WeCom).
// whatsappFileExt derives a clean file extension for a downloaded media file.
// Aliyun ChatApp media URLs carry an OSS signature query
// (e.g. https://.../xxx.jpg?expires=...&signature=...); calling filepath.Ext on the
// raw URL would yield ".jpg?expires=...&signature=..." and bake the query into the
// stored filename — breaking the served Content-Type and tripping multimodal models'
// "Unsupported image format". So strip the query via url.Parse first, then fall back
// to the declared mimeType (".bin" would fail the same way for vision models).
func whatsappFileExt(mediaUrl, mimeType string) string {
	if u, err := url.Parse(mediaUrl); err == nil {
		if ext := strings.ToLower(filepath.Ext(u.Path)); ext != `` && len(ext) <= 6 {
			return ext
		}
	}
	if mimeType != `` {
		if exts, err := mime.ExtensionsByType(mimeType); err == nil && len(exts) > 0 {
			return strings.ToLower(exts[0])
		}
	}
	return `.bin`
}

func getWhatsappFileByUrl(mediaUrl, mimeType string, adminUserId int) (string, string) {
	// Retry up to 3 times with 30s timeout each: Aliyun ChatApp media is served from
	// oss-ap-southeast-1 and can be slow; a single 60s attempt frequently hits
	// "context deadline exceeded". 30s per attempt keeps total under 90s.
	var bs []byte
	for attempt := 1; attempt <= 3; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, mediaUrl, nil)
		if err != nil {
			cancel()
			logs.Error(`whatsapp: build request failed url:%s err:%s`, mediaUrl, err.Error())
			return ``, ``
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			cancel()
			logs.Warning(`whatsapp: download failed attempt:%d url:%s err:%s`, attempt, mediaUrl, err.Error())
			continue
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			cancel()
			logs.Warning(`whatsapp: download bad status attempt:%d url:%s status:%d`, attempt, mediaUrl, resp.StatusCode)
			continue
		}
		bs, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		cancel()
		if err != nil || len(bs) == 0 {
			logs.Warning(`whatsapp: read body failed attempt:%d url:%s err:%v`, attempt, mediaUrl, err)
			bs = nil
			continue
		}
		break
	}
	if len(bs) == 0 {
		logs.Error(`whatsapp: download failed after 3 attempts url:%s`, mediaUrl)
		return ``, ``
	}

	data := string(bs)
	ext := whatsappFileExt(mediaUrl, mimeType)
	objectKey := fmt.Sprintf(`chat_ai/%d/received_message_images/%s/%s%s`,
		adminUserId, tool.Date(`Ym`), tool.MD5(data), ext)

	if err := tool.WriteFile(`internal/app/chatwiki/upload/`+objectKey, data); err != nil {
		logs.Error(`whatsapp: write file failed url:%s err:%s`, mediaUrl, err.Error())
		return ``, ``
	}
	return `【image_domain】/upload/` + objectKey, ``
}

// normalizeWhatsappPhone strips a leading '+' so numbers from Aliyun inbound
// (which may carry '+', e.g. "+8613800000000") compare equal to the stored
// app_id (saved without '+').
func normalizeWhatsappPhone(phone string) string {
	return strings.TrimPrefix(strings.TrimSpace(phone), `+`)
}

// WhatsappPush handles inbound WhatsApp callbacks from Aliyun ChatApp.
// The payload is a JSON array; each element is one inbound message.
func WhatsappPush(c *gin.Context) {
	accessKey := strings.TrimSpace(c.Param(`access_key`))

	// GET requests are health-check / verification probes — respond 200 immediately.
	if c.Request.Method == http.MethodGet {
		c.JSON(http.StatusOK, gin.H{`code`: 0, `msg`: `success`})
		return
	}

	appInfo, err := common.GetWechatAppInfo(`access_key`, accessKey)
	if err != nil {
		logs.Error(`whatsapp: get app info failed access_key:%s err:%s`, accessKey, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{`code`: 1, `msg`: err.Error()})
		return
	}
	// access_key is now channel(CustSpaceId)-level: multiple business numbers under the same
	// channel share one access_key. The record matched by access_key here is only used to verify
	// "the channel exists and is of type whatsapp", not to determine message ownership; which
	// number/robot each message belongs to is routed precisely by the To field in the loop below.
	if len(appInfo) == 0 || appInfo[`app_type`] != lib_define.AppWhatsapp {
		logs.Error(`whatsapp: app not found or type mismatch access_key:%s`, accessKey)
		c.JSON(http.StatusOK, gin.H{`code`: 1, `msg`: `app not found`})
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logs.Error(`whatsapp: read body failed err:%s`, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{`code`: 1, `msg`: err.Error()})
		return
	}

	var msgs []whatsappInboundMsg
	if err = tool.JsonDecodeUseNumber(string(body), &msgs); err != nil {
		logs.Error(`whatsapp: parse body failed body:%s err:%s`, string(body), err.Error())
		c.JSON(http.StatusOK, gin.H{`code`: 1, `msg`: `invalid body`})
		return
	}
	if len(msgs) == 0 {
		c.JSON(http.StatusOK, gin.H{`code`: 0, `msg`: `success`})
		return
	}

	for _, item := range msgs {
		// To is the message recipient (this system's business number). access_key is channel-level
		// and shared across numbers in the same channel, so we must locate the specific number record
		// by To (app_id is globally unique), then verify that record really belongs to this channel
		// (matching access_key) to route the message to the correct robot. This supports
		// "one channel, many numbers" and guards against stale-number / cross-channel mismatches.
		toRecord, recErr := common.GetWechatAppInfo(`app_id`, normalizeWhatsappPhone(item.To))
		if recErr != nil {
			logs.Error(`whatsapp: get app info by To failed msgid:%s to:%s err:%s`, item.MessageId, item.To, recErr.Error())
			continue
		}
		if len(toRecord) == 0 || toRecord[`app_type`] != lib_define.AppWhatsapp || toRecord[`access_key`] != accessKey {
			logs.Info(`whatsapp: skip mismatched To (likely status callback) msgid:%s to:%s access_key:%s`, item.MessageId, item.To, accessKey)
			continue
		}
		// Dedup: Aliyun retries pushing the same MsgId (more so under multi-IP/slow responses);
		// SetNX ensures we process it only once to avoid duplicate replies.
		if item.MessageId != `` {
			ok, dErr := define.Redis.SetNX(context.Background(),
				define.WhatsappInboundDedup+item.MessageId, `1`, 10*time.Minute).Result()
			if dErr != nil {
				logs.Error(`whatsapp: dedup setnx failed msgid:%s err:%s`, item.MessageId, dErr.Error())
			} else if !ok {
				logs.Info(`whatsapp: skip duplicate msgid:%s`, item.MessageId)
				continue
			}
		}
		adminUserId := cast.ToInt(toRecord[`admin_user_id`])
		nsqMsg := map[string]interface{}{
			`appid`:        toRecord[`app_id`],
			`ToUserName`:   item.To,
			`FromUserName`: item.From,
			`CreateTime`:   item.Timestamp.String(),
			`MsgId`:        item.MessageId,
		}

		switch item.Type {
		case `TEXT`:
			nsqMsg[`MsgType`] = `text`
			nsqMsg[`Content`] = item.Message

		case `IMAGE`, `AUDIO`, `VIDEO`, `STICKER`:
			// Normalise MsgType: IMAGE→image, AUDIO→voice, VIDEO→video, STICKER→image
			switch item.Type {
			case `IMAGE`, `STICKER`:
				nsqMsg[`MsgType`] = `image`
			case `AUDIO`:
				nsqMsg[`MsgType`] = `voice`
			case `VIDEO`:
				nsqMsg[`MsgType`] = `video`
			}

			// Parse embedded media JSON to get the download URL.
			var media whatsappMediaMessage
			if parseErr := tool.JsonDecodeUseNumber(item.Message, &media); parseErr != nil || media.Url == `` {
				logs.Error(`whatsapp: parse media json failed msgid:%s err:%v`, item.MessageId, parseErr)
				nsqMsg[`MsgType`] = `text`
				nsqMsg[`Content`] = `[媒体消息]`
				break
			}

			ossUrl, _ := getWhatsappFileByUrl(media.Url, media.MimeType, adminUserId)
			if ossUrl == `` {
				// Transfer failed — downgrade to text rather than exposing raw Aliyun URL.
				logs.Error(`whatsapp: media transfer failed msgid:%s url:%s`, item.MessageId, media.Url)
				nsqMsg[`MsgType`] = `text`
				nsqMsg[`Content`] = `[媒体消息]`
			} else {
				nsqMsg[`oss_url`] = ossUrl
			}

		case `DOCUMENT`:
			// Downgrade document to text: "filename(url)"
			var media whatsappMediaMessage
			_ = tool.JsonDecodeUseNumber(item.Message, &media)
			filename := media.Filename
			if filename == `` {
				filename = `文件`
			}
			nsqMsg[`MsgType`] = `text`
			if media.Url != `` {
				nsqMsg[`Content`] = fmt.Sprintf(`%s(%s)`, filename, media.Url)
			} else {
				nsqMsg[`Content`] = filename
			}

		default:
			// Non-conversational messages such as SYSTEM (e.g. customer_changed_number notifications)
			// should not be fed to the robot as a user question, so skip them.
			if item.Type == `SYSTEM` {
				logs.Info(`whatsapp: skip system message msgid:%s body:%s`, item.MessageId, item.Message)
				continue
			}
			logs.Warning(`whatsapp: unsupported type:%s msgid:%s`, item.Type, item.MessageId)
			nsqMsg[`MsgType`] = `text`
			nsqMsg[`Content`] = item.Message
		}

		common.PushNSQ(nsqMsg)
	}

	c.JSON(http.StatusOK, gin.H{`code`: 0, `msg`: `success`})
}
