#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_DIR="$(mktemp -d "${TMPDIR:-/tmp}/bus-demo.XXXXXX")"
trap 'rm -rf "$TMP_DIR"' EXIT
BIN="$TMP_DIR/bus"
GO_CACHE="${GOCACHE:-$TMP_DIR/go-build-cache}"
CAST_DIR="${BUS_DEMO_CAST_DIR:-$TMP_DIR}"
AGG="${AGG:-agg}"
if ! command -v "$AGG" >/dev/null 2>&1 && [[ -x "$HOME/.cargo/bin/agg" ]]; then
    AGG="$HOME/.cargo/bin/agg"
fi

for cmd in go expect asciinema "$AGG" magick; do
    if ! command -v "$cmd" >/dev/null 2>&1; then
        echo "error: '$cmd' is required to record the demos" >&2
        exit 1
    fi
done

mkdir -p "$ROOT_DIR/demo" "$CAST_DIR"
echo "Building bus..."
(
    cd "$ROOT_DIR"
    GOCACHE="$GO_CACHE" go build -o "$BIN" .
)

if [[ $# -eq 0 ]]; then
    set -- overview sections
fi

for name in "$@"; do
    CONFIG="$TMP_DIR/$name.yaml"
    COLS=82
    ROWS=14
    EXPECT_SCRIPT="$ROOT_DIR/scripts/arrivals.expect"
    ARGS="arrivals"
    case "$name" in
        overview)
            ROWS=18
            EXPECT_SCRIPT="$ROOT_DIR/scripts/overview.expect"
            cat >"$CONFIG" <<'YAML'
arrivals:
    - route: "135"
      stopPoint: Liverpool Street Station
      letter: F
      tflID: 490000138F
YAML
            ;;
        sections)
            cat >"$CONFIG" <<'YAML'
sections: [Home, Work]
arrivals:
    - route: "8"
      stopPoint: Liverpool Street Station
      letter: F
      tflID: 490000138F
      section: Home
    - route: "135"
      stopPoint: Liverpool Street Station
      letter: F
      tflID: 490000138F
      section: Home
    - route: "135"
      stopPoint: Westferry Station
      letter: WX
      tflID: 490002083A
      section: Work
    - route: D7
      stopPoint: Westferry Station
      letter: WX
      tflID: 490002083A
      section: Work
YAML
            ;;
        *) echo "error: unknown demo: $name" >&2; exit 1 ;;
    esac

    export BUS_DEMO_BIN="$BIN" BUS_CONFIG_PATH="$CONFIG"
    export BUS_DEMO_COLS="$COLS" BUS_DEMO_ROWS="$ROWS" BUS_DEMO_ARGS="$ARGS"
    export BUS_DEMO_EXPECT="$EXPECT_SCRIPT"
    CAST="$CAST_DIR/$name.cast"
    RENDERED="$TMP_DIR/$name.gif"
    OUT="$ROOT_DIR/demo/$name.png"
    if [[ "$name" == overview ]]; then
        OUT="$ROOT_DIR/demo/bus-demo.gif"
        RENDERED="$OUT"
    fi

    echo "Recording $name..."
    asciinema rec --overwrite --headless --return \
        --output-format asciicast-v2 \
        --window-size "${COLS}x${ROWS}" -i 1.4 \
        -c 'expect "$BUS_DEMO_EXPECT"' "$CAST"

    "$AGG" --quiet --renderer fontdue --font-family Menlo --theme github-dark \
        --cols "$COLS" --rows "$ROWS" --font-size 20 --line-height 1.0 \
        --fps-cap 30 --idle-time-limit 1.4 --last-frame-duration 3 --speed 0.85 \
        "$CAST" "$RENDERED"
    if [[ "$name" != overview ]]; then
        magick "$RENDERED" -coalesce -delete 0--2 "$OUT"
    fi
    echo "Wrote $OUT"
done
