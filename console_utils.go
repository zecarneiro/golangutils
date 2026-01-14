package golangutils

import (
	"bufio"
	"errors"
	"fmt"
	"golangutils/entities"
	"golangutils/enums"
	"os"
	"os/exec"
	"strings"
)

type ConsoleUtils struct {
	loggerUtils *LoggerUtils
	systemUtils *SystemUtils
	shellUtils  *ShellUtils
}

func NewConsoleUtilsDefault() *ConsoleUtils {
	loggerUtils := NewLoggerUtils()
	return &ConsoleUtils{
		loggerUtils: loggerUtils,
		systemUtils: NewSystemUtils(loggerUtils),
		shellUtils:  NewShellUtilsDefault(),
	}
}

func NewConsoleUtils(loggerUtils *LoggerUtils) *ConsoleUtils {
	return &ConsoleUtils{loggerUtils: loggerUtils, systemUtils: NewSystemUtils(loggerUtils)}
}

func (c *ConsoleUtils) fillCommand(command *entities.Command) {
	currentDir, _ := GetCurrentDir()
	command.EnvVars = Ternary(len(command.EnvVars) > 0, command.EnvVars, os.Environ())
	command.Cwd = Ternary(command.Cwd == ".", currentDir, command.Cwd)
}

func (c *ConsoleUtils) printCommand(command entities.Command) {
	if command.Verbose {
		c.loggerUtils.Prompt(fmt.Sprintf("%s %s", command.Cmd, strings.Join(command.Args, " ")))
	}
}

func (c *ConsoleUtils) detectShell(command entities.Command) (entities.Command, error) {
	platform := c.systemUtils.Platform()
	switch platform {
	// ───── Linux + macOS → bash or sh ─────
	case enums.DARWIN, enums.LINUX, enums.UNIX:
		cmd := c.shellUtils.BuildBashCmd(command.Cmd, command.Args)
		command.Cmd = cmd.Cmd
		command.Args = cmd.Args
		return command, nil
	// ───── Windows → PowerShell or CMD ─────
	case enums.WINDOWS:
		// Prefer PowerShell if available
		cmd := c.shellUtils.BuildPowershellCmd(command.Cmd, command.Args)
		if _, err := Which(cmd.Cmd); err != nil {
			cmd = c.shellUtils.BuildPromptCmd(command.Cmd, command.Args)
		}
		command.Cmd = cmd.Cmd
		command.Args = cmd.Args
		return command, nil
	}
	return command, errors.New(GetUnsupportedPlatformMsg())
}

func (c *ConsoleUtils) ExecRealTime(command entities.Command) error {
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

func (c *ConsoleUtils) Exec(command entities.Command) (string, error) {
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
	c.loggerUtils.EnableKeepLine()
	c.loggerUtils.Log(message)
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func (c *ConsoleUtils) Clear() {
	command := entities.Command{}
	if c.systemUtils.IsWindows() {
		command.Cmd = "cmd"
		command.Args = []string{"/c", "cls"}
	} else if c.systemUtils.IsLinux() {
		command.Cmd = "clear"
	}
	cmd := exec.Command(command.Cmd, command.Args...)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (c *ConsoleUtils) Chmod777(file string) {
	fileInfo := entities.FileInfo{}
	if c.systemUtils.IsWindows() {
		if IsDir(file) {
			info, err := ReadDirRecursive(file)
			if err != nil {
				c.loggerUtils.Error(err.Error())
				fileInfo = entities.FileInfo{}
			} else {
				fileInfo = info
			}
		} else if IsFile(file) {
			fileInfo = entities.FileInfo{Files: []string{file}}
		}
	} else if c.systemUtils.IsLinux() {
		fileInfo = entities.FileInfo{Files: []string{file}}
	}
	if len(fileInfo.Files) > 0 {
		c.loggerUtils.Info("Set full permission for '" + file + "'")
	}
	var command entities.Command
	if c.systemUtils.IsWindows() {
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
	} else if c.systemUtils.IsLinux() {
		for _, data := range fileInfo.Directories {
			command.UseShell = true
			command.Cmd = "chmod"
			command.Args = []string{"-R", "777", fmt.Sprintf(`"%s"`, &data)}
			c.ExecRealTime(command)
		}
	} else {
		c.loggerUtils.Error(GetNotImplementedYetMsg())
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
	if response == "Y" || response == "y" {
		return true
	} else if len(response) == 0 {
		return Ternary(isNoDefault, false, true)
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

func SetEnvBulk(envs []entities.EnvData) {
	for _, data := range envs {
		SetEnv(data.Key, data.Value)
	}
}

func UnsetEnvBulk(envs []entities.EnvData) {
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
