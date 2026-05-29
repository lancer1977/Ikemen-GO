# Render Probe Screen Mapping - Pass 6 Edge Band Plan

## Purpose

Pass 5 confirmed that the select screen can be painted through the instrumented select/batch draw path, but it mostly proved the portrait/roster picker area. Pass 6 isolates the top and bottom screen bands so we can learn which select-screen draw phase can safely paint HUD/header/footer overlays outside the portraits and cell picker.

## Run Mode

Run from the local app/build directory with:

```bash
IKEMEN_RENDER_PROBES=edge ./Ikemen_GO_Linux
```

This mode draws only the new top/bottom edge probes. Normal gameplay remains unchanged when `IKEMEN_RENDER_PROBES` is unset.

## What Changed

`external/script/start.lua` now draws paired top (`T*`) and bottom (`B*`) blocks at specific points in the select loop:

| Code | Phase | Color | Meaning |
|---|---|---|---|
| `T0 clear` / `B0 clear` | immediately after `clearColor` | red | earliest paint after clear; likely covered by backgrounds |
| `T1 bg0` / `B1 bg0` | after `bgDraw(..., 0)` | orange | paint above base background layer |
| `T2 title` / `B2 title` | after title draw | yellow-orange | paint after select title |
| `T3 after-menu` / `B3 after-menu` | after cursor/team/select menu block | cyan | paint after active picker/menu work, before names and completion/timer/hook |
| `T4 pre-bg1` / `B4 pre-bg1` | just before `bgDraw(..., 1)` | white | tests whether top/background layer 1 covers edge overlays |
| `T5 final` / `B5 final` | after `bgDraw(..., 1)` and select stats overlay, before `refresh()` | magenta | latest Lua-side paint before frame present; best candidate for topmost edge overlay |

## What To Observe

On the character select screen, report:

1. Which top blocks are visible: `T0` through `T5`.
2. Which bottom blocks are visible: `B0` through `B5`.
3. Whether any block is partially clipped at the screen edge.
4. Whether any block appears behind the title, roster, portraits, names, timer, or screenpack top layer.
5. Whether `T5 final` / `B5 final` sit above everything or get covered by something after they draw.

## Interpretation

- If only `T5`/`B5` are visible, then the final pre-refresh seam is the safe top/bottom overlay seam.
- If `T3`/`B3` remain visible, overlays can be drawn after menu/cursor work and before late screenpack top layers.
- If `T4`/`B4` disappear but `T5`/`B5` show, `bgDraw(..., 1)` is covering edge overlays and should be treated as a top-layer screenpack pass.
- If top and bottom differ, the screenpack may have asymmetric top/footer art or clipping.

## Likely Implementation Direction

If pass 6 confirms `T5 final` and `B5 final`, the production seam should be a small reusable select-screen overlay hook late in `start.f_selectScreen()`, after layer-1 background and stats overlay, before fade/refresh. That would support header/footer paint without coupling to roster cells, portraits, or FightScreen HUD internals.
