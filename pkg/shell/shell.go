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

func IsShell(shells []enums.ShellType) bool {
	return slices.Contains(shells, GetCurrentShell())
}

/* ----------------------------- POWERSHELL AREA ---------------------------- */
func IsPowerShell() bool {
	return GetCurrentShell().Equals(enums.PowerShell)
}

func GetPowershellCmd() string {
	cmd, err := console.Which("powershell.exe")
	if err == nil && cmd != "" {
		return cmd
	}
	cmd, err = console.Which("pwsh.exe")
	if err == nil && cmd != "" {
		return cmd
	}
	return ""
}

func IsPowershellInstalled() bool {
	return GetPowershellCmd() != ""
}

/* -------------------------------- BASH AREA ------------------------------- */
func IsBash() bool {
	return GetCurrentShell().Equals(enums.Bash)
}

func GetBashCmd() string {
	cmd, err := console.Which("bash")
	if err == nil && cmd != "" {
		return cmd
	}
	cmd, err = console.Which("sh")
	if err == nil && cmd != "" {
		return cmd
	}
	return ""
}

func IsBashInstalled() bool {
	return GetBashCmd() != ""
}

/* -------------------------------- FISH AREA ------------------------------- */
func IsFish() bool {
	return GetCurrentShell().Equals(enums.Fish)
}

func GetFishCmd() string {
	cmd, err := console.Which("fish")
	if err == nil && cmd != "" {
		return cmd
	}
	return ""
}

func IsFishInstalled() bool {
	return GetFishCmd() != ""
}

/* -------------------------------- ZSH AREA -------------------------------- */
func IsZsh() bool {
	return GetCurrentShell().Equals(enums.Zsh)
}

func GetZshCmd() string {
	cmd, err := console.Which("zsh")
	if err == nil && cmd != "" {
		return cmd
	}
	return ""
}

func IsZshInstalled() bool {
	return GetZshCmd() != ""
}

/* -------------------------------- KSH AREA -------------------------------- */
func IsKsh() bool {
	return GetCurrentShell().Equals(enums.Ksh)
}

func GetKshCmd() string {
	cmd, err := console.Which("ksh")
	if err == nil && cmd != "" {
		return cmd
	}
	return ""
}

func IsKshInstalled() bool {
	return GetKshCmd() != ""
}

/* -------------------------------- CMD AREA -------------------------------- */
func IsCmd() bool {
	return GetCurrentShell().Equals(enums.Cmd)
}

func GetPromptCMDCmd() string {
	cmd, err := console.Which("cmd.exe")
	if err == nil && cmd != "" {
		return cmd
	}
	return ""
}

func IsPromptCMDInstalled() bool {
	return GetPromptCMDCmd() != ""
}

/* ---------------------------- BUILD OTHERS AREA --------------------------- */
func BuildShellCmd(cmd string, args []string, isInteractive bool) models.Command {
	return BuildShellCmdByShell(cmd, args, isInteractive, enums.UnknownShell)
}

func BuildShellCmdByShell(cmd string, args []string, isInteractive bool, shellType enums.ShellType) models.Command {
	hasShellType := false
	if shellType != enums.UnknownShell {
		hasShellType = true
	}
	if (!hasShellType && IsPowerShell()) || shellType == enums.PowerShell {
		return buildPowershellCmd(cmd, args)
	} else if (!hasShellType && IsBash()) || shellType == enums.Bash {
		return buildBashCmd(cmd, args, isInteractive)
	} else if (!hasShellType && IsFish()) || shellType == enums.Fish {
		return buildFishCmd(cmd, args, isInteractive)
	} else if (!hasShellType && IsZsh()) || shellType == enums.Zsh {
		return buildZshCmd(cmd, args, isInteractive)
	} else if (!hasShellType && IsKsh()) || shellType == enums.Ksh {
		return buildKshCmd(cmd, args, isInteractive)
	} else if (!hasShellType && IsCmd()) || shellType == enums.Cmd {
		return buildPromptCmd(cmd, args)
	}
	return buildDefault(cmd, args)
}
