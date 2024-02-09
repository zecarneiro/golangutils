package jnoronha_golangutils

import (
	"errors"
	"jnoronha_golangutils/entities"
	"os"
	"os/user"
	"runtime"
	"strings"
)

func SystemInfo() entities.SystemInfo {
	info := entities.SystemInfo{
		TempDir:  os.TempDir(),
		Eol:      "\n",
		Platform: platform(),
		Cpu:      entities.CpuInfo{Cpu: runtime.NumCPU(), Arch: runtime.GOARCH},
	}
	currentUser, err := user.Current()
	if err == nil {
		info.UserInfo = *currentUser
		info.HomeDir = currentUser.HomeDir
	}
	hostname, err := os.Hostname()
	if err == nil {
		info.Hostname = hostname
	}
	if IsWindows() {
		info.Eol = "\r\n"
	}

	return info
}

func platform() int {
	if runtime.GOOS == "windows" {
		return entities.WINDOWS
	} else if runtime.GOOS == "darwin" {
		return entities.DARWIN
	} else if runtime.GOOS == "linux" {
		return entities.LINUX
	} else if runtime.GOOS == "unix" {
		return entities.UNIX
	} else {
		return entities.NONE
	}
}

func IsWindows() bool {
	return platform() == entities.WINDOWS
}

func IsDarwin() bool {
	return platform() == entities.DARWIN
}

func IsLinux() bool {
	return platform() == entities.LINUX
}

func IsUnix() bool {
	return platform() == entities.UNIX
}

func ValidateSystem() {
	if platform() == entities.NONE {
		ProcessError(errors.New("Unknown OS [" + strings.ToLower(runtime.GOOS) + "]"))
	}
}
