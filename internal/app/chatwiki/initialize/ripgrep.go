// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package initialize

import (
	"chatwiki/internal/app/chatwiki/define"
	"os/exec"

	"github.com/zhimaAi/go_tools/logs"
)

func initRipgrep() {
	// check ripgrep version
	if output, err := exec.Command(`rg`, `--version`).Output(); err == nil {
		logs.Info("ripgrep version:\r\n%s", string(output))
		return
	}
	// install ripgrep
	if define.Env == `local` {
		return
	}
	output, err := exec.Command(`dpkg`, `-i`, define.AppRoot+`data/ripgrep_14.1.1-1_amd64.deb`).Output()
	if err != nil {
		logs.Error(`install ripgrep error:%s`, err.Error())
		return
	}
	logs.Info("install ripgrep output:\r\n%s", string(output))
	// check ripgrep version
	if output, err = exec.Command(`rg`, `--version`).Output(); err == nil {
		logs.Info("ripgrep version:\r\n%s", string(output))
	} else {
		logs.Error(`get ripgrep version error:%s`, err.Error())
	}
}
