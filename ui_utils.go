package golangutils

import (
	toast "gopkg.in/toast.v1"
)

func Notify(appId string, title string, message string, icon string) {
	notify := toast.Notification{
		AppID:   appId,
		Title:   title,
		Message: message,
	}
	if len(icon) > 0 {
		notify.Icon = icon
	}
	err := notify.Push()
	if err != nil {
		ErrorLog(err.Error(), false)
	}
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
