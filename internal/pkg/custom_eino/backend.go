// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package custom_eino

import (
	"chatwiki/internal/pkg/llm_runner"
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino-ext/adk/backend/local"
	"github.com/cloudwego/eino/adk/filesystem"
)

type FileOperation string

const (
	FileOperationLsInfo   FileOperation = "ls_info"
	FileOperationRead     FileOperation = "read"
	FileOperationGrepRaw  FileOperation = "grep_raw"
	FileOperationGlobInfo FileOperation = "glob_info"
	FileOperationWrite    FileOperation = "write"
	FileOperationEdit     FileOperation = "edit"
)

type BackendConfig struct {
	WorkDir         string // task working directory for resolving relative paths
	ErrorAsContent  bool   // if true, validation/file errors are returned as content instead of error
	ValidateFile    func(op FileOperation, path string) error
	ValidateCommand func(command string) error
	ExecuteCommand  func(command string) (*filesystem.ExecuteResponse, error)
}

type Backend struct {
	*local.Local
	workDir         string // task working directory for resolving relative paths
	errorAsContent  bool   // if true, validation/file errors are returned as content instead of error
	validateFile    func(op FileOperation, path string) error
	validateCommand func(command string) error
	executeCommand  func(command string) (*filesystem.ExecuteResponse, error)
}

func NewBackend(ctx context.Context, cfg *BackendConfig) (*Backend, error) {
	if cfg == nil {
		cfg = &BackendConfig{}
	}

	localBackend, err := local.NewBackend(ctx, &local.Config{
		ValidateCommand: cfg.ValidateCommand,
	})
	if err != nil {
		return nil, err
	}

	return &Backend{
		Local:           localBackend,
		workDir:         cfg.WorkDir,
		errorAsContent:  cfg.ErrorAsContent,
		validateFile:    cfg.ValidateFile,
		validateCommand: cfg.ValidateCommand,
		executeCommand:  cfg.ExecuteCommand,
	}, nil
}

func (b *Backend) LsInfo(ctx context.Context, req *filesystem.LsInfoRequest) ([]filesystem.FileInfo, error) {
	req.Path = b.resolvePath(req.Path)
	if err := b.validateFileAccess(FileOperationLsInfo, req.Path); err != nil {
		return nil, err
	}
	return b.Local.LsInfo(ctx, req)
}

func (b *Backend) Read(ctx context.Context, req *filesystem.ReadRequest) (*filesystem.FileContent, error) {
	req.FilePath = b.resolvePath(req.FilePath)
	if err := b.validateFileAccess(FileOperationRead, req.FilePath); err != nil {
		if b.errorAsContent {
			return &filesystem.FileContent{
				Content: fmt.Sprintf("Error: %s", err.Error()),
			}, nil
		}
		return nil, err
	}
	content, err := b.Local.Read(ctx, req)
	if err != nil {
		if b.errorAsContent {
			return &filesystem.FileContent{
				Content: fmt.Sprintf("Error: %s", err.Error()),
			}, nil
		}
		return nil, err
	}
	return content, nil
}

func (b *Backend) GrepRaw(ctx context.Context, req *filesystem.GrepRequest) ([]filesystem.GrepMatch, error) {
	req.Path = b.resolvePath(req.Path)
	if err := b.validateFileAccess(FileOperationGrepRaw, req.Path); err != nil {
		return nil, err
	}
	return b.Local.GrepRaw(ctx, req)
}

func (b *Backend) GlobInfo(ctx context.Context, req *filesystem.GlobInfoRequest) ([]filesystem.FileInfo, error) {
	req.Path = b.resolvePath(req.Path)
	if err := b.validateFileAccess(FileOperationGlobInfo, req.Path); err != nil {
		return nil, err
	}
	return b.Local.GlobInfo(ctx, req)
}

func (b *Backend) Write(ctx context.Context, req *filesystem.WriteRequest) error {
	req.FilePath = b.resolvePath(req.FilePath)
	if err := b.validateFileAccess(FileOperationWrite, req.FilePath); err != nil {
		return err
	}
	return b.Local.Write(ctx, req)
}

func (b *Backend) Edit(ctx context.Context, req *filesystem.EditRequest) error {
	req.FilePath = b.resolvePath(req.FilePath)
	if err := b.validateFileAccess(FileOperationEdit, req.FilePath); err != nil {
		return err
	}
	return b.Local.Edit(ctx, req)
}

func (b *Backend) Execute(ctx context.Context, req *filesystem.ExecuteRequest) (*filesystem.ExecuteResponse, error) {
	if req == nil || strings.TrimSpace(req.Command) == `` {
		return nil, fmt.Errorf("command is required")
	}
	if b.validateCommand != nil {
		if err := b.validateCommand(req.Command); err != nil {
			if b.errorAsContent {
				exitCode := 1
				return &filesystem.ExecuteResponse{
					Output:   fmt.Sprintf("Error: %s", err.Error()),
					ExitCode: &exitCode,
				}, nil
			}
			return nil, err
		}
	}
	if b.executeCommand != nil {
		return b.executeCommand(req.Command)
	}
	return b.Local.Execute(ctx, req)
}

// resolvePath converts relative paths to full clawbot-relative paths.
// The local backend runs on the go_core process (not the llm_runner container),
// so relative paths like "input/foo.txt" would resolve against the wrong directory.
// By prepending workDir, the path becomes "clawbot/working_dir/<robot>/<task>/input/foo.txt"
// which resolves correctly from the go_core process root.
func (b *Backend) resolvePath(path string) string {
	path = strings.ReplaceAll(strings.TrimSpace(path), "\\", "/")
	if path == "" || path == "." || path == "./" {
		return b.workDir
	}
	// absolute path correction
	if strings.HasPrefix(path, `/`) {
		if rel := llm_runner.StripToClawbotAnchor(path, "clawbot/"); rel != `` {
			path = rel
		}
	}
	// Already a clawbot-relative or absolute path — don't modify
	if strings.HasPrefix(path, "clawbot/") || strings.HasPrefix(path, `/`) {
		return path
	}
	// Relative path — resolve against workDir
	if b.workDir != "" {
		return filepath.ToSlash(filepath.Join(b.workDir, path))
	}
	return path
}

func (b *Backend) validateFileAccess(op FileOperation, path string) error {
	if b.validateFile == nil {
		return nil
	}
	if strings.TrimSpace(path) == "" {
		path = "."
	}
	return b.validateFile(op, path)
}
