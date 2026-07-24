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

	"github.com/alibabacloud-go/tea/tea"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/filesystem"
	"github.com/cloudwego/eino/adk/middlewares/skill"
	"github.com/cloudwego/eino/adk/prebuilt/deep"
	einotool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/eino-contrib/jsonschema"
	"github.com/gin-contrib/sse"
	"github.com/matiasinsaurralde/go-e2b"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

func buildFileOperationPrompt(robot msql.Params, e2B bool) string {
	queryLocalDocsClose := cast.ToBool(robot[`query_local_docs_close`])
	writeFileEnabled := cast.ToBool(robot[`open_agent_write_file_tool`])
	executeEnabled := cast.ToBool(robot[`open_agent_execute_tool`])
	editFileEnabled := cast.ToBool(robot[`open_agent_edit_file_tool`])

	publicSkillsDir := define.PublicSkillsDir
	privateSkillsDir := strings.ReplaceAll(define.PrivateSkillsDir, `<robot_key>`, robot[`robot_key`])
	privateFileDir := strings.ReplaceAll(define.PrivateFileDir, `<robot_key>`, robot[`robot_key`])
	privateWorkDir := strings.ReplaceAll(define.PrivateWorkDir, `<robot_key>`, robot[`robot_key`])

	enabledTools := []string{`ls`, `read_file`, `grep`, `glob`}
	if writeFileEnabled {
		enabledTools = append(enabledTools, `write_file`)
	}
	if editFileEnabled {
		enabledTools = append(enabledTools, `edit_file`)
	}
	if executeEnabled {
		enabledTools = append(enabledTools, `execute`)
	}

	workDirPermissions := []string{`read`}
	if writeFileEnabled || editFileEnabled {
		workDirPermissions = append(workDirPermissions, `write`)
	}
	if executeEnabled && !e2B {
		workDirPermissions = append(workDirPermissions, `command execution`)
	}

	var prompt strings.Builder
	prompt.WriteString("## File and Command Operation Notes\n")
	prompt.WriteString("- Follow these rules over any conflicting default tool description. File tools only accept paths under the allowed `clawbot/...` directories listed below.\n")
	prompt.WriteString("- Available file/command tools for this session: `" + strings.Join(enabledTools, "`, `") + "`. Do not call file/command tools that are not listed here.\n")
	prompt.WriteString("- Read-only tools: `ls`, `read_file`, `grep`, and `glob`; use them to inspect allowed directories, read files, search file content, and match file paths.\n")
	if writeFileEnabled {
		prompt.WriteString("- Write tool: `write_file` may only create or overwrite files under the workspace directory. Before overwriting an existing file, read the original file with `read_file` first.\n")
	}
	if editFileEnabled {
		prompt.WriteString("- Edit tool: `edit_file` may only edit files under the workspace directory. Before editing, read the target file with `read_file`; `old_string` must be unique and preserve the original indentation.\n")
	}
	if executeEnabled && !e2B {
		prompt.WriteString("- Command execution tool: `execute` passes the full command string to `/bin/sh -c` inside the llm_runner container. Shell syntax is supported, but prefer one simple command; use chaining, pipes, redirects, variables, or command substitution only when the task requires them, and quote paths and user-provided values safely.\n")
		prompt.WriteString("- Basic command validation rejects these arguments when used with the leading `find` command: `-delete`, `-exec`, and `-execdir`. Do not use or attempt to bypass these blocked forms. Python and Node inline or module execution arguments are allowed.\n")
		prompt.WriteString("- Commands run from `/workspace`. Keep command file access within the accessible directories listed below: skill directories are read-only and the workspace directory is writable. Prefer relative `clawbot/...` paths so commands are portable.\n")
		prompt.WriteString("- Keep commands concise and deterministic. Avoid destructive or unrelated operations unless the user explicitly asks for them.\n")
		prompt.WriteString("- Prefer `read_file`, `grep`, and `glob` for normal file inspection. Use `execute` only when the task requires a command result.\n")
		_, _ = fmt.Fprintf(&prompt, "- Valid `execute` examples: `python3 -c \"print('hello')\"`, `node --eval \"console.log('hello')\"`, `python3 %s/test_python.py`, `ls %s`, `rg \"keyword\" %s`.\n", privateWorkDir, privateWorkDir, privateWorkDir)
		prompt.WriteString("- Blocked `execute` example: `find . -delete`.\n")
	}
	if executeEnabled && e2B {
		prompt.WriteString("- Command execution tool: `execute` runs inside the configured E2B sandbox. In E2B mode, commands bypass the local `ValidateCommand` denied-argument checks.\n")
		prompt.WriteString("- E2B command execution always passes the full command string to `bash -lc <command>`. Bash syntax such as `cd`, `&&`, `|`, `;`, environment variables, redirects, command substitution, and multi-line commands can be used when the task requires them.\n")
		prompt.WriteString("- File tools operate on the local Clawbot directories, while `execute` runs in a separate E2B filesystem. Do not assume files created or edited with file tools exist inside E2B; create E2B-side files within `execute` commands when needed.\n")
		prompt.WriteString("- Prefer `read_file`, `grep`, and `glob` for local file inspection. Use `execute` only when the task requires running a shell command, script, package manager, or runtime inside the E2B sandbox.\n")
		prompt.WriteString("- Keep E2B commands concise and deterministic; quote paths and user-provided values safely, and avoid destructive commands unless the user explicitly asks for them.\n")
		prompt.WriteString("- Valid E2B `execute` examples: `python3 -c \"print('hello')\"`, `printf 'keyword\\n' > /tmp/input.txt && grep \"keyword\" /tmp/input.txt`, `npm --version`.\n")
	}
	prompt.WriteString("- Accessible directories and permissions:\n")
	_, _ = fmt.Fprintf(&prompt, "  - Public skills directory (read-only): `%s`\n", publicSkillsDir)
	_, _ = fmt.Fprintf(&prompt, "  - Private skills directory (read-only): `%s`\n", privateSkillsDir)
	if !queryLocalDocsClose {
		_, _ = fmt.Fprintf(&prompt, "  - Local document reference directory (read-only): `%s`\n", privateFileDir)
	}
	_, _ = fmt.Fprintf(&prompt, "  - Workspace directory (%s): `%s`\n", strings.Join(workDirPermissions, ", "), privateWorkDir)
	if executeEnabled && e2B {
		prompt.WriteString("- File tool path arguments must be relative paths that start with one of the allowed directory prefixes above. For `execute`, use paths that exist inside the E2B sandbox, or create/download the needed files there first.\n")
	} else {
		prompt.WriteString("- File tool path arguments must start with one of the allowed directory prefixes above. Absolute paths are accepted only when they resolve back under an allowed `clawbot/...` directory; paths that escape the allowed directories or contain parent-directory traversal (`..`) are rejected.\n")
	}
	prompt.WriteString("- Always use the smallest sufficient operation. Do not write files, edit files, or execute commands unless the user explicitly requires it.\n")
	prompt.WriteString("- Valid path examples:\n")
	_, _ = fmt.Fprintf(&prompt, "  - `%s/skill_name/reference/filename.ext`\n", publicSkillsDir)
	_, _ = fmt.Fprintf(&prompt, "  - `%s/skill_name/reference/filename.ext`\n", privateSkillsDir)
	_, _ = fmt.Fprintf(&prompt, "  - `%s/filename.ext`", privateWorkDir)
	return prompt.String()
}

