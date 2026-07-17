// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/custom_eino"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/llm_runner"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/filesystem"
	"github.com/cloudwego/eino/adk/middlewares/skill"
	"github.com/cloudwego/eino/adk/prebuilt/deep"
	"github.com/cloudwego/eino/schema"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

type WebToSkillTaskInfo struct {
	TaskBatch     string
	AdminUserId   int
	ModelConfigId int
	UseModel      string
	Temperature   float32
	MaxToken      int
	Urls          []string
	CustomPrompt  string
	StopKey       string
}

type WebToSkillResult struct {
	DebugLog []any  `json:"debug_log"`
	ZipPath  string `json:"zip_path"`
}

func (task WebToSkillTaskInfo) StopRequested() bool {
	return IsWebToSkillTaskStopped(task.StopKey)
}

func DoWebToSkill(lang string, task WebToSkillTaskInfo) (WebToSkillResult, error) {
	// set language
	_ = adk.SetLanguage(adk.LanguageChinese)
	// init dirs
	if task.StopRequested() {
		return WebToSkillResult{}, errors.New(i18n.Show(lang, `web_to_skill_task_stopped`))
	}
	workDir := strings.ReplaceAll(define.WebToSkillWorkDir, `<task_batch>`, task.TaskBatch)
	_ = tool.MkDirAll(workDir)
	// ChatModel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cm, err := custom_eino.NewChatModel(ctx, &custom_eino.ChatModelConfig{
		Generate: func(ctx context.Context, input []*schema.Message, opts custom_eino.RuntimeOptions) (*schema.Message, error) {
			return webToSkillGenerate(ctx, input, opts, lang, task)
		},
	})
	if err != nil {
		return WebToSkillResult{}, errors.New(i18n.Show(lang, `web_to_skill_runner_error`, err.Error()))
	}
	// Backend
	backend, err := custom_eino.NewBackend(ctx, &custom_eino.BackendConfig{
		ValidateFile: nil, // do not verify file operations
		ValidateCommand: func(command string) error {
			return llm_runner.ValidateCommand(StripCdPrefix(command))
		},
		ExecuteCommand: func(command string) (*filesystem.ExecuteResponse, error) {
			if task.StopRequested() {
				return nil, errors.New(i18n.Show(lang, `web_to_skill_task_stopped`))
			}
			resp := llm_runner.RpcExecuteRunWithID(
				define.Config.WebService[`llm_runner_host`],
				``,
				GetWebToSkillTaskRunID(task.TaskBatch),
				command,
			)
			if task.StopRequested() {
				cancel()
				return nil, errors.New(i18n.Show(lang, `web_to_skill_task_stopped`))
			}
			if resp.IsError {
				if resp.ExitCode < 0 { // server execution error
					return &filesystem.ExecuteResponse{Output: resp.ErrorMsg, ExitCode: &resp.ExitCode}, nil
				}
				return nil, errors.New(i18n.Show(lang, `web_to_skill_runner_error`, resp.ErrorMsg))
			}
			return &filesystem.ExecuteResponse{Output: resp.Output, ExitCode: &resp.ExitCode}, nil
		},
	})
	if err != nil {
		return WebToSkillResult{}, errors.New(i18n.Show(lang, `web_to_skill_runner_error`, err.Error()))
	}
	// Skill
	skillBackend, err := skill.NewBackendFromFilesystem(ctx, &skill.BackendFromFilesystemConfig{
		Backend: backend, BaseDir: strings.ReplaceAll(define.SystemSkillsDir, `<biz_type>`, `alone_web`)})
	if err != nil {
		return WebToSkillResult{}, errors.New(i18n.Show(lang, `web_to_skill_runner_error`, err.Error()))
	}
	skillMiddleware, err := skill.NewMiddleware(ctx, &skill.Config{Backend: skillBackend})
	if err != nil {
		return WebToSkillResult{}, errors.New(i18n.Show(lang, `web_to_skill_runner_error`, err.Error()))
	}
	// Runner
	agent, err := deep.New(ctx, &deep.Config{
		ChatModel:              cm,
		Instruction:            "\n", // clear the default prompt words of the eino framework
		MaxIteration:           50,
		Backend:                backend,
		Shell:                  backend,
		WithoutGeneralSubAgent: true,
		Handlers:               []adk.ChatModelAgentMiddleware{skillMiddleware},
		ModelRetryConfig: &adk.ModelRetryConfig{
			MaxRetries: 3,
			ShouldRetry: func(_ context.Context, retryCtx *adk.RetryContext) *adk.RetryDecision {
				return &adk.RetryDecision{Retry: retryCtx.Err != nil, RewriteError: retryCtx.Err}
			},
		},
	})
	if err != nil {
		return WebToSkillResult{}, errors.New(i18n.Show(lang, `web_to_skill_runner_error`, err.Error()))
	}
	runner := adk.NewRunner(ctx, adk.RunnerConfig{Agent: agent, EnableStreaming: false})
	// events
	debugLog := make([]any, 0)
	events := runner.Run(ctx, []*schema.Message{schema.UserMessage(strings.Join(task.Urls, "\n"))})
	var sb strings.Builder
	var toolZipPath string
	stopRequested := false
	requestStop := func() {
		if stopRequested {
			return
		}
		stopRequested = true
		cancel()
		_ = llm_runner.RpcCancelRun(
			define.Config.WebService[`llm_runner_host`],
			GetWebToSkillTaskRunID(task.TaskBatch),
		)
	}
	for {
		if task.StopRequested() {
			requestStop()
		}
		event, ok := events.Next()
		if task.StopRequested() {
			requestStop()
		}
		if !ok {
			break
		}
		if stopRequested {
			continue
		}
		if event.Err != nil {
			return WebToSkillResult{DebugLog: debugLog, ZipPath: strings.TrimSpace(sb.String())}, errors.New(i18n.Show(lang, `web_to_skill_runner_error`, event.Err.Error()))
		}
		if event.Action != nil && event.Action.Interrupted != nil {
			continue
		}
		if event.Output == nil || event.Output.MessageOutput == nil {
			continue
		}
		message := event.Output.MessageOutput.Message
		role := event.Output.MessageOutput.Role
		if role == schema.Assistant || role == `` {
			for _, call := range message.ToolCalls {
				debugLog = append(debugLog, map[string]any{`type`: `tool_call`, `content`: call.Function})
			}
			sb.WriteString(message.Content)
		} else if role == schema.Tool {
			sb.Reset() // clear the result from the last time
			debugLog = append(debugLog, map[string]string{`type`: `tool_result`, `content`: message.Content})
			if candidate := normalizeWebToSkillZipPath(message.Content, task.TaskBatch); candidate != `` {
				toolZipPath = candidate
			}
		}
	}
	if stopRequested {
		return WebToSkillResult{DebugLog: debugLog, ZipPath: strings.TrimSpace(sb.String())}, errors.New(i18n.Show(lang, `web_to_skill_task_stopped`))
	}
	rawZipPath := strings.TrimSpace(sb.String())
	zipPath := normalizeWebToSkillZipPath(rawZipPath, task.TaskBatch)
	debugLog = append(debugLog, map[string]string{`type`: `llm_result`, `content`: rawZipPath})
	if zipPath == `` && toolZipPath != `` {
		zipPath = toolZipPath
		debugLog = append(debugLog, map[string]string{`type`: `zip_path_fallback`, `content`: toolZipPath})
	}
	return WebToSkillResult{DebugLog: debugLog, ZipPath: zipPath}, nil
}

