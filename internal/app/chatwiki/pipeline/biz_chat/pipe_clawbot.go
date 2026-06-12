// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package biz_chat

import (
	"chatwiki/internal/app/chatwiki/common"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/work_flow"
	"chatwiki/internal/pkg/custom_eino"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/llm_runner"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/filesystem"
	"github.com/cloudwego/eino/adk/middlewares/skill"
	"github.com/cloudwego/eino/adk/prebuilt/deep"
	einotool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/eino-contrib/jsonschema"
	"github.com/gin-contrib/sse"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

func buildFileOperationPrompt(robot msql.Params) string {
	prompt := `When using file operation tools (ls_info, read, grep_raw, glob_info, write, edit, execute), follow these rules:
1. Read-only tools: ls_info, read, grep_raw, glob_info.
2. Write tools: write, edit.
3. Command execution tool: execute.
- Execute receives one argv-style command only. Do not use shell syntax such as cd, &&, |, ;, &, <, >, backticks, $, redirects, command substitution, or multi-line commands.
- Execute command names must be one of: cat, find, grep, head, jq, ls, node, npm, pwd, python3, rg, tail, wc. Use the bare command name only, not a path to the command binary.
- Execute paths must be relative paths under ` + define.PrivateWorkDir + `. Do not use absolute paths such as /workspace/..., ~/..., or Windows drive paths, and do not use parent-directory traversal.
- For Python and Node, run files from ` + define.PrivateWorkDir + `; use python3 only and do not use python3 -c, python3 -m, node -e, node -p, --eval, or --print. For find, do not use -delete, -exec, or -execdir.
- Valid execute examples: python3 ` + define.PrivateWorkDir + `/test_python.py; ls ` + define.PrivateWorkDir + `; rg "keyword" ` + define.PrivateWorkDir + `.
- Invalid execute examples: cd ` + define.PrivateWorkDir + ` && python3 test_python.py; python3 /workspace/` + define.PrivateWorkDir + `/test_python.py.
4. Accessible directories and permissions:
- Public skills directory, read-only: ` + define.PublicSkillsDir + `
- Private skills directory, read-only: ` + define.PrivateSkillsDir + `
- Workspace directory, read, write, and command execution allowed: ` + define.PrivateWorkDir + `
5. Always use relative paths as tool arguments. Refuse access when a path is absolute, escapes the allowed directories, or uses parent-directory traversal.
6. Use the smallest sufficient operation. Do not write files or execute commands unless the task explicitly requires it.

Valid path examples:
- ` + define.PublicSkillsDir + `/skill_name/reference/xxx.md
- ` + define.PrivateSkillsDir + `/skills/skill_name/reference/xxx.md
- ` + define.PrivateWorkDir + `/temp.log`
	// replace robot_key
	return strings.ReplaceAll(prompt, `<robot_key>`, robot[`robot_key`])
}

func ValidateFile(robot msql.Params, op custom_eino.FileOperation, filePath string) error {
	filePath = strings.ReplaceAll(strings.TrimSpace(filePath), `\`, `/`)
	if filePath == `` {
		return fmt.Errorf(`file path is required`)
	}
	if strings.HasPrefix(filePath, `/`) {
		return fmt.Errorf(`absolute path is not allowed: %s`, filePath)
	}
	for _, part := range strings.Split(filePath, `/`) {
		if part == `..` {
			return fmt.Errorf(`parent-directory traversal is not allowed: %s`, filePath)
		}
	}
	var allowedPrefixes []string
	switch op {
	case custom_eino.FileOperationLsInfo, custom_eino.FileOperationRead, custom_eino.FileOperationGrepRaw, custom_eino.FileOperationGlobInfo:
		allowedPrefixes = []string{
			define.PublicSkillsDir,
			strings.ReplaceAll(define.PrivateSkillsDir, `<robot_key>`, robot[`robot_key`]),
			strings.ReplaceAll(define.PrivateWorkDir, `<robot_key>`, robot[`robot_key`]),
		}
		if !cast.ToBool(robot[`query_local_docs_close`]) {
			allowedPrefixes = append(allowedPrefixes, strings.ReplaceAll(define.PrivateFileDir, `<robot_key>`, robot[`robot_key`]))
		}
	case custom_eino.FileOperationWrite, custom_eino.FileOperationEdit:
		allowedPrefixes = []string{
			strings.ReplaceAll(define.PrivateWorkDir, `<robot_key>`, robot[`robot_key`]),
		}
	default:
		return fmt.Errorf(`%s operation is not allowed`, op)
	}
	for _, prefix := range allowedPrefixes {
		if filePath == prefix || strings.HasPrefix(filePath, prefix+`/`) {
			return nil
		}
	}
	return fmt.Errorf(`%s is not allowed for %s`, filePath, op)
}

func ValidateCommand(robotKey string, command string) error {
	args, err := llm_runner.PrepareCommand(command)
	if err != nil {
		return err
	}
	workDir := strings.ReplaceAll(define.PrivateWorkDir, `<robot_key>`, robotKey)
	return llm_runner.ValidateClawbotPathScope(args, workDir)
}

func doApplicationTypeClaw(in *ChatInParam, out *ChatOutParam) error {
	// set language
	_ = adk.SetLanguage(adk.LanguageChinese)
	// init dirs
	common.InitClawbotDirs(in.params.Robot[`robot_key`])
	// ChatModel
	ctx := context.Background()
	cm, err := custom_eino.NewChatModel(ctx, &custom_eino.ChatModelConfig{
		Generate: func(ctx context.Context, input []*schema.Message, opts custom_eino.RuntimeOptions) (*schema.Message, error) {
			return Generate(ctx, input, opts, in, out)
		},
		Stream: func(ctx context.Context, input []*schema.Message, opts custom_eino.RuntimeOptions) (*schema.StreamReader[*schema.Message], error) {
			return Stream(ctx, input, opts, in, out)
		},
	})
	if err != nil {
		return err
	}
	// Backend
	backend, err := custom_eino.NewBackend(ctx, &custom_eino.BackendConfig{
		ValidateFile: func(op custom_eino.FileOperation, path string) error {
			in.Stream(sse.Event{Event: `FileOperation`, Data: map[string]any{`op`: op, `path`: path}})
			return ValidateFile(in.params.Robot, op, path)
		},
		ValidateCommand: func(command string) error {
			in.Stream(sse.Event{Event: `ExecuteCommand`, Data: map[string]any{`command`: command}})
			return ValidateCommand(in.params.Robot[`robot_key`], command)
		},
		ExecuteCommand: func(command string) (*filesystem.ExecuteResponse, error) {
			resp := llm_runner.RpcExecuteRun(define.Config.WebService[`llm_runner_host`], command)
			if resp.IsError {
				if resp.ExitCode < 0 { // server execution error
					return &filesystem.ExecuteResponse{Output: resp.ErrorMsg, ExitCode: &resp.ExitCode}, nil
				}
				return nil, errors.New(resp.ErrorMsg)
			}
			return &filesystem.ExecuteResponse{Output: resp.Output, ExitCode: &resp.ExitCode}, nil
		},
	})
	if err != nil {
		return err
	}
	// BaseTool
	tools := make([]einotool.BaseTool, 0)
	kbSearchTool := newKbsearchTool(in.params, in.sessionId)
	if kbSearchTool != nil {
		tools = append(tools, kbSearchTool)
	}
	// related workflow
	workFlowFuncCall, _ := work_flow.BuildFunctionTools(in.params.Lang, in.params.Robot)
	for _, t := range workFlowFuncCall {
		result, err := custom_eino.ConvertProperties(t.Parameters.Properties)
		if err != nil {
			logs.Error(`ConvertProperties:` + err.Error())
			continue
		}
		sc := &jsonschema.Schema{Properties: result, Type: t.Parameters.Type, Required: t.Parameters.Required}
		tools = append(tools, custom_eino.BuildWorkFlowTool(t.Name, t.Description, schema.NewParamsOneOfByJSONSchema(sc),
			func(functionTool adaptor.FunctionToolCall) (string, error) {
				return doClawbotRelationWorkFlow(functionTool, in, out)
			}))
	}
	// Middleware
	handlers := make([]adk.ChatModelAgentMiddleware, 0)
	skillMiddleware, err := newSkillMiddleware(ctx, backend, in.params.Robot)
	if err != nil {
		return err
	}
	handlers = append(handlers, skillMiddleware)
	// Runner
	agent, err := deep.New(ctx, &deep.Config{
		ChatModel: cm,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools:               tools,
				ExecuteSequentially: true,
			},
		},
		MaxIteration:           20,
		Backend:                backend,
		Shell:                  backend,
		WithoutGeneralSubAgent: true,
		Handlers:               handlers,
		ModelRetryConfig: &adk.ModelRetryConfig{
			MaxRetries: 3,
			ShouldRetry: func(_ context.Context, retryCtx *adk.RetryContext) *adk.RetryDecision {
				return &adk.RetryDecision{Retry: retryCtx.Err != nil, RewriteError: retryCtx.Err}
			},
		},
	})
	if err != nil {
		return err
	}
	runner := adk.NewRunner(ctx, adk.RunnerConfig{Agent: agent, EnableStreaming: in.useStream})
	// history
	out.debugLog = append(out.debugLog, map[string]string{`type`: `prompt`, `content`: `【prompt】`})
	var history []*schema.Message
	contextList := common.BuildChatContextPair(in.params.Openid, cast.ToInt(in.params.Robot[`id`]),
		in.dialogueId, int(out.cMsgId), cast.ToInt(in.params.Robot[`context_pair`]))
	for i := range contextList {
		history = append(history, schema.UserMessage(contextList[i][`question`]))
		history = append(history, schema.AssistantMessage(contextList[i][`answer`], nil))
		out.debugLog = append(out.debugLog, map[string]string{`type`: `context_qa`, `question`: contextList[i][`question`], `answer`: contextList[i][`answer`]})
	}
	history = append(history, schema.UserMessage(in.params.Question))
	out.debugLog = append(out.debugLog, map[string]string{`type`: `cur_question`, `content`: in.params.Question})
	// events
	requestStartTime := time.Now()
	events := runner.Run(ctx, history)
	var sb strings.Builder
	for {
		event, ok := events.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			out.Error = event.Err
			common.SendDefaultUnknownQuestionPrompt(in.params, out.Error.Error(), in.chanStream, &out.content)
			return nil
		}
		if event.Action != nil && event.Action.Interrupted != nil {
			continue
		}
		if event.Output == nil || event.Output.MessageOutput == nil {
			continue
		}
		message, err := ConcatStreamMessages(in, event.Output.MessageOutput)
		if err != nil {
			out.Error = err
			common.SendDefaultUnknownQuestionPrompt(in.params, out.Error.Error(), in.chanStream, &out.content)
			return nil
		}
		role := event.Output.MessageOutput.Role
		if role == schema.Assistant || role == `` {
			for _, call := range message.ToolCalls {
				in.Stream(sse.Event{Event: `tool_call`, Data: call.Function})
				out.debugLog = append(out.debugLog, map[string]string{`type`: `tool_call`, `content`: tool.JsonEncodeNoError(call.Function)})
			}
			sb.WriteString(message.Content)
		} else if role == schema.Tool {
			sb.Reset() // clear the result from the last time
			in.Stream(sse.Event{Event: `tool_result`, Data: message})
			out.debugLog = append(out.debugLog, map[string]string{`type`: `tool_result`, `content`: tool.JsonEncodeNoError(message.Content)})
		}
	}
	out.chatResp.Result = sb.String()
	out.requestTime = time.Now().Sub(requestStartTime).Milliseconds()
	out.content = out.chatResp.Result
	return nil
}

func buildSystemPrompt(system string, in *ChatInParam, out *ChatOutParam) string {
	if len(in.params.Prompt) == 0 { //no custom is used
		prompt, promptStruct := in.params.Robot[`prompt`], in.params.Robot[`prompt_struct`]
		common.ReplaceChatVariables(in.params.Lang, in.sessionId, in.params.WorkFlowGlobal, &prompt, &promptStruct)
		in.params.Prompt = prompt
	}
	result := fmt.Sprintf("%s\n%s\n%s", in.params.Prompt, system, buildFileOperationPrompt(in.params.Robot))
	replacePromptPlaceholder(result, out)
	return result
}

func replacePromptPlaceholder(prompt string, out *ChatOutParam) {
	for idx := range out.debugLog {
		log, ok := out.debugLog[idx].(map[string]string)
		if !ok {
			continue
		}
		if log[`type`] == `prompt` && log[`content`] == `【prompt】` {
			log[`content`] = prompt
		}
		out.debugLog[idx] = log
	}
}

func Generate(ctx context.Context, input []*schema.Message, opts custom_eino.RuntimeOptions, in *ChatInParam, out *ChatOutParam) (*schema.Message, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	// reset messages
	var systemAppend bool
	out.messages = make([]adaptor.ZhimaChatCompletionMessage, 0) // clear
	out.messages = common.BuildOpenApiContent(in.params, out.messages)
	for _, message := range input {
		if message == nil {
			continue
		}
		result, err := custom_eino.ConvertMessage(*message)
		if err != nil {
			logs.Error(`ConvertMessage:` + err.Error())
			continue
		}
		if !systemAppend && message.Role == schema.System {
			systemAppend = true
			result.Content = buildSystemPrompt(result.Content, in, out)
		}
		out.messages = append(out.messages, result)
	}
	// reset functionTools
	out.functionTools = make([]adaptor.FunctionTool, 0) // clear
	for _, info := range opts.Common.Tools {
		if info == nil {
			continue
		}
		result, err := custom_eino.ConvertTools(*info)
		if err != nil {
			logs.Error(`ConvertTools:` + err.Error())
			continue
		}
		out.functionTools = append(out.functionTools, result)
	}
	// RequestChat
	chatResp, _, err := common.RequestChat(
		in.params.Lang,
		in.params.AdminUserId,
		in.params.Openid,
		in.params.Robot,
		in.params.AppType,
		cast.ToInt(in.params.Robot[`model_config_id`]),
		in.params.Robot[`use_model`],
		out.messages,
		out.functionTools,
		cast.ToFloat32(in.params.Robot[`temperature`]),
		cast.ToInt(in.params.Robot[`max_token`]),
	)
	if chatResp.PromptToken > 0 || chatResp.CompletionToken > 0 {
		out.chatResp.PromptToken += chatResp.PromptToken
		out.chatResp.CompletionToken += chatResp.CompletionToken
	}
	if err != nil {
		return nil, err
	}
	// AssistantMessage
	return custom_eino.ConvertChatResp(chatResp), nil
}

func Stream(ctx context.Context, input []*schema.Message, opts custom_eino.RuntimeOptions, in *ChatInParam, out *ChatOutParam) (*schema.StreamReader[*schema.Message], error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	// reset messages
	var systemAppend bool
	out.messages = make([]adaptor.ZhimaChatCompletionMessage, 0) // clear
	out.messages = common.BuildOpenApiContent(in.params, out.messages)
	for _, message := range input {
		if message == nil {
			continue
		}
		result, err := custom_eino.ConvertMessage(*message)
		if err != nil {
			logs.Error(`ConvertMessage:` + err.Error())
			continue
		}
		if !systemAppend && message.Role == schema.System {
			systemAppend = true
			result.Content = buildSystemPrompt(result.Content, in, out)
		}
		out.messages = append(out.messages, result)
	}
	// reset functionTools
	out.functionTools = make([]adaptor.FunctionTool, 0) // clear
	for _, info := range opts.Common.Tools {
		if info == nil {
			continue
		}
		result, err := custom_eino.ConvertTools(*info)
		if err != nil {
			logs.Error(`ConvertTools:` + err.Error())
			continue
		}
		out.functionTools = append(out.functionTools, result)
	}
	// RequestChatStream
	var streamErr error
	var totalResponse adaptor.ZhimaChatCompletionResponse
	chanStream := make(chan sse.Event)
	go func() {
		defer close(chanStream)
		totalResponse, _, streamErr = common.RequestChatStream(
			in.params.Lang,
			in.params.AdminUserId,
			in.params.Openid,
			in.params.Robot,
			in.params.AppType,
			cast.ToInt(in.params.Robot[`model_config_id`]),
			in.params.Robot[`use_model`],
			out.messages,
			out.functionTools,
			chanStream,
			cast.ToFloat32(in.params.Robot[`temperature`]),
			cast.ToInt(in.params.Robot[`max_token`]),
		)
		if totalResponse.PromptToken > 0 || totalResponse.CompletionToken > 0 {
			out.chatResp.PromptToken += totalResponse.PromptToken
			out.chatResp.CompletionToken += totalResponse.CompletionToken
		}
	}()
	// StreamReader
	reader, writer := schema.Pipe[*schema.Message](0)
	go func() {
		defer writer.Close()
		in.Stream(sse.Event{Event: `llm_rounds`, Data: `begin`})
		defer in.Stream(sse.Event{Event: `llm_rounds`, Data: `finish`})
		for event := range chanStream {
			if event.Event != `stream_raw` {
				continue
			}
			response, ok := event.Data.(adaptor.ZhimaChatCompletionResponse)
			if !ok {
				continue
			}
			if len(response.FunctionToolCalls) > 0 {
				continue
			}
			_ = writer.Send(custom_eino.ConvertChatResp(response), nil)
		}
		if streamErr != nil {
			_ = writer.Send(nil, streamErr)
		}
		if len(totalResponse.FunctionToolCalls) > 0 {
			_ = writer.Send(custom_eino.ConvertChatResp(totalResponse), nil)
		}
	}()
	return reader, nil
}

func newSkillMiddleware(ctx context.Context, backend *custom_eino.Backend, robot msql.Params) (adk.ChatModelAgentMiddleware, error) {
	privateBackend, err := skill.NewBackendFromFilesystem(ctx, &skill.BackendFromFilesystemConfig{Backend: backend, BaseDir: strings.ReplaceAll(define.PrivateSkillsDir, `<robot_key>`, robot[`robot_key`])})
	if err != nil {
		return nil, err
	}
	publicBackend, err := skill.NewBackendFromFilesystem(ctx, &skill.BackendFromFilesystemConfig{Backend: backend, BaseDir: define.PublicSkillsDir})
	if err != nil {
		return nil, err
	}
	return skill.NewMiddleware(ctx, &skill.Config{
		Backend: &custom_eino.LayeredSkillBackend{
			Backends: []skill.Backend{privateBackend, publicBackend},
			Filter: func(name string) bool {
				if name == `query-local-docs` && cast.ToBool(robot[`query_local_docs_close`]) {
					return true
				}
				return false
			},
		},
	})
}

func newKbsearchTool(params *define.ChatRequestParam, sessionId int) einotool.BaseTool {
	if cast.ToBool(params.Robot[`search_knowledge_close`]) {
		return nil // knowledge base search has been disabled
	}
	kbSearchTool := custom_eino.BuildKbsearchTool(func(query string) (string, error) {
		// Replace chat variable placeholders in metadata filter config (if enabled)
		common.ReplaceMetaSearchChatVariables(params.Lang, sessionId, &params.Robot)
		if len(params.LibraryIds) == 0 || !common.CheckIds(params.LibraryIds) { //no custom is used
			params.LibraryIds = params.Robot[`library_ids`]
		}
		var appId string
		if tool.InArrayString(params.AppType, lib_define.AppTypeList) {
			appId = params.AppInfo[`app_id`]
		}
		list, _, err := common.GetMatchLibraryParagraphList(
			params.Lang,
			params.Openid,
			params.AppType,
			appId,
			query,
			[]string{},
			params.LibraryIds,
			cast.ToInt(params.Robot[`top_k`]),
			cast.ToFloat64(params.Robot[`similarity`]),
			cast.ToInt(params.Robot[`search_type`]),
			params.Robot,
		)
		if err != nil {
			return ``, err
		}
		_, libraryContent := common.FormatSystemPrompt(params.Lang, ``, list)
		return libraryContent, nil
	})
	return kbSearchTool
}

func ConcatStreamMessages(in *ChatInParam, mv *adk.MessageVariant) (*schema.Message, error) {
	if !mv.IsStreaming {
		return mv.Message, nil
	}
	if mv.MessageStream == nil {
		return mv.Message, nil
	}
	mv.MessageStream.SetAutomaticClose()
	var frames []*schema.Message
	for {
		frame, err := mv.MessageStream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		if frame == nil {
			continue
		}
		in.Stream(sse.Event{Event: `stream_message`, Data: frame})
		frames = append(frames, frame)
	}
	if len(frames) > 0 {
		return schema.ConcatMessages(frames)
	}
	return mv.Message, nil
}

func doClawbotRelationWorkFlow(functionTool adaptor.FunctionToolCall, in *ChatInParam, out *ChatOutParam) (string, error) {
	workFlowRobot, workFlowGlobal := work_flow.ChooseWorkFlowRobot(cast.ToString(in.params.AdminUserId), []adaptor.FunctionToolCall{functionTool})
	if len(workFlowRobot) == 0 {
		return ``, errors.New(`failed to obtain workflow skill information`)
	}
	workFlowParams := work_flow.BuildWorkFlowParams(*in.params, workFlowRobot, workFlowGlobal, int(out.cMsgId), in.dialogueId, in.sessionId)
	workFlowParams.ImmediatelyReplyHandle = BuildImmediatelyReplyHandle(in, out)
	content, _, _, _, replyContentList, err := work_flow.CallWorkFlow(workFlowParams, nil, nil, nil)
	if len(replyContentList) > 0 {
		out.replyContentList = append(out.replyContentList, replyContentList...)
	}
	return content, err
}
