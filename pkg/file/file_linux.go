//go:build linux

package file

import (
	"strings"

	"golangutils/pkg/platform"
)

func IsHidden(path string) (bool, error) {
	if platform.IsUnix() || platform.IsLinux() {
		basename := Basename(path)
		return strings.HasPrefix(basename, ".") && basename != "." && basename != "..", nil
	}
	return false, nil
}
