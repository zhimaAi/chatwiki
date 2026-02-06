// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/pkg/lib_redis"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
)

// SemanticSplitter Semantic segmentation client
type SemanticSplitter struct {
	SemanticChunkSize      int                       // Maximum chunk size (by character count)
	SemanticChunkOverlap   int                       // Chunk overlap size (by character count)
	SemanticChunkThreshold int                       // Semantic breakpoint threshold (percentage, 1-100)
	AdminUserId            int                       // User ID
	ModelConfigId          int                       // Model configuration ID
	UseModel               string                    // Model being used
	GoRoutineNum           int                       // Goroutine number
	GetVector              func() ([]float64, error) // Function to get vector
}

// NewSemanticSplitterClient Creates a new semantic splitter client
func NewSemanticSplitterClient() *SemanticSplitter {
	return &SemanticSplitter{}
}

// SplitText Split text
func (c *SemanticSplitter) SplitText(text string) ([]string, error) {
	chunks, err := c.semanticSplit(text)
	if err != nil {
		return []string{text}, err
	}

	var result []string
	for _, chunk := range chunks {
		if utf8.RuneCountInString(chunk) > c.SemanticChunkSize {
			sentences := c.splitIntoSentences(chunk)
			subChunks := c.splitByBinaryDivision(sentences)
			result = append(result, subChunks...)
		} else {
			result = append(result, chunk)
		}
	}
	// Add overlapping content
	result = c.addOverlappingContent(result)
	return result, nil
}

// addOverlappingContent Add overlapping content to segments
func (c *SemanticSplitter) addOverlappingContent(chunks []string) []string {
	if len(chunks) <= 1 || c.SemanticChunkOverlap <= 0 {
		return chunks
	}

	result := make([]string, len(chunks))
	copy(result, chunks)

	for i := 1; i < len(result); i++ {
		prevChunk := chunks[i-1]

		// Split the previous chunk into sentences
		prevSentences := c.splitIntoSentences(prevChunk)
		if len(prevSentences) == 0 {
			continue
		}

		// Add sentences as overlapping content, approaching but not exceeding the set overlap size
		var overlappingSentences []string
		var overlapSize int

		// Add sentences from back to front
		for j := len(prevSentences) - 1; j >= 0; j-- {
			sentenceSize := utf8.RuneCountInString(prevSentences[j])

			// If adding the current sentence does not exceed the set overlap size, or no sentences have been added yet, add it
			if overlapSize+sentenceSize <= c.SemanticChunkOverlap || len(overlappingSentences) == 0 {
				overlappingSentences = append([]string{prevSentences[j]}, overlappingSentences...)
				overlapSize += sentenceSize
			} else {
				break
			}
		}

		// Add overlapping sentences to the front of the current chunk
		if len(overlappingSentences) > 0 {
			overlappingText := strings.Join(overlappingSentences, "")
			result[i] = overlappingText + result[i]
		}
	}

	return result
}

// semanticSplit Semantic segmentation
func (c *SemanticSplitter) semanticSplit(text string) ([]string, error) {
	sentences := c.splitIntoSentences(text)
	if len(sentences) <= 1 {
		return []string{text}, nil
	}
	vectors, err := c.generateEmbeddings(sentences)
	if err != nil {
		logs.Error("Generate embeddings error: %v", err)
		return []string{text}, err
	}
	differences := c.calculateDifferences(vectors)
	breakpoints := c.findBreakpoints(differences)
	chunks := c.formChunks(sentences, breakpoints)
	return chunks, nil
}

