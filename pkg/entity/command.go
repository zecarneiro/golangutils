package entity

type Command struct {
	Cmd              string
	Args             []string
	Cwd              string
	Verbose, IsThrow bool
	UseShell         bool
	EnvVars          []string
}
