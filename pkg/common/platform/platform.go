package platform

import (
	"fmt"
	"golangutils/pkg/common"
	"runtime"
	"strings"
)

func GetPlatform() PlatformType {
	if platform == nil {
		val := GetPlatformTypeFromValue(runtime.GOOS)
		platform = &val
	}
	return *platform
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
	return common.InArray(platforms, GetPlatform())
}

func GetUnknowOS() string {
	return fmt.Sprintf("Unknown OS [%s]", strings.ToLower(runtime.GOOS))
}
