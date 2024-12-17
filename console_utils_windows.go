//go:build windows
// +build windows

package golangutils

import (
	"os/exec"
	"syscall"
)

func (c *ConsoleUtils) setSysProAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
}
