//go:build linux

package system

import (
	"os"

	"golangutils/pkg/enums"
	"golangutils/pkg/env"
	"golangutils/pkg/platform"
	"golangutils/pkg/str"
)

func IsAdmin() bool {
	if platform.IsLinux() {
		return os.Geteuid() == 0
	}
	return false
}

func GetDesketopEnv() enums.DesktopEnvType {
	desktopEnv := enums.UnknownDE
	if !platform.IsLinux() {
		return desktopEnv
	}
	envDEArr := env.Get("XDG_CURRENT_DESKTOP")
	if len(envDEArr) <= 0 || len(envDEArr) > 2 {
		return desktopEnv
	}
	for _, envDE := range envDEArr {
		if str.Contains(envDE, enums.GnomeDE.String(), true) {
			desktopEnv = enums.GnomeDE
		} else if str.Contains(envDEArr[0], enums.KdeDE.String(), true) {
			desktopEnv = enums.KdeDE
		}
		if !enums.UnknownDE.Equals(desktopEnv) {
			break
		}
	}
	return desktopEnv
}
