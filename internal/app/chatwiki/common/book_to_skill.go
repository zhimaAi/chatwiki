// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"archive/zip"
	"chatwiki/internal/app/chatwiki/define"
	"chatwiki/internal/app/chatwiki/i18n"
	"chatwiki/internal/pkg/custom_eino"
	"chatwiki/internal/pkg/lib_define"
	"chatwiki/internal/pkg/llm_runner"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/filesystem"
	"github.com/cloudwego/eino/adk/middlewares/skill"
	"github.com/cloudwego/eino/adk/prebuilt/deep"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/gin-contrib/sse"
	"github.com/spf13/cast"
	"github.com/zhimaAi/go_tools/logs"
	"github.com/zhimaAi/go_tools/msql"
	"github.com/zhimaAi/go_tools/tool"
	"github.com/zhimaAi/llm_adaptor/adaptor"
)

// FileInfo represents an uploaded book file
type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Ext  string `json:"ext"`
	Link string `json:"link"`
}

// BookToSkillTask represents a book-to-skill generation task
type BookToSkillTask struct {
	Id               int        `json:"id"`
	AdminUserId      int        `json:"admin_user_id"`
	RobotId          int        `json:"robot_id"`
	RobotKey         string     `json:"robot_key"`
	SkillName        string     `json:"skill_name"`
	ModelConfigId    int        `json:"model_config_id"`
	UseModel         string     `json:"use_model"`
	Status           int        `json:"status"`
	SourceFiles      []FileInfo `json:"source_files"`
	TemplateContent  string     `json:"template_content"`
	SkillDir         string     `json:"skill_dir"`
	SkillFilePath    string     `json:"skill_file_path"`
	ErrorMsg         string     `json:"error_msg"`
	InstallStatus    int        `json:"install_status"`
	CurrentIteration int        `json:"current_iteration"`
	Lang             string     `json:"lang"`
	CreateTime       int64      `json:"create_time"`
	UpdateTime       int64      `json:"update_time"`
}

// BookToSkillExecutor executes the book-to-skill conversion using single-agent five-step (A→E) orchestration
type BookToSkillExecutor struct {
	Task           *BookToSkillTask
	CancelCtx      context.Context
	Cancel         context.CancelFunc
	iterationCount atomic.Int32 // cross-phase cumulative iteration count
	currentStep    int          // current phase hint (1=split, 2=batch, 0=merge/generate)
	tokenThreshold int          // compression threshold
	runLog         *os.File     // run.log file handle (task working dir)
	// shared resources (created once, reused across phases)
	cm      *custom_eino.ChatModel
	backend *custom_eino.Backend
	skillMW adk.ChatModelAgentMiddleware
}

// writeLog writes a line to run.log (for file).
func (e *BookToSkillExecutor) writeLog(line string) {
	if e.runLog != nil {
		_, _ = e.runLog.WriteString(line)
	}
}

// writeLogf is the formatted version of writeLog.
func (e *BookToSkillExecutor) writeLogf(format string, args ...any) {
	e.writeLog(fmt.Sprintf(format, args...))
}

// ========== Working directory helpers ==========

// bookToSkillWorkDir computes the task working directory (replace <robot_key> + <task_id>)
func bookToSkillWorkDir(robotKey string, taskId int) string {
	dir := strings.ReplaceAll(define.BookToSkillTaskDir, "<robot_key>", robotKey)
	return strings.ReplaceAll(dir, "<task_id>", cast.ToString(taskId))
}

// bookToSkillSubDir computes a task sub-directory from a dir constant
func bookToSkillSubDir(robotKey string, taskId int, dirConst string) string {
	dir := strings.ReplaceAll(dirConst, "<robot_key>", robotKey)
	return strings.ReplaceAll(dir, "<task_id>", cast.ToString(taskId))
}

func stopFlagPath(robotKey string, taskId int) string {
	return bookToSkillWorkDir(robotKey, taskId) + "/.stop"
}

// checkStopFlag returns true if .stop file exists
func checkStopFlag(robotKey string, taskId int) bool {
	_, err := os.Stat(stopFlagPath(robotKey, taskId))
	return err == nil
}

// CreateStopFlag creates the .stop signal file
func CreateStopFlag(robotKey string, taskId int) error {
	return os.WriteFile(stopFlagPath(robotKey, taskId), []byte("1"), 0644)
}

// ClearStopFlag removes the .stop signal file
func ClearStopFlag(robotKey string, taskId int) {
	_ = os.Remove(stopFlagPath(robotKey, taskId))
}

// ========== Execute: single-agent five-step orchestration ==========

// ErrTaskStopped is returned when .stop flag is detected
var ErrTaskStopped = errors.New("task stopped by user")

// Execute runs the book-to-skill conversion with the single-agent 5-step flow.
// It calls runSingleAgent (phases: split→batch→merge→generate) then collectSkill.
func (e *BookToSkillExecutor) Execute() error {
	if err := e.initResources(); err != nil {
		return err
	}

	rk := e.Task.RobotKey
	tid := e.Task.Id
	workDir := bookToSkillWorkDir(rk, tid)

	// Open run.log for detailed execution logging
	logPath := workDir + "/run.log"
	if f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644); err == nil {
		e.runLog = f
		defer func() { _ = f.Close() }()
	}

	// Clean up legacy architecture artifacts: if old process.txt exists, remove old artifacts to avoid interference with new pipeline
	oldProcessPath := workDir + "/process.txt"
	if _, err := os.Stat(oldProcessPath); err == nil {
		e.writeLogf("[Migration] Detected legacy process.txt, cleaning up\n")
		_ = os.Remove(oldProcessPath)
		_ = os.RemoveAll(workDir + "/chunks")
		_ = os.RemoveAll(workDir + "/chunks_raw")
		_ = os.RemoveAll(workDir + "/extracted")
		_ = os.RemoveAll(workDir + "/strategy")
	}

	e.writeLogf("[Start] Single-agent five-step execution, file_count=%d\n", len(e.Task.SourceFiles))

	// Execute single-agent five-step pipeline (A→E, batch map-reduce)
	if err := e.runSingleAgent(workDir); err != nil {
		if errors.Is(err, ErrTaskStopped) {
			e.writeLog("\n[Pause] Task stopped\n")
			return nil // graceful stop
		}
		return err
	}

	// Collect skill into skills_user
	if err := e.collectSkill(); err != nil {
		return fmt.Errorf("collect_skill failed: %w", err)
	}

	e.logWorkDirTree(workDir)
	e.writeLogf("[End] Execution complete, total_iterations=%d\n", e.iterationCount.Load())
	return nil
}

// logWorkDirTree traverses the task working directory and outputs directory/file structure log (no AI calls)
func (e *BookToSkillExecutor) logWorkDirTree(workDir string) {
	e.writeLogf("[DirTree] %s\n", workDir)
	_ = filepath.WalkDir(workDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		rel, _ := filepath.Rel(workDir, path)
		if rel == "." {
			return nil
		}
		if d.IsDir() {
			e.writeLogf("  [Dir] %s\n", rel)
		} else {
			info, _ := d.Info()
			e.writeLogf("  [File] %s (%d bytes)\n", rel, info.Size())
		}
		return nil
	})
}

