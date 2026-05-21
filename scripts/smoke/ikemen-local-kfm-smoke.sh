#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd -P)"
DEFAULT_BIN="$REPO_ROOT/Ikemen_GO_Linux"

usage() {
  cat <<'EOF'
Usage: scripts/smoke/ikemen-local-kfm-smoke.sh [--fixture-root PATH] [--p1-def PATH]
  [--p2-def PATH] [--stage PATH] [--time N] [--rounds N] [--p1-ai N]
  [--p2-ai N] [--timeout N] [--loadmotif PATH] [--use-repo-main] [--help]

AI is defaulted to level 8 for both players. Export IKEMEN_LOCAL_KFM_KEEP_WORKDIR=1
to keep a temp workdir for log/result inspection when a run fails.

Runs one kfm-style quick-vs match against local fixture data and emits a pass/fail
result quickly for iteration loops. Default timing is best-2-of-3 (2 rounds) at 10 seconds each.
EOF
}

IKEMEN_BIN="${IKEMEN_BIN:-$DEFAULT_BIN}"
IKEMEN_LOCAL_KFM_FIXTURE_ROOT="${IKEMEN_LOCAL_KFM_FIXTURE_ROOT:-/mnt/shared/games/ikemen2}"
IKEMEN_LOCAL_KFM_P1_DEF="${IKEMEN_LOCAL_KFM_P1_DEF:-chars/kfm/kfm.def}"
IKEMEN_LOCAL_KFM_P2_DEF="${IKEMEN_LOCAL_KFM_P2_DEF:-chars/kfm/kfm.def}"
IKEMEN_LOCAL_KFM_STAGE="${IKEMEN_LOCAL_KFM_STAGE:-}"
IKEMEN_LOCAL_KFM_TIME="${IKEMEN_LOCAL_KFM_TIME:-10}"
IKEMEN_LOCAL_KFM_ROUNDS="${IKEMEN_LOCAL_KFM_ROUNDS:-2}"
IKEMEN_LOCAL_KFM_P1_AI="${IKEMEN_LOCAL_KFM_P1_AI:-8}"
IKEMEN_LOCAL_KFM_P2_AI="${IKEMEN_LOCAL_KFM_P2_AI:-8}"
IKEMEN_LOCAL_KFM_TIMEOUT="${IKEMEN_LOCAL_KFM_TIMEOUT:-60}"
IKEMEN_LOCAL_KFM_KEEP_WORKDIR="${IKEMEN_LOCAL_KFM_KEEP_WORKDIR:-0}"
IKEMEN_LOCAL_KFM_USE_REPO_MAIN="${IKEMEN_LOCAL_KFM_USE_REPO_MAIN:-1}"
IKEMEN_LOCAL_KFM_LOADMOTIF="${IKEMEN_LOCAL_KFM_LOADMOTIF:-}"
IKEMEN_LOCAL_KFM_SUPPRESS_UI_ERROR="${IKEMEN_LOCAL_KFM_SUPPRESS_UI_ERROR:-1}"

log_error_patterns='attempt to call a non-function object|runtime error:|panic:|segmentation fault|fatal error|stack traceback'

die() {
  printf 'ERROR: %s\n' "$*" >&2
  exit 1
}

resolve_def_path() {
  local fixture_root="$1"
  local requested="$2"
  local request_base
  local -a candidates

  request_base="$(basename "$requested")"

  if [[ -f "$fixture_root/$requested" ]]; then
    echo "$fixture_root/$requested"
    return 0
  fi

  mapfile -t candidates < <(
    {
      find "$fixture_root/chars" -type f -name "$request_base" 2>/dev/null
      find "$fixture_root/data" -type f -name "$request_base" 2>/dev/null
    } | sort
  )

  if (( "${#candidates[@]}" > 0 )); then
    echo "${candidates[0]}"
    return 0
  fi

  return 1
}

resolve_stage() {
  local fixture_root="$1"
  local candidate="$2"
  if [[ -n "$candidate" ]]; then
    if [[ -f "$fixture_root/$candidate" ]]; then
      echo "$candidate"
      return 0
    fi
  fi

  if [[ -f "$fixture_root/stages/kfm.def" ]]; then
    echo "stages/kfm.def"
    return 0
  fi

  local stage
  stage="$(find "$fixture_root/stages" -maxdepth 1 -type f -name '*.def' | sort | head -n1 || true)"
  if [[ -n "$stage" ]]; then
    stage="${stage#"$fixture_root/"}"
    echo "$stage"
    return 0
  fi

  return 1
}

