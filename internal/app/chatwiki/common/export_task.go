// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/wechat/common"
	"chatwiki/internal/pkg/lib_define"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

func checkSourceExist(source uint) bool {
	for _, item := range define.ExportSourceList {
		if cast.ToUint(item[`source`]) == source {
			return true
		}
	}
	return false
}

func CreateExportTask(adminUserId, robotId, source uint, fileName string, params map[string]any) (int64, error) {
	if !checkSourceExist(source) {
		return 0, errors.New(`导出来源参数错误`)
	}
	paramsJson, err := tool.JsonEncode(params)
	if err != nil {
		return 0, err
	}
	insData := msql.Datas{
		`admin_user_id`: adminUserId,
		`robot_id`:      robotId,
		`file_name`:     fileName,
		`source`:        source,
		`params`:        paramsJson,
		`create_time`:   tool.Time2Int(),
		`update_time`:   tool.Time2Int(),
	}
	if source == define.ExportSourceLibFileDoc {
		insData[`library_id`] = params[`library_id`]
	}
	id, err := msql.Model(`chat_ai_export_task`, define.Postgres).Insert(insData, `id`)
	if err != nil {
		return 0, err
	}
	if err = AddJobs(define.ExportTaskTopic, cast.ToString(id)); err != nil {
		logs.Error(err.Error())
		return 0, err
	}
	return id, nil
}

func RunSessionExport(params map[string]any) (string, error) {
	//获取机器人信息
	robot, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`id`, cast.ToString(params[`robot_id`])).Where(`admin_user_id`, cast.ToString(params[`admin_user_id`])).Field(`robot_name`).Find()
	if err != nil {
		logs.Error(err.Error())
		return ``, err
	}
	if len(robot) == 0 {
		return ``, errors.New(`机器人信息不存在`)
	}
	//获取会话信息
	sessionModel := msql.Model(`chat_ai_session`, define.Postgres).Alias(`s`).
		Join(`chat_ai_dialogue d`, `s.dialogue_id=d.id`, `left`)
	sessionModel.Where(`s.admin_user_id`, cast.ToString(params[`admin_user_id`]))
	sessionModel.Where(`s.robot_id`, cast.ToString(params[`robot_id`]))
	sessionModel.Where(`d.is_background`, `0`)
	if len(cast.ToString(params[`app_type`])) > 0 {
		sessionModel.Where(`s.app_type`, cast.ToString(params[`app_type`]))
	}
	if len(cast.ToString(params[`app_id`])) > 0 {
		sessionModel.Where(`s.app_id`, cast.ToString(params[`app_id`]))
	}
	startTime, endTime := cast.ToInt(params[`start_time`]), cast.ToInt(params[`end_time`])
	if startTime > 0 && endTime > 0 && endTime >= startTime {
		sessionModel.Where(`s.last_chat_time`, `between`, fmt.Sprintf(`%d,%d`, startTime, endTime))
	}
	if len(cast.ToString(params[`name`])) > 0 {
		sessionModel.Join(`chat_ai_customer c`, fmt.Sprintf(`c.admin_user_id=%d AND c.openid=s.openid`, cast.ToUint(params[`admin_user_id`])), `left`)
		sessionModel.Where(`c.name`, `like`, cast.ToString(params[`name`]))
	}
	sessions, err := sessionModel.Order(`s.id DESC`).Field(`s.id,s.app_type,s.app_id,s.openid`).Select()
	if err != nil {
		logs.Error(err.Error())
		return ``, err
	}
	if len(sessions) == 0 {
		return ``, errors.New(`没有任何会话记录`)
	}
	//获取+组装数据
	channels := make(map[string]string)
	for _, item := range GetChannelList(cast.ToInt(params[`admin_user_id`]), cast.ToUint(params[`robot_id`])) {
		channels[fmt.Sprintf(`%s_%s`, item.AppType, item.AppId)] = item.AppName
	}
	messageModel := msql.Model(`chat_ai_message`, define.Postgres)
	data := make([]map[string]any, 0)
	for _, session := range sessions {
		list, err := messageModel.Where(`admin_user_id`, cast.ToString(params[`admin_user_id`])).
			Where(`robot_id`, cast.ToString(params[`robot_id`])).
			Where(`openid`, session[`openid`]).Where(`session_id`, session[`id`]).
			Order(`id`).Field(`id,is_customer,content,create_time`).Select()
		if err != nil {
			logs.Error(err.Error())
			return ``, err
		}
		customerName := `访客XXXX`
		if customer, _ := GetCustomerInfo(session[`openid`], cast.ToInt(params[`admin_user_id`])); len(customer) > 0 {
			customerName = customer[`name`]
		}
		if session[`app_type`] == lib_define.AppWechatKefu { //特殊处理
			session[`openid`], _ = common.GetExternalUserInfo(session[`openid`])
		}
		for _, item := range list {
			var sender = robot[`robot_name`]
			if cast.ToUint(item[`is_customer`]) > 0 {
				sender = customerName
			}
			data = append(data, map[string]any{
				`msgid`:       item[`id`],
				`openid`:      session[`openid`],
				`sender`:      sender,
				`session_id`:  session[`id`],
				`content`:     item[`content`],
				`create_time`: tool.Date(``, item[`create_time`]),
				`app_name`:    channels[fmt.Sprintf(`%s_%s`, session[`app_type`], session[`app_id`])],
			})
		}
	}
	//开始导出
	fields := tool.Fields{
		{Field: "msgid", Header: "msgid"},
		{Field: "openid", Header: "用户openid"},
		{Field: "sender", Header: "发送者"},
		{Field: "session_id", Header: "会话id"},
		{Field: "content", Header: "消息内容"},
		{Field: "create_time", Header: "消息发送时间"},
		{Field: "app_name", Header: "消息来源"},
	}
	filePath := `static/public/export/session`
	file, _, err := tool.ExcelExportPro(data, fields, `会话记录导出`, filePath)
	if err != nil {
		logs.Error(err.Error())
		return ``, err
	}
	if cast.ToUint(define.Config.OssConfig[`enable`]) > 0 { //put oss
		objectKey := fmt.Sprintf(`chat_ai/%v/export_session/%s%s`, params[`admin_user_id`], tool.Date(`Ym`), file[len(filePath):])
		if link, err := PutObjectFromFile(objectKey, file); err == nil {
			if err = os.Remove(file); err != nil {
				logs.Error(err.Error()) //remove local file
			}
			return link, nil
		} else {
			logs.Error(err.Error())
		}
	}
	return file[6:], nil
}

