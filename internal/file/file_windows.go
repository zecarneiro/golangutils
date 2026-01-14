//go:build windows
// +build windows

package golangutils

import "syscall"

func IsHiddenFile(path string) (bool, error) {
	systemUtils := NewSystemUtils()
	if systemUtils.IsWindows() {
		ptr, err := syscall.UTF16PtrFromString(path)
		if err != nil {
			return false, err
		}
		attributes, err := syscall.GetFileAttributes(ptr)
		if err != nil {
			return false, err
		}
		return attributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0, nil
	}
	return false, nil
}
