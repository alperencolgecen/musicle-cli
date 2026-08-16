//go:build engine_assets

package engine

import (
	"embed"
	"errors"
	"io/fs"
)

// ErrNoEmbeddedAssets is returned by Extract() when the binary is built without
// the engine_assets build tag (the embedded tools were not bundled at compile
// time). The bridge then falls back to the legacy pipeline.
var ErrNoEmbeddedAssets = errors.New("engine: gömülü araçlar bulunamadı (build sırasında 'make build' veya -tags engine_assets kullanın)")

// probeEmbedded verifies the embedded tools filesystem is usable at runtime.
func probeEmbedded() error {
	if _, err := fs.ReadFile(binFS, "engine_bin/.engine-stamp"); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return errors.New("engine: gömülü araçlar bozuk (.engine-stamp eksik)")
		}
		return err
	}
	return nil
}

//go:embed all:engine_bin
var binFS embed.FS
