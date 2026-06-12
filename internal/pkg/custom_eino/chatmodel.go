// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package custom_eino

import (
	"context"
	"errors"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

const (
	defaultType            = "Custom"
	defaultModelName       = "custom-local"
	defaultStreamChunkSize = 4
)

type RuntimeOptions struct {
	Common *model.Options
}

type GenerateFunc func(ctx context.Context, input []*schema.Message, opts RuntimeOptions) (*schema.Message, error)
type StreamFunc func(ctx context.Context, input []*schema.Message, opts RuntimeOptions) (*schema.StreamReader[*schema.Message], error)

type ChatModelConfig struct {
	Name  string
	Model string

	Temperature *float32
	MaxTokens   *int
	TopP        *float32
	Stop        []string

	StreamChunkSize int
	Generate        GenerateFunc
	Stream          StreamFunc
}

type ChatModel struct {
	name            string
	modelName       string
	temperature     *float32
	maxTokens       *int
	topP            *float32
	stop            []string
	streamChunkSize int
	generate        GenerateFunc
	stream          StreamFunc
}

func NewChatModel(_ context.Context, cfg *ChatModelConfig) (*ChatModel, error) {
	if cfg == nil {
		cfg = &ChatModelConfig{}
	}

	name := valueOrDefaultString(cfg.Name, defaultType)
	modelName := valueOrDefaultString(cfg.Model, defaultModelName)
	streamChunkSize := cfg.StreamChunkSize
	if streamChunkSize <= 0 {
		streamChunkSize = defaultStreamChunkSize
	}

	generate := cfg.Generate
	if generate == nil {
		return nil, errors.New("custom chat model generate function is required")
	}

	return &ChatModel{
		name:            name,
		modelName:       modelName,
		temperature:     clonePtr(cfg.Temperature),
		maxTokens:       clonePtr(cfg.MaxTokens),
		topP:            clonePtr(cfg.TopP),
		stop:            cloneSlice(cfg.Stop),
		streamChunkSize: streamChunkSize,
		generate:        generate,
		stream:          cfg.Stream,
	}, nil
}

func (cm *ChatModel) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	out, err := cm.generate(ctx, input, cm.runtimeOptions(opts...))
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, errors.New("custom chat model returned nil message")
	}
	if out.Role == "" {
		cp := *out
		cp.Role = schema.Assistant
		out = &cp
	}
	return out, nil
}

func (cm *ChatModel) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	if cm.stream != nil {
		out, err := cm.stream(ctx, input, cm.runtimeOptions(opts...))
		if err != nil {
			return nil, err
		}
		if out == nil {
			return nil, errors.New("custom chat model returned nil stream")
		}
		return out, nil
	}

	out, err := cm.Generate(ctx, input, opts...)
	if err != nil {
		return nil, err
	}
	return streamFromMessage(out, cm.streamChunkSize), nil
}

func (cm *ChatModel) GetType() string {
	return cm.name
}

func (cm *ChatModel) runtimeOptions(opts ...model.Option) RuntimeOptions {
	common := model.GetCommonOptions(&model.Options{
		Model:       ptr(cm.modelName),
		Temperature: clonePtr(cm.temperature),
		MaxTokens:   clonePtr(cm.maxTokens),
		TopP:        clonePtr(cm.topP),
		Stop:        cloneSlice(cm.stop),
	}, opts...)

	common.Stop = cloneSlice(common.Stop)
	common.Tools = cloneSlice(common.Tools)
	return RuntimeOptions{Common: common}
}

func streamFromMessage(msg *schema.Message, chunkSize int) *schema.StreamReader[*schema.Message] {
	if msg == nil {
		return schema.StreamReaderFromArray([]*schema.Message{})
	}
	if chunkSize <= 0 {
		chunkSize = defaultStreamChunkSize
	}

	if msg.Content == "" ||
		msg.ReasoningContent != "" ||
		len(msg.ToolCalls) > 0 ||
		len(msg.MultiContent) > 0 ||
		len(msg.UserInputMultiContent) > 0 ||
		len(msg.AssistantGenMultiContent) > 0 {
		return schema.StreamReaderFromArray([]*schema.Message{msg})
	}

	runes := []rune(msg.Content)
	chunks := make([]*schema.Message, 0, (len(runes)+chunkSize-1)/chunkSize)
	for start := 0; start < len(runes); start += chunkSize {
		end := min(start+chunkSize, len(runes))
		chunk := &schema.Message{
			Role:    msg.Role,
			Name:    msg.Name,
			Content: string(runes[start:end]),
		}
		if end == len(runes) {
			chunk.ResponseMeta = msg.ResponseMeta
			chunk.Extra = msg.Extra
		}
		chunks = append(chunks, chunk)
	}
	return schema.StreamReaderFromArray(chunks)
}

func valueOrDefaultString(value string, fallback string) string {
	if value = strings.TrimSpace(value); value != "" {
		return value
	}
	return fallback
}

func ptr[T any](value T) *T {
	return &value
}

func clonePtr[T any](in *T) *T {
	if in == nil {
		return nil
	}
	out := *in
	return &out
}

func cloneSlice[T any](in []T) []T {
	if in == nil {
		return nil
	}
	return append([]T(nil), in...)
}
