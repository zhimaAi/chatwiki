// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/llm_runner"
	"context"
	"errors"
	"fmt"
	"io"
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

var runningDocToSkillTasks sync.Map

func GetDocToSkillTaskStatusList(lang string) []map[string]any {
	return []map[string]any{
		{`status`: define.DocToSkillTaskStatusQueued, `status_name`: i18n.Show(lang, `doc_to_skill_status_queued`)},
		{`status`: define.DocToSkillTaskStatusRunning, `status_name`: i18n.Show(lang, `doc_to_skill_status_running`)},
		{`status`: define.DocToSkillTaskStatusSucceed, `status_name`: i18n.Show(lang, `doc_to_skill_status_succeed`)},
		{`status`: define.DocToSkillTaskStatusFailed, `status_name`: i18n.Show(lang, `doc_to_skill_status_failed`)},
		{`status`: define.DocToSkillTaskStatusStopping, `status_name`: i18n.Show(lang, `doc_to_skill_status_stopping`)},
		{`status`: define.DocToSkillTaskStatusStopped, `status_name`: i18n.Show(lang, `doc_to_skill_status_stopped`)},
	}
}

func GetDocToSkillTaskList(lang string, adminUserId int, filter define.DocToSkillTaskListFilter) (map[string]any, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Size <= 0 {
		filter.Size = define.DocToSkillTaskDefaultPageSize
	}
	m := msql.Model(define.TableDocToSkillTask, define.Postgres)
	m.Where(`admin_user_id`, cast.ToString(adminUserId)).
		Field(`id,task_batch,source_files,custom_prompt,model_config_id,use_model,temperature,max_token,status,skill_name,skill_description,file_name,file_url,file_size,err_msg,start_time,end_time,create_time,update_time`)
	if filter.Status != -1 {
		m.Where(`status`, cast.ToString(filter.Status))
	}
	list, total, err := m.Order(`id desc`).Paginate(filter.Page, filter.Size)
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return nil, errors.New(i18n.Show(lang, `sys_err`))
	}
	items := make([]define.DocToSkillTaskItem, 0, len(list))
	for _, row := range list {
		items = append(items, buildDocToSkillTaskItem(row, false))
	}
	return map[string]any{
		`status_map`: GetDocToSkillTaskStatusList(lang),
		`list`:       items,
		`total`:      total,
		`page`:       filter.Page,
		`size`:       filter.Size,
	}, nil
}

func CreateDocToSkillTask(lang string, adminUserId int, params define.DocToSkillTaskCreateParams, sourceFiles []*define.UploadInfo) (int64, error) {
	temperature := float32(define.DocToSkillTaskDefaultTemp)
	if params.Temperature != nil {
		temperature = *params.Temperature
	}
	maxToken := define.DocToSkillTaskDefaultMaxToken
	if params.MaxToken != nil {
		maxToken = *params.MaxToken
	}
	if err := checkDocToSkillCreateParams(lang, adminUserId, sourceFiles, params.ModelConfigId, params.UseModel, temperature, maxToken); err != nil {
		return 0, err
	}
	sourceJson, err := tool.JsonEncode(sourceFiles)
	if err != nil {
		logs.Error(err.Error())
		return 0, errors.New(i18n.Show(lang, `sys_err`))
	}
	now := tool.Time2Int()
	m := msql.Model(define.TableDocToSkillTask, define.Postgres)
	id, err := m.Insert(msql.Datas{
		`admin_user_id`:   adminUserId,
		`task_batch`:      uuid.NewString(),
		`source_files`:    sourceJson,
		`custom_prompt`:   strings.TrimSpace(params.CustomPrompt),
		`model_config_id`: params.ModelConfigId,
		`use_model`:       strings.TrimSpace(params.UseModel),
		`temperature`:     temperature,
		`max_token`:       maxToken,
		`status`:          define.DocToSkillTaskStatusQueued,
		`create_time`:     now,
		`update_time`:     now,
	}, `id`)
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return 0, errors.New(i18n.Show(lang, `sys_err`))
	}
	if err = enqueueDocToSkillTask(id); err != nil {
		markDocToSkillEnqueueFailed(id, err)
		return 0, errors.New(i18n.Show(lang, `sys_err`))
	}
	return id, nil
}

