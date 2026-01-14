package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"golangutils/pkg/common"
	"golangutils/pkg/console"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/system"
)

func GetCurrentShell() ShellType {
	// isWindows := runtime.GOOS == "windows"
	ppid := os.Getppid()

	// 1) Parent process (best-effort)
	if parent, err := system.GetParentProcessInfo(ppid); err == nil {
		shell := GetShellTypeFromValue(parent.Name)
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
			return PowerShell
		}
		if _, ok := os.LookupEnv("PSModulePath"); ok {
			return PowerShell
		}
		// 2) CMD
		if comspec, ok := os.LookupEnv("ComSpec"); ok && strings.HasSuffix(strings.ToLower(comspec), "cmd.exe") {
			return Cmd
		}

	} else {
		if s, ok := os.LookupEnv("SHELL"); ok {
			shell := GetShellTypeFromValue(filepath.Base(s))
			if shell.IsValid() {
				return shell
			}
		}

		// Linux fallback: /proc/<ppid>/comm
		comm := fmt.Sprintf("/proc/%d/comm", ppid)
		if data, err := os.ReadFile(comm); err == nil {
			shell := GetShellTypeFromValue(strings.TrimSpace(string(data)))
			if shell.IsValid() {
				return shell
			}
		}
	}
	// Unknown
	return common.Unknown
}

func GetAncestralShell(currentPPid int) ShellType {
	var shellList []ShellType = []ShellType{}
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
		shell := GetShellTypeFromValue(ancestralProcess.Name)
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
	return GetCurrentShell().Equals(Bash)
}

func IsZsh() bool {
	return GetCurrentShell().Equals(Zsh)
}

func IsFish() bool {
	return GetCurrentShell().Equals(Fish)
}

func IsCmd() bool {
	return GetCurrentShell().Equals(Cmd)
}

func IsPowerShell() bool {
	return GetCurrentShell().Equals(PowerShell)
}

func IsShell(shells []ShellType) bool {
	return slices.Contains(shells, GetCurrentShell())
}

func BuildBashCmd(cmd string, args []string) models.Command {
	command := models.Command{
		Cmd:  "sh",
		Args: append([]string{"-c", cmd}, args...),
	}
	if _, err := console.Which("bash"); err == nil {
		command.Cmd = "bash"
	}
	return command
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
