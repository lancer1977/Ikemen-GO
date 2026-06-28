# Setup

This section covers local installation, workflow entrypoints, and the self-hosted
runner hooks used by the smoke matrix.

## Prerequisites

Follow the platform-specific package lists in [`BUILDING.md`](../../BUILDING.md):

- Windows/MSYS2: `git`, `make`, `diffutils`, `pkg-config`, Go, the MinGW toolchain,
  `nasm`, `yasm`, `libxmp`, and `SDL2`
- Linux: `git`, `golang-go`, `pkg-config`, `make`, `nasm`, `yasm`,
  `build-essential`, `libxmp-dev`, and `libsdl2-dev`
- macOS: `git`, Go, `pkg-config`, `nasm`, `libxmp`, `sdl2`, and `molten-vk`

## Pages

- [Workflow Matrix](./workflows.md)
