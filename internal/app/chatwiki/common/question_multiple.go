// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	lib_define "chatwiki/internal/pkg/lib_define"
	"strings"

	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/v2/chat"
)

// ParseInputQuestion checks whether the input is a valid multimodal part array.
func ParseInputQuestion(question string) ([]chat.ContentPart, bool) {
	parts := make([]chat.ContentPart, 0)
	if err := tool.JsonDecodeUseNumber(question, &parts); err != nil || !validQuestionParts(parts) {
		return nil, false
	}
	return parts, true
}

func validQuestionParts(parts []chat.ContentPart) bool {
	if len(parts) == 0 {
		return false
	}
	for _, part := range parts {
		switch part.Type {
		case chat.ContentPartText:
			if part.Text == `` {
				return false
			}
		case chat.ContentPartImageURL:
			if part.ImageURL == nil || strings.TrimSpace(part.ImageURL.URL) == `` {
				return false
			}
		case chat.ContentPartInputAudio:
			if part.InputAudio == nil || strings.TrimSpace(part.InputAudio.Data) == `` || strings.TrimSpace(part.InputAudio.Format) == `` {
				return false
			}
		case chat.ContentPartVideoURL:
			if part.VideoURL == nil || strings.TrimSpace(part.VideoURL.URL) == `` {
				return false
			}
		default:
			return false
		}
	}
	return true
}

// AppendImageDomain appends the configured static-resource domain to relative links.
func AppendImageDomain(link string) string {
	if !IsUrl(link) && !strings.HasPrefix(link, `data:`) {
		link = define.Config.WebService[`image_domain`] + link
	}
	return link
}

func questionPartsAppendImageDomain(parts []chat.ContentPart) []chat.ContentPart {
	result := append([]chat.ContentPart(nil), parts...)
	for index, part := range result {
		switch part.Type {
		case chat.ContentPartImageURL:
			imageURL := *part.ImageURL
			imageURL.URL = AppendImageDomain(imageURL.URL)
			result[index].ImageURL = &imageURL
		case chat.ContentPartInputAudio:
			inputAudio := *part.InputAudio
			inputAudio.Data = AppendImageDomain(inputAudio.Data)
			result[index].InputAudio = &inputAudio
		case chat.ContentPartVideoURL:
			videoURL := *part.VideoURL
			videoURL.URL = AppendImageDomain(videoURL.URL)
			result[index].VideoURL = &videoURL
		}
	}
	return result
}

// ContentPartsAppendImageDomain normalizes persisted multimodal parts for business display and workflow variables.
func ContentPartsAppendImageDomain(parts []chat.ContentPart) []chat.ContentPart {
	return questionPartsAppendImageDomain(parts)
}

// ConvertQuestionMultiple converts encoded user questions into request content parts.
func ConvertQuestionMultiple(messages []chat.Message) []chat.Message {
	for index, message := range messages {
		if message.Role != chat.RoleUser || message.Content.Text == nil || message.Content.Parts != nil {
			continue
		}
		parts, ok := ParseInputQuestion(*message.Content.Text)
		if !ok {
			continue
		}
		messages[index].Content.Parts = questionPartsAppendImageDomain(parts)
		messages[index].Content.Text = nil
	}
	return messages
}

// convertUnsupportedQuestionParts replaces unsupported media with its JSON object text.
func convertUnsupportedQuestionParts(messages []chat.Message, support ModelInputSupport) []chat.Message {
	for index, message := range messages {
		if message.Role != chat.RoleUser || len(message.Content.Parts) == 0 {
			continue
		}
		parts := append([]chat.ContentPart(nil), message.Content.Parts...)
		for partIndex, part := range parts {
			var unsupported bool
			switch part.Type {
			case chat.ContentPartImageURL:
				unsupported = support.InputImage == 0
			case chat.ContentPartInputAudio:
				unsupported = support.InputVoice == 0
			case chat.ContentPartVideoURL:
				unsupported = support.InputVideo == 0
			}
			if unsupported {
				parts[partIndex] = chat.ContentPart{Type: chat.ContentPartText, Text: tool.JsonEncodeNoError(part)}
			}
		}
		messages[index].Content.Parts = parts
	}
	return messages
}

// GetQuestionByContentParts extracts the display text from multimodal input.
func GetQuestionByContentParts(parts []chat.ContentPart) string {
	for _, part := range parts {
		if part.Type == chat.ContentPartText && len(part.Text) > 0 {
			return part.Text
		}
	}
	for _, part := range parts {
		switch part.Type {
		case chat.ContentPartImageURL:
			return `[` + lib_define.MsgTypeNameMap[lib_define.MsgTypeImage] + `]`
		case chat.ContentPartInputAudio:
			return `[` + lib_define.MsgTypeNameMap[lib_define.MsgTypeVoice] + `]`
		case chat.ContentPartVideoURL:
			return `[` + lib_define.MsgTypeNameMap[lib_define.MsgTypeVideo] + `]`
		}
	}
	return ``
}

// GetFirstQuestionByInput extracts the first displayable question.
func GetFirstQuestionByInput(question string) string {
	if parts, ok := ParseInputQuestion(question); ok {
		return GetQuestionByContentParts(parts)
	}
	return question
}
