// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/message_service/define"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/wechat"
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

var PassiveMenuReg = regexp.MustCompile(`^chatwiki_passive_(\d+)_(\d+)$`)
var PassiveContentReg = regexp.MustCompile(`^查看问题【(\d+)】回复\((\d+)\):`)
var PassiveContentRegNew = regexp.MustCompile(`\((\d+)-(\d+)\)$`)

func MbSubstr(s string, start, length int) string {
	runes := []rune(s)
	if start >= len(runes) {
		return ""
	}
	if start+length > len(runes) {
		length = len(runes) - start
	}
	return string(runes[start : start+length])
}

func BuildMsgmenucontent(log msql.Params, serial int, image bool) string {
	var wordLimit = 30
	if image {
		wordLimit = 10 //图片列表...
	}
	question := MbSubstr(log[`content`], 0, wordLimit)
	if utf8.RuneCountInString(log[`content`]) > wordLimit {
		question += `...` //末尾追加省略号
	}
	return fmt.Sprintf(`%s(%s-%d)`, question, log[`id`], serial+1)
}

func CheckQueryAiReply(message map[string]any) (int64, int, bool) {
	bizmsgmenuid := cast.ToString(message[`bizmsgmenuid`])
	if len(bizmsgmenuid) > 0 {
		result := PassiveMenuReg.FindStringSubmatch(bizmsgmenuid)
		if len(result) == 3 {
			return cast.ToInt64(result[1]), cast.ToInt(result[2]), true
		}
	}
	content := cast.ToString(message[`Content`])
	if len(content) > 0 {
		result := PassiveContentReg.FindStringSubmatch(content)
		if len(result) == 3 {
			return cast.ToInt64(result[1]), cast.ToInt(result[2]) - 1, true
		}
		result = PassiveContentRegNew.FindStringSubmatch(content)
		if len(result) == 3 {
			return cast.ToInt64(result[1]), cast.ToInt(result[2]) - 1, true
		}
	}
	return 0, 0, false
}

// ReplaceXmlStr 响应xml数据占位符替换
func ReplaceXmlStr(xmlStr string, message map[string]any) string {
	xmlStr = strings.ReplaceAll(xmlStr, `:FromUserName`, cast.ToString(message[`ToUserName`]))
	xmlStr = strings.ReplaceAll(xmlStr, `:ToUserName`, cast.ToString(message[`FromUserName`]))
	xmlStr = strings.ReplaceAll(xmlStr, `:CreateTime`, tool.Time2String())
	//时间占位符替换
	dateType := map[string]string{
		`{{yyyy-MM-dd hh:mm:ss}}`: tool.Date(`Y-m-d H:i:s`),
		`{{MM-dd hh:mm:ss}}`:      tool.Date(`m-d H:i:s`),
		`{{yyyy-MM-dd hh:mm}}`:    tool.Date(`Y-m-d H:i`),
		`{{MM-dd hh:mm}}`:         tool.Date(`m-d H:i`),
		`{{yyyy-MM-dd}}`:          tool.Date(`Y-m-d`),
		`{{MM-dd}}`:               tool.Date(`m-d`),
		`{{MM月dd日}}`:              tool.Date(`n月j日`),
	}
	for k, v := range dateType {
		xmlStr = strings.ReplaceAll(xmlStr, k, v)
	}
	return xmlStr
}

var PassiveSubscribeMapLock = &sync.Mutex{}
var PassiveSubscribeMapChan = make(map[int64]chan struct{})

func PassiveSubscribeMapOp(doworkHandle func()) {
	PassiveSubscribeMapLock.Lock()
	defer PassiveSubscribeMapLock.Unlock()
	doworkHandle()
}

func WaitAiReply(robot msql.Params, message map[string]any) string {
	//先进入被动回复日志表拿到id
	m := msql.Model(`chat_ai_passive_reply_log`, define.Postgres)
	appid := cast.ToString(message[`appid`])
	msgid := cast.ToString(message[`MsgId`])
	logid, err := m.Insert(msql.Datas{
		`admin_user_id`: robot[`admin_user_id`],
		`app_id`:        appid,
		`msgid`:         msgid,
		`openid`:        message[`FromUserName`],
		`content`:       message[`Content`],
		`create_time`:   tool.Time2Int(),
		`update_time`:   tool.Time2Int(),
	}, `id`)
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
	}
	//推送消息进入NSQ,这里要注意先后顺序
	message[`passive_id`] = logid
	go PushNSQ(message)
	//等待AI回复内容
	wechatWait := cast.ToUint(define.Config.WebService[`wechat_wait`])
	select {
	case <-time.After(time.Second * time.Duration(wechatWait)):
	case <-GetAiReplyByMq(wechatWait, logid):
	}
	PassiveSubscribeMapOp(func() { //清理map+chan
		if channel, ok := PassiveSubscribeMapChan[logid]; ok {
			delete(PassiveSubscribeMapChan, logid)
			close(channel)
		}
	})
	return GetAiReply(robot, message, logid, 0)
}

func GetAiReplyByMq(wechatWait uint, logid int64) <-chan struct{} {
	if wechatWait == 0 {
		return nil
	}
	done := make(chan struct{})
	PassiveSubscribeMapOp(func() {
		PassiveSubscribeMapChan[logid] = done
	})
	return done
}

