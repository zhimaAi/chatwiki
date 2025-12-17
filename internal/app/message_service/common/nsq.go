// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/message_service/define"
	"time"
)

func AddJobs(topic, message string, delay ...time.Duration) error {
	topic = define.Env + `_` + topic
	return define.ProducerHandle.AddJobs(topic, message, delay...)
}
