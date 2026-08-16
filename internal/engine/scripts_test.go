package engine

import (
	"io/fs"
	"testing"
)

// TestBinEmbedded verifies the bundled tools are present in the embedded
// filesystem when the engine_assets tag is active. On a stub build this is a
// no-op (the FS is empty and probeEmbedded reports the missing stamp).
func TestBinEmbedded(t *testing.T) {
	if _, err := fs.ReadFile(binFS, "engine_bin/.engine-stamp"); err != nil {
		t.Skipf("tools not embedded in this build: %v", err)
	}
	for _, name := range []string{"yt-dlp", "ffmpeg", "ffprobe"} {
		if _, err := fs.ReadFile(binFS, "engine_bin/"+name); err != nil {
			t.Errorf("tool %q missing from embed: %v", name, err)
		}
	}
}
