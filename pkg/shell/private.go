package shell

import (
	"golangutils/pkg/enums"
	"golangutils/pkg/models"
	"golangutils/pkg/slice"
)

func buildDefault(cmd string, args []string) models.Command {
	return models.Command{
		Cmd:  cmd,
		Args: args,
	}
}

func buildLinuxUnixArgsCmd(cmd string, args []string, isInteractive bool) []string {
	argsCmd := append([]string{"-c", cmd}, args...)
	if isInteractive {
		argsCmd = append([]string{"-i"}, argsCmd...)
	}
	return argsCmd
}

func buildPowershellCmd(cmd string, args []string) models.Command {
	if IsPowershellInstalled() {
		return models.Command{
			Cmd:  GetPowershellCmd(),
			Args: append([]string{"-nologo", "-Command", cmd}, args...),
		}
	}
	return buildDefault(cmd, args)
}

func buildBashCmd(cmd string, args []string, isInteractive bool) models.Command {
	if IsBashInstalled() {
		return models.Command{
			Cmd:  GetBashCmd(),
			Args: buildLinuxUnixArgsCmd(cmd, args, isInteractive),
		}
	}
	return buildDefault(cmd, args)
}

func buildFishCmd(cmd string, args []string, isInteractive bool) models.Command {
	if IsFishInstalled() {
		return models.Command{
			Cmd:  GetFishCmd(),
			Args: buildLinuxUnixArgsCmd(cmd, args, isInteractive),
		}
	}
	return buildDefault(cmd, args)
}

func buildZshCmd(cmd string, args []string, isInteractive bool) models.Command {
	if IsZshInstalled() {
		return models.Command{
			Cmd:  GetZshCmd(),
			Args: buildLinuxUnixArgsCmd(cmd, args, isInteractive),
		}
	}
	return buildDefault(cmd, args)
}

func buildKshCmd(cmd string, args []string, isInteractive bool) models.Command {
	if IsKshInstalled() {
		return models.Command{
			Cmd:  GetKshCmd(),
			Args: buildLinuxUnixArgsCmd(cmd, args, isInteractive),
		}
	}
	return buildDefault(cmd, args)
}

func buildPromptCmd(cmd string, args []string) models.Command {
	if IsPromptCMDInstalled() {
		return models.Command{
			Cmd:  GetPromptCMDCmd(),
			Args: append([]string{"/C", cmd}, args...),
		}
	}
	return buildDefault(cmd, args)
}

func getShellAllArgsVarStrByShell(shellType enums.ShellType) string {
	switch shellType {
	case enums.Bash, enums.Zsh, enums.Ksh:
		return BashAllArgsVarStr
	case enums.Fish:
		return FishAllArgsVarStr
	case enums.PowerShell:
		return PowershellAllArgsVarStr
	case enums.Cmd:
		return CmdAllArgsVarStr
	}
	return ""
}

func updateShellCmdFound(shellType enums.ShellType, cmd string) {
	if shellCmdFound == nil {
		shellCmdFound = make(map[enums.ShellType]string)
	}
	shellCmdFound[shellType] = cmd
}

func existShellCmdFound(shellType enums.ShellType) bool {
	return slice.MapExistKey(shellCmdFound, shellType)
}

func getShellCmd(shellType enums.ShellType) string {
	if !existShellCmdFound(shellType) {
		return ""
	}
	return shellCmdFound[shellType]
}
