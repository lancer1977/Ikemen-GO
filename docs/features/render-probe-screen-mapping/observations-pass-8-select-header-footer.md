# Observations Pass 8: Select Header/Footer Screen-Space Text

## Goal

Iterate on the proven private-X character-select route and test whether explicit header/footer text can be drawn during the live Arcade select view.

## Change Under Test

- Added `renderProbeScreenBlockMode(...)`, a gated Lua-facing helper for screen-space probe blocks.
- Exposed the helper from Go through `src/script.go`.
- Called the helper from the select loop in `external/script/start.lua` after the stats overlay and before `refresh()`:
  - `SELECT HEADER SCREEN-SPACE` at the top edge.
  - `SELECT FOOTER SCREEN-SPACE` near the bottom edge.
- Kept this under the `IKEMEN_RENDER_PROBES=edge-screen` category gate.

## Important Iteration Notes

1. The first helper version drew immediately from Lua. It did not survive the select-loop draw ordering because subsequent queued screenpack draws flushed by `refresh()` covered it.
2. The working version queues the screen-space probe into the top Lua draw layer via `sys.luaQueueLayerDraw(2, ...)`, so it is flushed late enough to appear over the character-select view.
3. The runnable app tree at `/home/lancer1977/code/ikemen-app` must use the matching source scripts. Copying only `start.lua` from the source checkout produced a startup mismatch (`main.f_safeGameOption` call failed). Copying the full `external/script/*.lua` set fixed the runtime mismatch.
4. The route that reliably reached true Arcade character-select used `hold:Return:0.4` on the Arcade submenu before the final `hold:z:0.4` confirm.

## Evidence Command

From `/home/lancer1977/code/Ikemen-GO`:

```bash
OUT=/home/lancer1977/code/Ikemen-GO/artifacts/visual-probes/xvfb-edge-screen-header-footer-$(date +%Y%m%d-%H%M%S)
mkdir -p "$OUT"
/home/lancer1977/code/Ikemen-GO/scripts/visual/ikemen-xvfb-control-snapshots.sh \
  --workdir /home/lancer1977/code/ikemen-app \
  --bin ./Ikemen_GO_Linux \
  --probe-mode edge-screen \
  --output-dir "$OUT" \
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

## Verified Artifacts

- `/home/lancer1977/code/Ikemen-GO/artifacts/visual-probes/xvfb-edge-screen-header-footer-20260529-130322/20260529-130345-character-select.png`
- `/tmp/ikemen-edge-screen-pass7b/20260529-130819-04-after-enter-4.png`

Visual verification of that PNG showed:

- The view is the `ARCADE` character-select screen.
- The roster perspective grid is visible in the center.
- A selected character portrait/sprite is visible on the left.
- Stats labels are visible at the lower left and lower right.
- A cyan header block is visible across the top edge with red text `SELECT HEADER SCREEN-SPACE`.
- The updated negative-anchor pass draws a cyan footer block at the actual bottom edge.

## Current Mapping Takeaway

A late queued Lua layer is a viable seam for select-screen header/footer diagnostic text. The generic edge-band blocks proved the route and draw phases, but the explicit screen-space helper is the first pass that visibly places header/footer overlays over the live character-select screen. The final helper supports edge anchoring: width/height values of `0` fill to the right/bottom edge, and negative `x`/`y` values anchor from the right/bottom screen edge.

## Follow-up Notes

- The first footer probe used `y=218`, which produced a visible block in the lower-left/mid-screen area rather than the actual bottom edge. Switching to `y=-22`, `w=0`, `h=22` anchored the footer block to the bottom edge.
- The footer block itself is now bottom-edge visible; text placement can be tuned separately if production UI needs readable footer copy instead of just a colored bar.
- Keep all of this gated; do not leave unconditional header/footer debug text in normal gameplay.
