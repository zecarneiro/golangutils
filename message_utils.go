package golangutils

import (
	"runtime"
	"strings"
)

func GetInvalidPlatformMsg() string {
	return "Invalid Platform"
}

func GetNotImplementedYetMsg() string {
	return "Not impplemented yet!"
}

func GetUnknowOSMsg() string {
	return "Unknown OS [" + strings.ToLower(runtime.GOOS) + "]"
}