// initResources creates shared ChatModel + Backend + skillMiddleware (reused across phases)
func (e *BookToSkillExecutor) initResources() error {
	ctx := e.CancelCtx
	task := e.Task
	workDir := bookToSkillWorkDir(task.RobotKey, task.Id)
	e.tokenThreshold = define.BookToSkillCompressThreshold

	// 1. ChatModel (generate callback shared across phases, e.currentStep distinguishes logging)
	cm, err := custom_eino.NewChatModel(ctx, &custom_eino.ChatModelConfig{
		Generate: func(ctx context.Context, input []*schema.Message, opts custom_eino.RuntimeOptions) (*schema.Message, error) {
			return e.generate(ctx, input, opts)
		},
		Stream: func(ctx context.Context, input []*schema.Message, opts custom_eino.RuntimeOptions) (*schema.StreamReader[*schema.Message], error) {
			out, err := e.generate(ctx, input, opts)
			if err != nil {
				return nil, err
			}
			return schema.StreamReaderFromArray([]*schema.Message{out}), nil
		},
	})
	if err != nil {
		return fmt.Errorf("create ChatModel failed: %v", err)
	}
	e.cm = cm

	// 2. Backend (validateFile needs task_id, workDir for relative path resolution)
	backend, err := custom_eino.NewBackend(ctx, &custom_eino.BackendConfig{
		WorkDir:        workDir,
		ErrorAsContent: true,
		ValidateFile:   nil, // do not verify file operations
		ValidateCommand: func(command string) error {
			return llm_runner.ValidateCommand(StripCdPrefix(command))
		},
		ExecuteCommand: func(command string) (*filesystem.ExecuteResponse, error) {
			// Strip cd prefix — workDir is passed to RPC so commands
			// execute in the correct directory without manual cd.
			command = StripCdPrefix(command)

			e.writeLogf("[Exec] workDir=%s cmd=%s\n", workDir, truncateString(command, 100))
			resp := llm_runner.RpcExecuteRun(
				define.Config.WebService["llm_runner_host"], workDir, command)
			exitCode := resp.ExitCode
			if resp.IsError && exitCode < 0 {
				e.writeLogf("[ExecError] %s\n", resp.ErrorMsg)
				exitCode = 1
				return &filesystem.ExecuteResponse{
					Output:   fmt.Sprintf("Error: %s\n%s", resp.ErrorMsg, resp.Output),
					ExitCode: &exitCode,
				}, nil
			}
			outputLen := len(resp.Output)
			sanitized := strings.ReplaceAll(truncateString(strings.TrimSpace(resp.Output), 500), "\n", "\\n")
			e.writeLogf("[ExecDone] exit_code=%d output_len=%d content: %s\n", exitCode, outputLen, sanitized)
			return &filesystem.ExecuteResponse{
				Output: resp.Output, ExitCode: &exitCode,
			}, nil
		},
	})
	if err != nil {
		return fmt.Errorf("create Backend failed: %v", err)
	}
	e.backend = backend

	// 3. skillMiddleware (loads SKILL.md from skills_public/book-to-skill/)
	skillMW, err := newBookToSkillMiddleware(ctx, backend)
	if err != nil {
		return fmt.Errorf("load skill middleware failed: %v", err)
	}
	e.skillMW = skillMW

	return nil
}

// ========== Single-agent progress (in-memory phase tracking) ==========

// BookToSkillProgress tracks batch processing state within one run.
// Since resume is not supported, this is purely in-memory and never persisted.
type BookToSkillProgress struct {
	BatchSize    int
	TotalChunks  int
	TotalBatches int
}

func newProgress() *BookToSkillProgress {
	return &BookToSkillProgress{}
}

// listChunkFiles returns sorted chunk basenames (chunk_NNN.md) under chunks/
func listChunkFiles(workDir string) []string {
	matches, _ := filepath.Glob(workDir + "/chunks/chunk_*.md")
	if len(matches) == 0 {
		matches, _ = filepath.Glob(workDir + "/chunks/chunk_*.txt")
	}
	sort.Strings(matches)
	out := make([]string, 0, len(matches))
	for _, m := range matches {
		out = append(out, filepath.Base(m))
	}
	return out
}

// ========== Single-agent execution (5-step A→E with batch map-reduce) ==========

// runPhase runs a single deep-agent phase with a fresh conversation (no history
// accumulation across batches). stepHint: 1=split, 2=batch, 0=merge/generate.
func (e *BookToSkillExecutor) runPhase(phase string, stepHint int, prompt string) error {
	e.currentStep = stepHint
	ctx := e.CancelCtx

	agent, err := deep.New(ctx, &deep.Config{
		ChatModel:   e.cm,
		Instruction: "\n", // clear the default prompt words of the eino framework
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools:               nil,
				ExecuteSequentially: true,
			},
		},
		MaxIteration:           define.BookToSkillDeepMaxIteration,
		Backend:                e.backend,
		Shell:                  e.backend,
		WithoutGeneralSubAgent: true,
		Handlers:               []adk.ChatModelAgentMiddleware{e.skillMW},
		ModelRetryConfig: &adk.ModelRetryConfig{
			MaxRetries: 3,
			ShouldRetry: func(_ context.Context, retryCtx *adk.RetryContext) *adk.RetryDecision {
				if retryCtx.Err == nil {
					return &adk.RetryDecision{Retry: false}
				}
				errStr := retryCtx.Err.Error()
				if strings.Contains(errStr, "insufficient_quota") ||
					strings.Contains(errStr, "status code: 4") {
					return &adk.RetryDecision{Retry: false, RewriteError: retryCtx.Err}
				}
				return &adk.RetryDecision{Retry: true, RewriteError: retryCtx.Err}
			},
			BackoffFunc: func(_ context.Context, attempt int) time.Duration {
				d := 2 * time.Second * time.Duration(1<<(attempt-1))
				if d > 30*time.Second {
					d = 30 * time.Second
				}
				return d
			},
		},
	})
	if err != nil {
		return fmt.Errorf("create deep agent (%s) failed: %v", phase, err)
	}

	runner := adk.NewRunner(ctx, adk.RunnerConfig{Agent: agent, EnableStreaming: false})
	events := runner.Run(ctx, []*schema.Message{schema.UserMessage(prompt)})
	for {
		event, ok := events.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			e.writeLogf("[Error][%s] %v\n", phase, event.Err)
			return event.Err
		}
		e.logAgentEvent(event)
	}
	return nil
}

