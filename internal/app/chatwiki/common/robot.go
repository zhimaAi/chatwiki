// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/curl"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

func GetLibraryRobotInfo(adminUserId, libraryId int) ([]msql.Params, error) {
	robotInfo, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Field(`id,robot_key,robot_name,robot_intro,robot_avatar,application_type,library_ids`).
		Order(`id desc`).
		Select()
	if err != nil {
		return nil, err
	}
	robotData := []msql.Params{}
	if len(robotInfo) == 0 {
		return []msql.Params{}, nil
	}
	for _, item := range robotInfo {
		libraryIds := strings.Split(cast.ToString(item[`library_ids`]), ",")
		if tool.InArrayString(cast.ToString(libraryId), libraryIds) {
			robotData = append(robotData, item)
		}
	}
	return robotData, nil
}

func GetFormRobotInfo(adminUserId, formId int) ([]msql.Params, error) {
	robotInfo, err := msql.Model(`chat_ai_robot`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Field(`id,robot_key,robot_name,robot_intro,robot_avatar,application_type,form_ids`).
		Order(`id desc`).
		Select()
	if err != nil {
		return nil, err
	}
	robotData := []msql.Params{}
	if len(robotInfo) == 0 {
		return []msql.Params{}, nil
	}
	for _, item := range robotInfo {
		formIds := strings.Split(cast.ToString(item[`form_ids`]), ",")
		if tool.InArrayString(cast.ToString(formId), formIds) {
			robotData = append(robotData, item)
		}
	}
	return robotData, nil
}

type responseSaveRobot struct {
	Res  int    `json:"res"`
	Msg  string `json:"msg"`
	Data struct {
		Id       string `json:"id"`
		RobotKey string `json:"robot_key"`
	} `json:"data"`
}

func RobotAutoAdd(token string, adminUserId int) (msql.Params, error) {
	var (
		res responseSaveRobot
		err error
	)

	robotInfo := map[string]string{
		`robot_name`: `默认机器人`,
	}

	req := curl.Post(fmt.Sprintf(`http://127.0.0.1:%s/manage/saveRobot`, define.Config.WebService[`port`])).Header(`token`, token)
	for key, item := range robotInfo {
		req.Param(key, item)
	}
	if err = req.ToJSON(&res); err != nil {
		logs.Error(err.Error())
	}
	if cast.ToInt(res.Res) != define.StatusOK {
		return nil, fmt.Errorf(`%s`, cast.ToString(res.Msg))
	}
	return msql.Params{
		`id`:         cast.ToString(res.Data.Id),
		`robot_key`:  res.Data.RobotKey,
		`robot_name`: robotInfo[`robot_name`],
	}, nil
}

type RobotMessagesCacheBuildHandler struct{ RobotKey string }

func (h *RobotMessagesCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.robot_messages.%s`, h.RobotKey)
}
func SetRobotMessageCache(robotKey, question, messageId, conf string) error {
	var (
		expireTime  = 86400
		err         error
		cacheConfig = map[string]any{}
	)
	if err = tool.JsonDecode(conf, &cacheConfig); err != nil {
		logs.Error(err.Error())
		return err
	}
	if cast.ToInt(cacheConfig[`cache_switch`]) == define.SwitchOff {
		return nil
	}
	if cast.ToInt(cacheConfig[`valid_time`]) > 0 {
		expireTime = cast.ToInt(cacheConfig[`valid_time`])
	}

	// 构建缓存key，包含robotKey和question的MD5
	r := RobotMessagesCacheBuildHandler{RobotKey: robotKey}
	cacheKey := fmt.Sprintf(`%s.%s`, r.GetCacheKey(), tool.MD5(strings.TrimSpace(question)))

	// 使用Redis Set操作设置缓存
	if err := define.Redis.Set(context.Background(), cacheKey, messageId, time.Second*time.Duration(expireTime)).Err(); err != nil {
		return err
	}
	return nil
}

func HitRobotMessageCache(robotKey, question, conf string) (bool, string) {
	var (
		messageId   string
		err         error
		cacheConfig = map[string]any{}
	)
	if err = tool.JsonDecode(conf, &cacheConfig); err != nil {
		logs.Error(err.Error())
		return false, ``
	}
	if cast.ToInt(cacheConfig[`cache_switch`]) == define.SwitchOff {
		return false, ``
	}

	// 构建缓存key，包含robotKey和question的MD5
	r := RobotMessagesCacheBuildHandler{RobotKey: robotKey}
	cacheKey := fmt.Sprintf(`%s.%s`, r.GetCacheKey(), tool.MD5(strings.TrimSpace(question)))

	// 使用Redis Get操作获取缓存
	if messageId, err = define.Redis.Get(context.Background(), cacheKey).Result(); err != nil {
		return false, ``
	}
	return cast.ToInt(messageId) > 0, messageId
}

func BuildLibraryMessagesFromCache(robotKey, messageId string) ([]msql.Params, bool, error) {
	// 混合与知识库匹配
	// 查询chat_ai_message表的数据
	messageInfo, err := msql.Model("chat_ai_message", define.Postgres).
		Where("id", messageId).
		Field(`quote_file,content`).
		Find()
	if err != nil {
		logs.Error("查询消息缓存失败: %v", err)
		return nil, false, err
	}
	if len(messageInfo[`content`]) == 0 {
		return nil, false, nil
	}
	// 解析quote_file为json数组
	var quoteFile = make([]msql.Params, 0)
	if err := tool.JsonDecode(cast.ToString(messageInfo["quote_file"]), &quoteFile); err != nil {
		logs.Error("解析quote_file失败: %v", err)
		return nil, true, err
	}
	if len(quoteFile) <= 0 {
		return nil, true, nil
	}
	// 如果有引用文件，查询文件数据
	fileData, err := msql.Model(`chat_ai_answer_source`, define.Postgres).Alias(`s`).
		Join(`chat_ai_library_file_data d`, `d.id=s.paragraph_id`, `inner`).
		Where(`d.delete_time`, `0`).
		Where(`s.message_id`, messageId).Field(`d.*,s.similarity`).Select()
	if err != nil {
		logs.Error("查询引用文件数据失败: %v", err)
		return nil, true, err
	}
	for i, one := range fileData {
		fileInfo, _ := GetLibFileInfo(cast.ToInt(one[`file_id`]), 0)
		fileData[i][`file_name`] = fileInfo[`file_name`]
	}
	// 返回数据
	return changeListContent(fileData), true, nil
}

func ResponseMessagesFromCache(messageId string, useStream bool, chanStream chan sse.Event) (adaptor.ZhimaChatCompletionResponse, int64, error) {
	chatResp := adaptor.ZhimaChatCompletionResponse{}
	requestStartTime := time.Now()
	// 查询chat_ai_message表的数据
	content, err := msql.Model("chat_ai_message", define.Postgres).
		Where("id", messageId).
		Value(`content`)
	if err != nil {
		logs.Error("查询消息缓存失败: %v", err)
		return chatResp, 0, err
	}
	if len(content) == 0 {
		return chatResp, 0, errors.New("消息缓存不存在")
	}
	requestTime := time.Now().Sub(requestStartTime).Milliseconds()
	chanStream <- sse.Event{Event: `request_time`, Data: requestTime}
	// 如果使用流式输出
	if useStream {
		chanStream <- sse.Event{Event: `sending`, Data: content}
	}
	chatResp.Result = content
	// 返回数据
	return chatResp, requestTime, nil
}
