//go:build !engine_assets

// Stub build used during `go vet` and quick smoke tests when the embedded
// payload directories are absent. In a release build the build tag
// `engine_assets` is set by build.sh via -tags=engine_assets, which switches
// to embed_assets.go and binds the real assets.

package engine

import (
	"errors"
	"io/fs"
)

var ErrNoEmbeddedAssets = errors.New("engine: gömülü motor paketi derlenmedi (build tag engine_assets gerekli)")

// probeEmbedded always fails in the stub build so the wrapper code can
// gracefully return ErrNoEmbeddedAssets instead of touching undefined FSes.
func probeEmbedded() error {
	return ErrNoEmbeddedAssets
}

// venvFS and ffmpegFS are unused in the stub build; they are placeholders so
// other code in the package that references them still type-checks under
// `go vet` (these references are dead code that probeEmbedded() guards).
var (
	venvFS   = emptyFS{}
	ffmpegFS = emptyFS{}
)

// emptyFS satisfies fs.FS but returns ErrNotExist for any path. It only
// exists so stub builds compile.
type emptyFS struct{}

func (emptyFS) Open(_ string) (fs.File, error) { return nil, fs.ErrNotExist }
func (emptyFS) ReadDir(_ string) ([]fs.DirEntry, error) {
	return nil, fs.ErrNotExist
}
