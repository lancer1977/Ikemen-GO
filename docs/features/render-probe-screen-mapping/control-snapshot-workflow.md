# Control and Snapshot Workflow

Use this workflow when a visual probe pass needs repeatable movement through Ikemen GO screens plus timestamped evidence screenshots.

## Helper

```bash
scripts/visual/ikemen-control-snapshots.py --help
```

The helper can either launch the local runtime or attach to an already-running IKEMEN window. It uses the local X11 session tools that are already present on this box:

- `wmctrl` to find/focus the game window
- ImageMagick `import` to capture the game window
- `libX11` + `libXtst` via Python `ctypes` to send key events, so `xdotool` is not required

## Default probe pass

From the repo root:

```bash
scripts/visual/ikemen-control-snapshots.py \
  --probe-mode edge \
  --output-dir artifacts/visual-probes/render-probe-pass \
  --step wait:2 \
  --step snap:boot \
  --step key:Return \
  --step wait:3 \
  --step snap:after-enter
```

This launches `./Ikemen_GO_Linux` with `IKEMEN_RENDER_PROBES=edge`, focuses the first window whose title matches `(?i)ikemen`, sends Enter, and writes PNG screenshots under `artifacts/visual-probes/render-probe-pass/`.

## Runtime app tree pass

Use the runnable app tree when the source checkout is not the desired runtime root:

```bash
/home/lancer1977/code/Ikemen-GO/scripts/visual/ikemen-control-snapshots.py \
  --workdir /home/lancer1977/code/ikemen-app \
  --bin ./Ikemen_GO_Linux \
  --probe-mode edge \
  --output-dir /home/lancer1977/code/Ikemen-GO/artifacts/visual-probes/app-tree-pass \
  --step wait:2 \
  --step snap:boot \
  --step key:Return \
  --step wait:3 \
  --step snap:after-enter
```

## Headless/private-X pass

If the desktop X server is unavailable or reports `Maximum number of clients reached`, run the same helper inside a private Xvfb display:

```bash
/home/lancer1977/code/Ikemen-GO/scripts/visual/ikemen-xvfb-control-snapshots.sh \
  --workdir /home/lancer1977/code/ikemen-app \
  --bin ./Ikemen_GO_Linux \
  --probe-mode edge \
  --output-dir /home/lancer1977/code/Ikemen-GO/artifacts/visual-probes/xvfb-pass \
  --step wait:4 \
  --step snap:boot \
  --step key:Return \
  --step wait:2 \
  --step snap:after-enter
```

This starts `Xvfb` plus `metacity`, launches IKEMEN into that private display, captures real window PNGs, and tears the private display down afterward.

## Character-select route

For menu automation, prefer `hold:` over instant `key:` taps. IKEMEN/SDL can miss same-frame press/release events, while a short hold is reliably picked up by the game loop.

The verified route from the default boot/menu into Arcade character select is:

```bash
/home/lancer1977/code/Ikemen-GO/scripts/visual/ikemen-xvfb-control-snapshots.sh \
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

This reaches the `ARCADE` character-select screen with the random slot selected.

## Select header/footer route

The verified route for the pass-8 header/footer proof intentionally confirms Arcade with `Return` before using `z` to enter character select:

```bash
/home/lancer1977/code/Ikemen-GO/scripts/visual/ikemen-xvfb-control-snapshots.sh \
  --workdir /home/lancer1977/code/ikemen-app \
  --bin ./Ikemen_GO_Linux \
  --probe-mode edge-screen \
  --output-dir /home/lancer1977/code/Ikemen-GO/artifacts/visual-probes/select-header-footer-pass \
  --timeout 35 \
  --step wait:8 \
  --step snap:main-menu \
  --step hold:z:0.4 \
  --step wait:2 \
  --step snap:arcade-submenu \
  --step hold:Return:0.4 \
  --step wait:3 \
  --step snap:after-return \
  --step hold:z:0.4 \
  --step wait:5 \
  --step snap:character-select
```

In the verified pass, `20260529-130345-character-select.png` showed the live `ARCADE` character-select screen plus visible cyan/red `SELECT HEADER SCREEN-SPACE` and `SELECT FOOTER SCREEN-SPACE` probe text.

## Fight route

The verified route from boot/menu through character select into a live Arcade fight is:

```bash
/home/lancer1977/code/Ikemen-GO/scripts/visual/ikemen-xvfb-control-snapshots.sh \
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

This selects Arcade, enters the character-select flow, confirms the selected character path, and reaches an active fight. In the verified pass, the final fight screenshot showed Mike vs CPU-S Gill on the Ken Masters bridge stage with lifebars and timer visible.

## Attach to a manually positioned game window

Use attach mode when the game is already running or when launch has to be done manually from a separate runtime tree such as `/home/lancer1977/code/ikemen-app`:

```bash
scripts/visual/ikemen-control-snapshots.py \
  --no-launch \
  --output-dir artifacts/visual-probes/manual-pass \
  --step focus \
  --step snap:current-screen \
  --step key:Down \
  --step key:Return \
  --step wait:1 \
  --step snap:after-selection
```

## Common steps

- `wait:SECONDS` pauses for animation/transition settle time.
- `key:KEY` taps one X11 key name, for example `Return`, `Escape`, `Up`, `Down`, `Left`, `Right`, `a`, `b`, `x`, or `y`. Use this for non-critical navigation.
- `hold:KEY:SECONDS` holds one key long enough for IKEMEN/SDL polling to see it, for example `hold:z:0.25`. Prefer this for menu confirmation and cursor movement.
- `keys:KEY+KEY` sends a chord, for example `Alt+Return`.
- `text:VALUE` types simple text.
- `snap:LABEL` captures a named PNG.
- `focus` re-focuses the target window.

## Render-probe observation loop

1. Pick the probe mode: `IKEMEN_RENDER_PROBES=1` for full labels or `edge` for the top/bottom select-screen edge-band pass.
2. Build a short step list that reaches the target screen.
3. Capture at least one snapshot before and after each transition.
4. Record the visible probe codes in the render-probe observation notes.
5. Keep raw PNGs in `artifacts/visual-probes/...` and summarize the mapping in `docs/features/render-probe-screen-mapping/`.

## Notes and limitations

- This helper is X11-oriented. On Wayland, run from an X11 session or use an XWayland-visible game window.
- It intentionally avoids hard-coded game menu knowledge. Keep per-pass command sequences in shell history, notes, or a small wrapper when a route becomes stable.
- The helper terminates a process it launched unless `--keep-running` is provided. It never closes an attached window in `--no-launch` mode.
- Screenshot files are timestamped so repeated passes do not overwrite previous evidence.
