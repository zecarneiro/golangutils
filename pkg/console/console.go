package console

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/str"
)

func Confirm(message string, isNoDefault bool) bool {
	yesNoMsg := "[y/N]"
	if !isNoDefault {
		yesNoMsg = "[Y/n]"
	}
	fmt.Printf("%s %s: ", message, yesNoMsg)
	var response string
	fmt.Scanln(&response)
	response = strings.Trim(response, " ")
	if response == "Y" || response == "y" {
		return true
	} else if len(response) == 0 {
		return logic.Ternary(isNoDefault, false, true)
	}
	return false
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

func Which(cmd string) (string, error) {
	if str.IsEmpty(cmd) {
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

func WaitForAnyKeyPressed(message string) {
	logger.WithKeepLine(true)
	logger.Log(message)
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func Clear() {
	command := models.Command{}
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
