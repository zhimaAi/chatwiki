// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/lib_redis"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
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
	lang string,
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

	corpName := GetModelNameByDefine(lang, config[`model_define`])
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
	if err = SaveLlmRequestLogs(data); err != nil {
		logs.Error(err.Error())
		return err
	}

	if _type == `LLM` && len(robot) > 0 {
		err = statDailyRequestCount(adminUserId, robot, appType)
		if err != nil {
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

func GetTokenAppType(robot msql.Params) string {
	tokenAppType := define.TokenAppTypeOther
	if len(robot) == 0 {
		return tokenAppType
	}
	if len(robot) > 0 && cast.ToInt(robot[`id`]) > 0 {
		applicationType := cast.ToInt(robot[`application_type`])
		if applicationType == define.ApplicationTypeChat {
			tokenAppType = define.TokenAppTypeRobot
		} else if applicationType == define.ApplicationTypeFlow {
			tokenAppType = define.TokenAppTypeWorkflow
		}
	}
	return tokenAppType
}

func statTokenApp(_type string, adminUserId int, robot msql.Params, corpName, model string, promptToken, completionToken int) error {
	var oldPromptToken int
	var oldCompletionToken int
	if robot == nil {
		robot = msql.Params{}
	}
	robotId := cast.ToInt(robot[`id`])
	tokenAppType := GetTokenAppType(robot)
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
			`admin_user_id`:  cast.ToString(adminUserId),
			`token_app_type`: tokenAppType,
			`robot_id`:       cast.ToString(robotId),
			`corp`:           corpName,
			`model`:          model,
			`type`:           _type,
			`date`:           time.Now().Format(`2006-01-02`),
			`request_num`:    1,
			`create_time`:    tool.Time2Int(),
			`update_time`:    tool.Time2Int(),
		})
		if err != nil {
			return err
		}
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
	//today token app use
	//tokenUseIncr(adminUserId, tokenAppType, robotId, promptToken, completionToken, time.Now().Format("2006-01-02"))
	//token app use
	tokenAppUseIncr(adminUserId, tokenAppType, robotId, promptToken, completionToken)
	return nil
}

func tokenAppUseIncr(adminUserId int, tokenAppType string, robotId int, promptToken int, completionToken int) {
	if !TokenAppIsSwitchOpen(adminUserId, robotId, tokenAppType) {
		return
	}
	useCache := TokenAppUseCacheBuildHandler{
		AdminUserId:  adminUserId,
		TokenAppType: tokenAppType,
		RobotId:      robotId,
		DateYmd:      ``,
	}
	useToken, err := lib_redis.GetCacheIncrWithBuild(define.Redis, &useCache, time.Hour*30)
	if err != nil {
		logs.Error(`token app use cache get error：` + err.Error())
		return
	}
	_, err = msql.Model(`llm_token_app_limit`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`token_app_type`, cast.ToString(tokenAppType)).
		Where(`robot_id`, cast.ToString(robotId)).
		Update(msql.Datas{
			`use_token`: cast.ToString(cast.ToInt64(useToken) + cast.ToInt64(promptToken) + cast.ToInt64(completionToken)),
		})
	if err != nil {
		logs.Error(`token app use cache update error：` + err.Error())
		return
	}
	tokenUseIncr(adminUserId, tokenAppType, robotId, promptToken, completionToken, ``)
}

func tokenUseIncr(adminUserId int, tokenAppType string, robotId int, promptToken int, completionToken int, dateYmd string) {
	defer func() {
		if r := recover(); r != nil {
			logs.Error(`token app use cache incr panic error：` + fmt.Sprint(r))
		}
	}()
	useCache := TokenAppUseCacheBuildHandler{
		AdminUserId:  adminUserId,
		TokenAppType: tokenAppType,
		RobotId:      robotId,
		DateYmd:      dateYmd,
	}
	_, err := lib_redis.IncrCacheIncrWithBuild(define.Redis, &useCache, cast.ToInt64(promptToken+completionToken), time.Hour*30)
	if err != nil {
		logs.Error(`today token use incr error：` + err.Error())
	}
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

func GetTokenAppUseByDate(adminUserId, robotId int, tokenAppType, date string) (int64, error) {
	m := msql.Model(`llm_token_app_daily_stats`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`token_app_type`, tokenAppType).
		Where(`robot_id`, cast.ToString(robotId))
	if date == `` {
		date = time.Now().Format(`2006-01-02`)
	}
	stat, err := m.Where(`date`, date).
		Field(`sum(prompt_token) as prompt_token, sum(completion_token) as completion_token`).
		Find()
	if err != nil {
		return 0, err
	}
	return cast.ToInt64(stat[`prompt_token`]) + cast.ToInt64(stat[`completion_token`]), nil
}

func GetTokenAppLimitUse(adminUserId, robotId int, tokenAppType string) (int64, error) {
	m := msql.Model(`llm_token_app_limit`, define.Postgres)
	useToken, err := m.
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`token_app_type`, tokenAppType).
		Where(`robot_id`, cast.ToString(robotId)).
		Value(`use_token`)
	if err != nil {
		return 0, err
	}
	return cast.ToInt64(useToken), nil
}

