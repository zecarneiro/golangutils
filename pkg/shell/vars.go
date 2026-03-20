package shell

import "golangutils/pkg/enums"

const (
	PowershellAllArgsVarStr = "$args"
	CmdAllArgsVarStr        = "%*"
	BashAllArgsVarStr       = "$@"
	ZshAllArgsVarStr        = BashAllArgsVarStr
	KshAllArgsVarStr        = BashAllArgsVarStr
	FishAllArgsVarStr       = "$argv"
	UnsupportedMSG          = "Unsupported shell"
)

var shellCmdFound map[enums.ShellType]string
