// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"strings"
	"unicode/utf8"
)

// AiSpliter ai割客户端
type AiSpliter struct {
	AiChunkSize int // 最大分块大小（按字符数）
	Overlap     int // 分块重叠大小（按字符数）
}

// NewAiSpliterClient 创建一个新的语义分割客户端
func NewAiSpliterClient(aiChunkSize int) *AiSpliter {
	return &AiSpliter{
		AiChunkSize: aiChunkSize,
		Overlap:     50,
	}
}

// SplitText 分割文本
func (c *AiSpliter) SplitText(text string) ([]string, error) {
	var result []string
	contents := ""
	if utf8.RuneCountInString(text) <= c.AiChunkSize {
		result = append(result, text)
	} else {
		chunks := c.aiSplit(text)
		for k, chunk := range chunks {
			contents += chunk
			if utf8.RuneCountInString(contents) >= c.AiChunkSize {
				result = append(result, contents)
				contents = ""
			} else if k+1 == len(chunks) {
				result = append(result, contents)
			}
		}
	}
	return result, nil
}

func (c *AiSpliter) aiSplit(text string) []string {
	separators := []string{"。", "？", "?", "\n"}
	var sentences []string
	if len(text) == 0 {
		return sentences
	}
	startIdx := 0
	for i := 0; i < len(text); i++ {
		for _, sep := range separators {
			if i+len(sep) <= len(text) && text[i:i+len(sep)] == sep {
				if i > startIdx {
					sentence := strings.TrimSpace(text[startIdx:i+len(sep)]) + sep
					if len(sentence) > 0 {
						sentences = append(sentences, sentence)
					}
				}
				startIdx = i + len(sep)
				break
			}
		}
	}
	if startIdx < len(text) {
		lastSentence := strings.TrimSpace(text[startIdx:])
		if len(lastSentence) > 0 {
			sentences = append(sentences, lastSentence)
		}
	}

	return sentences
}

// addOverlappingContent 为分段添加重叠内容
func (c *AiSpliter) addOverlappingContent(chunks []string) []string {
	if len(chunks) <= 1 || c.Overlap <= 0 {
		return chunks
	}

	result := make([]string, len(chunks))
	copy(result, chunks)

	for i := 1; i < len(result); i++ {
		prevChunk := chunks[i-1]

		// 将前一个块分割成句子
		prevSentences := c.splitIntoSentences(prevChunk)
		if len(prevSentences) == 0 {
			continue
		}

		// 添加句子作为重叠内容，直到接近但不超过设定的重叠大小
		var overlappingSentences []string
		var overlapSize int

		// 从后向前添加句子
		for j := len(prevSentences) - 1; j >= 0; j-- {
			sentenceSize := utf8.RuneCountInString(prevSentences[j])

			// 如果添加当前句子后总大小不超过设定的重叠大小，或者尚未添加任何句子，则添加它
			if overlapSize+sentenceSize <= c.Overlap || len(overlappingSentences) == 0 {
				overlappingSentences = append([]string{prevSentences[j]}, overlappingSentences...)
				overlapSize += sentenceSize
			} else {
				break
			}
		}

		// 将重叠句子添加到当前块的前面
		if len(overlappingSentences) > 0 {
			overlappingText := strings.Join(overlappingSentences, "")
			result[i] = overlappingText + result[i]
		}
	}

	return result
}

// splitIntoSentences 将文本分割成句子
func (c *AiSpliter) splitIntoSentences(text string) []string {
	separators := []string{"。", "？", "?", "\n"}
	var sentences []string
	if len(text) == 0 {
		return sentences
	}
	startIdx := 0
	for i := 0; i < len(text); i++ {
		for _, sep := range separators {
			if i+len(sep) <= len(text) && text[i:i+len(sep)] == sep {
				if i > startIdx {
					sentence := strings.TrimSpace(text[startIdx : i+len(sep)])
					if len(sentence) > 0 {
						sentences = append(sentences, sentence)
					}
				}
				startIdx = i + len(sep)
				break
			}
		}
	}
	if startIdx < len(text) {
		lastSentence := strings.TrimSpace(text[startIdx:])
		if len(lastSentence) > 0 {
			sentences = append(sentences, lastSentence)
		}
	}

	return sentences
}
