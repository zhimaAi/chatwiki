// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/llm_runner"
	"context"
	"errors"
	"fmt"
	"net"
	neturl "net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

var runningWebToSkillTasks sync.Map

func GetWebToSkillTaskStatusList(lang string) []map[string]any {
	return []map[string]any{
		{`status`: define.WebToSkillTaskStatusQueued, `status_name`: i18n.Show(lang, `web_to_skill_status_queued`)},
		{`status`: define.WebToSkillTaskStatusRunning, `status_name`: i18n.Show(lang, `web_to_skill_status_running`)},
		{`status`: define.WebToSkillTaskStatusSucceed, `status_name`: i18n.Show(lang, `web_to_skill_status_succeed`)},
		{`status`: define.WebToSkillTaskStatusFailed, `status_name`: i18n.Show(lang, `web_to_skill_status_failed`)},
		{`status`: define.WebToSkillTaskStatusStopping, `status_name`: i18n.Show(lang, `web_to_skill_status_stopping`)},
		{`status`: define.WebToSkillTaskStatusStopped, `status_name`: i18n.Show(lang, `web_to_skill_status_stopped`)},
	}
}

func ParseWebToSkillURLText(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == `` {
		return nil
	}
	urls := make([]string, 0)
	if strings.HasPrefix(raw, `[`) {
		if err := tool.JsonDecodeUseNumber(raw, &urls); err == nil {
			return urls
		}
	}
	raw = strings.NewReplacer("\r\n", "\n", "\r", "\n").Replace(raw)
	parts := strings.Split(raw, "\n")
	if len(parts) == 1 {
		parts = strings.Split(raw, `,`)
	}
	for _, part := range parts {
		if part = strings.TrimSpace(part); part != `` {
			urls = append(urls, part)
		}
	}
	return urls
}

func GetWebToSkillTaskList(lang string, adminUserId int, filter define.WebToSkillTaskListFilter) (map[string]any, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Size <= 0 {
		filter.Size = define.WebToSkillTaskDefaultPageSize
	}
	m := msql.Model(define.TableWebToSkillTask, define.Postgres)
	m.Where(`admin_user_id`, cast.ToString(adminUserId)).
		Field(`id,task_batch,urls,custom_prompt,model_config_id,use_model,temperature,max_token,status,skill_name,skill_description,file_name,file_url,file_size,err_msg,start_time,end_time,create_time,update_time`)
	if filter.Status != -1 {
		m.Where(`status`, cast.ToString(filter.Status))
	}
	list, total, err := m.Order(`id desc`).Paginate(filter.Page, filter.Size)
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return nil, errors.New(i18n.Show(lang, `sys_err`))
	}
	items := make([]define.WebToSkillTaskItem, 0, len(list))
	for _, row := range list {
		items = append(items, buildWebToSkillTaskItem(row, false))
	}
	return map[string]any{
		`status_map`: GetWebToSkillTaskStatusList(lang),
		`list`:       items,
		`total`:      total,
		`page`:       filter.Page,
		`size`:       filter.Size,
	}, nil
}

