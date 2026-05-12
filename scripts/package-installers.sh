#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd -P)"
cd "$REPO_ROOT"

OUTPUT_ROOT="${OUTPUT_ROOT:-/mnt/syn1/games/Ikemen}"
BUILD_FFMPEG="${BUILD_FFMPEG:-yes}"
APP_VERSION="${APP_VERSION:-nightly}"
APP_BUILDTIME="${APP_BUILDTIME:-$(date '+%Y.%m.%d')}"
SCREENPACK_REPO="${SCREENPACK_REPO:-https://github.com/ikemen-engine/Ikemen-GO-Screenpack.git}"
SCREENPACK_REF="${SCREENPACK_REF:-master}"
SCREENPACK_DIR="${SCREENPACK_DIR:-}"

TARGETS=("linux" "windows")
DO_BUILD=1

usage() {
  cat <<'EOF'
Usage: ./scripts/package-installers.sh [options]

Builds and packages desktop installers into an output directory.
Default output root: /mnt/syn1/games/Ikemen

Options:
  --output-root PATH   Write final zip files to PATH
  --linux              Only package Linux
  --windows            Only package Windows
  --all                Package both Linux and Windows
  --no-build           Package from existing build outputs only
  -h, --help           Show this help

Environment overrides:
  OUTPUT_ROOT, APP_VERSION, APP_BUILDTIME, BUILD_FFMPEG,
  SCREENPACK_REPO, SCREENPACK_REF
EOF
}

cleanup_items=()
cleanup() {
  local item
  for item in "${cleanup_items[@]:-}"; do
    rm -rf "$item" 2>/dev/null || true
  done
}
trap cleanup EXIT

log() {
  printf '%s\n' "$*"
}

warn() {
  printf 'WARNING: %s\n' "$*" >&2
}

die() {
  printf 'ERROR: %s\n' "$*" >&2
  exit 1
}

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || die "Missing required command: $1"
}

host_is_windows() {
  case "$OSTYPE" in
    msys|cygwin) return 0 ;;
    *) return 1 ;;
  esac
}

host_is_linux() {
  case "$OSTYPE" in
    linux*) return 0 ;;
    *) return 1 ;;
  esac
}

require_writable_dir() {
  local dir="$1"
  mkdir -p "$dir" 2>/dev/null || return 1
  local probe="$dir/.codex-write-test.$$"
  if ! : > "$probe" 2>/dev/null; then
    return 1
  fi
  rm -f "$probe"
}

download_file() {
  local url="$1" out="$2"
  mkdir -p "$(dirname "$out")"
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$url" -o "$out"
  elif command -v wget >/dev/null 2>&1; then
    wget -qO "$out" "$url"
  else
    return 1
  fi
}

zip_dir() {
  local src_dir="$1" dest_zip="$2"
  rm -f "$dest_zip"
  if command -v zip >/dev/null 2>&1; then
    (cd "$src_dir" && zip -r "$dest_zip" . >/dev/null)
    return 0
  fi
  if command -v 7z >/dev/null 2>&1; then
    (cd "$src_dir" && 7z a -tzip "$dest_zip" . >/dev/null)
    return 0
  fi
  if command -v 7z.exe >/dev/null 2>&1; then
    (cd "$src_dir" && 7z.exe a -tzip "$dest_zip" . >/dev/null)
    return 0
  fi
  die "Need zip or 7z to create archives"
}

ensure_screenpack() {
  local dst="${SCREENPACK_DIR:-}"
  if [[ -n "$dst" && -d "$dst/.git" ]]; then
    echo "$dst"
    return 0
  fi

  if [[ -z "$dst" ]]; then
    dst="$(mktemp -d "$REPO_ROOT/build/screenpack.XXXXXX")"
  else
    rm -rf "$dst"
    mkdir -p "$(dirname "$dst")"
  fi

  if [[ -d "$SCREENPACK_REPO" ]]; then
    cp -a "$SCREENPACK_REPO" "$dst"
    echo "$dst"
    return 0
  fi

  if git clone --depth=1 -b "$SCREENPACK_REF" "$SCREENPACK_REPO" "$dst" >/dev/null 2>&1; then
    echo "$dst"
    return 0
  fi

  return 1
}

