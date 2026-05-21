#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd -P)"
DEFAULT_BIN="$REPO_ROOT/Ikemen_GO_Linux"

usage() {
  cat <<'EOF'
Usage: scripts/smoke/ikemen-workflow-matrix.sh [--char-sweep] [--char-sweep-only]
  [--char-sweep-output PATH] [--char-sweep-limit N] [--char-sweep-start N]
  [--char-sweep-time N] [--char-sweep-rounds N] [--char-sweep-timeout N]

Without arguments, runs the standard startup option matrix.
--char-sweep runs a per-character self-play compatibility check and writes JSON.
EOF
}

IKEMEN_BIN="${IKEMEN_BIN:-$DEFAULT_BIN}"
IKEMEN_WORKFLOW_TIMEOUT="${IKEMEN_WORKFLOW_TIMEOUT:-25}"
IKEMEN_WORKFLOW_CHAR_SWEEP="${IKEMEN_WORKFLOW_CHAR_SWEEP:-}"
IKEMEN_WORKFLOW_CHAR_SWEEP_ONLY="${IKEMEN_WORKFLOW_CHAR_SWEEP_ONLY:-}"
IKEMEN_WORKFLOW_CHAR_SWEEP_OUTPUT="${IKEMEN_WORKFLOW_CHAR_SWEEP_OUTPUT:-$REPO_ROOT/ikemen-char-sweep-results.json}"
IKEMEN_WORKFLOW_CHAR_SWEEP_LIMIT="${IKEMEN_WORKFLOW_CHAR_SWEEP_LIMIT:-0}"
IKEMEN_WORKFLOW_CHAR_SWEEP_START="${IKEMEN_WORKFLOW_CHAR_SWEEP_START:-0}"
IKEMEN_WORKFLOW_CHAR_SWEEP_TIME="${IKEMEN_WORKFLOW_CHAR_SWEEP_TIME:-30}"
IKEMEN_WORKFLOW_CHAR_SWEEP_ROUNDS="${IKEMEN_WORKFLOW_CHAR_SWEEP_ROUNDS:-3}"
IKEMEN_WORKFLOW_CHAR_SWEEP_TIMEOUT="${IKEMEN_WORKFLOW_CHAR_SWEEP_TIMEOUT:-75}"
IKEMEN_WORKFLOW_SUPPRESS_UI_ERROR="${IKEMEN_WORKFLOW_SUPPRESS_UI_ERROR:-1}"
IKEMEN_WORKFLOW_KEEP_WORK_ROOT="${IKEMEN_WORKFLOW_KEEP_WORK_ROOT:-0}"
IKEMEN_WORKFLOW_CHAR_SWEEP_LOG_DIR="${IKEMEN_WORKFLOW_CHAR_SWEEP_LOG_DIR:-}"
IKEMEN_WORKFLOW_P1_AI_LEVEL="${IKEMEN_WORKFLOW_P1_AI_LEVEL:-8}"
IKEMEN_WORKFLOW_P2_AI_LEVEL="${IKEMEN_WORKFLOW_P2_AI_LEVEL:-8}"
IKEMEN_WORKFLOW_FIXTURE_ROOT="${IKEMEN_WORKFLOW_FIXTURE_ROOT:-}"
IKEMEN_WORKFLOW_PLAYER_DEF="${IKEMEN_WORKFLOW_PLAYER_DEF:-chars/kfm/kfm.def}"
IKEMEN_WORKFLOW_STAGE_DEF="${IKEMEN_WORKFLOW_STAGE_DEF:-stages/kfm.def}"
IKEMEN_WORKFLOW_CHAR_DIR="${IKEMEN_WORKFLOW_CHAR_DIR:-chars}"
IKEMEN_WORKFLOW_CHARACTER_METADATA_FILE="${IKEMEN_WORKFLOW_CHARACTER_METADATA_FILE:-$REPO_ROOT/docs/features/workflow-smoke-matrix/character-metadata-overrides.tsv}"

while [[ "$#" -gt 0 ]]; do
  case "$1" in
    --char-sweep)
      IKEMEN_WORKFLOW_CHAR_SWEEP="1"
      ;;
    --char-sweep-only)
      IKEMEN_WORKFLOW_CHAR_SWEEP="1"
      IKEMEN_WORKFLOW_CHAR_SWEEP_ONLY="1"
      ;;
    --char-sweep-output)
      shift
      [[ "$#" -gt 0 ]] || { echo "Missing value for --char-sweep-output" >&2; usage; exit 1; }
      IKEMEN_WORKFLOW_CHAR_SWEEP_OUTPUT="$1"
      ;;
    --char-sweep-limit)
      shift
      [[ "$#" -gt 0 ]] || { echo "Missing value for --char-sweep-limit" >&2; usage; exit 1; }
      IKEMEN_WORKFLOW_CHAR_SWEEP_LIMIT="$1"
      ;;
    --char-sweep-start)
      shift
      [[ "$#" -gt 0 ]] || { echo "Missing value for --char-sweep-start" >&2; usage; exit 1; }
      IKEMEN_WORKFLOW_CHAR_SWEEP_START="$1"
      ;;
    --char-sweep-time)
      shift
      [[ "$#" -gt 0 ]] || { echo "Missing value for --char-sweep-time" >&2; usage; exit 1; }
      IKEMEN_WORKFLOW_CHAR_SWEEP_TIME="$1"
      ;;
    --char-sweep-rounds)
      shift
      [[ "$#" -gt 0 ]] || { echo "Missing value for --char-sweep-rounds" >&2; usage; exit 1; }
      IKEMEN_WORKFLOW_CHAR_SWEEP_ROUNDS="$1"
      ;;
    --char-sweep-timeout)
      shift
      [[ "$#" -gt 0 ]] || { echo "Missing value for --char-sweep-timeout" >&2; usage; exit 1; }
      IKEMEN_WORKFLOW_CHAR_SWEEP_TIMEOUT="$1"
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown argument: $1" >&2
      usage
      exit 1
      ;;
  esac
  shift
