// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"strings"
	"unicode/utf8"
)

// AiSpliter AI splitter client
type AiSpliter struct {
	AiChunkSize int // Maximum chunk size (by character count)
	Overlap     int // Chunk overlap size (by character count)
}

// NewAiSpliterClient creates a new semantic segmentation client
func NewAiSpliterClient(aiChunkSize int) *AiSpliter {
	return &AiSpliter{
		AiChunkSize: aiChunkSize,
		Overlap:     50,
	}
}

// SplitText splits text
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

// addOverlappingContent adds overlapping content to chunks
func (c *AiSpliter) addOverlappingContent(chunks []string) []string {
	if len(chunks) <= 1 || c.Overlap <= 0 {
		return chunks
	}

	result := make([]string, len(chunks))
	copy(result, chunks)

	for i := 1; i < len(result); i++ {
		prevChunk := chunks[i-1]

		// Split previous chunk into sentences
		prevSentences := c.splitIntoSentences(prevChunk)
		if len(prevSentences) == 0 {
			continue
		}

		// Add sentences as overlapping content until close to but not exceeding the set overlap size
		var overlappingSentences []string
		var overlapSize int

		// Add sentences from back to front
		for j := len(prevSentences) - 1; j >= 0; j-- {
			sentenceSize := utf8.RuneCountInString(prevSentences[j])

			// If adding current sentence doesn't exceed the set overlap size, or no sentences have been added yet, add it
			if overlapSize+sentenceSize <= c.Overlap || len(overlappingSentences) == 0 {
				overlappingSentences = append([]string{prevSentences[j]}, overlappingSentences...)
				overlapSize += sentenceSize
			} else {
				break
			}
		}

		// Add overlapping sentences to the beginning of current chunk
		if len(overlappingSentences) > 0 {
			overlappingText := strings.Join(overlappingSentences, "")
			result[i] = overlappingText + result[i]
		}
	}

	return result
}

// splitIntoSentences splits text into sentences
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