func RunLibFileDocExport(params map[string]any) (string, string, error) {
	libraryId := cast.ToInt(params[`library_id`])
	dataIds := cast.ToString(params[`data_ids`])
	exportType := cast.ToInt(params[`export_type`])
	adminUserId := cast.ToInt(params[`admin_user_id`])
	groupId := cast.ToInt(params[`group_id`])
	filePath := ""
	ext := `xlsx`
	fileName := ""
	switch exportType {
	case define.ExportQALibDocs:
		// 导出QA文档
		libraryInfo, err := GetLibraryInfo(libraryId, adminUserId)
		if err != nil {
			logs.Error(err.Error())
			return ``, ``, err
		}
		if len(libraryInfo) == 0 {
			return ``, ``, err
		}
		fileName = libraryInfo[`library_name`] + tool.Date(`YmdHis`) + `.` + ext
		id := 0
		pageSize := 500
		data := make([]msql.Params, 0)
		// 查询QA列表
		for {
			m := msql.Model("chat_ai_library_file_data", define.Postgres).
				Where("library_id", cast.ToString(libraryId)).
				Where("admin_user_id", cast.ToString(adminUserId)).
				Where(`delete_time`, `0`).
				Where(`id`, `>`, cast.ToString(id)).
				Field("id,question,similar_questions,answer,images")
			if len(dataIds) > 0 {
				m.Where(`id`, `in`, dataIds)
			}
			if groupId >= 0 {
				m.Where(`group_id`, cast.ToString(groupId))
			}
			// 分页查询
			list, err := m.Order("id ASC").Limit(pageSize).Select()
			if err != nil {
				logs.Error(err.Error())
				return ``, ``, err
			}
			if len(list) == 0 {
				break
			}
			id = cast.ToInt(list[len(list)-1][`id`])
			data = append(data, list...)
		}
		filePath, err = ExportFAQFileAllQA(data, ext, `library_file`)
		if err != nil {
			logs.Error(err.Error())
			return ``, ``, err
		}
	}
	return filePath[6:], fileName, nil
}
