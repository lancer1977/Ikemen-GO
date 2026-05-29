# Render Probe Screen Mapping - Observation Pass 4

## Context

Pass 4 ran the swarm-decorated probe build from `/home/lancer1977/code/ikemen-app`:

```bash
IKEMEN_RENDER_PROBES=swarm ./Ikemen_GO_Linux
```

## User Observation

On the select screen:

- A greenish bar is visible across roughly the first 3.3 character cells.
- The user could not make out words/codes.
- The bar appears to sit on an inner select box / roster region.
- User suspects the selected character may have a separate area.

## Interpretation

This is the first positive evidence that select-screen probe drawing is reaching the visible select surface.

Because words are not readable but a greenish probe rectangle is visible, the most likely causes are:

1. The probe background rectangle is visible, but the debug-font label is too small, clipped, recolored, or covered.
2. The visible green bar corresponds to one of the green swarm markers around the roster/cell draw path, most likely:
   - `SW30 batchDraw pre`
   - `SW31 batchDraw post`
   - or `SWB batch L<layer>` from the Go `batchDraw` queued layer callback.
3. Since it appears across the first several roster cells, it is probably aligned with the perspective roster/cell draw surface rather than the big selected-character portrait area.

## Evidence-Level Conclusion

The 3D-ish select screen is not fully detached from the instrumented path. At least one probe rectangle from the select/batch swarm pass is visible on the select surface.

The lack of readable words means the next and final diagnostic should avoid text-dependent labels. Use large colored/numbered blocks or oversized short codes (`A`, `B`, `C`, etc.) placed at distinct screen regions:

- roster/cell strip
- selected character portrait area
- face2/angled portrait area
- cursor/menu area
- top/background/final overlay

## Likely Ownership So Far

- Roster/cell surface: Lua `start.updateDrawList()` -> Go `batchDraw()` queued layer path, using perspective projection settings from the screenpack.
- Selected character portrait area: likely `start.f_drawPortraits()` / `main.f_animPosDraw()` path, separate from the cell batch list.
- Final/top layer: select `bgDraw(..., 1)`, stats overlay, hook, and refresh path still need text-free confirmation if precise top ordering matters.