resolve_motif() {
  local fixture_root="$1"
  local candidate="$2"
  if [[ -n "$candidate" ]]; then
    if [[ -f "$fixture_root/$candidate" ]]; then
      echo "$candidate"
      return 0
    fi
  fi

  local motif
  for motif in \
    data/mugen1/system.def \
    data/lbig/system.def \
    data/kfm/system.def \
    data/big/system.def \
    data/system.def; do
    if [[ -f "$fixture_root/$motif" ]]; then
      echo "$motif"
      return 0
    fi
  done

  motif="$(find "$fixture_root/data" -maxdepth 2 -type f -name system.def | sort | head -n1 || true)"
  if [[ -n "$motif" ]]; then
    motif="${motif#"$fixture_root/"}"
    echo "$motif"
    return 0
  fi

  return 1
}

resolve_fixture_root() {
  local requested="$1"
  local default_fallback="$2"

  if [[ -n "$requested" && -d "$requested" ]]; then
    echo "$requested"
    return 0
  fi

  if [[ -d "$default_fallback" ]]; then
    echo "$default_fallback"
    return 0
  fi

  for candidate in \
    /mnt/shared/games/ikemen2 \
    /mnt/shared/games/ikeman \
    /mnt/shared/Emu/ikemen \
    /mnt/shared/games/ikemen; do
    if [[ -d "$candidate" ]]; then
      if [[ -f "$candidate/chars/kfm/kfm.def" && -f "$candidate/stages/kfm.def" ]]; then
        echo "$candidate"
        return 0
      fi
    fi
  done

  for candidate in \
    /mnt/shared/games/ikemen2 \
    /mnt/shared/games/ikeman \
    /mnt/shared/Emu/ikemen \
    /mnt/shared/games/ikemen; do
    if [[ -d "$candidate" ]]; then
      echo "$candidate"
      return 0
    fi
  done

  return 1
}

if [[ "$#" -gt 0 ]]; then
  while [[ "$#" -gt 0 ]]; do
    case "$1" in
      --fixture-root)
        shift
        [[ "$#" -gt 0 ]] || die "--fixture-root requires a path"
        IKEMEN_LOCAL_KFM_FIXTURE_ROOT="$1"
        ;;
      --p1-def)
        shift
        [[ "$#" -gt 0 ]] || die "--p1-def requires a path"
        IKEMEN_LOCAL_KFM_P1_DEF="$1"
        ;;
      --p2-def)
        shift
        [[ "$#" -gt 0 ]] || die "--p2-def requires a path"
        IKEMEN_LOCAL_KFM_P2_DEF="$1"
        ;;
      --stage)
        shift
        [[ "$#" -gt 0 ]] || die "--stage requires a path"
        IKEMEN_LOCAL_KFM_STAGE="$1"
        ;;
      --time)
        shift
        [[ "$#" -gt 0 ]] || die "--time requires a value"
        IKEMEN_LOCAL_KFM_TIME="$1"
        ;;
      --rounds)
        shift
        [[ "$#" -gt 0 ]] || die "--rounds requires a value"
        IKEMEN_LOCAL_KFM_ROUNDS="$1"
        ;;
      --p1-ai)
        shift
        [[ "$#" -gt 0 ]] || die "--p1-ai requires a value"
        IKEMEN_LOCAL_KFM_P1_AI="$1"
        ;;
      --p2-ai)
        shift
        [[ "$#" -gt 0 ]] || die "--p2-ai requires a value"
        IKEMEN_LOCAL_KFM_P2_AI="$1"
        ;;
      --timeout)
        shift
        [[ "$#" -gt 0 ]] || die "--timeout requires a value"
        IKEMEN_LOCAL_KFM_TIMEOUT="$1"
        ;;
      --loadmotif)
        shift
        [[ "$#" -gt 0 ]] || die "--loadmotif requires a path"
        IKEMEN_LOCAL_KFM_LOADMOTIF="$1"
        ;;
      --use-repo-main)
        IKEMEN_LOCAL_KFM_USE_REPO_MAIN=1
        ;;
      -h|--help)
        usage
        exit 0
        ;;
      *)
        die "Unknown option: $1"
        ;;
    esac
    shift
  done
fi

IKEMEN_BIN="$(readlink -f "$IKEMEN_BIN")"
if [[ ! -x "$IKEMEN_BIN" ]]; then
  die "Ikemen binary not found or not executable: $IKEMEN_BIN"
fi

if [[ ! -d "$IKEMEN_LOCAL_KFM_FIXTURE_ROOT" ]]; then
  IKEMEN_LOCAL_KFM_FIXTURE_ROOT="$(resolve_fixture_root "$IKEMEN_LOCAL_KFM_FIXTURE_ROOT" /mnt/shared/games/ikemen2)"
fi
if [[ -z "${IKEMEN_LOCAL_KFM_FIXTURE_ROOT:-}" || ! -d "$IKEMEN_LOCAL_KFM_FIXTURE_ROOT" ]]; then
  die "Fixture root not found: $IKEMEN_LOCAL_KFM_FIXTURE_ROOT"
