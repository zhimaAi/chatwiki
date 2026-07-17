// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package llm_runner

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/zhimaAi/go_tools/logs"
)

type CommandRunner struct {
	mu     sync.Mutex
	active map[*exec.Cmd]string
}

const defaultCommandTimeout = 5 * time.Minute

func (r *CommandRunner) Run(req RpcRunRequest, resp *RpcRunResponse) (_ error) {
	fmt.Println(`request command:` + req.Command)
	defer func() {
		if resp.IsError {
			fmt.Println(`command:` + req.Command + `,error:` + resp.ErrorMsg)
		}
	}()
	timeout := commandTimeout()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", req.Command)
	configureCommandProcessGroup(cmd, cancel)
	if req.WorkDir != "" {
		cmd.Dir = req.WorkDir
	}
	runID := strings.TrimSpace(req.RunID)
	if runID != `` {
		r.register(cmd, runID)
		defer r.unregister(cmd)
	}
	var stdoutBuf, stderrBuf strings.Builder
	cmd.Stdout, cmd.Stderr = &stdoutBuf, &stderrBuf
	exitCode := 0
	err := cmd.Run()
	if err == nil {
		resp.Output = stdoutBuf.String()
		return // success
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		resp.IsError = true
		resp.ErrorMsg = fmt.Sprintf(`command timed out after %s`, timeout)
		resp.ExitCode = -1
		return // timed out
	}
	if errors.Is(ctx.Err(), context.Canceled) {
		resp.IsError = true
		resp.Cancelled = true
		resp.ErrorMsg = `command cancelled`
		resp.ExitCode = -1
		return
	}
	var exitError *exec.ExitError
	if !errors.As(err, &exitError) {
		resp.IsError = true
		resp.ErrorMsg = fmt.Sprintf(`failed to execute command: %v`, err)
		resp.ExitCode = -1
		return // failed to execute command
	}
	exitCode = exitError.ExitCode()
	stdoutStr := stdoutBuf.String()
	stderrStr := stderrBuf.String()
	parts := []string{fmt.Sprintf(`command exited with non-zero code %d`, exitCode)}
	if stdoutStr != "" {
		parts = append(parts, "[stdout]:\n"+strings.TrimSuffix(stdoutStr, ""))
	}
	if stderrStr != "" {
		parts = append(parts, "[stderr]:\n"+strings.TrimSuffix(stderrStr, ""))
	}
	resp.Output = strings.Join(parts, "\n")
	resp.ExitCode = exitCode
	return // command exited with non-zero code
}

func (r *CommandRunner) Cancel(req RpcCancelRequest, resp *RpcCancelResponse) (_ error) {
	runID := strings.TrimSpace(req.RunID)
	if runID == `` {
		resp.ErrorMsg = `run id is required`
		return
	}
	r.mu.Lock()
	commands := make([]*exec.Cmd, 0)
	for cmd, activeRunID := range r.active {
		if activeRunID == runID {
			commands = append(commands, cmd)
		}
	}
	r.mu.Unlock()
	if len(commands) == 0 {
		return
	}
	resp.Found = true
	var cancelErrors []error
	for _, cmd := range commands {
		if err := cmd.Cancel(); err != nil {
			cancelErrors = append(cancelErrors, err)
		}
	}
	if err := errors.Join(cancelErrors...); err != nil {
		resp.ErrorMsg = err.Error()
	}
	resp.Cancelled = true
	return
}

func (r *CommandRunner) register(cmd *exec.Cmd, runID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.active == nil {
		r.active = make(map[*exec.Cmd]string)
	}
	r.active[cmd] = runID
}

func (r *CommandRunner) unregister(cmd *exec.Cmd) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.active, cmd)
}

func commandTimeout() time.Duration {
	value := strings.TrimSpace(os.Getenv("LLM_RUNNER_COMMAND_TIMEOUT"))
	if value == "" {
		return defaultCommandTimeout
	}
	timeout, err := time.ParseDuration(value)
	if err != nil || timeout <= 0 {
		return defaultCommandTimeout
	}
	return timeout
}

func StartRpcService() {
	if err := rpc.RegisterName(RpcServiceName, &CommandRunner{}); err != nil {
		panic(err)
	}
	listener, err := net.Listen(`tcp`, RpcServicePort)
	if err != nil {
		panic(err)
	}
	defer func(listener net.Listener) {
		_ = listener.Close()
	}(listener)
	logs.Info(`command rpc server listening on %s`, RpcServicePort)
	for {
		conn, err := listener.Accept()
		if err != nil {
			logs.Error(`accept connection: %v`, err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
