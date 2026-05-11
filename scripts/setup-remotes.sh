#!/bin/bash

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd -P)"
cd "$REPO_ROOT"

FORK_URL="${1:-git@github.com:lancer1977/Ikemen-GO.git}"
UPSTREAM_URL="${2:-git@github.com:ikemen-engine/Ikemen-GO.git}"

if ! git remote get-url origin >/dev/null 2>&1; then
  git remote add origin "$FORK_URL"
else
  git remote set-url origin "$FORK_URL"
fi

if ! git remote get-url upstream >/dev/null 2>&1; then
  git remote add upstream "$UPSTREAM_URL"
else
  git remote set-url upstream "$UPSTREAM_URL"
fi

git remote -v