func CreateWebToSkillTask(lang string, adminUserId int, params define.WebToSkillTaskCreateParams) (int64, error) {
	temperature := float32(define.WebToSkillTaskDefaultTemp)
	if params.Temperature != nil {
		temperature = *params.Temperature
	}
	maxToken := define.WebToSkillTaskDefaultMaxToken
	if params.MaxToken != nil {
		maxToken = *params.MaxToken
	}
	if err := checkWebToSkillCreateParams(lang, adminUserId, params.Urls, params.ModelConfigId, params.UseModel, temperature, maxToken); err != nil {
		return 0, err
	}
	urlsJson, err := tool.JsonEncode(params.Urls)
	if err != nil {
		logs.Error(err.Error())
		return 0, errors.New(i18n.Show(lang, `sys_err`))
	}
	now := tool.Time2Int()
	taskBatch := uuid.NewString()
	m := msql.Model(define.TableWebToSkillTask, define.Postgres)
	id, err := m.Insert(msql.Datas{
		`admin_user_id`:   adminUserId,
		`task_batch`:      taskBatch,
		`urls`:            urlsJson,
		`custom_prompt`:   strings.TrimSpace(params.CustomPrompt),
		`model_config_id`: params.ModelConfigId,
		`use_model`:       strings.TrimSpace(params.UseModel),
		`temperature`:     temperature,
		`max_token`:       maxToken,
		`status`:          define.WebToSkillTaskStatusQueued,
		`create_time`:     now,
		`update_time`:     now,
	}, `id`)
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return 0, errors.New(i18n.Show(lang, `sys_err`))
	}
	if err = enqueueWebToSkillTask(id); err != nil {
		m = msql.Model(define.TableWebToSkillTask, define.Postgres)
		if _, updateErr := m.Where(`id`, cast.ToString(id)).Update(msql.Datas{
			`status`:      define.WebToSkillTaskStatusFailed,
			`err_msg`:     err.Error(),
			`update_time`: tool.Time2Int(),
			`end_time`:    tool.Time2Int(),
		}); updateErr != nil {
			logs.Error(`sql:%s,err:%s`, m.GetLastSql(), updateErr.Error())
		}
		return 0, errors.New(i18n.Show(lang, `sys_err`))
	}
	return id, nil
}

func StopWebToSkillTask(lang string, adminUserId int, id int64) (map[string]any, error) {
	info, err := getWebToSkillTaskRow(adminUserId, id)
	if err != nil {
		return nil, errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(info) == 0 {
		return nil, errors.New(i18n.Show(lang, `no_data`))
	}
	status := cast.ToInt(info[`status`])
	if status == define.WebToSkillTaskStatusStopped {
		return map[string]any{`stopped`: true, `status`: status}, nil
	}
	if status != define.WebToSkillTaskStatusQueued &&
		status != define.WebToSkillTaskStatusRunning &&
		status != define.WebToSkillTaskStatusStopping {
		return map[string]any{`stopped`: false, `status`: status}, nil
	}
	stopKey := GetWebToSkillTaskStopKey(id)
	if err = define.Redis.Set(context.Background(), stopKey, `1`, define.WebToSkillTaskStopKeyTTL).Err(); err != nil {
		logs.Error(err.Error())
		return nil, errors.New(i18n.Show(lang, `sys_err`))
	}
	if status != define.WebToSkillTaskStatusStopping {
		m := msql.Model(define.TableWebToSkillTask, define.Postgres)
		affected, updateErr := m.Where(`id`, cast.ToString(id)).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`status`, `in`, fmt.Sprintf(`%d,%d`, define.WebToSkillTaskStatusQueued, define.WebToSkillTaskStatusRunning)).
			Update(msql.Datas{`status`: define.WebToSkillTaskStatusStopping, `update_time`: tool.Time2Int()})
		if updateErr != nil {
			logs.Error(`sql:%s,err:%s`, m.GetLastSql(), updateErr.Error())
			return nil, errors.New(i18n.Show(lang, `sys_err`))
		}
		if affected == 0 {
			latest, queryErr := getWebToSkillTaskRow(adminUserId, id)
			if queryErr != nil {
				return nil, errors.New(i18n.Show(lang, `sys_err`))
			}
			latestStatus := cast.ToInt(latest[`status`])
			if latestStatus != define.WebToSkillTaskStatusQueued &&
				latestStatus != define.WebToSkillTaskStatusRunning &&
				latestStatus != define.WebToSkillTaskStatusStopping {
				if latestStatus != define.WebToSkillTaskStatusStopped {
					define.Redis.Del(context.Background(), stopKey)
				}
				return map[string]any{`stopped`: latestStatus == define.WebToSkillTaskStatusStopped, `status`: latestStatus}, nil
			}
		}
	}
	cancelResp := llm_runner.RpcCancelRun(
		define.Config.WebService[`llm_runner_host`],
		GetWebToSkillTaskRunID(info[`task_batch`]),
	)
	if cancelResp.ErrorMsg != `` {
		logs.Error(`cancel web-to-skill runner,task:%d,run_id:%s,err:%s`, id, GetWebToSkillTaskRunID(info[`task_batch`]), cancelResp.ErrorMsg)
	}
	return map[string]any{`stopped`: false, `status`: define.WebToSkillTaskStatusStopping}, nil
}

