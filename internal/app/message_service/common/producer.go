// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func pushEnterSession(appInfo msql.Params, msg map[string]interface{}) {
	message := map[string]interface{}{
		`ToUserName`:   appInfo[`app_id`],
		`FromUserName`: msg[`external_userid`],
		`CreateTime`:   msg[`send_time`],
		`MsgType`:      msg[`msgtype`],
		`Event`:        msg[`event_type`],
		`MsgId`:        msg[`msgid`],
		`open_kfid`:    msg[`open_kfid`],
		`appid`:        appInfo[`app_id`],
		`origin`:       msg[`origin`],
		`scene`:        msg[`scene`],
	}
	if msg[`scene_param`] != nil {
		message[`scene_param`] = msg[`scene_param`]
	}
	if msg[`welcome_code`] != nil {
		message[`welcome_code`] = msg[`welcome_code`]
	}
	if msg[`wechat_channels`] != nil {
		message[`wechat_channels`] = msg[`wechat_channels`]
	}
	PushNSQ(message)
}

func pushSendFail(appInfo msql.Params, msg map[string]interface{}) {
	message := map[string]interface{}{
		`ToUserName`:   appInfo[`app_id`],
		`FromUserName`: msg[`external_userid`],
		`CreateTime`:   msg[`send_time`],
		`MsgType`:      msg[`msgtype`],
		`Event`:        msg[`event_type`],
		`MsgId`:        msg[`msgid`],
		`open_kfid`:    msg[`open_kfid`],
		`appid`:        appInfo[`app_id`],
		`origin`:       msg[`origin`],
		`fail_msgid`:   msg[`fail_msgid`],
		`fail_type`:    msg[`fail_type`],
	}
	PushNSQ(message)
}

func pushUserRecallMsg(appInfo msql.Params, msg map[string]interface{}) {
	message := map[string]interface{}{
		`ToUserName`:   appInfo[`app_id`],
		`FromUserName`: msg[`external_userid`],
		`CreateTime`:   msg[`send_time`],
		`MsgType`:      msg[`msgtype`],
		`Event`:        msg[`event_type`],
		`MsgId`:        msg[`msgid`],
		`open_kfid`:    msg[`open_kfid`],
		`appid`:        appInfo[`app_id`],
		`origin`:       msg[`origin`],
		`recall_msgid`: msg[`recall_msgid`],
	}
	PushNSQ(message)
}

func pushText(appInfo msql.Params, msg map[string]interface{}) {
	message := map[string]interface{}{
		`ToUserName`:   appInfo[`app_id`],
		`FromUserName`: msg[`external_userid`],
		`CreateTime`:   msg[`send_time`],
		`MsgType`:      msg[`msgtype`],
		`Content`:      msg[`content`],
		`MsgId`:        msg[`msgid`],
		`open_kfid`:    msg[`open_kfid`],
		`appid`:        appInfo[`app_id`],
		`origin`:       msg[`origin`],
		`menu_id`:      msg[`menu_id`],
	}
	PushNSQ(message)
}

func pushFile(appInfo msql.Params, msg map[string]interface{}) {
	message := map[string]interface{}{
		`ToUserName`:   appInfo[`app_id`],
		`FromUserName`: msg[`external_userid`],
		`CreateTime`:   msg[`send_time`],
		`MsgType`:      msg[`msgtype`],
		`MediaId`:      msg[`media_id`],
		`MsgId`:        msg[`msgid`],
		`open_kfid`:    msg[`open_kfid`],
		`appid`:        appInfo[`app_id`],
		`origin`:       msg[`origin`],
	}
	PushNSQ(message)
}

func pushLocation(appInfo msql.Params, msg map[string]interface{}) {
	message := map[string]interface{}{
		`ToUserName`:   appInfo[`app_id`],
		`FromUserName`: msg[`external_userid`],
		`CreateTime`:   msg[`send_time`],
		`MsgType`:      msg[`msgtype`],
		`Location_X`:   msg[`latitude`],
		`Location_Y`:   msg[`longitude`],
		`Label`:        msg[`name`],
		`address`:      msg[`address`],
		`MsgId`:        msg[`msgid`],
		`open_kfid`:    msg[`open_kfid`],
		`appid`:        appInfo[`app_id`],
		`origin`:       msg[`origin`],
	}
	PushNSQ(message)
}

func pushMiniprogram(appInfo msql.Params, msg map[string]interface{}) {
	message := map[string]interface{}{
		`ToUserName`:   appInfo[`app_id`],
		`FromUserName`: msg[`external_userid`],
		`CreateTime`:   msg[`send_time`],
		`MsgType`:      `miniprogrampage`, //msg[`msgtype`]=miniprogram
		`Title`:        msg[`title`],
		`AppId`:        msg[`appid`],
		`PagePath`:     msg[`pagepath`],
		`ThumbMediaId`: msg[`thumb_media_id`],
		`MsgId`:        msg[`msgid`],
		`open_kfid`:    msg[`open_kfid`],
		`appid`:        appInfo[`app_id`],
		`origin`:       msg[`origin`],
	}
	PushNSQ(message)
}

