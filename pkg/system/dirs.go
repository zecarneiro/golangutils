package system

import "os"

func TempDir() string {
	return os.TempDir()
}

func HomeDir() string {
	home, _ := os.UserHomeDir()
	return home
}
