//go:build linux
// +build linux

package golangutils

import (
	"strings"
)

func IsHiddenFile(path string) (bool, error) {
	systemUtils := NewSystemUtils()
	if systemUtils.IsUnix() || systemUtils.IsLinux() {
		basename := Basename(path)
		return strings.HasPrefix(basename, ".") && basename != "." && basename != "..", nil
	}
	return false, nil
}
