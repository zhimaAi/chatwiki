//go:build linux

package initialize

import "golang.org/x/sys/unix"

func protectController() error {
	// The untrusted child shares our uid. Prevent /proc/<pid>/mem access
	// to the controller even on hosts whose ptrace_scope is permissive.
	return unix.Prctl(unix.PR_SET_DUMPABLE, 0, 0, 0, 0)
}
