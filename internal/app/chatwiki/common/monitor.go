// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

const ChatMonitorBaseTableName = `chat_ai_chat_monitor`

func SaveChatMonitor(data msql.Datas) {
	tableName := fmt.Sprintf(`%s_%s`, ChatMonitorBaseTableName, tool.Date(`Ym`))
	_, err := msql.Model(tableName, define.Postgres).Insert(data)
	if err == nil {
		return
	}
	var sqlerr *pq.Error
	if errors.As(err, &sqlerr) && sqlerr.Code == `42P01` { //表不存在
		//创建新表
		sql := fmt.Sprintf(`CREATE TABLE "%s" (LIKE "%s" INCLUDING ALL)`, tableName, ChatMonitorBaseTableName)
		_, err = msql.RawExec(define.Postgres, sql, nil)
		if err != nil {
			logs.Error(err.Error()) //创建新表出错了
		}
		//尝试重新插入
		if _, err = msql.Model(tableName, define.Postgres).Insert(data); err == nil {
			return
		}
	}
	logs.Error(err.Error())
}

type LibUseTime struct {
	QuestionOp int64 `json:"question_op"` //millisecond
	RecallTime int64 `json:"recall_time"` //millisecond
	RerankTime int64 `json:"rerank_time"` //millisecond
}

type NodeLog struct {
	StartTime int64        `json:"start_time"` //millisecond
	EndTime   int64        `json:"end_time"`   //millisecond
	NodeKey   string       `json:"node_key"`
	NodeName  string       `json:"node_name"`
	NodeType  int          `json:"node_type"`
	Output    SimpleFields `json:"output"`
	Error     error        `json:"error"`
	UseTime   int64        `json:"use_time"` //millisecond
}

type Monitor struct {
	params      *define.ChatRequestParam
	StartTime   int64      `json:"start_time"` //millisecond
	EndTime     int64      `json:"end_time"`   //millisecond
	Error       error      `json:"error"`
	LibUseTime  LibUseTime `json:"lib_use_time"`
	RequestTime int64      `json:"request_time"`  //millisecond
	LlmCallTime int64      `json:"llm_call_time"` //millisecond
	DebugLog    []any      `json:"debug_log"`
	NodeLogs    []NodeLog  `json:"node_logs"`
	AllUseTime  int64      `json:"all_use_time"` //millisecond
}

func (m *Monitor) Save(err error) {
	if len(m.params.Robot) == 0 {
		return //机器人信息为空,请求参数robot_key错误
	}
	m.EndTime = time.Now().UnixMilli()
	m.AllUseTime = m.EndTime - m.StartTime
	if err != nil {
		m.Error = err
	}
	//日志准入条件限制
	switch cast.ToInt(m.params.Robot[`application_type`]) {
	case define.ApplicationTypeChat:
		if m.Error == nil && m.AllUseTime < 2000 {
			return
		}
	case define.ApplicationTypeFlow:
		if m.Error == nil && m.AllUseTime < 5000 {
			return
		}
		m.LlmCallTime = m.RequestTime //特殊数据修正
	default:
		return
	}
	//开始记录
	if define.IsDev {
		logs.Other(`monitor`, tool.JsonEncodeNoError(m))
	}
	data := msql.Datas{
		`admin_user_id`:    m.params.AdminUserId,
		`robot_id`:         m.params.Robot[`id`],
		`robot_name`:       m.params.Robot[`robot_name`],
		`application_type`: m.params.Robot[`application_type`],
		`openid`:           m.params.Openid,
		`content`:          m.params.Question,
		`start_time`:       m.StartTime,
		`end_time`:         m.EndTime,
		`all_use_time`:     m.AllUseTime,
		`is_error`:         cast.ToInt(m.Error != nil),
		`error_msg`:        fmt.Sprintf(`%v`, m.Error),
		`question_op`:      m.LibUseTime.QuestionOp,
		`recall_time`:      m.LibUseTime.RecallTime,
		`rerank_time`:      m.LibUseTime.RerankTime,
		`request_time`:     m.RequestTime,
		`llm_call_time`:    m.LlmCallTime,
		`debug_log`:        tool.JsonEncodeNoError(m.DebugLog),
		`node_logs`:        tool.JsonEncodeNoError(m.NodeLogs),
		`create_time`:      tool.Time2Int(),
		`update_time`:      tool.Time2Int(),
	}
	SaveChatMonitor(data)
}

func NewMonitor(params *define.ChatRequestParam) *Monitor {
	return &Monitor{
		params:    params,
		StartTime: time.Now().UnixMilli(),
		DebugLog:  make([]any, 0),
		NodeLogs:  make([]NodeLog, 0),
	}
}
