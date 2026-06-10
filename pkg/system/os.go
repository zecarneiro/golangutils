package system

import (
	"fmt"
	"os/exec"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"golangutils/pkg/common"
	"golangutils/pkg/enums"
	"golangutils/pkg/file"
	"golangutils/pkg/platform"
	"golangutils/pkg/str"
)

func OSName() string {
	osName := common.GetUnknown("%s OS NAME")
	switch platform.GetPlatform() {
	case enums.Windows:
		cmd := exec.Command(getPwshCmd(), "-Command", "(Get-CimInstance -ClassName Win32_OperatingSystem).Caption")
		output, err := cmd.Output()
		if err == nil && len(output) > 0 {
			osName = strings.TrimSpace(string(output))
		}
	case enums.Linux:
		alreadySet := false
		file.ReadFileLineByLine("/etc/os-release", func(lineData string) {
			if strings.HasPrefix(lineData, "NAME=") && !alreadySet {
				osName = strings.Trim(strings.TrimPrefix(lineData, "NAME="), "\"")
				alreadySet = true
			}
		})
	case enums.Darwin:
		out, _ := exec.Command("sw_vers", "-productName").Output()
		ver, _ := exec.Command("sw_vers", "-productVersion").Output()
		osName = fmt.Sprintf("%s %s", strings.TrimSpace(string(out)), strings.TrimSpace(string(ver)))
	case enums.FreeBSD, enums.OpenBSD:
		out, _ := exec.Command("uname", "-sr").Output()
		osName = strings.TrimSpace(string(out))
	}
	return osName
}

func OSVersion() string {
	osVersion := common.GetUnknown("%s OS VERSION")
	switch platform.GetPlatform() {
	case enums.Windows:
		registryPath := `SOFTWARE\Microsoft\Windows NT\CurrentVersion`
		displayVersion := GetRegeditValueStr(registryPath, enums.LOCAL_MACHINE, "DisplayVersion")
		buildNumberStr := GetRegeditValueStr(registryPath, enums.LOCAL_MACHINE, "CurrentBuild")
		if !str.IsEmpty(displayVersion) && !str.IsEmpty(buildNumberStr) {
			winMajor := "10"
			buildNumber, _ := strconv.Atoi(buildNumberStr)
			if buildNumber >= 22000 {
				winMajor = "11"
			}
			osVersion = fmt.Sprintf(`%s %s`, winMajor, displayVersion)
		}
	case enums.Linux:
		alreadySet := false
		file.ReadFileLineByLine("/etc/os-release", func(lineData string) {
			if strings.HasPrefix(lineData, "VERSION_ID=") && !alreadySet {
				osVersion = strings.Trim(strings.TrimPrefix(lineData, "VERSION_ID="), "\"")
				alreadySet = true
			}
		})
	case enums.Darwin:
		ver, _ := exec.Command("sw_vers", "-productVersion").Output()
		osVersion = strings.TrimSpace(string(ver))
	case enums.FreeBSD, enums.OpenBSD:
		out, _ := exec.Command("uname", "-sr").Output()
		osVersion = strings.TrimSpace(string(out))
	}
	return osVersion
}

func GetOsType() enums.OSType {
	if osType == nil {
		val := enums.GetOSTypeFromValue(OSName())
		osType = &val
	}
	return *osType
}

func IsOsType(osTypeList []enums.OSType) bool {
	return slices.Contains(osTypeList, GetOsType())
}

func GetUnknowOS() string {
	return fmt.Sprintf("%s OS [%s]", common.Unknown, strings.ToLower(runtime.GOOS))
}

func OSFullName() string {
	osName := common.GetUnknown("%s OS NAME")
	switch platform.GetPlatform() {
	case enums.Windows:
		registryPath := `SOFTWARE\Microsoft\Windows NT\CurrentVersion`
		name := OSName()
		displayVersion := GetRegeditValueStr(registryPath, enums.LOCAL_MACHINE, "DisplayVersion")
		buildNumber := GetRegeditValueStr(registryPath, enums.LOCAL_MACHINE, "CurrentBuild")
		if osName != name && !str.IsEmpty(displayVersion) && !str.IsEmpty(buildNumber) {
			osName = fmt.Sprintf(`%s %s (Build %s)`, name, displayVersion, buildNumber)
		}
	case enums.Linux:
		alreadySet := false
		file.ReadFileLineByLine("/etc/os-release", func(lineData string) {
			if strings.HasPrefix(lineData, "PRETTY_NAME=") && !alreadySet {
				osName = strings.Trim(strings.TrimPrefix(lineData, "PRETTY_NAME="), "\"")
				alreadySet = true
			}
		})
	case enums.Darwin:
		osName = OSName()
	case enums.FreeBSD, enums.OpenBSD:
		osName = OSName()
	}
	return osName
}
