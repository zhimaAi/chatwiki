// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"strings"
	"unicode"
)

const (
	promptOptimizeOutputStart = iota
	promptOptimizeOutputChat
	promptOptimizeOutputPrompt
	promptOptimizeOutputBetween
	promptOptimizeOutputReply
	promptOptimizeOutputEnd
)

type promptOptimizeOutputParser struct {
	phase   int
	pending string
	prompt  strings.Builder
	reply   strings.Builder
}

// Hold a possible delimiter across chunks so no part of a reserved tag is sent to the UI.
func (p *promptOptimizeOutputParser) push(chunk string) (promptDelta, replyDelta string, valid bool) {
	promptOffset, replyOffset := p.prompt.Len(), p.reply.Len()
	text := p.pending + chunk
	p.pending = ""
	for text != "" {
		before, _, hasTagStart := strings.Cut(text, "<")
		if !p.acceptText(before) {
			return "", "", false
		}
		if !hasTagStart {
			break
		}
		text = text[len(before):]
		width := 0
		for _, tag := range [...]string{
			define.PromptOptimizePromptOpenTag, define.PromptOptimizePromptCloseTag,
			define.PromptOptimizeReplyOpenTag, define.PromptOptimizeReplyCloseTag,
		} {
			name := tag[:len(tag)-1]
			if len(text) <= len(name) {
				if strings.EqualFold(name[:len(text)], text) {
					p.pending = text
					return p.prompt.String()[promptOffset:], p.reply.String()[replyOffset:], true
				}
				continue
			}
			if !strings.EqualFold(text[:len(name)], name) {
				continue
			}
			switch text[len(name)] {
			case '>':
				if !strings.HasPrefix(text, tag) || !p.acceptTag(tag) {
					return "", "", false
				}
				width = len(tag)
			case ' ', '\t', '\r', '\n', '/':
				return "", "", false
			}
			break
		}
		if width == 0 {
			if !p.acceptText("<") {
				return "", "", false
			}
			width = 1
		}
		text = text[width:]
	}
	return p.prompt.String()[promptOffset:], p.reply.String()[replyOffset:], true
}

func (p *promptOptimizeOutputParser) acceptText(text string) bool {
	if p.phase == promptOptimizeOutputStart {
		text = strings.TrimLeftFunc(text, unicode.IsSpace)
		if text == "" {
			return true
		}
		p.phase = promptOptimizeOutputChat
	}
	switch p.phase {
	case promptOptimizeOutputPrompt:
		p.prompt.WriteString(text)
	case promptOptimizeOutputChat, promptOptimizeOutputReply:
		p.reply.WriteString(text)
	default:
		return strings.TrimSpace(text) == ""
	}
	return true
}

func (p *promptOptimizeOutputParser) acceptTag(tag string) bool {
	switch {
	case p.phase == promptOptimizeOutputStart && tag == define.PromptOptimizePromptOpenTag:
		p.phase = promptOptimizeOutputPrompt
	case p.phase == promptOptimizeOutputPrompt && tag == define.PromptOptimizePromptCloseTag:
		p.phase = promptOptimizeOutputBetween
	case p.phase == promptOptimizeOutputBetween && tag == define.PromptOptimizeReplyOpenTag:
		p.phase = promptOptimizeOutputReply
	case p.phase == promptOptimizeOutputReply && tag == define.PromptOptimizeReplyCloseTag:
		p.phase = promptOptimizeOutputEnd
	default:
		return false
	}
	return true
}

func (p *promptOptimizeOutputParser) finish() (define.PromptOptimizeOutput, bool) {
	// A trailing less-than sign can be ordinary text, but a partial reserved name cannot.
	if p.pending != "" && ((p.pending != "<" && p.pending != "</") || !p.acceptText(p.pending)) {
		return define.PromptOptimizeOutput{}, false
	}
	p.pending = ""
	result := define.PromptOptimizeOutput{
		Prompt: strings.TrimSpace(p.prompt.String()),
		Reply:  strings.TrimSpace(p.reply.String()),
	}
	switch p.phase {
	case promptOptimizeOutputChat:
		return result, result.Reply != ""
	case promptOptimizeOutputEnd:
		return result, result.Prompt != "" && result.Reply != ""
	default:
		return define.PromptOptimizeOutput{}, false
	}
}
