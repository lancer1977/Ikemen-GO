#!/bin/bash

set -euo pipefail

if ! command -v apt-get >/dev/null 2>&1; then
  echo "This helper expects an apt-based Linux system." >&2
  exit 1
fi

sudo apt-get update
sudo apt-get install -y \
  git \
  golang-go \
  pkg-config \
  build-essential \
  make \
  nasm \
  yasm \
  xvfb \
  libxmp-dev \
  libsdl2-dev \
  libgl-dev \
  libgtk-3-dev \
  ffmpeg \
  libavcodec-dev \
  libavformat-dev \
  libavutil-dev \
  libswscale-dev \
  libswresample-dev \
  libavfilter-dev

pkg-config --modversion \
  gl \
  sdl2 \
  gtk+-3.0 \
  libavformat \
  libavcodec \
  libavutil \
  libswresample \
  libswscale \
  libavfilter \
  libxmp

echo
echo "Dependencies installed. Ikemen-GO builds may require: GOEXPERIMENT=arenas"
