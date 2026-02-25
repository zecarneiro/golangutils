package system

import (
	"fmt"
	"runtime"
	"slices"
	"strings"

	"golangutils/pkg/common"
	"golangutils/pkg/enums"
)

func GetOsType() enums.OSType {
	if osType == nil {
		val := enums.GetOSTypeFromValue(OSName())
		osType = &val
	}
	return *osType
}

func IsOsType(osTypeList []enums.OSType) bool {
	return slices.Contains(osTypeList, GetOsType())
}

func GetUnknowOS() string {
	return fmt.Sprintf("%s OS [%s]", common.Unknown, strings.ToLower(runtime.GOOS))
}
