package conv

import (
	"fmt"
	"strconv"
	"strings"

	"golangutils/pkg/common"
)

func StringToInt(data string) (int, error) {
	return strconv.Atoi(data)
}

func IntToString(data int) string {
	return strconv.Itoa(data)
}

func BytesToMB(b int64) float64 {
	return float64(b) / 1024 / 1024
}

func BoolToString(data bool) string {
	return strconv.FormatBool(data)
}

func StringToBool(data string) (bool, error) {
	data = strings.Trim(data, "\r\n")
	data = strings.Trim(data, "\n")
	return strconv.ParseBool(data)
}

func ToString(data any) string {
	if common.IsNil(&data) {
		return ""
	}
	switch data.(type) {
	case string:
		return fmt.Sprintf("%s", data)
	default:
		return fmt.Sprintf("%v", data)
	}
}
