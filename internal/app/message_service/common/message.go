// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"fmt"
	"time"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/curl"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

type SyncMsgResp struct {
	Errcode    int                      `json:"errcode"`
	Errmsg     string                   `json:"errmsg"`
	NextCursor string                   `json:"next_cursor"`
	HasMore    int                      `json:"has_more"`
	MsgList    []map[string]interface{} `json:"msg_list"`
	SendTime   int
}

func SyncMsg(appInfo msql.Params, minTime int, batchtime int64) *SyncMsgResp {
	body, err := tool.JsonEncode(map[string]interface{}{
		`cursor`: GetMessageCursorCache(appInfo[`app_id`]),
		`token`:  GetSyncMsgTokenCache(appInfo[`app_id`]),
		`limit`:  1000,
	})
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	data := SyncMsgResp{}
	corpAccessToken := GetCorpAccessToken(appInfo)
	if len(corpAccessToken) == 0 {
		return nil
	}
	link := `https://qyapi.weixin.qq.com/cgi-bin/kf/sync_msg?access_token=` + corpAccessToken
	if err = curl.Post(link).Body(body).ToJSON(&data); err != nil {
		logs.Error(err.Error())
		return nil
	}
	if data.Errcode != 0 {
		logs.Error(`data:%v`, data)
		return nil
	}

	reqtime := time.Now().Format("2006-01-02 15:04:05.000")
	fmt.Println(fmt.Sprintf(`%s[sync_msg]appid:%s,batchtime:%d,body:%s,next_cursor:%s,has_more:%d,msg_len:%d`,
		reqtime, appInfo[`app_id`], batchtime, body, data.NextCursor, data.HasMore, len(data.MsgList)))

	SetMessageCursorCache(appInfo[`app_id`], data.NextCursor)
	for k, v := range data.MsgList {
		sendTime := cast.ToInt(v[`send_time`])
		if sendTime > data.SendTime {
			data.SendTime = sendTime
		}
		if sendTime < minTime {
			data.MsgList[k] = nil
			continue
		}
		msgtype := cast.ToString(v[`msgtype`])
		if info, ok := v[msgtype].(map[string]interface{}); ok {
			for kk, vv := range info {
				data.MsgList[k][kk] = vv
			}
			delete(data.MsgList[k], msgtype)
		} else {
			logs.Error(`one:%v`, v)
			data.MsgList[k] = nil
		}
	}
	return &data
}

func GetKfMsgOrEvent(appInfo msql.Params) {
	if !SetMessageRunningLock(appInfo[`app_id`]) {
		return
	}
	defer DelMessageRunningLock(appInfo[`app_id`])
	batchtime := time.Now().UnixNano()
	for {
		KeepMessageRunningLock(appInfo[`app_id`])
		minTime := cast.ToInt(appInfo[`create_time`])
		data := SyncMsg(appInfo, max(minTime, tool.Time2Int()-86400*3), batchtime)
		if data == nil {
			return
		}
		for _, msg := range data.MsgList {
			if msg != nil {
				go MessageOne(appInfo, msg)
			}
		}
		lasttime := GetMessageLasttimeCache(appInfo[`app_id`])
		if data.HasMore == 0 && (data.SendTime == 0 || data.SendTime >= lasttime) {
			return
		}
	}
}

func MessageOne(appInfo msql.Params, msg map[string]interface{}) {
	switch cast.ToString(msg[`msgtype`]) {
	case `event`:
		switch cast.ToString(msg[`event_type`]) {
		case `enter_session`:
			pushEnterSession(appInfo, msg)
		case `msg_send_fail`:
			pushSendFail(appInfo, msg)
		case `user_recall_msg`:
			pushUserRecallMsg(appInfo, msg)
		default:
			logs.Warning(`appid:%s,msg:%v`, appInfo[`app_id`], msg)
		}
	case `text`:
		pushText(appInfo, msg)
	case `image`, `voice`, `video`, `file`:
		pushFile(appInfo, msg)
	case `location`:
		pushLocation(appInfo, msg)
	case `miniprogram`:
		pushMiniprogram(appInfo, msg)
	case `link`:
		pushLink(appInfo, msg)
	case `business_card`:
		//abandon
	case `channels_shop_product`:
		pushChannelsShopProduct(appInfo, msg)
	case `channels_shop_order`:
		pushChannelsShopOrder(appInfo, msg)
	case `merged_msg`:
		pushMergedMsg(appInfo, msg)
	case `channels`:
		pushChannels(appInfo, msg)
	default:
		logs.Warning(`appid:%s,msg:%v`, appInfo[`app_id`], msg)
	}
}
