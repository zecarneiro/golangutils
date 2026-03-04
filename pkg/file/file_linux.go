//go:build linux

package file

import (
	"fmt"
	"strings"
	"syscall"

	"golangutils/pkg/platform"
)

func IsHidden(path string) (bool, error) {
	if platform.IsUnix() || platform.IsLinux() {
		basename := Basename(path)
		return strings.HasPrefix(basename, ".") && basename != "." && basename != "..", nil
	}
	return false, nil
}

func GetDevice(path string) (string, error) {
	var stat syscall.Stat_t
	if err := syscall.Stat(path, &stat); err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", stat.Dev), nil
}
