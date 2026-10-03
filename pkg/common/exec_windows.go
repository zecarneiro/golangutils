//go:build windows

package common

import (
	"os/exec"
	"syscall"
)

func PrepareExecCommand(cmd *exec.Cmd, background bool) *exec.Cmd {
	if background {
		// Config to hide the terminal windows on Windows OS
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow:    true,
			CreationFlags: 0x08000000, // CREATE_NO_WINDOW
		}
	}
	return cmd
}
