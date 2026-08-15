package bridge

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"MusicLeCLI/bridge/download/music"
	"MusicLeCLI/bridge/download/playlist"
	"MusicLeCLI/internal/engine"
)

// engineDisabled reports whether the user opted out of the embedded engine
// for this run (via env var). Useful for debugging the legacy pipeline.
func engineDisabled() bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv("MUSICLE_NO_ENGINE")))
	return v == "1" || v == "true" || v == "yes"
}

// isSpotifyURL heuristically detects Spotify links/URIs so the unified
// engine dispatcher can route to spotdl instead of yt-dlp.
func isSpotifyURL(url string) bool {
	u := strings.ToLower(strings.TrimSpace(url))
	return strings.Contains(u, "spotify.com") || strings.HasPrefix(u, "spotify:")
}

// isYouTubeURL detects YouTube links so we don't hand a YouTube URL to spotdl.
func isYouTubeURL(url string) bool {
	u := strings.ToLower(strings.TrimSpace(url))
	return strings.Contains(u, "youtube.com") || strings.Contains(u, "youtu.be")
}

// downloadYouTube downloads a YouTube URL using the embedded engine when
// available, falling back to the legacy pure-Go pipeline if the engine
// binary was not embedded (e.g. dev builds without the engine_assets tag).
func downloadYouTube(url, outputDir string) *Result {
	if url == "" {
		return &Result{Status: "error", Error: "invalid URL"}
	}

	CurrentDownload.Set(true, 0, "Bridge: YouTube indirme başlatılıyor...")

	if res, ok := runEngine(url, outputDir); ok {
		return res
	}

	filePath, err := music.DownloadYouTubeToFile(url, outputDir, func(pct int, msg string) {
		CurrentDownload.Set(true, float64(pct), fmt.Sprintf("YouTube: %s", msg))
	})
	if err != nil {
		CurrentDownload.Set(false, 0, fmt.Sprintf("Error: %v", err))
		return &Result{Status: "error", Error: fmt.Sprintf("YouTube download failed: %v", err)}
	}

	CurrentDownload.Set(false, 100, fmt.Sprintf("Saved: %s", filepath.Base(filePath)))

	meta := extractMetadata(filePath)
	if meta.Status == "error" {
		return &Result{
			Status:   "ok",
			Filename: filepath.Base(filePath),
			Title:    strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath)),
			Artist:   "Unknown",
		}
	}
	return meta
}

// runEngine is the unified yt-dlp / spotdl dispatcher. It tries the embedded
// engine first, routing by URL type, and returns (result, true) when the
// engine handled the URL. On any engine error or when the engine binary is
// not embedded it returns (nil, false) so the caller can fall back.
func runEngine(url, outputDir string) (*Result, bool) {
	if engineDisabled() {
		return nil, false
	}
	if _, err := engine.Extract(); err != nil {
		CurrentDownload.Set(true, 0, "Bridge: motor yok, eski yöntem kullanılıyor")
		return nil, false
	}

	if isSpotifyURL(url) {
		CurrentDownload.Set(true, 0, "Bridge: spotdl motoru başlatılıyor...")
		err := engine.DownloadSpotify(engine.SpotifyOptions{
			OutputDir: outputDir,
			URLs:      []string{url},
			Progress:  engineProgress,
		})
		if err != nil {
			CurrentDownload.Set(false, 0, fmt.Sprintf("Motor hatası, eski yöntem deneniyor: %v", err))
			return nil, false
		}
	} else {
		CurrentDownload.Set(true, 0, "Bridge: yt-dlp motoru başlatılıyor...")
		err := engine.DownloadYouTube(engine.YouTubeOptions{
			OutputDir: outputDir,
			URLs:      []string{url},
			Progress:  engineProgress,
		})
		if err != nil {
			CurrentDownload.Set(false, 0, fmt.Sprintf("Motor hatası, eski yöntem deneniyor: %v", err))
			return nil, false
		}
	}

	CurrentDownload.Set(false, 100, "Motor ile tamamlandı")
	return &Result{Status: "ok", Message: "gömülü motor ile indirildi"}, true
}

// engineProgress forwards the engine's progress events to the shared UI state.
func engineProgress(pct int, msg string) {
	CurrentDownload.Set(true, float64(pct), msg)
}

