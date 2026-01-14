package golangutils

import (
	"fmt"
	"golangutils/entities"
	"golangutils/enums"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type ShellUtils struct {
	systemUtils *SystemUtils
}

func NewShellUtilsDefault() *ShellUtils {
	return &ShellUtils{systemUtils: NewSystemUtilsDefault()}
}

func (s *ShellUtils) GetCurrentShell() enums.EShell {
	isWindows := runtime.GOOS == "windows"
	ppid := os.Getppid()

	// 1) Parent process (best-effort)
	if parent, err := s.systemUtils.GetParentProcessInfo(ppid); err == nil {
		shell := enums.EShellFromValue(parent.Name)
		if !shell.IsUnknown() {
			return shell
		}
	}

	// 2) PowerShell / CMD
	if isWindows {
		shell := s.GetAncestralShell(ppid)
		if !shell.IsUnknown() {
			return shell
		}
		// 1) PowerShell
		if _, ok := os.LookupEnv("PSExecutionPolicyPreference"); ok {
			return enums.POWERSHELL
		}
		if _, ok := os.LookupEnv("PSModulePath"); ok {
			return enums.POWERSHELL
		}
		// 2) CMD
		if comspec, ok := os.LookupEnv("ComSpec"); ok && strings.HasSuffix(strings.ToLower(comspec), "cmd.exe") {
			return enums.CMD
		}

	} else {
		if s, ok := os.LookupEnv("SHELL"); ok {
			shell := enums.EShellFromValue(filepath.Base(s))
			if !shell.IsUnknown() {
				return shell
			}
		}

		// Linux fallback: /proc/<ppid>/comm
		comm := fmt.Sprintf("/proc/%d/comm", ppid)
		if data, err := os.ReadFile(comm); err == nil {
			shell := enums.EShellFromValue(strings.TrimSpace(string(data)))
			if !shell.IsUnknown() {
				return shell
			}
		}
	}
	// Unknown
	return enums.UNKNOWN
}

func (s *ShellUtils) GetAncestralShell(currentPPid int) enums.EShell {
	var shellList []enums.EShell = []enums.EShell{}
	var ancestralProcess entities.ParentProcessInfo
	for {
		if currentPPid <= 4 {
			break
		}
		p, err_res := s.systemUtils.GetParentProcessInfo(currentPPid)
		if err_res != nil || p == nil || p.Pid == 0 {
			break
		}
		ancestralProcess = *p
		shell := enums.EShellFromValue(ancestralProcess.Name)
		if !shell.IsUnknown() {
			shellList = append(shellList, shell)
		}
		if ancestralProcess.PPid == currentPPid {
			break
		}
		currentPPid = ancestralProcess.PPid
	}
	sizeList := len(shellList)
	if sizeList == 0 {
		return enums.UNKNOWN
	}
	return shellList[sizeList-1]
}

func (s *ShellUtils) IsBash() bool {
	return s.GetCurrentShell().Equals(enums.BASH)
}

func (s *ShellUtils) IsZsh() bool {
	return s.GetCurrentShell().Equals(enums.ZSH)
}

func (s *ShellUtils) IsFish() bool {
	return s.GetCurrentShell().Equals(enums.FISH)
}
func (s *ShellUtils) IsCmd() bool {
	return s.GetCurrentShell().Equals(enums.CMD)
}

func (s *ShellUtils) IsPowerShell() bool {
	return s.GetCurrentShell().Equals(enums.POWERSHELL)
}

func (s *ShellUtils) IsShell(shells []enums.EShell) bool {
	return InArray(shells, s.GetCurrentShell())
}

func (s *ShellUtils) BuildBashCmd(cmd string, args []string) entities.Command {
	command := entities.Command{
		Cmd:  "sh",
		Args: append([]string{"-c", cmd}, args...),
	}
	if _, err := Which("bash"); err == nil {
		command.Cmd = "bash"
	}
	return command
}

func (s *ShellUtils) BuildPowershellCmd(cmd string, args []string) entities.Command {
	command := entities.Command{
		Cmd:  "powershell.exe",
		Args: append([]string{"-nologo", "-Command", cmd}, args...),
	}
	if _, err := Which("powershell.exe"); err != nil {
		if _, err := Which("pwsh.exe"); err == nil {
			command.Cmd = "pwsh.exe"
		}
	}
	return command
}

func (s *ShellUtils) BuildPromptCmd(cmd string, args []string) entities.Command {
	command := entities.Command{
		Cmd:  "cmd.exe",
		Args: append([]string{"/C", cmd}, args...),
	}
	return command
}
