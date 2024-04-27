package enums

import (
	"runtime"
	"strings"
)

const (
	INVALID_PLATFORM_MSG    = "Invalid Platform"
	NOT_IMPLEMENTED_YET_MSG = "Not impplemented yet!"
)

var (
	UNKNOW_OS_MSG = "Unknown OS [" + strings.ToLower(runtime.GOOS) + "]"
)
