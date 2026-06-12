// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package custom_eino

import (
	"context"
	"fmt"
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
	ValidateFile    func(op FileOperation, path string) error
	ValidateCommand func(command string) error
	ExecuteCommand  func(command string) (*filesystem.ExecuteResponse, error)
}

type Backend struct {
	*local.Local
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
		validateFile:    cfg.ValidateFile,
		validateCommand: cfg.ValidateCommand,
		executeCommand:  cfg.ExecuteCommand,
	}, nil
}

func (b *Backend) LsInfo(ctx context.Context, req *filesystem.LsInfoRequest) ([]filesystem.FileInfo, error) {
	if err := b.validateFileAccess(FileOperationLsInfo, req.Path); err != nil {
		return nil, err
	}
	return b.Local.LsInfo(ctx, req)
}

func (b *Backend) Read(ctx context.Context, req *filesystem.ReadRequest) (*filesystem.FileContent, error) {
	if err := b.validateFileAccess(FileOperationRead, req.FilePath); err != nil {
		return nil, err
	}
	return b.Local.Read(ctx, req)
}

func (b *Backend) GrepRaw(ctx context.Context, req *filesystem.GrepRequest) ([]filesystem.GrepMatch, error) {
	if err := b.validateFileAccess(FileOperationGrepRaw, req.Path); err != nil {
		return nil, err
	}
	return b.Local.GrepRaw(ctx, req)
}

func (b *Backend) GlobInfo(ctx context.Context, req *filesystem.GlobInfoRequest) ([]filesystem.FileInfo, error) {
	if err := b.validateFileAccess(FileOperationGlobInfo, req.Path); err != nil {
		return nil, err
	}
	return b.Local.GlobInfo(ctx, req)
}

func (b *Backend) Write(ctx context.Context, req *filesystem.WriteRequest) error {
	if err := b.validateFileAccess(FileOperationWrite, req.FilePath); err != nil {
		return err
	}
	return b.Local.Write(ctx, req)
}

func (b *Backend) Edit(ctx context.Context, req *filesystem.EditRequest) error {
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
			return nil, err
		}
	}
	if b.executeCommand != nil {
		return b.executeCommand(req.Command)
	}
	return b.Local.Execute(ctx, req)
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
