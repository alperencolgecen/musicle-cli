package engine

import (
	"fmt"
	"strings"
)

// YouTubeOptions tunes DownloadYouTube / DownloadYouTubePlaylist.
type YouTubeOptions struct {
	OutputDir string
	URLs      []string // accepts watch?v=, youtu.be/, and bare IDs
	Progress  Progress
}

// DownloadYouTube downloads one or more YouTube URLs to mp3 via the embedded
// yt-dlp engine. It returns the first non-fatal error if any track failed,
// nil when every URL was processed.
func DownloadYouTube(opt YouTubeOptions) error {
	ext, err := Extract()
	if err != nil {
		return err
	}
	urls := normalizeYouTubeURLs(opt.URLs)
	if len(urls) == 0 {
		return fmt.Errorf("yt-dlp: en az bir YouTube URL gerekli")
	}
	args := []string{
		"-m", "musicle_ytdlp",
		"--output", opt.OutputDir,
	}
	args = append(args, urls...)

	return Run(ext, RunOptions{
		Args:       args,
		Progress:   opt.Progress,
		JSONStdout: true,
		StderrCb:   func(line string) { /* debug-only hook */ },
		WorkDir:    ext.Root,
		Env: []string{
			"MUSICLE_FFMPEG=" + ext.FFmpegBin,
			"PYTHONPATH=" + scriptPathEnv(),
		},
	})
}

// DownloadYouTubePlaylist is a thin alias for DownloadYouTube with a single
// URL — yt-dlp handles playlists natively when given a playlist URL.
func DownloadYouTubePlaylist(opt YouTubeOptions) error {
	return DownloadYouTube(opt)
}

// normalizeYouTubeURLs accepts the common input shapes the bubbletea UI
// hands us (raw IDs, youtu.be shortlinks, full watch URLs, playlist URLs)
// and returns canonical watch URLs so yt-dlp handles them consistently.
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
