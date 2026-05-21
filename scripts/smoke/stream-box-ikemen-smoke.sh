#!/usr/bin/env bash
set -euo pipefail

STREAM_BOX_HOST="${STREAM_BOX_HOST:-stream-box}"
REMOTE_ROOT="${IKEMEN_STREAM_BOX_ROOT:-C:\\Apps\\mugen}"
REMOTE_APPS_ROOT="${IKEMEN_STREAM_BOX_APPS_ROOT:-C:\\Apps}"
TASK_NAME="${IKEMEN_STREAM_BOX_TASK_NAME:-IkemenWorkflowSmoke}"
LAUNCH_PS1="${IKEMEN_STREAM_BOX_LAUNCH_PS1:-$REMOTE_APPS_ROOT\\ikemen-workflow-smoke.ps1}"
STDOUT_LOG="${IKEMEN_STREAM_BOX_STDOUT_LOG:-$REMOTE_APPS_ROOT\\ikemen-workflow-smoke.stdout.log}"
STDERR_LOG="${IKEMEN_STREAM_BOX_STDERR_LOG:-$REMOTE_APPS_ROOT\\ikemen-workflow-smoke.stderr.log}"
PID_LOG="${IKEMEN_STREAM_BOX_PID_LOG:-$REMOTE_APPS_ROOT\\ikemen-workflow-smoke.pid.txt}"
EXIT_LOG="${IKEMEN_STREAM_BOX_EXIT_LOG:-$REMOTE_APPS_ROOT\\ikemen-workflow-smoke.exit.txt}"
WAIT_SECONDS="${IKEMEN_STREAM_BOX_WAIT_SECONDS:-8}"
RESULT_WAIT_SECONDS="${IKEMEN_STREAM_BOX_RESULT_WAIT_SECONDS:-12}"
STREAM_BOX_AI_LEVEL="${IKEMEN_STREAM_BOX_AI_LEVEL:-8}"
STREAM_BOX_TIME="${IKEMEN_STREAM_BOX_TIME:-5}"
STREAM_BOX_ROUNDS="${IKEMEN_STREAM_BOX_ROUNDS:-1}"
STREAM_BOX_MUTE_UI="${IKEMEN_STREAM_BOX_MUTE_UI:-1}"
STREAM_BOX_P1_DEF="${IKEMEN_STREAM_BOX_P1_DEF:-}"
STREAM_BOX_P2_DEF="${IKEMEN_STREAM_BOX_P2_DEF:-}"
RUN_ID="${IKEMEN_STREAM_BOX_RUN_ID:-$(date -u +%Y%m%d-%H%M%S)-$$}"
RESULT_DIR="${IKEMEN_STREAM_BOX_RESULT_DIR:-$REMOTE_APPS_ROOT\\ikemen-results}"
RESULT_FILE="${IKEMEN_STREAM_BOX_RESULT_FILE:-$RESULT_DIR\\ikemen-workflow-smoke.$RUN_ID.json}"

