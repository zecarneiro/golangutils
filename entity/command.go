package entity

import (
	"fmt"
	"golangutils/enum"
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

func (c *Command) GetCmdWithShell(shell Shell, platform int) Command {
	cmdStr := fmt.Sprintf("%s %s", c.Cmd, strings.Join(c.Args, " "))
	newCmd := c.copy()
	if platform == enum.WINDOWS {
		if !c.UsePowerShell {
			newCmd.Cmd = "cmd.exe"
			newCmd.Args = []string{"/c", cmdStr}
		} else {
			newCmd.Cmd = shell.Path
			newCmd.Args = []string{cmdStr}
		}
	} else if c.UseBash {
		newCmd.Cmd = shell.Path
		newCmd.Args = []string{shell.Arg, cmdStr}
	}
	return newCmd
}

func (c *Command) GetCommandFormated() string {
	return fmt.Sprintf("%s %s", c.Cmd, strings.Join(c.Args, " "))
}
