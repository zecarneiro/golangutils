package platform

import (
	"golangutils/pkg/enums"
	"runtime"
	"slices"
)

func GetPlatform() enums.PlatformType {
	if platformType == nil {
		val := enums.GetPlatformTypeFromValue(runtime.GOOS)
		platformType = &val
	}
	return *platformType
}

func IsWindows() bool {
	return GetPlatform() == enums.Windows
}

func IsLinux() bool {
	return GetPlatform() == enums.Linux
}

func IsDarwin() bool {
	return GetPlatform() == enums.Darwin
}

func IsUnix() bool {
	return GetPlatform() == enums.Unix
}

func IsPlatform(platforms []enums.PlatformType) bool {
	return slices.Contains(platforms, GetPlatform())
}
