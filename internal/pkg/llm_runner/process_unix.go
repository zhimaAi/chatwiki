//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris

package llm_runner

import (
	"errors"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

func configureCommandProcessGroup(cmd *exec.Cmd, cancel func()) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
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
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err == nil {
		err = syscall.Kill(-pgid, syscall.SIGKILL)
	}
	if err == nil || errors.Is(err, os.ErrProcessDone) || errors.Is(err, syscall.ESRCH) {
		return nil
	}
	if killErr := cmd.Process.Kill(); killErr == nil || errors.Is(killErr, os.ErrProcessDone) {
		return nil
	}
	return err
}
