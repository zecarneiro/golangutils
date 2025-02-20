package golangutils

import (
	"errors"
	"golangutils/entity"
	"golangutils/enum"
	"os"
	"os/user"
	"regexp"
	"runtime"
	"strings"
)

type SystemUtils struct {
	info        entity.SystemInfo
	loggerUtils *LoggerUtils
}

func NewSystemUtils(loggerUtils *LoggerUtils) SystemUtils {
	return SystemUtils{loggerUtils: loggerUtils}
}

func (s *SystemUtils) setDistroLinuxName() {
	osReleaseFile := "/etc/os-release"
	callback := func(line string, err error) {
		match, err := regexp.MatchString("^NAME=", line)
		if err == nil && match {
			lineArr := strings.Split(line, "=")
			s.info.PlatformName = lineArr[1]
			s.info.PlatformName = strings.Trim(s.info.PlatformName, "\"")
			s.info.PlatformName = strings.TrimSpace(s.info.PlatformName)
			s.info.PlatformName = strings.ToLower(s.info.PlatformName)
		}
	}
	ReadFileLineByLine(osReleaseFile, callback)
}

func (s *SystemUtils) setPlatform() {
	if s.info.PlatformName == "windows" {
		s.info.Platform = enum.WINDOWS
	} else if s.info.PlatformName == "darwin" {
		s.info.Platform = enum.DARWIN
	} else if s.info.PlatformName == "linux" {
		s.info.Platform = enum.LINUX
		s.setDistroLinuxName()
	} else {
		s.info.Platform = enum.NONE
	}
}

func (s *SystemUtils) getEol() string {
	if s.IsWindows() {
		return "\r\n"
	}
	return "\n"
}

func (s *SystemUtils) Info() entity.SystemInfo {
	if s.info == (entity.SystemInfo{}) {
		s.info = entity.SystemInfo{
			TempDir:      os.TempDir(),
			PlatformName: runtime.GOOS,
			Cpu:          entity.CpuInfo{Cpu: runtime.NumCPU(), Arch: runtime.GOARCH},
		}

		// Platform
		s.setPlatform()

		// EOL
		s.info.Eol = s.getEol()

		// User Info and Home dir
		currentUser, err := user.Current()
		if err == nil {
			s.info.UserInfo = *currentUser
			s.info.HomeDir = s.info.UserInfo.HomeDir
		}

		// Hostname
		hostname, err := os.Hostname()
		if err == nil {
			s.info.Hostname = hostname
		}
	}
	return s.info
}

func (s *SystemUtils) IsWindows() bool {
	return s.Info().Platform == enum.WINDOWS
}

func (s *SystemUtils) IsDarwin() bool {
	return s.Info().Platform == enum.DARWIN
}

func (s *SystemUtils) IsLinux() bool {
	return s.Info().Platform == enum.LINUX
}

func (s *SystemUtils) ValidatePlatform(canExit bool, logger LoggerUtils) bool {
	var isValid = s.Info().Platform != enum.NONE
	if !isValid {
		logger.Error(GetUnknowOSMsg())
		if canExit {
			os.Exit(1)
		}
	}
	return isValid
}

func (s *SystemUtils) Reboot(console ConsoleUtils) error {
	var cmd entity.Command
	if s.IsWindows() {
		cmd = entity.Command{
			Cmd:     "shutdown",
			Args:    []string{"/r", "/t", "0"},
			EnvVars: os.Environ(),
		}
	} else if s.IsLinux() {
		cmd = entity.Command{
			Cmd:     "sudo",
			Args:    []string{"shutdown", "-r", "now"},
			EnvVars: os.Environ(),
		}
	} else if s.IsDarwin() {
		return errors.New(GetNotImplementedYetMsg())
	}
	if console.Confirm("Will be restart PC. Continue", true) {
		console.ExecRealTime(cmd)
		os.Exit(0)
	}
	return nil
}

func HasOsArgs() bool {
	argsWithoutProg := os.Args[1:]
	return len(argsWithoutProg) > 0
}

func GetOsArgs() []string {
	if HasOsArgs() {
		return os.Args[1:]
	}
	return []string{}
}
