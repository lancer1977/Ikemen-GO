# Render Probe Screen Mapping

## Purpose

Map how Ikemen GO paints each visible screen by rendering short, unique probe codes at known draw points. The operator can run the local build, call out which codes are visible, where they appear, and what they cover or sit behind. Those observations become the evidence for choosing a safe integration point for future screen/HUD reuse work.

## Plan Type

This is a visual render-probe investigation. Use it when the draw stack is easier to understand from on-screen evidence than from static code reading alone.

### Loop

1. Add small, gated probe labels at suspected draw boundaries.
2. Run the game with probes enabled.
3. Visit one screen at a time.
4. Record which labels are visible, where they appear, and whether they are behind or in front of screen elements.
5. Compare observations against code order.
6. Move, remove, or add probes for the next pass.
7. Convert the final map into an implementation seam.

### Observation format

For each screen, record:

- Screen: title / select / versus / fight / pause / post-match / other
- Visible codes:
  - code
  - approximate position
  - in front of / behind what
  - hidden or clipped by what
- Notes:
  - resolution/aspect changes
  - fade behavior
  - cursor/portrait/lifebar ordering
  - whether labels persist across transition frames

## Current Probe Gate

Probes are disabled by default.

Enable them with:

```bash
IKEMEN_RENDER_PROBES=1 ./Ikemen_GO_Linux
```

or during a local Go run/build flow with the same environment variable set.

## Current Probe Families

### Go global gameplay draw probes

These live in `src/system.go` and map the main fight draw stack:

- `G0 system.draw start`
- `G1 after bg fill`
- `G2 after fight -1`
- `G3 after motif -1`
- `G4 after fight 0`
- `G5 after motif 0`
- `G6 after fight 1`
- `G7 after motif 1`
- `G8 after fight 2`
- `G9 after motif 2`
- `GA after motif 3/fade`

### Fight-screen probes

These live in `src/fightscreen.go` and map lifebar/HUD internals:

- `F<n> fight entry`
- `F<n> bars begin`
- `F<n> combo/action`
- `F<n> round`

`<n>` is the fight-screen layer number.

### Motif probes

These live in `src/motif.go` and map motif/screenpack overlay drawing:

- `M<n> motif entry`
- `M<n> menu`
- `M3 fade`

`<n>` is the motif layer number.

### Lua select/versus probes

These live in `external/script/start.lua` and map the script-driven roster and versus screens:

- `S0 select clear`
- `S1 select bg0`
- `S2 select title`
- `S3 select cells`
- `S4 select names`
- `S5 select bg1/top`
- `S6 select stats`
- `V0 versus clear`
- `V1 versus bg0`
- `V2 versus portraits`
- `V3 versus bg1/top`

## Implementation Notes

- `src/render_probe.go` owns the env gate and drawing helper.
- `src/script.go` exposes `renderProbe(label, x, y, r, g, b)` to Lua.
- The probe helper uses the existing debug font and a small colored rectangle so probes are visually distinct but cheap.
- The code path is inert unless `IKEMEN_RENDER_PROBES` is truthy: `1`, `true`, `yes`, or `on`.

## Initial Hypothesis

- Roster/select and versus screens are primarily Lua-driven in `external/script/start.lua`.
- Match HUD/lifebars are drawn through `FightScreen.draw(layerno)` from `System.draw`.
- Motif/screenpack overlays are shared through `Motif.draw(layerno)`.
- A future reusable screen overlay should probably piggyback on the motif layer path for screenpack-style UI, or on the fight-screen layer path only when it is explicitly battle HUD data.

## Control and Snapshot Workflow

Use [`control-snapshot-workflow.md`](control-snapshot-workflow.md) for the repeatable local X11 helper that can launch or attach to IKEMEN, send key input, and save named PNG screenshots for probe evidence passes.

The proven route notes live in [`observations-pass-7-control-snapshot-routes.md`](observations-pass-7-control-snapshot-routes.md). That pass established the private-X/Xvfb control path, the reliable `hold:z:0.25` menu confirmation input, and screenshot evidence from main menu through Arcade character select into a live fight.

Quick start:

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

For reliable menu automation, prefer the private-X wrapper plus `hold:` steps:

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
  --step snap:character-select
```

## Next Evidence Pass

Run the local build with `IKEMEN_RENDER_PROBES=1` or use the control/snapshot helper above, then capture observations for:

1. Select screen roster.
2. Versus screen portraits.
3. Fight screen health/life bars.
4. Pause/menu overlay during fight.
5. Fade/transition frames.

Pass 2 adds a high-visibility `SELECT PATH HIT` sentinel and moves select probes toward the upper-center of the screen. If `SELECT PATH HIT` is absent during a forced interactive select mode, the app is not reaching this Lua select loop.

Pass 2 also adds compact FightScreen subcomponent probes:

- `HP1/HP2 L<n> pre/post`
- `FACE1/FACE2 L<n> pre/post`
- `PWR1/PWR2 L<n> pre/post`
- `NAME1/NAME2 L<n> pre/post`

Use these to compare life bars, portraits, power bars, and names against each other.

Pass 6 adds a targeted top/bottom edge-band probe mode:

```bash
IKEMEN_RENDER_PROBES=edge ./Ikemen_GO_Linux
```

This draws paired `T*` and `B*` blocks across the top and bottom screen bands at select-loop phases from `clearColor` through the final pre-`refresh()` point. Use it to identify whether header/footer overlays should draw after `bgDraw(..., 0)`, after menu/cursor work, or at the final late Lua seam after `bgDraw(..., 1)` and stats.
