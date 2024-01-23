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
		commandInfo.Cmd = "cmd.exe"
		commandInfo.Args = append([]string{"/c", name}, commandInfo.Args...)
	} else if IsUnix() || IsDarwin() || IsLinux() {
		commandInfo.Cmd = "/bin/sh"
		commandInfo.Args = append([]string{"-c", name}, commandInfo.Args...)
	}
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
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}

	err = cmd.Start()
	fmt.Println("The command is running")
	if err != nil {
		fmt.Println(err)
	}

	// print the output of the subprocess
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()
}
