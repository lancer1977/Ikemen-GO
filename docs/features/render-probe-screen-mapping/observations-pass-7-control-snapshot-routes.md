# Observations Pass 7: Control/Snapshot Routes

## Purpose

Prove that IKEMEN can be controlled from automation, captured as screen evidence, and driven from boot through character select into a live fight without relying on the operator's active desktop session.

This pass is not a render-order probe by itself. It establishes the reusable visual-smoke route that future probe passes can use to reach specific screens and collect repeatable screenshots.

## Environment

- Source checkout: `/home/lancer1977/code/Ikemen-GO`
- Runnable app tree: `/home/lancer1977/code/ikemen-app`
- Runtime binary: `/home/lancer1977/code/ikemen-app/Ikemen_GO_Linux`
- Screenshot root: `/home/lancer1977/code/Ikemen-GO/artifacts/visual-probes/`
- Probe mode used for route validation: `IKEMEN_RENDER_PROBES=edge`
- Display strategy that worked: private `Xvfb` display with `metacity`

## Helpers

Two helpers now own the route:

- `scripts/visual/ikemen-control-snapshots.py`
  - Launches or attaches to an IKEMEN window.
  - Finds and focuses the window with `wmctrl`.
  - Sends key input via `libX11` + `libXtst` from Python `ctypes`.
  - Captures the game window with ImageMagick `import`.
- `scripts/visual/ikemen-xvfb-control-snapshots.sh`
  - Wraps the Python helper in a private `Xvfb` display.
  - Starts `metacity` as a lightweight window manager.
  - Avoids live desktop saturation or focus disruption.

## Key findings

### Live desktop X was not reliable enough for this pass

The first direct run against `DISPLAY=:0` failed before a game window could be captured:

```text
Maximum number of clients reached
Panic: gtk initialisation failed; presumably no X server is available
ERROR: No window matched '(?i)ikemen' within 20.0s. Last windows: no windows listed
```

`xdpyinfo -display :0` also failed with `Maximum number of clients reached`.

Conclusion: when visual smoke is the goal, do not block on the live desktop X server. Use the private-X wrapper instead.

### Private Xvfb path works

The Xvfb route produced real IKEMEN window captures:

```text
Matched window: 0x0060000e pid=<pid> title='Ikemen GO' geom=1280x720+0+0
snapshot: ...boot.png
snapshot: ...after-enter.png
```

The route uses Mesa llvmpipe under Xvfb:

```text
Using OpenGL 4.5 (Core Profile) Mesa ... (llvmpipe ...)
```

This is good enough for menu/screen traversal and probe screenshots.

### IKEMEN/SDL can miss instant key taps

Initial `key:z` and `key:Return` style instant taps were too easy for the game loop to miss. The helper now supports `hold:KEY:SECONDS`, and `hold:z:0.25` is the reliable confirm input for menu automation.

Use:

```bash
--step hold:z:0.25
```

instead of:

```bash
--step key:z
```

for menu confirmation, character confirmation, and cursor movement when reliability matters.

### Main menu can appear without input

A no-input pass showed that the title/load screen advances to the main menu after waiting long enough. The validated route therefore starts with `--step wait:8` and captures `main-menu` before sending confirmation.

## Verified route: boot to main menu

```bash
scripts/visual/ikemen-xvfb-control-snapshots.sh \
  --workdir /home/lancer1977/code/ikemen-app \
  --bin ./Ikemen_GO_Linux \
  --probe-mode edge \
  --output-dir /home/lancer1977/code/Ikemen-GO/artifacts/visual-probes/menu-pass \
  --timeout 25 \
  --step wait:8 \
  --step snap:main-menu
```

Evidence from the successful pass:

- `/home/lancer1977/code/Ikemen-GO/artifacts/visual-probes/xvfb-noinput-test-20260529-124501/20260529-124511-noinput-after-8s.png`

## Verified route: main menu to character select

```bash
scripts/visual/ikemen-xvfb-control-snapshots.sh \
  --workdir /home/lancer1977/code/ikemen-app \
  --bin ./Ikemen_GO_Linux \
  --probe-mode edge \
  --output-dir /home/lancer1977/code/Ikemen-GO/artifacts/visual-probes/character-select-pass \
  --timeout 25 \
  --step wait:8 \
  --step snap:main-menu \
  --step hold:z:0.25 \
  --step wait:2 \
  --step snap:arcade-submenu \
  --step hold:z:0.25 \
  --step wait:5 \
  --step snap:character-select
```

Evidence from the successful pass:

