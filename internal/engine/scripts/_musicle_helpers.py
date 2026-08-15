"""Shared helpers for the musicle-cli embedded download engine.

Both the yt-dlp and spotdl wrappers import this module to emit structured
progress events as JSON lines on stdout. The Go side (internal/engine/python.go)
parses these via bufio.Scanner and forwards them to the bubbletea UI.

Emitted events are documented in internal/engine/python.go (ProgressEvent).
"""

from __future__ import annotations

import json
import sys
import time
from typing import Optional


def emit(event: str, message: str = "", percent: float = -1.0,
         track: str = "", index: int = 0, total: int = 0) -> None:
    """Write a single JSON progress event to stdout and flush immediately."""
    payload = {
        "event": event,
        "message": message,
        "percent": percent,
        "track": track,
        "index": index,
        "total": total,
    }
    sys.stdout.write(json.dumps(payload, ensure_ascii=False) + "\n")
    sys.stdout.flush()


def progress_hook(percent: float, message: str) -> None:
    """Convenience wrapper used by the Go-side callbacks."""
    emit("progress", message=message, percent=percent)


def done(message: str = "Tamamlandı") -> None:
    emit("done", message=message, percent=100.0)


def error(message: str) -> None:
    emit("error", message=message, percent=-1.0)
    sys.stderr.write(message + "\n")
    sys.stderr.flush()


def track_done(track: str, index: int, total: int) -> None:
    pct = 100.0 * index / total if total > 0 else -1.0
    emit("track", message=track, percent=pct, track=track, index=index, total=total)


def start(message: str) -> None:
    emit("start", message=message, percent=0.0)


def throttle(seconds: float = 0.05):
    """Return a decorator that ensures rapid emit() calls don't flood stdout."""
    last = [0.0]

    def wrap(fn):
        def inner(*args, **kwargs):
            now = time.time()
            if now - last[0] < seconds:
                return None
            last[0] = now
            return fn(*args, **kwargs)

        return inner

    return wrap
