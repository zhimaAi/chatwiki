// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/zhimaAi/go_tools/logs"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"
)

func DetectMCPTransportType(url string) (int, error) {
	if strings.HasSuffix(url, "sse") {
		return define.McpClientTypeSse, nil
	}

	c := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Accept", "text/event-stream")

	resp, err := c.Do(req)
	if err != nil {
		logs.Error(err.Error())
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logs.Error(err.Error())
		}
	}(resp.Body)

	ct := resp.Header.Get("Content-Type")
	if strings.Contains(ct, "text/event-stream") {
		return define.McpClientTypeSse, nil
	}

	return define.McpClientTypeHttp, nil
}

func NewMcpClient(ctx context.Context, clientType int, serverURL string, headers string) (*client.Client, error) {
	var c *client.Client
	var err error

	var mapHeaders map[string]string
	if headers != `` {
		if err := json.Unmarshal([]byte(headers), &mapHeaders); err != nil {
			return nil, err
		}
	}

	switch clientType {
	case define.McpClientTypeSse:
		c, err = client.NewSSEMCPClient(serverURL, transport.WithHeaders(mapHeaders))
	case define.McpClientTypeHttp:
		c, err = client.NewStreamableHttpClient(serverURL, transport.WithHTTPHeaders(mapHeaders))
	default:
		return nil, fmt.Errorf("unknown client type: %d", clientType)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return c, nil
}

func ListTools(ctx context.Context, c *client.Client) ([]mcp.Tool, error) {
	err := c.Start(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start client: %w", err)
	}
	initReq := mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			ClientInfo: mcp.Implementation{
				Name:    "chatwiki-client",
				Version: "0.1.0",
			},
			Capabilities: mcp.ClientCapabilities{
				Roots: &struct {
					ListChanged bool `json:"listChanged,omitempty"`
				}{
					ListChanged: false,
				},
			},
		},
	}
	_, err = c.Initialize(ctx, initReq)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize client: %w", err)
	}

	listReq := mcp.ListToolsRequest{
		PaginatedRequest: mcp.PaginatedRequest{
			Request: mcp.Request{
				Method: "tools/list",
			},
			Params: mcp.PaginatedParams{
				Cursor: "", // 初始为空
			},
		},
	}
	tools, err := c.ListTools(ctx, listReq)
	if err != nil {
		return nil, err
	}
	return tools.Tools, nil
}

func ValidateMcpToolArguments(tool mcp.Tool, args map[string]interface{}) error {
	schema := tool.InputSchema

	// Check required properties
	if schema.Required != nil {
		for _, required := range schema.Required {
			if _, exists := args[required]; !exists {
				return fmt.Errorf("missing required argument: %s", required)
			}
		}
	}

	if schema.Properties != nil {
		for name, value := range args {
			propSchemaAny, exists := schema.Properties[name]
			if !exists {
				return fmt.Errorf("unknown argument: %s", name)
			}

			propSchema, ok := propSchemaAny.(map[string]any)
			if !ok {
				return fmt.Errorf("internal error: argument schema for %s is not a map[string]any", name)
			}

			if err := validateMcpValue(value, propSchema); err != nil {
				return fmt.Errorf("invalid argument %s: %w", name, err)
			}
		}
	}

	return nil
}

