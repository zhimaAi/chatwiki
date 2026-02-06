// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/zhimaAi/go_tools/logs"
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

func RunPython(mainFunc, batchNo string, params map[string]any) (string, error) {
	if err := checkPythonCodeSafety(mainFunc); err != nil {
		return ``, fmt.Errorf("code safety check failed: %w", err)
	}

	if params == nil {
		params = make(map[string]any)
	}

	paramsJson := tool.JsonEncodeNoError(params)

	indentedMainFunc := strings.ReplaceAll(mainFunc, "\n", "\n        ")

	logs.Debug("indentedMainFunc :%s", indentedMainFunc)

	wrapperCode := fmt.Sprintf(`import json
import sys
import math
import random
import datetime
import re
from typing import Any, Dict, List, Tuple, Optional

class SafeDict(dict):
    def __setitem__(self, key, value):
        raise PermissionError("Modifying globals is not allowed")

__builtins__ = SafeDict({
    'abs': abs,
    'all': all,
    'any': any,
    'bin': bin,
    'bool': bool,
    'chr': chr,
    'complex': complex,
    'dict': dict,
    'divmod': divmod,
    'enumerate': enumerate,
    'filter': filter,
    'float': float,
    'format': format,
    'frozenset': frozenset,
    'hex': hex,
    'int': int,
    'isinstance': isinstance,
    'issubclass': issubclass,
    'iter': iter,
    'len': len,
    'list': list,
    'map': map,
    'max': max,
    'min': min,
    'next': next,
    'oct': oct,
    'ord': ord,
    'pow': pow,
    'print': print,
    'range': range,
    'repr': repr,
    'reversed': reversed,
    'round': round,
    'set': set,
    'slice': slice,
    'sorted': sorted,
    'str': str,
    'sum': sum,
    'tuple': tuple,
    'type': type,
    'zip': zip,
    'Exception': Exception,
    'PermissionError': PermissionError,
    'TypeError': TypeError,
    'ValueError': ValueError,
    'KeyError': KeyError,
    'AttributeError': AttributeError,
    'IndexError': IndexError,
    'RuntimeError': RuntimeError,
})

def main_wrapper():
    try:
        input_data = json.loads('%s')
        %s
        result = main(**input_data)
        print(json.dumps(result))
    except PermissionError as e:
        print(json.dumps({"error": str(e)}))
        sys.exit(1)
    except Exception as e:
        print(json.dumps({"error": str(e)}))
        sys.exit(1)

if __name__ == "__main__":
    main_wrapper()
`, paramsJson, indentedMainFunc)

	logs.Debug("wrapperCode :%s", wrapperCode)

	cmd := exec.Command("python", "-c", wrapperCode)

	var output strings.Builder
	cmd.Stdout = &output
	cmd.Stderr = &output

	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case err := <-done:
		if err != nil {
			return ``, fmt.Errorf("python execution error: %s, output: %s", err, output.String())
		}
	case <-time.After(10 * time.Second):
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return ``, fmt.Errorf("python execution timeout")
	}

	result := strings.TrimSpace(output.String())
	return result, nil
}
