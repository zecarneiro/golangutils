package common

import (
	"fmt"
)

func GetUnknown(msg string) string {
	return fmt.Sprintf(msg, Unknown)
}
