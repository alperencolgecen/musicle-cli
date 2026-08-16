# Third-Party Notices

`musicle-cli` (MusicLeCLI) bundles a self-contained download engine built from
the following third-party projects. They are extracted at runtime into the
user cache directory and invoked directly (no Python).

These components are embedded into the binary via `go:embed` at build time
(see `scripts/prepare-engine.sh`). A build run with `-tags engine_assets`
produces a fully self-contained binary that extracts and runs them at startup
with no further downloads.

## Bundled components

| Component | Version | License | Source |
|-----------|---------|---------|--------|
| [yt-dlp](https://github.com/yt-dlp/yt-dlp) | 2026.07.04 | The Unlicense | https://github.com/yt-dlp/yt-dlp |
| [FFmpeg](https://ffmpeg.org/) (static build) | 7.0.2 | LGPL-2.1+ (with GPL components) | https://ffmpeg.org/ |

The versions above are the latest fetched by `scripts/prepare-engine.sh`; see
`assets/engine/manifest.yaml` for the pinned values.

## License summaries

### yt-dlp — The Unlicense
yt-dlp is released into the public domain. The full text of the Unlicense is
available at <https://unlicense.org/>.

### FFmpeg — LGPL-2.1+
FFmpeg is distributed under the GNU Lesser General Public License version 2.1
or later. When compiled with the optional GPL components, the resulting binary
is covered by the GNU General Public License version 2 or later. The full
license texts are available at <https://www.ffmpeg.org/legal.html>.

## Disclaimer

This software is provided for personal, lawful use only. The authors are not
responsible for any misuse of the bundled download engine or for any copyright
infringement resulting from downloaded content.