// runSingleAgent executes the single-agent 5-step flow (A→E) with batch map-reduce.
// Each phase/batch runs in a fresh deep agent; progress is tracked in-memory only.
func (e *BookToSkillExecutor) runSingleAgent(workDir string) error {
	rk := e.Task.RobotKey
	tid := e.Task.Id
	prog := newProgress()
	prog.BatchSize = define.BookToSkillBatchSize

	lang := e.Task.Lang
	if lang == "" {
		lang = define.BookToSkillDefaultLang
	}

	// Phase 1: split (steps A+B+C)
	if checkStopFlag(rk, tid) {
		return ErrTaskStopped
	}
	e.writeLog("[Phase1] Split (steps A/B/C)\n")
	if err := e.runPhase("split", 1, e.buildSinglePrompt("split", workDir, "", 0, 0, nil)); err != nil {
		return err
	}
	chunks := listChunkFiles(workDir)
	if len(chunks) == 0 {
		return fmt.Errorf(i18n.Show(lang, "bts_split_no_chunks"))
	}
	if len(chunks) > define.BookToSkillSplitMaxChunks {
		msg := fmt.Sprintf("[Warning] Split produced %d chunks exceeding recommended max %d; granularity may be too fine; re-split if batch phase is slow",
			len(chunks), define.BookToSkillSplitMaxChunks)
		e.writeLogf("%s\n", msg)
	}
	prog.TotalChunks = len(chunks)
	prog.TotalBatches = (len(chunks) + prog.BatchSize - 1) / prog.BatchSize
	e.writeLogf("[SplitDone] chunk_count=%d batch_count=%d\n", prog.TotalChunks, prog.TotalBatches)

	// Phase 2: batch extract (step D, map-reduce)
	chunks = listChunkFiles(workDir)
	if len(chunks) == 0 {
		return fmt.Errorf(i18n.Show(lang, "bts_chunks_dir_empty"))
	}
	prog.TotalChunks = len(chunks)
	prog.TotalBatches = (len(chunks) + prog.BatchSize - 1) / prog.BatchSize
	for batch := 1; batch <= prog.TotalBatches; batch++ {
		if checkStopFlag(rk, tid) {
			return ErrTaskStopped
		}
		start := (batch - 1) * prog.BatchSize
		end := start + prog.BatchSize
		if end > len(chunks) {
			end = len(chunks)
		}
		batchChunks := chunks[start:end]
		if len(batchChunks) == 0 {
			break
		}
		e.writeLogf("[Phase2] Batch %d/%d chunk=%s~%s\n", batch, prog.TotalBatches, batchChunks[0], batchChunks[len(batchChunks)-1])
		if err := e.runPhase(fmt.Sprintf("batch_%d", batch), 2, e.buildSinglePrompt("batch", workDir, "", batch, prog.TotalBatches, batchChunks)); err != nil {
			return err
		}
		// Validate batch output: summaries/ should contain files (AI may skip empty chunks, not all chunks require summaries)
		summariesDir := workDir + "/summaries"
		if entries, _ := os.ReadDir(summariesDir); len(entries) == 0 {
			return fmt.Errorf(i18n.Show(lang, "bts_batch_summaries_empty"), batch)
		}
	}

	// Phase 3: merge (step D')
	if checkStopFlag(rk, tid) {
		return ErrTaskStopped
	}
	// Clean up any stale summary_index.txt (files written by AI via write_file during batch phase
	// have owner different from container uid 10001, merge_summaries.py overwrite triggers PermissionError)
	_ = os.Remove(workDir + "/summary_index.txt")
	_ = os.Chmod(workDir, 0777)
	e.writeLog("[Phase3] Merge summary_index\n")
	if err := e.runPhase("merge", 0, e.buildSinglePrompt("merge", workDir, "", 0, prog.TotalBatches, nil)); err != nil {
		return err
	}
	// Validate merge output
	if _, statErr := os.Stat(workDir + "/summary_index.txt"); statErr != nil {
		return fmt.Errorf("%s: %w", i18n.Show(lang, "bts_merge_no_index"), statErr)
	}

	// Phase 4: generate SKILL.md (step E)
	if checkStopFlag(rk, tid) {
		return ErrTaskStopped
	}
	e.writeLog("[Phase4] Generate SKILL.md\n")
	if err := e.runPhase("generate", 0, e.buildSinglePrompt("generate", workDir, e.Task.SkillName, 0, 0, nil)); err != nil {
		return err
	}
	// Validate SKILL.md output
	skillMdPath := fmt.Sprintf("%s/skill/%s/SKILL.md", workDir, e.Task.SkillName)
	if _, statErr := os.Stat(skillMdPath); statErr != nil {
		return fmt.Errorf("%s: %w", i18n.Show(lang, "bts_generate_no_skillmd", skillMdPath), statErr)
	}

	return nil
}

// buildSinglePrompt constructs the phase-specific prompt for the single agent.
// SKILL.md now carries all workflow instructions, command/path conventions.
// Go side only injects runtime context: workDir, phase, batch info, chunk list.
// Prompt language is determined by task.Lang field (zh-CN or en-US).
func (e *BookToSkillExecutor) buildSinglePrompt(phase, workDir, skillName string, batchNo, totalBatches int, batchChunks []string) string {
	task := e.Task
	wd := workDir
	sn := task.SkillName
	if skillName != "" {
		sn = skillName
	}
	lang := task.Lang
	if lang == "" {
		lang = define.BookToSkillDefaultLang
	}

	templateContent := task.TemplateContent
	if templateContent == "" {
		templateContent = "Please read templates/md_extraction.md to get template rules"
	}

	header := fmt.Sprintf(define.GetBtsPrompt(lang, define.BtsPromptHeader),
		define.BookToSkillSplitMaxBytes/1024,
		define.BookToSkillSplitMaxDepth)

	body := define.GetBtsPhasePrompt(lang, phase)

	var prompt string
	switch phase {
	case "split":
		prompt = fmt.Sprintf(body,
			wd,
			define.BookToSkillSplitMaxBytes/1024, define.BookToSkillSplitMaxDepth,
			define.BookToSkillSplitMaxChunks, define.BookToSkillSplitMaxBytes/1024)
	case "batch":
		prompt = fmt.Sprintf(body,
			batchNo, totalBatches,
			wd, sn,
			len(batchChunks), strings.Join(batchChunks, ", "))
	case "merge":
		prompt = fmt.Sprintf(body, wd)
	case "generate":
		prompt = fmt.Sprintf(body, wd, sn)
	default:
		prompt = body
	}

	separator := "【Injected Template】"
	result := fmt.Sprintf("%s\n%s\n%s\n%s\n%s",
		header, separator, templateContent, prompt, define.GetBtsPrompt(lang, define.BtsPromptExecutePathRules))
	e.writeLogf("[Prompt] Phase=%s Length=%d\n", phase, len(result))
	return result
}

// ========== collectSkill ==========

func (e *BookToSkillExecutor) collectSkill() error {
	workDir := bookToSkillWorkDir(e.Task.RobotKey, e.Task.Id)
	skillDir := workDir + "/skill"
	entries, err := os.ReadDir(skillDir)
	if err != nil || len(entries) == 0 {
		return fmt.Errorf("skill output directory is empty: %s", skillDir)
	}

	// Find the directory that actually contains SKILL.md (not the first subdir)
	var genSkillDir string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		candidate := skillDir + "/" + entry.Name()
		if _, err := os.Stat(candidate + "/SKILL.md"); err == nil {
			genSkillDir = candidate
			break
		}
	}
	if genSkillDir == "" {
		// Fallback: pick the first subdirectory
		for _, entry := range entries {
			if entry.IsDir() {
				genSkillDir = skillDir + "/" + entry.Name()
				e.writeLogf("[Warning] SKILL.md not found, using first directory: %s\n", genSkillDir)
				break
			}
		}
	}
	if genSkillDir == "" {
		return fmt.Errorf("no skill directory found in: %s", skillDir)
	}
	usersDir := strings.ReplaceAll(define.UserSkillsDir, "<admin_user_id>", cast.ToString(e.Task.AdminUserId))
	usersSkillDir := fmt.Sprintf("%s/%d/%s", usersDir, e.Task.Id, filepath.Base(genSkillDir))
	_ = tool.MkDirAll(usersSkillDir)
	if err := copyDir(genSkillDir, usersSkillDir); err != nil {
		return fmt.Errorf("copy skill to skills_user failed: %v", err)
	}
	e.writeLogf("[Done] Skill saved to: %s\n", usersSkillDir)
	e.Task.SkillDir = usersSkillDir
	return nil
}

// ========== generate callback (iteration counting + context compression) ==========