generate_licenses() {
  local out="$1"
  local screenpack_dir="${2:-}"
  local tmp=""
  tmp="$(mktemp)"
  cleanup_items+=("$tmp")
  fetch_or_note() {
    local label="$1" url="$2"
    echo "$label"
    echo "-----------------------------"
    if download_file "$url" "$tmp"; then
      cat "$tmp"
    else
      echo "Unable to fetch $url"
    fi
    echo
  }
  {
    echo "============================="
    echo " Ikemen GO - LICENSES.txt"
    echo "============================="
    echo
    echo "Section 1: Ikemen GO (MIT)"
    echo "--------------------------"
    cat LICENCE.txt
    echo
    echo "Section 2: Bundled Assets"
    echo "------------------------"
    if [[ -n "$screenpack_dir" && -f "$screenpack_dir/LICENCE.txt" ]]; then
      cat "$screenpack_dir/LICENCE.txt"
    else
      echo "Screenpack assets are licensed separately."
      echo "See https://github.com/ikemen-engine/Ikemen-GO-Screenpack for details."
    fi
    echo
    echo "Section 3: Third-party runtime notes"
    echo "-----------------------------------"
    echo "FFmpeg is used under LGPL v2.1."
    fetch_or_note "FFmpeg LGPL v2.1" "https://raw.githubusercontent.com/FFmpeg/FFmpeg/release/7.1/COPYING.LGPLv2.1"
    fetch_or_note "FFmpeg LICENSE.md" "https://raw.githubusercontent.com/FFmpeg/FFmpeg/release/7.1/LICENSE.md"
    fetch_or_note "FFmpeg CREDITS" "https://raw.githubusercontent.com/FFmpeg/FFmpeg/release/7.1/CREDITS"
    fetch_or_note "libxmp README" "https://raw.githubusercontent.com/libxmp/libxmp/master/README"
    fetch_or_note "SDL2 LICENSE" "https://raw.githubusercontent.com/libsdl-org/SDL/main/LICENSE.txt"
    echo
  } > "$out"
}

stage_common_assets() {
  local deploy="$1"
  local screenpack_dir="${2:-}"

  mkdir -p "$deploy/lib"
  cp -a README.md "$deploy/"
  cp -a build/Ikemen_GO.command "$deploy/" 2>/dev/null || true
  cp -a LICENCE.txt "$deploy/" 2>/dev/null || true
  cp -a data "$deploy/"
  cp -a font "$deploy/"
  cp -a external "$deploy/"
  cp -a src/resources/defaultMotif.ini "$deploy/data/system.base.def"

  if [[ -n "$screenpack_dir" ]]; then
    for path in chars data font sound stages video; do
      if [[ -e "$screenpack_dir/$path" ]]; then
        cp -a "$screenpack_dir/$path" "$deploy/"
      fi
    done
  fi
}

stage_linux_package() {
  local root="$1"
  local screenpack_dir="${2:-}"
  local deploy="$root/deploy-linux"
  rm -rf "$deploy"
  mkdir -p "$deploy"

  cp -a Ikemen_GO_Linux "$deploy/"
  cp -a lib "$deploy/"
  cp -a external/icons/IkemenCylia_256.png "$deploy/ikemen-go.png" 2>/dev/null || true
  cat > "$deploy/Ikemen_GO.desktop" <<'EOF'
[Desktop Entry]
Name=Ikemen GO
Comment=An open-source fighting game engine
TryExec=Ikemen_GO_Linux
Exec=./Ikemen_GO_Linux
Icon=ikemen-go.png
Terminal=false
Type=Application
Categories=Game;ArcadeGame;
EOF

  stage_common_assets "$deploy" "$screenpack_dir"
  generate_licenses "$deploy/LICENSES.txt" "$screenpack_dir"
  zip_dir "$deploy" "$root/Ikemen_GO-dev-linux.zip"
}

