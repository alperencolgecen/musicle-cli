package engine

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

// TestSmokeExtractTools exercises the exact extraction path used at runtime:
// copyFS() copies the embedded bin tree into the cache directory. This
// validates the copy logic that makes the embedded yt-dlp/ffmpeg reachable at
// <Root>/engine_bin.
func TestSmokeExtractTools(t *testing.T) {
	fake := fstest.MapFS{
		"engine_bin/yt-dlp":      &fstest.MapFile{Data: []byte("#!/bin/sh"), Mode: 0o755},
		"engine_bin/ffmpeg":      &fstest.MapFile{Data: []byte("# ffmpeg"), Mode: 0o755},
		"engine_bin/ffprobe":     &fstest.MapFile{Data: []byte("# ffprobe"), Mode: 0o755},
		"engine_bin/.engine-stamp": &fstest.MapFile{Data: []byte("engine-v2\n")},
	}

	dst := t.TempDir()
	if err := copyFS(fake, "engine_bin", dst); err != nil {
		t.Fatalf("copyFS(engine_bin): %v", err)
	}

	want := []string{"yt-dlp", "ffmpeg", "ffprobe", ".engine-stamp"}
	for _, name := range want {
		if _, err := os.Stat(filepath.Join(dst, name)); err != nil {
			t.Errorf("tool %q not extracted to %s: %v", name, dst, err)
		}
	}
}