// requestChatStream calls the LLM API in streaming mode and returns the merged
// result via schema.ConcatMessages. A progress log is emitted every 10s so
// long-running generations don't appear frozen.
// Relies on context deadline (BookToSkillTaskTimeout = 60min) for overall timeout
// protection; the watchdog goroutine in requestChatStreamWithState closes the
// stream on context cancel, which unblocks stream.Read() immediately.
func (e *BookToSkillExecutor) requestChatStream(
	callName string,
	messages []adaptor.ZhimaChatCompletionMessage,
	functionTools []adaptor.FunctionTool,
	temperature float32,
	maxToken int,
) (*schema.Message, int64, error) {
	task := e.Task
	ctx := e.CancelCtx

	// Progress logging goroutine: log every 10s so operators can see the call is alive
	done := make(chan struct{})
	defer close(done)
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		start := time.Now()
		for {
			select {
			case <-ticker.C:
				elapsed := int(time.Since(start).Seconds())
				e.writeLogf("  [LLM] %s streaming... waited %ds\n", callName, elapsed)
			case <-done:
				return
			}
		}
	}()

	robot := msql.Params{
		"id":               cast.ToString(task.RobotId),
		"robot_key":        task.RobotKey,
		"application_type": cast.ToString(define.ApplicationTypeClaw),
		"model_config_id":  cast.ToString(task.ModelConfigId),
		"use_model":        task.UseModel,
	}
	lang := task.Lang
	if lang == "" {
		lang = define.BookToSkillDefaultLang
	}
	handler, err := GetModelCallHandler(lang, task.AdminUserId, task.ModelConfigId, task.UseModel, robot)
	if err != nil {
		return nil, 0, fmt.Errorf("get model handler: %w", err)
	}

	// Create channel to receive stream_raw events, then collect frames and merge
	// via schema.ConcatMessages (same pattern as pipe_clawbot.Stream).
	chanStream := make(chan sse.Event, 100)

	var (
		streamErr   error
		totalResp   adaptor.ZhimaChatCompletionResponse
		requestTime int64
		frames      []*schema.Message
	)

	go func() {
		defer close(chanStream)
		totalResp, requestTime, streamErr = handler.RequestChatStream(
			ctx, lang, task.AdminUserId, "",
			robot, lib_define.ChatClawClient,
			messages, functionTools,
			chanStream,
			temperature, maxToken,
		)
	}()

	// Collect stream_raw frames, then ConcatMessages merges them.
	for event := range chanStream {
		if event.Event != `stream_raw` {
			continue
		}
		response, ok := event.Data.(adaptor.ZhimaChatCompletionResponse)
		if !ok {
			continue
		}
		// Fallback for legacy providers that stream a single tool call without index.
		if len(response.ToolCalls) == 1 && response.ToolCalls[0].Index == nil {
			idx := 0
			response.ToolCalls[0].Index = &idx
		}
		frames = append(frames, custom_eino.ConvertChatResp(response))
	}

	if streamErr != nil {
		return nil, requestTime, streamErr
	}

	// Check if context was cancelled (stop flag / task timeout)
	if ctx.Err() != nil {
		return nil, requestTime, ctx.Err()
	}

	// Merge collected frames via ConcatMessages for proper tool call assembly
	if len(frames) > 0 {
		merged, err := schema.ConcatMessages(frames)
		if err != nil {
			return nil, requestTime, fmt.Errorf("concat stream frames: %w", err)
		}
		sanitizeMessageToolCalls(merged)
		return merged, requestTime, nil
	}

	// No stream_raw frames collected (edge case): fall back to direct conversion.
	msg := custom_eino.ConvertChatResp(totalResp)
	sanitizeMessageToolCalls(msg)
	return msg, requestTime, nil
}

// describeInput builds a concise action summary from the last few input messages,
// so logs show what the agent is actually doing rather than just token counts.
func describeInput(input []*schema.Message) string {
	if len(input) == 0 {
		return "(empty input)"
	}
	// Walk backwards to find the last assistant message (most recent LLM decision)
	for i := len(input) - 1; i >= 0; i-- {
		msg := input[i]
		if msg == nil {
			continue
		}
		if msg.Role == schema.Assistant {
			if len(msg.ToolCalls) > 0 {
				// Show tool names (max 3) + target hints from first arg
				parts := make([]string, 0, len(msg.ToolCalls))
				for j, tc := range msg.ToolCalls {
					if j >= 3 {
						parts = append(parts, fmt.Sprintf("…+%d", len(msg.ToolCalls)-3))
						break
					}
					argHint := extractArgHint(tc.Function.Name, tc.Function.Arguments)
					if argHint != "" {
						parts = append(parts, fmt.Sprintf("%s %s", tc.Function.Name, argHint))
					} else {
						parts = append(parts, tc.Function.Name)
					}
				}
				return fmt.Sprintf("→ %s", strings.Join(parts, ", "))
			}
			// No tool calls: LLM returned a text response (usually end of step)
			text := strings.TrimSpace(msg.Content)
			if len(text) > 0 {
				return fmt.Sprintf("→ [Text] %s", truncateString(text, 80))
			}
			return "→ (empty reply)"
		}
	}
	// No assistant message found: first iteration, starting fresh
	return "→ (first round)"
}

// extractArgHint extracts a short path/filename hint from tool call arguments
// for common BTS tools (read/write/execute).
func extractArgHint(toolName string, args string) string {
	if args == "" {
		return ""
	}
	switch toolName {
	case "read", "write":
		// Extract file_path from JSON like {"file_path":"chunks/chunk_012.md"}
		if idx := strings.Index(args, `"file_path"`); idx >= 0 {
			rest := args[idx+len(`"file_path"`):]
			if colonIdx := strings.Index(rest, ":"); colonIdx >= 0 {
				val := strings.TrimSpace(rest[colonIdx+1:])
				val = strings.Trim(val, `,"`)
				val = strings.TrimSpace(val)
				if len(val) > 0 {
					// Show just the filename portion for readability
					if slashIdx := strings.LastIndex(val, "/"); slashIdx >= 0 {
						val = val[slashIdx+1:]
					}
					return val
				}
			}
		}
	case "execute":
		// Extract command from JSON like {"command":"python3 scripts/splitter.py ..."}
		if idx := strings.Index(args, `"command"`); idx >= 0 {
			rest := args[idx+len(`"command"`):]
			if colonIdx := strings.Index(rest, ":"); colonIdx >= 0 {
				val := strings.TrimSpace(rest[colonIdx+1:])
				val = strings.Trim(val, `,"`)
				val = strings.TrimSpace(val)
				// Truncate long commands
				if len(val) > 60 {
					val = val[:57] + "..."
				}
				return val
			}
		}
	case "grep_raw", "grep":
		// Extract regex content from JSON like {"content":"Ch.*","file_path":"..."}
		if idx := strings.Index(args, `"content"`); idx >= 0 {
			rest := args[idx+len(`"content"`):]
			if colonIdx := strings.Index(rest, ":"); colonIdx >= 0 {
				val := strings.TrimSpace(rest[colonIdx+1:])
				val = strings.Trim(val, `,"`)
				val = strings.TrimSpace(val)
				if len(val) > 50 {
					val = val[:47] + "..."
				}
				return fmt.Sprintf("pattern=%s", val)
			}
		}
	}
	return ""
}

