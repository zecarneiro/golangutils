package env

import (
	"os"
	"slices"
	"strings"

	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/slice"
	"golangutils/pkg/str"
)

func GetSliceSeparator() string {
	// os.PathListSeparator é ';' on Windows and ':' on Linux/macOS
	return string(os.PathListSeparator)
}

func GetPathName() string {
	if platform.IsWindows() {
		return "Path"
	}
	return "PATH"
}

func Exists(name string) bool {
	_, exists := os.LookupEnv(name)
	return exists
}

func Get(name string) []string {
	valuesStr := os.Getenv(name)
	if str.IsEmpty(valuesStr) {
		return []string{}
	}
	return ConvValuesArr(valuesStr)
}

func Set(name string, values []string) {
	valuesStr := ConvValuesStr(slice.RemoveDuplicate(values))
	os.Setenv(name, valuesStr)
}

func SetBulk(envs []models.EnvData) {
	for _, data := range envs {
		Set(data.Key, data.Values)
	}
}

func Unset(name string) {
	os.Unsetenv(name)
}

func UnsetBulk(envs []models.EnvData) {
	for _, data := range envs {
		Unset(data.Key)
	}
}

func ConvValuesStr(values []string) string {
	valuesStr := slice.ArrayToStringBySep(values, GetSliceSeparator())
	return valuesStr
}

func ConvValuesArr(values string) []string {
	valuesArr := strings.Split(values, GetSliceSeparator())
	valuesArr = slice.FilterArray(valuesArr, func(val string) bool {
		return !str.IsEmpty(val)
	})
	return valuesArr
}

func HasValue(name, expectedValue string) bool {
	return slices.Contains(Get(name), strings.TrimSpace(expectedValue))
}

func ListFullInfo() map[string][]string {
	data := make(map[string][]string)
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		name := parts[0]
		data[name] = Get(name)
	}
	return data
}

func InsertOnPath(value string) {
	name := GetPathName()
	values := Get(name)
	values = append(values, value)
	Set(name, values)
}

func RemoveOnPath(value string) {
	name := GetPathName()
	values := Get(name)
	values = slice.FilterArray(values, func(val string) bool {
		return val != value
	})
	Set(name, values)
}
