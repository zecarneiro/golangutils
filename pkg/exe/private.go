package exe

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"golangutils/pkg/console"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/shell"
)

func fillCommand(command *models.Command) {
	currentDir, _ := file.GetCurrentDir()
	command.EnvVars = logic.Ternary(len(command.EnvVars) > 0, command.EnvVars, os.Environ())
	command.Cwd = logic.Ternary(command.Cwd == ".", currentDir, command.Cwd)
}

func printCommand(command models.Command) {
	if command.Verbose {
		logger.Prompt(fmt.Sprintf("%s %s", command.Cmd, strings.Join(command.Args, " ")))
	}
}

func detectShell(command models.Command) (models.Command, error) {
	switch platform.GetPlatform() {
	// ───── Linux + macOS → bash or sh ─────
	case platform.Darwin, platform.Linux, platform.Unix:
		cmd := shell.BuildBashCmd(command.Cmd, command.Args)
		command.Cmd = cmd.Cmd
		command.Args = cmd.Args
		return command, nil
	// ───── Windows → PowerShell or CMD ─────
	case platform.Windows:
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