// generate is the ChatModel Generate callback (shared across phases)
func (e *BookToSkillExecutor) generate(ctx context.Context, input []*schema.Message, opts custom_eino.RuntimeOptions) (*schema.Message, error) {
	task := e.Task

	// 0. Check stop flag before each LLM call (intra-step stop support)
	if checkStopFlag(task.RobotKey, task.Id) {
		e.writeLog("[Stop] generate detected .stop signal, interrupting agent\n")
		return nil, ErrTaskStopped
	}

	// 1. Iteration counting (cross-step cumulative)
	n := e.iterationCount.Add(1)
	if n%5 == 0 {
		go UpdateBookToSkillIteration(task.Id, int(n))
		e.writeLogf("[Iteration] DB updated iter=%d\n", n)
	}

	// 2. Context compression detection
	tokens := estimateTokens(input)
	actionHint := describeInput(input)
	e.writeLogf("[LLMReq] Round=%d Step=%d tokens=%d action=%s\n", n, e.currentStep, tokens, actionHint)

	// 2.5 Diagnostic: log phase + iteration info (no interception, just observability)
	if e.currentStep == 2 && (strings.Contains(actionHint, "read_file") || strings.Contains(actionHint, "glob_info")) {
		e.writeLogf("  [Reminder] batch phase: %s (remember to write)\n", actionHint)
	}

	if tokens > e.tokenThreshold {
		beforeTokens := tokens
		input = e.compressByLLM(ctx, input)
		afterTokens := estimateTokens(input)
		e.writeLogf("[Compress] before=%d after=%d\n", beforeTokens, afterTokens)
	}

	// 3. Original ConvertMessage + RequestChat logic
	var messages []adaptor.ZhimaChatCompletionMessage
	for _, msg := range input {
		converted, err := custom_eino.ConvertMessage(*msg)
		if err != nil {
			continue
		}
		messages = append(messages, converted)
	}

	var functionTools []adaptor.FunctionTool
	filterTools := []string{`write_todos`, `todo_write`}
	for _, toolInfo := range opts.Common.Tools {
		if toolInfo == nil {
			continue
		}
		if tool.InArray(toolInfo.Name, filterTools) {
			continue
		}
		converted, err := custom_eino.ConvertTools(*toolInfo)
		if err != nil {
			continue
		}
		functionTools = append(functionTools, converted)
	}

	callName := fmt.Sprintf("Round %d Step %d", n, e.currentStep)
	e.writeLogf("[LLMCall] %s msg_count=%d tool_count=%d temperature=0.7 maxToken=%d\n",
		callName, len(messages), len(functionTools), define.BookToSkillMaxToken)

	msg, requestTime, err := e.requestChatStream(
		callName,
		messages, functionTools,
		1,
		define.BookToSkillMaxToken,
	)
	if err != nil {
		e.writeLogf("[LLMError] %s request failed: %v\n", callName, err)
		return nil, err
	}

	// Log detailed LLM response
	if msg != nil {
		contentLen := len(msg.Content)
		toolCallCount := len(msg.ToolCalls)
		e.writeLogf("[LLMResp] %s elapsed=%dms content_len=%d tool_calls=%d\n",
			callName, requestTime, contentLen, toolCallCount)
		// Log full response content (truncated to 2000 for very long outputs)
		if contentLen > 0 {
			content := msg.Content
			if len(content) > 2000 {
				content = content[:2000] + "...(truncated)"
			}
			e.writeLogf("  [Response] %s\n", content)
		}
		// Log tool calls
		for i, tc := range msg.ToolCalls {
			argPreview := tc.Function.Arguments
			if len(argPreview) > 300 {
				argPreview = argPreview[:300] + "...(truncated)"
			}
			e.writeLogf("  [ToolCall%d] %s args=%s\n", i+1, tc.Function.Name, argPreview)
		}
	}

	return msg, nil
}

// ========== Context compression (LLM-based) ==========

func compressSystemPrompt(lang string) string {
	return define.GetBtsPrompt(lang, define.BtsPromptCompressSystem)
}
func compressUserPrompt(lang string) string {
	return define.GetBtsPrompt(lang, define.BtsPromptCompressUser)
}

// compressByLLM compresses middle history via a separate LLM call (no tools, no deep agent loop)
func (e *BookToSkillExecutor) compressByLLM(ctx context.Context, input []*schema.Message) []*schema.Message {
	if len(input) < 8 {
		return input // too few messages to compress
	}
	head := input[:2]            // system + first user task
	tail := input[len(input)-4:] // last 2 rounds (current chunk)
	middle := input[2 : len(input)-4]

	// Serialize middle history (truncate each message to avoid excessive length)
	var sb strings.Builder
	for _, msg := range middle {
		role := string(msg.Role)
		sb.WriteString(fmt.Sprintf("[%s] %s\n", role, truncateString(msg.Content, 500)))
	}

	lang := e.Task.Lang
	if lang == "" {
		lang = define.BookToSkillDefaultLang
	}

	// Build compression request (functionTools=nil, no tools)
	compressMsgs := []adaptor.ZhimaChatCompletionMessage{
		{Role: "system", Content: compressSystemPrompt(lang)},
		{Role: "user", Content: compressUserPrompt(lang) + sb.String()},
	}
	msg, _, err := e.requestChatStream(
		"Compress",
		compressMsgs, nil, 1, 2048,
	)
	if err != nil {
		e.writeLogf("[CompressFailed] Fallback to original history: %v\n", err)
		return input // compression failed, fallback to original
	}

	summary := msg.Content
	e.writeLogf("[CompressSummary] Length=%d Content:\n%s\n", len([]rune(summary)), truncateString(summary, 500))
	var summaryMsg *schema.Message
	summaryMsg = schema.UserMessage("【Processed Summary】\n" + summary)

	result := make([]*schema.Message, 0, len(head)+1+len(tail))
	result = append(result, head...)
	result = append(result, summaryMsg)
	result = append(result, tail...)
	return result
}

// estimateTokens roughly estimates token count (chars/3 for Chinese)
func estimateTokens(messages []*schema.Message) int {
	total := 0
	for _, msg := range messages {
		total += len([]rune(msg.Content))
		for _, tc := range msg.ToolCalls {
			total += len([]rune(tc.Function.Arguments))
		}
	}
	return total / 3
}

// StripCdPrefix removes a leading "cd <path> && " or "cd <path>; " prefix.
// LLM sometimes prepends this when trying to navigate to the working directory,
// but the working directory is already set by llm_runner.
func StripCdPrefix(cmd string) string {
	cmd = strings.TrimSpace(cmd)
	// Match "cd <path> && " or "cd <path> ; "
	if strings.HasPrefix(cmd, "cd ") {
		// Find the first "&&" or ";" after the path
		rest := cmd[3:] // after "cd "
		if idx := strings.Index(rest, "&&"); idx >= 0 {
			return strings.TrimSpace(rest[idx+2:])
		}
		if idx := strings.Index(rest, ";"); idx >= 0 {
			return strings.TrimSpace(rest[idx+1:])
		}
	}
	return cmd
}

// ========== Logging ==========

// logAgentEvent records an agent event to the log (concise, no full content)
func (e *BookToSkillExecutor) logAgentEvent(event *adk.AgentEvent) {
	if event.Action != nil && event.Action.Interrupted != nil {
		e.writeLogf("[Interrupted] %v\n", event.Action.Interrupted)
		return
	}
	if event.Output == nil || event.Output.MessageOutput == nil {
		return
	}
	mv := event.Output.MessageOutput
	if mv.Message == nil {
		return
	}
	msg := mv.Message
	switch mv.Role {
	case schema.Assistant:
		if len(msg.ToolCalls) > 0 {
			for _, call := range msg.ToolCalls {
				argLen := len(call.Function.Arguments)
				argHint := extractArgHint(call.Function.Name, call.Function.Arguments)
				if argHint != "" {
					e.writeLogf("[ToolCall] %s %s (args=%d)\n", call.Function.Name, argHint, argLen)
				} else {
					e.writeLogf("[ToolCall] %s args_len=%d\n", call.Function.Name, argLen)
				}
			}
		} else {
			textLen := len(msg.Content)
			snippet := truncateString(strings.TrimSpace(msg.Content), 60)
			e.writeLogf("[LLM] %s (len=%d)\n", snippet, textLen)
		}
	case schema.Tool:
		contentLen := len(msg.Content)
		toolName := msg.ToolName
		if toolName == "" {
			toolName = "?"
		}
		sanitized := strings.ReplaceAll(truncateString(strings.TrimSpace(msg.Content), 500), "\n", "\\n")
		e.writeLogf("[ToolResult] %s len=%d content: %s\n", toolName, contentLen, sanitized)
	}
}

// ========== Skill middleware ==========

