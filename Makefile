# Makefile — reproducible build, vet and test for musicle-cli
#
# The full binary needs CGO + ALSA dev headers (the oto audio backend). On a
# machine without those, `make build` fails only on the player package; the
# engine and download packages still build/vet/test via `make check`.

GO ?= go
VERSION ?= 1.0.0

LDFLAGS := -ldflags="-s -w -X main.version=$(VERSION)"

.PHONY: all engine build vet test check clean

all: build

# Build the embedded download engine (venv + spotdl + yt-dlp + ffmpeg).
engine:
	./scripts/prepare-engine.sh linux

# Full binary build (requires CGO + ALSA dev headers for the audio player).
# Depends on `engine` so the embedded venv+ffmpeg are prepared, then baked in
# with -tags engine_assets. The result is fully self-contained for end users.
build: engine
	CGO_ENABLED=1 $(GO) build -tags engine_assets $(LDFLAGS) -o build/musicle-cli .

# Vet the packages that do not depend on the audio player (buildable anywhere).
vet:
	$(GO) vet ./internal/... ./bridge/download/...

# Unit tests (engine package + download helpers).
test:
	$(GO) test ./internal/... ./bridge/download/...

# Fast local check: everything that compiles without the cgo audio backend.
check: vet test

clean:
	rm -rf build internal/engine/engine_venv internal/engine/engine_ffmpeg
