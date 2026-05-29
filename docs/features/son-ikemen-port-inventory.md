# Son's IKEMEN Work Port Inventory

## DreadBreadcrumb

This inventory captures the work that appears to come from the older/Lancero IKEMEN customization line and classifies what should be ported into the cleaner modern IKEMEN project now that the render-probe work has proven reliable select-screen drawing seams.

## Sources inspected

- Modern engine/source workspace: `/home/lancer1977/code/Ikemen-GO`
- Current runnable/probe install: `/home/lancer1977/code/ikemen-app`
- Older local IKEMEN install: `/home/lancer1977/Applications/ikemen`
- Wrapper / coordination repo: `/home/lancer1977/code/ikemen-plus`
- Prior feature docs:
  - `/home/lancer1977/code/ikemen-plus/docs/features/select-screen-stats-overlay/README.md`
  - `/home/lancer1977/code/ikemen-plus/docs/features/select-screen-stats-overlay/architecture.md`
  - `/home/lancer1977/code/ikemen-plus/docs/features/select-screen-stats-overlay/durability-notes.md`
  - `/home/lancer1977/code/ikemen-plus/docs/character-stats.md`
  - `/home/lancer1977/code/ikemen-plus/docs/lancero-consolidation.md`

## Current conclusion

The work should be treated as a three-layer product, not just a screen skin:

1. **Engine/runtime contract** — stats persistence, result JSON export, ranking/hiscore logic, command-line flags.
2. **Screen/UI layer** — select-screen W/L/tier display and future header/footer or fight-screen presentation.
3. **Distribution/publisher layer** — local Lancero install, protected roster/content folders, stable/dev binary selection, publish/rollback scripts.

The modern render-probe pass changed the risk profile: the select screen is now known to be drawable from a late queued Lua layer, so the stats display can be ported as a deliberate modern UI feature instead of preserved as a fragile old-build hack.

## Port classification

| Area | Evidence / files | Classification | Forward action |
| --- | --- | --- | --- |
| `save/stats.json` per-character records | `ikemen-plus/docs/lancero-consolidation.md`, `durability-notes.md`; current app has `/home/lancer1977/code/ikemen-app/save/stats.json` | **Keep as durable contract** | Define and version the stats schema. Do not tie it to a particular install folder. |
| JSON bridge | `getGameStatsJson()`, `setGameStatsJson()`, `-jsonlog`, `-resultfile`, `-jsonstdout` references in modern `src/script.go`, `src/main.go`, and `external/script/main.lua` | **Keep / stabilize** | Treat as the canonical import/export path for analysis and outside tooling. |
| Ranking / hiscore runtime support | `src/hiscore_rank.go`, `src/motif.go` | **Keep in Go/runtime** | Leave durable calculation/storage in engine code; UI should consume snapshots/results. |
| Select-screen stats overlay | `external/script/start.lua`; `ikemen-plus/docs/features/select-screen-stats-overlay/*` | **Port concept into modern seam** | Rebuild as thin Lua/UI rendering over modern `start.f_selectScreen()` draw order. Use private-X visual proofs. |
| Record/tier text fallback behavior | `start.f_getRecordText()` panic notes; docs require `0/0/U Tier` fallback | **Keep behavior, redesign defensively** | Missing stats/roster/JSON must never panic the select screen. |
| Placement grid / coordinate docs | `docs/character-stats.md` mentions `-placementgrid` / `-placement-grid` and 10x10 grid | **Keep as development tool, but verify current implementation** | Either port the flag/helper or replace it with the newer render-probe coordinate markers. |
| Roster/select.def management | `data/select.def`, possible `select-defs/active.def`; parser scripts in `ikemen-plus/scripts/select_entry_parser.py` | **Keep as tooling / content pipeline** | Do not merge into engine core. Keep parser and tests as external content-management tools. |
| Publisher scripts | `ikemen-plus/scripts/publish-to-lancero-test.sh`, publish/rollback/sync scripts | **Keep outside engine** | Use for local install lifecycle only. Must preserve `chars/`, `stages/`, `save/`, `sound/`, `video/`, `.git/`, and `data/select.def`. |
| Stable/dev local installs | Historical docs reference `~/apps/ikemen-stable` and `~/apps/ikemen-dev`; current machine has `/home/lancer1977/Applications/ikemen` and `/home/lancer1977/code/ikemen-app` | **Normalize docs before relying on paths** | Update wrapper docs/scripts if the canonical current paths are now `Applications/ikemen` and `code/ikemen-app`. |
| In-match W/L/tier display | Consolidation docs say preserve it; feature docs defer fight-screen support | **Defer / design separately** | Add only after a stable match-start snapshot hook exists. Do not poll stats every frame. |
| Old folder-specific hacks/backups | Backup files, direct mutable install edits | **Discard or archive** | Preserve as evidence only; do not port brittle path assumptions into modern IKEMEN. |

