package engine

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Progress is the callback signature used by wrappers to stream download
// progress back to the Go side. pct is 0..100, msg is human-readable.
type Progress func(pct int, msg string)

// ProgressEvent is a structured progress message emitted by the embedded
// Python helpers via stdout. The wrappers parse these and forward them to
// the Progress callback.
type ProgressEvent struct {
	Event   string  `json:"event"`             // "start" | "progress" | "track" | "done" | "error"
	Percent float64 `json:"percent,omitempty"` // 0..100
	Message string  `json:"message,omitempty"` // human readable
	Track   string  `json:"track,omitempty"`   // artist - title when relevant
	Index   int     `json:"index,omitempty"`
	Total   int     `json:"total,omitempty"`
}

// RunOptions tunes how the Python process is invoked.
type RunOptions struct {
	WorkDir    string            // cwd for the subprocess; default = output dir's parent
	Env        []string          // extra env vars appended after the default set
	Args       []string          // arguments passed after the helper script
	Progress   Progress          // optional progress callback
	JSONStdout bool              // when true, stdout is parsed as JSON ProgressEvent stream
	StderrCb   func(line string) // optional stderr line sink (for debug logs)
	Timeout    time.Duration     // 0 = no timeout
}

// Run executes the embedded python with the given RunOptions. It is the
// shared primitive used by both the yt-dlp and spotdl wrappers.
func Run(ext *Extracted, opt RunOptions) error {
	if ext == nil {
		return ErrNoEmbeddedAssets
	}
	if opt.Progress == nil {
		opt.Progress = func(int, string) {}
	}

	ctx := context.Background()
	if opt.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opt.Timeout)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, ext.PythonBin, opt.Args...)
	cmd.Env = buildEnv(ext, opt.Env)
	if opt.WorkDir != "" {
		cmd.Dir = opt.WorkDir
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("engine: stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("engine: stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("engine: başlatılamadı: %w", err)
	}

	var wg sync.WaitGroup

	if opt.JSONStdout {
		wg.Add(1)
		go func() {
			defer wg.Done()
			scanJSON(stdout, opt.Progress)
		}()
	} else {
		// Drain stdout so the process never blocks on a full pipe.
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = io.Copy(io.Discard, stdout)
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		scanStderr(stderr, opt.StderrCb)
	}()

	wg.Wait()
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("engine: çalıştırma hatası: %w", err)
	}
	return nil
}

// buildEnv returns a minimal environment pointing at the embedded python and
// ffmpeg binaries. PYTHONUNBUFFERED is set so progress events flush promptly.
func buildEnv(ext *Extracted, extra []string) []string {
	env := []string{
		"PYTHONUNBUFFERED=1",
		"PYTHONDONTWRITEBYTECODE=1",
		"PYTHONIOENCODING=utf-8",
		"MUSICLE_ENGINE=1",
	}
	// Prepend the venv bin to PATH so subprocesses (ffmpeg) inherit a clean lookup.
	if path := os.Getenv("PATH"); path != "" {
		env = append(env, "PATH="+filepath.Dir(ext.PythonBin)+string(os.PathListSeparator)+path)
	} else {
		env = append(env, "PATH="+filepath.Dir(ext.PythonBin))
	}
	env = append(env, extra...)
	return env
}

func scanJSON(r io.Reader, cb Progress) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || !strings.HasPrefix(line, "{") {
			continue
		}
		var ev ProgressEvent
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			continue
		}
		switch ev.Event {
		case "done":
			cb(100, ev.Message)
		case "error":
			cb(int(ev.Percent), "Hata: "+ev.Message)
		default:
			pct := int(ev.Percent)
			if pct < 0 {
				pct = 0
			}
			if pct > 100 {
				pct = 100
			}
			if ev.Message == "" && ev.Track != "" {
				ev.Message = ev.Track
			}
			cb(pct, ev.Message)
		}
	}
}

func scanStderr(r io.Reader, cb func(string)) {
	if cb == nil {
		_, _ = io.Copy(io.Discard, r)
		return
	}
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			cb(line)
		}
	}
}
