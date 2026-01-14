//go:build linux

package system

import (
	"os"

	"golangutils/pkg/platform"
)

func IsAdmin() bool {
	if platform.IsLinux() {
		return os.Geteuid() == 0
	}
	return false
}
