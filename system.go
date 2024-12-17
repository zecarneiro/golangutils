package golangutils

import (
	"errors"
	"os"
)

func Reboot() error {
	logger := NewLogger()
	systemInfo := NewSystemInfo()
	console := NewConsole(logger)
	var cmd Command
	if systemInfo.IsWindows() {
		cmd = Command{
			Cmd:     "shutdown",
			Args:    []string{"/r", "/t", "0"},
			EnvVars: os.Environ(),
		}
	} else if systemInfo.IsLinux() {
		cmd = Command{
			Cmd:     "sudo",
			Args:    []string{"shutdown", "-r", "now"},
			EnvVars: os.Environ(),
		}
	} else if systemInfo.IsDarwin() {
		return errors.New(GetNotImplementedYetMsg())
	}
	if console.Confirm("Will be restart PC. Continue", true) {
		console.ExecRealTime(cmd)
		os.Exit(0)
	}
	return nil
}
