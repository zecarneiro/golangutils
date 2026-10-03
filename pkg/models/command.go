package models

import "golangutils/pkg/enums"

type Command struct {
	Cmd                           string
	Args                          []string
	Cwd                           string
	Verbose, FullVerbose, IsThrow bool
	UseShell                      bool
	ShellToUse                    enums.ShellType
	EnvVars                       []string
	IsInteractiveShell            bool
	IsAsync                       bool // Only work for ExecRealTime
	Background                    bool // Command will run in backgound. If need user interraction or something like that, set false
}