func ValidateFile(robot msql.Params, op custom_eino.FileOperation, filePath string) error {
	filePath = strings.ReplaceAll(strings.TrimSpace(filePath), `\`, `/`)
	if filePath == `` {
		return fmt.Errorf(`file path is required`)
	}
	// Normalize absolute workspace paths (e.g. /home/ubuntu/clawbot/... discovered
	// via ls/pwd) back to the relative clawbot/... form so the scope check below
	// applies. The backend still opens the original (absolute) path, which resolves
	// correctly; this normalization is for validation only. Absolute paths without
	// a clawbot/ anchor remain rejected.
	if strings.HasPrefix(filePath, `/`) {
		if rel := llm_runner.StripToClawbotAnchor(filePath, "clawbot/"); rel != `` {
			filePath = rel
		} else {
			return fmt.Errorf(`absolute path is not allowed: %s`, filePath)
		}
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
	case custom_eino.FileOperationWrite:
		if !cast.ToBool(robot[`open_agent_write_file_tool`]) {
			return fmt.Errorf(`write_file tool is disabled`)
		}
		allowedPrefixes = []string{
			strings.ReplaceAll(define.PrivateWorkDir, `<robot_key>`, robot[`robot_key`]),
		}
	case custom_eino.FileOperationEdit:
		if !cast.ToBool(robot[`open_agent_edit_file_tool`]) {
			return fmt.Errorf(`edit_file tool is disabled`)
		}
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

func doApplicationTypeClaw(in *ChatInParam, out *ChatOutParam) error {
	// set language
	_ = adk.SetLanguage(adk.LanguageChinese)
	// init dirs
	common.InitClawbotDirs(in.params.Robot[`robot_key`])
	// ChatModel
	ctx := clawbotRunContext(in)
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
	// E2B sandbox
	in.e2bShell, err = newE2BShellFromConf(ctx, in.params.Robot)
	if err != nil {
		return err
	}
	if in.e2bShell != nil {
		defer func() {
			if err = in.e2bShell.Close(); err != nil {
				logs.Error(`close E2B sandbox failed: %v`, err)
			}
			in.e2bShell = nil // cleanup
		}()
	}
	// Backend
	backend, err := custom_eino.NewBackend(ctx, &custom_eino.BackendConfig{
		ValidateFile: func(op custom_eino.FileOperation, path string) error {
			in.Stream(sse.Event{Event: `FileOperation`, Data: map[string]any{`op`: op, `path`: path}})
			return ValidateFile(in.params.Robot, op, path)
		},
		ValidateCommand: func(command string) error {
			in.Stream(sse.Event{Event: `ExecuteCommand`, Data: map[string]any{`command`: command}})
			if !cast.ToBool(in.params.Robot[`open_agent_execute_tool`]) {
				return fmt.Errorf(`execute tool is disabled`)
			}
			if in.e2bShell != nil {
				return nil // E2B sandbox not check command
			}
			return llm_runner.ValidateCommand(command)
		},
		ExecuteCommand: func(command string) (*filesystem.ExecuteResponse, error) {
			if in.e2bShell != nil {
				return in.e2bShell.Execute(command) // E2B sandbox
			}
			resp := llm_runner.RpcExecuteRun(define.Config.WebService[`llm_runner_host`], ``, command)
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
	goodsLibRecommendTool := newGoodsLibRecommendTool(in.params)
	if goodsLibRecommendTool != nil {
		tools = append(tools, goodsLibRecommendTool)
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
		ChatModel:   cm,
		Instruction: "\n", // clear the default prompt words of the eino framework
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
		in.dialogueId, in.sessionId, int(out.cMsgId), cast.ToInt(in.params.Robot[`context_pair`]))
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
			if ctx.Err() != nil {
				break // client disconnect/stop: keep partial content, skip fallback
			}
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
			if ctx.Err() != nil {
				break // client disconnect/stop: keep partial content, skip fallback
			}
			out.Error = err
			common.SendDefaultUnknownQuestionPrompt(in.params, out.Error.Error(), in.chanStream, &out.content)
			return nil
		}
		role := event.Output.MessageOutput.Role
		if role == schema.Assistant || role == `` {
			for _, call := range message.ToolCalls {
				in.Stream(sse.Event{Event: `tool_call`, Data: call.Function}) // Deprecated: Get the full content from tool_call_full
				in.Stream(sse.Event{Event: `tool_call_full`, Data: call})
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
	system = `` // TODO: clear the default prompt words of the eino framework
	result := fmt.Sprintf("%s\n%s\n%s", in.params.Prompt, system, buildFileOperationPrompt(in.params.Robot, in.e2bShell != nil))
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
	filterTools := buildFilterTools(in.params.Robot)
	out.functionTools = make([]adaptor.FunctionTool, 0) // clear
	for _, info := range opts.Common.Tools {
		if info == nil {
			continue
		}
		if tool.InArray(info.Name, filterTools) {
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

func buildFilterTools(robot msql.Params) []string {
	filterTools := []string{`write_todos`, `todo_write`}
	if !cast.ToBool(robot[`open_agent_write_file_tool`]) {
		filterTools = append(filterTools, `write_file`)
	}
	if !cast.ToBool(robot[`open_agent_execute_tool`]) {
		filterTools = append(filterTools, `execute`)
	}
	if !cast.ToBool(robot[`open_agent_edit_file_tool`]) {
		filterTools = append(filterTools, `edit_file`)
	}
	return filterTools
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
	filterTools := buildFilterTools(in.params.Robot)
	out.functionTools = make([]adaptor.FunctionTool, 0) // clear
	for _, info := range opts.Common.Tools {
		if info == nil {
			continue
		}
		if tool.InArray(info.Name, filterTools) {
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
			ctx,
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
			// Fallback for legacy providers that stream a single tool call without index.
			// Multi-tool streaming without index cannot be merged reliably here.
			if len(response.ToolCalls) == 1 && response.ToolCalls[0].Index == nil {
				response.ToolCalls[0].Index = tea.Int(0)
			}
			_ = writer.Send(custom_eino.ConvertChatResp(response), nil)
		}
		if streamErr != nil {
			_ = writer.Send(nil, streamErr)
		}
	}()
	return reader, nil
}

func newE2BShellFromConf(ctx context.Context, robot msql.Params) (*custom_eino.E2BShell, error) {
	if !cast.ToBool(robot[`open_agent_execute_tool`]) {
		return nil, nil
	}
	info, err := common.GetE2bConfInfo(robot[`robot_key`])
	if err != nil {
		return nil, err
	}
	conf := common.BuildE2BConfParams(info)
	if !cast.ToBool(conf.SwitchStatus) {
		return nil, nil
	}
	return custom_eino.NewE2BShell(ctx,
		e2b.ClientConfig{APIKey: conf.ApiKey, APIBaseURL: conf.ApiBaseUrl, SandboxDomain: conf.SandboxDomain},
		e2b.SandboxConfig{Template: conf.Template, Timeout: conf.Timeout, Secure: true},
		e2b.WithTimeout(time.Duration(conf.CommandTimeout)*time.Second),
		e2b.WithUser(conf.CommandUser),
	)
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

func newGoodsLibRecommendTool(params *define.ChatRequestParam) einotool.BaseTool {
	if !cast.ToBool(params.Robot[`goods_lib_recommend_switch`]) {
		return nil // the recommendation of goods from the goods library has not been enabled
	}
	return custom_eino.BuildGoodsLibRecommendTool(func(query, searchType string, maxCount int) (string, error) {
		filter := define.GoodsLibListFilter{
			GroupID:      -1,
			GroupIDs:     params.Robot[`goods_lib_recommend_group_ids`],
			Keyword:      query,
			SwitchStatus: define.GoodsLibSwitchOn,
			Page:         1,
			Size:         maxCount,
		}
		list, total, err := common.GetGoodsLibLibraryList(params.Lang, params.AdminUserId, filter)
		if err != nil {
			return ``, err
		}
		return common.FormatGoodsLibRecommendResult(searchType, list, total), nil
	})
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
