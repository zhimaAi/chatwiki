// Copyright © 2016- 2024 Sesame Network Technology all right reserved

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

// SemanticSplitter 语义分割客户端
type SemanticSplitter struct {
	SemanticChunkSize      int                       // 最大分块大小（按字符数）
	SemanticChunkOverlap   int                       // 分块重叠大小（按字符数）
	SemanticChunkThreshold int                       // 语义断点阈值（百分比，1-100）
	AdminUserId            int                       // 用户ID
	ModelConfigId          int                       // 模型配置ID
	UseModel               string                    // 使用的模型
	GoRoutineNum           int                       // 协程数量
	GetVector              func() ([]float64, error) // 获取向量的函数
}

// NewSemanticSplitterClient 创建一个新的语义分割客户端
func NewSemanticSplitterClient() *SemanticSplitter {
	return &SemanticSplitter{}
}

// SplitText 分割文本
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
	// 添加重叠内容
	result = c.addOverlappingContent(result)
	return result, nil
}

// addOverlappingContent 为分段添加重叠内容
func (c *SemanticSplitter) addOverlappingContent(chunks []string) []string {
	if len(chunks) <= 1 || c.SemanticChunkOverlap <= 0 {
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
			if overlapSize+sentenceSize <= c.SemanticChunkOverlap || len(overlappingSentences) == 0 {
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

// semanticSplit 语义分割
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

// splitIntoSentences 将文本分割成句子
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

// VectorCacheBuildHandler 向量缓存处理程序
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

	// 解析JSON字符串到float64数组
	var vector []float64
	if err := json.Unmarshal([]byte(embedding), &vector); err != nil {
		logs.Error("Get vector error: %v", err)
		return nil, err
	}
	return vector, nil
}

// generateEmbeddings 为每个句子生成向量表示
func (c *SemanticSplitter) generateEmbeddings(sentences []string) ([][]float64, error) {
	var vectors = make([][]float64, len(sentences))

	// 使用通道来控制并发
	type vectorResult struct {
		index  int
		vector []float64
		err    error
	}
	resultChan := make(chan vectorResult, len(sentences))

	// 创建信号量来限制并发数量
	sem := make(chan struct{}, c.GoRoutineNum)

	// 启动协程生成向量
	for i, sentence := range sentences {
		go func(idx int, text string) {
			sem <- struct{}{}        // 获取信号量
			defer func() { <-sem }() // 释放信号量

			// 尝试从缓存获取向量
			var vector []float64
			var err error

			handler := &VectorCacheBuildHandler{
				AdminUserId:   c.AdminUserId,
				ModelConfigId: c.ModelConfigId,
				UseModel:      c.UseModel,
				Sentence:      text,
			}

			// 从缓存获取或生成向量，缓存5分钟
			err = lib_redis.GetCacheWithBuild(define.Redis, handler, &vector, time.Hour*24)

			resultChan <- vectorResult{
				index:  idx,
				vector: vector,
				err:    err,
			}
		}(i, sentence)
	}

	// 收集结果
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

	// 过滤掉空向量
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

// calculateDifferences 计算相邻向量之间的相似度差异
func (c *SemanticSplitter) calculateDifferences(vectors [][]float64) []float64 {
	var differences []float64

	for i := 0; i < len(vectors)-1; i++ {
		similarity := cosineSimilarity(vectors[i], vectors[i+1])
		difference := 1.0 - similarity
		differences = append(differences, difference)
	}

	return differences
}

// findBreakpoints 找出应该在哪里分割的断点
func (c *SemanticSplitter) findBreakpoints(differences []float64) []int {
	if len(differences) == 0 {
		return []int{}
	}

	// 使用百分位法确定阈值
	threshold := c.getThresholdPercentile(differences)

	// 寻找差异值大于阈值的位置作为断点
	var breakpoints []int
	for i, diff := range differences {
		if diff > threshold {
			breakpoints = append(breakpoints, i+1) // i+1 表示在第 i 和 i+1 句子之间断开
		}
	}

	return breakpoints
}

// getThresholdPercentile 获取差异值的百分位阈值
func (c *SemanticSplitter) getThresholdPercentile(differences []float64) float64 {
	if len(differences) == 0 {
		return 0
	}

	// 复制差异值数组并排序
	sortedDiffs := make([]float64, len(differences))
	copy(sortedDiffs, differences)
	sort.Float64s(sortedDiffs)

	// 计算阈值位置
	percentile := float64(c.SemanticChunkThreshold) / 100.0
	idx := int(math.Floor(float64(len(sortedDiffs)-1) * percentile))

	// 返回对应的阈值
	return sortedDiffs[idx]
}

// formChunks 根据断点形成文本块
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

// splitByBinaryDivision 使用二分法将句子数组分割成适当大小的块
func (c *SemanticSplitter) splitByBinaryDivision(sentences []string) []string {
	var chunks []string
	var currentChunk []string
	var currentLength int

	// 尝试将句子组合成块，不超过最大分块大小
	for _, sentence := range sentences {
		sentenceLength := utf8.RuneCountInString(sentence)

		// 如果单个句子就超过了最大长度，需要按字符分割
		if sentenceLength > c.SemanticChunkSize {
			if len(currentChunk) > 0 {
				chunks = append(chunks, strings.Join(currentChunk, ""))
				currentChunk = nil
				currentLength = 0
			}
			// 对超长句子进行字符分割
			subChunks := c.splitSentenceByChars(sentence)
			chunks = append(chunks, subChunks...)
			continue
		}

		// 判断添加当前句子是否会超出最大长度
		if currentLength+sentenceLength <= c.SemanticChunkSize {
			currentChunk = append(currentChunk, sentence)
			currentLength += sentenceLength
		} else {
			// 当前块已满，保存并开始新块
			if len(currentChunk) > 0 {
				chunks = append(chunks, strings.Join(currentChunk, ""))
				currentChunk = []string{sentence}
				currentLength = sentenceLength
			}
		}
	}

	// 添加最后一个块
	if len(currentChunk) > 0 {
		chunks = append(chunks, strings.Join(currentChunk, ""))
	}

	return chunks
}

// splitSentenceByChars 将单个句子按字符分割
func (c *SemanticSplitter) splitSentenceByChars(sentence string) []string {
	var result []string
	runes := []rune(sentence)

	// 使用最大分块大小进行分割
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

// cosineSimilarity 计算两个向量之间的余弦相似度
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
