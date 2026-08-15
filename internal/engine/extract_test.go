package engine

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestExtractStubBuild verifies that without the engine_assets tag the package
// returns the documented error instead of panicking.
func TestExtractStubBuild(t *testing.T) {
	if _, err := Extract(); err == nil {
		t.Skip("binary built with engine_assets; skipping stub probe")
	} else if err != ErrNoEmbeddedAssets {
		t.Fatalf("expected ErrNoEmbeddedAssets, got %v", err)
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
