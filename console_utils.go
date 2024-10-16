package golangutils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

/* -------------------------------------------------------------------------- */
/*                                 MODEL AREA                                 */
/* -------------------------------------------------------------------------- */
type CommandInfo struct {
	Cmd              string
	Args             []string
	Cwd              string
	Verbose, IsThrow bool
	UsePowerShell    bool
	UseBash          bool
	EnvVars          []string
}
/* ----------------------------- END MODEL AREA ----------------------------- */

func buildCmd(commandInfo CommandInfo) *exec.Cmd {
	var cmd *exec.Cmd
	command := AddShellCommand(commandInfo)
	if len(command.Args) > 0 {
		cmd = exec.Command(command.Cmd, command.Args...)
	} else {
		cmd = exec.Command(command.Cmd)
	}
	cmd.Dir = commandInfo.Cwd
	cmd.Env = commandInfo.EnvVars
	if commandInfo.Verbose {
		PromptLog(GetCommandToRun(command))
	}
	return cmd
}

func AddShellCommand(commandInfo CommandInfo) CommandInfo {
	ValidateSystem()
	cmdStr := fmt.Sprintf("%s %s", commandInfo.Cmd, strings.Join(commandInfo.Args, " "))
	if IsWindows() {
		if !commandInfo.UsePowerShell {
			commandInfo.Cmd = "cmd.exe"
			commandInfo.Args = []string{"/c", cmdStr}
		} else {
			commandInfo.Cmd = "powershell.exe"
			commandInfo.Args = []string{"-WindowStyle hidden", cmdStr}
		}
	} else if commandInfo.UseBash {
		commandInfo.Cmd = "/bin/bash"
		commandInfo.Args = []string{"-c", cmdStr}
	}
	return commandInfo
}

func GetCommandToRun(commandInfo CommandInfo) string {
	return fmt.Sprintf("%s %s", commandInfo.Cmd, strings.Join(commandInfo.Args, " "))
}

func Exec(commandInfo CommandInfo) Response[string] {
	var output []byte
	var err error
	cmd := buildCmd(commandInfo)
	output, err = cmd.CombinedOutput()
	outputStr := string(output[:])
	if commandInfo.IsThrow && err != nil {
		log.Fatal(err)
	}
	if commandInfo.Verbose {
		fmt.Println(outputStr)
	}
	return Response[string]{Data: outputStr, Error: err}
}

func ExecAsync(commandInfo CommandInfo, callback func(res Response[string])) {
	cmd := buildCmd(commandInfo)
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

func ExecRealTime(commandInfo CommandInfo) {
	cmd := buildCmd(commandInfo)
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

func ExecRealTimeAsync(commandInfo CommandInfo) {
	cmd := buildCmd(commandInfo)
	_, err := cmd.StdoutPipe()
	if err != nil {
		ErrorLog(err.Error(), false)
		return
	}

	// Start the command.
	if err := cmd.Start(); err != nil {
		ErrorLog(err.Error(), false)
		return
	}
}

func Which(cmd string) []string {
	commandInfo := CommandInfo{Verbose: false, IsThrow: false}
	if IsWindows() {
		commandInfo.Cmd = "Get-Command " + cmd + " | Select-Object -ExpandProperty Definition"
		commandInfo.UsePowerShell = true
	} else {
		commandInfo.Cmd = "which " + cmd
	}
	response := Exec(commandInfo)
	return strings.Split(response.Data, SysInfo().Eol)
}

func Confirm(message string, isNoDefault bool) bool {
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

func ConfirmOrExit(message string, isNoDefault bool) bool {
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

func WaitForAnyKeyPressed(message string) {
	LogLog(message, true)
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
