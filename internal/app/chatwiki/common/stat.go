// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

const StatsTypeDailyActiveUser = 1
const StatsTypeDailyNewUser = 2
const StatsTypeDailyMsgCount = 3
const StatsTypeDailyTokenCount = 4
const StatsTypeDailyLibraryTipCount = 5
const StatsTypeDailyAiMsgCount = 6

var statsMu sync.Mutex

func LlmLogRequest(
	_type string,
	adminUserId int,
	openid string,
	robot msql.Params,
	library msql.Params,
	config msql.Params,
	appType string,
	fileInfo msql.Params,
	model string,
	promptToken int,
	completionToken int,
	req interface{},
	resp interface{},
) error {
	statsMu.Lock()
	defer statsMu.Unlock()

	sourceRobot, err := tool.JsonEncode(robot)
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	sourceLibrary, err := tool.JsonEncode(library)
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	sourceFile, err := tool.JsonEncode(fileInfo)
	if err != nil {
		logs.Error(err.Error())
		return err
	}

	requestDetail, err := tool.JsonEncode(req)
	if err != nil {
		logs.Error(err.Error())
		return err
	}

	responseDetail, err := tool.JsonEncode(resp)
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	var corpName string
	for _, modelInfo := range GetModelList() {
		if len(config[`model_define`]) == 0 || modelInfo.ModelDefine == config[`model_define`] {
			corpName = modelInfo.ModelName
			break
		}
	}

	data := msql.Datas{
		`admin_user_id`:    adminUserId,
		`openid`:           openid,
		`corp`:             corpName,
		`model`:            model,
		`type`:             _type,
		`prompt_token`:     promptToken,
		`completion_token`: completionToken,
		`app_type`:         appType,
		`request_detail`:   requestDetail,
		`response_detail`:  responseDetail,
		`create_time`:      tool.Time2Int(),
		`update_time`:      tool.Time2Int(),
	}

	if len(robot) > 0 {
		data[`source_robot`] = sourceRobot
		data[`source_robot_id`] = cast.ToInt(robot[`id`]) //hotfix: pq: 无效的类型 integer 输入语法: ""
	}
	if len(library) > 0 {
		data[`source_library`] = sourceLibrary
	}
	if len(fileInfo) > 0 {
		data[`source_file`] = sourceFile
	}

	// record request logs
	_, err = msql.Model(`llm_request_logs`, define.Postgres).Insert(data)
	if err != nil {
		logs.Error(err.Error())
		return err
	}

	if _type == `LLM` && len(robot) > 0 {
		err = statDailyRequestCount(adminUserId, robot, appType)
		if err != nil {
			logs.Error(err.Error())
			return err
		}
		if err = statDailyActiveUser(adminUserId, robot, appType); err != nil {
			logs.Error(err.Error())
			return err
		}
		if err = statDailyNewUser(adminUserId, robot, appType, openid); err != nil {
			logs.Error(err.Error())
			return err
		}
	}

	if len(robot) > 0 {
		if err = statDailyTokenCount(adminUserId, robot, appType, promptToken, completionToken); err != nil {
			logs.Error(err.Error())
			return err
		}
	}

	if err = statToken(_type, adminUserId, corpName, model, promptToken, completionToken); err != nil {
		logs.Error(err.Error())
		return err
	}
	//app token stat
	if err = statTokenApp(_type, adminUserId, robot, corpName, model, promptToken, completionToken); err != nil {
		logs.Error(err.Error())
		return err
	}
	return nil
}

