package exe

import (
	"fmt"
	"os"
	"strings"

	"golangutils/pkg/enums"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/shell"

	"github.com/google/shlex"
)

func fillCommand(command *models.Command) {
	currentDir, _ := file.GetCurrentDir()
	command.Cwd = logic.Ternary(command.Cwd == ".", currentDir, command.Cwd)
}

func getEnv(command models.Command) []string {
	env := os.Environ()
	env = append(env, "CLICOLOR_FORCE=1", "FORCE_COLOR=1")
	if os.Getenv("TERM") == "" {
		env = append(env, "TERM=xterm-256color")
	}
	return append(env, command.EnvVars...)
}

func printCommand(command models.Command) {
	if command.Verbose {
		logger.Prompt(fmt.Sprintf("%s %s", command.Cmd, strings.Join(command.Args, " ")))
	}
}

func detectShell(command models.Command) models.Command {
	if command.ShellToUse == nil {
		if platform.IsWindows() || platform.IsLinux() {
			cmd := shell.BuildShellCmdByShell(command.Cmd, command.Args, false, logic.Ternary(platform.IsWindows(), enums.PowerShell, enums.Bash))
			command.Cmd = cmd.Cmd
			command.Args = cmd.Args
		} else {
			cmd := shell.BuildShellCmd(command.Cmd, command.Args, command.IsInteractiveShell)
			command.Cmd = cmd.Cmd
			command.Args = cmd.Args
		}
		return command
	}
	cmd := shell.BuildShellCmdByShell(command.Cmd, command.Args, command.IsInteractiveShell, *command.ShellToUse)
	command.Cmd = cmd.Cmd
	command.Args = cmd.Args
	return command
}

func buildNonShellCmd(command models.Command) (models.Command, error) {
	cmdParts, err := shlex.Split(command.Cmd)
	if err != nil {
		return command, err
	}
	command.Cmd = cmdParts[0]
	command.Args = append(cmdParts[1:], command.Args...)
	return command, nil
}
