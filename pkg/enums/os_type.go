package enums

import (
	"fmt"
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
	if val == WindowsSO.String() || strings.HasPrefix(val, WindowsSO.String()) || strings.HasPrefix(val, fmt.Sprintf("microsoft %s", WindowsSO.String())) {
		return WindowsSO
	} else if val == UbuntuSO.String() || strings.HasPrefix(val, UbuntuSO.String()) {
		return UbuntuSO
	} else if val == PopOsSO.String() || strings.HasPrefix(val, PopOsSO.String()) {
		return PopOsSO
	} else if val == LinuxMintSO.String() || strings.HasPrefix(val, LinuxMintSO.String()) {
		return LinuxMintSO
	}
	return UnknownSO
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
	return string(UbuntuSO)
}

func (p OSType) Equals(other OSType) bool {
	return p.String() == other.String()
}
