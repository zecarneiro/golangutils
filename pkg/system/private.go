package system

import (
	"golangutils/pkg/console"
	"golangutils/pkg/logic"
	"golangutils/pkg/platform"
)

func getPwshCmd() string {
	cmdPwsh, err := console.Which(logic.Ternary(platform.IsWindows(), "powershell.exe", "powershell"))
	if err != nil || cmdPwsh == "" {
		cmdPwsh = "powershell"
	}
	return cmdPwsh
}