func statTokenApp(_type string, adminUserId int, robot msql.Params, corpName, model string, promptToken, completionToken int) error {
	var oldPromptToken int
	var oldCompletionToken int
	robotId := cast.ToInt(robot[`id`])
	tokenAppType := define.TokenAppTypeOther
	if len(robot) > 0 && cast.ToInt(robot[`id`]) > 0 {
		applicationType := cast.ToInt(robot[`application_type`])
		if applicationType == define.ApplicationTypeChat {
			tokenAppType = define.TokenAppTypeRobot
		} else if applicationType == define.ApplicationTypeFlow {
			tokenAppType = define.TokenAppTypeWorkflow
		}
	}
	stats, err := msql.Model(`llm_token_app_daily_stats`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`token_app_type`, cast.ToString(tokenAppType)).
		Where(`date`, time.Now().Format(`2006-01-02`)).
		Where(`robot_id`, cast.ToString(robotId)).
		Field(`prompt_token, completion_token,request_num`).
		Find()
	if err != nil {
		return err
	}
	if len(stats) == 0 {
		_, err := msql.Model(`llm_token_app_daily_stats`, define.Postgres).Insert(msql.Datas{
			`admin_user_id`:    cast.ToString(adminUserId),
			`token_app_type`:   tokenAppType,
			`robot_id`:         cast.ToString(robotId),
			`corp`:             corpName,
			`model`:            model,
			`type`:             _type,
			`date`:             time.Now().Format(`2006-01-02`),
			`request_num`:      1,
			`prompt_token`:     cast.ToString(promptToken),
			`completion_token`: cast.ToString(completionToken),
			`create_time`:      tool.Time2Int(),
			`update_time`:      tool.Time2Int(),
		})
		if err != nil {
			return err
		}
		return nil
	}
	oldPromptToken = cast.ToInt(stats[`prompt_token`])
	oldCompletionToken = cast.ToInt(stats[`completion_token`])
	requestNum := cast.ToInt(stats[`request_num`])

	_, err = msql.Model(`llm_token_app_daily_stats`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`token_app_type`, cast.ToString(tokenAppType)).
		Where(`date`, time.Now().Format(`2006-01-02`)).
		Where(`robot_id`, cast.ToString(robotId)).
		Update(msql.Datas{
			`prompt_token`:     cast.ToString(promptToken + oldPromptToken),
			`completion_token`: cast.ToString(completionToken + oldCompletionToken),
			`request_num`:      requestNum + 1,
			`update_time`:      tool.Time2Int(),
		})
	if err != nil {
		return err
	}
	return nil
}

func statToken(_type string, adminUserId int, corpName string, model string, promptToken int, completionToken int) error {
	var oldPromptToken int
	var oldCompletionToken int
	stats, err := msql.Model(`llm_token_daily_stats`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`corp`, corpName).
		Where(`model`, model).
		Where(`type`, _type).
		Where(`date`, time.Now().Format(`2006-01-02`)).
		Field(`prompt_token, completion_token`).
		Find()
	if err != nil {
		return err
	}
	if len(stats) == 0 {
		_, err := msql.Model(`llm_token_daily_stats`, define.Postgres).Insert(msql.Datas{
			`admin_user_id`: cast.ToString(adminUserId),
			`corp`:          corpName,
			`model`:         model,
			`type`:          _type,
			`date`:          time.Now().Format(`2006-01-02`),
			`create_time`:   tool.Time2Int(),
			`update_time`:   tool.Time2Int(),
		})
		if err != nil {
			return err
		}
	} else {
		oldPromptToken = cast.ToInt(stats[`prompt_token`])
		oldCompletionToken = cast.ToInt(stats[`completion_token`])
	}

	_, err = msql.Model(`llm_token_daily_stats`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`corp`, corpName).
		Where(`model`, model).
		Where(`type`, _type).
		Where(`date`, time.Now().Format(`2006-01-02`)).
		Update(msql.Datas{
			`prompt_token`:     cast.ToString(promptToken + oldPromptToken),
			`completion_token`: cast.ToString(completionToken + oldCompletionToken),
			`update_time`:      tool.Time2Int(),
		})
	if err != nil {
		return err
	}
	return nil
}

