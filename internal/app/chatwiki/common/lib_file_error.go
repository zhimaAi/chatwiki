// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"strings"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/msql"
)

// AddLibFileFailureDetails also resolves legacy partial failures whose file error was not saved.
func AddLibFileFailureDetails(adminUserId, libraryId int, lang string, list []msql.Params) error {
	fileIds := make([]string, 0)
	for _, file := range list {
		if cast.ToInt(file[`status`]) == define.FileStatusPartException &&
			(strings.TrimSpace(file[`errmsg`]) == `` || file[`errmsg`] == `success`) {
			fileIds = append(fileIds, file[`id`])
		}
	}
	errors, err := GetLibFileErrors(adminUserId, libraryId, fileIds)
	if err != nil {
		return err
	}
	for _, file := range list {
		file[`failure_reason`], file[`failure_suggestion`] = ``, ``
		status := cast.ToInt(file[`status`])
		if status != define.FileStatusException && status != define.FileStatusPartException {
			continue
		}
		message := file[`errmsg`]
		if message == `success` || strings.TrimSpace(message) == `` {
			message = errors[file[`id`]]
		}
		reason, suggestion := LibFileErrorKeys(message)
		file[`failure_reason`] = i18n.Show(lang, reason)
		file[`failure_suggestion`] = i18n.Show(lang, suggestion)
		if status == define.FileStatusPartException {
			file[`errmsg`] = file[`failure_reason`]
		}
	}
	return nil
}

// GetLibFileErrors returns one stable error per file, including failures without an error message.
func GetLibFileErrors(adminUserId, libraryId int, fileIds []string) (map[string]string, error) {
	errors := make(map[string]string)
	if len(fileIds) == 0 {
		return errors, nil
	}
	vectors, err := msql.Model(`chat_ai_library_file_data_index`, define.Postgres).Alias(`i`).
		Join(`chat_ai_library_file_data d`, `i.data_id=d.id`, `inner`).
		Where(`i.admin_user_id`, cast.ToString(adminUserId)).Where(`i.library_id`, cast.ToString(libraryId)).
		Where(`i.file_id`, `in`, strings.Join(fileIds, `,`)).Where(`i.delete_time`, `0`).
		Where(`i.status`, cast.ToString(define.VectorStatusException)).
		Where(`d.admin_user_id`, cast.ToString(adminUserId)).Where(`d.library_id`, cast.ToString(libraryId)).
		WhereRaw(`d.file_id = i.file_id`).Where(`d.delete_time`, `0`).
		Group(`i.file_id`).Field(`i.file_id,COALESCE(MIN(NULLIF(i.errmsg,'')),'') AS errmsg`).Select()
	if err != nil {
		return nil, err
	}
	for _, row := range vectors {
		errors[row[`file_id`]] = row[`errmsg`]
	}
	splits, err := msql.Model(`chat_ai_library_file_data`, define.Postgres).
		Where(`admin_user_id`, cast.ToString(adminUserId)).Where(`library_id`, cast.ToString(libraryId)).
		Where(`file_id`, `in`, strings.Join(fileIds, `,`)).Where(`delete_time`, `0`).
		Where(`split_status`, cast.ToString(define.SplitStatusException)).
		Group(`file_id`).Field(`file_id,COALESCE(MIN(NULLIF(split_err_msg,'')),'') AS errmsg`).Select()
	if err != nil {
		return nil, err
	}
	for _, row := range splits {
		errors[row[`file_id`]] = row[`errmsg`]
	}
	return errors, nil
}

// libFileErrorRules maps a stored diagnostic to the i18n keys of its reason and suggestion.
// Our own errors keep their own reason; supplier text has no i18n and is grouped below them.
var libFileErrorRules = []struct {
	reason, suggestion string
	patterns           []string
}{
	{`doc_empty`, `document_invalid`, nil},
	{`excel_less_row`, `document_invalid`, nil},
	{`model_config_id_invalid`, `model_config_invalid`, nil},
	{`lib_file_failure_file_not_exist`, `document_invalid`, []string{`file not exist`}},
	{`lib_file_failure_model_not_config`, `model_config_invalid`, []string{`) not config`}},
	{`lib_file_failure_model_not_valid`, `model_config_invalid`, []string{`model not valid`}},
	{`lib_file_failure_usage_limit`, `usage_limit`, []string{`token usage exceeded`}},
	{`lib_file_failure_vector_empty`, `model_request_failed`, []string{`get vector return nil`}},

	{`lib_file_failure_insufficient_credit`, `insufficient_credit`, []string{`there are no integral available`, `ocr_integral_not_enough`, `insufficient_quota`, `insufficient balance`, `insufficient credit`}},
	{`lib_file_failure_model_auth_failed`, `model_auth_failed`, []string{`invalid_api_key`, `invalid_iam_token`, `invalid api key`, `incorrect api key`, `authentication failed`, `unauthorized`}},
	{`lib_file_failure_request_limited`, `request_limited`, []string{`rate_limit`, `rate limit`, `too many requests`}},
	{`lib_file_failure_request_timeout`, `request_timeout`, []string{`timeout`, `timed out`, `context deadline exceeded`, `connection refused`, `connection reset`, `no such host`}},
	{`lib_file_failure_content_too_long`, `content_too_long`, []string{`context_length_exceeded`, `maximum context length`, `input is too long`, `input too long`}},
	// Wraps a supplier detail, so it only applies when the detail matched nothing above.
	{`lib_file_failure_model_unavailable`, `model_config_invalid`, []string{`self-owned model unavailable`}},
}

// LibFileErrorKeys returns the i18n keys of the reason and suggestion for a stored diagnostic.
// i18n errors are stored as their key or translated text, so each reason also matches those.
func LibFileErrorKeys(message string) (string, string) {
	message = strings.ToLower(message)
	for _, rule := range libFileErrorRules {
		patterns := append([]string{rule.reason, i18n.Show(define.LangZhCn, rule.reason), i18n.Show(define.LangEnUs, rule.reason)}, rule.patterns...)
		for _, pattern := range patterns {
			if strings.Contains(message, strings.ToLower(pattern)) {
				return rule.reason, `lib_file_suggestion_` + rule.suggestion
			}
		}
	}
	return `lib_file_failure_unknown`, `lib_file_suggestion_unknown`
}
