package engine

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"testing/fstest"
)

// TestSmokeExtractScripts exercises the exact extraction path used at runtime:
// doExtract() copies the embedded scriptsFS rooted at "scripts" into the cache
// directory. This validates the copy logic that makes the embedded Python
// helpers reachable at <Root>/scripts (the PYTHONPATH the engine advertises).
func TestSmokeExtractScripts(t *testing.T) {
	fake := fstest.MapFS{
		"scripts/musicle_ytdlp.py":    &fstest.MapFile{Data: []byte("# yt-dlp wrapper")},
		"scripts/musicle_spotdl.py":   &fstest.MapFile{Data: []byte("# spotdl wrapper")},
		"scripts/_musicle_helpers.py": &fstest.MapFile{Data: []byte("# helpers")},
	}

	dst := t.TempDir()
	if err := copyFS(fake, "scripts", dst); err != nil {
		t.Fatalf("copyFS(scripts): %v", err)
	}

	want := []string{
		"musicle_ytdlp.py",
		"musicle_spotdl.py",
		"_musicle_helpers.py",
	}
	for _, name := range want {
		if _, err := os.Stat(filepath.Join(dst, name)); err != nil {
			t.Errorf("script %q not extracted to %s: %v", name, dst, err)
		}
	}
}

// TestSmokeNoEngineFallback confirms that, when the engine is not embedded (the
// default stub build), Extract() reports ErrNoEmbeddedAssets so the bridge's
// runEngine() can fall back to the legacy pipeline instead of panicking.
func TestSmokeNoEngineFallback(t *testing.T) {
	t.Setenv("MUSICLE_CACHE_DIR", t.TempDir())
	// Force a fresh extract attempt regardless of prior package state.
	extracted = nil
	extractOnce = sync.Once{}

	_, err := Extract()
	if err == nil {
		t.Skip("engine embedded in this build; skipping stub fallback check")
	}
	if err != ErrNoEmbeddedAssets {
		t.Fatalf("expected ErrNoEmbeddedAssets, got %v", err)
	}
}