usage() {
  cat <<'EOF'
Usage: scripts/smoke/stream-box-ikemen-smoke.sh [--check] [--sync-main-lua PATH] [--sync-script SRC RELPATH] [--launch] [--post-check]

Validates the stream-box Ikemen test deploy at C:\Apps\mugen.

Options:
  --check              Verify the deploy root and startup guard markers.
  --sync-main-lua PATH Copy a tested main.lua into the stream-box deploy.
  --sync-script SRC RELPATH
                       Copy a local file into the deploy tree at RELPATH.
  --launch             Launch a visible quick-vs smoke through Task Scheduler.
  --post-check         Check process status and fail on new Ikemen.log crash text.
  --time NUM           Override quick-vs fight time in seconds.
  --rounds NUM         Override rounds to win.
  --p1-def PATH        Override P1 char def path.
  --p2-def PATH        Override P2 char def path.
  --ai-level NUM       Override AI level used by both players (default: 8).
  --mute-ui            Force UI crash dialogs off and rely on logs for errors.
  -h, --help           Show this help.

Environment:
  STREAM_BOX_HOST                 SSH alias or host. Default: stream-box
  IKEMEN_STREAM_BOX_ROOT          Remote Ikemen root. Default: C:\Apps\mugen
  IKEMEN_STREAM_BOX_APPS_ROOT     Remote apps root. Default: C:\Apps
  IKEMEN_STREAM_BOX_WAIT_SECONDS  Seconds to wait after launch. Default: 8
  IKEMEN_STREAM_BOX_RESULT_WAIT_SECONDS  Seconds to wait for the result file. Default: 12
  IKEMEN_STREAM_BOX_TIME          Quick-vs round time in seconds. Default: 5
  IKEMEN_STREAM_BOX_ROUNDS        Rounds to win. Default: 1
  IKEMEN_STREAM_BOX_AI_LEVEL      AI level for P1/P2 quick-vs. Default: 8
  IKEMEN_STREAM_BOX_MUTE_UI       Set to 0 to allow crash popups. Default: 1
  IKEMEN_STREAM_BOX_P1_DEF        Optional P1 def. Auto-discovered when blank.
  IKEMEN_STREAM_BOX_P2_DEF        Optional P2 def. Auto-discovered when blank.
EOF
}

die() {
  printf 'ERROR: %s\n' "$*" >&2
  exit 1
}

run_remote() {
  ssh "$STREAM_BOX_HOST" "$@"
}

run_remote_ps() {
  local script="$1"
  local encoded
  encoded="$(printf '%s' "$script" | iconv -f UTF-8 -t UTF-16LE | base64 -w0)"
  run_remote "powershell -NoProfile -EncodedCommand $encoded"
}

remote_file() {
  local relative="$1"
  printf '%s\\%s' "$REMOTE_ROOT" "$relative"
}

check_required_file() {
  local path="$1"
  run_remote "if not exist \"$path\" exit /b 1" >/dev/null || die "Missing stream-box file: $path"
}

check_marker() {
  local marker="$1"
  local path
  path="$(remote_file "external\\script\\main.lua")"
  if ! run_remote "findstr /n /c:\"$marker\" \"$path\"" >/dev/null; then
    die "Missing startup guard marker in deployed main.lua: $marker"
  fi
}

check_any_motif() {
  local candidate
  for candidate in \
    "$(remote_file "data\\ikemen1\\system.def")" \
    "$(remote_file "data\\lbig\\system.def")" \
    "$(remote_file "data\\mugen1\\system.def")" \
    "$(remote_file "data\\system.def")"; do
    if run_remote "if exist \"$candidate\" exit /b 0 else exit /b 1" >/dev/null; then
      return 0
    fi
  done
  die "Missing stream-box motif: expected one of data\\ikemen1, data\\lbig, data\\mugen1, or data\\system.def"
}

check_deploy() {
  check_required_file "$(remote_file "Ikemen_GO.exe")"
  check_required_file "$(remote_file "external\\script\\main.lua")"
  check_any_motif

  check_marker "safeGameOption"
  check_marker "safeCommandLineValue"
  check_marker "safeLoadMotif"
  check_marker "safeAnimGetPreloadedCharData"
  check_marker "safeAnimNew"
}

resolve_default_launch_def() {
  local resolved
  local script
  script="$(cat <<'PS'
$root = 'C:\Apps\mugen\chars'
$def = Get-ChildItem -Path $root -Recurse -Filter '*.def' -File |
  Where-Object {
    $n = $_.Name.ToLower()
    -not ($n -match 'ending|select|fight|gofx|intro')
  } |
  Select-Object -First 1 -ExpandProperty FullName
if ([string]::IsNullOrWhiteSpace($def)) {
  exit 2
}
$prefix = 'C:\Apps\mugen\'
if ($def.StartsWith($prefix, [System.StringComparison]::OrdinalIgnoreCase)) {
  $def = $def.Substring($prefix.Length)
}
Write-Output ($def -replace '/', '\')
PS
)"
  resolved="$(run_remote_ps "$script" 2>/dev/null || true)"
  resolved="$(printf '%s' "$resolved" | tr -d '\r' | awk 'NF {print; exit}')"
  [[ -n "$resolved" ]] || die "Unable to auto-discover a usable stream-box character .def"
  printf '%s\n' "$resolved"
}

