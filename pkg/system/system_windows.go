//go:build windows

package system

import (
	"golangutils/pkg/common/platform"
	"syscall"
)

func IsAdminWindows() bool {
	if platform.IsWindows() {
		shell32 := syscall.NewLazyDLL("shell32.dll")
		isUserAnAdmin := shell32.NewProc("IsUserAnAdmin")
		r1, _, _ := isUserAnAdmin.Call()
		return r1 != 0
	}
	return false
}
