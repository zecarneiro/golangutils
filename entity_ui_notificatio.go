//go:build !windows
// +build !windows

package golangutils

import "os"

type UiNotify struct {
	console Console
}

func NewUiNotify(console Console) UiNotify {
	if console == (Console{}) {
		console = NewConsole(Logger{})
	}
	return UiNotify{console: console}
}

func (u *UiNotify) Notify(appId string, title string, message string, icon string) {
	args := []string{"-a", appId}
	if len(icon) > 0 {
		args = append(args, "-i", icon)
	}
	args = append(args, "-t", "60000", title, message)
	u.console.ExecRealTime(Command{Cmd: "notify-send", Args: args, Verbose: false, IsThrow: false, EnvVars: os.Environ(), UseBash: false})
}

func (u *UiNotify) Ok(appId string, message string, icon string) {
	u.Notify(appId, "Success", message, icon)
}

func (u *UiNotify) Info(appId string, message string, icon string) {
	u.Notify(appId, "Information", message, icon)
}

func (u *UiNotify) Warn(appId string, message string, icon string) {
	u.Notify(appId, "Warning", message, icon)
}

func (u *UiNotify) Error(appId string, message string, icon string) {
	u.Notify(appId, "Error", message, icon)
}
