//go:build windows

package system

import (
	"syscall"

	"golangutils/pkg/enums"
	"golangutils/pkg/platform"

	"golang.org/x/sys/windows/registry"
)

func IsAdmin() bool {
	if platform.IsWindows() {
		shell32 := syscall.NewLazyDLL("shell32.dll")
		isUserAnAdmin := shell32.NewProc("IsUserAnAdmin")
		r1, _, _ := isUserAnAdmin.Call()
		return r1 != 0
	}
	return false
}

func GetDesketopEnv() enums.DesktopEnvType {
	if !platform.IsWindows() {
		return enums.UnknownDE
	}
	return enums.WindowsDE
}

func GetRegeditValueStr(regPath string, typeRegedit enums.RegistryType, keyName string) string {
	var typeRegeditKey registry.Key
	switch typeRegedit {
	case enums.CLASSES_ROOT:
		typeRegeditKey = registry.CLASSES_ROOT
	case enums.CURRENT_USER:
		typeRegeditKey = registry.CURRENT_USER
	case enums.LOCAL_MACHINE:
		typeRegeditKey = registry.LOCAL_MACHINE
	case enums.USERS:
		typeRegeditKey = registry.USERS
	}
	k, err := registry.OpenKey(typeRegeditKey, regPath, registry.QUERY_VALUE)
	if err != nil {
		return ""
	}
	defer k.Close()
	value, _, _ := k.GetStringValue(keyName)
	return value
}
