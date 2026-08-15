package engine

import "embed"

// scriptsFS holds the helper Python modules (musicle_ytdlp.py,
// musicle_spotdl.py, _musicle_helpers.py) that the embedded interpreter
// runs via `python -m <module>`. Unlike the venv/ffmpeg payloads these are
// plain source checked into the repo, so they are embedded unconditionally
// (no build tag) and always available.
//
// At runtime doExtract() writes them to <Root>/scripts and scriptPathEnv()
// points PYTHONPATH there.
//
//go:embed all:scripts
var scriptsFS embed.FS
