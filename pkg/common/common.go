package common

import (
	"runtime"
)

func Eol() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

func IsNil(arg any) bool {
	return arg == nil
}
