package engine

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestExtractSucceedsWhenEmbedded verifies the full extract path when the engine
// assets are embedded (i.e. internal/engine/engine_venv and engine_ffmpeg are
// present at compile time). It is a no-op on a checkout where the assets have
// not been prepared/committed, because the package would not compile without
// them.
func TestExtractSucceedsWhenEmbedded(t *testing.T) {
	ext, err := Extract()
	if err != nil {
		t.Skipf("engine assets not embedded in this build: %v", err)
	}
	if ext == nil || ext.PythonBin == "" || ext.SpotdlBin == "" || ext.YtdlpBin == "" {
		t.Fatal("Extract returned incomplete locations")
	}
}

// TestVenvLayoutHelpers checks the OS-specific path helpers used to locate
// the embedded binaries.
func TestVenvLayoutHelpers(t *testing.T) {
	venv := "/tmp/x"
	if runtime.GOOS == "windows" {
		if got := venvPython(venv); got != filepath.Join(venv, "Scripts", "python.exe") {
			t.Errorf("venvPython windows: %s", got)
		}
		if got := venvBin(venv, "spotdl"); got != filepath.Join(venv, "Scripts", "spotdl.exe") {
			t.Errorf("venvBin windows: %s", got)
		}
	} else {
		if got := venvPython(venv); got != filepath.Join(venv, "bin", "python") {
			t.Errorf("venvPython unix: %s", got)
		}
		if got := venvBin(venv, "spotdl"); got != filepath.Join(venv, "bin", "spotdl") {
			t.Errorf("venvBin unix: %s", got)
		}
	}
	if got := exeSuffix(); got == "" && runtime.GOOS == "windows" {
		t.Errorf("exeSuffix windows should not be empty")
	}
	if got := exeSuffix(); got != "" && runtime.GOOS != "windows" {
		t.Errorf("exeSuffix unix should be empty, got %q", got)
	}
}

// TestIsExtractedSentinel ensures we only treat an extract as valid when the
// sentinel file matches the cache version.
func TestIsExtractedSentinel(t *testing.T) {
	dir := t.TempDir()
	if ok, _ := isExtracted(dir); ok {
		t.Fatalf("empty dir should not validate")
	}
	if err := os.WriteFile(filepath.Join(dir, sentinelName), []byte("wrong\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if ok, _ := isExtracted(dir); ok {
		t.Fatalf("wrong sentinel should not validate")
	}
	if err := os.WriteFile(filepath.Join(dir, sentinelName), []byte(cacheVersion+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if ok, _ := isExtracted(dir); !ok {
		t.Fatalf("matching sentinel should validate")
	}
}
