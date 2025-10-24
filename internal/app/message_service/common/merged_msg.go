// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"fmt"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func disposeMergedMsg(appInfo msql.Params, item any) []map[string]any {
	list := make([]map[string]any, 0)
	if rows, ok1 := item.([]any); ok1 {
		for _, row := range rows {
			temp := make(map[string]any)
			if msg, ok2 := row.(map[string]any); ok2 && msg[`msg_content`] != nil && tool.JsonDecode(cast.ToString(msg[`msg_content`]), &temp) == nil {
				msgtype := cast.ToString(temp[`msgtype`])
				if info, ok3 := temp[msgtype].(map[string]any); ok3 {
					for kk, vv := range info {
						msg[kk] = vv
					}
				} else {
					msg[`Content`] = fmt.Sprintf(`Error(%v):%v`, msg[`msgtype`], msg[`msg_content`])
					msg[`msgtype`] = `text`
				}
				delete(msg, `msg_content`)
				disposeOneMsg(appInfo, msg)
				list = append(list, msg)
			} else {
				list = append(list, map[string]any{
					`CreateTime`:  tool.Time2Int(),
					`MsgType`:     `text`,
					`sender_name`: `system`,
					`Content`:     `parsing failure`,
				})
			}
		}
	} else {
		list = append(list, map[string]any{
			`CreateTime`:  tool.Time2Int(),
			`MsgType`:     `text`,
			`sender_name`: `system`,
			`Content`:     `parsing failure`,
		})
	}
	return list
}

func disposeOneMsg(appInfo msql.Params, msg map[string]any) {
	msgtype := cast.ToString(msg[`msgtype`])
	switch msgtype {
	case `event`:
		msg[`Content`] = fmt.Sprintf(`unsupported event type(%v)`, msg[`event_type`])
		msg[`msgtype`] = `text`
	case `text`:
		msg[`Content`] = msg[`content`]
		delete(msg, `content`)
	case `image`, `voice`, `video`, `file`:
		msg[`MediaId`] = msg[`media_id`]
		delete(msg, `media_id`)
	case `location`:
		msg[`Location_X`] = msg[`latitude`]
		msg[`Location_Y`] = msg[`longitude`]
		msg[`Label`] = msg[`name`]
		delete(msg, `latitude`)
		delete(msg, `longitude`)
		delete(msg, `name`)
	case `miniprogram`:
		msg[`msgtype`] = `miniprogrampage`
		msg[`Title`] = msg[`title`]
		msg[`AppId`] = msg[`appid`]
		msg[`PagePath`] = msg[`pagepath`]
		msg[`ThumbMediaId`] = msg[`thumb_media_id`]
		delete(msg, `title`)
		delete(msg, `appid`)
		delete(msg, `pagepath`)
		delete(msg, `thumb_media_id`)
	case `link`:
		msg[`msgtype`] = `multiple`
		msg[`description`] = msg[`desc`]
		msg[`thumb_url`] = msg[`pic_url`]
		delete(msg, `desc`)
		delete(msg, `pic_url`)
	case `channels_shop_product`:
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
		msg[`Content`] = content
		delete(msg, `product_id`)
		delete(msg, `head_image`)
		delete(msg, `title`)
		delete(msg, `sales_price`)
		delete(msg, `shop_nickname`)
		delete(msg, `shop_head_image`)
	case `channels_shop_order`:
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
		msg[`Content`] = content
		delete(msg, `order_id`)
		delete(msg, `product_titles`)
		delete(msg, `price_wording`)
		delete(msg, `state`)
		delete(msg, `image_url`)
		delete(msg, `shop_nickname`)
	case `merged_msg`:
		msg[`msgtype`] = `wxkf_` + msgtype
		content, err := tool.JsonEncode(disposeMergedMsg(appInfo, msg[`item`]))
		if err != nil {
			logs.Error(`msg:%v,err:%v`, msg, err)
		}
		msg[`Content`] = content
		delete(msg, `item`)
	case `channels`:
		msg[`msgtype`] = `wxkf_` + msgtype
		content, err := tool.JsonEncode(map[string]interface{}{
			`sub_type`: msg[`sub_type`],
			`nickname`: msg[`nickname`],
			`title`:    msg[`title`],
		})
		if err != nil {
			logs.Error(`msg:%v,err:%v`, msg, err)
		}
		msg[`Content`] = content
		delete(msg, `sub_type`)
		delete(msg, `nickname`)
		delete(msg, `title`)
	default:
		msg[`Content`] = fmt.Sprintf(`unsupported message type(%v)`, msgtype)
		msg[`msgtype`] = `text`
	}
	msg[`CreateTime`] = msg[`send_time`]
	msg[`MsgType`] = msg[`msgtype`]
	delete(msg, `send_time`)
	delete(msg, `msgtype`)
}
