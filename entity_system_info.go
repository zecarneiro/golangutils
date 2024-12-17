package golangutils

import (
	"os"
	"os/user"
	"runtime"
	"strings"
)

type SystemInfo struct {
	userInfo                   user.User
	homeDir, tempDir, hostname string
	platform                   int
	platformName               string
	isSetPlatform              bool
	eol                        string
	uptime                     float64
	cpu                        CpuInfo
}

func NewSystemInfo() SystemInfo {
	return SystemInfo{}
}

func (s *SystemInfo) isValidUser() bool {
	s.GetUserInfo()
	return s.userInfo != (user.User{})
}

func (s *SystemInfo) setDistroLinuxName() {
	cmd := NewConsole(Logger{})
	res := cmd.Exec(Command{Cmd: "lsb_release -a | grep 'Distributor ID' | awk '{print $3}'"})
	if !res.HasError() {
		s.platformName = strings.TrimSpace(res.Data)
		s.platformName = strings.ToLower(s.platformName)
	}
}

func (s *SystemInfo) GetUserInfo() user.User {
	if s.userInfo == (user.User{}) {
		currentUser, err := user.Current()
		if err == nil {
			s.userInfo = *currentUser
		}
	}
	return s.userInfo
}

func (s *SystemInfo) GetHomeDir() string {
	if len(s.homeDir) == 0 && s.isValidUser() {
		s.homeDir = s.GetUserInfo().HomeDir
	}
	return s.homeDir
}

func (s *SystemInfo) GetTempDir() string {
	if len(s.tempDir) == 0 {
		s.tempDir = os.TempDir()
	}
	return s.tempDir
}

func (s *SystemInfo) GetHostname() string {
	if len(s.hostname) == 0 {
		hostname, err := os.Hostname()
		if err == nil {
			s.hostname = hostname
		}
	}
	return s.hostname
}

func (s *SystemInfo) GetPlatform() int {
	if !s.isSetPlatform {
		s.platformName = runtime.GOOS
		if s.platformName == "windows" {
			s.platform = WINDOWS
		} else if s.platformName == "darwin" {
			s.platform = DARWIN
		} else if s.platformName == "linux" {
			s.platform = LINUX
			s.setDistroLinuxName()
		} else {
			s.platform = NONE
		}
		s.isSetPlatform = true
	}
	return s.platform
}

func (s *SystemInfo) IsWindows() bool {
	return s.GetPlatform() == WINDOWS
}

func (s *SystemInfo) IsDarwin() bool {
	return s.GetPlatform() == DARWIN
}

func (s *SystemInfo) IsLinux() bool {
	return s.GetPlatform() == LINUX
}

func (s *SystemInfo) GetEol() string {
	if len(s.eol) == 0 {
		s.eol = "\n"
		if s.IsWindows() {
			s.eol = "\r\n"
		}
	}
	return s.eol
}

func (s *SystemInfo) GetCpu() CpuInfo {
	if s.cpu == (CpuInfo{}) {
		s.cpu = CpuInfo{Cpu: runtime.NumCPU(), Arch: runtime.GOARCH}
	}
	return s.cpu
}

func (s *SystemInfo) ValidatePlatform(canExit bool, logger Logger) bool {
	var isValid = s.GetPlatform() != NONE
	if !isValid {
		logger.Error(GetUnknowOSMsg())
		if canExit {
			os.Exit(1)
		}
	}
	return isValid
}

func (s *SystemInfo) GetPlatformName() string {
	s.GetPlatform()
	return s.platformName
}
