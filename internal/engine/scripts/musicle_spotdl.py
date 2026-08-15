"""spotdl wrapper used by musicle-cli's embedded download engine.

Invocation:

    python -m musicle_spotdl --output DIR URL [URL ...]

spotdl supports URLs that are tracks, albums, or playlists. We don't pre-
classify them — spotdl does that internally. Output goes to --output; we
also override the audio format and ffmpeg location to use the embedded
binaries.

The wrapper streams JSON progress events on stdout for the Go bridge.
"""

from __future__ import annotations

import argparse
import os
import sys
import traceback
from pathlib import Path
from typing import List

sys.path.insert(0, str(Path(__file__).resolve().parent))

from _musicle_helpers import (  # noqa: E402
    done,
    emit,
    error,
    start,
)


def _resolve_ffmpeg() -> str | None:
    env = os.environ.get("MUSICLE_FFMPEG")
    if env and Path(env).is_file():
        return env
    here = Path(__file__).resolve().parent
    for cand in ("ffmpeg", "ffmpeg.exe"):
        p = here / cand
        if p.is_file():
            return str(p)
    return None


def main(argv: List[str]) -> int:
    parser = argparse.ArgumentParser(prog="musicle_spotdl")
    parser.add_argument("--output", required=True, help="Hedef klasör")
    parser.add_argument("--query", action="store_true",
                        help="URL'leri arama sorgusu olarak yorumla")
    parser.add_argument("urls", nargs="+", help="Spotify URL'leri veya arama sorguları")
    args = parser.parse_args(argv)

    out = Path(args.output)
    out.mkdir(parents=True, exist_ok=True)

    ffmpeg = _resolve_ffmpeg()
    if ffmpeg is None:
        error("ffmpeg bulunamadı")
        return 2

    try:
        import spotdl  # type: ignore  # noqa: F401
        from spotdl.download import downloader as spotdl_downloader  # type: ignore
        from spotdl.utils.config import get_config  # type: ignore
    except ImportError as exc:
        error(f"spotdl yüklenemedi: {exc}")
        return 2

    start(f"spotdl: {len(args.urls)} girdi işlenecek")

    try:
        config = get_config()
        config["output"] = str(out) + os.sep
        config["format"] = "mp3"
        config["bitrate"] = "192k"
        config["ffmpeg"] = ffmpeg
        config["simple_tui"] = False
        config["log_level"] = "ERROR"

        # spotdl 4.x exposes a Downloader class.
        downloader_cls = getattr(spotdl_downloader, "Downloader", None)
        if downloader_cls is None:
            error("spotdl.Downloader bulunamadı (uyumsuz sürüm?)")
            return 2

        downloader = downloader_cls(audio_provider=None, lyrics_provider=None, ffmpeg=ffmpeg)
        # Re-apply config on the downloader instance too.
        if hasattr(downloader, "config"):
            downloader.config.update(config)

        queries = args.urls
        total = len(queries)
        if total == 0:
            done()
            return 0

        # Progress callback shim — spotdl 4.x's Downloader doesn't expose a
        # direct hook, so we emit coarse progress before/after each URL.
        successes = 0
        for i, query in enumerate(queries, start=1):
            emit("progress",
                 message=f"[{i}/{total}] spotdl: {query}",
                 percent=(i - 1) * 100.0 / total)
            try:
                songs = downloader.search(query) if args.query else downloader.search([query])
                if not songs:
                    error(f"Sonuç bulunamadı: {query}")
                    continue
                results = downloader.download_songs(songs)
                ok = sum(1 for r in results if getattr(r, "status", "") == "success")
                successes += ok
            except Exception as exc:  # noqa: BLE001
                error(f"spotdl hata: {exc}")
                traceback.print_exc(file=sys.stderr)

        emit("progress",
             message=f"{successes}/{total} başarılı",
             percent=100.0)
        done()
        return 0 if successes == total else 1
    except Exception as exc:  # noqa: BLE001
        error(f"spotdl çalıştırılamadı: {exc}")
        traceback.print_exc(file=sys.stderr)
        return 1


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
