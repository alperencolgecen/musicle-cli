package engine

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// Extracted holds the on-disk paths to a fully extracted engine. It is the
// value returned by Extract() and is what the wrappers operate on.
type Extracted struct {
	Root         string // cache root, e.g. ~/.cache/musicle/engine-v1
	PythonBin    string // absolute path to the embedded python interpreter
	SpotdlBin    string // absolute path to the spotdl console script wrapper
	YtdlpBin     string // absolute path to the yt-dlp console script wrapper
	FFmpegBin    string // absolute path to the ffmpeg binary
	FFprobeBin   string // absolute path to the ffprobe binary
	SitePackages string // PYTHONPATH entry: venv/lib/python*/site-packages
}

// ErrNoEmbeddedAssets is defined in embed_assets.go (always compiled) so the
// same identifier is available to all callers. It is returned by Extract() when
// the embedded engine assets are missing or corrupt.

// cacheVersion must match the value written by scripts/prepare-engine.sh.
const cacheVersion = "engine-v1"

// sentinelName marks a successful extract; its presence short-circuits Extract.
const sentinelName = ".extracted"

var (
	extractOnce sync.Once
	extracted   *Extracted
	extractErr  error
)

// Extract returns the on-disk locations of the embedded engine, extracting
// the embedded venv + ffmpeg to UserCacheDir on first call. Subsequent calls
// return the cached result.
//
// If the embedded engine assets are absent or corrupt, Extract returns
// ErrNoEmbeddedAssets.
func Extract() (*Extracted, error) {
	extractOnce.Do(func() {
		extracted, extractErr = doExtract()
	})
	return extracted, extractErr
}

// MustExtract is like Extract but panics on error. Intended for use after a
// successful startup probe.
func MustExtract() *Extracted {
	e, err := Extract()
	if err != nil {
		panic(err)
	}
	return e
}

func doExtract() (*Extracted, error) {
	// Probe the embedded FS variables. Without the build tag the stub
	// implementation returns ErrNoEmbeddedAssets.
	if err := probeEmbedded(); err != nil {
		return nil, err
	}

	cacheRoot, err := cacheDir()
	if err != nil {
		return nil, err
	}

	root := filepath.Join(cacheRoot, cacheVersion)

	// Fast path: already extracted.
	if valid, _ := isExtracted(root); valid {
		return buildExtracted(root)
	}

	// Slow path: wipe and re-extract.
	if err := os.RemoveAll(root); err != nil {
		return nil, fmt.Errorf("engine: eski önbellek temizlenemedi: %w", err)
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, fmt.Errorf("engine: cache dizini oluşturulamadı: %w", err)
	}

	venvDst := filepath.Join(root, "venv")
	ffmpegDst := filepath.Join(root, "ffmpeg")
	scriptsDst := filepath.Join(root, "scripts")

	if err := copyFS(venvFS, "engine_venv", venvDst); err != nil {
		return nil, fmt.Errorf("engine: venv çıkarılamadı: %w", err)
	}
	if err := copyFS(ffmpegFS, "engine_ffmpeg", ffmpegDst); err != nil {
		return nil, fmt.Errorf("engine: ffmpeg çıkarılamadı: %w", err)
	}
	// Helper scripts must be on disk so the embedded python can import them.
	if err := copyFS(scriptsFS, "scripts", scriptsDst); err != nil {
		return nil, fmt.Errorf("engine: scriptler çıkarılamadı: %w", err)
	}

	// chmod +x on ffmpeg binaries on unix.
	if runtime.GOOS != "windows" {
		_ = os.Chmod(filepath.Join(ffmpegDst, "ffmpeg"), 0o755)
		_ = os.Chmod(filepath.Join(ffmpegDst, "ffprobe"), 0o755)
		_ = os.Chmod(filepath.Join(venvDst, "bin", "python"), 0o755)
		_ = os.Chmod(filepath.Join(venvDst, "bin", "spotdl"), 0o755)
		_ = os.Chmod(filepath.Join(venvDst, "bin", "yt-dlp"), 0o755)
	}

	// Stamp the cache as valid.
	if err := os.WriteFile(filepath.Join(root, sentinelName), []byte(cacheVersion+"\n"), 0o644); err != nil {
		return nil, fmt.Errorf("engine: sentinel yazılamadı: %w", err)
	}

	return buildExtracted(root)
}

func cacheDir() (string, error) {
	if d := os.Getenv("MUSICLE_CACHE_DIR"); d != "" {
		return d, nil
	}
	base, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("engine: kullanıcı cache dizini bulunamadı: %w", err)
	}
	return filepath.Join(base, "musicle"), nil
}

func isExtracted(root string) (bool, error) {
	data, err := os.ReadFile(filepath.Join(root, sentinelName))
	if err != nil {
		return false, err
	}
	return string(data) == cacheVersion+"\n", nil
}

func buildExtracted(root string) (*Extracted, error) {
	venv := filepath.Join(root, "venv")
	ffmpeg := filepath.Join(root, "ffmpeg")

	py := venvPython(venv)
	spotdl := venvBin(venv, "spotdl")
	ytdlp := venvBin(venv, "yt-dlp")

	for _, p := range []string{py, spotdl, ytdlp} {
		if _, err := os.Stat(p); err != nil {
			return nil, fmt.Errorf("engine: beklenen binary eksik: %s (%w)", p, err)
		}
	}

	return &Extracted{
		Root:         root,
		PythonBin:    py,
		SpotdlBin:    spotdl,
		YtdlpBin:     ytdlp,
		FFmpegBin:    filepath.Join(ffmpeg, "ffmpeg"+exeSuffix()),
		FFprobeBin:   filepath.Join(ffmpeg, "ffprobe"+exeSuffix()),
		SitePackages: venvSitePackages(venv),
	}, nil
}

func venvPython(venv string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(venv, "Scripts", "python.exe")
	}
	return filepath.Join(venv, "bin", "python")
}

func venvBin(venv, name string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(venv, "Scripts", name+".exe")
	}
	return filepath.Join(venv, "bin", name)
}

func venvSitePackages(venv string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(venv, "Lib", "site-packages")
	}
	// unix: lib/pythonX.Y/site-packages — pick whichever exists.
	candidates, _ := filepath.Glob(filepath.Join(venv, "lib", "python*", "site-packages"))
	if len(candidates) > 0 {
		return candidates[0]
	}
	return filepath.Join(venv, "lib", "site-packages")
}

func exeSuffix() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}

// copyFS recursively copies an embed.FS subtree rooted at srcRoot into dst.
// Permissions follow fs.FileInfo when available; directories are 0o755,
// regular files are 0o644 (the caller chmods executables afterwards).
func copyFS(files fs.FS, srcRoot, dst string) error {
	return fs.WalkDir(files, srcRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		// Skip the synthetic root.
		if path == srcRoot {
			return nil
		}
		rel, err := filepath.Rel(srcRoot, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		target := filepath.Join(dst, filepath.FromSlash(rel))

		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		// Regular file (or symlink): copy bytes.
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		in, err := files.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		info, err := d.Info()
		if err != nil {
			info = nil
		}
		mode := os.FileMode(0o644)
		if info != nil && info.Mode()&0o111 != 0 {
			mode = 0o755
		}

		out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
		if err != nil {
			return err
		}
		if _, err := io.Copy(out, in); err != nil {
			out.Close()
			return err
		}
		return out.Close()
	})
}
