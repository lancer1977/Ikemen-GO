#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd -P)"

case "$OSTYPE" in
  linux*)
    exec "$SCRIPT_DIR/package-installers.sh" --linux "$@"
    ;;
  msys|cygwin)
    exec "$SCRIPT_DIR/package-installers.sh" --windows "$@"
    ;;
  *)
    echo "ERROR: Unsupported host '$OSTYPE'. Use deploy-local-linux.sh or deploy-local-windows.sh on a matching host." >&2
    exit 1
    ;;
esac
