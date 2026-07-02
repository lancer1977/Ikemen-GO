#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd -P)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd -P)"

APP_ROOT="${APP_ROOT:-$HOME/apps/ikemen-dev}"
SCREENPACK_REPO="${SCREENPACK_REPO:-https://github.com/ikemen-engine/Ikemen-GO-Screenpack.git}"
SCREENPACK_REF="${SCREENPACK_REF:-master}"
SCREENPACK_DIR="${SCREENPACK_DIR:-}"
SCREENPACK_CACHE="${SCREENPACK_CACHE:-$HOME/apps/_sources/Ikemen-GO-Screenpack}"
DO_BUILD=1
CHECK_DEPS=0

usage() {
  cat <<'EOF'
Usage: ./scripts/deploy-ikemen-dev.sh [options]

Builds Ikemen-GO for Linux and deploys a runnable tree to ~/apps/ikemen-dev.

Options:
  --app-root PATH        Destination runtime folder. Default: ~/apps/ikemen-dev
  --screenpack-dir PATH  Use an existing Ikemen-GO-Screenpack checkout/archive tree
  --screenpack-ref REF   Git ref to checkout when cloning/updating screenpack. Default: master
  --check-deps           Report local build/runtime dependency state and exit
  --no-build             Deploy an already-built ./Ikemen_GO_Linux
  -h, --help             Show this help

Environment overrides:
  APP_ROOT, SCREENPACK_REPO, SCREENPACK_REF, SCREENPACK_DIR, SCREENPACK_CACHE

After deploy:
  ~/apps/ikemen-dev/run-ikemen.sh
EOF
}

die() {
  printf 'ERROR: %s\n' "$*" >&2
  exit 1
}

log() {
  printf '%s\n' "$*"
}

MISSING_PACKAGES=()

add_missing_package() {
  local package="$1"
  [[ -n "$package" ]] || return 0
  MISSING_PACKAGES+=("$package")
}

print_install_missing() {
  if [[ "${#MISSING_PACKAGES[@]}" -eq 0 ]]; then
    printf 'preflight|installMissing|none\n'
    return 0
  fi

  local packages
  packages="$(printf '%s\n' "${MISSING_PACKAGES[@]}" | sort -u | tr '\n' ' ' | sed 's/ $//')"
  printf 'preflight|installMissing|sudo apt-get install -y %s\n' "$packages"
}

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || die "Missing required command: $1"
}

check_command() {
  local name="$1"
  local package_hint="$2"
  if command -v "$name" >/dev/null 2>&1; then
    printf 'dependency|command|%s|present|%s\n' "$name" "$(command -v "$name")"
  else
    add_missing_package "$package_hint"
    printf 'dependency|command|%s|missing|%s\n' "$name" "$package_hint"
  fi
}

check_pkg_config() {
  local name="$1"
  local package_hint="$2"
  if command -v pkg-config >/dev/null 2>&1 && pkg-config --exists "$name"; then
    printf 'dependency|pkg-config|%s|present|%s\n' "$name" "$(pkg-config --modversion "$name" 2>/dev/null || true)"
  else
    add_missing_package "$package_hint"
    printf 'dependency|pkg-config|%s|missing|%s\n' "$name" "$package_hint"
  fi
}

check_file() {
  local label="$1"
  local path="$2"
  if [[ -e "$path" ]]; then
    printf 'dependency|file|%s|present|%s\n' "$label" "$path"
  else
    printf 'dependency|file|%s|missing|%s\n' "$label" "$path"
  fi
}

check_deploy_manifest_hash() {
  local manifest="$APP_ROOT/deploy-manifest.json"
  local binary="$APP_ROOT/Ikemen_GO_Linux"
  if [[ ! -f "$manifest" ]]; then
    printf 'preflight|deployManifestHash|missing|%s\n' "$manifest"
    return 0
  fi
  if [[ ! -f "$binary" ]]; then
    printf 'preflight|deployManifestHash|missing-binary|%s\n' "$binary"
    return 0
  fi
  python3 - "$manifest" "$binary" <<'PY'
import hashlib
import json
import pathlib
import sys

manifest = pathlib.Path(sys.argv[1])
binary = pathlib.Path(sys.argv[2])
data = json.loads(manifest.read_text(encoding="utf-8"))
expected = ((data.get("binary") or {}).get("sha256") or "").strip()
h = hashlib.sha256()
with binary.open("rb") as handle:
    for chunk in iter(lambda: handle.read(1024 * 1024), b""):
        h.update(chunk)
actual = h.hexdigest()
status = "match" if expected == actual else "mismatch"
print(f"preflight|deployManifestHash|{status}|expected={expected}|actual={actual}")
PY
}

