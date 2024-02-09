package entities

import "os/user"

type CpuInfo struct {
	Cpu  int
	Arch string
}

type SystemInfo struct {
	TempDir, HomeDir, Hostname, Eol string
	Platform                        int
	Uptime                          float64
	UserInfo                        user.User
	Cpu                             CpuInfo
}

/*



   platform: EPlatformType,
   userInfo: os.UserInfo<string>,
   cpus: os.CpuInfo[],
*/
