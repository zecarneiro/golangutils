//go:build linux

package system

import (
	"golangutils/pkg/common/platform"
	"os"
)

func IsAdminLinux() bool {
	if platform.IsLinux() {
		return os.Geteuid() == 0
	}
	return false
}
