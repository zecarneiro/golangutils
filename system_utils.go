package golangutils

import (
	"errors"
	"fmt"
	"golangutils/entity"
	"golangutils/enum"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type SystemUtils struct {
	loggerUtils *LoggerUtils
}

func NewSystemUtilsDefault() *SystemUtils {
	return &SystemUtils{loggerUtils: NewLoggerUtils()}
}

func NewSystemUtils(loggerUtils *LoggerUtils) *SystemUtils {
	return &SystemUtils{loggerUtils: loggerUtils}
}

func (s *SystemUtils) Eol() string {
	if s.IsWindows() {
		return "\r\n"
	}
	return "\n"
}

func (s *SystemUtils) Platform() enum.EPlatform {
	platform := runtime.GOOS
	return enum.EPlatformFromValue(platform)
}

func (s *SystemUtils) IsWindows() bool {
	return s.Platform() == enum.WINDOWS
}

func (s *SystemUtils) IsLinux() bool {
	return s.Platform() == enum.LINUX
}

func (s *SystemUtils) IsDarwin() bool {
	return s.Platform() == enum.DARWIN
}

func (s *SystemUtils) IsUnix() bool {
	return s.Platform() == enum.UNIX
}

func (s *SystemUtils) IsPlatform(platforms []enum.EPlatform) bool {
	return InArray(platforms, s.Platform())
}

func (s *SystemUtils) HomeDir() string {
	home, _ := os.UserHomeDir()
	return home
}

func (s *SystemUtils) TempDir() string {
	return os.TempDir()
}

func (s *SystemUtils) OSName() string {
	osName := GetUnknownMsg("%s OS NAME")
	switch s.Platform() {
	case enum.WINDOWS:
		cmd := exec.Command("powershell", "-Command", "(Get-CimInstance -ClassName Win32_OperatingSystem).Caption")
		output, err := cmd.Output()
		if err == nil && len(output) > 0 {
			osName = strings.TrimSpace(string(output))
		}
	case enum.LINUX:
		alreadySet := false
		ReadFileLineByLine("/etc/os-release", func(lineData string, err error) {
			if strings.HasPrefix(lineData, "PRETTY_NAME=") && !alreadySet {
				osName = strings.Trim(strings.TrimPrefix(lineData, "PRETTY_NAME="), "\"")
				alreadySet = true
			}
		})
	case enum.DARWIN:
		out, _ := exec.Command("sw_vers", "-productName").Output()
		ver, _ := exec.Command("sw_vers", "-productVersion").Output()
		osName = fmt.Sprintf("%s %s", strings.TrimSpace(string(out)), strings.TrimSpace(string(ver)))
	case enum.FREEBSD, enum.OPENBSD:
		out, _ := exec.Command("uname", "-sr").Output()
		osName = strings.TrimSpace(string(out))
	}
	return osName
}

func (s *SystemUtils) Reboot() error {
	var cmd entity.Command
	if s.IsWindows() {
		cmd = entity.Command{
			Cmd:  "shutdown",
			Args: []string{"/r", "/t", "0"},
		}
	} else if s.IsLinux() {
		cmd = entity.Command{
			Cmd:  "sudo",
			Args: []string{"shutdown", "-r", "now"},
		}
	} else if s.IsDarwin() {
		return errors.New(GetNotImplementedYetMsg())
	}
	console := NewConsoleUtilsDefault()
	if Confirm("Will be restart PC. Continue", true) {
		console.ExecRealTime(cmd)
		os.Exit(0)
	}
	return nil
}

func GetParentProcess(ppid int) (string, error) {
	systemUtils := NewSystemUtilsDefault()
	switch systemUtils.Platform() {
	case enum.LINUX, enum.DARWIN, enum.UNIX:
		out, err := exec.Command("ps", "-p", strconv.Itoa(ppid), "-o", "comm=").Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(out)), nil
	case enum.WINDOWS:
		out, err := exec.Command(
			"powershell",
			"-Command",
			fmt.Sprintf("(Get-Process -Id %d).Name", ppid),
		).Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(out)), nil
	}
	return "", errors.New(GetUnsupportedPlatformMsg())
}

func EnvVarExists(name string) bool {
	_, exists := os.LookupEnv(name)
	return exists
}

func EnvVarValuesAsList(name string) []string {
	value := os.Getenv(name)
	if value == "" {
		return []string{}
	}
	// os.PathListSeparator é ';' no Windows e ':' no Linux/macOS
	parts := strings.Split(value, string(os.PathListSeparator))
	result := []string{}
	for _, part := range parts {
		result = append(result, part)
	}
	return result
}

func EnvVarHasValue(name, expectedValue string) bool {
	return InArray(EnvVarValuesAsList(name), strings.TrimSpace(expectedValue))
}

func EnvVarList() map[string][]string {
	data := make(map[string][]string)
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		name := parts[0]
		data[name] = EnvVarValuesAsList(name)
	}
	return data
}
