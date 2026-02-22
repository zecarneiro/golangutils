package system

import (
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/enums"
	"runtime"
	"strings"
)

func GetOsType() enums.OSType {
	if osType == nil {
		val := enums.GetOSTypeFromValue(OSName())
		osType = &val
	}
	return *osType
}

func GetUnknowOS() string {
	return fmt.Sprintf("%s OS [%s]", common.Unknown, strings.ToLower(runtime.GOOS))
}
