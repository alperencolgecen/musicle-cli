package engine

import (
	"fmt"
	"strings"
)

// DownloadYouTube downloads one or more YouTube URLs to mp3 (320k) with an
// embedded cover via the bundled yt-dlp + ffmpeg. It returns the first
// non-fatal error if any URL failed, nil when every URL was processed.
func DownloadYouTube(opt YouTubeOptions) error {
	ext, err := Extract()
	if err != nil {
		return err
	}
	urls := normalizeYouTubeURLs(opt.URLs)
	if len(urls) == 0 {
		return fmt.Errorf("yt-dlp: en az bir YouTube URL gerekli")
	}
	return downloadWithYTDLP(ext, urls, opt.OutputDir, opt.Progress)
}

// DownloadYouTubePlaylist is a thin alias for DownloadYouTube with a single
// URL — yt-dlp handles playlists natively when given a playlist URL.
func DownloadYouTubePlaylist(opt YouTubeOptions) error {
	return DownloadYouTube(opt)
}

// normalizeYouTubeURLs trims whitespace and drops empty entries; yt-dlp is
// tolerant of any URL shape, so we don't rewrite the URL itself.
func normalizeYouTubeURLs(in []string) []string {
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