ensure_launch_defs() {
  if [[ -z "$STREAM_BOX_P1_DEF" && -z "$STREAM_BOX_P2_DEF" ]]; then
    STREAM_BOX_P1_DEF="$(resolve_default_launch_def)"
    STREAM_BOX_P2_DEF="$STREAM_BOX_P1_DEF"
  fi
  if [[ -z "$STREAM_BOX_P1_DEF" && -n "$STREAM_BOX_P2_DEF" ]]; then
    STREAM_BOX_P1_DEF="$STREAM_BOX_P2_DEF"
  fi
  if [[ -z "$STREAM_BOX_P2_DEF" && -n "$STREAM_BOX_P1_DEF" ]]; then
    STREAM_BOX_P2_DEF="$STREAM_BOX_P1_DEF"
  fi
  if [[ -z "$STREAM_BOX_P1_DEF" || -z "$STREAM_BOX_P2_DEF" ]]; then
    die "launch requires --p1-def/--p2-def (or IKEMEN_STREAM_BOX_P1_DEF/IKEMEN_STREAM_BOX_P2_DEF)"
  fi

  check_required_file "$(remote_file "$STREAM_BOX_P1_DEF")"
  if [[ "$STREAM_BOX_P2_DEF" != "$STREAM_BOX_P1_DEF" ]]; then
    check_required_file "$(remote_file "$STREAM_BOX_P2_DEF")"
  fi
}

sync_main_lua() {
  local src="$1"
  [[ -f "$src" ]] || die "Local main.lua not found: $src"
  scp "$src" "$STREAM_BOX_HOST:$(remote_file "external\\script\\main.lua")" >/dev/null
}

sync_script() {
  local src="$1"
  local rel="$2"
  [[ -f "$src" ]] || die "Local script not found: $src"
  scp "$src" "$STREAM_BOX_HOST:$(remote_file "$rel")" >/dev/null
}

