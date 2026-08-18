//go:build windows

package browser

import (
	"encoding/csv"
	"os/exec"
	"strconv"
	"strings"
)

func listBrowserProcesses() ([]Browser, error) {
	out, err := exec.Command("tasklist", "/fo", "csv", "/nh").Output()
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(strings.NewReader(string(out)))
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	var res []Browser
	for _, row := range rows {
		if len(row) < 2 {
			continue
		}
		img := strings.Trim(row[0], `"`)
		if !isBrowserBin(img) {
			continue
		}
		pid, err := strconv.Atoi(strings.Trim(row[1], `"`))
		if err != nil {
			continue
		}
		res = append(res, Browser{Name: img, PID: pid, ExePath: img})
	}
	return res, nil
}

// EnableDebugging is a no-op on Windows: remote debugging must be enabled by
// launching the browser with --remote-debugging-port=9222.
func EnableDebugging(b *Browser) error {
	return nil
}
