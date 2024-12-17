//go:build !windows
// +build !windows

package golangutils

import "os/exec"

func setSysProAttr(cmd *exec.Cmd) {}
