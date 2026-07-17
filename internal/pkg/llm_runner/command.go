// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package llm_runner

import (
	"fmt"
	"strings"

	"github.com/zhimaAi/go_tools/tool"
)

var deniedArgsByCommand = map[string][]string{
	"find": {"-delete", "-exec", "-execdir"},
}

func ValidateCommand(command string) error {
	args := strings.Fields(command)
	if len(args) < 2 {
		return nil
	}
	cmd := strings.ToLower(normalizeCommandArg(args[0]))
	denieds := deniedArgsByCommand[cmd]
	if len(denieds) == 0 {
		return nil
	}
	for _, arg := range args[1:] {
		arg = normalizeCommandArg(arg)
		if tool.InArray(arg, denieds) {
			return fmt.Errorf("%s is not allowed for %s", arg, cmd)
		}
	}
	return nil
}

func normalizeCommandArg(arg string) string {
	arg = strings.TrimSpace(arg)
	arg = strings.Trim(arg, `"'`)
	return strings.TrimSpace(arg)
}

// StripToClawbotAnchor converts an absolute workspace path (e.g.
// /workspace/clawbot/working_dir/<robot_key>/<task_id>/foo or
// /home/ubuntu/clawbot/...) into the relative form (clawbot/working_dir/...)
// by stripping everything before the anchor (e.g. "clawbot/") path anchor.
// It returns "" when the path does not anchor under the given anchor, so the
// caller keeps rejecting it as an absolute path. The subsequent scope check
// (allowedPrefix) still enforces that the path stays within the allowed
// directory, so this only relaxes the absolute-path rejection for paths that
// legitimately live under the workspace. Shared by execute-command and
// file-operation validators.
func StripToClawbotAnchor(path string, anchor string) string {
	idx := strings.Index(path, anchor)
	if idx <= 0 || path[idx-1] != '/' {
		return ""
	}
	return path[idx:]
}
