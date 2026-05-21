#!/usr/bin/env bash
set -euo pipefail

SOURCE_ROOT="${IKEMEN_SOURCE_ROOT:-$HOME/apps/ikemen-source}"
DEST_ROOT="${IKEMEN_LANCERO_TEST_ROOT:-$HOME/apps/ikemen-dev}"
DRY_RUN=0

usage() {
	cat <<'USAGE'
Usage: build/lancero/sync-content-to-lancero-dev.sh [--dry-run]

Copies these content folders from:
  ~/apps/ikemen-source
to:
  ~/apps/ikemen-dev

Synced folders:
  chars/
  stages/
  sound/
  music/

Only the selected folders are updated. Other content in the destination is left alone.
USAGE
}

while (($#)); do
	case "$1" in
		--dry-run) DRY_RUN=1 ;;
		-h|--help) usage; exit 0 ;;
		*) echo "Unknown argument: $1" >&2; usage >&2; exit 2 ;;
	esac
	shift
done

SOURCE_ROOT="$(cd "$SOURCE_ROOT" && pwd -P)"
DEST_ROOT="$(cd "$DEST_ROOT" && pwd -P)"

log() {
	printf '%s\n' "$*"
}

sync_tree() {
	local rel="$1"
	local src="$SOURCE_ROOT/$rel"
	local dst="$DEST_ROOT/$rel"

	if [[ ! -d "$src" ]]; then
		log "skip missing $rel"
		return 0
	fi

	mkdir -p "$dst"
	log "sync $rel/"

	if ((DRY_RUN)); then
		rsync -a --delete --dry-run "$src"/ "$dst"/
		return 0
	fi

	rsync -a --delete "$src"/ "$dst"/
}

if ((DRY_RUN)); then
	log "Dry run: no files will be changed."
fi

sync_tree chars
sync_tree stages
sync_tree sound
sync_tree music
