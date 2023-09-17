package jnoronhautils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type ICommandInfo struct {
	Cmd              string
	Args             []string
	Cwd              string
	Verbose, IsThrow bool
}

func Exec(commandInfo ICommandInfo) (string, error) {
	var output []byte
	var cmd *exec.Cmd
	var err error
	var command string = fmt.Sprintf("%s %s", commandInfo.Cmd, strings.Join(commandInfo.Args, " "))
	if IsWindows() {
		cmd = exec.Command("powershell", "-command", command)
	} else {
		cmd = exec.Command(command)
	}
	cmd.Dir = commandInfo.Cwd
	if commandInfo.Verbose {
		PromptLog(command)
	}
	output, err = cmd.CombinedOutput()
	outputStr := string(output[:])
	if commandInfo.IsThrow && err != nil {
		log.Fatal(err)
	}
	if commandInfo.Verbose {
		fmt.Println(outputStr)
	}
	return outputStr, err
}

func ExecRealTime(commandInfo ICommandInfo) {
	var cmd *exec.Cmd
	var command string = fmt.Sprintf("%s %s", commandInfo.Cmd, strings.Join(commandInfo.Args, " "))
	if IsWindows() {
		cmd = exec.Command("powershell", "-command", command)
	} else {
		cmd = exec.Command(command)
	}
	cmd.Dir = commandInfo.Cwd
	if commandInfo.Verbose {
		PromptLog(command)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if commandInfo.IsThrow && err != nil {
		log.Fatal(err)
	}
}