fi

if [[ "$IKEMEN_LOCAL_KFM_USE_REPO_MAIN" == "1" ]]; then
  if [[ ! -f "$REPO_ROOT/external/script/main.lua" ]]; then
    die "Missing repo external/script/main.lua: $REPO_ROOT/external/script/main.lua"
  fi
else
  if [[ ! -f "$IKEMEN_LOCAL_KFM_FIXTURE_ROOT/external/script/main.lua" ]]; then
    die "Missing external/script/main.lua in fixture root: $IKEMEN_LOCAL_KFM_FIXTURE_ROOT"
  fi
fi

if [[ ! -f "$IKEMEN_LOCAL_KFM_FIXTURE_ROOT/save/config.json" ]]; then
  printf '[warn] Missing save/config.json in fixture root: %s\n' "$IKEMEN_LOCAL_KFM_FIXTURE_ROOT" >&2
fi

WORK_ROOT="$(mktemp -d "${TMPDIR:-/tmp}/ikemen-local-kfm.XXXXXX")"
cleanup() {
  if [[ "$IKEMEN_LOCAL_KFM_KEEP_WORKDIR" != "1" ]]; then
    rm -rf "$WORK_ROOT"
  else
    printf 'workdir preserved for inspection: %s\n' "$WORK_ROOT"
  fi
}
trap cleanup EXIT

for dir in data chars stages external font sound; do
  if [[ -d "$IKEMEN_LOCAL_KFM_FIXTURE_ROOT/$dir" ]]; then
    cp -a "$IKEMEN_LOCAL_KFM_FIXTURE_ROOT/$dir" "$WORK_ROOT/$dir"
  fi
done

if [[ -d "$IKEMEN_LOCAL_KFM_FIXTURE_ROOT/save" ]]; then
  mkdir -p "$WORK_ROOT/save"
  cp -a "$IKEMEN_LOCAL_KFM_FIXTURE_ROOT/save/." "$WORK_ROOT/save/"
fi

if [[ "$IKEMEN_LOCAL_KFM_USE_REPO_MAIN" == "1" ]]; then
  rm -rf "$WORK_ROOT/external/script"
  cp -a "$REPO_ROOT/external/script" "$WORK_ROOT/external/"
else
  cp -a "$IKEMEN_LOCAL_KFM_FIXTURE_ROOT/external/script" "$WORK_ROOT/external/"
fi

if [[ ! -f "$WORK_ROOT/external/script/main.lua" ]]; then
  die "external/script/main.lua missing in workdir."
fi

if [[ ! -f "$WORK_ROOT/save/config.json" ]]; then
  printf '[warn] save/config.json missing after copy; creating minimal placeholder for smoke run.\n' >&2
  printf '{}\n' > "$WORK_ROOT/save/config.json"
fi

if [[ -d "$WORK_ROOT/data" ]]; then
  if [[ ! -f "$WORK_ROOT/data/gofx/gofx.def" && -f "$WORK_ROOT/data/gofx.def" ]]; then
    mkdir -p "$WORK_ROOT/data/gofx"
    if [[ ! -f "$WORK_ROOT/data/gofx/gofx.def" ]]; then
      cp -a "$WORK_ROOT/data/gofx.def" "$WORK_ROOT/data/gofx/gofx.def"
    fi
    for asset in gofx.sff gofx.air gofx.snd gofx.pcx; do
      if [[ -f "$WORK_ROOT/data/$asset" ]]; then
        cp -a "$WORK_ROOT/data/$asset" "$WORK_ROOT/data/gofx/$asset"
      fi
    done
  fi
fi

P1_DEF="$(resolve_def_path "$WORK_ROOT" "$IKEMEN_LOCAL_KFM_P1_DEF" || true)"
P2_DEF="$(resolve_def_path "$WORK_ROOT" "$IKEMEN_LOCAL_KFM_P2_DEF" || true)"
if [[ -z "$P1_DEF" ]]; then
  die "Cannot resolve P1 def: $IKEMEN_LOCAL_KFM_P1_DEF"
fi
if [[ -z "$P2_DEF" ]]; then
  die "Cannot resolve P2 def: $IKEMEN_LOCAL_KFM_P2_DEF"
fi

P1_DEF="${P1_DEF#"$WORK_ROOT/"}"
P2_DEF="${P2_DEF#"$WORK_ROOT/"}"

RESOLVED_STAGE="$(resolve_stage "$WORK_ROOT" "$IKEMEN_LOCAL_KFM_STAGE" || true)"
if [[ -z "$RESOLVED_STAGE" ]]; then
  die "No stage file found in fixture. Set one explicitly with --stage."
