package enums

import (
	"strings"

	"golangutils/pkg/common"
)

type RegistryType string

const (
	CLASSES_ROOT    RegistryType = "CLASSES_ROOT"
	CURRENT_USER    RegistryType = "CURRENT_USER"
	LOCAL_MACHINE   RegistryType = "LOCAL_MACHINE"
	USERS           RegistryType = "USERS"
	UnknownRegistry RegistryType = common.Unknown
)

func GetRegistryTypeFromValue(value string) RegistryType {
	val := strings.ToUpper(strings.TrimSpace(value))
	switch val {
	case "CLASSES_ROOT":
		return CLASSES_ROOT
	case "CURRENT_USER":
		return CURRENT_USER
	case "LOCAL_MACHINE":
		return LOCAL_MACHINE
	case "USERS":
		return USERS
	default:
		return UnknownRegistry
	}
}

func (r RegistryType) IsValid() bool {
	switch r {
	case CLASSES_ROOT, CURRENT_USER, LOCAL_MACHINE, USERS:
		return true
	default:
		return false
	}
}

func (r RegistryType) String() string {
	if r.IsValid() {
		return string(r)
	}
	return string(UnknownRegistry)
}

func (r RegistryType) Equals(other RegistryType) bool {
	return r.String() == other.String()
}