func validateMcpValue(value interface{}, schema map[string]any) error {
	schemaType, ok := schema["type"].(string)
	if !ok {
		return fmt.Errorf("schema missing type")
	}

	switch schemaType {
	case "string":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("expected string, got %T", value)
		}
	//case "number":
	//	if _, ok := value.(float64); !ok {
	//		return fmt.Errorf("expected number, got %T", value)
	//	}
	//case "integer":
	//	if _, ok := value.(float64); !ok {
	//		return fmt.Errorf("expected integer, got %T", value)
	//	}
	case "boolean":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("expected boolean, got %T", value)
		}
	case "array":
		if _, ok := value.([]interface{}); !ok {
			return fmt.Errorf("expected array, got %T", value)
		}
	case "object":
		if _, ok := value.(map[string]interface{}); !ok {
			return fmt.Errorf("expected object, got %T", value)
		}
	}

	return nil
}
func NormalizeArgumentsBySchema(args map[string]any, schema *mcp.ToolInputSchema) (map[string]any, error) {
	if schema == nil || schema.Properties == nil {
		return args, nil
	}

	normalized := make(map[string]any, len(args))
	requiredSet := make(map[string]struct{}, len(schema.Required))

	// 建立 required 字段快速查找表
	for _, name := range schema.Required {
		requiredSet[name] = struct{}{}
	}

	for name, val := range args {
		propSchemaAny, exists := schema.Properties[name]
		if !exists {
			// schema 中未定义的字段直接跳过
			continue
		}

		propSchema, ok := propSchemaAny.(map[string]any)
		if !ok {
			normalized[name] = val
			continue
		}

		// 如果是非必填字段并且值为空，则跳过
		if _, isRequired := requiredSet[name]; !isRequired {
			if isEmptyValue(val) {
				continue
			}
		}

		schemaType, _ := propSchema["type"].(string)

		switch schemaType {
		case "object", "array":
			// 若值是字符串，看起来像 JSON，则尝试解析
			if str, ok := val.(string); ok {
				var parsed any
				if err := json.Unmarshal([]byte(str), &parsed); err == nil {
					normalized[name] = parsed
				} else {
					return nil, fmt.Errorf("argument %q should be valid JSON: %v", name, err)
				}
			} else {
				normalized[name] = val
			}
		default:
			normalized[name] = val
		}
	}

	return normalized, nil
}

// 判断是否为空值（适用于各种可能类型）
func isEmptyValue(v any) bool {
	if v == nil {
		return true
	}

	switch val := v.(type) {
	case string:
		return strings.TrimSpace(val) == ""
	case []any:
		return len(val) == 0
	case map[string]any:
		return len(val) == 0
	case bool:
		return false
	case float64, int, int64:
		return false
	default:
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Slice, reflect.Map, reflect.Array:
			return rv.Len() == 0
		case reflect.Ptr, reflect.Interface:
			return rv.IsNil()
		default:
			return false
		}
	}
}

// CallTool 调用MCP工具
func CallTool(ctx context.Context, c *client.Client, selectedTool mcp.Tool, arguments map[string]any) (string, error) {
	err := c.Start(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to start client: %w", err)
	}
	initReq := mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			ClientInfo: mcp.Implementation{
				Name:    "chatwiki-client",
				Version: "0.1.0",
			},
			Capabilities: mcp.ClientCapabilities{
				Roots: &struct {
					ListChanged bool `json:"listChanged,omitempty"`
				}{
					ListChanged: false,
				},
			},
		},
	}
	_, err = c.Initialize(ctx, initReq)
	if err != nil {
		return "", fmt.Errorf("failed to initialize client: %w", err)
	}

	normalizedArgs, err := NormalizeArgumentsBySchema(arguments, &selectedTool.InputSchema)
	if err != nil {
		return "", fmt.Errorf("参数格式转换错误: %v", err)
	}

	// 2. 校验参数合法性
	if err := ValidateMcpToolArguments(selectedTool, normalizedArgs); err != nil {
		return "", fmt.Errorf("参数校验失败: %v", err)
	}

	// 3. 执行调用
	result, err := c.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      selectedTool.Name,
			Arguments: normalizedArgs,
		},
	})

	if err != nil {
		return "", fmt.Errorf("调用mcp工具出错: %v", err.Error())
	}
	//if result.IsError {
	//	return "", errors.New("调用mcp工具失败")
	//}

	// 声明用于存储最终文本内容的变量
	var strContent string

	for _, content := range result.Content {
		// 使用类型断言检查 content 的具体类型
		switch c := content.(type) {
		case mcp.TextContent:
			strContent = c.Text
			goto FoundContent
		case mcp.ImageContent:
			// TODO: 后续兼容 ImageContent 类型的逻辑
		case mcp.AudioContent:
			// TODO: 后续兼容 AudioContent 类型的逻辑
		default:
			return "", fmt.Errorf("未知返回类型: %T", content)
		}
	}

FoundContent:
	if strContent == "" {
		return "", errors.New("调用mcp工具返回结果中未找到TextContent内容")
	}

	return strContent, nil
}
