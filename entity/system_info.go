package entity

import (
	"os/user"
)

type CpuInfo struct {
	Cpu  int
	Arch string
}

type SystemInfo struct {
	UserInfo                   user.User
	HomeDir, TempDir, Hostname string
	Platform                   int
	PlatformName               string
	Eol                        string
	Uptime                     float64
	Cpu                        CpuInfo
}
