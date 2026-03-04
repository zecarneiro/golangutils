package exe

import (
	"fmt"
	"os"
	"os/exec"

	"golangutils/pkg/common"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
)

func ExecRealTime(command models.Command) error {
	if command.Verbose {
		printCommand(command)
	}
	fillCommand(&command)
	if command.UseShell {
		command = detectShell(command)
	} else {
		cmd, err := buildNonShellCmd(command)
		if err != nil {
			return err
		}
		command = cmd
	}
	if command.FullVerbose {
		printCommand(command)
	}
	cmdResult := exec.Command(command.Cmd, command.Args...)
	cmdResult.Env = getEnv(command)
	cmdResult.Dir = command.Cwd
	cmdResult.Stdout = os.Stdout
	cmdResult.Stderr = os.Stderr
	cmdResult.Stdin = os.Stdin
	return cmdResult.Run()
}

func Exec(command models.Command) (string, error) {
	if command.Verbose {
		printCommand(command)
	}
	fillCommand(&command)
	if command.UseShell {
		command = detectShell(command)
	} else {
		cmd, err := buildNonShellCmd(command)
		if err != nil {
			return "", err
		}
		command = cmd
	}
	if command.FullVerbose {
		printCommand(command)
	}
	cmdResult := exec.Command(command.Cmd, command.Args...)
	cmdResult.Env = getEnv(command)
	cmdResult.Dir = command.Cwd
	output, err := cmdResult.CombinedOutput()
	if len(output) > 0 {
		return string(output), err
	}
	return "", err
}

func Chmod777(filepath string, verbose bool) error {
	filepath = file.ResolvePath(filepath)
	fileInfo := models.FileInfo{}
	if file.IsDir(filepath) {
		info, err := file.ReadDirRecursive(filepath)
		if err != nil {
			return err
		} else {
			fileInfo = info
		}
	} else if file.IsFile(filepath) {
		fileInfo = models.FileInfo{Files: []string{filepath}}
	} else {
		return fmt.Errorf("%s given file: %s", common.Unknown, filepath)
	}
	if verbose && (len(fileInfo.Files) > 0 || len(fileInfo.Directories) > 0) {
		logger.Info("Set full permission for '" + filepath + "'")
	}
	if platform.IsWindows() {
		for _, data := range fileInfo.Files {
			err := setFullAccessPowerShell(data)
			if err != nil {
				return err
			}
		}
		for _, data := range fileInfo.Directories {
			err := setFullAccessPowerShell(data)
			if err != nil {
				return err
			}
		}
	} else if platform.IsLinux() {
		for _, data := range fileInfo.Directories {
			err := os.Chmod(data, 0o777)
			if err != nil {
				return err
			}
		}
		for _, data := range fileInfo.Files {
			err := os.Chmod(data, 0o777)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return fmt.Errorf(common.NotImplementedYetMSG)
}

func GetExecutable() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return execPath, nil
}

func GetExecutableDir() (string, error) {
	execPath, err := GetExecutable()
	if err != nil {
		return "", err
	}
	return file.Dirname(execPath), nil
}
