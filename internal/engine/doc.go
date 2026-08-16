// Package engine bundles a fully self-contained download engine (yt-dlp +
// static ffmpeg/ffprobe) into the musicle-cli binary via go:embed.
//
// The flow:
//
//  1. scripts/prepare-engine.sh downloads yt-dlp and a static ffmpeg/ffprobe
//     into internal/engine/engine_bin/ (no Python, no venv), then go:embed
//     bakes them into the binary with -tags engine_assets.
//
//  2. At startup, Extract() unpacks those binaries into
//     UserCacheDir/musicle/engine-v2/ (a single-shot copy guarded by a
//     sentinel file). Subsequent runs reuse the cache.
//
//  3. The wrappers in this package invoke the embedded yt-dlp to fetch audio
//     and the embedded ffmpeg to mux it into a 320k MP3 with an embedded
//     cover (APIC). Spotify URLs are resolved by yt-dlp itself, so no Spotify
//     API or token is involved.
package engine