done

if [[ -z "$IKEMEN_WORKFLOW_FIXTURE_ROOT" && -d /mnt/shared/games/ikemen2 ]]; then
  IKEMEN_WORKFLOW_FIXTURE_ROOT="/mnt/shared/games/ikemen2"
elif [[ -z "$IKEMEN_WORKFLOW_FIXTURE_ROOT" && -d /mnt/shared/Emu/ikemen ]]; then
  IKEMEN_WORKFLOW_FIXTURE_ROOT="/mnt/shared/Emu/ikemen"
fi

if [[ "$IKEMEN_WORKFLOW_CHAR_SWEEP" == "1" && -z "$IKEMEN_WORKFLOW_FIXTURE_ROOT" ]]; then
  echo "ERROR: --char-sweep requires IKEMEN_WORKFLOW_FIXTURE_ROOT to locate character defs." >&2
  exit 1
fi

IKEMEN_BIN="$(readlink -f "$IKEMEN_BIN")"
if [[ ! -x "$IKEMEN_BIN" ]]; then
  echo "ERROR: Ikemen binary not found or not executable: $IKEMEN_BIN" >&2
  echo "Set IKEMEN_BIN to a built Linux binary, e.g. ./Ikemen_GO_Linux" >&2
  exit 1
fi

WORK_ROOT="$(mktemp -d "${TMPDIR:-/tmp}/ikemen-workflow.XXXXXX")"
cleanup() {
  if [[ "$IKEMEN_WORKFLOW_KEEP_WORK_ROOT" == "1" ]]; then
    echo "INFO: preserved workflow temp tree at $WORK_ROOT" >&2
    return
  fi
  rm -rf "$WORK_ROOT"
}
trap cleanup EXIT

bootstrap_from_repo() {
  for dir in data external font lib save; do
    if [[ "$dir" == "save" ]]; then
      mkdir -p "$WORK_ROOT/save"
      cp -a "$REPO_ROOT/save/." "$WORK_ROOT/save/" 2>/dev/null || true
    else
      ln -s "$REPO_ROOT/$dir" "$WORK_ROOT/$dir"
    fi
  done
}

bootstrap_from_fixture() {
  local fixture_root="$1"
  cp -a "$fixture_root/." "$WORK_ROOT/"
  cp -a "$REPO_ROOT/external/script/." "$WORK_ROOT/external/script/"
}

bootstrap_from_fixture_char_sweep() {
  local fixture_root="$1"
  local dir
  local -a link_dirs=(chars stages sound video font lib)

  mkdir -p "$WORK_ROOT"

  # Character sweeps only need a subset of the runtime tree.
  # Link large content directories instead of copying full fixture roots.
  for dir in "${link_dirs[@]}"; do
    if [[ -d "$fixture_root/$dir" ]]; then
      ln -s "$fixture_root/$dir" "$WORK_ROOT/$dir"
    elif [[ -d "$REPO_ROOT/$dir" ]]; then
      ln -s "$REPO_ROOT/$dir" "$WORK_ROOT/$dir"
    fi
  done

  # Keep data writable in the temp root so missing core scripts can be backfilled
  # without mutating the fixture source tree.
  if [[ -d "$fixture_root/data" ]]; then
    cp -a "$fixture_root/data" "$WORK_ROOT/data"
  else
    ln -s "$REPO_ROOT/data" "$WORK_ROOT/data"
  fi

  # Use repo external tree so script/runtime changes under test are active.
  ln -s "$REPO_ROOT/external" "$WORK_ROOT/external"

  mkdir -p "$WORK_ROOT/save"
  cp -a "$REPO_ROOT/save/." "$WORK_ROOT/save/" 2>/dev/null || true
}

ensure_required_data_scripts() {
  local required=(
    action.zss
    demo.zss
    dizzy.zss
    functions.zss
    guardbreak.zss
    score.zss
    system.zss
    tag.zss
    training.zss
  )
  local name
  mkdir -p "$WORK_ROOT/data"
  for name in "${required[@]}"; do
    if [[ ! -f "$WORK_ROOT/data/$name" && -f "$REPO_ROOT/data/$name" ]]; then
      cp -a "$REPO_ROOT/data/$name" "$WORK_ROOT/data/$name"
    fi
  done
}

