package system

import (
	"errors"
	"fmt"
	"os/exec"
	"os/user"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"golangutils/pkg/common"
	"golangutils/pkg/console"
	"golangutils/pkg/enums"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
	"golangutils/pkg/str"
	"golangutils/pkg/timec"
)

func Reboot() error {
	if console.Confirm("Will be restart the PC. Continue?", true) {
		shutdownCmdStr := "shutdown"
		shutdownCmd, _ := console.Which(shutdownCmdStr)
		if !str.IsEmpty(shutdownCmd) {
			logger.Info(fmt.Sprintf(`Founded %s: %s`, shutdownCmdStr, shutdownCmd))
			timec.Sleep(1)
			if platform.IsWindows() {
				return common.PrepareExecCommand(exec.Command(shutdownCmd, "/r", "/t", "0", "/f"), true).Run()
			} else if platform.IsLinux() || platform.IsDarwin() {
				err := common.PrepareExecCommand(exec.Command("sudo", shutdownCmd, "-r", "now"), true).Run()
				if err != nil {
					if platform.IsLinux() {
						systemctlCmdStr := "systemctl"
						systemctlCmd, _ := console.Which(systemctlCmdStr)
						logger.Error(fmt.Errorf(`Failed to run: %s`, shutdownCmd))
						logger.Info(fmt.Sprintf(`Trying %s`, systemctlCmdStr))
						timec.Sleep(1)
						if !str.IsEmpty(shutdownCmd) {
							logger.Info(fmt.Sprintf(`Founded %s: %s`, shutdownCmdStr, shutdownCmd))
							err = common.PrepareExecCommand(exec.Command("sudo", systemctlCmd, "reboot", "-i"), true).Run()
						} else {
							return fmt.Errorf("Not found %s", systemctlCmdStr)
						}
					}
					return err
				}
			} else {
				return errors.New(common.NotImplementedYetMSG)
			}
		} else {
			return fmt.Errorf("Not found %s", shutdownCmdStr)
		}
	}
	return nil
}

func Shutdown() error {
	var cmd *exec.Cmd
	if console.Confirm("Will be shutdown the PC. Continue?", true) {
		shutdownCmd, err := console.Which("shutdown")
		if err != nil {
			return err
		}
		if platform.IsWindows() {
			cmd = common.PrepareExecCommand(exec.Command(shutdownCmd, "/s", "/t", "0"), true)
		} else if platform.IsLinux() {
			cmd = common.PrepareExecCommand(exec.Command("sudo", shutdownCmd, "-h", "now"), true)
		} else if platform.IsDarwin() {
			return errors.New(common.NotImplementedYetMSG)
		}
		return cmd.Run()
	}
	return nil
}

func GetParentProcessInfo(ppid int) (*models.ParentProcessInfo, error) {
	var parentInfo *models.ParentProcessInfo
	switch platform.GetPlatform() {
	case enums.Linux, enums.Darwin, enums.Unix:
		out, err := common.PrepareExecCommand(exec.Command("ps", "-p", strconv.Itoa(ppid), "-o", "ppid=,comm="), true).Output()
		if err != nil {
			return nil, err
		} else {
			fields := strings.Fields(string(out))
			if len(fields) >= 2 {
				parentPID, _ := strconv.Atoi(fields[0])
				name := fields[1]
				return &models.ParentProcessInfo{
					Pid:  ppid,
					PPid: parentPID,
					Name: name,
				}, nil
			}
		}
	case enums.Windows:
		args := []string{
			"-Command",
			fmt.Sprintf(
				"Get-CimInstance Win32_Process -Filter 'ProcessId = %d' | Select-Object Name, ParentProcessId | ForEach-Object { \"$($_.Name),$($_.ParentProcessId)\" }",
				ppid,
			),
		}
		out, err := common.PrepareExecCommand(exec.Command(getPwshCmd(), args...), true).Output()
		if err != nil {
			return parentInfo, err
		}
		parts := strings.Split(strings.TrimSpace(string(out)), ",")
		if len(parts) >= 2 {
			parentPID, _ := strconv.Atoi(parts[1])
			name := parts[0]
			return &models.ParentProcessInfo{
				Pid:  ppid,
				PPid: parentPID,
				Name: name,
			}, nil
		}
	}
	return nil, errors.New(platform.UnsupportedMSG)
}

func GetAncestralProcessInfo(currentPPid int) (*models.ParentProcessInfo, error) {
	var err error
	var ancestralProcess *models.ParentProcessInfo
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

func IsDesktopEnv(desktopEnv enums.DesktopEnvType) bool {
	currentDesktopEnv := GetDesketopEnv()
	return currentDesktopEnv.Equals(desktopEnv)
}

func IsDesktopEnvs(desktopEnvs []enums.DesktopEnvType) bool {
	currentDesktopEnv := GetDesketopEnv()
	return slices.Contains(desktopEnvs, currentDesktopEnv)
}

func IsValidUserHomeDir(verbose bool) bool {
	var status bool
	homeDir := HomeDir()
	homeDirBasename := file.Basename(homeDir)
	if strings.Contains(homeDirBasename, " ") {
		status = false
	} else {
		invalidRegex := `[!@#$%^&*(),?"":{}|<>=´]|[à-ü]|[À-Ü]`
		re := regexp.MustCompile(invalidRegex)
		status = logic.Ternary(re.MatchString(homeDirBasename), false, true)
	}

	if verbose && !status {
		logger.Info(fmt.Sprintf("Your full home dir is: %s", homeDir))
		logger.Info(fmt.Sprintf("Your home dir is: %s", homeDirBasename))
		logger.Title("Usernames must")
		logger.Prompt("Start with an alphabetic character")
		logger.Prompt("Not contain empty spaces or @")
		logger.Prompt("Contain only valid Unix Characters - letters, numbers, '-', '.', and '_'")
		if platform.IsWindows() {
			logger.Prompt("Length: 20 characters or fewer for Windows")
			logger.Prompt("Be different from the device host name on Windows")
			logger.Prompt("When setting for a Windows device, usernames can't end with a period (.) or else they will not appear on the device login screen")
		}
	}
	return status
}

func UserInfo() user.User {
	currentUser, err := user.Current()
	logic.ProcessError(err)
	return *currentUser
}
