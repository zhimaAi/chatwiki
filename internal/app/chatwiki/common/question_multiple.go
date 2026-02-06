// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"

	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

// ParseInputQuestion Parse input question to check if it's multimodal JSON format
func ParseInputQuestion(question string) (adaptor.QuestionMultiple, bool) {
	questionMultiple := make(adaptor.QuestionMultiple, 0)
	err := tool.JsonDecodeUseNumber(question, &questionMultiple)
	if err == nil && len(questionMultiple) > 0 {
		return questionMultiple, true
	}
	return nil, false
}

// AppendImageDomain Unified processing of appending static resource domain
func AppendImageDomain(link string) string {
	if !IsUrl(link) {
		link = define.Config.WebService[`image_domain`] + link
	}
	return link
}

// QuestionMultipleAppendImageDomain Append static resource domain for multimodal input
func QuestionMultipleAppendImageDomain(questionMultiple adaptor.QuestionMultiple) adaptor.QuestionMultiple {
	for i, item := range questionMultiple {
		switch item.Type {
		case adaptor.TypeImage:
			questionMultiple[i].ImageUrl.Url = AppendImageDomain(item.ImageUrl.Url)
		case adaptor.TypeAudio:
			questionMultiple[i].InputAudio.Data = AppendImageDomain(item.InputAudio.Data)
		case adaptor.TypeVideo:
			questionMultiple[i].VedioUrl.Url = AppendImageDomain(item.VedioUrl.Url)
		}
	}
	return questionMultiple
}

// ConvertQuestionMultiple Convert to multimodal input structure
func ConvertQuestionMultiple(messages []adaptor.ZhimaChatCompletionMessage) []adaptor.ZhimaChatCompletionMessage {
	for i, message := range messages {
		if message.Role != `user` {
			continue
		}
		if questionMultiple, ok := ParseInputQuestion(message.Content); ok {
			messages[i].SetQuestionMultiple(QuestionMultipleAppendImageDomain(questionMultiple))
		}
	}
	return messages
}

// GetQuestionByQuestionMultiple Extract input question from multimodal input structure
func GetQuestionByQuestionMultiple(questionMultiple adaptor.QuestionMultiple) string {
	for _, item := range questionMultiple {
		if item.Type == adaptor.TypeText && len(item.Text) > 0 {
			return item.Text
		}
	}
	return ``
}

// GetFirstQuestionByInput Extract first question from user input
func GetFirstQuestionByInput(question string) string {
	if questionMultiple, ok := ParseInputQuestion(question); ok {
		return GetQuestionByQuestionMultiple(questionMultiple)
	}
	return question
}