run_dependency_check() {
  MISSING_PACKAGES=()

  printf 'preflight|repoRoot|%s\n' "$REPO_ROOT"
  printf 'preflight|appRoot|%s\n' "$APP_ROOT"
  printf 'preflight|screenpackDir|%s\n' "${SCREENPACK_DIR:-}"
  printf 'preflight|screenpackCache|%s\n' "$SCREENPACK_CACHE"

  check_command bash "bash"
  check_command git "git"
  check_command go "golang-go"
  check_command gcc "build-essential"
  check_command make "make"
  check_command nasm "nasm"
  check_command yasm "yasm"
  check_command pkg-config "pkg-config"
  check_command ffmpeg "ffmpeg"
  check_command ffprobe "ffmpeg"
  check_command Xvfb "xvfb"

  check_pkg_config gl "libgl-dev"
  check_pkg_config gtk+-3.0 "libgtk-3-dev"
  check_pkg_config sdl2 "libsdl2-dev"
  check_pkg_config libavformat "libavformat-dev"
  check_pkg_config libavcodec "libavcodec-dev"
  check_pkg_config libavutil "libavutil-dev"
  check_pkg_config libswscale "libswscale-dev"
  check_pkg_config libswresample "libswresample-dev"
  check_pkg_config libavfilter "libavfilter-dev"
  check_pkg_config libxmp "libxmp-dev"

  check_file built-binary "$REPO_ROOT/Ikemen_GO_Linux"
  check_file deployed-binary "$APP_ROOT/Ikemen_GO_Linux"
  check_file deployed-runner "$APP_ROOT/run-ikemen.sh"
  check_file deployed-motif "$APP_ROOT/data/ikemen1/system.def"
  check_file deployed-kfm-character "$APP_ROOT/chars/kfm/kfm.def"
  check_file deployed-kfm-stage "$APP_ROOT/stages/kfm.def"
  check_file deployed-manifest "$APP_ROOT/deploy-manifest.json"
  check_deploy_manifest_hash

  print_install_missing
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --app-root)
      APP_ROOT="${2:-}"
      shift 2
      ;;
    --screenpack-dir)
      SCREENPACK_DIR="${2:-}"
      shift 2
      ;;
    --screenpack-ref)
      SCREENPACK_REF="${2:-}"
      shift 2
      ;;
    --check-deps)
      CHECK_DEPS=1
      shift
      ;;
    --no-build)
      DO_BUILD=0
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      die "Unknown argument: $1"
      ;;
  esac
done

[[ -n "$APP_ROOT" ]] || die "APP_ROOT cannot be empty"

if [[ "$CHECK_DEPS" == "1" ]]; then
  run_dependency_check
  exit 0
fi

ensure_screenpack() {
  if [[ -n "$SCREENPACK_DIR" ]]; then
    [[ -d "$SCREENPACK_DIR" ]] || die "SCREENPACK_DIR does not exist: $SCREENPACK_DIR"
    printf '%s\n' "$SCREENPACK_DIR"
    return 0
  fi

  if [[ -d "$REPO_ROOT/../Ikemen-GO-Screenpack/.git" ]]; then
    printf '%s\n' "$REPO_ROOT/../Ikemen-GO-Screenpack"
    return 0
  fi

  need_cmd git
  mkdir -p "$(dirname "$SCREENPACK_CACHE")"
  if [[ -d "$SCREENPACK_CACHE/.git" ]]; then
    log "==> Updating screenpack cache: $SCREENPACK_CACHE" >&2
    git -C "$SCREENPACK_CACHE" fetch --depth=1 origin "$SCREENPACK_REF" >&2
    git -C "$SCREENPACK_CACHE" checkout -q FETCH_HEAD >&2
  else
    rm -rf "$SCREENPACK_CACHE"
    log "==> Cloning screenpack: $SCREENPACK_REPO" >&2
    git clone --depth=1 -b "$SCREENPACK_REF" "$SCREENPACK_REPO" "$SCREENPACK_CACHE" >&2
  fi
  printf '%s\n' "$SCREENPACK_CACHE"
}

check_build_deps() {
  local missing=()
  need_cmd go
  need_cmd pkg-config
  need_cmd gcc

  pkg-config --exists sdl2 || missing+=("libsdl2-dev")
  pkg-config --exists libavformat || missing+=("libavformat-dev")
  pkg-config --exists libavcodec || missing+=("libavcodec-dev")
  pkg-config --exists libavutil || missing+=("libavutil-dev")
  pkg-config --exists libswscale || missing+=("libswscale-dev")
  pkg-config --exists libswresample || missing+=("libswresample-dev")
  pkg-config --exists libavfilter || missing+=("libavfilter-dev")
  pkg-config --exists gl || missing+=("libgl-dev")
  pkg-config --exists gtk+-3.0 || missing+=("libgtk-3-dev")

  if [[ "${#missing[@]}" -gt 0 ]]; then
    printf 'ERROR: Missing native build packages: %s\n' "${missing[*]}" >&2
    printf 'Install them with:\n  sudo apt install -y %s\n' "${missing[*]}" >&2
    exit 1
  fi
}

sync_dir() {
  local src="$1"
  local dst="$2"
  if command -v rsync >/dev/null 2>&1; then
    rsync -a "$src"/ "$dst"/
  else
    mkdir -p "$dst"
    cp -a "$src"/. "$dst"/
  fi
}