if [[ -n "$IKEMEN_WORKFLOW_FIXTURE_ROOT" && -d "$IKEMEN_WORKFLOW_FIXTURE_ROOT" ]]; then
  if [[ "$IKEMEN_WORKFLOW_CHAR_SWEEP_ONLY" == "1" ]]; then
    bootstrap_from_fixture_char_sweep "$IKEMEN_WORKFLOW_FIXTURE_ROOT"
  else
    bootstrap_from_fixture "$IKEMEN_WORKFLOW_FIXTURE_ROOT"
  fi
else
  bootstrap_from_repo
fi

ensure_required_data_scripts

write_workflow_config() {
  local motif_path=""
  for candidate in data/ikemen1/system.def data/lbig/system.def data/mugen1/system.def data/system.def; do
    if [[ -f "$WORK_ROOT/$candidate" ]]; then
      motif_path="$candidate"
      break
    fi
  done

  if [[ -z "$motif_path" ]]; then
    echo "ERROR: Could not find a usable motif inside the workflow runtime tree." >&2
    exit 1
  fi

  mkdir -p "$WORK_ROOT/save"
  cat > "$WORK_ROOT/save/config.ini" <<EOF
[Config]
Motif = $motif_path
System = external/script/main.lua

[Debug]
Font = font/debug.def
FontScale = 0.5
StartStage = stages/stage1.def
EOF
}

write_workflow_config

ensure_font_alias() {
  local alias_path="$WORK_ROOT/font/arial.ttf"
  if [[ ! -e "$alias_path" ]]; then
    mkdir -p "$WORK_ROOT/font"
    if [[ -f "$WORK_ROOT/font/Open_Sans/OpenSans-Regular.ttf" ]]; then
      ln -s "Open_Sans/OpenSans-Regular.ttf" "$alias_path"
    elif [[ -f "$WORK_ROOT/font/mssansserif.ttf" ]]; then
      ln -s "mssansserif.ttf" "$alias_path"
    fi
  fi
}

ensure_font_alias

RUNNER=()
if [[ -n "${IKEMEN_RUNNER:-}" ]]; then
  read -r -a RUNNER <<< "$IKEMEN_RUNNER"
elif [[ -z "${DISPLAY:-}" ]] && command -v xvfb-run >/dev/null 2>&1; then
  RUNNER=(xvfb-run -a)
fi

COMMON_ARGS=()
if [[ "$IKEMEN_WORKFLOW_SUPPRESS_UI_ERROR" == "1" ]]; then
  COMMON_ARGS+=(-noerrordialog)
fi
log_error_patterns='attempt to call a non-function object|runtime error:|panic:|segmentation fault|fatal error|stack traceback'
failed=0
char_sweep_failed=0
char_results=()
quick_player_one=""
quick_player_two=""
quick_stage=""

run_binary_case() {
  local timeout_seconds="$1"
  local log_file="$2"
  shift 2
  local run_status=0
  local start_ts
  local end_ts
  local elapsed_seconds

  start_ts=$(date +%s)
  (
    cd "$WORK_ROOT"
    if [[ "$IKEMEN_WORKFLOW_SUPPRESS_UI_ERROR" == "1" ]]; then
      IKEMEN_SUPPRESS_ERROR_DIALOG=1 timeout "${timeout_seconds}s" "${RUNNER[@]}" "$IKEMEN_BIN" "${COMMON_ARGS[@]}" "$@"
    else
      timeout "${timeout_seconds}s" "${RUNNER[@]}" "$IKEMEN_BIN" "${COMMON_ARGS[@]}" "$@"
    fi
  ) >"$log_file" 2>&1
  run_status=$?
  end_ts=$(date +%s)
  elapsed_seconds=$((end_ts - start_ts))
  printf '%s|%s|%s\n' "$run_status" "$elapsed_seconds" "$log_file"
}

escape_json() {
  local value="$1"
  value="${value//\\/\\\\}"
  value="${value//\"/\\\"}"
  value="${value//$'\n'/\\n}"
  value="${value//$'\r'/\\r}"
  value="${value//$'\t'/\\t}"
  printf '%s' "$value"
}

normalize_character_id() {
  local value="$1"
  value="$(printf '%s' "$value" | tr '[:upper:]' '[:lower:]' | sed -E 's/[^a-z0-9]+/-/g; s/^-+//; s/-+$//')"
  printf '%s' "$value"
}

slugify_tag() {
  local value="$1"
  value="$(printf '%s' "$value" | tr '[:upper:]' '[:lower:]' | sed -E 's/[^a-z0-9]+/-/g; s/^-+//; s/-+$//')"
  printf '%s' "$value"
}

