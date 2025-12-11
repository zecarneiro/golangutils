package golangutils

import (
	"fmt"
	"golangutils/entity"
	"golangutils/enum"
	"os"
	"runtime"
	"strings"
)

type ShellUtils struct {
	CurrentShell enum.EShell
}

func NewShellUtils() *ShellUtils {
	shell := ShellUtils{}
	shell.CurrentShell = shell.getCurrentShell()
	return &shell
}

func (s *ShellUtils) getCurrentShell() enum.EShell {
	// 2) Parent process (best-effort)
	ppid := os.Getppid()
	if name, err := GetParentProcess(ppid); err == nil {
		shell := enum.EShellFromValue(name)
		if !shell.IsUnknown() {
			return shell
		}
	}
	// 2) Environment hints (PowerShell / CMD)
	if comspec, ok := os.LookupEnv("ComSpec"); ok && strings.HasSuffix(strings.ToLower(comspec), "cmd.exe") {
		return enum.CMD
	}
	if _, ok := os.LookupEnv("PSModulePath"); ok {
		return enum.POWERSHELL
	}
	if _, ok := os.LookupEnv("PSExecutionPolicyPreference"); ok {
		return enum.POWERSHELL
	}
	// 3) Linux fallback: /proc/<ppid>/comm
	if runtime.GOOS == "linux" {
		comm := fmt.Sprintf("/proc/%d/comm", ppid)
		if data, err := os.ReadFile(comm); err == nil {
			shell := enum.EShellFromValue(strings.TrimSpace(string(data)))
			if !shell.IsUnknown() {
				return shell
			}
		}
	}
	// Unknown
	return enum.UNKNOWN
}

func (s *ShellUtils) IsBash() bool {
	return s.CurrentShell.Equals(enum.BASH)
}

func (s *ShellUtils) IsZsh() bool {
	return s.CurrentShell.Equals(enum.ZSH)
}

func (s *ShellUtils) IsFish() bool {
	return s.CurrentShell.Equals(enum.FISH)
}
func (s *ShellUtils) IsCmd() bool {
	return s.CurrentShell.Equals(enum.CMD)
}

func (s *ShellUtils) IsPowerShell() bool {
	return s.CurrentShell.Equals(enum.POWERSHELL)
}

func (s *ShellUtils) BuildBashCmd(cmd string, args []string) entity.Command {
	command := entity.Command{
		Cmd:  "sh",
		Args: append([]string{"-c", cmd}, args...),
	}
	if _, err := Which("bash"); err == nil {
		command.Cmd = "bash"
	}
	return command
}

func (s *ShellUtils) BuildPowershellCmd(cmd string, args []string) entity.Command {
	command := entity.Command{
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

func (s *ShellUtils) BuildPromptCmd(cmd string, args []string) entity.Command {
	command := entity.Command{
		Cmd:  "cmd.exe",
		Args: append([]string{"/C", cmd}, args...),
	}
	return command
}