func buildWebToSkillSystemPrompt(task WebToSkillTaskInfo) string {
	workDir := strings.ReplaceAll(define.WebToSkillWorkDir, `<task_batch>`, task.TaskBatch)
	prompt := fmt.Sprintf(`You are the ChatWiki web-to-skill generation agent.

Proactively load and follow the $web-to-skill skill to convert the URL or URL list in the user message into a dedicated reusable skill zip.

The llm_runner environment already includes Python, Playwright, Chromium, beautifulsoup4, lxml, and jieba. Do not install, upgrade, or reinstall Python packages or browser binaries.

Writable task directory:
%[1]s

Keep all generated files, temporary files, and the final zip under this directory.

Final zip path format:
%[1]s/generate_skill/<skill-name>.zip

On success, output only the generated zip path. Do not include explanations, Markdown, or any other text.`, workDir)
	if customPrompt := strings.TrimSpace(task.CustomPrompt); customPrompt != `` {
		prompt += fmt.Sprintf(`

Additional user requirements:
Apply the following requirements when using the skill, unless they conflict with the writable task directory or final output format above.

<custom_prompt>
%s
</custom_prompt>`, customPrompt)
	}
	return prompt
}

func webToSkillGenerate(ctx context.Context, input []*schema.Message, opts custom_eino.RuntimeOptions, lang string, task WebToSkillTaskInfo) (*schema.Message, error) {
	if err := ctx.Err(); err != nil {
		return nil, errors.New(i18n.Show(lang, `web_to_skill_runner_error`, err.Error()))
	}
	if task.StopRequested() {
		return nil, errors.New(i18n.Show(lang, `web_to_skill_task_stopped`))
	}
	// messages
	var systemAppend bool
	messages := make([]adaptor.ZhimaChatCompletionMessage, 0)
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
			result.Content = buildWebToSkillSystemPrompt(task)
		}
		messages = append(messages, result)
	}
	// functionTools
	functionTools := make([]adaptor.FunctionTool, 0)
	filterTools := []string{`write_todos`, `todo_write`}
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
		functionTools = append(functionTools, result)
	}
	// RequestChat
	chatResp, _, err := RequestChat(
		lang,
		task.AdminUserId,
		cast.ToString(task.AdminUserId),
		nil,
		lib_define.AppYunPc,
		task.ModelConfigId,
		task.UseModel,
		messages,
		functionTools,
		task.Temperature,
		task.MaxToken,
	)
	if err != nil {
		return nil, errors.New(i18n.Show(lang, `web_to_skill_runner_error`, err.Error()))
	}
	// AssistantMessage
	return custom_eino.ConvertChatResp(chatResp), nil
}
