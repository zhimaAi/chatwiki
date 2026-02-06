// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package define

const (
	MessageTypeAll     = 0
	MessageTypeSpecify = 1
)
const (
	RuleTypeMessageType = `receive_reply_message_type` // message type rule
	RuleTypeDuration    = `receive_reply_duration`     // time range rule

	RuleTypeSubscribeDefault  = `subscribe_reply_default`
	RuleTypeSubscribeDuration = `subscribe_reply_duration`
	RuleTypeSubscribeSource   = `subscribe_reply_source`
)

const (
	DurationTypeWeek      = `week`       // duration type: week
	DurationTypeDay       = `day`        // duration type: day
	DurationTypeTimeRange = `time_range` // time range rule

)