func RegenerateWebToSkillTask(lang string, adminUserId int, id int64) (int64, error) {
	info, err := getWebToSkillTaskRow(adminUserId, id)
	if err != nil {
		return 0, errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(info) == 0 {
		return 0, errors.New(i18n.Show(lang, `no_data`))
	}
	status := cast.ToInt(info[`status`])
	if status == define.WebToSkillTaskStatusQueued ||
		status == define.WebToSkillTaskStatusRunning ||
		status == define.WebToSkillTaskStatusStopping {
		return 0, errors.New(i18n.Show(lang, `task_running`))
	}
	urls := decodeWebToSkillURLs(info[`urls`])
	if err = checkWebToSkillCreateParams(lang, adminUserId, urls, cast.ToInt(info[`model_config_id`]), info[`use_model`],
		cast.ToFloat32(info[`temperature`]), cast.ToInt(info[`max_token`])); err != nil {
		return 0, err
	}
	define.Redis.Del(context.Background(), GetWebToSkillTaskStopKey(id))
	now := tool.Time2Int()
	m := msql.Model(define.TableWebToSkillTask, define.Postgres)
	affected, err := m.Where(`id`, cast.ToString(id)).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`status`, `in`, fmt.Sprintf(`%d,%d,%d`, define.WebToSkillTaskStatusSucceed, define.WebToSkillTaskStatusFailed, define.WebToSkillTaskStatusStopped)).
		Update(msql.Datas{
			`task_batch`:        uuid.NewString(),
			`status`:            define.WebToSkillTaskStatusQueued,
			`skill_name`:        ``,
			`skill_description`: ``,
			`file_name`:         ``,
			`file_url`:          ``,
			`file_size`:         0,
			`debug_log`:         ``,
			`err_msg`:           ``,
			`start_time`:        0,
			`end_time`:          0,
			`update_time`:       now,
		})
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return 0, errors.New(i18n.Show(lang, `sys_err`))
	}
	if affected == 0 {
		return 0, errors.New(i18n.Show(lang, `task_running`))
	}
	if err = enqueueWebToSkillTask(id); err != nil {
		m = msql.Model(define.TableWebToSkillTask, define.Postgres)
		if _, updateErr := m.Where(`id`, cast.ToString(id)).Where(`admin_user_id`, cast.ToString(adminUserId)).
			Update(msql.Datas{
				`status`:      define.WebToSkillTaskStatusFailed,
				`err_msg`:     err.Error(),
				`update_time`: tool.Time2Int(),
				`end_time`:    tool.Time2Int(),
			}); updateErr != nil {
			logs.Error(`sql:%s,err:%s`, m.GetLastSql(), updateErr.Error())
		}
		return 0, errors.New(i18n.Show(lang, `sys_err`))
	}
	return id, nil
}

