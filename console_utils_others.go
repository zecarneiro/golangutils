//go:build !windows
// +build !windows

package golangutils

import "os/exec"

func (c *ConsoleUtils) setSysProAttr(cmd *exec.Cmd) {}
