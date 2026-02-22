package enums

import (
	"golangutils/pkg/common"
	"strings"
)

type OSType string

const (
	WindowsSO   OSType = "windows"
	UbuntuSO    OSType = "ubuntu"
	PopOsSO     OSType = "pop!_os"
	LinuxMintSO OSType = "linux mint"
	UnknownSO   OSType = common.Unknown
)

func GetOSTypeFromValue(value string) OSType {
	val := strings.ToLower(strings.TrimSpace(value))
	if val == WindowsSO.String() || strings.HasPrefix(strings.ToLower(value), WindowsSO.String()) {
		return WindowsSO
	} else if val == UbuntuSO.String() || strings.HasPrefix(strings.ToLower(value), UbuntuSO.String()) {
		return UbuntuSO
	} else if val == PopOsSO.String() || strings.HasPrefix(strings.ToLower(value), PopOsSO.String()) {
		return PopOsSO
	} else if val == LinuxMintSO.String() || strings.HasPrefix(strings.ToLower(value), LinuxMintSO.String()) {
		return LinuxMintSO
	}
	return UbuntuSO
}

func (p OSType) IsValid() bool {
	switch p {
	case WindowsSO, UbuntuSO, PopOsSO, LinuxMintSO:
		return true
	default:
		return false
	}
}

func (p OSType) String() string {
	if p.IsValid() {
		return string(p)
	}
	return UbuntuSO.String()
}

func (p OSType) Equals(other OSType) bool {
	return p.String() == other.String()
}
