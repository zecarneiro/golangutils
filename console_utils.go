package golangutils

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"golangutils/entity"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

type ConsoleUtils struct {
	loggerUtils *LoggerUtils
	systemUtils SystemUtils
}

func NewConsoleUtils(loggerUtils *LoggerUtils) ConsoleUtils {
	if *loggerUtils == (LoggerUtils{}) {
		*loggerUtils = NewLoggerUtils()
	}
	return ConsoleUtils{loggerUtils: loggerUtils, systemUtils: NewSystemUtils(loggerUtils)}
}

func (c *ConsoleUtils) buildCmd(command entity.Command, useSysProAttr bool) *exec.Cmd {
	var cmd *exec.Cmd
	cmdWithShell := command.GetCmdWithShell(c.GetShellPath(), c.systemUtils.info.Platform)
	if len(cmdWithShell.Args) > 0 {
		cmd = exec.Command(cmdWithShell.Cmd, cmdWithShell.Args...)
	} else {
		cmd = exec.Command(cmdWithShell.Cmd)
	}
	if useSysProAttr {
		c.setSysProAttr(cmd)
	}
	cmd.Dir = command.Cwd
	cmd.Env = command.EnvVars
	if command.Verbose {
		c.loggerUtils.Prompt(cmdWithShell.GetCommandFormated())
	}
	return cmd
}

func (c *ConsoleUtils) Exec(command entity.Command) entity.Response[string] {
	var output []byte
	var err error
	cmd := c.buildCmd(command, true)
	output, err = cmd.CombinedOutput()
	outputStr := string(output[:])
	if command.IsThrow && err != nil {
		log.Fatal(err)
	}
	if command.Verbose {
		fmt.Println(outputStr)
	}
	return entity.Response[string]{Data: outputStr, Error: err}
}

func (c *ConsoleUtils) ExecAsync(command entity.Command, callback func(res entity.Response[string])) {
	cmd := c.buildCmd(command, true)
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		if callback != nil {
			callback(entity.Response[string]{Data: "", Error: err})
		}
		return
	}
	// Start the command.
	if err := cmd.Start(); err != nil {
		if callback != nil {
			callback(entity.Response[string]{Data: "", Error: err})
		}
		return
	}
	// Copy all output from the running command to stdout.
	go func() {
		io.Copy(os.Stdout, stdoutPipe)
		if err := cmd.Wait(); err != nil {
			if callback != nil {
				callback(entity.Response[string]{Data: "", Error: err})
			}
			return
		}
		buf := new(bytes.Buffer)
		buf.ReadFrom(stdoutPipe)
		if callback != nil {
			callback(entity.Response[string]{Data: buf.String(), Error: nil})
		}
	}()
}

func (c *ConsoleUtils) ExecRealTime(command entity.Command) {
	cmd := c.buildCmd(command, false)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	// Start the command
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error starting command:", err)
		return
	}
}

func (c *ConsoleUtils) ExecRealTimeAsync(command entity.Command) {
	cmd := c.buildCmd(command, false)
	_, err := cmd.StdoutPipe()
	if err != nil {
		c.loggerUtils.Error(err.Error())
		return
	}

	// Start the command.
	if err := cmd.Start(); err != nil {
		c.loggerUtils.Error(err.Error())
		return
	}
}

func (c *ConsoleUtils) ExecCommand(command entity.Command) (*os.ProcessState, error) {
	shell := c.GetShellPath()
	if len(shell.Path) == 0 {
		return nil, errors.New("Must provide one of shell type")
	}
	cmd := command.GetCommandFormated()
	if command.Verbose {
		c.loggerUtils.Prompt(cmd)
	}
	var attr = os.ProcAttr{
		Dir:   command.Cwd,
		Env:   command.EnvVars,
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}
	process, err := os.StartProcess(shell.Path, []string{shell.Path, shell.Arg, cmd}, &attr)
	if err == nil {
		return process.Wait()
	}
	return nil, err
}

func (c *ConsoleUtils) ExecCommandAsync(command entity.Command) error {
	shell := c.GetShellPath()
	if len(shell.Path) == 0 {
		return errors.New("Must provide one of shell type")
	}
	cmd := command.GetCommandFormated()
	if command.Verbose {
		c.loggerUtils.Prompt(cmd)
	}
	var attr = os.ProcAttr{
		Dir:   command.Cwd,
		Env:   command.EnvVars,
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}
	process, err := os.StartProcess(shell.Path, []string{shell.Path, shell.Arg, cmd}, &attr)
	if err == nil {
		return process.Release()
	}
	return err
}

func (c *ConsoleUtils) Which(cmd string) []string {
	command := entity.Command{Verbose: false, IsThrow: false}
	if c.systemUtils.IsWindows() {
		command.Cmd = "Get-Command " + cmd + " | Select-Object -ExpandProperty Definition"
		command.UsePowerShell = true
	} else {
		command.Cmd = "which " + cmd
	}
	response := c.Exec(command)
	return strings.Split(response.Data, c.systemUtils.Info().Eol)
}

func (c *ConsoleUtils) Confirm(message string, isNoDefault bool) bool {
	yesNoMsg := "[y/N]"
	if !isNoDefault {
		yesNoMsg = "[Y/n]"
	}
	fmt.Printf("%s %s?: ", message, yesNoMsg)
	var response string
	fmt.Scanln(&response)
	if response == "Y" || response == "y" {
		return true
	}
	return false
}

func (c *ConsoleUtils) ConfirmOrExit(message string, isNoDefault bool) bool {
	yesNoMsg := "[y/N/0(Exit)]"
	if !isNoDefault {
		yesNoMsg = "[Y/n/0(Exit)]"
	}
	fmt.Printf("%s %s?: ", message, yesNoMsg)
	var response string
	fmt.Scanln(&response)
	if response == "0" {
		os.Exit(0)
	}
	if response == "Y" || response == "y" {
		return true
	}
	return false
}

func (c *ConsoleUtils) WaitForAnyKeyPressed(message string) {
	c.loggerUtils.EnableKeepLine()
	c.loggerUtils.Log(message)
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func (c *ConsoleUtils) Clear() {
	command := entity.Command{}
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

func (c *ConsoleUtils) GetShellPath() entity.Shell {
	if c.systemUtils.IsWindows() {
		ps, _ := exec.LookPath("powershell.exe")
		return entity.Shell{Path: ps, Arg: ""}
	} else if c.systemUtils.IsLinux() {
		return entity.Shell{Path: "/bin/bash", Arg: "-c"}
	}
	return entity.Shell{}
}

func SetEnv(key string, value string) {
	os.Setenv(key, value)
}

func UnsetEnv(key string) {
	os.Unsetenv(key)
}
