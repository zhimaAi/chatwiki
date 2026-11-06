// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/i18n"
	"errors"
	"unicode/utf8"

	"github.com/spf13/cast"
	"github.com/zhimaAi/llm_adaptor/v2/chat"
)

const ReasoningEffortMaxLength = 50

// ToReasoningEffort keeps legacy switches authoritative and defaults enabled legacy data to medium.
func ToReasoningEffort(enableThinking any, reasoningEffort string) chat.ReasoningEffort {
	if !cast.ToBool(cast.ToUint(enableThinking)) {
		return chat.ReasoningEffortNone
	}
	if reasoningEffort == "" {
		return chat.ReasoningEffortMedium
	}
	return chat.ReasoningEffort(reasoningEffort)
}

func ValidateReasoningEffort(lang, reasoningEffort, field string) error {
	if utf8.RuneCountInString(reasoningEffort) > ReasoningEffortMaxLength {
		return errors.New(i18n.Show(lang, "param_invalid", field))
	}
	return nil
}
