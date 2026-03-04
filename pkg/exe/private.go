package exe

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golangutils/pkg/enums"
	"golangutils/pkg/env"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/shell"
	"golangutils/pkg/str"

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
	if command.Verbose || command.FullVerbose {
		logger.Prompt(fmt.Sprintf("%s %s", command.Cmd, strings.Join(command.Args, " ")))
	}
}

func detectShell(command models.Command) models.Command {
	if !command.ShellToUse.IsValid() {
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
	cmd := shell.BuildShellCmdByShell(command.Cmd, command.Args, command.IsInteractiveShell, command.ShellToUse)
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

func setFullAccessPowerShell(filepath string) error {
	username := env.Get("username")
	if str.IsEmpty(username) {
		username = env.Get("USER")
	}
	script := fmt.Sprintf(`
$acl = Get-Acl -Path "%s"
$rule = [security.accesscontrol.filesystemaccessrule]::new("%s", "FullControl", "Allow")
$acl.AddAccessRule($rule)
$acl | Set-Acl
`, filepath, username)
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-Command", script)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
