# Render Probe Screen Mapping - Observation Pass 5

## Context

Pass 5 ran the large color-block probe mode from `/home/lancer1977/code/ikemen-app`:

```bash
IKEMEN_RENDER_PROBES=swarm-block ./Ikemen_GO_Linux
```

## User Observation

On the character select screen:

- A green square is visible over the left character.
- There may be a red/pink letter/box near the bottom-right of the character select picker.
- The user thought it might read `10` or `E0`, but could not confirm without glasses.
- The red/pink box looked new compared with the previous pass.

## Interpretation

The color-block pass confirms that non-text probes are visible on the select screen.

The green square over the left character most likely corresponds to one of the roster/cell batch markers:

- `A` bright green: before `batchDraw(staticDrawList)`
- `B` dark green: after `batchDraw(staticDrawList)`
- or `E<layer>` cyan/green: Go-side queued `batchDraw` layer callback before projected animation draw

Because it appears over the left character/roster picker area, the character select picker is very likely owned by the projected roster/cell batch path.

The possible red/pink box near the bottom-right of the picker most likely corresponds to one of the late/final markers:

- `Z` magenta: final marker before `refresh()`
- possibly another block being transformed/covered near the picker area

## Evidence-Level Conclusion

The 3D-ish select screen is reachable by the instrumented render path. The visible picker/character grid is not a separate opaque renderer; it is at least overlapping with the Lua select loop and/or Go `batchDraw` queued layer path.

Best current ownership map:

- Character select picker / inner roster box: `start.updateDrawList()` -> `batchDraw(staticDrawList)` -> Go queued layer animation draw.
- Selected character / left character visual area: likely overlaps with the same projected/batch draw path or the selected portrait path, but the green block points most strongly at the batch/cell path.
- Late/top marker near bottom-right: likely final select-loop/top/refresh path, if the visible red/pink block is `Z`.

## Practical Next Step

If this evidence is enough, stop probing and use this ownership map for implementation decisions:

- Roster/grid/select picker effects should hook into the select draw-list / batchDraw path.
- Battle health/power/portrait effects should hook into FightScreen.
- Global top overlays should use late select/motif/top-layer points.

If exact per-block identity is required later, rerun with glasses or replace letters with very large single-color full-screen quadrant markers.
