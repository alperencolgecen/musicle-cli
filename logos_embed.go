package main

import "embed"

// logoFS embeds the brand logos so they always resolve regardless of the
// process working directory. The Spotify/YouTube connector cards and the
// connect view render these directly from memory.
//
//go:embed assets/Spotify_logo.png assets/Youtube_logo.png
var logoFS embed.FS
