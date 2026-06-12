// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package llm_runner

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"
)

var allowedCommands = map[string]struct{}{
	"cat":     {},
	"find":    {},
	"grep":    {},
	"head":    {},
	"jq":      {},
	"ls":      {},
	"node":    {},
	"npm":     {},
	"pwd":     {},
	"python":  {},
	"python3": {},
	"rg":      {},
	"tail":    {},
	"wc":      {},
}

var deniedArgsByCommand = map[string][]string{
	"find":    {"-delete", "-exec", "-execdir"},
	"node":    {"-e", "--eval", "-p", "--print"},
	"python":  {"-c", "-m"},
	"python3": {"-c", "-m"},
}

// PrepareCommand parses a command string into argv and applies runner-side
// checks that do not depend on a specific robot.
func PrepareCommand(command string) ([]string, error) {
	args, err := SplitCommand(command)
	if err != nil {
		return nil, err
	}
	if err = ValidateCommandArgs(args); err != nil {
		return nil, err
	}
	return args, nil
}

// SplitCommand is a small shlex-style parser. It supports whitespace
// separation, single quotes, double quotes, and backslash escaping, while
// rejecting shell control syntax outside quoted strings.
func SplitCommand(command string) ([]string, error) {
	command = strings.TrimSpace(command)
	if command == "" {
		return nil, fmt.Errorf("command is required")
	}
	if containsWindowsDrivePath(command) {
		return nil, fmt.Errorf("absolute path is not allowed")
	}

	var args []string
	var cur strings.Builder
	var quote rune
	var inToken bool
	rs := []rune(command)
	for i := 0; i < len(rs); i++ {
		ch := rs[i]
		if quote == 0 && isShellControlRune(ch) {
			return nil, fmt.Errorf("shell control syntax is not allowed: %q", ch)
		}
		if quote == 0 && unicode.IsSpace(ch) {
			if inToken {
				args = append(args, cur.String())
				cur.Reset()
				inToken = false
			}
			continue
		}
		if ch == '\'' || ch == '"' {
			if quote == 0 {
				quote = ch
				inToken = true
				continue
			}
			if quote == ch {
				quote = 0
				inToken = true
				continue
			}
		}
		if ch == '\\' && i+1 < len(rs) {
			i++
			cur.WriteRune(rs[i])
			inToken = true
			continue
		}
		cur.WriteRune(ch)
		inToken = true
	}
	if quote != 0 {
		return nil, fmt.Errorf("unterminated quoted string")
	}
	if inToken {
		args = append(args, cur.String())
	}
	if len(args) == 0 {
		return nil, fmt.Errorf("command is required")
	}
	return args, nil
}

func ValidateCommandArgs(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("command is required")
	}
	cmd := commandName(args[0])
	if _, ok := allowedCommands[cmd]; !ok {
		return fmt.Errorf("command is not allowed: %s", args[0])
	}
	for _, denied := range deniedArgsByCommand[cmd] {
		for _, arg := range args[1:] {
			if arg == denied || strings.HasPrefix(arg, denied+"=") {
				return fmt.Errorf("%s is not allowed for %s", denied, cmd)
			}
		}
	}
	if err := validatePathCandidates(args, "clawbot/working_dir"); err != nil {
		return err
	}
	return nil
}

func ValidateClawbotPathScope(args []string, allowedPrefix string) error {
	return validatePathCandidates(args, allowedPrefix)
}

func validatePathCandidates(args []string, allowedPrefix string) error {
	allowedPrefix = normalizePathToken(allowedPrefix)
	for _, arg := range args {
		for _, candidate := range pathCandidates(arg) {
			candidate = normalizePathToken(candidate)
			if candidate == "" || looksLikeURL(candidate) {
				continue
			}
			if isAbsolutePath(candidate) {
				return fmt.Errorf("absolute path is not allowed: %s", candidate)
			}
			if hasParentTraversal(candidate) {
				return fmt.Errorf("parent-directory traversal is not allowed: %s", candidate)
			}
			if strings.HasPrefix(candidate, "clawbot/") && candidate != allowedPrefix && !strings.HasPrefix(candidate, allowedPrefix+"/") {
				return fmt.Errorf("%s is not allowed in command", candidate)
			}
		}
	}
	return nil
}

func pathCandidates(arg string) []string {
	candidates := []string{arg}
	for _, sep := range []string{"=", ":"} {
		if idx := strings.Index(arg, sep); idx >= 0 && idx+1 < len(arg) {
			candidates = append(candidates, arg[idx+1:])
		}
	}
	return candidates
}

func commandName(cmd string) string {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return ""
	}
	if strings.ContainsAny(cmd, `/\`) {
		return ""
	}
	return strings.ToLower(filepath.Base(cmd))
}

func normalizePathToken(token string) string {
	token = strings.TrimSpace(token)
	token = strings.Trim(token, `"'`)
	token = strings.ReplaceAll(token, `\`, `/`)
	return strings.TrimRight(token, `/`)
}

func hasParentTraversal(path string) bool {
	for _, part := range strings.Split(path, "/") {
		if part == ".." {
			return true
		}
	}
	return false
}

func isAbsolutePath(path string) bool {
	if path == "" {
		return false
	}
	if strings.HasPrefix(path, "/") || strings.HasPrefix(path, "~") {
		return true
	}
	if len(path) >= 3 && unicode.IsLetter(rune(path[0])) && path[1] == ':' && path[2] == '/' {
		return true
	}
	return runtime.GOOS == "windows" && strings.HasPrefix(path, "//")
}

func looksLikeURL(path string) bool {
	idx := strings.Index(path, "://")
	if idx <= 0 {
		return false
	}
	for _, ch := range path[:idx] {
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '+' && ch != '-' && ch != '.' {
			return false
		}
	}
	return true
}

func containsWindowsDrivePath(s string) bool {
	rs := []rune(s)
	for i := 0; i+2 < len(rs); i++ {
		if unicode.IsLetter(rs[i]) && rs[i+1] == ':' && (rs[i+2] == '\\' || rs[i+2] == '/') {
			return true
		}
	}
	return false
}

func isShellControlRune(ch rune) bool {
	switch ch {
	case '|', ';', '&', '<', '>', '`', '$', '\n', '\r':
		return true
	default:
		return false
	}
}
