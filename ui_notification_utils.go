//go:build !windows
// +build !windows

package golangutils

import "os"

func Notify(appId string, title string, message string, icon string) {
	args := []string{"-a", appId}
	if len(icon) > 0 {
		args = append(args, "-i", icon)
	}
	args = append(args, "-t", "60000", title, message)
	ExecRealTime(CommandInfo{Cmd: "notify-send", Args: args, Verbose: false, IsThrow: false, EnvVars: os.Environ(), UseBash: false})
}

func OkNotify(appId string, message string, icon string) {
	Notify(appId, "Success", message, icon)
}

func InfoNotify(appId string, message string, icon string) {
	Notify(appId, "Information", message, icon)
}

func WarnNotify(appId string, message string, icon string) {
	Notify(appId, "Warning", message, icon)
}

func ErrorNotify(appId string, message string, icon string) {
	Notify(appId, "Error", message, icon)
}
