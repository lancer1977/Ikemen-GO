#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
HELPER="$SCRIPT_DIR/ikemen-control-snapshots.py"
XVFB_SCREEN="${IKEMEN_XVFB_SCREEN:-1280x720x24}"

usage() {
  cat <<'EOF'
Usage: scripts/visual/ikemen-xvfb-control-snapshots.sh [ikemen-control-snapshots.py args...]

Runs the IKEMEN control/snapshot helper inside a private Xvfb display with a
minimal metacity window manager. This avoids failures from the desktop X server
being unavailable or out of client slots while still producing real window PNGs.

Environment:
  IKEMEN_XVFB_SCREEN   Xvfb screen spec, default 1280x720x24

Example:
  scripts/visual/ikemen-xvfb-control-snapshots.sh \
    --workdir /home/lancer1977/code/ikemen-app \
    --bin ./Ikemen_GO_Linux \
    --probe-mode edge \
    --output-dir artifacts/visual-probes/xvfb-pass \
    --step wait:4 \
    --step snap:boot \
    --step key:Return \
    --step wait:2 \
    --step snap:after-enter
EOF
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
fi

for tool in xvfb-run metacity; do
  if ! command -v "$tool" >/dev/null 2>&1; then
    echo "ERROR: required tool not found: $tool" >&2
    exit 1
  fi
done

if [[ ! -x "$HELPER" ]]; then
  echo "ERROR: helper is not executable: $HELPER" >&2
  exit 1
fi

xvfb-run -a -s "-screen 0 $XVFB_SCREEN" bash -lc '
  set -euo pipefail
  metacity --sm-disable --replace >/tmp/ikemen-metacity.log 2>&1 &
  wm_pid=$!
  cleanup() {
    kill "$wm_pid" >/dev/null 2>&1 || true
  }
  trap cleanup EXIT
  sleep 1
  exec "$@"
' bash "$HELPER" "$@"
