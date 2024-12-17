package entity

import (
	"fmt"
	"os/exec"
	"strings"
)

type Command struct {
	Cmd              string
	Args             []string
	Cwd              string
	Verbose, IsThrow bool
	UsePowerShell    bool
	UseBash          bool
	EnvVars          []string
}

func (c *Command) copy() Command {
	return Command{
		Cmd:           c.Cmd,
		Args:          c.Args,
		Cwd:           c.Cwd,
		Verbose:       c.Verbose,
		IsThrow:       c.IsThrow,
		UsePowerShell: c.UsePowerShell,
		UseBash:       c.UseBash,
		EnvVars:       c.EnvVars,
	}
}

func (c *Command) getShell(shell string) string {
	shellRes, err := exec.LookPath("powershell.exe")
	if err != nil {
		return shell
	}
	return shellRes
}

func (c *Command) GetCmdWithShell() Command {
	systemInfo := NewSystemInfo()
	cmdStr := fmt.Sprintf("%s %s", c.Cmd, strings.Join(c.Args, " "))
	newCmd := c.copy()
	if systemInfo.IsWindows() {
		if !c.UsePowerShell {
			newCmd.Cmd = c.getShell("cmd.exe")
			newCmd.Args = []string{"/c", cmdStr}
		} else {
			newCmd.Cmd = c.getShell("powershell.exe")
			newCmd.Args = []string{cmdStr}
		}
	} else if c.UseBash {
		newCmd.Cmd = "/bin/bash"
		newCmd.Args = []string{"-c", cmdStr}
	}
	return newCmd
}

func (c *Command) GetCommandFormated() string {
	return fmt.Sprintf("%s %s", c.Cmd, strings.Join(c.Args, " "))
}
