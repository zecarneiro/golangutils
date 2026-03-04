package common

import (
	"crypto/rand"
	"fmt"
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

func GenerateTempName(prefixo string) string {
	bStr := ""
	// Create 4 random bytes (8 chars hexadecimais)
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err == nil {
		bStr = fmt.Sprintf("%x", b)
	}
	return fmt.Sprintf("%s-%s", prefixo, bStr)
}