// newBookToSkillMiddleware creates a skill middleware for book-to-skill
func newBookToSkillMiddleware(ctx context.Context, backend *custom_eino.Backend) (adk.ChatModelAgentMiddleware, error) {
	skillBackend, err := skill.NewBackendFromFilesystem(ctx, &skill.BackendFromFilesystemConfig{
		Backend: backend,
		BaseDir: define.BookToSkillTemplateDir,
	})
	if err != nil {
		return nil, err
	}
	return skill.NewMiddleware(ctx, &skill.Config{
		Backend: &custom_eino.LayeredSkillBackend{
			Backends: []skill.Backend{skillBackend},
		},
	})
}

// ========== Utility functions ==========

// ZipDir creates a zip archive of the source directory at the destination path
func ZipDir(srcDir, dstZip string) error {
	outFile, err := os.Create(dstZip)
	if err != nil {
		return fmt.Errorf("create zip file failed: %v", err)
	}
	defer outFile.Close()

	zw := zip.NewWriter(outFile)
	defer zw.Close()

	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(relPath)
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})
	if err != nil {
		return err
	}
	return zw.Close()
}

// truncateString truncates a string for logging
func truncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}

// ========== Init & Prepare ==========

// InitBookToSkillTaskDir creates task working directory + copies scripts/templates.
// Pre-creates skill/<skillName>/resources/ so the AI can write output files
// directly without needing to run mkdir scripts.
// Called after CreateBookToSkillTask (task_id available)
func InitBookToSkillTaskDir(robotKey string, taskId int, skillName string) error {
	workDir := bookToSkillWorkDir(robotKey, taskId)
	// Create root work dir with 0777 so llm_runner (uid 10001) can write
	_ = tool.MkDirAll(workDir)
	_ = os.Chmod(workDir, 0777)
	// Create all subdirectories, chmod 0777 so llm_runner (uid 10001) can write
	for _, dirConst := range []string{
		define.WorkInputDir, define.WorkScriptsDir, define.WorkTemplatesDir,
		define.WorkChunksDir, define.WorkOutputDir, define.WorkSummariesDir,
	} {
		dir := bookToSkillSubDir(robotKey, taskId, dirConst)
		_ = tool.MkDirAll(dir)
		_ = os.Chmod(dir, 0777)
	}
	// Pre-create skill output directories so the AI can write files directly
	// without needing Python mkdir scripts (which cause path duplication bugs).
	skillBase := workDir + "/skill/" + skillName
	skillResources := skillBase + "/resources"
	_ = tool.MkDirAll(skillBase)
	_ = os.Chmod(skillBase, 0777)
	_ = tool.MkDirAll(skillResources)
	_ = os.Chmod(skillResources, 0777)
	// Pre-create input_md/ for document conversion output
	inputMdDir := workDir + "/input_md"
	_ = tool.MkDirAll(inputMdDir)
	_ = os.Chmod(inputMdDir, 0777)
	// Copy entire template (SKILL.md + scripts + templates) to task directory
	templateDir := define.BookToSkillTemplateDir
	if err := copyDir(templateDir, workDir); err != nil {
		return fmt.Errorf("copy template failed: %v", err)
	}
	// Remove stale legacy directories (dictionaries/keywords from old
	// pre_extract.py, no longer used in new architecture; __pycache__ is useless if copied)
	for _, stale := range []string{"dictionaries", "keywords"} {
		_ = os.RemoveAll(workDir + "/" + stale)
	}
	_ = os.RemoveAll(workDir + "/scripts/__pycache__")
	return nil
}

// PrepareBookToSkillFiles copies uploaded files to working_dir/input
func PrepareBookToSkillFiles(task *BookToSkillTask, uploadInfos []*define.UploadInfo) error {
	inputDir := bookToSkillSubDir(task.RobotKey, task.Id, define.WorkInputDir)
	_ = tool.MkDirAll(inputDir)

	if len(uploadInfos) == 0 {
		// No files to prepare — keep existing input/ and SourceFiles untouched
		return nil
	}

	task.SourceFiles = make([]FileInfo, 0, len(uploadInfos))
	for i, ui := range uploadInfos {
		locFile := GetFileByLink(ui.Link)
		if locFile == "" {
			return fmt.Errorf("cannot access file: %s", ui.Name)
		}
		destPath := fmt.Sprintf("%s/%d_%s", inputDir, i, ui.Name)
		data, err := os.ReadFile(locFile)
		if err != nil {
			return fmt.Errorf("read file failed: %v", err)
		}
		if err := os.WriteFile(destPath, data, 0644); err != nil {
			return fmt.Errorf("write to working_dir failed: %v", err)
		}
		task.SourceFiles = append(task.SourceFiles, FileInfo{
			Name: ui.Name, Size: ui.Size, Ext: ui.Ext, Link: ui.Link,
		})
	}

	// Persist SourceFiles to DB so retries can find the original files
	sourceFilesJson, _ := json.Marshal(task.SourceFiles)
	table := getBookToSkillTable()
	_, _ = table.Where("id", cast.ToString(task.Id)).Update(msql.Datas{
		"source_files": string(sourceFilesJson),
		"update_time":  time.Now().Unix(),
	})

	return nil
}

// ========== Task lifecycle ==========

// RunBookToSkillTaskSync executes the book-to-skill pipeline synchronously.
// Called by NSQ consumer (already asynchronous) or StartBookToSkillTask.
// Supports stop: checks .stop flag before/during steps.
func RunBookToSkillTaskSync(task *BookToSkillTask) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(define.BookToSkillTaskTimeout)*time.Second)
	defer cancel()

	// Re-check DB status: if task was stopped while still Pending, skip execution
	if currentTask, err := GetBookToSkillTaskById(task.Id); err == nil && currentTask.Status == define.BookToSkillStatusStopped {
		// Task stopped before execution — write a note to run.log if possible
		workDir := bookToSkillWorkDir(task.RobotKey, task.Id)
		if f, e := os.OpenFile(workDir+"/run.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); e == nil {
			_, _ = f.WriteString(fmt.Sprintf("[Skip] Task was stopped in Pending stage, skipping execution\n"))
			_ = f.Close()
		}
		return
	}

	// Clear stale .stop flag (created by previous stop request)
	ClearStopFlag(task.RobotKey, task.Id)

	UpdateBookToSkillTaskStatus(task.Id, define.BookToSkillStatusRunning)

	executor := &BookToSkillExecutor{
		Task:      task,
		CancelCtx: ctx,
		Cancel:    cancel,
	}

	err := executor.Execute()
	// Sync final iteration count to DB
	UpdateBookToSkillIteration(task.Id, int(executor.iterationCount.Load()))

	// Check if task was stopped (Execute returns nil on graceful stop)
	if checkStopFlag(task.RobotKey, task.Id) {
		// .stop flag still exists — task was stopped, status already set by StopBookToSkillTask
		executor.writeLog("[Stop] Task stopped\n")
		return
	}

	if err != nil {
		failReason := err.Error()
		executor.writeLogf("\n[Failed] Reason: %s\n", failReason)
		UpdateBookToSkillTaskFailed(task.Id, failReason)
	} else {
		UpdateBookToSkillTaskSuccess(task.Id, task.SkillDir)
	}

	// Cleanup task working directory temp files (keep chunks/ + skill/)
	// Only delete input/ on success — keep it on failure so retry can reuse
	workDir := bookToSkillWorkDir(task.RobotKey, task.Id)
	if err == nil {
		_ = os.RemoveAll(workDir + "/input")
	}
}

