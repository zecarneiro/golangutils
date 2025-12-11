package golangutils

import (
	"fmt"
	"runtime"
	"strings"
)

func GetInvalidPlatformMsg() string {
	return "Invalid Platform"
}

func GetUnsupportedPlatformMsg() string {
	return "Unsupported Platform"
}

func GetNotImplementedYetMsg() string {
	return "Not impplemented yet!"
}

func GetUnknowOSMsg() string {
	return "Unknown OS [" + strings.ToLower(runtime.GOOS) + "]"
}

func GetUnknownMsg(msg string) string {
	return fmt.Sprintf(msg, "Unknown")
}
