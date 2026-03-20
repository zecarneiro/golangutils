package system

import (
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/platform"
	"golangutils/pkg/str"
	"os"
)

func TempDir() string {
	return os.TempDir()
}

func HomeDir() string {
	home, _ := os.UserHomeDir()
	return home
}

func HomeUserConfigDir() string {
	config_dir := file.ResolvePath(HomeDir() + "/.config")
	logger.Error(file.CreateDirectory(config_dir, true))
	return config_dir
}

func HomeUserLocalDir() string {
	config_dir := file.JoinPath(HomeDir(), ".local")
	logger.Error(file.CreateDirectory(config_dir, true))
	return config_dir
}

func HomeUserOptDir() string {
	opt_dir := file.JoinPath(HomeUserLocalDir(), "opt")
	logger.Error(file.CreateDirectory(opt_dir, true))
	return opt_dir
}

func HomeUserBinDir() string {
	bin_dir := file.JoinPath(HomeUserLocalDir(), "bin")
	logger.Error(file.CreateDirectory(bin_dir, true))
	return bin_dir
}

func HomeUserStartupDir() string {
	startup_dir := ""
	if platform.IsWindows() {
		startup_dir = file.JoinPath(HomeDir(), "Start Menu", "Programs", "Startup")
	} else if platform.IsLinux() {
		startup_dir = file.JoinPath(HomeUserConfigDir(), "autostart")
	}
	if !str.IsEmpty(startup_dir) {
		logger.Error(file.CreateDirectory(startup_dir, true))
	}
	return startup_dir
}

func HomeUserTempDir() string {
	temp_dir := ""
	if platform.IsWindows() {
		temp_dir = file.ResolvePath(TempDir())
	} else if platform.IsLinux() {
		temp_dir = file.ResolvePath(HomeUserLocalDir() + "/tmp")
	}
	if !str.IsEmpty(temp_dir) {
		logger.Error(file.CreateDirectory(temp_dir, true))
	}
	return temp_dir
}
