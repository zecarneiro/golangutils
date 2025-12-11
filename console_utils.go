package golangutils

import (
	"bufio"
	"errors"
	"fmt"
	"golangutils/entity"
	"golangutils/enum"
	"os"
	"os/exec"
	"strings"
)

type ConsoleUtils struct {
	logger *LoggerUtils
	system *SystemUtils
}

func NewConsoleUtilsDefault() *ConsoleUtils {
	logger := NewLoggerUtils()
	return &ConsoleUtils{logger: logger, system: NewSystemUtils(logger)}
}

func NewConsoleUtils(logger *LoggerUtils, system *SystemUtils) *ConsoleUtils {
	return &ConsoleUtils{logger: logger, system: system}
}

func (c *ConsoleUtils) fillCommand(command *entity.Command) {
	currentDir, _ := GetCurrentDir()
	command.EnvVars = Ternary(len(command.EnvVars) > 0, command.EnvVars, os.Environ())
	command.Cwd = Ternary(command.Cwd == ".", currentDir, command.Cwd)
}

func (c *ConsoleUtils) printCommand(command entity.Command) {
	if command.Verbose {
		c.logger.Prompt(fmt.Sprintf("%s %s", command.Cmd, strings.Join(command.Args, " ")))
	}
}

func (c *ConsoleUtils) detectShell(command entity.Command) (entity.Command, error) {
	platform := c.system.Platform()
	shellUtils := NewShellUtils()
	switch platform {
	// ───── Linux + macOS → bash or sh ─────
	case enum.DARWIN, enum.LINUX, enum.UNIX:
		cmd := shellUtils.BuildBashCmd(command.Cmd, command.Args)
		command.Cmd = cmd.Cmd
		command.Args = cmd.Args
		return command, nil
	// ───── Windows → PowerShell or CMD ─────
	case enum.WINDOWS:
		// Prefer PowerShell if available
		cmd := shellUtils.BuildPowershellCmd(command.Cmd, command.Args)
		if _, err := Which(cmd.Cmd); err != nil {
			cmd = shellUtils.BuildPromptCmd(command.Cmd, command.Args)
		}
		command.Cmd = cmd.Cmd
		command.Args = cmd.Args
		return command, nil
	}
	return command, errors.New(GetUnsupportedPlatformMsg())
}

func (c *ConsoleUtils) ExecRealTime(command entity.Command) error {
	c.fillCommand(&command)
	if command.UseShell {
		cmd, err := c.detectShell(command)
		if err != nil {
			return err
		}
		command = cmd
	}
	c.printCommand(command)
	cmdResult := exec.Command(command.Cmd, command.Args...)
	cmdResult.Env = command.EnvVars
	cmdResult.Dir = command.Cwd
	cmdResult.Stdout = os.Stdout
	cmdResult.Stderr = os.Stderr
	cmdResult.Stdin = os.Stdin
	return cmdResult.Run()
}

func (c *ConsoleUtils) Exec(command entity.Command) (string, error) {
	c.printCommand(command)
	c.fillCommand(&command)
	if command.UseShell {
		cmd, err := c.detectShell(command)
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

func (c *ConsoleUtils) WaitForAnyKeyPressed(message string) {
	c.logger.EnableKeepLine()
	c.logger.Log(message)
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func (c *ConsoleUtils) Clear() {
	command := entity.Command{}
	if c.system.IsWindows() {
		command.Cmd = "cmd"
		command.Args = []string{"/c", "cls"}
	} else if c.system.IsLinux() {
		command.Cmd = "clear"
	}
	cmd := exec.Command(command.Cmd, command.Args...)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (c *ConsoleUtils) Chmod777(file string) {
	fileInfo := entity.FileInfo{}
	if c.system.IsWindows() {
		if IsDir(file) {
			info, err := ReadDirRecursive(file)
			if err != nil {
				c.logger.Error(err.Error())
				fileInfo = entity.FileInfo{}
			} else {
				fileInfo = info
			}
		} else if IsFile(file) {
			fileInfo = entity.FileInfo{Files: []string{file}}
		}
	} else if c.system.IsLinux() {
		fileInfo = entity.FileInfo{Files: []string{file}}
	}
	if len(fileInfo.Files) > 0 {
		c.logger.Info("Set full permission for '" + file + "'")
	}
	var command entity.Command
	if c.system.IsWindows() {
		for _, data := range fileInfo.Files {
			command.UseShell = true
			command.Cmd = "Unblock-File"
			command.Args = []string{"-Path", fmt.Sprintf(`"%s"`, &data)}
			c.ExecRealTime(command)
		}
		for _, data := range fileInfo.Directories {
			command.UseShell = true
			command.Cmd = "Unblock-File"
			command.Args = []string{"-Path", fmt.Sprintf(`"%s"`, &data)}
			c.ExecRealTime(command)
		}
	} else if c.system.IsLinux() {
		for _, data := range fileInfo.Directories {
			command.UseShell = true
			command.Cmd = "chmod"
			command.Args = []string{"-R", "777", fmt.Sprintf(`"%s"`, &data)}
			c.ExecRealTime(command)
		}
	} else {
		c.logger.Error(GetNotImplementedYetMsg())
	}
}

func Confirm(message string, isNoDefault bool) bool {
	yesNoMsg := "[y/N]"
	if !isNoDefault {
		yesNoMsg = "[Y/n]"
	}
	fmt.Printf("%s %s?: ", message, yesNoMsg)
	var response string
	fmt.Scanln(&response)
	response = strings.Trim(response, " ")
	if len(response) == 0 || response == "Y" || response == "y" {
		return true
	}
	return false
}

func ConfirmOrExit(message string, isNoDefault bool) bool {
	yesNoMsg := "[y/N/0(Exit)]"
	if !isNoDefault {
		yesNoMsg = "[Y/n/0(Exit)]"
	}
	fmt.Printf("%s %s?: ", message, yesNoMsg)
	var response string
	fmt.Scanln(&response)
	response = strings.Trim(response, " ")
	if response == "0" {
		os.Exit(0)
	} else if len(response) == 0 || response == "Y" || response == "y" {
		return true
	}
	return false
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
