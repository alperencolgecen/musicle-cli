//go:build !engine_assets

package engine

import (
	"embed"
	"errors"
)

// ErrNoEmbeddedAssets is returned by Extract() when the binary is built without
// the engine_assets build tag (i.e. the embedded Python engine was not bundled
// at compile time). The bridge's runEngine() then falls back to the legacy
// pipeline.
var ErrNoEmbeddedAssets = errors.New("engine: gömülü motor varlıkları bulunamadı (build sırasında 'make build' veya -tags engine_assets kullanın)")

// venvFS/ffmpegFS are declared here as empty filesystems so extract.go compiles
// in the non-tagged build. Extract() never reaches doExtract() because
// probeEmbedded() short-circuits with ErrNoEmbeddedAssets.
var venvFS embed.FS

var ffmpegFS embed.FS

func probeEmbedded() error {
	return ErrNoEmbeddedAssets
}
