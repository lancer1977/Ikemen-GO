# Render Probe Screen Mapping - Pass 3 Plan

## Purpose

Reduce the pass-2 probe noise and isolate the missing select-screen evidence.

## Changes

- `IKEMEN_RENDER_PROBES=1`, `true`, `yes`, or `on` now maps to `all`.
- `IKEMEN_RENDER_PROBES=all` shows every probe category.
- `IKEMEN_RENDER_PROBES=select` shows only select probes.
- `IKEMEN_RENDER_PROBES=versus` shows only versus probes.
- `IKEMEN_RENDER_PROBES=global` shows only broad `System.draw` ladder probes.
- `IKEMEN_RENDER_PROBES=motif` shows only motif/menu/fade probes.
- `IKEMEN_RENDER_PROBES=fight` shows general fight probes and all `fight-*` subcategories.
- `IKEMEN_RENDER_PROBES=fight-health` shows life/healthbar pre/post probes.
- `IKEMEN_RENDER_PROBES=fight-face` shows portrait/face pre/post probes.
- `IKEMEN_RENDER_PROBES=fight-power` shows powerbar pre/post probes.
- `IKEMEN_RENDER_PROBES=fight-name` shows name pre/post probes.

Select instrumentation now has three path sentinels:

- `SELECT FUNC ENTER` at the top of `start.f_selectScreen()`.
- `SELECT EARLY RETURN` immediately before the select-screen early return.
- `SELECT PATH HIT` inside the draw loop.

## Current Relaunch

The app was relaunched from `/home/lancer1977/code/ikemen-app` with:

```bash
IKEMEN_RENDER_PROBES=select ./Ikemen_GO_Linux
```

This should suppress the fight-screen probe noise and show only select probes if the select path is visible.

## What To Observe

On any select/menu path, look for:

- `SELECT FUNC ENTER`
- `SELECT EARLY RETURN`
- `SELECT PATH HIT`
- `S0 select clear`
- `S1 select bg0`
- `S2 select title`
- `S3 select cells`
- `S4 select names`
- `S5 select bg1/top`
- `S6 select stats` if using the newer source script path

Interpretation:

- `SELECT FUNC ENTER` only: the function is called, but the label is likely overwritten before frame present or the early return path is not visible long enough.
- `SELECT FUNC ENTER` + `SELECT EARLY RETURN`: select function is called but bypassed by `main.selectMenu` / `selScreenEnd` state.
- `SELECT PATH HIT`: the select draw loop is active.
- No select labels at all: current visible flow does not call/render this select function, or labels are not presented before another frame clears them.
