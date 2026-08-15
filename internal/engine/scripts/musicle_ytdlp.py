"""yt-dlp wrapper used by musicle-cli's embedded download engine.

Invocation (Go side builds the argv):

    python -m musicle_ytdlp --output DIR URL [URL ...]

Behaviour:
    * Extracts audio as MP3 192 kbps via the embedded ffmpeg.
    * Embeds ID3 tags (title, artist, album, track number).
    * Embeds thumbnail as cover art when available.
    * Streams JSON progress events on stdout for the Go bridge.

Exit codes:
    0  success
    1  any track failed; remaining tracks were skipped
    2  fatal setup error (no ffmpeg, no URLs, ...)
"""

from __future__ import annotations

import argparse
import os
import sys
import traceback
from pathlib import Path
from typing import List

# Make sibling helpers importable regardless of how Python is invoked.
sys.path.insert(0, str(Path(__file__).resolve().parent))

from _musicle_helpers import (  # noqa: E402
    done,
    emit,
    error,
    start,
    track_done,
)


def _resolve_ffmpeg() -> str | None:
    """Locate ffmpeg: prefer $MUSICLE_FFMPEG, then sibling ./ffmpeg, then PATH."""
    env = os.environ.get("MUSICLE_FFMPEG")
    if env and Path(env).is_file():
        return env
    here = Path(__file__).resolve().parent
    for cand in ("ffmpeg", "ffmpeg.exe"):
        p = here / cand
        if p.is_file():
            return str(p)
    return None


def _run_one(url: str, output_dir: str, ffmpeg: str | None, index: int, total: int) -> bool:
    try:
        import yt_dlp  # type: ignore
    except ImportError as exc:
        error(f"yt-dlp yüklenemedi: {exc}")
        return False

    start(f"yt-dlp: {url}")
    outtmpl = str(Path(output_dir) / "%(title)s [%(id)s].%(ext)s")
    ydl_opts = {
        "format": "bestaudio/best",
        "outtmpl": outtmpl,
        "noplaylist": True,
        "quiet": True,
        "no_warnings": True,
        "writethumbnail": False,
        "postprocessors": [
            {
                "key": "FFmpegExtractAudio",
                "preferredcodec": "mp3",
                "preferredquality": "192",
            },
            {"key": "FFmpegMetadata"},
            {"key": "EmbedThumbnail"},
        ],
        "progress_hooks": [],
    }
    if ffmpeg:
        ydl_opts["ffmpeg_location"] = ffmpeg

    last_pct = [-1]

    def hook(d: dict) -> None:
        status = d.get("status")
        if status == "downloading":
            total_bytes = d.get("total_bytes") or d.get("total_bytes_estimate") or 0
            downloaded = d.get("downloaded_bytes") or 0
            pct = (downloaded * 100.0 / total_bytes) if total_bytes else -1.0
            # Downsample to integer percent changes to avoid spamming stdout.
            ipct = int(pct)
            if ipct == last_pct[0]:
                return
            last_pct[0] = ipct
            msg = f"[{index}/{total}] İndiriliyor %{ipct}"
            emit("progress", message=msg, percent=pct)
        elif status == "finished":
            emit("progress", message=f"[{index}/{total}] Dönüştürülüyor", percent=95.0)
        elif status == "error":
            error(f"yt-dlp hook hatası: {d.get('error', 'bilinmiyor')}")

    ydl_opts["progress_hooks"].append(hook)

    try:
        with yt_dlp.YoutubeDL(ydl_opts) as ydl:
            info = ydl.extract_info(url, download=True)
            title = info.get("title", url) if isinstance(info, dict) else url
            track_done(title, index, total)
            return True
    except Exception as exc:  # noqa: BLE001
        error(f"yt-dlp indirme hatası: {exc}")
        traceback.print_exc(file=sys.stderr)
        return False


def main(argv: List[str]) -> int:
    parser = argparse.ArgumentParser(prog="musicle_ytdlp")
    parser.add_argument("--output", required=True, help="Hedef klasör")
    parser.add_argument("urls", nargs="+", help="İndirilecek YouTube URL'leri")
    args = parser.parse_args(argv)

    out = Path(args.output)
    out.mkdir(parents=True, exist_ok=True)

    ffmpeg = _resolve_ffmpeg()
    if ffmpeg is None:
        error("ffmpeg bulunamadı (MUSICLE_FFMPEG veya gömülü ./ffmpeg)")
        return 2

    start(f"yt-dlp: {len(args.urls)} URL işlenecek")
    failures = 0
    for i, url in enumerate(args.urls, start=1):
        if not _run_one(url, str(out), ffmpeg, i, len(args.urls)):
            failures += 1

    if failures:
        emit("done", message=f"{len(args.urls) - failures}/{len(args.urls)} başarılı", percent=100.0)
        return 1
    done()
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
