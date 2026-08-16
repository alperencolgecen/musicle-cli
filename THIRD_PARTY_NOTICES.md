# Third-Party Notices

`musicle-cli` (MusicLeCLI) is licensed under the **Apache License 2.0**
(see [`LICENSE`](LICENSE)). This file documents every third-party project it
depends on or bundles, together with their licenses, as required by the
respective license terms and Apache-2.0 §4(d).

## Bundled runtime engine

A self-contained download engine is embedded into the binary via `go:embed`
at build time (see `scripts/prepare-engine.sh`). It is extracted at runtime
into the user cache directory and invoked directly (no Python).

| Component | Version | License | Source |
|-----------|---------|---------|--------|
| [yt-dlp](https://github.com/yt-dlp/yt-dlp) | 2026.07.04 | The Unlicense | https://github.com/yt-dlp/yt-dlp |
| [FFmpeg](https://ffmpeg.org/) (static build) | 7.0.2 | LGPL-2.1+ (with GPL components) | https://ffmpeg.org/ |

The versions above are the latest fetched by `scripts/prepare-engine.sh`; see
`assets/engine/manifest.yaml` for the pinned values.

## Direct Go dependencies

| Module | Version | License | Source |
|--------|---------|---------|--------|
| [github.com/charmbracelet/bubbletea](https://pkg.go.dev/github.com/charmbracelet/bubbletea) | v1.3.10 | MIT | https://github.com/charmbracelet/bubbletea |
| [github.com/charmbracelet/bubbles](https://pkg.go.dev/github.com/charmbracelet/bubbles) | v1.0.0 | MIT | https://github.com/charmbracelet/bubbles |
| [github.com/charmbracelet/lipgloss](https://pkg.go.dev/github.com/charmbracelet/lipgloss) | v1.1.0 | MIT | https://github.com/charmbracelet/lipgloss |
| [github.com/dhowden/tag](https://pkg.go.dev/github.com/dhowden/tag) | v0.0.0-20240417053706 | BSD-2-Clause | https://github.com/dhowden/tag |
| [github.com/gopxl/beep](https://pkg.go.dev/github.com/gopxl/beep) | v1.4.1 | MIT | https://github.com/gopxl/beep |
| [github.com/ncruces/zenity](https://pkg.go.dev/github.com/ncruces/zenity) | v0.10.14 | MIT | https://github.com/ncruces/zenity |

## Indirect Go dependencies

| Module | License |
|--------|---------|
| github.com/akavel/rsrc | MIT |
| github.com/atotto/clipboard | BSD-3-Clause |
| github.com/aymanbagabas/go-osc52/v2 | MIT |
| github.com/charmbracelet/colorprofile | MIT |
| github.com/charmbracelet/x/ansi | MIT |
| github.com/charmbracelet/x/cellbuf | MIT |
| github.com/charmbracelet/x/term | MIT |
| github.com/clipperhouse/displaywidth | MIT |
| github.com/clipperhouse/stringish | MIT |
| github.com/clipperhouse/uax29/v2 | MIT |
| github.com/dchest/jsmin | MIT |
| github.com/ebitengine/oto/v3 | Apache-2.0 |
| github.com/ebitengine/purego | Apache-2.0 |
| github.com/erikgeiser/coninput | MIT |
| github.com/hajimehoshi/go-mp3 | Apache-2.0 |
| github.com/icza/bitio | Apache-2.0 |
| github.com/josephspurrier/goversioninfo | MIT |
| github.com/lucasb-eyer/go-colorful | MIT |
| github.com/mattn/go-isatty | MIT |
| github.com/mattn/go-localereader | MIT |
| github.com/mattn/go-runewidth | MIT |
| github.com/mewkiz/flac | The Unlicense |
| github.com/mewkiz/pkg | The Unlicense |
| github.com/muesli/ansi | MIT |
| github.com/muesli/cancelreader | MIT |
| github.com/muesli/termenv | MIT |
| github.com/pkg/errors | BSD-2-Clause |
| github.com/randall77/makefat | The Unlicense |
| github.com/rivo/uniseg | MIT |
| github.com/xo/terminfo | MIT |
| golang.org/x/image | BSD-3-Clause |
| golang.org/x/sys | BSD-3-Clause |
| golang.org/x/text | BSD-3-Clause |

## License summaries

- **MIT** — Permission is hereby granted, free of charge, to any person
  obtaining a copy of this software and associated documentation files, to
  deal in the Software without restriction. License text:
  <https://opensource.org/license/mit>.
- **BSD-2-Clause** — Redistribution and use in source and binary forms, with
  or without modification, are permitted provided that conditions are met.
  License text: <https://opensource.org/license/bsd-2-clause>.
- **BSD-3-Clause** — As BSD-2-Clause, with the additional clause that neither
  the name of the copyright holder nor contributors may be used to endorse or
  promote products derived from this software.
  License text: <https://opensource.org/license/bsd-3-clause>.
- **Apache-2.0** — Apache License, Version 2.0, January 2004.
  License text: <https://www.apache.org/licenses/LICENSE-2.0>.
- **The Unlicense** — A public-domain dedication; the software is released
  into the public domain. License text: <https://unlicense.org/>.
- **LGPL-2.1+ (FFmpeg)** — GNU Lesser General Public License version 2.1 or
  later; when compiled with optional GPL components, the resulting binary is
  covered by the GPL version 2 or later. More at:
  <https://www.ffmpeg.org/legal.html>.

The full license texts of the Go dependencies are available in the Go module
cache (`$GOPATH/pkg/mod/...`) and at each project's upstream repository.

## Disclaimer

This software is provided for personal, lawful use only. The authors are not
responsible for any misuse of the bundled download engine or for any copyright
infringement resulting from downloaded content.