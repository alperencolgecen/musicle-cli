//go:build !engine_assets

// Stub build used during `go vet` and quick smoke tests when the embedded
// payload directories are absent. In a release build the build tag
// `engine_assets` is set by build.sh via -tags=engine_assets, which switches
// to embed_assets.go and binds the real assets.

package engine

import "errors"

// ErrNoEmbeddedAssets is returned when the binary was compiled without the
// engine_assets build tag. This happens when the engine payload was not
// prepared before the Go build.
var ErrNoEmbeddedAssets = errors.New("engine: gömülü motor paketi derlenmedi (build tag engine_assets gerekli)")

// venvFS / ffmpegFS are nil here; any caller that hits them will panic.
// Real usage always flows through Extract(), which checks the build tag
// first and returns ErrNoEmbeddedAssets before touching the FS variables.
