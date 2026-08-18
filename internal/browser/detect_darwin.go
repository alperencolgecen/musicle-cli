//go:build darwin

package browser

import (
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func listBrowserProcesses() ([]Browser, error) {
	out, err := exec.Command("ps", "-eo", "pid=,command=").Output()
	if err != nil {
		return nil, err
	}
	var res []Browser
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		pid, err := strconv.Atoi(fields[0])
		if err != nil {
			continue
		}
		cmd := strings.Join(fields[1:], " ")
		exe := filepath.Base(fields[1])
		if !isBrowserBin(exe) {
			continue
		}
		b := Browser{Name: exe, PID: pid, ExePath: fields[1]}
		b.UserDataDir = parseUserDataSetDir(strings.Fields(cmd))
		res = append(res, b)
	}
	return res, nil
}

// EnableDebugging asks a running browser to open its remote debugging server
// by sending SIGUSR1 (a Chromium-supported signal on POSIX).
func EnableDebugging(b *Browser) error {
	return syscall.Kill(b.PID, syscall.SIGUSR1)
}
