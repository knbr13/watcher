//go:build !windows

package main

import (
	"os/exec"
	"syscall"
)

// setProcAttrs puts the command in its own process group so killProcessTree
// can terminate it along with any children it spawned.
func setProcAttrs(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func killProcessTree(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return nil
	}
	return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
}
