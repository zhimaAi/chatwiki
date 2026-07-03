// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package custom_eino

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cloudwego/eino/adk/filesystem"
	"github.com/matiasinsaurralde/go-e2b"
)

type E2BShell struct {
	ctx           context.Context
	client        *e2b.Client
	sandboxConfig e2b.SandboxConfig
	runOptions    []e2b.RunOption

	mu      sync.Mutex
	sandbox *e2b.Sandbox
}

func NewE2BShell(ctx context.Context, clientConfig e2b.ClientConfig, sandboxConfig e2b.SandboxConfig, runOptions ...e2b.RunOption) (*E2BShell, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	client, err := e2b.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}

	return &E2BShell{
		ctx:           ctx,
		client:        client,
		sandboxConfig: sandboxConfig,
		runOptions:    runOptions,
	}, nil
}

func (s *E2BShell) Execute(command string) (*filesystem.ExecuteResponse, error) {
	if strings.TrimSpace(command) == "" {
		return nil, fmt.Errorf("command is required")
	}

	sb, err := s.ensureSandbox()
	if err != nil {
		return nil, err
	}

	result, err := sb.Commands.RunWithContext(s.ctx, "bash", []string{"-lc", command}, s.runOptions...)
	if err != nil {
		return nil, err
	}

	return commandExecuteResponse(result.Stdout, result.Stderr, result.ExitCode), nil
}

func (s *E2BShell) Close() error {
	s.mu.Lock()
	sb := s.sandbox
	s.sandbox = nil
	s.mu.Unlock()

	if sb == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return sb.CloseWithContext(ctx)
}

func (s *E2BShell) ensureSandbox() (*e2b.Sandbox, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.sandbox != nil {
		return s.sandbox, nil
	}

	sb, err := s.client.NewSandbox(s.ctx, s.sandboxConfig)
	if err != nil {
		return nil, err
	}

	s.sandbox = sb
	return sb, nil
}

func commandExecuteResponse(stdout, stderr string, exitCode int) *filesystem.ExecuteResponse {
	output := stdout
	if exitCode != 0 {
		parts := []string{fmt.Sprintf("command exited with non-zero code %d", exitCode)}
		if stdout != "" {
			parts = append(parts, "[stdout]:\n"+strings.TrimRight(stdout, "\r\n"))
		}
		if stderr != "" {
			parts = append(parts, "[stderr]:\n"+strings.TrimRight(stderr, "\r\n"))
		}
		output = strings.Join(parts, "\n")
	} else if output == "" && stderr != "" {
		output = "[stderr]:\n" + strings.TrimRight(stderr, "\r\n")
	}
	return &filesystem.ExecuteResponse{Output: output, ExitCode: &exitCode}
}
