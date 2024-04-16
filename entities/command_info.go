package entities

type CommandInfo struct {
	Cmd              string
	Args             []string
	Cwd              string
	Verbose, IsThrow bool
	UsePowerShell bool
}