func statDailyActiveUser(adminUserId, robotId int, appType string) error {
	count, err := msql.Model(`chat_ai_session`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, cast.ToString(robotId)).
		Where(`app_type`, appType).
		Where(`create_time`, `>=`, cast.ToString(tool.GetTimestamp(0))).
		Value(`count(distinct openid)`)
	if err != nil {
		return err
	}

	row, err := msql.Model(`llm_request_daily_stats`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, cast.ToString(robotId)).
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
			`robot_id`:      robotId,
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

func statDailyNewUser(adminUserId, robotId int, appType, openid string) error {
	oldUser, err := msql.Model(`chat_ai_session`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, cast.ToString(robotId)).
		Where(`app_type`, appType).
		Where(`openid`, openid).
		Where(`create_time`, `<`, cast.ToString(tool.GetTimestamp(0))).
		Field(`id`).Find()
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	if len(oldUser) > 0 {
		return nil
	}

	key := fmt.Sprintf(`daily_new_user_%d_%d_%s_%s`, adminUserId, robotId, appType, tool.Date(`Ymd`))
	_, err = define.Redis.SAdd(context.Background(), key, openid).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error(err.Error())
		return err
	}
	_, err = define.Redis.Expire(context.Background(), key, 24*time.Hour).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		logs.Error(err.Error())
		return err
	}
	newUserCount, err := define.Redis.SCard(context.Background(), key).Result()
	if err != nil {
		logs.Error(err.Error())
		return err
	}

	row, err := msql.Model(`llm_request_daily_stats`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, cast.ToString(robotId)).
		Where(`app_type`, appType).
		Where(`date`, time.Now().Format(`2006-01-02`)).
		Where(`type`, cast.ToString(StatsTypeDailyNewUser)).
		Find()
	if len(row) == 0 {
		_, err = msql.Model(`llm_request_daily_stats`, define.Postgres).Insert(msql.Datas{
			`admin_user_id`: adminUserId,
			`robot_id`:      robotId,
			`app_type`:      appType,
			`date`:          time.Now().Format(`2006-01-02`),
			`type`:          StatsTypeDailyNewUser,
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
			Where(`robot_id`, cast.ToString(robotId)).
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
	robotId := cast.ToInt(robot[`id`])
	row, err := msql.Model(`llm_request_daily_stats`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, cast.ToString(robotId)).
		Where(`app_type`, appType).
		Where(`date`, time.Now().Format(`2006-01-02`)).
		Where(`type`, cast.ToString(StatsTypeDailyMsgCount)).
		Find()
	if err != nil {
		return err
	}
	if len(row) == 0 {
		_, err = msql.Model(`llm_request_daily_stats`, define.Postgres).Insert(msql.Datas{
			`admin_user_id`: adminUserId,
			`robot_id`:      robotId,
			`date`:          time.Now().Format(`2006-01-02`),
			`app_type`:      appType,
			`type`:          StatsTypeDailyMsgCount,
			`amount`:        2,
			`create_time`:   tool.Time2Int(),
			`update_time`:   tool.Time2Int(),
		})
		if err != nil {
			return err
		}
	} else {
		_, err = msql.Model(`llm_request_daily_stats`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`robot_id`, cast.ToString(robotId)).
			Where(`app_type`, appType).
			Where(`date`, time.Now().Format(`2006-01-02`)).
			Where(`type`, cast.ToString(StatsTypeDailyMsgCount)).
			Update(msql.Datas{
				`amount`:      cast.ToInt(row[`amount`]) + 2,
				`update_time`: tool.Time2Int(),
			})
		if err != nil {
			return err
		}
	}
	return nil
}

func statDailyTokenCount(adminUserId int, robot msql.Params, appType string, promptToken, completionToken int) error {
	robotId := cast.ToInt(robot[`id`])
	row, err := msql.Model(`llm_request_daily_stats`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, cast.ToString(robotId)).
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
				`admin_user_id`: adminUserId,
				`robot_id`:      robotId,
				`date`:          time.Now().Format(`2006-01-02`),
				`app_type`:      appType,
				`type`:          StatsTypeDailyTokenCount,
				`amount`:        promptToken + completionToken,
				`create_time`:   tool.Time2Int(),
				`update_time`:   tool.Time2Int(),
			})
		if err != nil {
			return err
		}
	} else {
		_, err = msql.Model(`llm_request_daily_stats`, define.Postgres).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`robot_id`, cast.ToString(robotId)).
			Where(`app_type`, appType).
			Where(`date`, time.Now().Format(`2006-01-02`)).
			Where(`type`, cast.ToString(StatsTypeDailyTokenCount)).
			Update(msql.Datas{
				"amount":      cast.ToInt(row[`amount`]) + promptToken + completionToken,
				"update_time": tool.Time2Int(),
			})
		if err != nil {
			return err
		}
	}
	return nil
}

func statDailyRequestLibraryTip(adminUserId int, robot msql.Params, appType, statType string) {
	robotId := cast.ToInt(robot[`id`])
	row, err := msql.Model(`llm_request_daily_stats`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`robot_id`, cast.ToString(robotId)).
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
			`robot_id`:      robotId,
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
			Where(`robot_id`, cast.ToString(robotId)).
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
