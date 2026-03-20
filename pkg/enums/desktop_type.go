package enums

import (
	"golangutils/pkg/common"
	"strings"
)

type DesktopEnvType string

const (
	GnomeDE   DesktopEnvType = "GNOME"
	KdeDE     DesktopEnvType = "KDE"
	WindowsDE DesktopEnvType = "WINDOWS"
	UnknownDE DesktopEnvType = common.Unknown
)

func GetDesktopEnvTypeFromValue(value string) DesktopEnvType {
	val := strings.ToLower(strings.TrimSpace(value))
	switch val {
	case "gnome":
		return GnomeDE
	case "kde":
		return KdeDE
	case "windows":
		return WindowsDE
	default:
		return UnknownDE
	}
}

func (s DesktopEnvType) IsValid() bool {
	switch s {
	case GnomeDE, KdeDE, WindowsDE:
		return true
	default:
		return false
	}
}

func (s DesktopEnvType) String() string {
	if s.IsValid() {
		return string(s)
	}
	return string(UnknownDE)
}

func (s DesktopEnvType) Equals(other DesktopEnvType) bool {
	return s.String() == other.String()
}
