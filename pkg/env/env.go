package env

import (
	"os"
	"slices"
	"strings"

	"golangutils/pkg/models"
)

func SetEnv(key string, value string) {
	os.Setenv(key, value)
}

func UnsetEnv(key string) {
	os.Unsetenv(key)
}

func SetEnvBulk(envs []models.EnvData) {
	for _, data := range envs {
		SetEnv(data.Key, data.Value)
	}
}

func UnsetEnvBulk(envs []models.EnvData) {
	for _, data := range envs {
		UnsetEnv(data.Key)
	}
}

func EnvVarExists(name string) bool {
	_, exists := os.LookupEnv(name)
	return exists
}

func EnvVarValuesAsList(name string) []string {
	value := os.Getenv(name)
	if value == "" {
		return []string{}
	}
	// os.PathListSeparator é ';' no Windows e ':' no Linux/macOS
	parts := strings.Split(value, string(os.PathListSeparator))
	result := []string{}
	for _, part := range parts {
		result = append(result, part)
	}
	return result
}

func EnvVarHasValue(name, expectedValue string) bool {
	return slices.Contains(EnvVarValuesAsList(name), strings.TrimSpace(expectedValue))
}

func EnvVarList() map[string][]string {
	data := make(map[string][]string)
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		name := parts[0]
		data[name] = EnvVarValuesAsList(name)
	}
	return data
}
