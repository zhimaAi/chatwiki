// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"context"
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