func StopDocToSkillTask(lang string, adminUserId int, id int64) (map[string]any, error) {
	info, err := getDocToSkillTaskRow(adminUserId, id)
	if err != nil {
		return nil, errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(info) == 0 {
		return nil, errors.New(i18n.Show(lang, `no_data`))
	}
	status := cast.ToInt(info[`status`])
	if status == define.DocToSkillTaskStatusStopped {
		return map[string]any{`stopped`: true, `status`: status}, nil
	}
	if status != define.DocToSkillTaskStatusQueued && status != define.DocToSkillTaskStatusRunning && status != define.DocToSkillTaskStatusStopping {
		return map[string]any{`stopped`: false, `status`: status}, nil
	}
	stopKey := GetDocToSkillTaskStopKey(id)
	if err = define.Redis.Set(context.Background(), stopKey, `1`, define.DocToSkillTaskStopKeyTTL).Err(); err != nil {
		logs.Error(err.Error())
		return nil, errors.New(i18n.Show(lang, `sys_err`))
	}
	if status != define.DocToSkillTaskStatusStopping {
		m := msql.Model(define.TableDocToSkillTask, define.Postgres)
		affected, updateErr := m.Where(`id`, cast.ToString(id)).
			Where(`admin_user_id`, cast.ToString(adminUserId)).
			Where(`status`, `in`, fmt.Sprintf(`%d,%d`, define.DocToSkillTaskStatusQueued, define.DocToSkillTaskStatusRunning)).
			Update(msql.Datas{`status`: define.DocToSkillTaskStatusStopping, `update_time`: tool.Time2Int()})
		if updateErr != nil {
			logs.Error(`sql:%s,err:%s`, m.GetLastSql(), updateErr.Error())
			return nil, errors.New(i18n.Show(lang, `sys_err`))
		}
		if affected == 0 {
			latest, queryErr := getDocToSkillTaskRow(adminUserId, id)
			if queryErr != nil {
				return nil, errors.New(i18n.Show(lang, `sys_err`))
			}
			latestStatus := cast.ToInt(latest[`status`])
			if latestStatus != define.DocToSkillTaskStatusQueued && latestStatus != define.DocToSkillTaskStatusRunning && latestStatus != define.DocToSkillTaskStatusStopping {
				if latestStatus != define.DocToSkillTaskStatusStopped {
					define.Redis.Del(context.Background(), stopKey)
				}
				return map[string]any{`stopped`: latestStatus == define.DocToSkillTaskStatusStopped, `status`: latestStatus}, nil
			}
		}
	}
	cancelResp := llm_runner.RpcCancelRun(define.Config.WebService[`llm_runner_host`], GetDocToSkillTaskRunID(info[`task_batch`]))
	if cancelResp.ErrorMsg != `` {
		logs.Error(`cancel doc-to-skill runner,task:%d,run_id:%s,err:%s`, id, GetDocToSkillTaskRunID(info[`task_batch`]), cancelResp.ErrorMsg)
	}
	return map[string]any{`stopped`: false, `status`: define.DocToSkillTaskStatusStopping}, nil
}

func RegenerateDocToSkillTask(lang string, adminUserId int, id int64) (int64, error) {
	info, err := getDocToSkillTaskRow(adminUserId, id)
	if err != nil {
		return 0, errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(info) == 0 {
		return 0, errors.New(i18n.Show(lang, `no_data`))
	}
	status := cast.ToInt(info[`status`])
	if status == define.DocToSkillTaskStatusQueued || status == define.DocToSkillTaskStatusRunning || status == define.DocToSkillTaskStatusStopping {
		return 0, errors.New(i18n.Show(lang, `task_running`))
	}
	sourceFiles := decodeDocToSkillSourceFiles(info[`source_files`])
	if err = checkDocToSkillCreateParams(lang, adminUserId, sourceFiles, cast.ToInt(info[`model_config_id`]), info[`use_model`], cast.ToFloat32(info[`temperature`]), cast.ToInt(info[`max_token`])); err != nil {
		return 0, err
	}
	define.Redis.Del(context.Background(), GetDocToSkillTaskStopKey(id))
	now := tool.Time2Int()
	m := msql.Model(define.TableDocToSkillTask, define.Postgres)
	affected, err := m.Where(`id`, cast.ToString(id)).
		Where(`admin_user_id`, cast.ToString(adminUserId)).
		Where(`status`, `in`, fmt.Sprintf(`%d,%d,%d`, define.DocToSkillTaskStatusSucceed, define.DocToSkillTaskStatusFailed, define.DocToSkillTaskStatusStopped)).
		Update(msql.Datas{
			`task_batch`:        uuid.NewString(),
			`status`:            define.DocToSkillTaskStatusQueued,
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
	if err = enqueueDocToSkillTask(id); err != nil {
		markDocToSkillEnqueueFailed(id, err)
		return 0, errors.New(i18n.Show(lang, `sys_err`))
	}
	return id, nil
}

func GetDocToSkillTaskDetail(lang string, adminUserId int, id int64) (*define.DocToSkillTaskItem, error) {
	info, err := getDocToSkillTaskRow(adminUserId, id)
	if err != nil {
		return nil, errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(info) == 0 {
		return nil, errors.New(i18n.Show(lang, `no_data`))
	}
	item := buildDocToSkillTaskItem(info, true)
	return &item, nil
}

func GetDocToSkillTaskDownloadFile(lang string, adminUserId int, id int64) (string, string, error) {
	info, err := getDocToSkillTaskRow(adminUserId, id)
	if err != nil {
		return ``, ``, errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(info) == 0 || cast.ToInt(info[`status`]) != define.DocToSkillTaskStatusSucceed || info[`file_url`] == `` {
		return ``, ``, errors.New(i18n.Show(lang, `no_data`))
	}
	file := GetFileByLink(info[`file_url`])
	if file == `` {
		return ``, ``, errors.New(i18n.Show(lang, `no_data`))
	}
	fileName := info[`file_name`]
	if fileName == `` {
		fileName = fmt.Sprintf(`doc-to-skill-%d.zip`, id)
	}
	return file, fileName, nil
}

func InstallDocToSkillTask(lang string, adminUserId int, id int64, overwrite bool) (*ClawbotUserSkillItem, error) {
	info, err := getDocToSkillTaskRow(adminUserId, id)
	if err != nil {
		return nil, errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(info) == 0 || cast.ToInt(info[`status`]) != define.DocToSkillTaskStatusSucceed || info[`file_url`] == `` || info[`skill_name`] == `` || info[`skill_description`] == `` {
		return nil, errors.New(i18n.Show(lang, `no_data`))
	}
	if cast.ToInt64(info[`file_size`]) > int64(define.MaxSkillZipSize) {
		return nil, errors.New(i18n.Show(lang, `clawbot_skill_zip_too_big`))
	}
	file := GetFileByLink(info[`file_url`])
	if file == `` {
		return nil, errors.New(i18n.Show(lang, `no_data`))
	}
	fileName := info[`file_name`]
	if fileName == `` {
		fileName = fmt.Sprintf(`doc-to-skill-%d.zip`, id)
	}
	item, errKey, err := InstallUserClawbotSkillZip(adminUserId, fileName, file, info[`skill_name`], info[`skill_description`], overwrite)
	if err != nil {
		logs.Error(err.Error())
		return nil, errors.New(i18n.Show(lang, `sys_err`))
	}
	if errKey != `` {
		return nil, errors.New(i18n.Show(lang, errKey))
	}
	return item, nil
}

func RunDocToSkillTask(id int64) (returnErr error) {
	if id <= 0 {
		return errors.New(`invalid doc-to-skill task id`)
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			panicErr := fmt.Errorf(`doc-to-skill task panic: %v`, recovered)
			logs.Error(`task:%d,err:%s`, id, panicErr.Error())
			_ = finishDocToSkillTask(id, define.DocToSkillTaskStatusFailed, msql.Datas{`err_msg`: panicErr.Error()})
			returnErr = panicErr
		}
	}()
	if _, loaded := runningDocToSkillTasks.LoadOrStore(id, struct{}{}); loaded {
		return nil
	}
	defer runningDocToSkillTasks.Delete(id)
	info, err := getDocToSkillTaskRow(0, id)
	if err != nil || len(info) == 0 {
		return err
	}
	stopKey := GetDocToSkillTaskStopKey(id)
	status := cast.ToInt(info[`status`])
	if status == define.DocToSkillTaskStatusStopping || IsDocToSkillTaskStopped(stopKey) {
		return setDocToSkillTaskStopped(id, info[`task_batch`], nil)
	}
	if status == define.DocToSkillTaskStatusRunning {
		return finishDocToSkillTask(id, define.DocToSkillTaskStatusFailed, msql.Datas{`err_msg`: `Task interrupted by service restart or failure; please retry manually.`})
	}
	if status != define.DocToSkillTaskStatusQueued {
		return nil
	}
	m := msql.Model(define.TableDocToSkillTask, define.Postgres)
	affected, err := m.Where(`id`, cast.ToString(id)).Where(`status`, cast.ToString(define.DocToSkillTaskStatusQueued)).Update(msql.Datas{
		`status`: define.DocToSkillTaskStatusRunning, `start_time`: tool.Time2Int(), `update_time`: tool.Time2Int(),
	})
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return err
	}
	if affected == 0 {
		latest, queryErr := getDocToSkillTaskRow(0, id)
		if queryErr != nil {
			return queryErr
		}
		if cast.ToInt(latest[`status`]) == define.DocToSkillTaskStatusStopping || IsDocToSkillTaskStopped(stopKey) {
			return setDocToSkillTaskStopped(id, info[`task_batch`], nil)
		}
		return nil
	}
	defer cleanupDocToSkillTaskWorkDir(info[`task_batch`])
	sourceFiles := decodeDocToSkillSourceFiles(info[`source_files`])
	prepared, prepareErr := prepareDocToSkillInput(info[`task_batch`], sourceFiles)
	if prepareErr != nil {
		return finishDocToSkillTask(id, define.DocToSkillTaskStatusFailed, msql.Datas{`err_msg`: prepareErr.Error()})
	}
	if IsDocToSkillTaskStopped(stopKey) {
		return setDocToSkillTaskStopped(id, info[`task_batch`], nil)
	}
	result, runErr := DoDocToSkill(define.LangZhCn, DocToSkillTaskInfo{
		TaskBatch:     info[`task_batch`],
		AdminUserId:   cast.ToInt(info[`admin_user_id`]),
		ModelConfigId: cast.ToInt(info[`model_config_id`]),
		UseModel:      info[`use_model`],
		Temperature:   cast.ToFloat32(info[`temperature`]),
		MaxToken:      cast.ToInt(info[`max_token`]),
		SourceFiles:   prepared,
		CustomPrompt:  info[`custom_prompt`],
		StopKey:       stopKey,
	})
	if IsDocToSkillTaskStopped(stopKey) {
		return setDocToSkillTaskStopped(id, info[`task_batch`], result.DebugLog)
	}
	if runErr != nil {
		return finishDocToSkillTask(id, define.DocToSkillTaskStatusFailed, msql.Datas{`debug_log`: tool.JsonEncodeNoError(result.DebugLog), `err_msg`: runErr.Error()})
	}
	if result.ZipPath == `` {
		return finishDocToSkillTask(id, define.DocToSkillTaskStatusFailed, msql.Datas{`debug_log`: tool.JsonEncodeNoError(result.DebugLog), `err_msg`: `generated zip path is empty`})
	}
	skillName, skillDescription, metaErr := ReadClawbotSkillZipMeta(result.ZipPath)
	if metaErr != nil {
		return finishDocToSkillTask(id, define.DocToSkillTaskStatusFailed, msql.Datas{`debug_log`: tool.JsonEncodeNoError(result.DebugLog), `err_msg`: metaErr.Error()})
	}
	uploadInfo, saveErr := saveGeneratedSkillZipFile(cast.ToInt(info[`admin_user_id`]), `doc_to_skill`, result.ZipPath, define.DocToSkillTaskZipAllowExt)
	if IsDocToSkillTaskStopped(stopKey) {
		return setDocToSkillTaskStopped(id, info[`task_batch`], result.DebugLog)
	}
	if saveErr != nil {
		return finishDocToSkillTask(id, define.DocToSkillTaskStatusFailed, msql.Datas{`debug_log`: tool.JsonEncodeNoError(result.DebugLog), `err_msg`: saveErr.Error()})
	}
	return finishDocToSkillTask(id, define.DocToSkillTaskStatusSucceed, msql.Datas{
		`skill_name`: skillName, `skill_description`: skillDescription, `file_name`: uploadInfo.Name,
		`file_url`: uploadInfo.Link, `file_size`: uploadInfo.Size, `debug_log`: tool.JsonEncodeNoError(result.DebugLog), `err_msg`: `SUCCEED`,
	})
}

func GetDocToSkillTaskStopKey(id int64) string {
	return define.DocToSkillTaskStopKeyPrefix + cast.ToString(id)
}

func GetDocToSkillTaskRunID(taskBatch string) string {
	return `doc-to-skill:` + strings.TrimSpace(taskBatch)
}

func IsDocToSkillTaskStopped(stopKey string) bool {
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

func getDocToSkillTaskRow(adminUserId int, id int64) (msql.Params, error) {
	m := msql.Model(define.TableDocToSkillTask, define.Postgres)
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

func buildDocToSkillTaskItem(row msql.Params, withLog bool) define.DocToSkillTaskItem {
	item := define.DocToSkillTaskItem{
		ID: cast.ToInt64(row[`id`]), TaskBatch: row[`task_batch`], SourceFiles: decodeDocToSkillSourceFiles(row[`source_files`]),
		CustomPrompt: row[`custom_prompt`], ModelConfigId: cast.ToInt(row[`model_config_id`]), UseModel: row[`use_model`],
		Temperature: cast.ToFloat32(row[`temperature`]), MaxToken: cast.ToInt(row[`max_token`]), Status: cast.ToInt(row[`status`]),
		SkillName: row[`skill_name`], SkillDescription: row[`skill_description`], FileName: row[`file_name`], FileUrl: row[`file_url`],
		FileSize: cast.ToInt(row[`file_size`]), ErrMsg: row[`err_msg`], StartTime: cast.ToInt(row[`start_time`]), EndTime: cast.ToInt(row[`end_time`]),
		CreateTime: cast.ToInt(row[`create_time`]), UpdateTime: cast.ToInt(row[`update_time`]),
	}
	if withLog && row[`debug_log`] != `` {
		_ = tool.JsonDecodeUseNumber(row[`debug_log`], &item.DebugLog)
	}
	return item
}

func decodeDocToSkillSourceFiles(raw string) []*define.UploadInfo {
	files := make([]*define.UploadInfo, 0)
	if raw != `` {
		_ = tool.JsonDecodeUseNumber(raw, &files)
	}
	return files
}

func checkDocToSkillCreateParams(lang string, adminUserId int, sourceFiles []*define.UploadInfo, modelConfigId int, useModel string, temperature float32, maxToken int) error {
	if !CheckModelIsValid(adminUserId, modelConfigId, strings.TrimSpace(useModel), Llm) {
		return errors.New(i18n.Show(lang, `param_invalid`, `use_model`))
	}
	if temperature < 0 || temperature > 2 {
		return errors.New(i18n.Show(lang, `param_invalid`, `temperature`))
	}
	if maxToken <= 0 {
		return errors.New(i18n.Show(lang, `param_invalid`, `max_token`))
	}
	if len(sourceFiles) == 0 || len(sourceFiles) > define.DocToSkillTaskMaxFileCount {
		return errors.New(i18n.Show(lang, `param_invalid`, `files`))
	}
	for _, file := range sourceFiles {
		if file == nil || strings.TrimSpace(file.Link) == `` || !tool.InArrayString(strings.ToLower(file.Ext), define.DocToSkillTaskAllowExt) {
			return errors.New(i18n.Show(lang, `param_invalid`, `files`))
		}
		name, ok := docToSkillOriginalFileName(file.Name)
		nameExt := strings.TrimPrefix(strings.ToLower(filepath.Ext(name)), `.`)
		fileExt := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(file.Ext)), `.`)
		if !ok || nameExt != fileExt {
			return errors.New(i18n.Show(lang, `param_invalid`, `files`))
		}
	}
	return nil
}

func docToSkillOriginalFileName(raw string) (string, bool) {
	name := filepath.Base(strings.TrimSpace(raw))
	return name, name != `` && name != `.` && name != `..` && !strings.HasPrefix(name, `.`)
}

func prepareDocToSkillInput(taskBatch string, sourceFiles []*define.UploadInfo) ([]string, error) {
	workDir, err := getDocToSkillTaskWorkDir(taskBatch)
	if err != nil {
		return nil, err
	}
	inputDir := filepath.Join(workDir, `input`)
	if err = os.RemoveAll(workDir); err != nil {
		return nil, err
	}
	if err := tool.MkDirAll(inputDir); err != nil {
		return nil, err
	}
	prepared := make([]string, 0, len(sourceFiles))
	for index, source := range sourceFiles {
		localPath := GetFileByLink(source.Link)
		if localPath == `` {
			return nil, fmt.Errorf(`source file not found: %s`, source.Name)
		}
		name, ok := docToSkillOriginalFileName(source.Name)
		if !ok {
			return nil, fmt.Errorf(`invalid source file name: %s`, source.Name)
		}
		target := filepath.Join(inputDir, fmt.Sprintf(`%03d-%s`, index+1, name))
		if err := copyDocToSkillFile(localPath, target); err != nil {
			return nil, fmt.Errorf(`copy source file %s: %w`, source.Name, err)
		}
		prepared = append(prepared, filepath.ToSlash(target))
	}
	return prepared, nil
}

func copyDocToSkillFile(source, target string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(writer, reader)
	closeErr := writer.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}

func enqueueDocToSkillTask(id int64) error {
	return AddJobs(define.DocToSkillTaskTopic, cast.ToString(id))
}

func markDocToSkillEnqueueFailed(id int64, enqueueErr error) {
	m := msql.Model(define.TableDocToSkillTask, define.Postgres)
	if _, err := m.Where(`id`, cast.ToString(id)).Update(msql.Datas{
		`status`: define.DocToSkillTaskStatusFailed, `err_msg`: enqueueErr.Error(), `update_time`: tool.Time2Int(), `end_time`: tool.Time2Int(),
	}); err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
	}
}

func setDocToSkillTaskStopped(id int64, taskBatch string, debugLog []any) error {
	defer cleanupDocToSkillTaskWorkDir(taskBatch)
	data := msql.Datas{`status`: define.DocToSkillTaskStatusStopped, `err_msg`: `STOPPED`, `end_time`: tool.Time2Int(), `update_time`: tool.Time2Int()}
	if debugLog != nil {
		data[`debug_log`] = tool.JsonEncodeNoError(debugLog)
	}
	m := msql.Model(define.TableDocToSkillTask, define.Postgres)
	_, err := m.Where(`id`, cast.ToString(id)).Where(`status`, `in`, fmt.Sprintf(`%d,%d,%d`, define.DocToSkillTaskStatusQueued, define.DocToSkillTaskStatusRunning, define.DocToSkillTaskStatusStopping)).Update(data)
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
	}
	return err
}

func finishDocToSkillTask(id int64, status int, data msql.Datas) error {
	data[`status`] = status
	data[`end_time`] = tool.Time2Int()
	data[`update_time`] = tool.Time2Int()
	m := msql.Model(define.TableDocToSkillTask, define.Postgres)
	affected, err := m.Where(`id`, cast.ToString(id)).Where(`status`, cast.ToString(define.DocToSkillTaskStatusRunning)).Update(data)
	if err != nil {
		logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
		return err
	}
	if affected == 0 {
		info, queryErr := getDocToSkillTaskRow(0, id)
		if queryErr != nil {
			return queryErr
		}
		if cast.ToInt(info[`status`]) == define.DocToSkillTaskStatusStopping {
			debugLog := make([]any, 0)
			if raw := cast.ToString(data[`debug_log`]); raw != `` {
				_ = tool.JsonDecodeUseNumber(raw, &debugLog)
			}
			return setDocToSkillTaskStopped(id, info[`task_batch`], debugLog)
		}
	}
	return nil
}

func normalizeDocToSkillZipPath(raw, taskBatch string) string {
	raw = strings.Trim(strings.TrimSpace(raw), "`'\" \n\r\t")
	for _, line := range strings.Split(raw, "\n") {
		line = strings.Trim(line, "`'\" \r\t")
		if strings.HasSuffix(strings.ToLower(line), `.zip`) {
			raw = line
			break
		}
	}
	raw = strings.TrimPrefix(strings.TrimPrefix(raw, `/workspace/`), `workspace/`)
	if !strings.HasSuffix(strings.ToLower(raw), `.zip`) {
		return ``
	}
	workDir, err := getDocToSkillTaskWorkDir(taskBatch)
	if err != nil {
		return ``
	}
	clean := filepath.Clean(filepath.FromSlash(raw))
	rel, err := filepath.Rel(workDir, clean)
	if err != nil || rel == `.` || strings.HasPrefix(rel, `..`) || !tool.IsFile(clean) {
		return ``
	}
	return clean
}

func cleanupDocToSkillTaskWorkDir(taskBatch string) {
	workDir, err := getDocToSkillTaskWorkDir(taskBatch)
	if err != nil {
		logs.Error(`invalid doc-to-skill task batch:%s,err:%s`, taskBatch, err.Error())
		return
	}
	if tool.IsDir(workDir) {
		if err = os.RemoveAll(workDir); err != nil {
			logs.Error(`remove doc-to-skill work dir:%s,err:%s`, workDir, err.Error())
		}
	}
}

func getDocToSkillTaskWorkDir(taskBatch string) (string, error) {
	taskBatch = strings.TrimSpace(taskBatch)
	if _, err := uuid.Parse(taskBatch); err != nil {
		return ``, fmt.Errorf(`invalid task batch: %w`, err)
	}
	baseDir := filepath.Clean(filepath.Dir(filepath.FromSlash(define.DocToSkillWorkDir)))
	workDir := filepath.Clean(filepath.Join(baseDir, taskBatch))
	rel, err := filepath.Rel(baseDir, workDir)
	if err != nil || rel == `.` || strings.HasPrefix(rel, `..`) {
		return ``, errors.New(`task work directory is outside doc-to-skill root`)
	}
	return workDir, nil
}
