// Copyright © 2016- 2024 Sesame Network Technology all right reserved

package gemini

import (
	"bufio"
	"bytes"
	"chatwiki/internal/app/chatwiki/llm/common"
	"encoding/json"
	"errors"
	"io"
)

type ChatCompletionRequest struct {
	Contents         []Content        `json:"contents"`
	GenerationConfig GenerationConfig `json:"generationConfig,omitempty"`
}

type GenerationConfig struct {
	StopSequences   []string `json:"stopSequences,omitempty"`
	CandidateCount  int      `json:"candidateCount,omitempty"`
	MaxOutputTokens int      `json:"maxOutputTokens,omitempty"`
	Temperature     float64  `json:"temperature,omitempty"`
	TopP            int      `json:"topP,omitempty"`
	TopK            int      `json:"topK,omitempty"`
}

type ChatCompletionResponse struct {
	Candidates    []Candidate   `json:"candidates,omitempty"`
	UsageMetadata UsageMetadata `json:"usageMetadata,omitempty"`
}

type Candidate struct {
	Content               Content                `json:"content"`
	FinishReason          string                 `json:"finishReason"`
	SafetyRatings         []SafetyRating         `json:"safetyRatings"`
	CitationMetadata      CitationMetadata       `json:"citationMetadata"`
	TokenCount            int                    `json:"tokenCount"`
	GroundingAttributions []GroundingAttribution `json:"groundingAttributions"`
	Index                 int                    `json:"index"`
}
type CitationMetadata struct {
	CitationSources []CitationSource `json:"citationSources"`
}
type UsageMetadata struct {
	PromptTokenCount     int `json:"promptTokenCount"`
	CandidatesTokenCount int `json:"candidatesTokenCount"`
	TotalTokenCount      int `json:"totalTokenCount"`
}
type CitationSource struct {
	StartIndex int    `json:"startIndex"`
	EndIndex   int    `json:"endIndex"`
	Uri        string `json:"uri"`
	License    string `json:"license"`
}

type SafetyRating struct {
	Caregory    string `json:"caregory"`
	Probability string `json:"probability"`
	Blocked     bool   `json:"blocked"`
}

type GroundingAttribution struct {
	SourceId AttributionSourceId `json:"sourceId"`
	Content  Content             `json:"content"`
}

type AttributionSourceId struct {
	GroundingPassage       GroundingPassageId     `json:"groundingPassage"`
	SemanticRetrieverChunk SemanticRetrieverChunk `json:"semanticRetrieverChunk"`
}

type GroundingPassageId struct {
	PassageId string `json:"passageId"`
	PartIndex int    `json:"partIndex"`
}
type SemanticRetrieverChunk struct {
	Source string `json:"source"`
	Chunk  string `json:"chunk"`
}
type ChatCompletionStream struct {
	*common.StreamReader[ChatCompletionResponse]
}

func (c *ChatCompletionStream) Recv() (ChatCompletionResponse, error) {
	var response ChatCompletionResponse
	for {
		rawLine, readErr := readUntilDelimiter(c.StreamReader.Reader, []byte(",\r\n"))
		if readErr == io.EOF {
			c.StreamReader.IsFinished = true
			return *new(ChatCompletionResponse), io.EOF
		}
		noSpaceLine := bytes.TrimSpace(rawLine)
		if string(noSpaceLine) == "]" {
			c.StreamReader.IsFinished = true
			return *new(ChatCompletionResponse), io.EOF
		}
		if bytes.HasPrefix(noSpaceLine, []byte("[")) {
			noSpaceLine = bytes.TrimPrefix(noSpaceLine, []byte("["))
		}
		if bytes.HasSuffix(noSpaceLine, []byte("]")) {
			noSpaceLine = bytes.TrimSuffix(noSpaceLine, []byte("]"))
		}
		unmarshalErr := json.Unmarshal(noSpaceLine, &response)
		if unmarshalErr != nil {
			return *new(ChatCompletionResponse), unmarshalErr
		}
		if len(response.Candidates) <= 0 || len(response.Candidates[0].Content.Parts) <= 0 {
			return ChatCompletionResponse{}, errors.New("gemini response candidates is empty")
		}
		return response, nil
	}
}

func readUntilDelimiter(r *bufio.Reader, delim []byte) ([]byte, error) {
	var buf []byte
	for {
		b, err := r.ReadByte()
		if err != nil {
			if err == io.EOF && len(buf) > 0 {
				return buf, nil
			}
			return buf, err
		}
		buf = append(buf, b)
		if len(buf) >= len(delim) && bytes.Equal(buf[len(buf)-len(delim):], delim) {
			return buf[:len(buf)-len(delim)], nil
		}
	}
}
