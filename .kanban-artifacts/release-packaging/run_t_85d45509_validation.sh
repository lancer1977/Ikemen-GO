#!/usr/bin/env bash
set -euo pipefail
STAMP=$(date +%Y%m%d_%H%M%S)
ART=.kanban-artifacts/release-packaging/t_85d45509_run_${STAMP}
mkdir -p "$ART"
LOG="$ART/validation.log"
{
  echo "artifact_dir=$PWD/$ART"
  echo "timestamp=$STAMP"
  echo
  echo '## bash syntax'
  bash -n scripts/package-installers.sh scripts/deploy-local.sh scripts/deploy-local-linux.sh scripts/deploy-local-windows.sh
  echo "exit=$?"
  echo
  echo '## make dry run'
  make -n installers deploy-local deploy-local-windows
  echo "exit=$?"
  echo
  echo '## Windows host guard on Linux (expected failure)'
  WIN_OUT=$(mktemp -d /tmp/ikemen-win-test-XXXXXX)
  set +e
  ./scripts/package-installers.sh --windows --no-build --output-root "$WIN_OUT" >"$ART/windows-host-guard.stdout" 2>"$ART/windows-host-guard.stderr"
  WIN_RC=$?
  set -e
  echo "rc=$WIN_RC output_root=$WIN_OUT"
  head -80 "$ART/windows-host-guard.stderr"
  if [ "$WIN_RC" -eq 0 ]; then echo 'ERROR: windows guard unexpectedly succeeded'; exit 10; fi
  echo
  echo '## screenpack hard fail (expected failure)'
  SP_OUT=$(mktemp -d /tmp/ikemen-screenpack-hard-fail-XXXXXX)
  set +e
  SCREENPACK_REPO=/tmp/does-not-exist-ikemen-screenpack ./scripts/package-installers.sh --linux --no-build --output-root "$SP_OUT" >"$ART/screenpack-hard-fail.stdout" 2>"$ART/screenpack-hard-fail.stderr"
  SP_RC=$?
  set -e
  echo "rc=$SP_RC output_root=$SP_OUT"
  head -80 "$ART/screenpack-hard-fail.stderr"
  if [ "$SP_RC" -eq 0 ]; then echo 'ERROR: screenpack hard fail unexpectedly succeeded'; exit 11; fi
  echo
  echo '## Linux package no-build'
  LINUX_OUT=$(mktemp -d /tmp/ikemen-linux-package-validate-XXXXXX)
  ./scripts/package-installers.sh --linux --no-build --output-root "$LINUX_OUT"
  echo "output_root=$LINUX_OUT"
  unzip -l "$LINUX_OUT/Ikemen_GO-dev-linux.zip" >"$ART/linux-zip-list.txt"
  for required in Ikemen_GO_Linux data/ikemen1/system.def external/script/main.lua LICENSES.txt; do
    if ! grep -q " $required$" "$ART/linux-zip-list.txt"; then echo "ERROR: missing $required"; exit 12; fi
    echo "found $required"
  done
  echo
  echo '## partial screenpack escape hatch'
  PART_OUT=$(mktemp -d /tmp/ikemen-partial-package-validate-XXXXXX)
  ALLOW_PARTIAL_SCREENPACK=1 SCREENPACK_REPO=/tmp/does-not-exist-ikemen-screenpack ./scripts/package-installers.sh --linux --no-build --output-root "$PART_OUT" >"$ART/partial-screenpack.stdout" 2>"$ART/partial-screenpack.stderr"
  echo "output_root=$PART_OUT"
  grep -q 'WARNING: Could not clone or copy the screenpack tree; installer assets will be partial' "$ART/partial-screenpack.stderr"
  test -f "$PART_OUT/Ikemen_GO-dev-linux.zip"
  echo 'partial screenpack warning and archive verified'
  echo
  echo '## workflow static check'
  python3 .kanban-artifacts/release-packaging/workflow_static_check.py
} 2>&1 | tee "$LOG"
printf '\nARTIFACT_DIR=%s\n' "$PWD/$ART"