## File-change stocktake

Key file families differ substantially between the older local install and the current probe install:

- `external/script/start.lua`: old 4421 lines vs current 4155 lines; large semantic drift.
- `external/script/main.lua`: old 4161 lines vs current 4087 lines; large semantic drift.
- `data/select.def`: old 393 lines vs current 354 lines; content/layout drift.
- `save/stats.json`: old install contains an effectively empty `{}` file; current app has a populated 14-line JSON sample.

Because the Lua files have heavy upstream drift, the port should **not** be a blind file copy. It should be a feature extraction:

1. Identify the stable data contract.
2. Recreate the display behavior in current `start.f_selectScreen()`.
3. Use current render-probe route to capture visual proof.
4. Keep wrapper/content tooling outside the engine repo.

## Modern implementation target

The clean forward project should expose these small, stable units:

### 1. Stats data contract

- Source: `save/stats.json` or an explicit `-stats` path.
- Required behavior:
  - load safely if missing/empty/partial;
  - provide per-character W/L/tier lookup;
  - write updated results after matches;
  - optionally export/import through JSON bridge.

### 2. Select-screen stats widget

- Draw location: modern select-screen Lua path, likely `external/script/start.lua` near `start.f_selectScreen()`.
- Rendering seam: late enough to remain visible. The pass-8 render-probe proof showed a late queued Lua layer can draw over the live character-select view.
- Display:
  - compact black-backed `stats` label;
  - `Win - Loss - Tier` line;
  - mirrored P1/P2 placement;
  - safe fallback `0 - 0 - U` or `0/0/U Tier`.

### 3. Optional dev placement helper

- Either port `-placementgrid` or replace it with the newer render-probe coordinate/block helper.
- Keep it gated by a flag/env var.
- Use screenshots as evidence for placement changes.

### 4. Future fight-screen hook

- Add later, not in the first port.
- Preferred shape: one match-start snapshot emitted after roster/match setup, then render cached values.
- Avoid live polling `save/stats.json` during combat.

## Verification path

Use the proven private-X workflow:

1. Build modern IKEMEN.
2. Copy/publish into `/home/lancer1977/code/ikemen-app` or the chosen dev runtime.
3. Launch with private Xvfb wrapper.
4. Route to Arcade select using the verified hold cadence.
5. Capture screenshots showing the stats widget and any header/footer markers.
6. Confirm missing/empty stats files do not panic.
7. Confirm a match result updates or exports through the JSON bridge.

## Immediate next port slices

1. **Schema slice** — document the current and desired `save/stats.json` shape with example records.
2. **Lookup slice** — add/verify a safe per-character lookup helper with missing-data fallback.
3. **UI slice** — render a modern select-screen stats box using the proven late Lua draw seam.
4. **Visual proof slice** — capture private-X screenshots for P1/P2 boxes.
5. **Publisher cleanup slice** — decide whether `ikemen-plus` paths should point at `/home/lancer1977/code/ikemen-app` instead of stale `~/apps/ikemen-dev` references.

## Notes / risks

- The wrapper repo currently says `gamerepo0` is a symlink to the engine repo, but in this environment `readlink -f gamerepo0` did not resolve during inspection. Treat that as a wrapper-doc drift item before relying on wrapper publish commands.
- Some `ikemen-plus` tracked files include `__pycache__` artifacts. Those are not part of the durable port and should be ignored/removed in a cleanup slice.
- Existing modern `go test ./src` has unrelated `reisen`/FFmpeg build issues; use `GOEXPERIMENT=arenas ./build/build.sh Linux` plus private-X visual proof as the current practical validation path.
