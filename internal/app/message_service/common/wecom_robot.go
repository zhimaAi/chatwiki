// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/pkg/lib_define"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/zhimaAi/go_tools/curl"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

func BuildWelcomesReply(welcomes string) string {
	menuJson := struct {
		Content string `json:"content"`
	}{}
	_ = tool.JsonDecodeUseNumber(welcomes, &menuJson)
	return menuJson.Content
}

func WecomMsgEncrypt(message map[string]any, receiveid string) (map[string]any, error) {
	jsonStr, err := tool.JsonEncode(message)
	if err != nil {
		return nil, err
	}
	key, err := tool.Base64Decode(lib_define.AesKey + `=`)
	if err != nil {
		return nil, err
	}
	jsonStr = tool.Random(16) + string(Uint32ToBytes(uint32(len(jsonStr)))) + jsonStr + receiveid
	bs, err := AesEncrypt([]byte(jsonStr), []byte(key))
	if err != nil {
		return nil, err
	}
	encrypt := tool.Base64Encode(string(bs))
	timestamp, nonce := tool.Time2String(), tool.Random(10)
	sign := GenerateSignature(timestamp, nonce, encrypt)
	result := map[string]any{
		`encrypt`:      encrypt,
		`msgsignature`: sign,
		`timestamp`:    timestamp,
		`nonce`:        nonce,
	}
	return result, nil
}

func GetWecomFileByUrl(link string, userId int) (string, string) {
	request := curl.Get(link)
	resp, err := request.Response()
	if err != nil {
		logs.Error(`请求失败:%s/%s`, link, err.Error())
		return ``, ``
	}
	if resp.StatusCode != http.StatusOK {
		logs.Error(`响应失败:%s/%d`, link, resp.StatusCode)
		return ``, ``
	}
	content, err := request.String()
	if err != nil {
		logs.Error(`请求失败:%s/%s`, link, err.Error())
		return ``, ``
	}
	if len(content) == 0 {
		logs.Error(`响应为空:%s/%d`, link, resp.StatusCode)
		return ``, ``
	}
	key, _ := tool.Base64Decode(lib_define.AesKey + `=`)
	bs, err := AesDecrypt([]byte(content), []byte(key))
	if err != nil {
		logs.Error(`解密失败:%s/%s`, link, err.Error())
		return ``, ``
	}
	var filename string
	objectKey := fmt.Sprintf("chat_ai/%d/received_message_images/%s/%s", userId, tool.Date(`Ym`), tool.MD5(content))
	_, params, _ := mime.ParseMediaType(resp.Header.Get(`Content-Disposition`))
	if len(params) > 0 {
		filename, _ = url.QueryUnescape(params[`filename`])
		objectKey += filepath.Ext(params[`filename`])
	}
	if err = tool.WriteFile(`internal/app/chatwiki/upload/`+objectKey, string(bs)); err != nil {
		logs.Error(`写入失败:%s/%s`, link, err.Error())
		return ``, ``
	}
	return `【image_domain】/upload/` + objectKey, filename
}

type WecomMsg struct {
	Msgtype string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
	Image struct {
		Url string `json:"url"`
	} `json:"image"`
}

type MixedItem struct {
	Type     string           `json:"type"`
	Text     string           `json:"text,omitempty"`
	ImageUrl adaptor.ImageUrl `json:"image_url,omitzero"`
}

// ParseMsgItem 解析图文混排
func ParseMsgItem(msgItem any, userId int) string {
	list := make([]WecomMsg, 0)
	err := tool.JsonDecodeUseNumber(tool.JsonEncodeNoError(msgItem), &list)
	if err != nil {
		logs.Error(`解析失败:%v/%s`, msgItem, err.Error())
	}
	multiple := make([]MixedItem, 0)
	for _, msg := range list {
		switch msg.Msgtype {
		case `text`:
			multiple = append(multiple, MixedItem{Type: adaptor.TypeText, Text: msg.Text.Content})
		case `image`:
			if fileurl, _ := GetWecomFileByUrl(msg.Image.Url, userId); len(fileurl) > 0 {
				multiple = append(multiple, MixedItem{Type: adaptor.TypeImage, ImageUrl: adaptor.ImageUrl{Url: fileurl}})
			}
		}
	}
	return tool.JsonEncodeNoError(multiple)
}
