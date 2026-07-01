#!/usr/bin/env bash
set -u -o pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
cd "$ROOT"

RUN_GO=0
RUN_PYTHON=1
RUN_SHELL_SYNTAX=1
RUN_SMOKE=0
RUN_STREAM_BOX=0
KEEP_GOING=1

usage() {
  cat <<'EOF'
Usage: ./test.sh [options]

Runs the Ikemen-GO test suite from the repository root.

Default suite:
  - Bash syntax checks for repo test/smoke/deploy scripts
  - Python unittest discovery under scripts/tests

Optional checks and integration suites:
  --go-check    Run Go package compile/test check with GOEXPERIMENT=arenas
  --smoke       Run local Ikemen smoke tests if prerequisites exist
  --stream-box  Run Windows stream-box smoke checks over SSH if available
  --all         Run default suite plus --go-check, --smoke, and --stream-box

Selection options:
  --no-go       Skip Go package compile/test check even with --all
  --no-python   Skip Python tests
  --no-shell    Skip Bash syntax checks
  --fail-fast   Stop at the first failing section
  -h, --help    Show this help

Environment overrides:
  IKEMEN_BIN                         Binary for smoke scripts. Default: ./Ikemen_GO_Linux
  GOEXPERIMENT                       Defaults to arenas for go test
  IKEMEN_WORKFLOW_FIXTURE_ROOT       Fixture root for workflow smoke/char sweeps
  IKEMEN_LOCAL_KFM_FIXTURE_ROOT      Fixture root for local KFM smoke
  STREAM_BOX_HOST                    SSH alias for stream-box smoke. Default: stream-box
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --go-check)
      RUN_GO=1
      ;;
    --smoke)
      RUN_SMOKE=1
      ;;
    --stream-box)
      RUN_STREAM_BOX=1
      ;;
    --all)
      RUN_GO=1
      RUN_SMOKE=1
      RUN_STREAM_BOX=1
      ;;
    --no-go)
      RUN_GO=0
      ;;
    --no-python)
      RUN_PYTHON=0
      ;;
    --no-shell)
      RUN_SHELL_SYNTAX=0
      ;;
    --fail-fast)
      KEEP_GOING=0
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "ERROR: Unknown option: $1" >&2
      usage >&2
      exit 2
      ;;
  esac
  shift
done

failures=0
skips=0

section() {
  printf '\n==> %s\n' "$*"
}

pass() {
  printf 'PASS: %s\n' "$*"
}

skip() {
  skips=$((skips + 1))
  printf 'SKIP: %s\n' "$*"
}

fail() {
  failures=$((failures + 1))
  printf 'FAIL: %s\n' "$*" >&2
  if [[ "$KEEP_GOING" -eq 0 ]]; then
    exit 1
  fi
}

run_cmd() {
  local name="$1"
  shift
  section "$name"
  printf '+ '
  printf '%q ' "$@"
  printf '\n'
  if "$@"; then
    pass "$name"
  else
    fail "$name"
  fi
}

run_shell_syntax() {
  if [[ "$RUN_SHELL_SYNTAX" -ne 1 ]]; then
    skip "Bash syntax checks disabled"
    return
  fi

  section "Bash syntax checks"
  local files=()
  local path
  for path in \
    test.sh \
    build/*.sh \
    scripts/*.sh \
    scripts/smoke/*.sh \
    scripts/stream-box/*.sh; do
    if compgen -G "$path" >/dev/null; then
      while IFS= read -r file; do
        files+=("$file")
      done < <(compgen -G "$path" | sort)
    fi
  done

  if [[ "${#files[@]}" -eq 0 ]]; then
    skip "No Bash scripts found"
    return
  fi

  local failed=0
  for path in "${files[@]}"; do
    if bash -n "$path"; then
      printf 'ok   %s\n' "$path"
    else
      printf 'bad  %s\n' "$path" >&2
      failed=1
    fi
  done

  if [[ "$failed" -eq 0 ]]; then
    pass "Bash syntax checks (${#files[@]} files)"
  else
    fail "Bash syntax checks"
  fi
}

run_python_tests() {
  if [[ "$RUN_PYTHON" -ne 1 ]]; then
    skip "Python tests disabled"
    return
  fi
  if [[ ! -d scripts/tests ]]; then
    skip "No scripts/tests directory"
    return
  fi
  run_cmd "Python unittest suite" env PYTHONDONTWRITEBYTECODE=1 PYTHONPATH="$ROOT" python3 -m unittest discover -s scripts/tests -p 'test_*.py'
}

run_go_tests() {
  if [[ "$RUN_GO" -ne 1 ]]; then
    skip "Go package compile/test check not requested; use --go-check or --all"
    return
  fi
  if ! command -v go >/dev/null 2>&1; then
    skip "go is not on PATH"
    return
  fi

  # Ikemen-GO currently imports Go's arena package, so enable the experiment unless
  # the caller intentionally overrides GOEXPERIMENT.
  run_cmd "Go package tests" env GOEXPERIMENT="${GOEXPERIMENT:-arenas}" go test ./...
}

have_local_smoke_prereqs() {
  local bin="${IKEMEN_BIN:-$ROOT/Ikemen_GO_Linux}"
  [[ -x "$bin" ]] || return 1
  [[ -d "${IKEMEN_WORKFLOW_FIXTURE_ROOT:-/mnt/shared/games/ikemen2}" || -d /mnt/shared/Emu/ikemen ]] || return 1
  return 0
}

run_local_smoke() {
  if [[ "$RUN_SMOKE" -ne 1 ]]; then
    skip "Local smoke tests not requested; use --smoke or --all"
    return
  fi
  if ! have_local_smoke_prereqs; then
    skip "Local smoke prerequisites missing: executable Ikemen_GO_Linux and fixture root"
    return
  fi

  run_cmd "Workflow smoke matrix" scripts/smoke/ikemen-workflow-matrix.sh
  run_cmd "Local KFM smoke" scripts/smoke/ikemen-local-kfm-smoke.sh
}

run_stream_box_smoke() {
  if [[ "$RUN_STREAM_BOX" -ne 1 ]]; then
    skip "Stream-box smoke not requested; use --stream-box or --all"
    return
  fi
  local host="${STREAM_BOX_HOST:-stream-box}"
  if ! command -v ssh >/dev/null 2>&1; then
    skip "ssh is not on PATH"
    return
  fi
  if ! ssh -o BatchMode=yes -o ConnectTimeout=5 "$host" "echo ok" >/dev/null 2>&1; then
    skip "Stream-box SSH host unavailable: $host"
    return
  fi

  run_cmd "Stream-box smoke check" scripts/smoke/stream-box-ikemen-smoke.sh --check --post-check
}

run_shell_syntax
run_python_tests
run_go_tests
run_local_smoke
run_stream_box_smoke

printf '\n==> Test summary\n'
printf 'Failures: %d\n' "$failures"
printf 'Skipped:  %d\n' "$skips"

if [[ "$failures" -ne 0 ]]; then
  exit 1
fi
exit 0
