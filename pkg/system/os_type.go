package system

import (
	"strings"
)

type OSType string

const (
	Windows   OSType = "windows"
	Ubuntu    OSType = "ubuntu"
	PopOs     OSType = "pop!_os"
	LinuxMint OSType = "linux mint"
)

func GetOSTypeFromValue(value string) OSType {
	val := strings.ToLower(strings.TrimSpace(value))
	if val == Windows.String() || strings.HasPrefix(strings.ToLower(value), Windows.String()) {
		return Windows
	} else if val == Ubuntu.String() || strings.HasPrefix(strings.ToLower(value), Ubuntu.String()) {
		return Ubuntu
	} else if val == PopOs.String() || strings.HasPrefix(strings.ToLower(value), PopOs.String()) {
		return PopOs
	} else if val == LinuxMint.String() || strings.HasPrefix(strings.ToLower(value), LinuxMint.String()) {
		return LinuxMint
	}
	return ""
}

func (p OSType) IsValid() bool {
	switch p {
	case Windows, Ubuntu, PopOs:
		return true
	default:
		return false
	}
}

func (p OSType) String() string {
	if p.IsValid() {
		return string(p)
	}
	return ""
}

func (p OSType) Equals(other OSType) bool {
	return p.String() == other.String()
}
