package golangutils

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"golangutils/entity"
	"os"
	"os/exec"
	"strings"
)

type Console struct {
	logger     Logger
	systemInfo SystemInfo
}

func NewConsole(logger Logger) Console {
	if logger == (Logger{}) {
		logger = NewLogger()
	}
	return Console{logger: logger, systemInfo: NewSystemInfo()}
}

func (c *Console) buildCmd(command Command, useSysProAttr bool) *exec.Cmd {
	var cmd *exec.Cmd
	cmdWithShell := command.GetCmdWithShell()
	if len(cmdWithShell.Args) > 0 {
		cmd = exec.Command(cmdWithShell.Cmd, cmdWithShell.Args...)
	} else {
		cmd = exec.Command(cmdWithShell.Cmd)
	}
	if useSysProAttr {
		setSysProAttr(cmd)
	}
	cmd.Dir = command.Cwd
	cmd.Env = command.EnvVars
	if command.Verbose {
		c.logger.Prompt(cmdWithShell.GetCommandFormated())
	}
	return cmd
}

func (c *Console) Exec(command Command) Response[string] {
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
	return Response[string]{Data: outputStr, Error: err}
}

func (c *Console) ExecAsync(command Command, callback func(res Response[string])) {
	cmd := c.buildCmd(command, true)
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		if callback != nil {
			callback(Response[string]{Data: "", Error: err})
		}
		return
	}
	// Start the command.
	if err := cmd.Start(); err != nil {
		if callback != nil {
			callback(Response[string]{Data: "", Error: err})
		}
		return
	}
	// Copy all output from the running command to stdout.
	go func() {
		io.Copy(os.Stdout, stdoutPipe)
		if err := cmd.Wait(); err != nil {
			if callback != nil {
				callback(Response[string]{Data: "", Error: err})
			}
			return
		}
		buf := new(bytes.Buffer)
		buf.ReadFrom(stdoutPipe)
		if callback != nil {
			callback(Response[string]{Data: buf.String(), Error: nil})
		}
	}()
}

func (c *Console) ExecRealTime(command entity.Command) {
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

func (c *Console) ExecRealTimeAsync(command Command) {
	cmd := c.buildCmd(command, false)
	_, err := cmd.StdoutPipe()
	if err != nil {
		c.logger.Error(err.Error())
		return
	}

	// Start the command.
	if err := cmd.Start(); err != nil {
		c.logger.Error(err.Error())
		return
	}
}

func (c *Console) Which(cmd string) []string {
	command := Command{Verbose: false, IsThrow: false}
	if c.systemInfo.IsWindows() {
		command.Cmd = "Get-Command " + cmd + " | Select-Object -ExpandProperty Definition"
		command.UsePowerShell = true
	} else {
		command.Cmd = "which " + cmd
	}
	response := c.Exec(command)
	return strings.Split(response.Data, c.systemInfo.GetEol())
}

func (c *Console) Confirm(message string, isNoDefault bool) bool {
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

func (c *Console) ConfirmOrExit(message string, isNoDefault bool) bool {
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

func (c *Console) WaitForAnyKeyPressed(message string) {
	c.logger.EnableKeepLine()
	c.logger.Log(message)
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func (c *Console) ExecCommandAsync(command Command) error {
	var shell, shellArg string
	if command.UseBash {
		shell = "/bin/bash"
		shellArg = "-c"
	} else {
		return errors.New("Must provide one of shell type")
	}
	cmd := command.GetCommandFormated()
	if command.Verbose {
		c.logger.Prompt(cmd)
	}
	var attr = os.ProcAttr{
		Dir:   command.Cwd,
		Env:   command.EnvVars,
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}
	process, err := os.StartProcess(shell, []string{shell, shellArg, cmd}, &attr)
	if err == nil {
		return process.Release()
	}
	return err
}

func (c *Console) ExecCommand(command Command) (*os.ProcessState, error) {
	shell, shellArg := c.GetShellPath()
	if len(shell) == 0 {
		return nil, errors.New("Must provide one of shell type")
	}
	cmd := command.GetCommandFormated()
	if command.Verbose {
		c.logger.Prompt(cmd)
	}
	var attr = os.ProcAttr{
		Dir:   command.Cwd,
		Env:   command.EnvVars,
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}
	process, err := os.StartProcess(shell, []string{shell, shellArg, cmd}, &attr)
	if err == nil {
		return process.Wait()
	}
	return nil, err
}

func (c *Console) Clear() {
	command := Command{}
	if c.systemInfo.IsWindows() {
		command.Cmd = "cmd"
		command.Args = []string{"/c", "cls"}
	} else if c.systemInfo.IsLinux() {
		command.Cmd = "clear"
	}
	cmd := exec.Command(command.Cmd, command.Args...)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (c *Console) GetShellPath() (string, string) {
	if c.systemInfo.IsWindows() {
		ps, _ := exec.LookPath("powershell.exe")
		return ps, ""
	} else if c.systemInfo.IsLinux() {
		return "/bin/bash", "-c"
	}
	return "", ""
}
