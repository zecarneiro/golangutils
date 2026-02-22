package models

type Command struct {
	Cmd                string
	Args               []string
	Cwd                string
	Verbose, IsThrow   bool
	UseShell           bool
	EnvVars            []string
	IsInteractiveShell bool
}
