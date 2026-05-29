# Render Probe Screen Mapping - Observation Pass 1

## Context

Source probe build was copied into the local app runtime at `/home/lancer1977/code/ikemen-app` and run with:

```bash
IKEMEN_RENDER_PROBES=1 ./Ikemen_GO_Linux
```

Observation source: user visual report from Demo Mode fight screen.

## Visible Probe Layout

### Global `System.draw` column

On the very left edge, the global probes appear as a vertical column:

- `after fight 0`
- `after motif 0`
- `after fight 1`
- `after motif 1`
- `after fight 2`
- `after motif 2`
- `after motif 3/fade`

Notes:

- The column is on the left edge and is partly cramped/cut off.
- There are colorful bars in the upper-left corner near/over this area, likely existing engine debug visualizations.

Inference:

- The main battle `System.draw` sequence is active.
- Fight and motif layer calls are interleaved as expected.
- Probe positions need to be shifted inward for the next pass.

### FightScreen probes

Observed:

- `F0 fight entry`
- `F1 fight entry`
- `F2 fight entry`
- `F2 Round` appears at the bottom of a list with `F0 fight entry`.
- `F1 Fight` appears over the P1 power bar.
- `F0 fight entry` appears behind the bar.

Inference:

- `FightScreen.draw` is called for layers 0, 1, and 2 during battle.
- At least one FightScreen probe is behind/under part of the HUD bar stack.
- Some FightScreen-layer text can land over power-meter space, so that area is not safe for readable probes without moving them.

### Motif probes

Observed:

- `M1 motif entry` appears around the P2 power meter area.
- `M0 motif entry` appears above the health bar and in front.
- Nothing notable appeared on the right side despite earlier expectation of right-side labels.

Inference:

- Motif layer probes are visible in battle.
- Motif 0 can render visibly in front of the healthbar region, depending on placement/layer interaction.
- Current right-side/mid-screen probe placement is not where the user sees it or is being occluded/off useful view; next pass should use clearer anchor positions.

### Existing colored bars / portrait overlap

Observed:

- A green bar appears over the P1 portrait.
- Yellow and orange bars above it appear behind the P1 picture.

Inference:

- The P1 portrait is not uniformly above or below the whole colored-bar/debug stack.
- The existing bar stack likely spans multiple draw depths or overlaps with HUD/portrait draw order in a non-trivial way.
- This is strong evidence that lifebar/portrait internals need more granular probes around the healthbar, portrait, and powerbar draw calls.

## Select Screen Note

User did not think any select-screen probes were visible in this pass.

Interpretation:

- This pass was observed in Demo Mode / fight runtime, so the normal interactive select screen may not have been visited.
- If a select screen was briefly shown, the current `S*` probes may be hidden by the local app's older Lua flow, screen timing, or placement.
- The lack of select labels should not yet be treated as proof that select drawing is uninstrumented; it is an unresolved observation target.

Next select-specific check:

- Launch a mode that forces the interactive select screen instead of Demo Mode.
- Watch specifically for `S0` through `S5` before entering a fight.
- If still absent, add probes earlier in the local app's select flow and/or add a title/menu probe to prove Lua `renderProbe` is callable outside battle.

## First-Pass Conclusions

1. Battle HUD rendering is layered and active through `FightScreen.draw(layerno)` for multiple layers.
2. Motif rendering is also active during battle through `Motif.draw(layerno)`.
3. `System.draw` interleaves fight and motif layers, then ends with motif layer 3/fade.
4. Healthbar/powerbar/portrait sub-order is not resolved yet; current probes are too coarse and overlap live HUD/debug visuals.
5. Select-screen instrumentation is not confirmed yet; first pass was fight-focused and likely skipped the interactive select path.
6. The next pass should add granular probes around lifebar, portrait, name, and powerbar calls, and should move labels away from the far-left edge and lower edge.

## Next Probe Pass

Recommended changes:

- Shift the global `G*` column rightward so full labels are readable.
- Move `F* fight entry` labels away from powerbar and healthbar regions.
- Add dedicated probes immediately before/after these FightScreen sections:
  - life bar draw
  - face/portrait draw
  - name draw
  - power bar draw
  - round/action/combo draw
- Use compact labels like:
  - `HP pre`
  - `HP post`
  - `FACE pre`
  - `FACE post`
  - `PWR pre`
  - `PWR post`
- Use separate fixed columns for P1 and P2 to avoid ambiguity.

## Design Implication

For battle-time roster/healthbar-style UI, piggybacking on `FightScreen.draw(layerno)` is still the strongest candidate because the evidence shows battle bars, portraits, power meters, combo/action text, and round text are owned by the fight-screen path.

For non-battle screenpack UI, continue treating Lua/motif as the likely owner.
