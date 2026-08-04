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

func PauseWithMsg(message string) {
	if len(message) == 0 {
		message = "Press Enter to continue..."
	}
	fmt.Print(message)
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n') // waits for Enter
}

func Pause() {
	PauseWithMsg("")
}

func WhichIgnoreError(cmd string) string {
	bin, err := Which(cmd)
	if err != nil {
		return cmd
	}
	return bin
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

func CmdExists(cmd string) bool {
	bin := WhichIgnoreError(cmd)
	return !str.IsEmpty(bin)
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

func ReadUserInput(message string) string {
	reader := bufio.NewReader(os.Stdin)
	if !str.IsEmpty(message) {
		fmt.Printf(`%s: `, message)
	}
	userInput, err := reader.ReadString('\n')
	if err != nil {
		return err.Error()
	}
	return strings.TrimSpace(userInput)
}

func ReadBashUserInput(message string) string {
	bashPath := WhichIgnoreError("bash")
	script := `read -e user_input; echo "$user_input"`
	if str.IsEmpty(bashPath) {
		bashPath = "/bin/bash"
	}
	if !str.IsEmpty(message) {
		fmt.Printf(`%s: `, message)
	}
	cmdResult := exec.Command(bashPath, "-c", script)
	cmdResult.Stdin = os.Stdin
	cmdResult.Stderr = os.Stderr
	output, err := cmdResult.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}