func GetWebToSkillTaskDetail(lang string, adminUserId int, id int64) (*define.WebToSkillTaskItem, error) {
	info, err := getWebToSkillTaskRow(adminUserId, id)
	if err != nil {
		return nil, errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(info) == 0 {
		return nil, errors.New(i18n.Show(lang, `no_data`))
	}
	item := buildWebToSkillTaskItem(info, true)
	return &item, nil
}

func GetWebToSkillTaskDownloadFile(lang string, adminUserId int, id int64) (string, string, error) {
	info, err := getWebToSkillTaskRow(adminUserId, id)
	if err != nil {
		return ``, ``, errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(info) == 0 || cast.ToInt(info[`status`]) != define.WebToSkillTaskStatusSucceed || len(info[`file_url`]) == 0 {
		return ``, ``, errors.New(i18n.Show(lang, `no_data`))
	}
	file := GetFileByLink(info[`file_url`])
	if file == `` {
		return ``, ``, errors.New(i18n.Show(lang, `no_data`))
	}
	fileName := info[`file_name`]
	if fileName == `` {
		fileName = fmt.Sprintf(`web-to-skill-%d.zip`, id)
	}
	return file, fileName, nil
}

func InstallWebToSkillTask(lang string, adminUserId int, id int64, overwrite bool) (*ClawbotUserSkillItem, error) {
	info, err := getWebToSkillTaskRow(adminUserId, id)
	if err != nil {
		return nil, errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(info) == 0 || cast.ToInt(info[`status`]) != define.WebToSkillTaskStatusSucceed ||
		info[`file_url`] == `` || info[`skill_name`] == `` || info[`skill_description`] == `` {
		return nil, errors.New(i18n.Show(lang, `no_data`))
	}
	file := GetFileByLink(info[`file_url`])
	if file == `` {
		return nil, errors.New(i18n.Show(lang, `no_data`))
	}
	fileName := info[`file_name`]
	if fileName == `` {
		fileName = fmt.Sprintf(`web-to-skill-%d.zip`, id)
	}
	item, errKey, err := InstallUserClawbotSkillZip(
		adminUserId,
		fileName,
		file,
		info[`skill_name`],
		info[`skill_description`],
		overwrite,
	)
	if err != nil {
		logs.Error(err.Error())
		return nil, errors.New(i18n.Show(lang, `sys_err`))
	}
	if errKey != `` {
		return nil, errors.New(i18n.Show(lang, errKey))
	}
	return item, nil
}

func RunWebToSkillTask(id int64) error {
	if id <= 0 {
		return errors.New(`invalid web-to-skill task id`)
	}
	if _, loaded := runningWebToSkillTasks.LoadOrStore(id, struct{}{}); loaded {
		return nil
	}
	defer runningWebToSkillTasks.Delete(id)

	info, err := getWebToSkillTaskRow(0, id)
	if err != nil {
		return err
	}
	if len(info) == 0 {
		return nil
	}
	stopKey := GetWebToSkillTaskStopKey(id)
	status := cast.ToInt(info[`status`])
	if status == define.WebToSkillTaskStatusStopping || IsWebToSkillTaskStopped(stopKey) {
		return setWebToSkillTaskStopped(id, info[`task_batch`], nil)
	}
	if status == define.WebToSkillTaskStatusRunning {
		return finishWebToSkillTask(id, define.WebToSkillTaskStatusFailed, msql.Datas{
			`err_msg`: `Task interrupted by service restart or failure; please retry manually.`,
		})
	}
	if status != define.WebToSkillTaskStatusQueued {
		return nil
	}
	m := msql.Model(define.TableWebToSkillTask, define.Postgres)
	affected, err := m.Where(`id`, cast.ToString(id)).
		Where(`status`, cast.ToString(define.WebToSkillTaskStatusQueued)).
		Update(msql.Datas{
			`status`:      define.WebToSkillTaskStatusRunning,
			`start_time`:  tool.Time2Int(),
			`update_time`: tool.Time2Int(),
		})
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return err
	}
	if affected == 0 {
		latest, queryErr := getWebToSkillTaskRow(0, id)
		if queryErr != nil {
			return queryErr
		}
		if cast.ToInt(latest[`status`]) == define.WebToSkillTaskStatusStopping || IsWebToSkillTaskStopped(stopKey) {
			return setWebToSkillTaskStopped(id, info[`task_batch`], nil)
		}
		return nil
	}
	defer cleanupWebToSkillTaskWorkDir(info[`task_batch`])
	if IsWebToSkillTaskStopped(stopKey) {
		return setWebToSkillTaskStopped(id, info[`task_batch`], nil)
	}
	result, runErr := DoWebToSkill(define.LangZhCn, WebToSkillTaskInfo{
		TaskBatch:     info[`task_batch`],
		AdminUserId:   cast.ToInt(info[`admin_user_id`]),
		ModelConfigId: cast.ToInt(info[`model_config_id`]),
		UseModel:      info[`use_model`],
		Temperature:   cast.ToFloat32(info[`temperature`]),
		MaxToken:      cast.ToInt(info[`max_token`]),
		Urls:          decodeWebToSkillURLs(info[`urls`]),
		CustomPrompt:  info[`custom_prompt`],
		StopKey:       stopKey,
	})
	if IsWebToSkillTaskStopped(stopKey) {
		return setWebToSkillTaskStopped(id, info[`task_batch`], result.DebugLog)
	}
	if runErr != nil {
		return finishWebToSkillTask(id, define.WebToSkillTaskStatusFailed, msql.Datas{
			`debug_log`: tool.JsonEncodeNoError(result.DebugLog),
			`err_msg`:   runErr.Error(),
		})
	}
	if result.ZipPath == `` {
		return finishWebToSkillTask(id, define.WebToSkillTaskStatusFailed, msql.Datas{
			`debug_log`: tool.JsonEncodeNoError(result.DebugLog),
			`err_msg`:   `generated zip path is empty`,
		})
	}
	skillName, skillDescription, metaErr := ReadClawbotSkillZipMeta(result.ZipPath)
	if metaErr != nil {
		return finishWebToSkillTask(id, define.WebToSkillTaskStatusFailed, msql.Datas{
			`debug_log`: tool.JsonEncodeNoError(result.DebugLog),
			`err_msg`:   metaErr.Error(),
		})
	}
	uploadInfo, saveErr := saveWebToSkillZipFile(cast.ToInt(info[`admin_user_id`]), result.ZipPath)
	if IsWebToSkillTaskStopped(stopKey) {
		return setWebToSkillTaskStopped(id, info[`task_batch`], result.DebugLog)
	}
	if saveErr != nil {
		return finishWebToSkillTask(id, define.WebToSkillTaskStatusFailed, msql.Datas{
			`debug_log`: tool.JsonEncodeNoError(result.DebugLog),
			`err_msg`:   saveErr.Error(),
		})
	}
	return finishWebToSkillTask(id, define.WebToSkillTaskStatusSucceed, msql.Datas{
		`skill_name`:        skillName,
		`skill_description`: skillDescription,
		`file_name`:         uploadInfo.Name,
		`file_url`:          uploadInfo.Link,
		`file_size`:         uploadInfo.Size,
		`debug_log`:         tool.JsonEncodeNoError(result.DebugLog),
		`err_msg`:           `SUCCEED`,
	})
}

func GetWebToSkillTaskStopKey(id int64) string {
	return define.WebToSkillTaskStopKeyPrefix + cast.ToString(id)
}

func GetWebToSkillTaskRunID(taskBatch string) string {
	return `web-to-skill:` + strings.TrimSpace(taskBatch)
}

func IsWebToSkillTaskStopped(stopKey string) bool {
	if stopKey == `` || define.Redis == nil {
		return false
	}
	exists, err := define.Redis.Exists(context.Background(), stopKey).Result()
	if err != nil {
		logs.Error(err.Error())
		return false
	}
	return exists > 0
}

func getWebToSkillTaskRow(adminUserId int, id int64) (msql.Params, error) {
	m := msql.Model(define.TableWebToSkillTask, define.Postgres)
	m.Where(`id`, cast.ToString(id))
	if adminUserId > 0 {
		m.Where(`admin_user_id`, cast.ToString(adminUserId))
	}
	info, err := m.Find()
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
	}
	return info, err
}

func buildWebToSkillTaskItem(row msql.Params, withLog bool) define.WebToSkillTaskItem {
	item := define.WebToSkillTaskItem{
		ID:               cast.ToInt64(row[`id`]),
		TaskBatch:        row[`task_batch`],
		Urls:             decodeWebToSkillURLs(row[`urls`]),
		CustomPrompt:     row[`custom_prompt`],
		ModelConfigId:    cast.ToInt(row[`model_config_id`]),
		UseModel:         row[`use_model`],
		Temperature:      cast.ToFloat32(row[`temperature`]),
		MaxToken:         cast.ToInt(row[`max_token`]),
		Status:           cast.ToInt(row[`status`]),
		SkillName:        row[`skill_name`],
		SkillDescription: row[`skill_description`],
		FileName:         row[`file_name`],
		FileUrl:          row[`file_url`],
		FileSize:         cast.ToInt(row[`file_size`]),
		ErrMsg:           row[`err_msg`],
		StartTime:        cast.ToInt(row[`start_time`]),
		EndTime:          cast.ToInt(row[`end_time`]),
		CreateTime:       cast.ToInt(row[`create_time`]),
		UpdateTime:       cast.ToInt(row[`update_time`]),
	}
	if withLog && row[`debug_log`] != `` {
		_ = tool.JsonDecodeUseNumber(row[`debug_log`], &item.DebugLog)
	}
	return item
}

func decodeWebToSkillURLs(raw string) []string {
	urls := make([]string, 0)
	if raw != `` {
		_ = tool.JsonDecodeUseNumber(raw, &urls)
	}
	return urls
}

func checkWebToSkillCreateParams(lang string, adminUserId int, urls []string, modelConfigId int, useModel string, temperature float32, maxToken int) error {
	if !CheckModelIsValid(adminUserId, modelConfigId, strings.TrimSpace(useModel), Llm) {
		return errors.New(i18n.Show(lang, `param_invalid`, `use_model`))
	}
	if temperature < 0 || temperature > 2 {
		return errors.New(i18n.Show(lang, `param_invalid`, `temperature`))
	}
	if maxToken <= 0 {
		return errors.New(i18n.Show(lang, `param_invalid`, `max_token`))
	}
	if len(urls) == 0 {
		return errors.New(i18n.Show(lang, `param_invalid`, `urls`))
	}
	for i, item := range urls {
		item = strings.TrimSpace(item)
		if item == `` {
			return errors.New(i18n.Show(lang, `param_invalid`, `urls`))
		}
		if err := validateWebToSkillPublicURL(lang, item); err != nil {
			logs.Error(`invalid web-to-skill url,err:%s`, err.Error())
			return err
		}
		urls[i] = item
	}
	return nil
}

func validateWebToSkillPublicURL(lang, rawURL string) error {
	parsed, err := neturl.ParseRequestURI(rawURL)
	if err != nil || parsed == nil || (parsed.Scheme != `http` && parsed.Scheme != `https`) {
		return fmt.Errorf(`%s:%s`, rawURL, i18n.Show(lang, `invalid_url`))
	}
	host := parsed.Hostname()
	if host == `` {
		return fmt.Errorf(`%s:%s`, rawURL, i18n.Show(lang, `web_to_skill_url_host_required`))
	}
	if ip := net.ParseIP(host); ip != nil {
		if isBlockedWebToSkillIP(ip) {
			return fmt.Errorf(`%s:%s`, rawURL, i18n.Show(lang, `web_to_skill_non_public_address`, ip.String()))
		}
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), define.WebToSkillTaskDNSLookupTimeout)
	defer cancel()
	addresses, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return fmt.Errorf(`%s:%s`, rawURL, i18n.Show(lang, `web_to_skill_dns_lookup_failed`, err.Error()))
	}
	if len(addresses) == 0 {
		return fmt.Errorf(`%s:%s`, rawURL, i18n.Show(lang, `web_to_skill_dns_no_addresses`))
	}
	for _, address := range addresses {
		if isBlockedWebToSkillIP(address.IP) {
			return fmt.Errorf(`%s:%s`, rawURL, i18n.Show(lang, `web_to_skill_non_public_address`, address.IP.String()))
		}
	}
	return nil
}

