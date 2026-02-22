package exe

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"golangutils/pkg/console"
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

func detectShell(command models.Command) (models.Command, error) {
	switch platform.GetPlatform() {
	// ───── Linux + macOS → bash or sh ─────
	case enums.Darwin, enums.Linux, enums.Unix:
		cmd := shell.BuildOthersCmd(command.Cmd, command.Args, command.IsInteractiveShell)
		command.Cmd = cmd.Cmd
		command.Args = cmd.Args
		return command, nil
	// ───── Windows → PowerShell or CMD ─────
	case enums.Windows:
		// Prefer PowerShell if available
		cmd := shell.BuildPowershellCmd(command.Cmd, command.Args)
		if _, err := console.Which(cmd.Cmd); err != nil {
			cmd = shell.BuildPromptCmd(command.Cmd, command.Args)
		}
		command.Cmd = cmd.Cmd
		command.Args = cmd.Args
		return command, nil
	}
	return command, errors.New(platform.UnsupportedMSG)
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
