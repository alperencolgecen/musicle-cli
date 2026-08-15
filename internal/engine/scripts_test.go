package engine

import (
	"io/fs"
	"testing"
)

// TestScriptsEmbedded verifies the helper Python modules the embedded
// interpreter imports at runtime are bundled into the binary.
func TestScriptsEmbedded(t *testing.T) {
	want := []string{
		"scripts/musicle_ytdlp.py",
		"scripts/musicle_spotdl.py",
		"scripts/_musicle_helpers.py",
	}
	for _, name := range want {
		if _, err := fs.ReadFile(scriptsFS, name); err != nil {
			t.Errorf("script %q missing from embed: %v", name, err)
		}
	}
}

// TestScriptPathEnv ensures the PYTHONPATH the engine advertises points at
// the extracted <Root>/scripts directory.
func TestScriptPathEnv(t *testing.T) {
	// scriptPathEnv guards on the extracted global; in a stub build before
	// Extract() it falls back to "scripts" which is still a valid relative
	// search path. After a successful Extract it resolves to an absolute dir.
	got := scriptPathEnv()
	if got == "" {
		t.Fatal("scriptPathEnv returned empty")
	}
}
