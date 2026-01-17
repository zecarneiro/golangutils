package system

import (
	"errors"
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/common/platform"
	"golangutils/pkg/entity"
	"golangutils/pkg/file"
	"golangutils/pkg/generic"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Eol() string {
	if platform.IsWindows() {
		return "\r\n"
	}
	return "\n"
}

func OSName() string {
	osName := common.GetUnknown("%s OS NAME")
	switch platform.GetPlatform() {
	case platform.Windows:
		cmd := exec.Command("powershell", "-Command", "(Get-CimInstance -ClassName Win32_OperatingSystem).Caption")
		output, err := cmd.Output()
		if err == nil && len(output) > 0 {
			osName = strings.TrimSpace(string(output))
		}
	case platform.Linux:
		alreadySet := false
		file.ReadFileLineByLine("/etc/os-release", func(lineData string, err error) {
			if strings.HasPrefix(lineData, "PRETTY_NAME=") && !alreadySet {
				osName = strings.Trim(strings.TrimPrefix(lineData, "PRETTY_NAME="), "\"")
				alreadySet = true
			}
		})
	case platform.Darwin:
		out, _ := exec.Command("sw_vers", "-productName").Output()
		ver, _ := exec.Command("sw_vers", "-productVersion").Output()
		osName = fmt.Sprintf("%s %s", strings.TrimSpace(string(out)), strings.TrimSpace(string(ver)))
	case platform.FreeBSD, platform.OpenBSD:
		out, _ := exec.Command("uname", "-sr").Output()
		osName = strings.TrimSpace(string(out))
	}
	return osName
}

func GetOsType() OSType {
	return GetOSTypeFromValue(OSName())
}

func Reboot() error {
	var cmd *exec.Cmd
	if generic.Confirm("Will be restart PC. Continue", true) {
		if platform.IsWindows() {
			cmd = exec.Command("shutdown", "/r", "/t", "0", "/f")
		} else if platform.IsLinux() {
			cmd = exec.Command("sudo", "shutdown", "-r", "now")
		} else if platform.IsDarwin() {
			return errors.New(common.NotImplementedYetMSG)
		}
		return cmd.Run()
	}
	return nil
}

func GetParentProcessInfo(ppid int) (*entity.ParentProcessInfo, error) {
	var parentInfo *entity.ParentProcessInfo
	switch platform.GetPlatform() {
	case platform.Linux, platform.Darwin, platform.Unix:
		out, err := exec.Command("ps", "-p", strconv.Itoa(ppid), "-o", "ppid=,comm=").Output()
		if err != nil {
			return nil, err
		} else {
			fields := strings.Fields(string(out))
			if len(fields) >= 2 {
				parentPID, _ := strconv.Atoi(fields[0])
				name := fields[1]
				return &entity.ParentProcessInfo{
					Pid:  ppid,
					PPid: parentPID,
					Name: name,
				}, nil
			}
		}
	case platform.Windows:
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
			return &entity.ParentProcessInfo{
				Pid:  ppid,
				PPid: parentPID,
				Name: name,
			}, nil
		}
	}
	return nil, errors.New(platform.UnsupportedMSG)
}

func GetAncestralProcessInfo(currentPPid int) (*entity.ParentProcessInfo, error) {
	var err error
	var ancestralProcess *entity.ParentProcessInfo
	for {
		if currentPPid <= 4 {
			break
		}
		p, err_res := GetParentProcessInfo(currentPPid)
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
	return common.InArray(EnvVarValuesAsList(name), strings.TrimSpace(expectedValue))
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