// RetryBookToSkillTask clears the task working directory, resets all progress/result
// fields in DB, re-initializes the task directory, and re-prepares source files
// for a fresh start. Called by the retry HTTP handler.
func RetryBookToSkillTask(task *BookToSkillTask) error {
	workDir := bookToSkillWorkDir(task.RobotKey, task.Id)

	// 1. Clean old artifacts selectively — keep input/ (source files are immutable)
	for _, dir := range []string{"chunks", "chunks_raw", "summaries", "extracted", "skill", "input_md", "strategy"} {
		_ = os.RemoveAll(workDir + "/" + dir)
	}
	for _, file := range []string{"run.log", "summary_index.txt", "process.txt", "book_name.txt"} {
		_ = os.Remove(workDir + "/" + file)
	}

	// 1b. Clean old skill artifact in skills_user (from previous failed run)
	if task.SkillDir != "" {
		_ = os.RemoveAll(task.SkillDir)
		_ = os.RemoveAll(task.SkillDir + ".zip")
	}

	// 2. Clear stop flag (in case task was in stopped state)
	ClearStopFlag(task.RobotKey, task.Id)

	// 3. Reset DB: status→Pending, clear result fields
	table := getBookToSkillTable()
	_, err := table.Where("id", cast.ToString(task.Id)).Update(msql.Datas{
		"status":            define.BookToSkillStatusPending,
		"error_msg":         "",
		"skill_dir":         "",
		"current_iteration": "0",
		"update_time":       time.Now().Unix(),
	})
	if err != nil {
		return fmt.Errorf("reset task db failed: %w", err)
	}

	// 4. Re-init task working directory (create dirs + copy template scripts)
	if err := InitBookToSkillTaskDir(task.RobotKey, task.Id, task.SkillName); err != nil {
		return fmt.Errorf("re-init task dir failed: %w", err)
	}

	// 5. Re-prepare source files from stored FileInfo
	uploadInfos := make([]*define.UploadInfo, 0, len(task.SourceFiles))
	for _, sf := range task.SourceFiles {
		uploadInfos = append(uploadInfos, &define.UploadInfo{
			Name: sf.Name,
			Size: sf.Size,
			Ext:  sf.Ext,
			Link: sf.Link,
		})
	}
	if err := PrepareBookToSkillFiles(task, uploadInfos); err != nil {
		return fmt.Errorf("re-prepare files failed: %w", err)
	}

	return nil
}

// StopBookToSkillTask stops a running task
func StopBookToSkillTask(taskId int, cancel context.CancelFunc) error {
	if cancel != nil {
		cancel()
	}
	return updateTaskStatus(taskId, define.BookToSkillStatusStopped, "")
}

// ========== Task CRUD ==========

// CreateBookToSkillTask inserts a new task record.
// customTemplate is optional: if non-empty, used directly; if empty, read default template from disk.
func CreateBookToSkillTask(task *BookToSkillTask, customTemplate string) (int, error) {
	sourceFilesJson, _ := json.Marshal(task.SourceFiles)
	now := time.Now().Unix()

	// Read template content: use customTemplate if provided, otherwise load default
	var templateContent []byte
	if customTemplate != "" {
		templateContent = []byte(customTemplate)
	} else {
		templatePath := define.BookToSkillTemplateDir + "/templates/md_extraction.md"
		var err error
		templateContent, err = os.ReadFile(templatePath)
		if err != nil {
			logs.Debug("[BTS] Failed to read template file: %v, path=%s", err, templatePath)
			templateContent = []byte("")
		}
	}

	lang := task.Lang
	if lang == "" {
		lang = define.BookToSkillDefaultLang
	}

	table := getBookToSkillTable()
	id, err := table.Insert(msql.Datas{
		"admin_user_id":    task.AdminUserId,
		"robot_id":         task.RobotId,
		"robot_key":        task.RobotKey,
		"skill_name":       task.SkillName,
		"model_config_id":  task.ModelConfigId,
		"use_model":        task.UseModel,
		"status":           define.BookToSkillStatusPending,
		"source_files":     string(sourceFilesJson),
		"template_content": string(templateContent),
		"lang":             lang,
		"create_time":      now,
		"update_time":      now,
	}, "id")
	return int(id), err
}

// UpdateBookToSkillTaskStatus updates task status to running
func UpdateBookToSkillTaskStatus(taskId int, status int) {
	table := getBookToSkillTable()
	_, _ = table.Where("id", cast.ToString(taskId)).Update(msql.Datas{
		"status":      status,
		"update_time": time.Now().Unix(),
	})
}

// UpdateBookToSkillTaskSuccess marks task as success
func UpdateBookToSkillTaskSuccess(taskId int, skillDir string) {
	table := getBookToSkillTable()
	_, _ = table.Where("id", cast.ToString(taskId)).Update(msql.Datas{
		"status":      define.BookToSkillStatusSuccess,
		"skill_dir":   skillDir,
		"update_time": time.Now().Unix(),
	})
}

// UpdateBookToSkillTaskFailed marks task as failed with error
func UpdateBookToSkillTaskFailed(taskId int, errorMsg string) {
	table := getBookToSkillTable()
	_, _ = table.Where("id", cast.ToString(taskId)).Update(msql.Datas{
		"status":      define.BookToSkillStatusFailed,
		"error_msg":   errorMsg,
		"update_time": time.Now().Unix(),
	})
}

// UpdateBookToSkillIteration updates current_iteration field
func UpdateBookToSkillIteration(taskId int, iteration int) {
	table := getBookToSkillTable()
	_, _ = table.Where("id", cast.ToString(taskId)).Update(msql.Datas{
		"current_iteration": cast.ToString(iteration),
		"update_time":       time.Now().Unix(),
	})
}

func updateTaskStatus(taskId int, status int, errorMsg string) error {
	data := msql.Datas{
		"status":      status,
		"update_time": time.Now().Unix(),
	}
	if errorMsg != "" {
		data["error_msg"] = errorMsg
	}
	table := getBookToSkillTable()
	_, err := table.Where("id", cast.ToString(taskId)).Update(data)
	return err
}

// GetBookToSkillTask retrieves a task by id (with adminUserId check)
func GetBookToSkillTask(taskId int, adminUserId int) (*BookToSkillTask, error) {
	table := getBookToSkillTable()
	list, err := table.Where("id", cast.ToString(taskId)).
		Where("admin_user_id", cast.ToString(adminUserId)).
		Limit(0, 1).Select()
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("task not found")
	}
	return rowToTask(list[0]), nil
}

// GetBookToSkillTaskById retrieves a task by id only (for internal/NSQ consumer use)
func GetBookToSkillTaskById(taskId int) (*BookToSkillTask, error) {
	table := getBookToSkillTable()
	list, err := table.Where("id", cast.ToString(taskId)).
		Limit(0, 1).Select()
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("task not found")
	}
	return rowToTask(list[0]), nil
}

// GetBookToSkillTaskList retrieves task list with pagination
func GetBookToSkillTaskList(adminUserId, robotId, status, page, pageSize int) ([]*BookToSkillTask, int, error) {
	tableC := getBookToSkillTable().
		Where("admin_user_id", cast.ToString(adminUserId))
	tableL := getBookToSkillTable().
		Where("admin_user_id", cast.ToString(adminUserId))
	if robotId > 0 {
		tableC.Where("robot_id", cast.ToString(robotId))
		tableL.Where("robot_id", cast.ToString(robotId))
	}
	if status > 0 {
		tableC.Where("status", cast.ToString(status))
		tableL.Where("status", cast.ToString(status))
	}

	total, _ := tableC.Count()
	list, err := tableL.Limit((page-1)*pageSize, pageSize).Order("id DESC").Select()
	if err != nil {
		return nil, 0, err
	}

	tasks := make([]*BookToSkillTask, 0, len(list))
	for _, row := range list {
		tasks = append(tasks, rowToTask(row))
	}
	return tasks, total, nil
}