func isBlockedWebToSkillIP(ip net.IP) bool {
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() || ip.IsMulticast()
}

func enqueueWebToSkillTask(id int64) error {
	if err := AddJobs(define.WebToSkillTaskTopic, cast.ToString(id)); err != nil {
		logs.Error(err.Error())
		return err
	}
	return nil
}

func setWebToSkillTaskStopped(id int64, taskBatch string, debugLog []any) error {
	defer cleanupWebToSkillTaskWorkDir(taskBatch)
	data := msql.Datas{
		`status`:      define.WebToSkillTaskStatusStopped,
		`err_msg`:     `STOPPED`,
		`end_time`:    tool.Time2Int(),
		`update_time`: tool.Time2Int(),
	}
	if debugLog != nil {
		data[`debug_log`] = tool.JsonEncodeNoError(debugLog)
	}
	m := msql.Model(define.TableWebToSkillTask, define.Postgres)
	_, err := m.Where(`id`, cast.ToString(id)).
		Where(`status`, `in`, fmt.Sprintf(`%d,%d,%d`, define.WebToSkillTaskStatusQueued, define.WebToSkillTaskStatusRunning, define.WebToSkillTaskStatusStopping)).
		Update(data)
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return err
	}
	return nil
}

func finishWebToSkillTask(id int64, status int, data msql.Datas) error {
	data[`status`] = status
	data[`end_time`] = tool.Time2Int()
	data[`update_time`] = tool.Time2Int()
	m := msql.Model(define.TableWebToSkillTask, define.Postgres)
	affected, err := m.Where(`id`, cast.ToString(id)).
		Where(`status`, cast.ToString(define.WebToSkillTaskStatusRunning)).
		Update(data)
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return err
	}
	if affected == 0 {
		info, queryErr := getWebToSkillTaskRow(0, id)
		if queryErr != nil {
			return queryErr
		}
		if cast.ToInt(info[`status`]) == define.WebToSkillTaskStatusStopping {
			debugLog := make([]any, 0)
			if raw := cast.ToString(data[`debug_log`]); raw != `` {
				_ = tool.JsonDecodeUseNumber(raw, &debugLog)
			}
			return setWebToSkillTaskStopped(id, info[`task_batch`], debugLog)
		}
	}
	return nil
}