def_info_value() {
  local def_path="$1"
  local wanted_key="$2"
  awk -v wanted="$(printf '%s' "$wanted_key" | tr '[:upper:]' '[:lower:]')" '
    function trim(s) {
      gsub(/^[ \t]+|[ \t]+$/, "", s)
      return s
    }
    {
      line = $0
      sub(/;.*/, "", line)
      if (line ~ /^[[:space:]]*\[/) {
        sec = tolower(line)
        in_info = (sec ~ /^[[:space:]]*\[info\][[:space:]]*$/)
        next
      }
      if (!in_info) next
      pos = index(line, "=")
      if (pos == 0) next
      key = trim(substr(line, 1, pos - 1))
      val = trim(substr(line, pos + 1))
      gsub(/\r/, "", val)
      if (tolower(key) != wanted) next
      if (val ~ /^".*"$/) {
        sub(/^"/, "", val)
        sub(/"$/, "", val)
      }
      print val
      exit
    }
  ' "$def_path"
}

metadata_lookup() {
  local character_id="$1"
  local -n out_display_name="$2"
  local -n out_series="$3"
  local -n out_genre="$4"
  local -n out_tags_csv="$5"
  local -n out_notes="$6"
  out_display_name=""
  out_series="unknown"
  out_genre="fighting"
  out_tags_csv=""
  out_notes=""
  if [[ ! -f "$IKEMEN_WORKFLOW_CHARACTER_METADATA_FILE" ]]; then
    return
  fi
  while IFS='|' read -r id display_name series genre tags_csv notes; do
    [[ -n "$id" ]] || continue
    [[ "$id" =~ ^[[:space:]]*# ]] && continue
    id="$(printf '%s' "$id" | tr -d '\r' | sed -E 's/^[[:space:]]+|[[:space:]]+$//g')"
    if [[ "$id" != "$character_id" ]]; then
      continue
    fi
    out_display_name="$(printf '%s' "${display_name:-}" | tr -d '\r' | sed -E 's/^[[:space:]]+|[[:space:]]+$//g')"
    out_series="$(printf '%s' "${series:-unknown}" | tr -d '\r' | sed -E 's/^[[:space:]]+|[[:space:]]+$//g')"
    out_genre="$(printf '%s' "${genre:-fighting}" | tr -d '\r' | sed -E 's/^[[:space:]]+|[[:space:]]+$//g')"
    out_tags_csv="$(printf '%s' "${tags_csv:-}" | tr -d '\r' | sed -E 's/^[[:space:]]+|[[:space:]]+$//g')"
    out_notes="$(printf '%s' "${notes:-}" | tr -d '\r' | sed -E 's/^[[:space:]]+|[[:space:]]+$//g')"
    return
  done < "$IKEMEN_WORKFLOW_CHARACTER_METADATA_FILE"
}

build_tags_json() {
  local series="$1"
  local genre="$2"
  local tags_csv="$3"
  local -a tags=()
  local -A seen=()
  local item
  local slug

  slug="$(slugify_tag "$genre")"
  if [[ -n "$slug" && -z "${seen[$slug]:-}" ]]; then
    tags+=("$slug")
    seen[$slug]=1
  fi

  if [[ "$series" != "unknown" ]]; then
    slug="$(slugify_tag "$series")"
    if [[ -n "$slug" && -z "${seen[$slug]:-}" ]]; then
      tags+=("$slug")
      seen[$slug]=1
    fi
  fi

  IFS=',' read -r -a raw_tags <<< "$tags_csv"
  for item in "${raw_tags[@]}"; do
    item="$(printf '%s' "$item" | sed -E 's/^[[:space:]]+|[[:space:]]+$//g')"
    [[ -n "$item" ]] || continue
    slug="$(slugify_tag "$item")"
    [[ -n "$slug" ]] || continue
    if [[ -z "${seen[$slug]:-}" ]]; then
      tags+=("$slug")
      seen[$slug]=1
    fi
  done

  local out="["
  local i
  for i in "${!tags[@]}"; do
    if (( i > 0 )); then
      out+=","
    fi
    out+="\"$(escape_json "${tags[$i]}")\""
  done
  out+="]"
  printf '%s' "$out"
}

run_case() {
  local name="$1"
  shift
  local log_file="$WORK_ROOT/${name}.log"
  printf '==> [%s]\n' "$name"
  local run_meta
  local run_status
  local elapsed_seconds
  run_meta="$(run_binary_case "$IKEMEN_WORKFLOW_TIMEOUT" "$log_file" "$@")"
  run_status="$(echo "$run_meta" | cut -d'|' -f1)"
  elapsed_seconds="$(echo "$run_meta" | cut -d'|' -f2)"

  printf '    elapsed=%ss\n' "$elapsed_seconds"

  if [[ "$run_status" -ne 0 && "$run_status" -ne 124 ]]; then
    echo "FAIL: $name exited with status $run_status" >&2
    tail -n 20 "$log_file" >&2 || true
    failed=1
    return
  fi

  if [[ "$run_status" -eq 124 ]]; then
    echo "FAIL: $name hit timeout (${IKEMEN_WORKFLOW_TIMEOUT}s)" >&2
    failed=1
    return
  fi

  if grep -Eqi "$log_error_patterns" "$log_file"; then
    echo "FAIL: $name hit a startup/runtime error" >&2
    tail -n 20 "$log_file" >&2 || true
    failed=1
    return
  fi

  if [[ -s "$log_file" ]]; then
    echo "OK: $name"
    tail -n 5 "$log_file" || true
  else
    echo "OK: $name (no output)"
  fi
}

run_result_case() {
  local name="$1"
  local result_file="$2"
  shift 2
  rm -f "$result_file"
  run_case "$name" "$@"
  if [[ "$failed" -ne 0 ]]; then
    return
  fi
  if [[ ! -s "$result_file" ]]; then
    echo "FAIL: $name did not write a result file: $result_file" >&2
    failed=1
    return
  fi
  if ! grep -q '"statsLog"' "$result_file"; then
    echo "FAIL: $name wrote a result file without the expected JSON payload: $result_file" >&2
    failed=1
    return
  fi
  echo "OK: $name result file -> $result_file"
}

normalize_result_path() {
  local file_path="$1"
  if [[ "$file_path" == "$WORK_ROOT"/* ]]; then
    file_path="${file_path#"$WORK_ROOT/"}"
  fi
  printf '%s\n' "$file_path"
}

is_excluded_def() {
  local def_path="$1"
  local def_name
  def_name="$(basename "$def_path" | tr '[:upper:]' '[:lower:]')"

  if [[ "$def_name" == *ending*.def || "$def_name" == *select*.def || "$def_name" == *fight*.def || "$def_name" == *gofx*.def || "$def_name" == *intro*.def ]]; then
    return 0
  fi
  return 1
}

pick_character_def() {
  local char_dir="$1"
  local char_name="$2"
  local preferred="$char_dir/$char_name.def"
  local -a candidates
  local candidate
  local char_name_normalized

  char_name_normalized="$(printf '%s' "$char_name" | tr '[:upper:]' '[:lower:]' | tr -d '[:space:]' | tr -d '_' | tr -d '-' | tr -d '.')"

  if [[ -f "$preferred" ]]; then
    echo "$preferred"
    return 0
  fi

  mapfile -t candidates < <(
    find "$char_dir" -maxdepth 1 -type f -name '*.def' \
      | sort
  )

  if (( ${#candidates[@]} == 0 )); then
    return 1
  fi

  for candidate in "${candidates[@]}"; do
    if is_excluded_def "$candidate"; then
      continue
    fi

    local base_name
    local base_normalized
    base_name="$(basename "$candidate" .def)"
    base_normalized="$(printf '%s' "${base_name,,}" | tr -d '[:space:]' | tr -d '_' | tr -d '-' | tr -d '.')"

    if [[ "${base_normalized}" == *ai*patch* ]]; then
      continue
    fi

    if [[ "$base_name" == "$char_name" || "$base_normalized" == "$char_name_normalized" ]]; then
      echo "$candidate"
      return 0
    fi
  done

  for candidate in "${candidates[@]}"; do
    if is_excluded_def "$candidate"; then
      continue
    fi

    local candidate_basename
    candidate_basename="$(basename "$candidate")"
    if [[ "${candidate_basename,,}" == *ending*.def || "${candidate_basename,,}" == *select*.def || "${candidate_basename,,}" == *fight*.def || "${candidate_basename,,}" == *gofx*.def || "${candidate_basename,,}" == *intro*.def ]]; then
      continue
    fi

    if [[ "$candidate_basename" == *".def" ]]; then
      echo "$candidate"
      return 0
    fi
  done

  return 1
}

run_char_sweep_case() {
  local char_name="$1"
  local def_path="$2"
  local result_file="$3"
  shift 3

  local log_file="$WORK_ROOT/char-sweep-${char_name}.log"
  local run_status
  local elapsed_seconds
  local run_meta
  local status="ok"
  local compatible="true"
  local error_hint=""
  local entry
  local has_result_file="false"
  local log_signature=""
  local character_id
  local character_display_name
  local character_internal_name
  local def_display_name
  local def_internal_name
  local character_series
  local character_genre
  local character_tags_csv
  local character_notes
  local character_tags_json
  local character_has_long_intro="false"
  local folder_normalized

  character_id="$(normalize_character_id "$char_name")"
  folder_normalized="$character_id"
  def_internal_name="$(def_info_value "$def_path" "name" || true)"
  def_display_name="$(def_info_value "$def_path" "displayname" || true)"
  def_internal_name="$(printf '%s' "$def_internal_name" | tr -d '\r' | sed -E 's/^[[:space:]]+|[[:space:]]+$//g; s/^"//; s/"$//')"
  def_display_name="$(printf '%s' "$def_display_name" | tr -d '\r' | sed -E 's/^[[:space:]]+|[[:space:]]+$//g; s/^"//; s/"$//')"
  character_internal_name="$def_internal_name"
  character_display_name="$def_display_name"
  if [[ -z "$character_display_name" ]]; then
    character_display_name="$character_internal_name"
  fi
  if [[ -z "$character_display_name" ]]; then
    character_display_name="$char_name"
  fi
  metadata_lookup "$character_id" character_display_name character_series character_genre character_tags_csv character_notes
  if [[ -z "$character_display_name" ]]; then
    character_display_name="$def_display_name"
  fi
  if [[ -z "$character_display_name" ]]; then
    character_display_name="$def_internal_name"
  fi
  if [[ -z "$character_display_name" ]]; then
    character_display_name="$char_name"
  fi
  character_tags_json="$(build_tags_json "$character_series" "$character_genre" "$character_tags_csv")"
  if printf '%s' "$character_tags_csv" | tr '[:upper:]' '[:lower:]' | grep -q 'long-intro'; then
    character_has_long_intro="true"
  fi

  run_meta="$(run_binary_case "$IKEMEN_WORKFLOW_CHAR_SWEEP_TIMEOUT" "$log_file" "$@")"
  run_status="$(echo "$run_meta" | cut -d'|' -f1)"
  elapsed_seconds="$(echo "$run_meta" | cut -d'|' -f2)"

  if [[ "$run_status" -eq 124 ]]; then
    status="timeout"
    compatible="false"
    error_hint="match run timed out after ${IKEMEN_WORKFLOW_CHAR_SWEEP_TIMEOUT}s"
  elif [[ "$run_status" -ne 0 ]]; then
    status="error"
    compatible="false"
    error_hint="exit code $run_status"
  fi

  if [[ "$compatible" == "false" ]]; then
    log_signature="$(grep -Eim1 'open data/demo.zss: no such file|attempt to call a non-function object|runtime error:|Panic:|fatal error|segmentation fault|invalid animationType|stack traceback' "$log_file" | head -n 1 || true)"
    if [[ -n "$log_signature" ]]; then
      error_hint="${error_hint}: ${log_signature}"
    fi
  fi

  if [[ "$compatible" == "true" ]] && grep -Eqi "$log_error_patterns" "$log_file"; then
    status="error"
    compatible="false"
    error_hint="runtime/crash text in logs"
  fi

  if [[ -s "$result_file" ]]; then
    has_result_file="true"
    if ! grep -q '"statsLog"' "$result_file"; then
      status="error"
      compatible="false"
      error_hint="result file missing expected statsLog payload"
    fi
  fi

  if [[ "$has_result_file" == "false" && "$compatible" == "true" ]]; then
    status="missing-result"
    compatible="false"
    error_hint="result file was not created"
  fi

  if [[ "$compatible" == "true" ]]; then
    entry="$(printf '    {"character":"%s","characterId":"%s","folderNormalized":"%s","displayName":"%s","internalName":"%s","series":"%s","genre":"%s","tags":%s,"notes":"%s","hasLongIntro":%s,"path":"%s","status":"%s","compatible":%s,"elapsedSeconds":%s,"exitCode":%s,"resultExists":%s,"errorHint":null}' \
      "$(escape_json "$char_name")" \
      "$(escape_json "$character_id")" \
      "$(escape_json "$folder_normalized")" \
      "$(escape_json "$character_display_name")" \
      "$(escape_json "$character_internal_name")" \
      "$(escape_json "$character_series")" \
      "$(escape_json "$character_genre")" \
      "$character_tags_json" \
      "$(escape_json "$character_notes")" \
      "$character_has_long_intro" \
      "$(escape_json "$(normalize_result_path "$def_path")")" \
      "$status" \
      "$compatible" \
      "$elapsed_seconds" \
      "$run_status" \
      "$has_result_file")"
  else
    entry="$(printf '    {"character":"%s","characterId":"%s","folderNormalized":"%s","displayName":"%s","internalName":"%s","series":"%s","genre":"%s","tags":%s,"notes":"%s","hasLongIntro":%s,"path":"%s","status":"%s","compatible":%s,"elapsedSeconds":%s,"exitCode":%s,"resultExists":%s,"errorHint":"%s"}' \
      "$(escape_json "$char_name")" \
      "$(escape_json "$character_id")" \
      "$(escape_json "$folder_normalized")" \
      "$(escape_json "$character_display_name")" \
      "$(escape_json "$character_internal_name")" \
      "$(escape_json "$character_series")" \
      "$(escape_json "$character_genre")" \
      "$character_tags_json" \
      "$(escape_json "$character_notes")" \
      "$character_has_long_intro" \
      "$(escape_json "$(normalize_result_path "$def_path")")" \
      "$status" \
      "$compatible" \
      "$elapsed_seconds" \
      "$run_status" \
      "$has_result_file" \
      "$(escape_json "$error_hint")")"
  fi

  char_results+=("$entry")
  if [[ "$compatible" == "true" ]]; then
    echo "OK: [${char_name}] compatible (exit=$run_status, elapsed=${elapsed_seconds}s)"
  else
    echo "FAIL: [${char_name}] incompatible (${status}) - ${error_hint}" >&2
    char_sweep_failed=1
  fi

  if [[ -n "$IKEMEN_WORKFLOW_CHAR_SWEEP_LOG_DIR" ]]; then
    cp -f "$log_file" "$IKEMEN_WORKFLOW_CHAR_SWEEP_LOG_DIR/${character_id}.log"
  fi
}

run_character_sweep() {
  local char_root
  if [[ "$IKEMEN_WORKFLOW_CHAR_DIR" = /* ]]; then
    char_root="$IKEMEN_WORKFLOW_CHAR_DIR"
  else
    char_root="$WORK_ROOT/$IKEMEN_WORKFLOW_CHAR_DIR"
  fi
  local output_file="$IKEMEN_WORKFLOW_CHAR_SWEEP_OUTPUT"
  local sweep_results_dir="$WORK_ROOT/char-sweep-results"
  local limit="$IKEMEN_WORKFLOW_CHAR_SWEEP_LIMIT"
  local start_index="$IKEMEN_WORKFLOW_CHAR_SWEEP_START"
  local chars_tested=0
  local chars_total=0
  local chars_skipped=0
  local chars_offset_skipped=0
  local chars_considered=0
  local char_dir
  local char_name
  local def_path
  local result_file
  local result_entry_count=0
  local compatible_count=0
  local incompatible_count=0
  local def_basename
  local char_sweep_log_dir="$IKEMEN_WORKFLOW_CHAR_SWEEP_LOG_DIR"

  if [[ -n "$char_sweep_log_dir" ]]; then
    mkdir -p "$char_sweep_log_dir"
  fi

  if [[ ! -d "$char_root" ]]; then
    echo "ERROR: character sweep root not found: $char_root" >&2
    char_sweep_failed=1
    return
  fi

  mkdir -p "$sweep_results_dir"
  mkdir -p "$(dirname "$output_file")"

  for char_dir in "$char_root"/*; do
    [[ -d "$char_dir" ]] || continue
    chars_total=$((chars_total + 1))
    chars_considered=$((chars_considered + 1))

    if [[ "$chars_considered" -le "$start_index" ]]; then
      chars_offset_skipped=$((chars_offset_skipped + 1))
      continue
    fi

    if [[ "$limit" -gt 0 && "$chars_tested" -ge "$limit" ]]; then
      chars_skipped=$((chars_skipped + 1))
      continue
    fi

    char_name="$(basename "$char_dir")"
    def_path="$(pick_character_def "$char_dir" "$char_name" || true)"
    if [[ -z "$def_path" ]]; then
      echo "WARN: skip ${char_name}: no usable .def file found" >&2
      chars_skipped=$((chars_skipped + 1))
      continue
    fi
    def_basename="$(basename "$def_path")"
    if is_excluded_def "$def_basename"; then
      echo "WARN: skip ${char_name}: excluded def selected ($def_basename)" >&2
      chars_skipped=$((chars_skipped + 1))
      continue
    fi

    result_file="$sweep_results_dir/${char_name}.json"
    run_char_sweep_case \
      "$char_name" \
      "$def_path" \
      "$result_file" \
      -windowed -nosound -nojoy -jsonstdout \
      -p1 "$def_path" \
      -p2 "$def_path" \
      -p1.ai "$IKEMEN_WORKFLOW_P1_AI_LEVEL" \
      -p2.ai "$IKEMEN_WORKFLOW_P2_AI_LEVEL" \
      -rounds "$IKEMEN_WORKFLOW_CHAR_SWEEP_ROUNDS" \
      -time "$IKEMEN_WORKFLOW_CHAR_SWEEP_TIME" \
      -debugstartup \
      -resultfile "$result_file"

    chars_tested=$((chars_tested + 1))
    result_entry_count=$((result_entry_count + 1))
  done

  for entry in "${char_results[@]}"; do
    if [[ "$entry" == *'"compatible":true'* ]]; then
      compatible_count=$((compatible_count + 1))
    else
      incompatible_count=$((incompatible_count + 1))
    fi
  done

  {
    printf '{\n'
    printf '  "schema": "ikemen-char-sweep/v2",\n'
    printf '  "generatedAt": "%s",\n' "$(date -u +'%Y-%m-%dT%H:%M:%SZ')"
    printf '  "characterMetadataFile": "%s",\n' "$(escape_json "$(normalize_result_path "$IKEMEN_WORKFLOW_CHARACTER_METADATA_FILE")")"
    printf '  "workRoot": "%s",\n' "$(escape_json "$(normalize_result_path "$WORK_ROOT")")"
    printf '  "characterDir": "%s",\n' "$(escape_json "$(normalize_result_path "$char_root")")"
    printf '  "timeSeconds": %s,\n' "$IKEMEN_WORKFLOW_CHAR_SWEEP_TIME"
    printf '  "roundsToWin": %s,\n' "$IKEMEN_WORKFLOW_CHAR_SWEEP_ROUNDS"
    printf '  "timeoutSeconds": %s,\n' "$IKEMEN_WORKFLOW_CHAR_SWEEP_TIMEOUT"
    printf '  "startIndex": %s,\n' "$start_index"
    printf '  "offsetSkippedCharacters": %s,\n' "$chars_offset_skipped"
    printf '  "totalCharacters": %s,\n' "$chars_total"
    printf '  "testedCharacters": %s,\n' "$chars_tested"
    printf '  "skippedCharacters": %s,\n' "$chars_skipped"
    printf '  "compatibleCount": %s,\n' "$compatible_count"
    printf '  "incompatibleCount": %s,\n' "$incompatible_count"
    printf '  "results": [\n'
    if (( result_entry_count > 0 )); then
      for i in "${!char_results[@]}"; do
        if (( i + 1 < result_entry_count )); then
          printf '    %s,\n' "${char_results[$i]}"
        else
          printf '    %s\n' "${char_results[$i]}"
        fi
      done
    fi
    printf '  ]\n'
    printf '}\n'
  } > "$output_file"

  echo "Character compatibility report: $output_file"
  if [[ "$char_sweep_failed" -ne 0 ]]; then
    failed=1
  fi
}

cases=(
  "startup||"
  "windowed-nosound||-windowed -nosound"
  "windowed-jsonstdout||-windowed -nosound -jsonstdout"
  "ailevel||-ailevel 8"
  "speedtest||-speedtest 100"
  "framerate||-framerate 120"
  "debug-toggles||-debug -togglelifebars -maxpowermode"
  "window-size||-windowed -width 1280 -height 720"
  "volume||-setvolume 0"
)

if [[ -n "$IKEMEN_WORKFLOW_FIXTURE_ROOT" && -d "$IKEMEN_WORKFLOW_FIXTURE_ROOT" ]]; then
  if [[ -f "$IKEMEN_WORKFLOW_FIXTURE_ROOT/$IKEMEN_WORKFLOW_PLAYER_DEF" ]]; then
    quick_player_one="$IKEMEN_WORKFLOW_PLAYER_DEF"
    quick_player_two="$IKEMEN_WORKFLOW_PLAYER_DEF"
    if [[ -f "$IKEMEN_WORKFLOW_FIXTURE_ROOT/$IKEMEN_WORKFLOW_STAGE_DEF" ]]; then
      quick_stage="$IKEMEN_WORKFLOW_STAGE_DEF"
    else
      quick_stage=""
    fi
  else
    quick_player_one=""
    quick_stage=""
  fi
fi

if [[ -n "$quick_player_one" ]]; then
  cases+=("quickvs||-windowed -nosound -nojoy -jsonstdout -p1 $quick_player_one -p2 $quick_player_two -p1.ai $IKEMEN_WORKFLOW_P1_AI_LEVEL -p2.ai $IKEMEN_WORKFLOW_P2_AI_LEVEL -rounds 1 -time 5 -debugstartup")
  if [[ -n "$quick_stage" ]]; then
    cases+=("quickvs-stage-override||-windowed -nosound -nojoy -jsonstdout -p1 $quick_player_one -p2 $quick_player_two -p1.ai $IKEMEN_WORKFLOW_P1_AI_LEVEL -p2.ai $IKEMEN_WORKFLOW_P2_AI_LEVEL -s $quick_stage -rounds 1 -time 5 -debugstartup -width 1280 -height 720")
  fi
else
  echo "NOTE: quick-vs workflow skipped; set IKEMEN_WORKFLOW_FIXTURE_ROOT to a tree containing the selected character def (IKEMEN_WORKFLOW_PLAYER_DEF)." >&2
fi

quickvs_ai_args=(-p1.ai "$IKEMEN_WORKFLOW_P1_AI_LEVEL" -p2.ai "$IKEMEN_WORKFLOW_P2_AI_LEVEL")

if [[ "$IKEMEN_WORKFLOW_CHAR_SWEEP_ONLY" != "1" ]]; then
  for case_spec in "${cases[@]}"; do
    name="${case_spec%%||*}"
    args="${case_spec#*||}"
    # shellcheck disable=SC2206
    argv=($args)
    run_case "$name" "${argv[@]}"
  done

  if [[ -n "$quick_player_one" ]]; then
    run_result_case "quickvs-resultfile" "$WORK_ROOT/quickvs-result.json" \
      -windowed -nosound -nojoy -jsonstdout \
      -p1 "$quick_player_one" -p2 "$quick_player_two" "${quickvs_ai_args[@]}" \
      -rounds 1 -time 5 -debugstartup \
      -resultfile "$WORK_ROOT/quickvs-result.json"
    run_result_case "quickvs-ai-palette" "$WORK_ROOT/quickvs-ai-palette.json" \
      -windowed -nosound -nojoy -jsonstdout \
      -p1 "$quick_player_one" -p2 "$quick_player_two" \
      -p1.ai 8 -p2.ai 7 -p1.color 1 -p2.color 2 \
      -rounds 1 -time 5 -debugstartup \
      -resultfile "$WORK_ROOT/quickvs-ai-palette.json"

    if [[ -n "$quick_stage" ]]; then
      run_result_case "quickvs-stage-override-resultfile" "$WORK_ROOT/quickvs-stage-override-result.json" \
        -windowed -nosound -nojoy -jsonstdout \
        -p1 "$quick_player_one" -p2 "$quick_player_two" "${quickvs_ai_args[@]}" \
        -s "$quick_stage" -rounds 1 -time 5 -debugstartup \
        -resultfile "$WORK_ROOT/quickvs-stage-override-result.json"
    fi
  fi
fi

if [[ "$IKEMEN_WORKFLOW_CHAR_SWEEP" == "1" ]]; then
  echo "==> character sweep"
  run_character_sweep
fi

if [[ "$failed" -ne 0 ]]; then
  echo "Startup/workflow smoke failed. Logs are in $WORK_ROOT" >&2
  exit 1
fi

echo "All workflow smoke cases passed."
