#!/usr/bin/env bash
# prepare-engine.sh — fetch the self-contained download engine for musicle-cli.
#
# Downloads a standalone yt-dlp binary plus a static ffmpeg into
# internal/engine/engine_bin/ (no Python, no venv). `go:embed` (see
# internal/engine/embed_assets.go) then bakes them into the binary when built
# with -tags engine_assets.
#
# Usage:
#   ./scripts/prepare-engine.sh [linux|windows|darwin]
#
# The chosen OS is the host where the binaries will run.

set -euo pipefail

OS_TARGET="${1:-linux}"
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BIN_DIR="$ROOT/internal/engine/engine_bin"

mkdir -p "$BIN_DIR"

echo "==> Target OS: $OS_TARGET"
echo "==> Engine bin dir: $BIN_DIR"

# ---------- yt-dlp (standalone, no Python runtime required) ----------
fetch_ytdlp() {
  local out="$BIN_DIR/yt-dlp"
  if [ -x "$out" ]; then
    echo "==> Reusing existing yt-dlp"
    return
  fi
  echo "==> Downloading yt-dlp for $OS_TARGET..."
  case "$OS_TARGET" in
    linux)   URL="https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux" ;;
    windows) URL="https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe" ;;
    darwin)  URL="https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_macos" ;;
    *) echo "Unknown OS: $OS_TARGET" >&2; exit 1 ;;
  esac
  curl -sL "$URL" -o "$out"
  chmod +x "$out"
  echo "    yt-dlp: $("$out" --version 2>/dev/null || echo '?')"
}

# ---------- ffmpeg (static, only what's used for MP3 muxing) ----------
fetch_ffmpeg() {
  local ffmpeg_out="$BIN_DIR/ffmpeg"
  if [ -x "$ffmpeg_out" ]; then
    echo "==> Reusing existing ffmpeg"
    return
  fi
  echo "==> Downloading static ffmpeg for $OS_TARGET..."
  local TMP
  TMP="$(mktemp -d)"
  case "$OS_TARGET" in
    linux)
      URL="https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz"
      curl -sL "$URL" -o "$TMP/ffmpeg.tar.xz"
      tar -xJf "$TMP/ffmpeg.tar.xz" -C "$TMP"
      cp "$(find "$TMP" -type f -name ffmpeg | head -1)" "$ffmpeg_out"
      ;;
    windows)
      URL="https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip"
      curl -sL "$URL" -o "$TMP/ffmpeg.zip"
      unzip -q "$TMP/ffmpeg.zip" -d "$TMP"
      cp "$(find "$TMP" -type f -iname 'ffmpeg.exe' | head -1)" "$ffmpeg_out"
      ;;
    darwin)
      curl -sL "https://evermeet.cx/ffmpeg/getrelease/zip" -o "$TMP/ffmpeg.zip"
      unzip -q "$TMP/ffmpeg.zip" -d "$BIN_DIR"
      ;;
    *)
      echo "Unknown OS: $OS_TARGET" >&2
      rm -rf "$TMP"
      exit 1
      ;;
  esac
  chmod +x "$ffmpeg_out"
  rm -rf "$TMP"
  echo "    ffmpeg: $("$ffmpeg_out" -version 2>/dev/null | head -1 || echo '?')"
}

fetch_ytdlp
fetch_ffmpeg

# Compress the embedded binaries when UPX is available. Static ffmpeg packs
# ~31%, which keeps the final binary small; yt-dlp is already dense and UPX
# reports no gain, so only ffmpeg is processed. If UPX is missing the build
# still works, just with a larger binary.
if command -v upx >/dev/null 2>&1; then
  echo "==> UPX sıkıştırması (ffmpeg)..."
  upx -q --best "$BIN_DIR/ffmpeg" 2>/dev/null || echo "    ffmpeg UPX başarısız, atlanıyor"
else
  echo "==> upx bulunamadı; ffmpeg sıkıştırılmadan bırakılıyor (daha büyük binary)"
fi

# Stamp the cache version so the Go side can verify freshness. Must match the
# cacheVersion constant in internal/engine/extract.go.
echo "engine-v2" > "$BIN_DIR/.engine-stamp"

echo "==> Done. Embedded engine ready at $BIN_DIR (baked in via -tags engine_assets)."
