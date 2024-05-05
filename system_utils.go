package jnoronhautils

import (
	"errors"
	"jnoronhautils/entities"
	"jnoronhautils/enums"
	"os"
	"os/user"
	"runtime"
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

func ValidateSystem() {
	if platform() == entities.NONE {
		ErrorLog(enums.UNKNOW_OS_MSG, false)
		os.Exit(1)
	}
}

func Reboot() error {
	var cmdInfo entities.CommandInfo
	if IsWindows() {
		cmdInfo = entities.CommandInfo{
			Cmd:     "shutdown",
			Args:    []string{"/r", "/t", "0"},
			EnvVars: os.Environ(),
		}
	} else if IsLinux() {
		cmdInfo = entities.CommandInfo{
			Cmd:     "sudo",
			Args:    []string{"shutdown", "-r", "now"},
			EnvVars: os.Environ(),
		}
	} else if IsDarwin() {
		return errors.New(enums.NOT_IMPLEMENTED_YET_MSG)
	}
	if Confirm("Will be restart PC. Continue", true) {
		ExecRealTime(cmdInfo)
		os.Exit(0)
	}
	return nil
}
