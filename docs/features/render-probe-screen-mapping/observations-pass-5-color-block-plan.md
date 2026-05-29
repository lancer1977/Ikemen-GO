# Render Probe Screen Mapping - Pass 5 Color Block Plan

## Purpose

Words were unreadable on the 3D-ish select screen, but a probe rectangle was visible. Pass 5 uses large color blocks and one-letter identifiers instead of long labels.

## Run Mode

```bash
IKEMEN_RENDER_PROBES=swarm-block ./Ikemen_GO_Linux
```

## Block Legend

Lua select loop blocks:

- `A` bright green: before roster/cell `batchDraw(staticDrawList)`.
- `B` dark green: after roster/cell `batchDraw(staticDrawList)`.
- `C` yellow: before selected portrait block (`start.f_drawPortraits` calls from select loop).
- `D` orange: after selected portrait block.
- `F` blue: before team/select menu and active cursor handling.
- `N` purple: before selected-name drawing.
- `G` white: before select background layer 1/top layer draw.
- `Z` magenta: final marker before `refresh()`.

Go batch draw queued-layer blocks:

- `E<layer>` cyan/green: inside the queued `batchDraw` layer callback, immediately before the projected animation draw for that layer.

## What To Report

On select screen, report which colored blocks are visible and what they overlap:

- roster cells / inner select box
- selected character portrait area
- angled/secondary portraits
- cursor/team menu
- names
- top overlay/background

If `A/B/E*` are visible on the roster but `C/D` are not, the roster is the projected batch path and selected portraits are either covered or drawn elsewhere.

If `C/D` are visible on selected portrait areas, `start.f_drawPortraits()` owns that surface.

If `G/Z` are visible on top of everything, the late top/background/refresh path is usable for global overlays.

If only colored rectangles but no letters are visible, rely on color and location rather than text.
