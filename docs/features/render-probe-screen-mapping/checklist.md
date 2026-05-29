# Render Probe Screen Mapping Checklist

## Setup

- [x] Document the visual render-probe plan type.
- [x] Add an environment-gated render probe helper.
- [x] Expose `renderProbe(...)` to Lua.
- [x] Add probes to the main gameplay draw stack.
- [x] Add probes to fight-screen/HUD drawing.
- [x] Add probes to motif drawing.
- [x] Add probes to select and versus screen Lua drawing.
- [x] Run the probed build locally with probes enabled.
- [x] Record first-pass observations from the operator.
- [x] Record second-pass select/fight observations from the operator.
- [x] Record fourth-pass swarm select observation from the operator.
- [x] Record fifth-pass color-block select observation from the operator.
- [x] Add sixth-pass edge-band probes for top/bottom select-screen paint seams.
- [x] Run sixth-pass edge-band probes under Xvfb and record observations.
- [x] Add reusable X11/Xvfb control and snapshot helpers.
- [x] Prove boot/load screen to main-menu screenshot capture.
- [x] Prove main menu to Arcade character-select route.
- [x] Prove character-select route can continue into a live fight.
- [x] Record pass-7 control/snapshot route observations and evidence paths.
- [x] Add and verify explicit screen-space select header/footer text under the edge probe gate.
- [x] Record pass-8 select header/footer evidence path.

## Observation Targets

- [x] Select screen: top band and bottom band; explicit screen-space helper shows readable header/footer text over the live Arcade select view.
- [x] Select screen: route reaches Arcade character select under Xvfb using held confirmation inputs.
- [ ] Select screen: roster cells, cursors, names, title, stats overlay; convert current visual proof into a more complete mapping table.
- [x] Fight screen: route reaches a live fight with lifebars, timer, stage, and spawned characters visible.
- [ ] Versus screen: background, portraits, names, stage, top background.
- [ ] Fight screen: compare `HP*`, `FACE*`, `PWR*`, and `NAME*` pre/post labels.
- [ ] Pause/menu overlay during a fight.
- [ ] Fade and transition frames.

## Mapping Output

- [ ] Build a table of screen -> visible probes -> inferred draw phase.
- [ ] Identify which surfaces are Lua/script-driven.
- [ ] Identify which surfaces are Go/FightScreen-driven.
- [ ] Identify which motif layers are safe for reusable overlay work.
- [ ] Decide whether the target feature should hook motif, fight-screen, Lua screen code, or a new thin shared overlay seam.

## Cleanup Gate

Before merging any production branch:

- [ ] Keep the probes gated, or remove them entirely if the mapping is complete.
- [ ] Preserve the map in docs if it remains useful.
- [ ] Do not leave unconditional visual debug output in normal gameplay.
