//go:build linux

package common

import "os/exec"

func PrepareExecCommand(cmd *exec.Cmd, background bool) *exec.Cmd {
	return cmd
}
