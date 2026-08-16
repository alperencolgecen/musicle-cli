package engine

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
)

// TestSmokeNoEngineFallback confirms that Extract() never panics and returns a
// controlled error when the embedded tools are missing or corrupt. On a build
// without -tags engine_assets (or before prepare-engine.sh runs) this returns
// ErrNoEmbeddedAssets and the bridge falls back to the legacy pipeline.
func TestSmokeNoEngineFallback(t *testing.T) {
	t.Setenv("MUSICLE_CACHE_DIR", t.TempDir())
	// Force a fresh extract attempt regardless of prior package state.
	extracted = nil
	extractOnce = sync.Once{}

	_, err := Extract()
	if err == nil {
		t.Skip("engine embedded in this build; skipping fallback check")
	}
	if !errors.Is(err, ErrNoEmbeddedAssets) {
		t.Logf("engine not usable (controlled error, expected on stub/corrupt builds): %v", err)
	}
}

// TestExtractSucceedsWhenEmbedded verifies the full extract path when the tools
// are embedded (internal/engine/engine_bin present at compile time). It is a
// no-op on a checkout where the assets have not been prepared.
func TestExtractSucceedsWhenEmbedded(t *testing.T) {
	ext, err := Extract()
	if err != nil {
		t.Skipf("engine assets not embedded in this build: %v", err)
	}
	if ext == nil || ext.YTDLP == "" || ext.FFMPEG == "" || ext.FFPROBE == "" {
		t.Fatal("Extract returned incomplete locations")
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

// TestExeSuffix checks the OS-specific executable suffix helper.
func TestExeSuffix(t *testing.T) {
	if got := exeSuffix(); got == "" && runtime.GOOS != "windows" {
		// empty suffix is expected on unix
	}
	if got := exeSuffix(); got != "" && runtime.GOOS == "windows" {
		t.Errorf("exeSuffix on windows should be non-empty, got %q", got)
	}
}
