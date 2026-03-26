package shell

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"golangutils/pkg/common"
	"golangutils/pkg/console"
	"golangutils/pkg/enums"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/str"
	"golangutils/pkg/system"
)

func GetShellProfileFile(shellType enums.ShellType) string {
	shells := map[enums.ShellType]string{
		enums.Bash: file.JoinPath(system.HomeDir(), ".bashrc"),
		enums.Zsh:  file.JoinPath(system.HomeDir(), ".zshrc"),
		enums.Fish: file.JoinPath(system.HomeUserConfigDir(), "fish/config.fish"),
		enums.Ksh:  file.JoinPath(system.HomeDir(), ".kshrc"),
	}
	if platform.IsWindows() {
		shells[enums.PowerShell] = file.JoinPath(system.HomeDir(), "Documents/WindowsPowerShell/Microsoft.PowerShell_profile.ps1")
		cmd := exec.Command(GetPowershellCmd(), "-NoProfile", "-Command", "$PROFILE")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err == nil {
			outputShellFile := strings.TrimSpace(out.String())
			if !str.IsEmpty(outputShellFile) {
				shells[enums.PowerShell] = outputShellFile
			}
		}
	} else if platform.IsPlatform([]enums.PlatformType{enums.Linux, enums.Darwin}) {
		shells[enums.PowerShell] = file.JoinPath(system.HomeUserConfigDir(), "powershell/Microsoft.PowerShell_profile.ps1")
	}
	shellFile := shells[shellType]
	return logic.Ternary(str.IsEmpty(shellFile), string(common.Unknown), shellFile)
}

func GetCurrentShellSimple() enums.ShellType {
	if !str.IsEmpty(os.Getenv("BASH_VERSION")) || os.Getenv("BASH_VERSION") == "msys" || !str.IsEmpty(os.Getenv("MSYSTEM")) {
		return enums.Bash
	} else if !str.IsEmpty(os.Getenv("ZSH_VERSION")) {
		return enums.Zsh
	} else if !str.IsEmpty(os.Getenv("KSH_VERSION")) {
		return enums.Ksh
	} else if !str.IsEmpty(os.Getenv("FISH_VERSION")) {
		return enums.Fish
	} else if !str.IsEmpty(os.Getenv("PSModulePath")) {
		return enums.PowerShell
	} else if !str.IsEmpty(os.Getenv("ComSpec")) {
		return enums.Cmd
	}
	return GetCurrentShell()
}

func GetCurrentShell() enums.ShellType {
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
	return enums.UnknownShell
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

func IsShellSimple(shells []enums.ShellType) bool {
	return slices.Contains(shells, GetCurrentShellSimple())
}

func GetShellAllArgsVarStr() string {
	return getShellAllArgsVarStrByShell(GetCurrentShell())
}

func GetShellAllArgsVarStrSimple() string {
	return getShellAllArgsVarStrByShell(GetCurrentShellSimple())
}

/* ----------------------------- POWERSHELL AREA ---------------------------- */
func IsPowerShell() bool {
	return GetCurrentShell().Equals(enums.PowerShell)
}

func IsPowerShellSimple() bool {
	return GetCurrentShellSimple().Equals(enums.PowerShell)
}

func GetPowershellCmd() string {
	shellType := enums.PowerShell
	if existShellCmdFound(shellType) {
		return getShellCmd(shellType)
	}
	cmd, err := console.Which(logic.Ternary(platform.IsWindows(), "powershell.exe", "powershell"))
	if err == nil && cmd != "" {
		updateShellCmdFound(shellType, cmd)
		return cmd
	}
	cmd, err = console.Which(logic.Ternary(platform.IsWindows(), "pwsh.exe", "pwsh"))

	if err == nil && cmd != "" {
		updateShellCmdFound(shellType, cmd)
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

func IsBashSimple() bool {
	return GetCurrentShellSimple().Equals(enums.Bash)
}

func GetBashCmd() string {
	shellType := enums.Bash
	if existShellCmdFound(shellType) {
		return getShellCmd(shellType)
	}
	cmd, err := console.Which(logic.Ternary(platform.IsWindows(), "bash.exe", "bash"))
	if err == nil && cmd != "" {
		updateShellCmdFound(shellType, cmd)
		return cmd
	}
	cmd, err = console.Which(logic.Ternary(platform.IsWindows(), "sh.exe", "sh"))
	if err == nil && cmd != "" {
		updateShellCmdFound(shellType, cmd)
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

func IsFishSimple() bool {
	return GetCurrentShellSimple().Equals(enums.Fish)
}

func GetFishCmd() string {
	shellType := enums.Fish
	if existShellCmdFound(shellType) {
		return getShellCmd(shellType)
	}
	cmd, err := console.Which("fish")
	if err == nil && cmd != "" {
		updateShellCmdFound(shellType, cmd)
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

func IsZshSimple() bool {
	return GetCurrentShellSimple().Equals(enums.Zsh)
}

func GetZshCmd() string {
	shellType := enums.Zsh
	if existShellCmdFound(shellType) {
		return getShellCmd(shellType)
	}
	cmd, err := console.Which("zsh")
	if err == nil && cmd != "" {
		updateShellCmdFound(shellType, cmd)
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

func IsKshSimple() bool {
	return GetCurrentShellSimple().Equals(enums.Ksh)
}

func GetKshCmd() string {
	shellType := enums.Ksh
	if existShellCmdFound(shellType) {
		return getShellCmd(shellType)
	}
	cmd, err := console.Which("ksh")
	if err == nil && cmd != "" {
		updateShellCmdFound(shellType, cmd)
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

func IsCmdSimple() bool {
	return GetCurrentShellSimple().Equals(enums.Cmd)
}

func GetPromptCMDCmd() string {
	shellType := enums.Cmd
	if existShellCmdFound(shellType) {
		return getShellCmd(shellType)
	}
	cmd, err := console.Which("cmd.exe")
	if err == nil && cmd != "" {
		updateShellCmdFound(shellType, cmd)
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
	if (!hasShellType && IsPowerShellSimple()) || shellType == enums.PowerShell {
		return buildPowershellCmd(cmd, args)
	} else if (!hasShellType && IsBashSimple()) || shellType == enums.Bash {
		return buildBashCmd(cmd, args, isInteractive)
	} else if (!hasShellType && IsFishSimple()) || shellType == enums.Fish {
		return buildFishCmd(cmd, args, isInteractive)
	} else if (!hasShellType && IsZshSimple()) || shellType == enums.Zsh {
		return buildZshCmd(cmd, args, isInteractive)
	} else if (!hasShellType && IsKshSimple()) || shellType == enums.Ksh {
		return buildKshCmd(cmd, args, isInteractive)
	} else if (!hasShellType && IsCmdSimple()) || shellType == enums.Cmd {
		return buildPromptCmd(cmd, args)
	}
	return buildDefault(cmd, args)
}

func GetShellCmd(shellType enums.ShellType) string {
	switch shellType {
	case enums.PowerShell:
		return GetPowershellCmd()
	case enums.Bash:
		return GetBashCmd()
	case enums.Zsh:
		return GetZshCmd()
	case enums.Fish:
		return GetFishCmd()
	case enums.Ksh:
		return GetKshCmd()
	case enums.Cmd:
		return GetPromptCMDCmd()
	}
	return logic.Ternary(platform.IsWindows(), GetPowershellCmd(), GetBashCmd())
}
