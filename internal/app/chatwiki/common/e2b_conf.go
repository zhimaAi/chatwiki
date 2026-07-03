// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/custom_eino"
	"chatwiki/internal/pkg/lib_redis"
	"context"
	"errors"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/matiasinsaurralde/go-e2b"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
)

type E2bConfCacheBuildHandler struct{ RobotKey string }

func (h *E2bConfCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf(`chatwiki.robot_e2b_conf.%s`, h.RobotKey)
}
func (h *E2bConfCacheBuildHandler) GetCacheData() (any, error) {
	return msql.Model(define.TableChatAiRobotE2bConf, define.Postgres).Where(`robot_key`, h.RobotKey).Find()
}
func GetE2bConfInfo(robotKey string) (msql.Params, error) {
	result := make(msql.Params)
	err := lib_redis.GetCacheWithBuild(define.Redis, &E2bConfCacheBuildHandler{RobotKey: robotKey}, &result, time.Hour)
	return result, err
}

func GetE2bConf(lang string, adminUserId int, robotKey string) (*define.E2bConfParams, error) {
	if err := checkE2bConfRobot(lang, adminUserId, robotKey); err != nil {
		return nil, err
	}
	info, err := GetE2bConfInfo(robotKey)
	if err != nil {
		logs.Error(err.Error())
		return nil, errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(info) == 0 {
		return &define.E2bConfParams{RobotKey: robotKey}, nil
	}
	info[`api_key`] = maskE2bApiKey(info[`api_key`])
	return BuildE2BConfParams(info), nil
}

func BuildE2BConfParams(info msql.Params) *define.E2bConfParams {
	return &define.E2bConfParams{
		RobotKey:       info[`robot_key`],
		SwitchStatus:   cast.ToInt(info[`switch_status`]),
		ApiKey:         info[`api_key`],
		ApiBaseUrl:     info[`api_base_url`],
		SandboxDomain:  info[`sandbox_domain`],
		Template:       info[`template`],
		Timeout:        cast.ToInt(info[`timeout`]),
		CommandTimeout: cast.ToInt(info[`command_timeout`]),
		CommandUser:    info[`command_user`],
	}
}

func maskE2bApiKey(apiKey string) string {
	runes := []rune(apiKey)
	if len(runes) == 0 {
		return ``
	}
	if len(runes) <= 8 {
		return `****`
	}
	return string(runes[:4]) + `****` + string(runes[len(runes)-4:])
}

func SaveE2bConf(lang string, adminUserId int, params define.E2bConfParams) (*define.E2bConfParams, error) {
	if err := checkE2bConfRobot(lang, adminUserId, params.RobotKey); err != nil {
		return nil, err
	}
	existing, err := GetE2bConfInfo(params.RobotKey)
	if err != nil {
		logs.Error(err.Error())
		return nil, errors.New(i18n.Show(lang, `sys_err`))
	}
	data := msql.Datas{
		`switch_status`: params.SwitchStatus,
		`update_time`:   time.Now().Unix(),
	}
	if cast.ToBool(params.SwitchStatus) {
		if len(existing) > 0 && params.ApiKey == maskE2bApiKey(existing[`api_key`]) {
			params.ApiKey = existing[`api_key`]
		}
		if err = checkE2bConfParams(lang, params); err != nil {
			return nil, err
		}
		if err = CheckE2bConfConnection(params); err != nil {
			return nil, err
		}
		data[`api_key`] = params.ApiKey
		data[`api_base_url`] = params.ApiBaseUrl
		data[`sandbox_domain`] = params.SandboxDomain
		data[`template`] = params.Template
		data[`timeout`] = params.Timeout
		data[`command_timeout`] = params.CommandTimeout
		data[`command_user`] = params.CommandUser
	}
	m := msql.Model(define.TableChatAiRobotE2bConf, define.Postgres)
	if len(existing) == 0 {
		data[`admin_user_id`] = adminUserId
		data[`robot_key`] = params.RobotKey
		data[`create_time`] = time.Now().Unix()
		if _, err = m.Insert(data); err != nil {
			logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
			return nil, errors.New(i18n.Show(lang, `sys_err`))
		}
	} else {
		if _, err = m.Where(`robot_key`, params.RobotKey).Update(data); err != nil {
			logs.Error(`sql:%s,err:%s`, m.GetLastSql(), err.Error())
			return nil, errors.New(i18n.Show(lang, `sys_err`))
		}
	}
	lib_redis.DelCacheData(define.Redis, &E2bConfCacheBuildHandler{RobotKey: params.RobotKey})
	return GetE2bConf(lang, adminUserId, params.RobotKey)
}

func checkE2bConfRobot(lang string, adminUserId int, robotKey string) error {
	if !CheckRobotKey(robotKey) {
		return errors.New(i18n.Show(lang, `param_invalid`, `robot_key`))
	}
	robot, err := GetRobotInfo(robotKey)
	if err != nil {
		logs.Error(err.Error())
		return errors.New(i18n.Show(lang, `sys_err`))
	}
	if len(robot) == 0 || cast.ToInt(robot[`admin_user_id`]) != adminUserId {
		return errors.New(i18n.Show(lang, `no_data`))
	}
	return nil
}

func checkE2bConfParams(lang string, params define.E2bConfParams) error {
	if len(params.ApiKey) == 0 || utf8.RuneCountInString(params.ApiKey) > define.E2bConfApiKeyMaxLen {
		return errors.New(i18n.Show(lang, `param_invalid`, `api_key`))
	}
	if len(params.ApiBaseUrl) == 0 || utf8.RuneCountInString(params.ApiBaseUrl) > define.E2bConfApiBaseUrlMaxLen {
		return errors.New(i18n.Show(lang, `param_invalid`, `api_base_url`))
	}
	if len(params.SandboxDomain) == 0 || utf8.RuneCountInString(params.SandboxDomain) > define.E2bConfSandboxDomainMaxLen {
		return errors.New(i18n.Show(lang, `param_invalid`, `sandbox_domain`))
	}
	if len(params.Template) == 0 || utf8.RuneCountInString(params.Template) > define.E2bConfTemplateMaxLen {
		return errors.New(i18n.Show(lang, `param_invalid`, `template`))
	}
	if params.Timeout <= 0 {
		return errors.New(i18n.Show(lang, `param_invalid`, `timeout`))
	}
	if params.CommandTimeout <= 0 {
		return errors.New(i18n.Show(lang, `param_invalid`, `command_timeout`))
	}
	if len(params.CommandUser) == 0 || utf8.RuneCountInString(params.CommandUser) > define.E2bConfCommandUserMaxLen {
		return errors.New(i18n.Show(lang, `param_invalid`, `command_user`))
	}
	return nil
}

func CheckE2bConfConnection(conf define.E2bConfParams) error {
	timeout := min(30, max(10, conf.Timeout, conf.CommandTimeout))
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	e2bShell, err := custom_eino.NewE2BShell(ctx,
		e2b.ClientConfig{APIKey: conf.ApiKey, APIBaseURL: conf.ApiBaseUrl, SandboxDomain: conf.SandboxDomain},
		e2b.SandboxConfig{Template: conf.Template, Timeout: conf.Timeout, Secure: true},
		e2b.WithTimeout(time.Duration(conf.CommandTimeout)*time.Second),
		e2b.WithUser(conf.CommandUser),
	)
	if err != nil {
		return err
	}
	defer func() {
		if err = e2bShell.Close(); err != nil {
			logs.Error(`close E2B sandbox failed: %v`, err)
		}
	}()
	// Run bash's no-op builtin to verify that the sandbox can execute commands.
	resp, err := e2bShell.Execute(`:`)
	if err != nil {
		return err
	}
	if resp.ExitCode != nil && *resp.ExitCode != 0 {
		return fmt.Errorf(`E2B sandbox execute command failed: %s`, resp.Output)
	}
	return nil
}
