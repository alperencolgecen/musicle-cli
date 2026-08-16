# Third-Party Notices

`musicle-cli` (MusicLeCLI) bundles a self-contained download engine built from
the following third-party projects. They are extracted at runtime into the
user cache directory and invoked via an embedded Python interpreter.

When the binary is built **without** the `engine_assets` build tag, these
components are not bundled and the application falls back to its legacy
pure-Go download pipeline.

## Bundled components

| Component | Version | License | Source |
|-----------|---------|---------|--------|
| [yt-dlp](https://github.com/yt-dlp/yt-dlp) | 2024.12.13 | The Unlicense | https://github.com/yt-dlp/yt-dlp |
| [spotDL](https://github.com/spotDL/spotify-downloader) (spotdl) | 4.4.0 | MIT | https://github.com/spotDL/spotify-downloader |
| [FFmpeg](https://ffmpeg.org/) (static build) | 6.1.1 | LGPL-2.1+ (with GPL components) | https://ffmpeg.org/ |
| CPython (venv interpreter) | 3.11 | PSF License | https://www.python.org/ |

## License summaries

### yt-dlp — The Unlicense
yt-dlp is released into the public domain. The full text of the Unlicense is
available at <https://unlicense.org/>.

### spotDL — MIT License
Copyright (c) 2021 spotDL Developers.

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

### FFmpeg — LGPL-2.1+
FFmpeg is distributed under the GNU Lesser General Public License version 2.1
or later. When compiled with the optional GPL components, the resulting binary
is covered by the GNU General Public License version 2 or later. The full
license texts are available at <https://www.ffmpeg.org/legal.html>.

### CPython — PSF License Agreement
The embedded Python interpreter is distributed under the Python Software
Foundation License Agreement. The full text is available at
<https://www.python.org/psf/license/>.

## Disclaimer

This software is provided for personal, lawful use only. The authors are not
responsible for any misuse of the bundled download engine or for any copyright
infringement resulting from downloaded content.
