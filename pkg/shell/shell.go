package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"golangutils/pkg/common"
	"golangutils/pkg/console"
	"golangutils/pkg/enums"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/system"
)

func GetCurrentShell() enums.ShellType {
	// isWindows := runtime.GOOS == "windows"
	ppid := os.Getppid()

	// 1) Parent process (best-effort)
	if parent, err := system.GetParentProcessInfo(ppid); err == nil {
		shell := enums.GetShellTypeFromValue(parent.Name)
		if shell.IsValid() {
			return shell
		}
	}

	// 2) PowerShell / CMD
	if platform.IsWindows() {
		shell := GetAncestralShell(ppid)
		if shell.IsValid() {
			return shell
		}
		// 1) PowerShell
		if _, ok := os.LookupEnv("PSExecutionPolicyPreference"); ok {
			return enums.PowerShell
		}
		if _, ok := os.LookupEnv("PSModulePath"); ok {
			return enums.PowerShell
		}
		// 2) CMD
		if comspec, ok := os.LookupEnv("ComSpec"); ok && strings.HasSuffix(strings.ToLower(comspec), "cmd.exe") {
			return enums.Cmd
		}

	} else {
		if s, ok := os.LookupEnv("SHELL"); ok {
			shell := enums.GetShellTypeFromValue(filepath.Base(s))
			if shell.IsValid() {
				return shell
			}
		}

		// Linux fallback: /proc/<ppid>/comm
		comm := fmt.Sprintf("/proc/%d/comm", ppid)
		if data, err := os.ReadFile(comm); err == nil {
			shell := enums.GetShellTypeFromValue(strings.TrimSpace(string(data)))
			if shell.IsValid() {
				return shell
			}
		}
	}
	// Unknown
	return common.Unknown
}

func GetAncestralShell(currentPPid int) enums.ShellType {
	var shellList []enums.ShellType = []enums.ShellType{}
	var ancestralProcess models.ParentProcessInfo
	for {
		if currentPPid <= 4 {
			break
		}
		p, err_res := system.GetParentProcessInfo(currentPPid)
		if err_res != nil || p == nil || p.Pid == 0 {
			break
		}
		ancestralProcess = *p
		shell := enums.GetShellTypeFromValue(ancestralProcess.Name)
		if shell.IsValid() {
			shellList = append(shellList, shell)
		}
		if ancestralProcess.PPid == currentPPid {
			break
		}
		currentPPid = ancestralProcess.PPid
	}
	sizeList := len(shellList)
	if sizeList == 0 {
		return common.Unknown
	}
	return shellList[sizeList-1]
}

func IsBash() bool {
	return GetCurrentShell().Equals(enums.Bash)
}

func IsZsh() bool {
	return GetCurrentShell().Equals(enums.Zsh)
}

func IsKsh() bool {
	return GetCurrentShell().Equals(enums.Ksh)
}

func IsFish() bool {
	return GetCurrentShell().Equals(enums.Fish)
}

func IsCmd() bool {
	return GetCurrentShell().Equals(enums.Cmd)
}

func IsPowerShell() bool {
	return GetCurrentShell().Equals(enums.PowerShell)
}

func IsShell(shells []enums.ShellType) bool {
	return slices.Contains(shells, GetCurrentShell())
}

func BuildOthersCmd(cmd string, args []string, isInteractive bool) models.Command {
	if IsFish() {
		return buildFishCmd(cmd, args, isInteractive)
	} else if IsKsh() {
		return buildKshCmd(cmd, args, isInteractive)
	} else if IsZsh() {
		return buildZshCmd(cmd, args, isInteractive)
	}
	return buildBashCmd(cmd, args, isInteractive)
}

func BuildPowershellCmd(cmd string, args []string) models.Command {
	command := models.Command{
		Cmd:  "powershell.exe",
		Args: append([]string{"-nologo", "-Command", cmd}, args...),
	}
	if _, err := console.Which("powershell.exe"); err != nil {
		if _, err := console.Which("pwsh.exe"); err == nil {
			command.Cmd = "pwsh.exe"
		}
	}
	return command
}

func BuildPromptCmd(cmd string, args []string) models.Command {
	command := models.Command{
		Cmd:  "cmd.exe",
		Args: append([]string{"/C", cmd}, args...),
	}
	return command
}
