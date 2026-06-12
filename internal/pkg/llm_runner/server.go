// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package llm_runner

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"os/exec"
	"strings"
	"time"

	"github.com/zhimaAi/go_tools/logs"
)

type CommandRunner struct{}

const commandTimeout = 5 * time.Minute

func (r *CommandRunner) Run(req RpcRunRequest, resp *RpcRunResponse) (_ error) {
	fmt.Println(`request command:` + req.Command)
	defer func() {
		if resp.IsError {
			fmt.Println(`command:` + req.Command + `,error:` + resp.ErrorMsg)
		}
	}()
	args, err := PrepareCommand(req.Command)
	if err != nil {
		resp.IsError = true
		resp.ErrorMsg = err.Error()
		resp.ExitCode = -1
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	var stdoutBuf, stderrBuf strings.Builder
	cmd.Stdout, cmd.Stderr = &stdoutBuf, &stderrBuf
	exitCode := 0
	err = cmd.Run()
	if err == nil {
		resp.Output = stdoutBuf.String()
		return // success
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		resp.IsError = true
		resp.ErrorMsg = fmt.Sprintf(`command timed out after %s`, commandTimeout)
		resp.ExitCode = -1
		return // timed out
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
