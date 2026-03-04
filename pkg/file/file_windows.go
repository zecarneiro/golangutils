//go:build windows

package file

import (
	"path/filepath"
	"strings"
	"syscall"

	"golangutils/pkg/platform"
)

func IsHidden(path string) (bool, error) {
	if platform.IsWindows() {
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

func GetDevice(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	vol := filepath.VolumeName(absPath) // ex: "C:"
	return strings.ToUpper(vol), nil
}
