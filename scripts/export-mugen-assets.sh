#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
ASSET_TOOLS_ROOT="${ASSET_TOOLS_ROOT:-$REPO_ROOT/../Api.Ikemen/scripts/asset-normalization}"
MANIFEST_SCRIPT="${MANIFEST_SCRIPT:-$ASSET_TOOLS_ROOT/generate-mugen-character-manifests.py}"

SOURCE_ROOT="${SOURCE_ROOT:-/mnt/syn1/games/Mugen}"
DEST_ROOT="${DEST_ROOT:-/mnt/syn1/games/Ikemen}"
ASSET_ROOT="$DEST_ROOT/test-assets"
CHAR_DEST="$ASSET_ROOT/chars"
STAGE_DEST="$ASSET_ROOT/stages"

if [[ ! -d "$SOURCE_ROOT" ]]; then
  echo "ERROR: source root not found: $SOURCE_ROOT" >&2
  exit 1
fi

mkdir -p "$CHAR_DEST" "$STAGE_DEST"

declare -A seen_stage_file=()
declare -i copied_chars=0
declare -i copied_stage_files=0
declare -i copied_stage_sets=0

ensure_root_link() {
  local link_path="$1"
  local target="$2"
  if [[ -e "$link_path" && ! -L "$link_path" ]]; then
    return 0
  fi
  if [[ ! -e "$link_path" ]]; then
    ln -s "test-assets/$target" "$link_path"
  fi
}

copy_character_dir() {
  local src_dir="$1"
  local name
  name="$(basename "$src_dir")"
  local dst_dir="$CHAR_DEST/$name"
  mkdir -p "$dst_dir"
  rsync -a --ignore-existing "$src_dir"/ "$dst_dir"/
}

copy_stage_file() {
  local src_file="$1"
  local base
  base="$(basename "$src_file")"
  if [[ -n "${seen_stage_file[$base]:-}" ]]; then
    return 0
  fi
  seen_stage_file["$base"]=1
  if [[ -e "$STAGE_DEST/$base" ]]; then
    return 0
  fi
  rsync -a --ignore-existing "$src_file" "$STAGE_DEST/"
  copied_stage_files+=1
}

export_chars_from_tree() {
  local tree="$1"
  [[ -d "$tree" ]] || return 0
  while IFS= read -r -d '' src_dir; do
    if find "$src_dir" -type f -name '*.def' -print -quit | grep -q .; then
      copy_character_dir "$src_dir"
      copied_chars+=1
    fi
  done < <(find "$tree" -mindepth 1 -maxdepth 1 -type d -print0 | sort -z)
}

export_chars_from_split_runs() {
  local base="$SOURCE_ROOT/Characters-split/runs"
  [[ -d "$base" ]] || return 0
  while IFS= read -r -d '' split_root; do
    while IFS= read -r -d '' src_dir; do
      if find "$src_dir" -type f -name '*.def' -print -quit | grep -q .; then
        copy_character_dir "$src_dir"
        copied_chars+=1
      fi
    done < <(find "$split_root/split/chars" -mindepth 1 -maxdepth 1 -type d -print0 2>/dev/null | sort -z)
  done < <(find "$base" -mindepth 1 -maxdepth 1 -type d -print0 | sort -z)
}

export_stages_from_tree() {
  local tree="$1"
  [[ -d "$tree" ]] || return 0
  while IFS= read -r -d '' stage_dir; do
    while IFS= read -r -d '' src_file; do
      copy_stage_file "$src_file"
    done < <(find "$stage_dir" -maxdepth 1 -type f \( -name '*.def' -o -name '*.sff' -o -name '*.air' \) -print0 | sort -z)
    copied_stage_sets+=1
  done < <(find "$tree" -type d -name stages -print0 | sort -z)
}

ensure_root_link "$DEST_ROOT/chars" "chars"
ensure_root_link "$DEST_ROOT/stages" "stages"

export_chars_from_tree "$SOURCE_ROOT/Characters"
export_chars_from_split_runs
export_stages_from_tree "$SOURCE_ROOT"

if [[ "${SKIP_CHARACTER_MANIFESTS:-0}" != "1" && -x "$MANIFEST_SCRIPT" ]]; then
  "$MANIFEST_SCRIPT" --root "$CHAR_DEST"
fi

cat <<EOF
Export complete.
  Characters exported: $copied_chars
  Stage directories visited: $copied_stage_sets
  Stage files copied: $copied_stage_files
  Asset root: $ASSET_ROOT
EOF
