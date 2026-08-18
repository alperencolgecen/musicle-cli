package browser

import "testing"

func TestCapTracks(t *testing.T) {
	if got := capTracks(nil); len(got) != 0 {
		t.Errorf("capTracks(nil) len = %d, want 0", len(got))
	}
	small := make([]Track, 10)
	if got := capTracks(small); len(got) != 10 {
		t.Errorf("capTracks(10) len = %d, want 10", len(got))
	}
	big := make([]Track, 150)
	if got := capTracks(big); len(got) != MaxTracks {
		t.Errorf("capTracks(150) len = %d, want %d", len(got), MaxTracks)
	}
}

func TestIsBrowserBin(t *testing.T) {
	yes := []string{"chrome", "google-chrome-stable", "msedge.exe", "Brave.exe", "chromium-browser"}
	for _, b := range yes {
		if !isBrowserBin(b) {
			t.Errorf("isBrowserBin(%q) = false, want true", b)
		}
	}
	no := []string{"firefox", "node", "explorer.exe", "spotify"}
	for _, b := range no {
		if isBrowserBin(b) {
			t.Errorf("isBrowserBin(%q) = true, want false", b)
		}
	}
}

func TestParseUserDataSetDir(t *testing.T) {
	fields := []string{"/usr/bin/chrome", "--user-data-dir", "/home/u/.config/chrome", "--remote-debugging-port=9222"}
	if got := parseUserDataSetDir(fields); got != "/home/u/.config/chrome" {
		t.Errorf("parseUserDataSetDir = %q, want %q", got, "/home/u/.config/chrome")
	}
	eq := []string{"/usr/bin/chrome", "--user-data-dir=/tmp/x"}
	if got := parseUserDataSetDir(eq); got != "/tmp/x" {
		t.Errorf("parseUserDataSetDir = %q, want %q", got, "/tmp/x")
	}
	if got := parseUserDataSetDir([]string{"nothing"}); got != "" {
		t.Errorf("parseUserDataSetDir = %q, want empty", got)
	}
}
