package platform

import (
	"fmt"
	"runtime"
	"slices"
	"strings"
)

func GetPlatform() PlatformType {
	if platformType == nil {
		val := GetPlatformTypeFromValue(runtime.GOOS)
		platformType = &val
	}
	return *platformType
}

func IsWindows() bool {
	return GetPlatform() == Windows
}

func IsLinux() bool {
	return GetPlatform() == Linux
}

func IsDarwin() bool {
	return GetPlatform() == Darwin
}

func IsUnix() bool {
	return GetPlatform() == Unix
}

func IsPlatform(platforms []PlatformType) bool {
	return slices.Contains(platforms, GetPlatform())
}

func GetUnknowOS() string {
	return fmt.Sprintf("Unknown OS [%s]", strings.ToLower(runtime.GOOS))
}
