//go:build linux
// +build linux

package golangutils

import "os"

func (s *SystemUtils) IsAdmin() bool {
	if s.IsLinux() {
		return os.Geteuid() == 0
	}
	return false
}