- Main menu: `/home/lancer1977/code/Ikemen-GO/artifacts/visual-probes/xvfb-character-select-test-20260529-124723/20260529-124733-main-menu.png`
- Arcade submenu: `/home/lancer1977/code/Ikemen-GO/artifacts/visual-probes/xvfb-character-select-test-20260529-124723/20260529-124736-arcade-submenu.png`
- Character select: `/home/lancer1977/code/Ikemen-GO/artifacts/visual-probes/xvfb-character-select-test-20260529-124723/20260529-124742-after-single-mode.png`

Observed character-select state:

- `ARCADE` screen reached.
- Character grid visible.
- Random slot selected.
- This proves the route reaches the Lua select-screen loop.

## Verified route: character select to fight

```bash
scripts/visual/ikemen-xvfb-control-snapshots.sh \
  --workdir /home/lancer1977/code/ikemen-app \
  --bin ./Ikemen_GO_Linux \
  --probe-mode edge \
  --output-dir /home/lancer1977/code/Ikemen-GO/artifacts/visual-probes/fight-pass \
  --timeout 25 \
  --step wait:8 \
  --step snap:main-menu \
  --step hold:z:0.25 \
  --step wait:2 \
  --step snap:arcade-submenu \
  --step hold:z:0.25 \
  --step wait:5 \
  --step snap:character-select \
  --step hold:z:0.25 \
  --step wait:3 \
  --step snap:after-select-1 \
  --step hold:z:0.25 \
  --step wait:5 \
  --step snap:after-select-2 \
  --step hold:z:0.25 \
  --step wait:8 \
  --step snap:fight
```

Evidence from the successful pass:

- Main menu: `/home/lancer1977/code/Ikemen-GO/artifacts/visual-probes/xvfb-fight-test-20260529-125009/20260529-125019-main-menu.png`
- Arcade submenu: `/home/lancer1977/code/Ikemen-GO/artifacts/visual-probes/xvfb-fight-test-20260529-125009/20260529-125022-arcade-submenu.png`
- Character select: `/home/lancer1977/code/Ikemen-GO/artifacts/visual-probes/xvfb-fight-test-20260529-125009/20260529-125028-character-select.png`
- Post-selection transition: `/home/lancer1977/code/Ikemen-GO/artifacts/visual-probes/xvfb-fight-test-20260529-125009/20260529-125032-after-select-1.png`
- Versus/loading path: `/home/lancer1977/code/Ikemen-GO/artifacts/visual-probes/xvfb-fight-test-20260529-125009/20260529-125038-after-select-2.png`
- Fight: `/home/lancer1977/code/Ikemen-GO/artifacts/visual-probes/xvfb-fight-test-20260529-125009/20260529-125047-after-select-3.png`

Observed fight state:

- Active fight reached.
- Mike vs CPU-S Gill.
- Lifebars visible.
- Timer visible.
- Ken Masters bridge/yacht background loaded.
- Characters spawned and idle in match.

## Runtime warnings seen during successful passes

These appeared during successful runs and did not block the route:

```text
WARNING: Animation missing sprite 9000,3 from data/ikemen1/system.sff
WARNING: Duplicate action key in data/fight.def: 510 (ignored)
Failed to open BGM: open sound/Title.mp3: no such file or directory
Failed to open BGM: open sound/Select.mp3: no such file or directory
Failed to open BGM: open sound/Versus.mp3: no such file or directory
WARNING: Gill: Duplicate 'pos' parameter in state 250
```

Treat these as route-tolerated warnings unless a later pass is specifically validating assets/audio/config hygiene.

## Recommended default for future visual probes

Use the private-X wrapper and `hold:`-based confirmation:

```bash
scripts/visual/ikemen-xvfb-control-snapshots.sh \
  --workdir /home/lancer1977/code/ikemen-app \
  --bin ./Ikemen_GO_Linux \
  --probe-mode edge \
  --output-dir /home/lancer1977/code/Ikemen-GO/artifacts/visual-probes/$(date +%Y%m%d-%H%M%S)-probe-pass \
  --timeout 25 \
  --step wait:8 \
  --step snap:main-menu \
  --step hold:z:0.25 \
  --step wait:2 \
  --step snap:arcade-submenu \
  --step hold:z:0.25 \
  --step wait:5 \
  --step snap:character-select
```

Extend that route with additional `hold:` and `snap:` steps for versus, fight, pause, and post-match mapping.

## Follow-up opportunities

- Add a small named shell wrapper for common routes if they become frequent (`menu`, `character-select`, `fight`).
- Add OCR or simple image-shape checks later if we want CI-style pass/fail assertions instead of manual screenshot inspection.
- Consider a route that moves off the random slot and selects a deterministic character/stage if the current random/default behavior becomes too variable for probe comparison.
- Continue the render-probe mapping from the now-proven route into:
  - select-screen edge bands
  - versus portraits
  - fight HUD/lifebars
  - pause overlay
  - transition/fade frames