launch_visible() {
  local local_launcher
  local_launcher="$(mktemp "${TMPDIR:-/tmp}/ikemen-stream-box-launch.XXXXXX.ps1")"
  cat > "$local_launcher" <<EOF
\$ErrorActionPreference = 'SilentlyContinue'
\$streamBoxMuteUi = '$STREAM_BOX_MUTE_UI'
New-Item -ItemType Directory -Force -Path '$REMOTE_APPS_ROOT' | Out-Null
New-Item -ItemType Directory -Force -Path '$RESULT_DIR' | Out-Null
if (\$streamBoxMuteUi -ne '0') {
  \$env:IKEMEN_SUPPRESS_ERROR_DIALOG = '1'
  New-Item -Path 'HKCU:\\Software\\Microsoft\\Windows\\Windows Error Reporting' -Force | Out-Null
  Set-ItemProperty -Path 'HKCU:\\Software\\Microsoft\\Windows\\Windows Error Reporting' -Name DontShowUI -Type DWord -Value 1
  New-Item -Path 'HKCU:\\Software\\Microsoft\\Windows\\Windows Error Reporting\\Consent' -Force | Out-Null
  Set-ItemProperty -Path 'HKCU:\\Software\\Microsoft\\Windows\\Windows Error Reporting\\Consent' -Name DefaultConsent -Type DWord -Value 0
}
Stop-Process -Name Ikemen_GO -Force -ErrorAction SilentlyContinue
Remove-Item '$REMOTE_ROOT\\Ikemen.log' -Force -ErrorAction SilentlyContinue
Remove-Item '$STDOUT_LOG' -Force -ErrorAction SilentlyContinue
Remove-Item '$STDERR_LOG' -Force -ErrorAction SilentlyContinue
Remove-Item '$PID_LOG' -Force -ErrorAction SilentlyContinue
Remove-Item '$EXIT_LOG' -Force -ErrorAction SilentlyContinue
Remove-Item '$RESULT_FILE' -Force -ErrorAction SilentlyContinue
\$ikemenArgs = @(
  '-windowed',
  '-nojoy',
  '-nomusic',
  '-nosound',
  '-jsonstdout',
  '-p1.ai', '$STREAM_BOX_AI_LEVEL',
  '-p2.ai', '$STREAM_BOX_AI_LEVEL',
  '-resultfile',
  '$RESULT_FILE',
  '-p1', '$STREAM_BOX_P1_DEF',
  '-p2', '$STREAM_BOX_P2_DEF',
  '-rounds', '$STREAM_BOX_ROUNDS',
  '-time', '$STREAM_BOX_TIME',
  '-debugstartup'
)
if (\$streamBoxMuteUi -ne '0') {
  \$ikemenArgs += '-noerrordialog'
}
Push-Location '$REMOTE_ROOT'
try {
  \$proc = Start-Process -FilePath '$REMOTE_ROOT\\Ikemen_GO.exe' -WorkingDirectory '$REMOTE_ROOT' -ArgumentList \$ikemenArgs -RedirectStandardOutput '$STDOUT_LOG' -RedirectStandardError '$STDERR_LOG' -PassThru
  \$proc.Id | Set-Content -Path '$PID_LOG' -Encoding ASCII
  if (\$proc.WaitForExit(15000)) {
    \$proc.ExitCode | Set-Content -Path '$EXIT_LOG' -Encoding ASCII
  } else {
    'still-running' | Set-Content -Path '$EXIT_LOG' -Encoding ASCII
  }
} finally {
  Pop-Location
}
EOF
  scp "$local_launcher" "$STREAM_BOX_HOST:$LAUNCH_PS1" >/dev/null
  rm -f "$local_launcher"

  run_remote "schtasks /Create /TN \"$TASK_NAME\" /SC ONCE /ST 23:59 /TR \"powershell -NoProfile -ExecutionPolicy Bypass -WindowStyle Hidden -File $LAUNCH_PS1\" /F" >/dev/null
  run_remote "schtasks /Run /TN \"$TASK_NAME\"" >/dev/null
  sleep "$WAIT_SECONDS"
  run_remote "schtasks /Delete /TN \"$TASK_NAME\" /F" >/dev/null 2>&1 || true
}

