//go:build linux

package browser

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func listBrowserProcesses() ([]Browser, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}
	var out []Browser
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(e.Name())
		if err != nil {
			continue
		}
		data, err := os.ReadFile(filepath.Join("/proc", e.Name(), "cmdline"))
		if err != nil {
			continue
		}
		fields := strings.Fields(strings.ReplaceAll(string(data), "\x00", " "))
		if len(fields) == 0 {
			continue
		}
		if !isBrowserBin(filepath.Base(fields[0])) {
			continue
		}
		b := Browser{Name: filepath.Base(fields[0]), PID: pid, ExePath: fields[0]}
		b.UserDataDir = parseUserDataSetDir(fields)
		out = append(out, b)
	}
	return out, nil
}

// EnableDebugging asks a running browser to open its remote debugging server
// by sending SIGUSR1 (a Chromium-supported signal on POSIX).
func EnableDebugging(b *Browser) error {
	return syscall.Kill(b.PID, syscall.SIGUSR1)
}
