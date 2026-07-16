#!/usr/bin/env bash
set -euo pipefail

mapfile -d '' go_files < <(git ls-files -z -- '*.go')

if ((${#go_files[@]} == 0)); then
  echo "No tracked Go files found."
  exit 0
fi

gofmt_output="$(mktemp)"
trap 'rm -f "${gofmt_output}"' EXIT

if ! gofmt -l "${go_files[@]}" >"${gofmt_output}"; then
  echo "gofmt could not parse or inspect the tracked Go files." >&2
  exit 1
fi

mapfile -t unformatted <"${gofmt_output}"

if ((${#unformatted[@]} > 0)); then
  echo "The following tracked Go files are not formatted with gofmt:" >&2
  printf '  %s\n' "${unformatted[@]}" >&2
  echo "Run: gofmt -w <files listed above>" >&2
  exit 1
fi

echo "All ${#go_files[@]} tracked Go files are formatted."
