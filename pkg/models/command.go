package models

import "golangutils/pkg/enums"

type Command struct {
	Cmd                string
	Args               []string
	Cwd                string
	Verbose, IsThrow   bool
	UseShell           bool
	ShellToUse         *enums.ShellType
	EnvVars            []string
	IsInteractiveShell bool
}
