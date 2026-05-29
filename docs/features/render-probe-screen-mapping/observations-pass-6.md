# Render Probe Screen Mapping - Observation Pass 6

## Context

Pass 6 explored whether the character select screen can paint dedicated top and bottom edge bands instead of only the portrait and roster/cell picker areas.

The first local desktop launch on `DISPLAY=:0` failed because the X server was saturated:

```text
Maximum number of clients reached
gtk initialisation failed; presumably no X server is available
```

The pass was then run headlessly with Xvfb and screenshots captured from `/tmp/ikemen-edge-pass6/`.

## Run Harness

The runnable app directory was `/home/lancer1977/code/ikemen-app`.

The source pass mode is:

```bash
IKEMEN_RENDER_PROBES=edge ./Ikemen_GO_Linux
```

Because the existing runnable binary predates the newly-added `edge` category behavior, the exploratory runtime also used `IKEMEN_RENDER_PROBES=all` and `IKEMEN_RENDER_PROBES=swarm-block` to verify whether same-frame block drawing was visible through the already-working probe category.

## Screens Reached

Automated key input through Xvfb reached:

1. title/logo screen
2. main menu
3. arcade mode submenu
4. character select screen
5. selected-character state after additional confirm/right input

Captured select screenshots:

```text
/tmp/ikemen-edge-pass6/03-after-enter-3.png
/tmp/ikemen-edge-pass6/04-after-enter-4.png
/tmp/ikemen-edge-pass6/05-after-enter-5.png
/tmp/ikemen-edge-pass6/06-after-right-right.png
```

## Observations

### Confirmed visible again

The previously-known select/cell probes still render:

- Green block near the left selected-character/portrait area.
- Green block at the top/left portion of the roster/cell picker.
- Red/pink `E0`-style marker near the lower-right part of the roster/cell picker.

This reconfirms that the visible roster/portrait picker is reachable through the select/batch draw path.

### Edge probes were not visible

Dedicated top/bottom probes were not visible in the screenshots:

- `T0` / `B0`
- `T1` / `B1`
- `T2` / `B2`
- `T3` / `B3`
- `T4` / `B4`
- `T5` / `B5`

The first attempt used larger/offscreen-friendly labels and positions; the second attempt adjusted to the 320x240 logical coordinate space. Neither produced visible top/bottom edge blocks.

### Final seam test also did not show

A runtime-only experiment temporarily drew large final-seam `swarm-block` probes after `bgDraw(..., 1)` at multiple y positions. Those final-seam probes also did not appear, even though earlier select/batch probes remained visible.

That means the issue is not only the new `edge` category. Late Lua-side drawing after the layer-1 background path is not currently producing visible probe blocks in this select-screen state.

## Evidence-Level Conclusion

Current evidence says:

- The visible character/roster picker is still owned by the select draw-list / `batchDraw` path.
- The top and bottom header/footer bands are not proven paintable from the same Lua probe calls used for roster/portrait evidence.
- Late Lua paint after `bgDraw(motif.selectbgdef.BGDef, 1)` is suspect: either it is clipped, not presented in this screenpack state, affected by render state, or otherwise not a reliable header/footer seam.

## Updated Hypothesis

For top and bottom select-screen overlays, the next useful probe is not more Lua `renderProbeBlockMode(...)` at arbitrary coordinates.

The next pass should test one of these seams:

1. A Go-side screen-space overlay helper that explicitly resets render state and paints in `sys.scrrect` coordinates after Lua select drawing but before `refresh()`.
2. A motif/screenpack layer entry rather than ad hoc Lua blocks.
3. A targeted probe inside/around the existing visible `batchDraw` queued-layer callback, but positioned in logical top/bottom screen-space to see whether the projected batch path can own header/footer art.

## Practical Next Step

Build a `screen-space` probe helper separate from the current generic `renderProbeBlockMode` path. It should draw with explicit orthographic screen coordinates and no inherited select/screenpack state assumptions.

Then expose it to Lua as something like:

```lua
renderProbeScreenBlock('edge-screen', 'TOP', 0, 0, 320, 24, 255, 0, 255)
renderProbeScreenBlock('edge-screen', 'BOTTOM', 0, 216, 320, 24, 255, 0, 255)
```

If that works, the production overlay seam should be a dedicated screen-space select overlay hook, not the existing portrait/cell batch path.
