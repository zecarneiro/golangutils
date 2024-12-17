//go:build windows
// +build windows

package golangutils

import (
	toast "gopkg.in/toast.v1"
)

type UiNotifyWindows struct {
	logger Logger
}

func NewUiNotifyWindows(logger Logger) UiNotifyWindows {
	if logger == (Logger{}) {
		logger = NewLogger()
	}
	return UiNotifyWindows{logger: logger}
}

func (u *UiNotifyWindows) Notify(appId string, title string, message string, icon string) {
	notify := toast.Notification{
		AppID:   appId,
		Title:   title,
		Message: message,
	}
	if len(icon) > 0 {
		notify.Icon = icon
	}
	if err := notify.Push(); err != nil {
		u.logger.Error(err.Error())
	}
}

func (u *UiNotifyWindows) Ok(appId string, message string, icon string) {
	u.Notify(appId, "Success", message, icon)
}

func (u *UiNotifyWindows) Info(appId string, message string, icon string) {
	u.Notify(appId, "Information", message, icon)
}

func (u *UiNotifyWindows) Warn(appId string, message string, icon string) {
	u.Notify(appId, "Warning", message, icon)
}

func (u *UiNotifyWindows) Error(appId string, message string, icon string) {
	u.Notify(appId, "Error", message, icon)
}
