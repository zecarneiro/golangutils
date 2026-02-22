package enums

import (
	"golangutils/pkg/common"
	"strings"
)

type ShellType string

const (
	PowerShell   ShellType = "pwsh"
	Bash         ShellType = "bash"
	Zsh          ShellType = "zsh"
	Fish         ShellType = "fish"
	Ksh          ShellType = "ksh"
	Cmd          ShellType = "cmd"
	UnknownShell ShellType = common.Unknown
)

func GetShellTypeFromValue(value string) ShellType {
	val := strings.ToLower(strings.TrimSpace(value))
	switch val {
	case "pwsh", "powershell", "pwsh.exe", "powershell.exe":
		return PowerShell
	case "bash", "bash.exe":
		return Bash
	case "zsh":
		return Zsh
	case "fish":
		return Fish
	case "ksh":
		return Ksh
	case "cmd", "cmd.exe":
		return Cmd
	default:
		return UnknownShell
	}
}

func (s ShellType) IsValid() bool {
	switch s {
	case PowerShell, Bash, Zsh, Fish, Ksh, Cmd:
		return true
	default:
		return false
	}
}

func (s ShellType) String() string {
	if s.IsValid() {
		return string(s)
	}
	return UnknownShell.String()
}

func (s ShellType) Equals(other ShellType) bool {
	return s.String() == other.String()
}
