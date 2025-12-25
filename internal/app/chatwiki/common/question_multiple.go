// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"strings"

	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

// ParseInputQuestion 解析输入的问题是否为多模态json格式
func ParseInputQuestion(question string) (adaptor.QuestionMultiple, bool) {
	questionMultiple := make(adaptor.QuestionMultiple, 0)
	err := tool.JsonDecodeUseNumber(question, &questionMultiple)
	if err == nil && len(questionMultiple) > 0 {
		return questionMultiple, true
	}
	return nil, false
}

// AppendImageDomain 统一处理追加静态资源的域名
func AppendImageDomain(link string) string {
	if !IsUrl(link) {
		link = define.Config.WebService[`image_domain`] + link
	}
	return link
}

// QuestionMultipleAppendImageDomain 多模态输入追加静态资源的域名
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

// ConvertQuestionMultiple 转换成多模态输入结构
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

// GetQuestionByQuestionMultiple 从多模态输入结构抽取输入的问题
func GetQuestionByQuestionMultiple(questionMultiple adaptor.QuestionMultiple) string {
	for _, item := range questionMultiple {
		if item.Type == adaptor.TypeText && len(item.Text) > 0 {
			return item.Text
		}
	}
	return ``
}

// GetFirstQuestionByInput 从用户输入抽取第一个问题
func GetFirstQuestionByInput(question string) string {
	if questionMultiple, ok := ParseInputQuestion(question); ok {
		return GetQuestionByQuestionMultiple(questionMultiple)
	}
	return question
}

// GetQuestionByContent 支持用户content为string或者array<object>类型
func GetQuestionByContent(content any) string {
	if question, ok := content.(string); ok {
		return strings.TrimSpace(question)
	}
	return tool.JsonEncodeNoError(content)
}
