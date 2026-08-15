// Package engine bundles a fully self-contained Python download engine
// (Python interpreter + spotdl + yt-dlp + static ffmpeg) into the musicle-cli
// binary via go:embed.
//
// The flow:
//
//  1. `scripts/prepare-engine.sh` builds a venv at assets/engine/venv and
//     downloads a static ffmpeg into assets/engine/ffmpeg/<os>-<arch>/,
//     then copies both trees to internal/engine/engine_{venv,ffmpeg}/ so
//     `go:embed` can pick them up.
//
//  2. At startup, Extract() unpacks those trees into
//     UserCacheDir/musicle/engine-v1/ (a single-shot copy guarded by a
//     sentinel file). Subsequent runs reuse the cache.
//
//  3. The wrappers in this package invoke the embedded python with the
//     right arguments (spotdl / yt-dlp) and stream progress as JSON.
package engine
