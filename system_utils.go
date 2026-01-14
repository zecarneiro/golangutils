package golangutils

import (
	"errors"
	"fmt"
	"golangutils/entities"
	"golangutils/enums"
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

func (s *SystemUtils) Platform() enums.EPlatform {
	platform := runtime.GOOS
	return enums.EPlatformFromValue(platform)
}

func (s *SystemUtils) IsWindows() bool {
	return s.Platform() == enums.WINDOWS
}

func (s *SystemUtils) IsLinux() bool {
	return s.Platform() == enums.LINUX
}

func (s *SystemUtils) IsDarwin() bool {
	return s.Platform() == enums.DARWIN
}

func (s *SystemUtils) IsUnix() bool {
	return s.Platform() == enums.UNIX
}

func (s *SystemUtils) IsPlatform(platforms []enums.EPlatform) bool {
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
	case enums.WINDOWS:
		cmd := exec.Command("powershell", "-Command", "(Get-CimInstance -ClassName Win32_OperatingSystem).Caption")
		output, err := cmd.Output()
		if err == nil && len(output) > 0 {
			osName = strings.TrimSpace(string(output))
		}
	case enums.LINUX:
		alreadySet := false
		ReadFileLineByLine("/etc/os-release", func(lineData string, err error) {
			if strings.HasPrefix(lineData, "PRETTY_NAME=") && !alreadySet {
				osName = strings.Trim(strings.TrimPrefix(lineData, "PRETTY_NAME="), "\"")
				alreadySet = true
			}
		})
	case enums.DARWIN:
		out, _ := exec.Command("sw_vers", "-productName").Output()
		ver, _ := exec.Command("sw_vers", "-productVersion").Output()
		osName = fmt.Sprintf("%s %s", strings.TrimSpace(string(out)), strings.TrimSpace(string(ver)))
	case enums.FREEBSD, enums.OPENBSD:
		out, _ := exec.Command("uname", "-sr").Output()
		osName = strings.TrimSpace(string(out))
	}
	return osName
}

func (s *SystemUtils) Reboot() error {
	var cmd entities.Command
	if s.IsWindows() {
		cmd = entities.Command{
			Cmd:  "shutdown",
			Args: []string{"/r", "/t", "0"},
		}
	} else if s.IsLinux() {
		cmd = entities.Command{
			Cmd:  "sudo",
			Args: []string{"shutdown", "-r", "now"},
		}
	} else if s.IsDarwin() {
		return errors.New(GetNotImplementedYetMsg())
	}
	console := NewConsoleUtils()
	if Confirm("Will be restart PC. Continue", true) {
		console.ExecRealTime(cmd)
		os.Exit(0)
	}
	return nil
}

func (s *SystemUtils) GetParentProcessInfo(ppid int) (*entities.ParentProcessInfo, error) {
	var parentInfo *entities.ParentProcessInfo
	switch s.Platform() {
	case enums.LINUX, enums.DARWIN, enums.UNIX:
		out, err := exec.Command("ps", "-p", strconv.Itoa(ppid), "-o", "ppid=,comm=").Output()
		if err != nil {
			return nil, err
		} else {
			fields := strings.Fields(string(out))
			if len(fields) >= 2 {
				parentPID, _ := strconv.Atoi(fields[0])
				name := fields[1]
				return &entities.ParentProcessInfo{
					Pid:  ppid,
					PPid: parentPID,
					Name: name,
				}, nil
			}
		}
	case enums.WINDOWS:
		out, err := exec.Command(
			"powershell",
			"-Command",
			fmt.Sprintf(
				"Get-CimInstance Win32_Process -Filter 'ProcessId = %d' | Select-Object Name, ParentProcessId | ForEach-Object { \"$($_.Name),$($_.ParentProcessId)\" }",
				ppid,
			),
		).Output()
		if err != nil {
			return parentInfo, err
		}
		parts := strings.Split(strings.TrimSpace(string(out)), ",")
		if len(parts) >= 2 {
			parentPID, _ := strconv.Atoi(parts[1])
			name := parts[0]
			return &entities.ParentProcessInfo{
				Pid:  ppid,
				PPid: parentPID,
				Name: name,
			}, nil
		}
	}
	return nil, errors.New(GetUnsupportedPlatformMsg())
}

func (s *SystemUtils) GetAncestralProcessInfo(currentPPid int) (*entities.ParentProcessInfo, error) {
	var err error
	var ancestralProcess *entities.ParentProcessInfo
	for {
		if currentPPid <= 4 {
			break
		}
		p, err_res := s.GetParentProcessInfo(currentPPid)
		if err_res != nil || p == nil || p.Pid == 0 {
			break
		}
		ancestralProcess = p
		if ancestralProcess.PPid == currentPPid {
			break
		}
		currentPPid = ancestralProcess.PPid
	}
	return ancestralProcess, err
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
