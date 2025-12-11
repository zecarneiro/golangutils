package enum

import (
	"strings"
)

type EPlatform string

const (
	UNKNOWN_PLATFORM = "UNKNOWN"
	UNIX             = "unix"
	DARWIN           = "darwin"
	LINUX            = "linux"
	WINDOWS          = "windows"
	FREEBSD          = "freebsd"
	OPENBSD          = "openbsd"
)

func EPlatformFromValue(value string) EPlatform {
	val := strings.ToLower(strings.TrimSpace(value))
	switch val {
	case "windows":
		return WINDOWS
	case "darwin":
		return DARWIN
	case "linux":
		return LINUX
	case "freebsd":
		return FREEBSD
	case "openbsd":
		return OPENBSD
	case "netbsd", "dragonfly", "solaris":
		return UNIX
	}
	return UNKNOWN_PLATFORM
}

func (s EPlatform) IsValid() bool {
	switch s {
	case UNKNOWN_PLATFORM, UNIX, DARWIN, LINUX, WINDOWS, FREEBSD, OPENBSD:
		return true
	default:
		return false
	}
}

func (s EPlatform) String() string {
	if s.IsValid() {
		return string(s)
	}
	return UNKNOWN_PLATFORM
}

func (s EPlatform) Equals(other EPlatform) bool {
	return s.String() == other.String()
}

func (s EPlatform) IsUnknown() bool {
	return s.String() == UNKNOWN_PLATFORM || !s.IsValid()
}