// InstallBookToSkill installs a BookToSkill-generated skill to a robot.
// Uses the standard skill management flow: insert into chat_ai_clawbot_user_skill
// then call syncOneRobotSkill (which handles DB insert + file copy + backup/rollback).
func InstallBookToSkill(task *BookToSkillTask, targetRobotKey string) error {
	if task.SkillDir == "" || !tool.IsDir(task.SkillDir) {
		return fmt.Errorf("skill files not found")
	}

	// 1. Parse SKILL.md for description (frontmatter name may be Chinese book title,
	//    e.g. "Legend of Sword"; system name comes from task.SkillName, always ASCII).
	skillMdPath := findSkillMdInDir(task.SkillDir)
	var description, frontName string
	if skillMdPath != "" {
		if content, rerr := tool.ReadFile(skillMdPath); rerr == nil {
			frontName, description, _, _ = parseSkillFrontmatter(content)
		}
	}

	// 2. Use task.SkillName as system directory name (must match SkillNameRegexp)
	systemName := task.SkillName
	if !define.SkillNameRegexp.MatchString(systemName) {
		systemName = fmt.Sprintf("bts_%d", task.Id)
	}
	remarkName := frontName
	if remarkName == "" {
		remarkName = systemName
	}

	// 3. Copy skill to user skills directory (required by syncOneRobotSkill)
	adminUserId := task.AdminUserId
	InitClawbotUserDirs(adminUserId)
	destUserDir := filepath.Join(userSkillsDir(adminUserId), systemName)
	if err := copyDir(task.SkillDir, destUserDir); err != nil {
		return fmt.Errorf("copy to user skills failed: %v", err)
	}

	// Rewrite SKILL.md frontmatter name in the copied dir if it differs from systemName
	if systemName != frontName && frontName != "" {
		if mdPath := findSkillMdInDir(destUserDir); mdPath != "" {
			_ = rewriteSkillMdFrontmatter(mdPath, systemName, "")
		}
	}

	// 4. Insert/update chat_ai_clawbot_user_skill
	now := tool.Time2Int()
	intro := description
	if len([]rune(intro)) > 500 {
		intro = string([]rune(intro)[:500])
	}
	userSkillId, err := msql.Model(define.TableChatAiClawbotUserSkill, define.Postgres).Insert(msql.Datas{
		"admin_user_id":    adminUserId,
		"skill_name":       systemName,
		"remark_name":      remarkName,
		"intro":            intro,
		"description":      description,
		"file_size":        0,
		"origin_file_name": task.SkillName,
		"create_time":      now,
		"update_time":      now,
	}, "id")
	if err != nil {
		// Unique constraint (admin_user_id, skill_name): already exists, get existing id
		existing, qerr := msql.Model(define.TableChatAiClawbotUserSkill, define.Postgres).
			Where("admin_user_id", cast.ToString(adminUserId)).
			Where("skill_name", systemName).
			Find()
		if qerr != nil || len(existing) == 0 {
			return fmt.Errorf("insert user skill failed: %v", err)
		}
		userSkillId = cast.ToInt64(existing["id"])
	}

	// 5. Build userRow and check existing robot skill for update
	userRow := msql.Params{
		"id":          cast.ToString(userSkillId),
		"skill_name":  systemName,
		"remark_name": remarkName,
		"intro":       intro,
		"description": description,
		"file_size":   "0",
	}
	existingRobot, _ := msql.Model(define.TableChatAiClawbotSkill, define.Postgres).
		Where("robot_key", targetRobotKey).
		Where("source_type", cast.ToString(define.SkillSourceTypeUpload)).
		Where("skill_name", systemName).
		Find()
	var oldRow msql.Params
	if len(existingRobot) > 0 {
		oldRow = existingRobot
	}

	// 6. Sync user skill → robot skill (DB insert + file copy with backup/rollback)
	errKey, err := syncOneRobotSkill(adminUserId, targetRobotKey, userRow, oldRow)
	if err != nil || errKey != "" {
		return fmt.Errorf("sync robot skill failed: %v (key=%s)", err, errKey)
	}

	// 7. Update task install status
	table := getBookToSkillTable()
	_, err = table.Where("id", cast.ToString(task.Id)).Update(msql.Datas{
		"install_status": define.BookToSkillInstallInstalled,
		"update_time":    time.Now().Unix(),
	})
	return err
}

// sanitizeMessageToolCalls fixes common JSON escape issues in LLM-generated tool call arguments.
// AI models sometimes embed regex patterns (grep, sed) with backslash escapes like \d, \+, \w
// that are not valid JSON escape sequences, causing eino's tool node to fail on unmarshal.
func sanitizeMessageToolCalls(msg *schema.Message) {
	if msg == nil {
		return
	}
	for i := range msg.ToolCalls {
		msg.ToolCalls[i].Function.Arguments = sanitizeToolCallArguments(msg.ToolCalls[i].Function.Arguments)
	}
}

// sanitizeToolCallArguments attempts to fix JSON with invalid escape sequences.
// Returns the original string if it's already valid JSON or if fixing fails.
func sanitizeToolCallArguments(args string) string {
	if args == "" {
		return args
	}
	// Fast path: if it's already valid JSON, return as-is
	var test map[string]any
	if json.Unmarshal([]byte(args), &test) == nil {
		return args
	}

	// Try to fix: escape backslashes that precede non-JSON-escape characters
	fixed := fixJSONEscapes(args)
	if json.Unmarshal([]byte(fixed), &test) == nil {
		return fixed
	}

	// Fix failed, return original (let the error propagate naturally)
	return args
}

// fixJSONEscapes adds an extra backslash before backslash sequences that are
// not valid JSON escape sequences. Valid JSON escapes: \" \\ \/ \b \f \n \r \t \u
func fixJSONEscapes(s string) string {
	var b strings.Builder
	b.Grow(len(s) * 2)
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			next := s[i+1]
			switch next {
			case '"', '\\', '/', 'b', 'f', 'n', 'r', 't', 'u':
				// Valid JSON escape - keep as-is
				b.WriteByte('\\')
				b.WriteByte(next)
				i++
			default:
				// Not a valid JSON escape - add extra backslash
				b.WriteByte('\\')
				b.WriteByte('\\')
				b.WriteByte(next)
				i++
			}
		} else {
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

func getBookToSkillTable() *msql.Builder {
	return msql.Model("chat_ai_book_to_skill_task", define.Postgres)
}

func rowToTask(row msql.Params) *BookToSkillTask {
	task := &BookToSkillTask{
		Id:               cast.ToInt(row["id"]),
		AdminUserId:      cast.ToInt(row["admin_user_id"]),
		RobotId:          cast.ToInt(row["robot_id"]),
		RobotKey:         cast.ToString(row["robot_key"]),
		SkillName:        cast.ToString(row["skill_name"]),
		ModelConfigId:    cast.ToInt(row["model_config_id"]),
		UseModel:         cast.ToString(row["use_model"]),
		Status:           cast.ToInt(row["status"]),
		SkillDir:         cast.ToString(row["skill_dir"]),
		SkillFilePath:    cast.ToString(row["skill_file_path"]),
		ErrorMsg:         cast.ToString(row["error_msg"]),
		TemplateContent:  cast.ToString(row["template_content"]),
		InstallStatus:    cast.ToInt(row["install_status"]),
		CurrentIteration: cast.ToInt(row["current_iteration"]),
		Lang:             cast.ToString(row["lang"]),
		CreateTime:       cast.ToInt64(row["create_time"]),
		UpdateTime:       cast.ToInt64(row["update_time"]),
	}
	if sourceFilesStr := cast.ToString(row["source_files"]); sourceFilesStr != "" {
		var files []FileInfo
		if err := json.Unmarshal([]byte(sourceFilesStr), &files); err == nil {
			task.SourceFiles = files
		}
	}
	return task
}
