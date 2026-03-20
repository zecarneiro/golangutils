//go:build windows

package console

import (
	"golangutils/pkg/platform"
	"os"

	"golang.org/x/sys/windows"
)

func EnableFeatures() {
	if platform.IsWindows() {
		handle := windows.Handle(os.Stdout.Fd())
		var mode uint32
		if err := windows.GetConsoleMode(handle, &mode); err == nil {
			// Adiciona a flag de processamento de sequências virtuais (ANSI)
			mode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
			windows.SetConsoleMode(handle, mode)
		}
	}
}
