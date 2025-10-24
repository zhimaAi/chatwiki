// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package common

import (
	"chatwiki/internal/app/message_service/define"
	"time"
)

func AddJobs(topic, message string, delay ...time.Duration) error {
	topic = define.Env + `_` + topic
	return define.ProducerHandle.AddJobs(topic, message, delay...)
}