// downloadSpotify downloads a Spotify URL. Handles both track and playlist URLs.
func downloadSpotify(url, outputDir string) *Result {
	if url == "" {
		return &Result{Status: "error", Error: "invalid URL"}
	}

	CurrentDownload.Set(true, 0, "Starting...")

	if res, ok := runEngine(url, outputDir); ok {
		return res
	}

	urlLower := strings.ToLower(url)
	isPlaylist := strings.Contains(urlLower, "/playlist/") || strings.Contains(urlLower, "spotify:playlist:")
	isAlbum := strings.Contains(urlLower, "/album/") || strings.Contains(urlLower, "spotify:album:")
	isTrack := strings.Contains(urlLower, "/track/") || strings.Contains(urlLower, "spotify:track:")

	if isPlaylist || isAlbum {
		return downloadSpotifyPlaylist(url, outputDir)
	}

	if !isTrack {
		// Try collection (playlist/album) detection as fallback
		return downloadSpotifyPlaylist(url, outputDir)
	}

	// Single track: fetch Spotify metadata, search YouTube, download
	spTrack, err := music.FetchSpotifyTrack(url)
	if err != nil {
		CurrentDownload.Set(false, 0, "Error")
		return &Result{Status: "error", Error: fmt.Sprintf("Spotify metadata: %v", err)}
	}

	query := spTrack.Artist + " - " + spTrack.Title
	CurrentDownload.Set(true, 5, fmt.Sprintf("Searching YouTube: %s", query))

	videoID, _, err := music.SearchYouTubeTrack(query)
	if err != nil {
		CurrentDownload.Set(false, 0, "Error")
		return &Result{Status: "error", Error: fmt.Sprintf("YouTube search: %v", err)}
	}

	CurrentDownload.Set(true, 10, "Downloading from YouTube...")

	// Try direct MP3 download via yt-dlp + ffmpeg first (most reliable)
	_, mp3Data, err := music.DownloadYouTubeTrackDirectMP3(videoID, func(pct int, msg string) {
		CurrentDownload.Set(true, 10+float64(pct)*80/100, msg)
	})
	if err == nil && mp3Data != nil {
		CurrentDownload.Set(true, 90, "Writing ID3 tag...")

		spTrack.StreamURL = "https://www.youtube.com/watch?v=" + videoID
		spTrack.Format = "mp3"
		spTrack.ContentLen = int64(len(mp3Data))

		tagged, err := music.SaveRawAsMP3(mp3Data, spTrack, outputDir, func(pct int, msg string) {
			CurrentDownload.Set(true, 90+float64(pct)*10/100, msg)
		})
		if err != nil {
			CurrentDownload.Set(false, 0, "Error")
			return &Result{Status: "error", Error: fmt.Sprintf("tag: %v", err)}
		}

		CurrentDownload.Set(false, 100, "Done")
		return &Result{
			Status:   "ok",
			Filename: filepath.Base(tagged),
			Title:    spTrack.Title,
			Artist:   spTrack.Artist,
			Duration: spTrack.DurationSec,
		}
	}

	// Fallback: download raw WebM and convert
	CurrentDownload.Set(true, 10, "Direct MP3 failed, trying raw download...")
	_, rawAudio, err := music.DownloadYouTubeTrack(videoID, func(pct int, msg string) {
		CurrentDownload.Set(true, 10+float64(pct)*30/100, msg)
	})
	if err != nil {
		CurrentDownload.Set(false, 0, "Error")
		return &Result{Status: "error", Error: fmt.Sprintf("YouTube download: %v", err)}
	}

	CurrentDownload.Set(true, 40, "Converting to MP3...")

	spTrack.StreamURL = "https://www.youtube.com/watch?v=" + videoID
	spTrack.Format = "webm"
	spTrack.ContentLen = int64(len(rawAudio))

	filePath, err := music.SaveRawAsMP3(rawAudio, spTrack, outputDir, func(pct int, msg string) {
		CurrentDownload.Set(true, 40+float64(pct)*60/100, msg)
	})
	if err != nil {
		CurrentDownload.Set(false, 0, "Error")
		return &Result{Status: "error", Error: fmt.Sprintf("convert: %v", err)}
	}

	CurrentDownload.Set(false, 100, "Done")

	return &Result{
		Status:   "ok",
		Filename: filepath.Base(filePath),
		Title:    spTrack.Title,
		Artist:   spTrack.Artist,
		Duration: spTrack.DurationSec,
	}
}

// downloadSpotifyPlaylist also routes through the embedded engine first.
func downloadSpotifyPlaylist(spotifyURL, outputDir string) *Result {
	if res, ok := runEngine(spotifyURL, outputDir); ok {
		return res
	}

	files, err := playlist.DownloadSpotifyPlaylist(spotifyURL, outputDir, func(pct int, msg string) {
		CurrentDownload.Set(true, float64(pct), msg)
	})
	if err != nil {
		CurrentDownload.Set(false, 0, "Error")
		return &Result{Status: "error", Error: err.Error()}
	}

	CurrentDownload.Set(false, 100, fmt.Sprintf("Downloaded %d songs", len(files)))

	songs := make([]Result, 0, len(files))
	for _, f := range files {
		meta := extractMetadata(f)
		if meta.Status == "ok" {
			songs = append(songs, *meta)
		}
	}

	return &Result{
		Status:  "ok",
		Message: fmt.Sprintf("Downloaded %d song(s)", len(songs)),
		Songs:   songs,
	}
}