fi

RESOLVED_LOADMOTIF=""
if [[ -n "$IKEMEN_LOCAL_KFM_LOADMOTIF" ]]; then
  RESOLVED_LOADMOTIF="$(resolve_motif "$WORK_ROOT" "$IKEMEN_LOCAL_KFM_LOADMOTIF" || true)"
  if [[ -z "$RESOLVED_LOADMOTIF" ]]; then
    die "No motif system.def found in fixture. Set one with --loadmotif."
  fi
fi

RESULT_FILE="$WORK_ROOT/quickvs-result.json"
RUN_LOG="$WORK_ROOT/quickvs.log"
RUNNER=()
if [[ -z "${DISPLAY:-}" ]] && command -v xvfb-run >/dev/null 2>&1; then
  RUNNER=(xvfb-run -a)
fi

printf 'fixture-root=%s\n' "$IKEMEN_LOCAL_KFM_FIXTURE_ROOT"
printf 'p1=%s\n' "$P1_DEF"
printf 'p2=%s\n' "$P2_DEF"
printf 'stage=%s\n' "$RESOLVED_STAGE"
if [[ -n "$RESOLVED_LOADMOTIF" ]]; then
  printf 'loadmotif=%s\n' "$RESOLVED_LOADMOTIF"
  printf 'motif-override=%s\n' "$RESOLVED_LOADMOTIF"
fi
printf 'rounds=%s\n' "$IKEMEN_LOCAL_KFM_ROUNDS"
printf 'time=%s\n' "$IKEMEN_LOCAL_KFM_TIME"
printf 'ai=%s/%s\n' "$IKEMEN_LOCAL_KFM_P1_AI" "$IKEMEN_LOCAL_KFM_P2_AI"
printf 'timeout=%s\n' "$IKEMEN_LOCAL_KFM_TIMEOUT"

start_ts=$(date +%s)
set +e
LAUNCH_ARGS=(
  -windowed
  -nosound
  -nojoy
  -jsonstdout
  -nomusic
  -debugstartup
  -p1
  "$P1_DEF"
  -p2
  "$P2_DEF"
  -p1.ai
  "$IKEMEN_LOCAL_KFM_P1_AI"
  -p2.ai
  "$IKEMEN_LOCAL_KFM_P2_AI"
  -s
  "$RESOLVED_STAGE"
  -rounds
  "$IKEMEN_LOCAL_KFM_ROUNDS"
  -time
  "$IKEMEN_LOCAL_KFM_TIME"
  -resultfile
  "$RESULT_FILE"
)
if [[ "$IKEMEN_LOCAL_KFM_SUPPRESS_UI_ERROR" == "1" ]]; then
  LAUNCH_ARGS+=(-noerrordialog)
fi
if [[ -n "$IKEMEN_LOCAL_KFM_LOADMOTIF" ]]; then
  LAUNCH_ARGS+=(-r "$RESOLVED_LOADMOTIF")
fi
(
  cd "$WORK_ROOT"
  IKEMEN_SUPPRESS_ERROR_DIALOG="$IKEMEN_LOCAL_KFM_SUPPRESS_UI_ERROR" \
  timeout "${IKEMEN_LOCAL_KFM_TIMEOUT}s" "${RUNNER[@]}" "$IKEMEN_BIN" \
    "${LAUNCH_ARGS[@]}"
) >"$RUN_LOG" 2>&1
run_status=$?
set -e
end_ts=$(date +%s)
elapsed=$((end_ts - start_ts))

printf 'elapsed=%ss\n' "$elapsed"

if [[ "$run_status" -eq 124 ]]; then
  echo "FAIL: match run hit timeout (${IKEMEN_LOCAL_KFM_TIMEOUT}s)"
  tail -n 80 "$RUN_LOG" || true
  exit 1
fi

if [[ "$run_status" -ne 0 ]]; then
  echo "FAIL: match exited with status $run_status"
  tail -n 80 "$RUN_LOG" || true
  exit 1
fi

if grep -Eqi "$log_error_patterns" "$RUN_LOG"; then
  echo "FAIL: match hit startup/runtime error"
  tail -n 80 "$RUN_LOG" || true
  exit 1
fi

if [[ ! -s "$RESULT_FILE" ]]; then
  echo "FAIL: match did not produce result file: $RESULT_FILE"
  tail -n 80 "$RUN_LOG" || true
  exit 1
fi

if ! grep -q '"statsLog"' "$RESULT_FILE"; then
  echo "FAIL: result file is missing expected JSON payload"
  head -n 20 "$RESULT_FILE" || true
  exit 1
fi

echo "OK: local kfm-style match completed."
echo "result: $RESULT_FILE"
echo "log: $RUN_LOG"
echo "workdir: $WORK_ROOT"
