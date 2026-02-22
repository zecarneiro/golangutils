package env

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"golangutils/pkg/models"
	"golangutils/pkg/slice"
	"golangutils/pkg/str"
)

var pathName = "PATH"

func GetSliceSeparator() string {
	// os.PathListSeparator é ';' no Windows e ':' no Linux/macOS
	return string(os.PathListSeparator)
}

func Get(name string) string {
	return os.Getenv(name)
}

func Set(name string, value string) {
	os.Setenv(name, value)
}

func Unset(name string) {
	os.Unsetenv(name)
}

func SetBulk(envs []models.EnvData) {
	for _, data := range envs {
		Set(data.Key, data.Value)
	}
}

func UnsetBulk(envs []models.EnvData) {
	for _, data := range envs {
		Unset(data.Key)
	}
}

func VarExists(name string) bool {
	_, exists := os.LookupEnv(name)
	return exists
}

func VarValuesAsList(name string) []string {
	value := Get(name)
	if str.IsEmpty(value) {
		return []string{}
	}
	parts := strings.Split(value, GetSliceSeparator())
	result := []string{}
	for _, part := range parts {
		result = append(result, part)
	}
	return result
}

func VarHasValue(name, expectedValue string) bool {
	return slices.Contains(VarValuesAsList(name), strings.TrimSpace(expectedValue))
}

func VarList() map[string][]string {
	data := make(map[string][]string)
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		name := parts[0]
		data[name] = VarValuesAsList(name)
	}
	return data
}

func GetPath() []string {
	return VarValuesAsList(pathName)
}

func GetPathStr() string {
	return Get(pathName)
}

func InsertOnPath(value string) {
	if !VarHasValue(pathName, value) {
		currentPath := GetPathStr()
		newPath := fmt.Sprintf("%s%s%s", value, GetSliceSeparator(), currentPath)
		Set(pathName, newPath)
	}
}

func RemoveOnPath(value string) {
	if VarHasValue(pathName, value) {
		listValues := GetPath()
		newListValues := slice.FilterArray(listValues, func(val string) bool {
			return val != value
		})
		newPath := slice.ArrayToStringBySep(newListValues, GetSliceSeparator())
		Set(pathName, newPath)
	}
}
