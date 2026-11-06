// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/zhimaAi/go_tools/tool"
)

var (
	dangerousKeywords = []string{
		"import os", "import sys", "import subprocess", "import shutil",
		"import urllib", "import requests", "import http", "import socket",
		"import ctypes", "import pickle", "import marshal", "import eval",
		"import exec", "import compile", "from os import", "from sys import",
		"from subprocess import", "from shutil import", "from urllib import",
		"from requests import", "from http import", "from socket import",
		"from ctypes import", "from pickle import", "from marshal import",
		"__import__", "eval(", "exec(", "compile(",
		"open(", "file(", "os.system", "os.popen", "os.spawn",
		"subprocess.call", "subprocess.run", "subprocess.Popen",
		"shutil.rmtree", "shutil.copy", "shutil.move",
		"urllib.request", "requests.get", "requests.post",
		"socket.socket", "socket.connect",
	}
)

func checkPythonCodeSafety(code string) error {
	lowerCode := strings.ToLower(code)

	for _, keyword := range dangerousKeywords {
		if strings.Contains(lowerCode, strings.ToLower(keyword)) {
			return fmt.Errorf("dangerous operation detected: %s", keyword)
		}
	}

	dangerousPatterns := []string{
		`import\s+\*`,
		`__import__\s*\(`,
		`eval\s*\(`,
		`exec\s*\(`,
		`compile\s*\(`,
		`open\s*\(`,
		`file\s*\(`,
		`globals\s*\(\)`,
		`locals\s*\(\)`,
		`__builtins__`,
		`__dict__`,
	}

	for _, pattern := range dangerousPatterns {
		matched, err := regexp.MatchString(pattern, lowerCode)
		if err != nil {
			return fmt.Errorf("safety check error: %v", err)
		}
		if matched {
			return fmt.Errorf("dangerous pattern detected: %s", pattern)
		}
	}

	return nil
}

//go:embed python_runner.py
var pythonRunner string

func RunPython(mainFunc, batchNo string, params map[string]any) (string, error) {
	return RunPythonContext(context.Background(), mainFunc, batchNo, params)
}

func RunPythonContext(parent context.Context, mainFunc, batchNo string, params map[string]any) (string, error) {
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		return "", fmt.Errorf("python sandbox requires Linux amd64")
	}
	if err := checkPythonCodeSafety(mainFunc); err != nil {
		return "", fmt.Errorf("code safety check failed: %w", err)
	}
	if params == nil {
		params = make(map[string]any)
	}
	body, err := tool.JsonEncode(struct {
		MainFunc string         `json:"main_func"`
		Params   map[string]any `json:"params"`
	}{mainFunc, params})
	if err != nil {
		return "", fmt.Errorf("encode python request: %w", err)
	}
	ctx, cancel := context.WithTimeout(parent, 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "/usr/local/bin/python", "-I", "-B", "-c", pythonRunner)
	cmd.Env = []string{"PATH=/usr/local/bin:/usr/bin:/bin", "LANG=C.UTF-8", "HOME=/nonexistent"}
	// Preserve the configured timezone without inheriting unrelated host variables.
	if tz, ok := os.LookupEnv("TZ"); ok {
		cmd.Env = append(cmd.Env, "TZ="+tz)
	}
	cmd.Stdin = strings.NewReader(body)
	cmd.WaitDelay = time.Second
	output := &pythonOutput{}
	cmd.Stdout, cmd.Stderr = output, output
	err = cmd.Run() // Run waits for process reaping, including cancellation.
	if ctx.Err() != nil {
		return "", fmt.Errorf("python execution stopped: %w", ctx.Err())
	}
	if output.truncated {
		return "", fmt.Errorf("python output exceeds 1 MiB")
	}
	if err != nil {
		return "", fmt.Errorf("python execution error: %w, output: %s", err, output.text.String())
	}
	return strings.TrimSpace(output.text.String()), nil
}

type pythonOutput struct {
	text      strings.Builder
	truncated bool
}

func (w *pythonOutput) Write(p []byte) (int, error) {
	n := len(p)
	remaining := (1 << 20) - w.text.Len()
	if len(p) > remaining {
		w.truncated = true
		p = p[:remaining]
	}
	_, _ = w.text.Write(p)
	return n, nil
}