// splitIntoSentences Split text into sentences
func (c *SemanticSplitter) splitIntoSentences(text string) []string {
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

// VectorCacheBuildHandler Vector cache handler
type VectorCacheBuildHandler struct {
	AdminUserId   int
	ModelConfigId int
	UseModel      string
	Sentence      string
}

func (h *VectorCacheBuildHandler) GetCacheKey() string {
	return fmt.Sprintf("chatwiki.vector_cache.%d.%d.%s.%x", h.AdminUserId, h.ModelConfigId, h.UseModel, tool.MD5(h.Sentence))
}

func (h *VectorCacheBuildHandler) GetCacheData() (any, error) {
	embedding, err := GetVector2000(
		define.LangEnUs,
		h.AdminUserId,
		cast.ToString(h.AdminUserId),
		msql.Params{},
		msql.Params{},
		msql.Params{},
		cast.ToInt(h.ModelConfigId),
		h.UseModel,
		h.Sentence,
	)
	if err != nil {
		return nil, err
	}

	// Parse JSON string to float64 array
	var vector []float64
	if err := json.Unmarshal([]byte(embedding), &vector); err != nil {
		logs.Error("Get vector error: %v", err)
		return nil, err
	}
	return vector, nil
}

// generateEmbeddings Generate vector representations for each sentence
func (c *SemanticSplitter) generateEmbeddings(sentences []string) ([][]float64, error) {
	var vectors = make([][]float64, len(sentences))

	// Use channel to control concurrency
	type vectorResult struct {
		index  int
		vector []float64
		err    error
	}
	resultChan := make(chan vectorResult, len(sentences))

	// Create semaphore to limit concurrent number
	sem := make(chan struct{}, c.GoRoutineNum)

	// Start goroutines to generate vectors
	for i, sentence := range sentences {
		go func(idx int, text string) {
			sem <- struct{}{}        // Acquire semaphore
			defer func() { <-sem }() // Release semaphore

			// Try to get vector from cache
			var vector []float64
			var err error

			handler := &VectorCacheBuildHandler{
				AdminUserId:   c.AdminUserId,
				ModelConfigId: c.ModelConfigId,
				UseModel:      c.UseModel,
				Sentence:      text,
			}

			// Get or generate vector from cache, cached for 5 minutes
			err = lib_redis.GetCacheWithBuild(define.Redis, handler, &vector, time.Hour*24)

			resultChan <- vectorResult{
				index:  idx,
				vector: vector,
				err:    err,
			}
		}(i, sentence)
	}

	// Collect results
	var lastError error
	validVectors := 0

	for i := 0; i < len(sentences); i++ {
		r := <-resultChan
		if r.err != nil {
			lastError = r.err
			continue
		}

		vectors[r.index] = r.vector
		validVectors++
	}

	// Filter out empty vectors
	resultVectors := make([][]float64, 0, validVectors)
	for _, v := range vectors {
		if len(v) > 0 {
			resultVectors = append(resultVectors, v)
		}
	}

	if len(resultVectors) == 0 {
		if lastError != nil {
			return nil, lastError
		}
		return nil, errors.New("failed to generate embeddings for all sentences")
	}

	return resultVectors, nil
}

// calculateDifferences Calculate similarity differences between adjacent vectors
func (c *SemanticSplitter) calculateDifferences(vectors [][]float64) []float64 {
	var differences []float64

	for i := 0; i < len(vectors)-1; i++ {
		similarity := cosineSimilarity(vectors[i], vectors[i+1])
		difference := 1.0 - similarity
		differences = append(differences, difference)
	}

	return differences
}

// findBreakpoints Find breakpoints indicating where to split
func (c *SemanticSplitter) findBreakpoints(differences []float64) []int {
	if len(differences) == 0 {
		return []int{}
	}

	// Determine threshold using percentile method
	threshold := c.getThresholdPercentile(differences)

	// Find positions where difference values exceed the threshold as breakpoints
	var breakpoints []int
	for i, diff := range differences {
		if diff > threshold {
			breakpoints = append(breakpoints, i+1) // i+1 means breaking between sentence i and i+1
		}
	}

	return breakpoints
}

// getThresholdPercentile Get percentile threshold of difference values
func (c *SemanticSplitter) getThresholdPercentile(differences []float64) float64 {
	if len(differences) == 0 {
		return 0
	}

	// Copy difference value array and sort
	sortedDiffs := make([]float64, len(differences))
	copy(sortedDiffs, differences)
	sort.Float64s(sortedDiffs)

	// Calculate threshold position
	percentile := float64(c.SemanticChunkThreshold) / 100.0
	idx := int(math.Floor(float64(len(sortedDiffs)-1) * percentile))

	// Return the corresponding threshold
	return sortedDiffs[idx]
}

// formChunks Form text chunks based on breakpoints
func (c *SemanticSplitter) formChunks(sentences []string, breakpoints []int) []string {
	var chunks []string
	if len(breakpoints) == 0 {
		return []string{strings.Join(sentences, "")}
	}
	startIdx := 0
	for _, bp := range breakpoints {
		chunk := strings.Join(sentences[startIdx:bp], "")
		if len(chunk) > 0 {
			chunks = append(chunks, chunk)
		}
		startIdx = bp
	}
	if startIdx < len(sentences) {
		lastChunk := strings.Join(sentences[startIdx:], "")
		if len(lastChunk) > 0 {
			chunks = append(chunks, lastChunk)
		}
	}
	return chunks
}

// splitByBinaryDivision Use binary division to split sentence array into appropriately sized chunks
func (c *SemanticSplitter) splitByBinaryDivision(sentences []string) []string {
	var chunks []string
	var currentChunk []string
	var currentLength int

	// Try to combine sentences into chunks, not exceeding maximum chunk size
	for _, sentence := range sentences {
		sentenceLength := utf8.RuneCountInString(sentence)

		// If a single sentence exceeds the maximum length, it needs to be split by characters
		if sentenceLength > c.SemanticChunkSize {
			if len(currentChunk) > 0 {
				chunks = append(chunks, strings.Join(currentChunk, ""))
				currentChunk = nil
				currentLength = 0
			}
			// Perform character splitting on overly long sentences
			subChunks := c.splitSentenceByChars(sentence)
			chunks = append(chunks, subChunks...)
			continue
		}

		// Determine if adding the current sentence would exceed the maximum length
		if currentLength+sentenceLength <= c.SemanticChunkSize {
			currentChunk = append(currentChunk, sentence)
			currentLength += sentenceLength
		} else {
			// Current chunk is full, save and start a new chunk
			if len(currentChunk) > 0 {
				chunks = append(chunks, strings.Join(currentChunk, ""))
				currentChunk = []string{sentence}
				currentLength = sentenceLength
			}
		}
	}

	// Add the last chunk
	if len(currentChunk) > 0 {
		chunks = append(chunks, strings.Join(currentChunk, ""))
	}

	return chunks
}

// splitSentenceByChars Split a single sentence by characters
func (c *SemanticSplitter) splitSentenceByChars(sentence string) []string {
	var result []string
	runes := []rune(sentence)

	// Split using maximum chunk size
	for i := 0; i < len(runes); i += c.SemanticChunkSize {
		end := i + c.SemanticChunkSize
		if end > len(runes) {
			end = len(runes)
		}
		result = append(result, string(runes[i:end]))
		if end == len(runes) {
			break
		}
	}

	return result
}

// cosineSimilarity Calculate cosine similarity between two vectors
func cosineSimilarity(v1, v2 []float64) float64 {
	if len(v1) != len(v2) {
		return 0
	}
	var dotProduct, magnitude1, magnitude2 float64
	for i := 0; i < len(v1); i++ {
		dotProduct += v1[i] * v2[i]
		magnitude1 += v1[i] * v1[i]
		magnitude2 += v2[i] * v2[i]
	}
	magnitude1 = math.Sqrt(magnitude1)
	magnitude2 = math.Sqrt(magnitude2)

	if magnitude1 == 0 || magnitude2 == 0 {
		return 0
	}

	return dotProduct / (magnitude1 * magnitude2)
}
