package engine

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestIntegrationDownloadReal exercises the full embedded pipeline against a
// real YouTube video ("Me at the zoo", 19s): Extract() -> yt-dlp metadata +
// bestaudio + thumbnail -> ffmpeg mux to 320k mp3 with an APIC cover. It
// requires the engine to be embedded (built with -tags engine_assets) and
// network access, so it is skipped unless MUSICLE_INTEGRATION=1.
func TestIntegrationDownloadReal(t *testing.T) {
	if os.Getenv("MUSICLE_INTEGRATION") == "" {
		t.Skip("MUSICLE_INTEGRATION=1 set to run network download test")
	}
	ext, err := Extract()
	if err != nil {
		t.Skipf("engine assets not embedded in this build: %v", err)
	}
	if ext == nil || ext.YTDLP == "" || ext.FFMPEG == "" {
		t.Fatal("Extract returned incomplete locations")
	}

	out := t.TempDir()
	err = DownloadYouTube(YouTubeOptions{
		OutputDir: out,
		URLs:      []string{"https://www.youtube.com/watch?v=jNQXAC9IVRw"},
		Progress:  func(pct int, msg string) { t.Logf("%3d%% %s", pct, msg) },
	})
	if err != nil {
		t.Fatalf("DownloadYouTube: %v", err)
	}

	mp3s, _ := filepath.Glob(filepath.Join(out, "*.mp3"))
	if len(mp3s) == 0 {
		t.Fatal("engine üretilen mp3 bulunamadı")
	}
	f := mp3s[0]
	data, err := os.ReadFile(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) < 200_000 {
		t.Fatalf("mp3 şüpheli şekilde küçük (%d bayt)", len(data))
	}
	if !bytes.HasPrefix(data, []byte("ID3")) {
		t.Errorf("mp3 ID3 başlığı taşımıyor (kapak gömülü değil)")
	}
	if !bytes.Contains(data, []byte("APIC")) {
		t.Errorf("mp3 APIC kapak çerçevesi içermiyor")
	}
	t.Logf("üretilen dosya: %s (%d bayt, ID3+APIC)", filepath.Base(f), len(data))
}

// TestSpotifyResolveTrack verifies the API-free Spotify -> YouTube query
// resolution against the live embed page (network-gated).
func TestSpotifyResolveTrack(t *testing.T) {
	if os.Getenv("MUSICLE_INTEGRATION") == "" {
		t.Skip("MUSICLE_INTEGRATION=1 set to run network test")
	}
	q, err := spotifyQueryFor("https://open.spotify.com/track/4uLU6hMCjMI75M1A2tKUQC")
	if err != nil {
		t.Fatalf("spotify track resolve: %v", err)
	}
	if !strings.Contains(q, "Never Gonna Give You Up") || !strings.Contains(q, "Rick Astley") {
		t.Fatalf("beklenmeyen sorgu: %q", q)
	}
	t.Logf("çözülen sorgu: %q", q)

	// Collections must be routed to the bridge fallback, not resolved here.
	if _, err := spotifyQueryFor("https://open.spotify.com/album/1ATL5GLyefJaxhQzSPVrLX"); err == nil {
		t.Fatalf("albüm sorgusu hata vermemeli (bridge fallback'e düşmeli)")
	}
}

// TestIntegrationSpotifyTrackDownload runs a full API-free Spotify track
// download through the embedded engine (embed-page scrape -> ytsearch1 ->
// mp3 320k with cover). Network-gated.
func TestIntegrationSpotifyTrackDownload(t *testing.T) {
	if os.Getenv("MUSICLE_INTEGRATION") == "" {
		t.Skip("MUSICLE_INTEGRATION=1 set to run network download test")
	}
	ext, err := Extract()
	if err != nil {
		t.Skipf("engine assets not embedded in this build: %v", err)
	}
	if ext == nil || ext.YTDLP == "" || ext.FFMPEG == "" {
		t.Fatal("Extract returned incomplete locations")
	}

	out := t.TempDir()
	err = DownloadSpotify(SpotifyOptions{
		OutputDir: out,
		URLs:      []string{"https://open.spotify.com/track/4uLU6hMCjMI75M1A2tKUQC"},
		Progress:  func(pct int, msg string) { t.Logf("%3d%% %s", pct, msg) },
	})
	if err != nil {
		t.Fatalf("DownloadSpotify: %v", err)
	}

	mp3s, _ := filepath.Glob(filepath.Join(out, "*.mp3"))
	if len(mp3s) == 0 {
		t.Fatal("engin üretilen mp3 bulunamadı")
	}
	data, err := os.ReadFile(mp3s[0])
	if err != nil {
		t.Fatal(err)
	}
	if len(data) < 2_000_000 {
		t.Fatalf("mp3 şüpheli şekilde küçük (%d bayt)", len(data))
	}
	if !bytes.HasPrefix(data, []byte("ID3")) {
		t.Errorf("mp3 ID3 başlığı taşımıyor")
	}
	if !bytes.Contains(data, []byte("APIC")) {
		t.Errorf("mp3 APIC kapak çerçevesi içermiyor")
	}
	t.Logf("üretilen dosya: %s (%d bayt, ID3+APIC)", filepath.Base(mp3s[0]), len(data))
}