post_check() {
  local status
  status="$(run_remote "powershell -NoProfile -Command \"Get-Process Ikemen_GO -ErrorAction SilentlyContinue | Select-Object Id,ProcessName,Responding,StartTime | Format-Table -AutoSize\"" || true)"
  if [[ -n "$status" ]]; then
    printf '%s\n' "$status"
  else
    printf 'No live Ikemen_GO process found after launch.\n'
  fi

  local log_tail
  log_tail="$(run_remote "powershell -NoProfile -Command \"if (Test-Path '$REMOTE_ROOT\\Ikemen.log') { Get-Content '$REMOTE_ROOT\\Ikemen.log' -Tail 80 }\"" || true)"
  if [[ -n "$log_tail" ]]; then
    printf '%s\n' "$log_tail"
  fi

  local stderr_tail
  stderr_tail="$(run_remote "powershell -NoProfile -Command \"if (Test-Path '$STDERR_LOG') { Get-Content '$STDERR_LOG' -Tail 80 }\"" || true)"
  if [[ -n "$stderr_tail" ]]; then
    printf '%s\n' "$stderr_tail"
  fi

  local exit_state
  exit_state="$(run_remote "powershell -NoProfile -Command \"if (Test-Path '$EXIT_LOG') { Get-Content '$EXIT_LOG' }\"" || true)"
  if [[ -n "$exit_state" ]]; then
    printf 'Exit state: %s\n' "$exit_state"
  fi

  local result_candidates=("$RESULT_FILE" "$REMOTE_ROOT\\save\\last-match.json")
  local result_found_path="$RESULT_FILE"
  local result_text=""
  local result_attempt
  local result_candidate
  for result_candidate in "${result_candidates[@]}"; do
    for result_attempt in $(seq 1 "$RESULT_WAIT_SECONDS"); do
      result_text="$(run_remote "powershell -NoProfile -Command \"if (Test-Path '$result_candidate') { Get-Content '$result_candidate' -Raw }\"" || true)"
      if [[ -n "$result_text" ]]; then
        result_found_path="$result_candidate"
        break 2
      fi
      sleep 1
    done
  done
  if [[ -z "$result_text" ]]; then
    die "stream-box quick-vs result file was not created: $RESULT_FILE"
  fi
  if ! printf '%s\n' "$result_text" | grep -q '"statsLog"'; then
    die "stream-box quick-vs result file does not contain the expected JSON payload: $RESULT_FILE"
  fi

  if printf '%s\n%s\n' "$log_tail" "$stderr_tail" | grep -Eiq 'attempt to call a non-function object|runtime error:|panic:|segmentation fault|fatal error|stack traceback'; then
    die "stream-box Ikemen log contains a startup/runtime crash"
  fi

  if [[ "$result_found_path" != "$RESULT_FILE" ]]; then
    run_remote "powershell -NoProfile -Command \"if (Test-Path '$result_found_path') { Copy-Item '$result_found_path' '$RESULT_FILE' -Force }\""
  fi
}

do_check=0
do_launch=0
do_post_check=0
sync_main_lua_path=""
sync_script_args=()

if [[ "$#" -eq 0 ]]; then
  do_check=1
fi

while [[ "$#" -gt 0 ]]; do
  case "$1" in
    --check)
      do_check=1
      ;;
    --sync-main-lua)
      shift
      [[ "$#" -gt 0 ]] || die "--sync-main-lua requires a path"
      sync_main_lua_path="$1"
      ;;
    --sync-script)
      shift
      [[ "$#" -gt 1 ]] || die "--sync-script requires a source path and a remote relative path"
      sync_script_args+=("$1" "$2")
      shift
      ;;
    --launch)
      do_launch=1
      ;;
    --post-check)
      do_post_check=1
      ;;
    --time)
      shift
      [[ "$#" -gt 0 ]] || die "--time requires a value"
      STREAM_BOX_TIME="$1"
      ;;
    --rounds)
      shift
      [[ "$#" -gt 0 ]] || die "--rounds requires a value"
      STREAM_BOX_ROUNDS="$1"
      ;;
    --p1-def)
      shift
      [[ "$#" -gt 0 ]] || die "--p1-def requires a value"
      STREAM_BOX_P1_DEF="$1"
      ;;
    --p2-def)
      shift
      [[ "$#" -gt 0 ]] || die "--p2-def requires a value"
      STREAM_BOX_P2_DEF="$1"
      ;;
    --ai-level)
      shift
      [[ "$#" -gt 0 ]] || die "--ai-level requires a value"
      STREAM_BOX_AI_LEVEL="$1"
      ;;
    --mute-ui)
      STREAM_BOX_MUTE_UI="1"
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

if (( do_check )); then
  check_deploy
fi

if [[ -n "$sync_main_lua_path" ]]; then
  sync_main_lua "$sync_main_lua_path"
fi

if (( ${#sync_script_args[@]} > 0 )); then
  for (( i=0; i<${#sync_script_args[@]}; i+=2 )); do
    sync_script "${sync_script_args[$i]}" "${sync_script_args[$i+1]}"
  done
fi

if (( do_launch )); then
  check_deploy
  ensure_launch_defs
  launch_visible
fi

if (( do_post_check )); then
  post_check
fi

printf 'stream-box Ikemen smoke checks completed.\n'