func normalizeWebToSkillZipPath(raw, taskBatch string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.Trim(raw, "`'\" \n\r\t")
	for _, line := range strings.Split(raw, "\n") {
		line = strings.Trim(line, "`'\" \r\t")
		if strings.HasSuffix(strings.ToLower(line), `.zip`) {
			raw = line
			break
		}
	}
	raw = strings.TrimPrefix(raw, `/workspace/`)
	raw = strings.TrimPrefix(raw, `workspace/`)
	if !strings.HasSuffix(strings.ToLower(raw), `.zip`) {
		return ``
	}
	workDir := strings.ReplaceAll(define.WebToSkillWorkDir, `<task_batch>`, taskBatch)
	clean := filepath.Clean(filepath.FromSlash(raw))
	workDir = filepath.Clean(filepath.FromSlash(workDir))
	rel, err := filepath.Rel(workDir, clean)
	if err != nil || rel == `.` || strings.HasPrefix(rel, `..`) {
		return ``
	}
	if !tool.IsFile(clean) {
		return ``
	}
	return clean
}

func cleanupWebToSkillTaskWorkDir(taskBatch string) {
	workDir := strings.ReplaceAll(define.WebToSkillWorkDir, `<task_batch>`, taskBatch)
	if tool.IsDir(workDir) {
		if err := os.RemoveAll(workDir); err != nil {
			logs.Error(`remove web-to-skill work dir:%s,err:%s`, workDir, err.Error())
		}
	}
}

func saveWebToSkillZipFile(adminUserId int, zipPath string) (*define.UploadInfo, error) {
	info, err := os.Stat(zipPath)
	if err != nil {
		return nil, err
	}
	if info.IsDir() || info.Size() <= 0 {
		return nil, errors.New(`invalid zip file size`)
	}
	ext := strings.ToLower(strings.TrimLeft(filepath.Ext(zipPath), `.`))
	if !tool.InArrayString(ext, define.WebToSkillTaskZipAllowExt) {
		return nil, errors.New(ext + ` not allow`)
	}
	bs, err := os.ReadFile(zipPath)
	if err != nil {
		return nil, err
	}
	if len(bs) == 0 {
		return nil, errors.New(`file content is empty`)
	}
	content := string(bs)
	objectKey := fmt.Sprintf(`chat_ai/%v/%s/%s/%s.%s`, adminUserId, `web_to_skill`, tool.Date(`Ym`), tool.MD5(content), ext)
	link, err := WriteFileByString(objectKey, content)
	if err != nil {
		return nil, err
	}
	return &define.UploadInfo{Name: filepath.Base(zipPath), Size: info.Size(), Ext: ext, Link: link}, nil
}
