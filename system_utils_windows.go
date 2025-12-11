//go:build windows
// +build windows

package golangutils

import (
	"syscall"
)

func (s *SystemUtils) IsAdmin() bool {
	if s.IsWindows() {
		shell32 := syscall.NewLazyDLL("shell32.dll")
		isUserAnAdmin := shell32.NewProc("IsUserAnAdmin")
		r1, _, _ := isUserAnAdmin.Call()
		return r1 != 0
	}
	return false
}
