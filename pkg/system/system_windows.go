//go:build windows

package system

import (
	"syscall"

	"golangutils/pkg/platform"
)

func IsAdmin() bool {
	if platform.IsWindows() {
		shell32 := syscall.NewLazyDLL("shell32.dll")
		isUserAnAdmin := shell32.NewProc("IsUserAnAdmin")
		r1, _, _ := isUserAnAdmin.Call()
		return r1 != 0
	}
	return false
}
