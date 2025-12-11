package enum

import "strings"

type EShell string

const (
	UNKNOWN_SHELL = "UNKNOWN"
	POWERSHELL    = "pwsh"
	BASH          = "bash"
	ZSH           = "zsh"
	FISH          = "fish"
	KSH           = "ksh"
	CMD           = "cmd"
)

func EShellFromValue(value string) EShell {
	val := strings.ToLower(strings.TrimSpace(value))
	switch val {
	case "pwsh", "powershell", "pwsh.exe", "powershell.exe":
		return POWERSHELL
	case "bash":
		return BASH
	case "zsh":
		return ZSH
	case "fish":
		return FISH
	case "ksh":
		return KSH
	case "cmd", "cmd.exe":
		return CMD
	default:
		return UNKNOWN_SHELL
	}
}

func (s EShell) IsValid() bool {
	switch s {
	case UNKNOWN_SHELL, POWERSHELL, BASH, ZSH, FISH, KSH, CMD:
		return true
	default:
		return false
	}
}

func (s EShell) String() string {
	if s.IsValid() {
		return string(s)
	}
	return UNKNOWN_SHELL
}

func (s EShell) Equals(other EShell) bool {
	return s.String() == other.String()
}

func (s EShell) IsUnknown() bool {
	return s.String() == UNKNOWN_SHELL || !s.IsValid()
}
