#!/usr/bin/env bash
set -euo pipefail

SOURCE_ROOT="${IKEMEN_SOURCE_ROOT:-$HOME/apps/ikemen-source}"
DEST_ROOT="${IKEMEN_LANCERO_TEST_ROOT:-$HOME/apps/ikemen-dev}"
CHAR_NAME="${1:-MK1_LIU-KANG}"
CHAR_STAGE="${2:-stages/kfm.def}"
DRY_RUN=0
POSITIONAL=()

usage() {
	cat <<'USAGE'
Usage: build/lancero/add-source-char-to-lancero-dev.sh [--dry-run] [character [stage]]

Adds one character from ~/apps/ikemen-source to ~/apps/ikemen-dev/data/select.def
and writes launchers into the dev install:
  start.sh
  start.cmd

The generated start.sh understands stable and dev modes.

Defaults:
  character = MK1_LIU-KANG
  stage     = stages/kfm.def
USAGE
}

while (($#)); do
	case "$1" in
		--dry-run) DRY_RUN=1 ;;
		-h|--help) usage; exit 0 ;;
		*) POSITIONAL+=("$1") ;;
	esac
	shift
done

CHAR_NAME="${POSITIONAL[0]:-MK1_LIU-KANG}"
CHAR_STAGE="${POSITIONAL[1]:-stages/kfm.def}"

SOURCE_ROOT="$(cd "$SOURCE_ROOT" && pwd -P)"
DEST_ROOT="$(cd "$DEST_ROOT" && pwd -P)"
SELECT_DEF="$DEST_ROOT/data/select.def"
SOURCE_CHAR_DIR="$SOURCE_ROOT/chars/$CHAR_NAME"

log() {
	printf '%s\n' "$*"
}

write_file() {
	local path="$1"
	local content="$2"
	local tmp
	tmp="$(mktemp)"
	cat >"$tmp" <<EOF
$content
EOF
	if [[ -e "$path" ]] && cmp -s "$tmp" "$path"; then
		rm -f "$tmp"
		log "unchanged $(basename "$path")"
		return 0
	fi
	log "write $(basename "$path")"
	if ((DRY_RUN)); then
		rm -f "$tmp"
		return 0
	fi
	mkdir -p "$(dirname "$path")"
	mv "$tmp" "$path"
}

insert_select_entry() {
	local entry="$CHAR_NAME, $CHAR_STAGE"
	if grep -Fq "${CHAR_NAME}," "$SELECT_DEF" || grep -Fxq "$CHAR_NAME" "$SELECT_DEF"; then
		log "select.def already contains $CHAR_NAME"
		return 0
	fi

	local tmp
	tmp="$(mktemp)"
	awk -v entry="$entry" '
		BEGIN { inserted = 0 }
		{
			print
			if (!inserted && $0 == "randomselect") {
				print entry
				inserted = 1
			}
		}
		END {
			if (!inserted) {
				print entry
			}
		}
	' "$SELECT_DEF" >"$tmp"

	if [[ -e "$SELECT_DEF" ]] && cmp -s "$tmp" "$SELECT_DEF"; then
		rm -f "$tmp"
		log "select.def unchanged"
		return 0
	fi

	log "update select.def"
	if ((DRY_RUN)); then
		rm -f "$tmp"
		return 0
	fi
	mv "$tmp" "$SELECT_DEF"
}

write_dev_start_sh() {
	write_file "$DEST_ROOT/start.sh" '#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"

export LD_LIBRARY_PATH="$PWD/lib:${LD_LIBRARY_PATH:-}"
MODE="stable"

case "${1:-}" in
	dev|stable)
		MODE="$1"
		shift
	;;
esac

pick_bin() {
	case "$MODE" in
		dev)
			for candidate in ./Ikemen_GO_Linux_dev ./Ikemen_GO_dev ./Ikemen_GO_Linux ./Ikemen_GO; do
				[[ -x "$candidate" ]] && { printf '%s\n' "$candidate"; return 0; }
			done
		;;
		stable)
			for candidate in ./Ikemen_GO_Linux ./Ikemen_GO; do
				[[ -x "$candidate" ]] && { printf '%s\n' "$candidate"; return 0; }
			done
		;;
	esac
	return 1
}

BIN="$(pick_bin)" || {
	echo "ERROR: no runnable IKEMEN binary found for mode '$MODE'." >&2
	exit 1
}

exec "$BIN" "$@"
'
	if ((! DRY_RUN)); then
		chmod +x "$DEST_ROOT/start.sh"
	fi
}

write_dev_start_cmd() {
	write_file "$DEST_ROOT/start.cmd" '@echo off
setlocal
cd /d "%~dp0"

set "PATH=%~dp0lib;%PATH%"
set "MODE=stable"

if /I "%~1"=="dev" (
	set "MODE=dev"
	shift
) else if /I "%~1"=="stable" (
	shift
)

set "ARGS=%1 %2 %3 %4 %5 %6 %7 %8 %9"

if /I "%MODE%"=="dev" (
	if exist "%~dp0Ikemen_GO_Linux_dev" (
		"%~dp0Ikemen_GO_Linux_dev" %ARGS%
		exit /b %errorlevel%
	)
	if exist "%~dp0Ikemen_GO_dev" (
		"%~dp0Ikemen_GO_dev" %ARGS%
		exit /b %errorlevel%
	)
)

if exist "%~dp0Ikemen_GO_Linux" (
	"%~dp0Ikemen_GO_Linux" %ARGS%
	exit /b %errorlevel%
)

if exist "%~dp0Ikemen_GO" (
	"%~dp0Ikemen_GO" %ARGS%
	exit /b %errorlevel%
)

echo ERROR: no runnable IKEMEN executable found for mode %MODE%.
exit /b 1
'
}

if ((DRY_RUN)); then
	log "Dry run: no files will be changed."
fi

if [[ ! -d "$SOURCE_CHAR_DIR" ]]; then
	echo "ERROR: source character directory not found: $SOURCE_CHAR_DIR" >&2
	exit 1
fi

if [[ ! -f "$SELECT_DEF" ]]; then
	echo "ERROR: select.def not found: $SELECT_DEF" >&2
	exit 1
fi

insert_select_entry
write_dev_start_sh
write_dev_start_cmd