func pushLink(appInfo msql.Params, msg map[string]interface{}) {
	message := map[string]interface{}{
		`ToUserName`:   appInfo[`app_id`],
		`FromUserName`: msg[`external_userid`],
		`CreateTime`:   msg[`send_time`],
		`MsgType`:      `multiple`, //msg[`msgtype`]=link
		`title`:        msg[`title`],
		`description`:  msg[`desc`],
		`url`:          msg[`url`],
		`thumb_url`:    msg[`pic_url`],
		`MsgId`:        msg[`msgid`],
		`open_kfid`:    msg[`open_kfid`],
		`appid`:        appInfo[`app_id`],
		`origin`:       msg[`origin`],
	}
	PushNSQ(message)
}

func pushChannelsShopProduct(appInfo msql.Params, msg map[string]interface{}) {
	content, err := tool.JsonEncode(map[string]interface{}{
		`product_id`:      msg[`product_id`],
		`head_image`:      msg[`head_image`],
		`title`:           msg[`title`],
		`sales_price`:     msg[`sales_price`],
		`shop_nickname`:   msg[`shop_nickname`],
		`shop_head_image`: msg[`shop_head_image`],
	})
	if err != nil {
		logs.Error(`msg:%v,err:%v`, msg, err)
	}
	message := map[string]interface{}{
		`ToUserName`:   appInfo[`app_id`],
		`FromUserName`: msg[`external_userid`],
		`CreateTime`:   msg[`send_time`],
		`MsgType`:      msg[`msgtype`],
		`Content`:      content,
		`MsgId`:        msg[`msgid`],
		`open_kfid`:    msg[`open_kfid`],
		`appid`:        appInfo[`app_id`],
		`origin`:       msg[`origin`],
	}
	PushNSQ(message)
}

func pushChannelsShopOrder(appInfo msql.Params, msg map[string]interface{}) {
	content, err := tool.JsonEncode(map[string]interface{}{
		`order_id`:       msg[`order_id`],
		`product_titles`: msg[`product_titles`],
		`price_wording`:  msg[`price_wording`],
		`state`:          msg[`state`],
		`image_url`:      msg[`image_url`],
		`shop_nickname`:  msg[`shop_nickname`],
	})
	if err != nil {
		logs.Error(`msg:%v,err:%v`, msg, err)
	}
	message := map[string]interface{}{
		`ToUserName`:   appInfo[`app_id`],
		`FromUserName`: msg[`external_userid`],
		`CreateTime`:   msg[`send_time`],
		`MsgType`:      msg[`msgtype`],
		`Content`:      content,
		`MsgId`:        msg[`msgid`],
		`open_kfid`:    msg[`open_kfid`],
		`appid`:        appInfo[`app_id`],
		`origin`:       msg[`origin`],
	}
	PushNSQ(message)
}

func pushMergedMsg(appInfo msql.Params, msg map[string]interface{}) {
	content, err := tool.JsonEncode(disposeMergedMsg(appInfo, msg[`item`]))
	if err != nil {
		logs.Error(`msg:%v,err:%v`, msg, err)
	}
	message := map[string]interface{}{
		`ToUserName`:   appInfo[`app_id`],
		`FromUserName`: msg[`external_userid`],
		`CreateTime`:   msg[`send_time`],
		`MsgType`:      `wxkf_` + cast.ToString(msg[`msgtype`]),
		`title`:        msg[`title`],
		`Content`:      content, //json
		`MsgId`:        msg[`msgid`],
		`open_kfid`:    msg[`open_kfid`],
		`appid`:        appInfo[`app_id`],
		`origin`:       msg[`origin`],
	}
	PushNSQ(message)
}

func pushChannels(appInfo msql.Params, msg map[string]interface{}) {
	content, err := tool.JsonEncode(map[string]interface{}{
		`sub_type`: msg[`sub_type`],
		`nickname`: msg[`nickname`],
		`title`:    msg[`title`],
	})
	if err != nil {
		logs.Error(`msg:%v,err:%v`, msg, err)
	}
	message := map[string]interface{}{
		`ToUserName`:   appInfo[`app_id`],
		`FromUserName`: msg[`external_userid`],
		`CreateTime`:   msg[`send_time`],
		`MsgType`:      `wxkf_` + cast.ToString(msg[`msgtype`]),
		`Content`:      content, //json
		`MsgId`:        msg[`msgid`],
		`open_kfid`:    msg[`open_kfid`],
		`appid`:        appInfo[`app_id`],
		`origin`:       msg[`origin`],
	}
	PushNSQ(message)
}