stage_windows_package() {
  local root="$1"
  local screenpack_dir="${2:-}"
  local deploy="$root/deploy-windows"
  rm -rf "$deploy"
  mkdir -p "$deploy"

  cp -a Ikemen_GO.exe "$deploy/"
  cp -a lib "$deploy/"
  stage_common_assets "$deploy" "$screenpack_dir"
  generate_licenses "$deploy/LICENSES.txt" "$screenpack_dir"
  zip_dir "$deploy" "$root/Ikemen_GO-dev-windows.zip"
}

build_linux() {
  log "==> Building Linux binary"
  BUILD_FFMPEG="$BUILD_FFMPEG" APP_VERSION="$APP_VERSION" APP_BUILDTIME="$APP_BUILDTIME" ./build/build.sh Linux
}

build_windows() {
  log "==> Building Windows binary"
  BUILD_FFMPEG="$BUILD_FFMPEG" APP_VERSION="$APP_VERSION" APP_BUILDTIME="$APP_BUILDTIME" ./build/build.sh Win64
}

ensure_target_supported() {
  local target="$1"
  case "$target" in
    linux)
      host_is_linux || die "Linux packaging requires a Linux host. Use scripts/deploy-local-linux.sh on Linux."
      ;;
    windows)
      if ! host_is_windows; then
        die "Windows packaging requires Windows or MSYS2/MINGW64. Use scripts/deploy-local-windows.sh on that host."
      fi
      command -v x86_64-w64-mingw32-gcc >/dev/null 2>&1 || die "Windows packaging needs x86_64-w64-mingw32-gcc on PATH."
      ;;
  esac
}

main() {
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --output-root)
        OUTPUT_ROOT="${2:-}"
        shift 2
        ;;
      --linux)
        TARGETS=("linux")
        shift
        ;;
      --windows)
        TARGETS=("windows")
        shift
        ;;
      --all)
        TARGETS=("linux" "windows")
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

  [[ -z "$OUTPUT_ROOT" ]] && die "OUTPUT_ROOT cannot be empty"

  need_cmd git
  need_cmd go
  need_cmd pkg-config

  if ! require_writable_dir "$OUTPUT_ROOT"; then
    die "Output root is not writable: $OUTPUT_ROOT"
  fi

  local screenpack_dir=""
  if screenpack_dir="$(ensure_screenpack)"; then
    :
  else
    warn "Could not clone or copy the screenpack tree; installer assets will be partial"
    screenpack_dir=""
  fi

  local failures=0
  local target
  for target in "${TARGETS[@]}"; do
    ensure_target_supported "$target"
    case "$target" in
      linux)
        if [[ "$DO_BUILD" -eq 1 ]]; then
          if ! build_linux; then
            warn "Linux build failed"
            failures=1
            continue
          fi
        fi
        if ! stage_linux_package "$OUTPUT_ROOT" "$screenpack_dir"; then
          warn "Linux packaging failed"
          failures=1
          continue
        fi
        log "==> Linux installer written to $OUTPUT_ROOT/Ikemen_GO-dev-linux.zip"
        ;;
      windows)
        if [[ "$DO_BUILD" -eq 1 ]]; then
          if ! build_windows; then
            warn "Windows build failed"
            failures=1
            continue
          fi
        fi
        if ! stage_windows_package "$OUTPUT_ROOT" "$screenpack_dir"; then
          warn "Windows packaging failed"
          failures=1
          continue
        fi
        log "==> Windows installer written to $OUTPUT_ROOT/Ikemen_GO-dev-windows.zip"
        ;;
      *)
        die "Unknown target: $target"
        ;;
    esac
  done

  if [[ "$failures" -ne 0 ]]; then
    die "One or more installer targets failed"
  fi
  if [[ "$screenpack_dir" == "$REPO_ROOT"/build/screenpack.* ]]; then
    rm -rf "$screenpack_dir" 2>/dev/null || true
  fi
}

main "$@"
