package console

import (
	"bufio"
	"errors"
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/common/platform"
	"golangutils/pkg/entity"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"os"
	"os/exec"
	"strings"
)

func fillCommand(command *entity.Command) {
	currentDir, _ := file.GetCurrentDir()
	command.EnvVars = common.Ternary(len(command.EnvVars) > 0, command.EnvVars, os.Environ())
	command.Cwd = common.Ternary(command.Cwd == ".", currentDir, command.Cwd)
}

func printCommand(command entity.Command) {
	if command.Verbose {
		logger.Prompt(fmt.Sprintf("%s %s", command.Cmd, strings.Join(command.Args, " ")))
	}
}

func detectShell(command entity.Command) (entity.Command, error) {
	switch platform.GetPlatform() {
	// ───── Linux + macOS → bash or sh ─────
	case platform.Darwin, platform.Linux, platform.Unix:
		cmd := BuildBashCmd(command.Cmd, command.Args)
		command.Cmd = cmd.Cmd
		command.Args = cmd.Args
		return command, nil
	// ───── Windows → PowerShell or CMD ─────
	case platform.Windows:
		// Prefer PowerShell if available
		cmd := BuildPowershellCmd(command.Cmd, command.Args)
		if _, err := Which(cmd.Cmd); err != nil {
			cmd = BuildPromptCmd(command.Cmd, command.Args)
		}
		command.Cmd = cmd.Cmd
		command.Args = cmd.Args
		return command, nil
	}
	return command, errors.New(platform.UnsupportedMSG)
}

func ExecRealTime(command entity.Command) error {
	fillCommand(&command)
	if command.UseShell {
		cmd, err := detectShell(command)
		if err != nil {
			return err
		}
		command = cmd
	}
	printCommand(command)
	cmdResult := exec.Command(command.Cmd, command.Args...)
	cmdResult.Env = command.EnvVars
	cmdResult.Dir = command.Cwd
	cmdResult.Stdout = os.Stdout
	cmdResult.Stderr = os.Stderr
	cmdResult.Stdin = os.Stdin
	return cmdResult.Run()
}

func Exec(command entity.Command) (string, error) {
	printCommand(command)
	fillCommand(&command)
	if command.UseShell {
		cmd, err := detectShell(command)
		if err != nil {
			return "", err
		}
		command = cmd
	}
	cmdResult := exec.Command(command.Cmd, command.Args...)
	cmdResult.Env = command.EnvVars
	cmdResult.Dir = command.Cwd
	output, err := cmdResult.CombinedOutput()
	if len(output) > 0 {
		return string(output), err
	}
	return "", err
}

func WaitForAnyKeyPressed(message string) {
	logger.WithKeepLine(true)
	logger.Log(message)
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func Clear() {
	command := entity.Command{}
	if platform.IsWindows() {
		command.Cmd = "cmd"
		command.Args = []string{"/c", "cls"}
	} else if platform.IsLinux() {
		command.Cmd = "clear"
	}
	cmd := exec.Command(command.Cmd, command.Args...)
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		os.Stdout.WriteString("\x1b[H\x1b[2J")
	}
}

func Chmod777(filepath string) error {
	filepath = file.ResolvePath(filepath)
	fileInfo := entity.FileInfo{}
	if file.IsDir(filepath) {
		info, err := file.ReadDirRecursive(filepath)
		if err != nil {
			return err
		} else {
			fileInfo = info
		}
	} else if file.IsFile(filepath) {
		fileInfo = entity.FileInfo{Files: []string{filepath}}
	} else {
		return fmt.Errorf("%s given file: %s", common.Unknown, filepath)
	}
	if len(fileInfo.Files) > 0 || len(fileInfo.Directories) > 0 {
		logger.Info("Set full permission for '" + filepath + "'")
	}
	var command entity.Command
	if platform.IsWindows() {
		for _, data := range fileInfo.Files {
			command.UseShell = true
			command.Cmd = "Unblock-File"
			command.Args = []string{"-Path", fmt.Sprintf(`"%s"`, &data)}
			ExecRealTime(command)
		}
		for _, data := range fileInfo.Directories {
			command.UseShell = true
			command.Cmd = "Unblock-File"
			command.Args = []string{"-Path", fmt.Sprintf(`"%s"`, &data)}
			ExecRealTime(command)
		}
	} else if platform.IsLinux() {
		for _, data := range fileInfo.Directories {
			err := os.Chmod(data, 0777)
			if err != nil {
				return err
			}
		}
		for _, data := range fileInfo.Files {
			err := os.Chmod(data, 0777)
			if err != nil {
				return err
			}
		}
	}
	return fmt.Errorf(common.NotImplementedYetMSG)
}

func Which(cmd string) (string, error) {
	if cmd == "" {
		return "", nil
	}
	path, err := exec.LookPath(cmd)
	if err != nil {
		return "", err
	}
	return path, nil
}

func WhichByCmds(cmds []string) (string, []error) {
	var errors []error
	for _, cmd := range cmds {
		result, err := Which(cmd)
		if err != nil {
			errors = append(errors, err)
		} else {
			if len(result) > 0 {
				return result, nil
			}
		}

	}
	return "", errors
}

func SetEnv(key string, value string) {
	os.Setenv(key, value)
}

func UnsetEnv(key string) {
	os.Unsetenv(key)
}

func SetEnvBulk(envs []entity.EnvData) {
	for _, data := range envs {
		SetEnv(data.Key, data.Value)
	}
}

func UnsetEnvBulk(envs []entity.EnvData) {
	for _, data := range envs {
		UnsetEnv(data.Key)
	}
}

func HasArgs() bool {
	argsWithoutProg := os.Args[1:]
	return len(argsWithoutProg) > 0
}

func GetArgsList() []string {
	if HasArgs() {
		return os.Args[1:]
	}
	return []string{}
}

func CountArgs() int {
	return len(GetArgsList())
}

func Pause(message string) {
	if len(message) == 0 {
		message = "Press Enter to continue..."
	}
	fmt.Print(message)
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n') // waits for Enter
}