func statDailyActiveUser(adminUserId int, robot msql.Params, appType string) error {
	now := time.Now()
	todayStartTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	count, err := msql.Model(`llm_request_logs`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`source_robot_id`, robot[`id`]).
		Where(`app_type`, appType).
		Where(`create_time`, `>=`, cast.ToString(todayStartTime.Unix())).
		Value(`count(distinct concat(app_type, '-', openid))`)
	if err != nil {
		return err
	}

	row, err := msql.Model(`llm_request_daily_stats`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, robot[`id`]).
		Where(`app_type`, appType).
		Where(`date`, time.Now().Format(`2006-01-02`)).
		Where(`type`, cast.ToString(StatsTypeDailyActiveUser)).
		Find()
	if err != nil {
		return err
	}
	if len(row) == 0 {
		_, err = msql.Model(`llm_request_daily_stats`, define.Postgres).Insert(msql.Datas{
			`admin_user_id`: adminUserId,
			`robot_id`:      robot[`id`],
			`app_type`:      appType,
			`date`:          time.Now().Format(`2006-01-02`),
			`type`:          StatsTypeDailyActiveUser,
			`amount`:        count,
			`create_time`:   tool.Time2Int(),
			`update_time`:   tool.Time2Int(),
		})
		if err != nil {
			return err
		}
	} else {
		_, err = msql.Model(`llm_request_daily_stats`, define.Postgres).Where(`id`, row[`id`]).Update(msql.Datas{
			`amount`:      count,
			`update_time`: tool.Time2Int(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func statDailyNewUser(adminUserId int, robot msql.Params, appType, openid string) error {
	now := time.Now()
	todayStartTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	count, err := msql.Model(`llm_request_logs`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`source_robot_id`, robot[`id`]).
		Where(`app_type`, appType).
		Where(`create_time`, `<`, cast.ToString(todayStartTime.Unix())).
		Count()
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	if count > 0 {
		return nil
	}

	ctx := context.Background()
	key := fmt.Sprintf(`daily_new_user_%d_%s_%s_%s`, adminUserId, robot[`id`], appType, tool.Date(`Ymd`))
	exists, err := define.Redis.Exists(ctx, key).Result()
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	if exists == 0 {
		rows, err := msql.Model(`llm_request_logs`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`source_robot_id`, robot[`id`]).
			Where(`app_type`, appType).
			Where(`create_time`, `>=`, cast.ToString(todayStartTime.Unix())).
			Field(`openid`).
			Select()
		var openidList []interface{}
		for _, row := range rows {
			openidList = append(openidList, row[`openid`])
		}
		_, err = define.Redis.SAdd(ctx, key, openidList...).Result()
		if err != nil {
			logs.Error(err.Error())
			return err
		}
		_, err = define.Redis.Expire(ctx, key, 24*time.Hour).Result()
		if err != nil {
			logs.Error(err.Error())
			return err
		}
	} else {
		_, err = define.Redis.SAdd(ctx, key, openid).Result()
		if err != nil {
			logs.Error(err.Error())
			return err
		}
	}

	newUserCount, err := define.Redis.SCard(ctx, key).Result()
	if err != nil {
		logs.Error(err.Error())
		return err
	}

	row, err := msql.Model(`llm_request_daily_stats`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, robot[`id`]).
		Where(`app_type`, appType).
		Where(`date`, time.Now().Format(`2006-01-02`)).
		Where(`type`, cast.ToString(StatsTypeDailyNewUser)).
		Find()
	if len(row) == 0 {
		_, err = msql.Model(`llm_request_daily_stats`, define.Postgres).Insert(msql.Datas{
			`admin_user_id`: cast.ToString(adminUserId),
			`robot_id`:      robot[`id`],
			`app_type`:      appType,
			`date`:          time.Now().Format(`2006-01-02`),
			`type`:          cast.ToString(StatsTypeDailyNewUser),
			`amount`:        newUserCount,
			`create_time`:   tool.Time2Int(),
			`update_time`:   tool.Time2Int(),
		})
		if err != nil {
			return err
		}
	} else {
		_, err = msql.Model(`llm_request_daily_stats`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`robot_id`, robot[`id`]).
			Where(`app_type`, appType).
			Where(`date`, time.Now().Format(`2006-01-02`)).
			Where(`type`, cast.ToString(StatsTypeDailyNewUser)).
			Update(msql.Datas{
				`amount`:      newUserCount,
				`update_time`: tool.Time2Int(),
			})
		if err != nil {
			return err
		}
	}

	return nil
}

func statDailyRequestCount(adminUserId int, robot msql.Params, appType string) error {
	row, err := msql.Model(`llm_request_daily_stats`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, robot[`id`]).
		Where(`app_type`, appType).
		Where(`date`, time.Now().Format(`2006-01-02`)).
		Where(`type`, cast.ToString(StatsTypeDailyMsgCount)).
		Find()
	if err != nil {
		return err
	}
	if len(row) == 0 {
		_, err = msql.Model(`llm_request_daily_stats`, define.Postgres).Insert(msql.Datas{
			`admin_user_id`: cast.ToString(adminUserId),
			`robot_id`:      robot[`id`],
			`date`:          time.Now().Format(`2006-01-02`),
			`app_type`:      appType,
			`type`:          cast.ToString(StatsTypeDailyMsgCount),
			`amount`:        "2",
			`create_time`:   tool.Time2Int(),
			`update_time`:   tool.Time2Int(),
		})
		if err != nil {
			return err
		}
	} else {
		_, err = msql.Model(`llm_request_daily_stats`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`robot_id`, robot[`id`]).
			Where(`app_type`, appType).
			Where(`date`, time.Now().Format(`2006-01-02`)).
			Where(`type`, cast.ToString(StatsTypeDailyMsgCount)).
			Update(msql.Datas{
				`amount`:      cast.ToString(cast.ToInt(row[`amount`]) + 2),
				`update_time`: tool.Time2Int(),
			})
		if err != nil {
			return err
		}
	}
	return nil
}

func statDailyTokenCount(adminUserId int, robot msql.Params, appType string, promptToken, completionToken int) error {
	row, err := msql.Model(`llm_request_daily_stats`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, robot[`id`]).
		Where(`app_type`, appType).
		Where(`date`, time.Now().Format(`2006-01-02`)).
		Where(`type`, cast.ToString(StatsTypeDailyTokenCount)).
		Find()
	if err != nil {
		return err
	}
	if len(row) == 0 {
		_, err = msql.Model(`llm_request_daily_stats`, define.Postgres).
			Insert(msql.Datas{
				`admin_user_id`: cast.ToString(adminUserId),
				`robot_id`:      robot[`id`],
				`date`:          time.Now().Format(`2006-01-02`),
				`app_type`:      appType,
				`type`:          cast.ToString(StatsTypeDailyTokenCount),
				`amount`:        cast.ToString(promptToken + completionToken),
				`create_time`:   tool.Time2Int(),
				`update_time`:   tool.Time2Int(),
			})
		if err != nil {
			return err
		}
	} else {
		_, err = msql.Model(`llm_request_daily_stats`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`robot_id`, robot[`id`]).
			Where(`app_type`, appType).
			Where(`date`, time.Now().Format(`2006-01-02`)).
			Where(`type`, cast.ToString(StatsTypeDailyTokenCount)).
			Update(msql.Datas{
				"amount":      cast.ToString(cast.ToInt(row[`amount`]) + promptToken + completionToken),
				"update_time": tool.Time2Int(),
			})
		if err != nil {
			return err
		}
	}
	return nil
}

func statDailyRequestLibraryTip(adminUserId int, robot msql.Params, appType, statType string) {
	row, err := msql.Model(`llm_request_daily_stats`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, robot[`id`]).
		Where(`app_type`, appType).
		Where(`date`, time.Now().Format(`2006-01-02`)).
		Where(`type`, statType).
		Find()
	if err != nil {
		logs.Error(err.Error())
		return
	}
	if len(row) == 0 {
		_, err = msql.Model(`llm_request_daily_stats`, define.Postgres).Insert(msql.Datas{
			`admin_user_id`: cast.ToString(adminUserId),
			`robot_id`:      robot[`id`],
			`date`:          time.Now().Format(`2006-01-02`),
			`app_type`:      appType,
			`type`:          statType,
			`amount`:        1,
			`create_time`:   tool.Time2Int(),
			`update_time`:   tool.Time2Int(),
		})
		if err != nil {
			logs.Error(err.Error())
			return
		}
	} else {
		_, err = msql.Model(`llm_request_daily_stats`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`robot_id`, robot[`id`]).
			Where(`app_type`, appType).
			Where(`date`, time.Now().Format(`2006-01-02`)).
			Where(`type`, statType).
			Update(msql.Datas{
				`amount`:      cast.ToString(cast.ToInt(row[`amount`]) + 1),
				`update_time`: tool.Time2Int(),
			})
		if err != nil {
			logs.Error(err.Error())
			return
		}
	}
	return
}

func StatAiTipAnalyse(userId, robotId int, startDate, endDate, lang, channel string) (map[string]any, error) {
	types := fmt.Sprintf(`%d,%d`, StatsTypeDailyAiMsgCount, StatsTypeDailyLibraryTipCount)
	condition := fmt.Sprintf(`admin_user_id = %d and robot_id = %d and type in(%s) 
		and date >= '%s' and date <= '%s'`, userId, robotId, types, startDate, endDate)
	if len(channel) > 0 {
		condition = condition + fmt.Sprintf(`and app_type = '%s'`, channel)
	}
	m := msql.Model(`llm_request_daily_stats`, define.Postgres).
		Where(condition).
		Group(`date,type`).
		Order(`date asc`).
		Field(`date,sum(amount) as amount,type`)
	result, err := m.Select()
	if err != nil {
		logs.Error(err.Error())
		return nil, errors.New(i18n.Show(lang, `sys_err`))
	}
	list := make([]map[string]any, 0)
	libraryHitRate, messageTotal, libraryHitTotal, libraryMissCount := 0, 0, 0, 0
	fillBackDateRange(startDate, endDate, func(tempDate string) {
		boolFind := false
		msgCount := 0
		libraryHitCount := 0
		for _, resultVal := range result {
			if DbDateToDateFormat(resultVal[`date`]) == tempDate {
				if resultVal[`type`] == cast.ToString(StatsTypeDailyAiMsgCount) {
					msgCount = cast.ToInt(resultVal[`amount`])
					messageTotal += msgCount
				} else if resultVal[`type`] == cast.ToString(StatsTypeDailyLibraryTipCount) {
					libraryHitCount = cast.ToInt(resultVal[`amount`])
					libraryHitTotal += libraryHitCount
				}
				boolFind = true
			}
		}
		libraryHitRateVal := 0
		if !boolFind {
			list = append(list, map[string]any{
				`date`:               tempDate,
				`message_total`:      0,
				`library_hit_total`:  0,
				`library_miss_total`: 0,
				`library_hit_rate`:   libraryHitRateVal,
			})
		} else {
			if msgCount > 0 {
				libraryHitRateVal = int(float64(libraryHitCount) / float64(msgCount) * 100)
			}
			list = append(list, map[string]any{
				`date`:               tempDate,
				`message_total`:      msgCount,
				`library_hit_total`:  libraryHitCount,
				`library_miss_total`: msgCount - libraryHitCount,
				`library_hit_rate`:   libraryHitRateVal,
			})
		}
	})
	if messageTotal > 0 {
		libraryHitRate = int(float64(libraryHitTotal) / float64(messageTotal) * 100)
	}
	libraryMissCount = messageTotal - libraryHitTotal
	data := make(map[string]any)
	data[`chart_list`] = list
	data[`header`] = map[string]any{
		`library_hit_rate`:   libraryHitRate,
		`message_total`:      messageTotal,
		`library_hit_total`:  libraryHitTotal,
		`library_miss_total`: libraryMissCount,
	}
	return data, nil
}

func fillBackDateRange(startDate, endDate string, back func(string)) {
	tempDate := ``
	for {
		if tempDate == `` {
			tempDate = startDate
		} else {
			date, err := time.Parse("2006-01-02", tempDate)
			if err != nil {
				break
			}
			nextDate := date.AddDate(0, 0, 1)
			tempDate = nextDate.Format("2006-01-02")
		}
		back(tempDate)
		if tempDate >= endDate {
			break
		}
	}
}

func DbDateToDateFormat(date string, formats ...string) string {
	format := "2006-01-02"
	if len(formats) > 0 {
		format = formats[0]
	}
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		logs.Error(err.Error())
		return date
	}
	return t.Format(format)
}
