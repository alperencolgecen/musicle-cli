package engine

// Progress is the callback signature used by wrappers to stream download
// progress back to the Go side. pct is 0..100, msg is human-readable.
type Progress func(pct int, msg string)

// YouTubeOptions tunes DownloadYouTube / DownloadYouTubePlaylist.
type YouTubeOptions struct {
	OutputDir string
	URLs      []string // accepts watch?v=, youtu.be/, and bare IDs
	Progress  Progress
}

// SpotifyOptions tunes DownloadSpotify / DownloadSpotifyPlaylist.
type SpotifyOptions struct {
	OutputDir string
	URLs      []string // track, album, playlist, or search queries
	Progress  Progress
}
