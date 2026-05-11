#!/bin/bash

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd -P)"
cd "$REPO_ROOT"

if ! git diff --quiet || ! git diff --cached --quiet; then
  echo "Working tree must be clean before syncing develop." >&2
  exit 1
fi

git fetch upstream develop
git checkout develop
git merge --ff-only upstream/develop
