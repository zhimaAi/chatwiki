// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package define

const (
	MessageTypeAll     = 0
	MessageTypeSpecify = 1
)
const (
	RuleTypeMessageType = `receive_reply_message_type` // 消息类型规则
	RuleTypeDuration    = `receive_reply_duration`     //时间范围规则

	RuleTypeSubscribeDefault  = `subscribe_reply_default`
	RuleTypeSubscribeDuration = `subscribe_reply_duration`
	RuleTypeSubscribeSource   = `subscribe_reply_source`
)

const (
	DurationTypeWeek      = `week`       //时间类型：week:周
	DurationTypeDay       = `day`        //时间类型：day:天
	DurationTypeTimeRange = `time_range` //时间范围规则

)
