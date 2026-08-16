package engine

import (
	"fmt"
	"path/filepath"
	"strings"
)

// DownloadSpotify downloads Spotify URLs (track / album / playlist) to mp3
// (320k) with an embedded cover. yt-dlp resolves the Spotify page and fetches
// the audio from YouTube, so this reuses the same pipeline as YouTube — no
// Spotify API or token is involved.
func DownloadSpotify(opt SpotifyOptions) error {
	ext, err := Extract()
	if err != nil {
		return err
	}
	urls := sanitizeSpotifyURLs(opt.URLs)
	if len(urls) == 0 {
		return fmt.Errorf("yt-dlp: en az bir Spotify URL veya sorgu gerekli")
	}
	return downloadWithYTDLP(ext, urls, opt.OutputDir, opt.Progress)
}

// DownloadSpotifyPlaylist is the public name for playlist / album flows.
// Kept as a separate entry point so the bridge layer can dispatch on intent.
func DownloadSpotifyPlaylist(opt SpotifyOptions) error {
	return DownloadSpotify(opt)
}

// sanitizeSpotifyURLs trims whitespace and drops empty entries.
func sanitizeSpotifyURLs(in []string) []string {
	out := make([]string, 0, len(in))
	for _, raw := range in {
		u := strings.TrimSpace(raw)
		if u == "" {
			continue
		}
		out = append(out, u)
	}
	return out
}

// scriptPathEnv is retained only to avoid breaking callers that referenced it;
// the new engine has no Python helpers, so it simply returns the cache root.
func scriptPathEnv() string {
	if extracted == nil {
		return "engine_bin"
	}
	return filepath.Join(extracted.Root, "engine_bin")
}
