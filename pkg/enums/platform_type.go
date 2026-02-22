package enums

import (
	"golangutils/pkg/common"
	"strings"
)

type PlatformType string

const (
	Unix            PlatformType = "unix"
	Darwin          PlatformType = "darwin"
	Linux           PlatformType = "linux"
	Windows         PlatformType = "windows"
	FreeBSD         PlatformType = "freebsd"
	OpenBSD         PlatformType = "openbsd"
	UnknownPlatform PlatformType = common.Unknown
)

func GetPlatformTypeFromValue(value string) PlatformType {
	val := strings.ToLower(strings.TrimSpace(value))
	switch val {
	case "windows":
		return Windows
	case "darwin":
		return Darwin
	case "linux":
		return Linux
	case "freebsd":
		return FreeBSD
	case "openbsd":
		return OpenBSD
	case "netbsd", "dragonfly", "solaris":
		return Unix
	}
	return UnknownPlatform
}

func (p PlatformType) IsValid() bool {
	switch p {
	case Unix, Darwin, Linux, Windows, FreeBSD, OpenBSD:
		return true
	default:
		return false
	}
}

func (p PlatformType) String() string {
	if p.IsValid() {
		return string(p)
	}
	return UnknownPlatform.String()
}

func (p PlatformType) Equals(other PlatformType) bool {
	return p.String() == other.String()
}
