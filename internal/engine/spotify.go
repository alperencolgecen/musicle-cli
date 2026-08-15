package engine

import (
	"fmt"
	"path/filepath"
	"strings"
)

// SpotifyOptions tunes DownloadSpotify / DownloadSpotifyPlaylist.
type SpotifyOptions struct {
	OutputDir string
	URLs      []string // track, album, playlist, or search queries
	AsQuery   bool     // treat each URL as a free-text search query
	Progress  Progress
}

// DownloadSpotify downloads Spotify URLs (track / album / playlist) via the
// embedded spotdl engine. spotdl detects the URL type internally so a single
// helper covers all three cases plus search queries.
func DownloadSpotify(opt SpotifyOptions) error {
	ext, err := Extract()
	if err != nil {
		return err
	}
	urls := sanitizeSpotifyURLs(opt.URLs)
	if len(urls) == 0 {
		return fmt.Errorf("spotdl: en az bir Spotify URL veya sorgu gerekli")
	}

	args := []string{
		"-m", "musicle_spotdl",
		"--output", opt.OutputDir,
	}
	if opt.AsQuery {
		args = append(args, "--query")
	}
	args = append(args, urls...)

	return Run(ext, RunOptions{
		Args:       args,
		Progress:   opt.Progress,
		JSONStdout: true,
		WorkDir:    ext.Root,
		Env: []string{
			"MUSICLE_FFMPEG=" + ext.FFmpegBin,
			"PYTHONPATH=" + scriptPathEnv(),
		},
	})
}

// DownloadSpotifyPlaylist is the public name for playlist / album flows.
// Kept as a separate entry point so the bridge layer can dispatch on intent.
func DownloadSpotifyPlaylist(opt SpotifyOptions) error {
	return DownloadSpotify(opt)
}

// sanitizeSpotifyURLs trims whitespace and drops empty entries; spotdl is
// tolerant of any URL shape, so we don't normalise the URI itself.
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

// scriptPathEnv returns a PYTHONPATH value that lets the embedded python
// find musicle_ytdlp.py / musicle_spotdl.py / _musicle_helpers.py.
//
// The helpers live alongside the extracted venv at <Root>/scripts. The venv
// itself is the active interpreter, so its site-packages already resolve;
// we only need to expose our helper directory.
func scriptPathEnv() string {
	if extracted == nil {
		// Best-effort fallback if scriptPathEnv is invoked before Extract().
		// The Run() that follows will surface a clearer error.
		return "scripts"
	}
	return filepath.Join(extracted.Root, "scripts") //nolint:gosec // computed from cacheRoot
}
