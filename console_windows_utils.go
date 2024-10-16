//go:build windows
// +build windows

package golangutils

import (
	"os/exec"
	"syscall"
)

/* ----------------------------- END MODEL AREA ----------------------------- */

func setSysProAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
}