func PassiveSubscribeCall() {
	for message := range define.Redis.Subscribe(context.Background(), lib_define.RedisPrefixPassiveSubscribe).Channel() {
		PassiveSubscribeMapOp(func() {
			if channel, ok := PassiveSubscribeMapChan[cast.ToInt64(message.Payload)]; ok {
				select {
				case channel <- struct{}{}: //发送通知
				default:
					logs.Warning(`chan写入异常:%s`, message.Payload)
				}
			}
		})
	}
}

func GetAiReply(robot msql.Params, message map[string]any, logid int64, serial int) string {
	appid := cast.ToString(message[`appid`])
	openid := cast.ToString(message[`FromUserName`])
	m := msql.Model(`chat_ai_passive_reply_log`, define.Postgres)
	log, err := m.Where(`id`, cast.ToString(logid)).
		Where(`app_id`, appid).Where(`openid`, openid).Find()
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
	}
	if len(log) == 0 { //数据非法
		return `回复内容不存在...`
	}
	if !cast.ToBool(log[`status`]) { //AI回复未完成
		bizmsgmenuid := fmt.Sprintf(`chatwiki_passive_%d_%d`, logid, 0)
		msgmenucontent := BuildMsgmenucontent(log, 0, false)
		return fmt.Sprintf("%s\r\n<a href=\"weixin://bizmsgmenu?msgmenucontent=%s&msgmenuid=%s\">%s</a>",
			robot[`wechat_not_verify_hand_get_reply`], msgmenucontent, bizmsgmenuid, robot[`wechat_not_verify_hand_get_word`])
	}
	replys, images := make([]string, 0), make([]string, 0)
	_ = tool.JsonDecodeUseNumber(log[`replys`], &replys)
	_ = tool.JsonDecodeUseNumber(log[`images`], &images)
	if serial >= len(replys)+len(images) { //数据非法
		return `回复内容不存在...`
	}
	echo := append(replys, images...)[serial]
	//判断还有没有下一页
	if serial+1 < len(replys) {
		nextBizmsgmenuid := fmt.Sprintf(`chatwiki_passive_%d_%d`, logid, serial+1)
		nextMsgmenucontent := BuildMsgmenucontent(log, serial+1, false)
		jointBizmsgmenu := fmt.Sprintf("\r\n<a href=\"weixin://bizmsgmenu?msgmenucontent=%s&msgmenuid=%s\">%s</a>",
			nextMsgmenucontent, nextBizmsgmenuid, robot[`wechat_not_verify_hand_get_next`])
		if len(echo)+len(jointBizmsgmenu) < 2048 {
			echo += jointBizmsgmenu
		}
	} else if serial+1 == len(replys) && len(images) > 0 { //最后一个文本,将图片全部输出
		for idx := range images {
			bizmsgmenuid := fmt.Sprintf(`chatwiki_passive_%d_%d`, logid, len(replys)+idx)
			msgmenucontent := BuildMsgmenucontent(log, len(replys)+idx, true)
			showContent := fmt.Sprintf(`点击查看图片(%d/%d)`, idx+1, len(images))
			if len(images) == 1 { //只有一张图片的时候
				showContent = `点击查看回复图片`
			}
			jointBizmsgmenu := fmt.Sprintf("\r\n<a href=\"weixin://bizmsgmenu?msgmenucontent=%s&msgmenuid=%s\">%s</a>",
				msgmenucontent, bizmsgmenuid, showContent)
			if len(echo)+len(jointBizmsgmenu) < 2048 {
				echo += jointBizmsgmenu
			}
		}
	}
	return echo
}

func BuildXmlStr(appInfo msql.Params, message map[string]any, echo string) string {
	var xmlStr string
	if IsUrl(echo) && LinkExists(echo) { //内容是图片的
		var mediaId string
		var err error
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()
		done := DoWorkWithContext(ctx, func() {
			mediaId, err = GetMediaIdByUrl(appInfo, echo)
		})
		if done != nil {
			err = done //任务未完成,将超时error赋给err
		}
		if err == nil {
			xmlStr = strings.ReplaceAll(define.PassiveReplyImageTemplate, `:MediaId`, mediaId)
			return ReplaceXmlStr(xmlStr, message)
		} else {
			logs.Error(`appid:%s,image:%s,err:%s`, appInfo[`app_id`], echo, err.Error())
		}
	}
	//用文本消息回复出来(包含图片失败了的)
	xmlStr = strings.ReplaceAll(define.PassiveReplyTextTemplate, `:Content`, echo)
	return ReplaceXmlStr(xmlStr, message)
}

func GetMediaIdByUrl(appInfo msql.Params, image string) (string, error) {
	cacheKey := fmt.Sprintf(lib_define.RedisPrefixMediaUpload, appInfo[`app_id`], tool.MD5(image))
	mediaId := define.Redis.Get(context.Background(), cacheKey).Val()
	if len(mediaId) > 0 {
		return mediaId, nil
	}
	app, err := wechat.GetApplication(appInfo)
	if err != nil {
		return ``, err
	}
	mediaId, _, err = app.UploadTempImage(GetFileByLink(image))
	if err != nil {
		return ``, err
	}
	err = define.Redis.Set(context.Background(), cacheKey, mediaId, 70*time.Hour).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error(err.Error())
	}
	return mediaId, nil
}
