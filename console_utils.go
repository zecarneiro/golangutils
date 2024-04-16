package jnoronha_golangutils

import (
	"bufio"
	"fmt"
	"jnoronha_golangutils/entities"
	"log"
	"os/exec"
	"strings"
)

func addShellCommand(commandInfo entities.CommandInfo) entities.CommandInfo {
	ValidateSystem()
	name := commandInfo.Cmd
	if IsWindows() {
		if !commandInfo.UsePowerShell {
			commandInfo.Cmd = "cmd.exe"
			commandInfo.Args = append([]string{"/c",  name}, commandInfo.Args...)
		} else {
			commandInfo.Cmd = "powershell.exe"
			commandInfo.Args = append([]string{name}, commandInfo.Args...)
		}
		
	} /* else if IsUnix() || IsDarwin() || IsLinux() {
		commandInfo.Cmd = "/bin/sh"
		commandInfo.Args = append([]string{"-c", name}, commandInfo.Args...)
	}*/
	return commandInfo
}

func Exec(commandInfo entities.CommandInfo) entities.Response[string] {
	var output []byte
	var cmd *exec.Cmd
	var err error
	var commandStr string = fmt.Sprintf("%s %s", commandInfo.Cmd, strings.Join(commandInfo.Args, " "))
	command := addShellCommand(commandInfo)
	cmd = exec.Command(command.Cmd, command.Args...)
	cmd.Dir = commandInfo.Cwd
	if commandInfo.Verbose {
		PromptLog(commandStr)
	}
	output, err = cmd.CombinedOutput()
	outputStr := string(output[:])
	if commandInfo.IsThrow && err != nil {
		log.Fatal(err)
	}
	if commandInfo.Verbose {
		fmt.Println(outputStr)
	}
	return entities.Response[string]{Data: outputStr, Error: err}
}

func ExecRealTime(commandInfo entities.CommandInfo) {
	var cmd *exec.Cmd
	var commandStr string = fmt.Sprintf("%s %s", commandInfo.Cmd, strings.Join(commandInfo.Args, " "))
	command := addShellCommand(commandInfo)
	cmd = exec.Command(command.Cmd, command.Args...)
	cmd.Dir = command.Cwd
	if commandInfo.Verbose {
		PromptLog(commandStr)
	}

	// Create a pipe to capture the command's output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating StdoutPipe:", err)
		return
	}

	// Start the command
	err = cmd.Start()
	if err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	// Create a scanner to read from the command's output in real-time
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	// Wait for the command to finish
	err = cmd.Wait()
	if err != nil {
		fmt.Println("Error waiting for command:", err)
	}

	// Close the pipe
	stdout.Close()
}

func Which(cmd string) []string {
	commandInfo := entities.CommandInfo{Verbose: false, IsThrow: false}
	if IsWindows() {
		commandInfo.Cmd = "Get-Command " + cmd + " | Select-Object -ExpandProperty Definition"
		commandInfo.UsePowerShell = true
	} else {
		commandInfo.Cmd = "which " + cmd
	}
	response := Exec(commandInfo)
	return strings.Split(response.Data, SystemInfo().Eol);
} 
