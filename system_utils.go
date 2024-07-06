package golangutils

import (
	"errors"
	"os"
	"os/user"
	"runtime"
)

/* -------------------------------------------------------------------------- */
/*                                 MODEL AREA                                 */
/* -------------------------------------------------------------------------- */
type CpuInfo struct {
	Cpu  int
	Arch string
}

type SystemInfo struct {
	TempDir, HomeDir, Hostname, Eol string
	Platform                        int
	Uptime                          float64
	UserInfo                        user.User
	Cpu                             CpuInfo
}
/* ----------------------------- END MODEL AREA ----------------------------- */

const (
	NONE    = 0
	UNIX    = 1
	DARWIN  = 2
	LINUX   = 3
	WINDOWS = 4
)

func SysInfo() SystemInfo {
	info := SystemInfo{
		TempDir:  os.TempDir(),
		Eol:      "\n",
		Platform: platform(),
		Cpu:      CpuInfo{Cpu: runtime.NumCPU(), Arch: runtime.GOARCH},
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
		return WINDOWS
	} else if runtime.GOOS == "darwin" {
		return DARWIN
	} else if runtime.GOOS == "linux" {
		return LINUX
	} else {
		return NONE
	}
}

func IsWindows() bool {
	return platform() == WINDOWS
}

func IsDarwin() bool {
	return platform() == DARWIN
}

func IsLinux() bool {
	return platform() == LINUX
}

func ValidateSystem() {
	if platform() == NONE {
		ErrorLog(UNKNOW_OS_MSG, false)
		os.Exit(1)
	}
}

func Reboot() error {
	var cmdInfo CommandInfo
	if IsWindows() {
		cmdInfo = CommandInfo{
			Cmd:     "shutdown",
			Args:    []string{"/r", "/t", "0"},
			EnvVars: os.Environ(),
		}
	} else if IsLinux() {
		cmdInfo = CommandInfo{
			Cmd:     "sudo",
			Args:    []string{"shutdown", "-r", "now"},
			EnvVars: os.Environ(),
		}
	} else if IsDarwin() {
		return errors.New(NOT_IMPLEMENTED_YET_MSG)
	}
	if Confirm("Will be restart PC. Continue", true) {
		ExecRealTime(cmdInfo)
		os.Exit(0)
	}
	return nil
}
