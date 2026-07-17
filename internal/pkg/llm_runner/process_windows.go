//go:build windows

package llm_runner

import (
	"errors"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

func configureCommandProcessGroup(cmd *exec.Cmd, cancel func()) {
	var once sync.Once
	cmd.Cancel = func() error {
		once.Do(cancel)
		return terminateCommandProcessGroup(cmd)
	}
}

func terminateCommandProcessGroup(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	err := exec.Command(`taskkill`, `/PID`, strconv.Itoa(cmd.Process.Pid), `/T`, `/F`).Run()
	if err == nil {
		return nil
	}
	if killErr := cmd.Process.Kill(); killErr == nil || errors.Is(killErr, os.ErrProcessDone) {
		return nil
	}
	return err
}
