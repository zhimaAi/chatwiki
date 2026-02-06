// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_define"
	"fmt"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func createNewReceiver(params *define.ChatRequestParam, sessionId int64) {
	if len(params.Customer) > 0 && cast.ToInt(params.Customer[`is_background`]) > 0 {
		return // Do not display chat test sessions in real-time sessions
	}
	var appId string
	if len(params.ChatBaseParam.AppInfo) > 0 {
		appId = params.ChatBaseParam.AppInfo[`app_id`]
	}
	channels := make(map[string]string)
	for _, item := range GetChannelList(params.ChatBaseParam.AdminUserId, cast.ToUint(params.ChatBaseParam.Robot[`id`])) {
		channels[fmt.Sprintf(`%s_%s`, item.AppType, item.AppId)] = item.AppName
	}
	question := GetFirstQuestionByInput(params.Question) // Special handling for multimodal input
	appName := channels[fmt.Sprintf(`%s_%s`, params.ChatBaseParam.AppType, appId)]
	data := msql.Datas{
		`admin_user_id`:     params.ChatBaseParam.AdminUserId,
		`robot_id`:          params.ChatBaseParam.Robot[`id`],
		`robot_key`:         params.ChatBaseParam.Robot[`robot_key`],
		`openid`:            params.ChatBaseParam.Openid,
		`session_id`:        sessionId,
		`last_chat_time`:    tool.Time2Int(),
		`last_chat_message`: MbSubstr(question, 0, 1000),
		`app_type`:          params.ChatBaseParam.AppType,
		`app_id`:            appId,
		`come_from`:         tool.JsonEncodeNoError(map[string]string{`robot_name`: params.ChatBaseParam.Robot[`robot_name`], `app_name`: appName}),
		`unread`:            0,
		`create_time`:       tool.Time2Int(),
		`update_time`:       tool.Time2Int(),
		`rel_user_id`:       params.RelUserId,
	}
	if len(params.Customer) > 0 {
		data[`nickname`] = params.Customer[`nickname`]
		data[`name`] = params.Customer[`name`]
		data[`avatar`] = params.Customer[`avatar`]
	} else {
		data[`nickname`] = lib_define.DefaultCustomerName
		data[`name`] = lib_define.DefaultCustomerName
		data[`avatar`] = define.DefaultCustomerAvatar
	}
	// Get associated user info
	if params.RelUserId > 0 {
		FillRelUserInfo(data, params.RelUserId)
	}
	m := msql.Model(`chat_ai_receiver`, define.Postgres)
	id, err := m.Insert(data, `id`)
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return
	}
	// WebSocket notify
	ReceiverChangeNotify(params.ChatBaseParam.AdminUserId, `create`, ToStringMap(data, `id`, id))
}

func updateReceiver(sessionId int, lastChat msql.Datas, isCustomer int) {
	lastChat[`update_time`] = tool.Time2Int()
	m := msql.Model(`chat_ai_receiver`, define.Postgres)
	info, err := m.Where(`session_id`, cast.ToString(sessionId)).Find()
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return
	}
	if len(info) == 0 {
		return //no receiver info
	}
	lastChat[`rel_user_id`] = cast.ToInt(info[`rel_user_id`])
	if isCustomer == define.MsgFromCustomer {
		lastChat[`unread`] = cast.ToInt(info[`unread`]) + 1
	}
	// Sync update CustomerInfo
	if customer, _ := GetCustomerInfo(info[`openid`], cast.ToInt(info[`admin_user_id`])); len(customer) > 0 {
		lastChat[`nickname`] = customer[`nickname`]
		lastChat[`name`] = customer[`name`]
		lastChat[`avatar`] = customer[`avatar`]
	}
	// Get associated user info
	relUserId := cast.ToInt(lastChat[`rel_user_id`])
	if relUserId > 0 {
		FillRelUserInfo(lastChat, relUserId)
	}
	if _, err = m.Where(`id`, info[`id`]).Update(lastChat); err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
	}
	for key, val := range lastChat {
		info[key] = cast.ToString(val)
	}
	// WebSocket notify
	ReceiverChangeNotify(cast.ToInt(info[`admin_user_id`]), `update`, info)
}

func FillRelUserInfo(info msql.Datas, relUserId int, relUserInfos ...msql.Params) {
	if relUserId > 0 {
		var relUserInfo msql.Params
		if len(relUserInfos) == 0 {
			relUserInfo = GetUserInfo(relUserId)
		} else {
			relUserInfo = relUserInfos[0]
		}
		if len(relUserInfo) > 0 {
			info[`name`] = relUserInfo[`nick_name`]
			info[`avatar`] = relUserInfo[`avatar`]
			info[`nickname`] = relUserInfo[`nick_name`]
			if info[`avatar`] == `` {
				info[`avatar`] = `/upload/default/robot_avatar.svg`
			}
		}
	}
}

func FillRelUserInfo2(info msql.Params, relUserId int, relUserInfos ...msql.Params) {
	if relUserId > 0 {
		var relUserInfo msql.Params
		if len(relUserInfos) == 0 {
			relUserInfo = GetUserInfo(relUserId)
		} else {
			relUserInfo = relUserInfos[0]
		}
		if len(relUserInfo) > 0 {
			info[`name`] = relUserInfo[`nick_name`]
			info[`avatar`] = relUserInfo[`avatar`]
			info[`nickname`] = relUserInfo[`nick_name`]
			if info[`avatar`] == `` {
				info[`avatar`] = `/upload/default/robot_avatar.svg`
			}
		}
	}
}

func DeleteReceiver() {
	m := msql.Model(`chat_ai_receiver`, define.Postgres)
	list, err := m.Where(`last_chat_time`, `<=`, cast.ToString(tool.Time2Int()-GetSessionSecond())).ColumnObj(`admin_user_id`, `id`)
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return
	}
	if len(list) == 0 {
		return // No closed receiver
	}
	for id, adminUserId := range list {
		if _, err = m.Where(`id`, id).Delete(); err != nil {
			logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		}
		// WebSocket notify
		ReceiverChangeNotify(cast.ToInt(adminUserId), `delete`, id)
	}
}