copy_if_exists() {
  local src="$1"
  local dst="$2"
  if [[ -e "$src" ]]; then
    cp -a "$src" "$dst"
  fi
}

write_deploy_manifest() {
  local dst="$1"
  local screenpack="$2"
  local staged_binary="$dst/Ikemen_GO_Linux"
  local deployed_binary="$APP_ROOT/Ikemen_GO_Linux"
  python3 - "$dst/deploy-manifest.json" "$REPO_ROOT" "$APP_ROOT" "$screenpack" "$staged_binary" "$deployed_binary" "$DO_BUILD" "${GOEXPERIMENT:-arenas}" "${CGO_ENABLED:-1}" <<'PY'
import datetime as dt
import hashlib
import json
import pathlib
import subprocess
import sys

manifest_path, repo_root, app_root, screenpack, staged_binary, deployed_binary, do_build, goexperiment, cgo_enabled = sys.argv[1:]
repo = pathlib.Path(repo_root)
staged_binary_path = pathlib.Path(staged_binary)

def git_value(*args: str) -> str | None:
    try:
        return subprocess.check_output(["git", "-C", str(repo), *args], text=True, stderr=subprocess.DEVNULL).strip()
    except Exception:
        return None

def sha256(path: pathlib.Path) -> str | None:
    if not path.exists():
        return None
    h = hashlib.sha256()
    with path.open("rb") as handle:
        for chunk in iter(lambda: handle.read(1024 * 1024), b""):
            h.update(chunk)
    return h.hexdigest()

status = git_value("status", "--short")
payload = {
    "schema": "ikemen-go/dev-deploy-manifest/v1",
    "deployedAtUtc": dt.datetime.now(dt.timezone.utc).isoformat().replace("+00:00", "Z"),
    "appRoot": app_root,
    "sourceRoot": str(repo),
    "sourceCommit": git_value("rev-parse", "HEAD"),
    "sourceBranch": git_value("rev-parse", "--abbrev-ref", "HEAD"),
    "sourceDirty": bool(status),
    "sourceStatusShort": status.splitlines() if status else [],
    "screenpackRoot": screenpack,
    "builtDuringDeploy": do_build == "1",
    "goExperiment": goexperiment,
    "cgoEnabled": cgo_enabled,
    "binary": {
        "path": deployed_binary,
        "size": staged_binary_path.stat().st_size if staged_binary_path.exists() else None,
        "sha256": sha256(staged_binary_path),
    },
}
pathlib.Path(manifest_path).write_text(json.dumps(payload, indent=2) + "\n", encoding="utf-8")
PY
}

cd "$REPO_ROOT"

screenpack="$(ensure_screenpack)"

if [[ "$DO_BUILD" -eq 1 ]]; then
  check_build_deps
  log "==> Building Ikemen_GO_Linux"
  GOEXPERIMENT="${GOEXPERIMENT:-arenas}" CGO_ENABLED="${CGO_ENABLED:-1}" go build -o Ikemen_GO_Linux ./src
else
  [[ -x "$REPO_ROOT/Ikemen_GO_Linux" ]] || die "Missing executable ./Ikemen_GO_Linux; remove --no-build or build first."
fi

mkdir -p "$(dirname "$APP_ROOT")"
tmp_root="$(mktemp -d "${APP_ROOT}.tmp.XXXXXX")"
cleanup() {
  rm -rf "$tmp_root" 2>/dev/null || true
}
trap cleanup EXIT

log "==> Staging runtime tree: $tmp_root"
mkdir -p "$tmp_root"

for path in chars data font sound stages video; do
  if [[ -e "$screenpack/$path" ]]; then
    copy_if_exists "$screenpack/$path" "$tmp_root/"
  fi
done

for path in data font external; do
  if [[ -e "$REPO_ROOT/$path" ]]; then
    if [[ -d "$tmp_root/$path" ]]; then
      sync_dir "$REPO_ROOT/$path" "$tmp_root/$path"
    else
      copy_if_exists "$REPO_ROOT/$path" "$tmp_root/"
    fi
  fi
done

copy_if_exists "$REPO_ROOT/README.md" "$tmp_root/"
copy_if_exists "$REPO_ROOT/LICENCE.txt" "$tmp_root/"
copy_if_exists "$REPO_ROOT/Ikemen_GO_Linux" "$tmp_root/"
copy_if_exists "$REPO_ROOT/src/resources/defaultMotif.ini" "$tmp_root/data/system.base.def"
write_deploy_manifest "$tmp_root" "$screenpack"

cat > "$tmp_root/run-ikemen.sh" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
exec ./Ikemen_GO_Linux "$@"
EOF
chmod +x "$tmp_root/run-ikemen.sh" "$tmp_root/Ikemen_GO_Linux"

rm -rf "$APP_ROOT"
mv "$tmp_root" "$APP_ROOT"
trap - EXIT

log "==> Deployed Ikemen-GO dev runtime to $APP_ROOT"
log "==> Run it with: $APP_ROOT/run-ikemen.sh"
