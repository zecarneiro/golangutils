//go:build windows
// +build windows

package golangutils

import (
	toast "gopkg.in/toast.v1"
)

type UiNotifyUtilsWindows struct {
	loggerUtils LoggerUtils
}

func NewUiNotifyWindows(loggerUtils LoggerUtils) UiNotifyUtilsWindows {
	if loggerUtils == (LoggerUtils{}) {
		loggerUtils = NewLoggerUtils()
	}
	return UiNotifyUtilsWindows{loggerUtils: loggerUtils}
}

func (u *UiNotifyUtilsWindows) Notify(appId string, title string, message string, icon string) {
	notify := toast.Notification{
		AppID:   appId,
		Title:   title,
		Message: message,
	}
	if len(icon) > 0 {
		notify.Icon = icon
	}
	if err := notify.Push(); err != nil {
		u.loggerUtils.Error(err.Error())
	}
}

func (u *UiNotifyUtilsWindows) Ok(appId string, message string, icon string) {
	u.Notify(appId, "Success", message, icon)
}

func (u *UiNotifyUtilsWindows) Info(appId string, message string, icon string) {
	u.Notify(appId, "Information", message, icon)
}

func (u *UiNotifyUtilsWindows) Warn(appId string, message string, icon string) {
	u.Notify(appId, "Warning", message, icon)
}

func (u *UiNotifyUtilsWindows) Error(appId string, message string, icon string) {
	u.Notify(appId, "Error", message, icon)
}
