//go:build !windows
// +build !windows

package golangutils

import (
	"golangutils/entity"
	"os"
)

type UiNotifyUtils struct {
	consoleUtils ConsoleUtils
}

func NewUiNotifyUtils(consoleUtils ConsoleUtils) UiNotifyUtils {
	if consoleUtils == (ConsoleUtils{}) {
		consoleUtils = NewConsole(LoggerUtils{})
	}
	return UiNotifyUtils{consoleUtils: consoleUtils}
}

func (u *UiNotifyUtils) Notify(appId string, title string, message string, icon string) {
	args := []string{"-a", appId}
	if len(icon) > 0 {
		args = append(args, "-i", icon)
	}
	args = append(args, "-t", "60000", title, message)
	u.consoleUtils.ExecRealTime(entity.Command{Cmd: "notify-send", Args: args, Verbose: false, IsThrow: false, EnvVars: os.Environ(), UseBash: false})
}

func (u *UiNotifyUtils) Ok(appId string, message string, icon string) {
	u.Notify(appId, "Success", message, icon)
}

func (u *UiNotifyUtils) Info(appId string, message string, icon string) {
	u.Notify(appId, "Information", message, icon)
}

func (u *UiNotifyUtils) Warn(appId string, message string, icon string) {
	u.Notify(appId, "Warning", message, icon)
}

func (u *UiNotifyUtils) Error(appId string, message string, icon string) {
	u.Notify(appId, "Error", message, icon)
}
