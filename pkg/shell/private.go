package shell

import (
	"golangutils/pkg/console"
	"golangutils/pkg/models"
)

func isShellCommandExists(cmd string) bool {
	if _, err := console.Which(cmd); err == nil {
		return true
	}
	return false
}

func buildLinuxUnixArgsCmd(cmd string, args []string, isInteractive bool) []string {
	argsCmd := append([]string{"-c", cmd}, args...)
	if isInteractive {
		argsCmd = append([]string{"-i"}, argsCmd...)
	}
	return argsCmd
}

func buildBashCmd(cmd string, args []string, isInteractive bool) models.Command {
	command := models.Command{Args: buildLinuxUnixArgsCmd(cmd, args, isInteractive)}
	if isShellCommandExists("bash") {
		command.Cmd = "bash"
	} else if isShellCommandExists("sh") {
		command.Cmd = "sh"
	} else {
		command.Cmd = cmd
		command.Args = args
	}
	return command
}

func buildFishCmd(cmd string, args []string, isInteractive bool) models.Command {
	command := models.Command{Args: buildLinuxUnixArgsCmd(cmd, args, isInteractive)}
	if isShellCommandExists("fish") {
		command.Cmd = "fish"
	} else {
		command.Cmd = cmd
		command.Args = args
	}
	return command
}

func buildKshCmd(cmd string, args []string, isInteractive bool) models.Command {
	command := models.Command{Args: buildLinuxUnixArgsCmd(cmd, args, isInteractive)}
	if isShellCommandExists("ksh") {
		command.Cmd = "ksh"
	} else {
		command.Cmd = cmd
		command.Args = args
	}
	return command
}

func buildZshCmd(cmd string, args []string, isInteractive bool) models.Command {
	command := models.Command{Args: buildLinuxUnixArgsCmd(cmd, args, isInteractive)}
	if isShellCommandExists("zsh") {
		command.Cmd = "zsh"
	} else {
		command.Cmd = cmd
		command.Args = args
	}
	return command
}
