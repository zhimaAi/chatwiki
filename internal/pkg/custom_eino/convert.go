// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package custom_eino

import (
	"github.com/cloudwego/eino/schema"
	"github.com/eino-contrib/jsonschema"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
	"github.com/zhimaAi/llm_adaptor/basics"
)

func StructConvert[In, Out any](in In) (out Out, err error) {
	var jsonStr string
	jsonStr, err = tool.JsonEncode(in)
	if err != nil {
		return
	}
	err = tool.JsonDecodeUseNumber(jsonStr, &out)
	return
}

func ConvertMessage(message schema.Message) (adaptor.ZhimaChatCompletionMessage, error) {
	return StructConvert[schema.Message, adaptor.ZhimaChatCompletionMessage](message)
}

func ConvertTools(info schema.ToolInfo) (result adaptor.FunctionTool, err error) {
	params, err := info.ParamsOneOf.ToJSONSchema()
	if err != nil {
		return
	}
	result.Name = info.Name
	result.Description = info.Desc
	if params != nil {
		result.Parameters = adaptor.Parameters{
			Type:       params.Type,
			Properties: params.Properties,
			Required:   params.Required,
		}
	}
	return
}

func ConvertToolCalls(toolCalls basics.ToolCalls) ([]schema.ToolCall, error) {
	return StructConvert[basics.ToolCalls, []schema.ToolCall](toolCalls)
}

func ConvertChatResp(chatResp adaptor.ZhimaChatCompletionResponse) *schema.Message {
	var content string
	var toolCalls []schema.ToolCall
	if len(chatResp.ToolCalls) > 0 {
		toolCalls, _ = ConvertToolCalls(chatResp.ToolCalls)
	} else {
		content = chatResp.Result
		toolCalls = nil
	}
	msg := schema.AssistantMessage(content, toolCalls)
	msg.ReasoningContent = chatResp.ReasoningContent
	if len(toolCalls) > 0 {
		msg.ResponseMeta = &schema.ResponseMeta{FinishReason: `tool_calls`}
	}
	return msg
}

func ConvertProperties(properties any) (*orderedmap.OrderedMap[string, *jsonschema.Schema], error) {
	return StructConvert[any, *orderedmap.OrderedMap[string, *jsonschema.Schema]](properties)
}
