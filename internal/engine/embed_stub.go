//go:build !engine_assets

package engine

import (
	"embed"
	"errors"
)

// ErrNoEmbeddedAssets is returned by Extract() when the binary is built without
// the engine_assets build tag (the embedded tools were not bundled at compile
// time). The bridge then falls back to the legacy pipeline.
var ErrNoEmbeddedAssets = errors.New("engine: gömülü araçlar bulunamadı (build sırasında 'make build' veya -tags engine_assets kullanın)")

// binFS is declared here as an empty filesystem so extract.go compiles in the
// non-tagged build. Extract() short-circuits with ErrNoEmbeddedAssets.
var binFS embed.FS

func probeEmbedded() error {
	return ErrNoEmbeddedAssets
}
