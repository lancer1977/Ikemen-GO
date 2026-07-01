# Ikemen GO Launch Contract

## Summary

This feature documents the launch surface that ChannelCheevos and related
tools can rely on when starting Ikemen GO.

The current shape is intentionally split:

- ChannelCheevos owns identity, roster selection, viewer ownership, and launch
  policy.
- Ikemen GO owns runtime execution, stage/motif selection, and final match
  result emission.
- stdout remains a fast local/debug lane, while a per-match result file is the
  canonical bridge integration path.
- The final result payload and live snapshot both carry end-state signals:
  `fightEnded` in the result contract and `matchOver` in the live contract.
- The existing JSON output shape is documented as
  [Result Contract V1](./result-contract-v1.md) and mirrored by the
  `contracts/Ikemen.Go.Contracts` C# model library.
- Live player-visible overlays use a separate file snapshot contract at
  `-overlayfile` (default `save/live_overlays.json`) so sidecars can inject
  text, emoji, and PNG-backed image sprites during a fight without rewriting
  the main result snapshot.

## Engine Seams

The engine exposes a small set of named Lua hooks that are stable attachment
points for launch-adjacent behavior:

- `game.challenger_init` / `game.challenger`
- `game.continue_init` / `game.continue`
- `game.hiscore_init` / `game.hiscore`
- `game.victory_init` / `game.victory`
- `game.result_init` / `game.result`

These hooks are the places to attach menu, transition, and post-match behavior
without coupling new behavior to the battle loop itself.

The file-backed seams are the other stable attachment points:

- `-combateventsfile` for combat telemetry
- `-matcheventsfile` for round/match timeline events
- `-commandinboxfile` and `-commandresultsfile` for gameplay command relay
- `-livedatafile` and `-resultfile` for active/live and terminal snapshots
- `-overlayfile` for runtime-published player-visible overlays

## Current State

- [x] Document the existing CLI launch surface
- [x] Preserve the current quick-vs and debug stdout path
- [x] Add canonical `-resultfile` output for bridge consumption
- [x] Keep `-jsonlog` and `-jsonstdout` as developer-friendly lanes
- [x] Define the existing result/live JSON shape as a V1 contract
- [x] Add a file-backed live overlay transport for in-fight text/image injects
- [ ] Add a structured JSON launch payload if the flag surface becomes too wide
- [ ] Add live round-event transport for match-by-match streaming

## Supported Customization Matrix

### Bridge-only now

These concerns should stay in ChannelCheevos or another orchestrator.

- player identity
- viewer ownership
- roster ordering
- palette choice policy
- stage choice policy
- match history persistence
- run IDs and per-run file naming

### Engine-supported now

These are already expressible through Ikemen GO’s CLI and runtime output.

- `-p<n>` player selection
- `-p<n>.ai`
- `-p<n>.color` / `-p<n>.pal`
- `-p<n>.life`
- `-p<n>.lifeMax`
- `-p<n>.power`
- `-p<n>.dizzyPoints`
- `-p<n>.guardPoints`
- `-tmode1` / `-tmode2`
- `-time`
- `-rounds`
- `-draws`
- `-s`
- `-windowed`
- `-width` / `-height`
- `-setvolume`
- `-nojoy`
- `-nomusic`
- `-nosound`
- `-debug`
- `-debugstartup`
- `-togglelifebars`
- `-maxpowermode`
- `-jsonlog`
- `-jsonstdout`
- `-nojsonlog`
- `-nojsonstdout`
- `-resultfile`
- `-livedatafile`
- `-overlayfile`

### Engine change later

These are better treated as future work unless the current CLI surface stops
being enough.

- first-class JSON launch payloads
- named pipe result transport
- localhost callback transport
- live round-by-round event streaming
- automatic exit modes for automation-only launches

### Out of scope

These should stay out of Ikemen GO itself.

- ChannelCheevos UI state
- viewer authentication policy
- long-term match-history policy
- bridge-side strategy selection
- stream overlay presentation

## Notes

- The existing JSON stdout path is still useful for local smoke and debugging.
- The file-based result path should be unique per match so the bridge can treat
  it as the source of truth.
- The bridge can keep using the same flag surface even if the roster policy or
  result transport changes later.
- Consumer apps should treat the file result as canonical and can use
  `fightEnded`/`matchOver` to distinguish an active fight from a completed one
  without parsing message text.
- Consumer apps that use C# can reference `Ikemen.Go.Contracts.V1` for strong
  models of the current result and live snapshot payloads.

## Launch Mapping

ChannelCheevos request fields map to the following launch contract:

- Match participant identity:
  - `p1.def`, `p2.def` → `-p1 <path>`, `-p2 <path>`
- AI settings:
  - `p1.ai`, `p2.ai` → `-p<n>.ai <level>`
- Palette:
  - `p1.palette`, `p2.palette` → `-p<n>.color <palette_index>`
- Life/power overrides:
  - `p1.life`, `p2.life` → `-p<n>.life <value>`
  - `p1.power`, `p2.power` → `-p<n>.power <value>`
- Match pacing:
  - `match.roundTime` → `-time <seconds>`
  - `match.rounds` → `-rounds <count>`
- Team mode:
  - `match.teamMode[1]`, `match.teamMode[2]` → `-tmode1 <mode>`, `-tmode2 <mode>`
- Stage/motif:
  - `stage.path` → `-s <path>`
  - `system.motif` overrides remain in `save/config.ini` as `Motif = ...`
- Launcher strategy:
  - stream vs local result handling via `-resultfile <run-scoped path>`
  - optional debug mirror lane through `-jsonstdout`

Stable run metadata convention:

- Generate one unique result file path per launch
  (`ikemen-workflow-smoke.<runId>.json`, `.../ikemen-results/<runId>.json`, etc.)
- Keep both `-jsonstdout` and `-resultfile` enabled for local/dev validation,
  but treat `-resultfile` as canonical for bridge integration.
- Use `-livedatafile` when a caller needs pollable in-flight match state such as
  health, current round score, or `matchOver` during the fight. The live file
  is also refreshed once the fight ends so the last snapshot records that the
  match is over and the result has been finalized.
- Use `-overlayfile` when a sidecar wants the engine to render player-visible
  overlay requests. The file is treated as a snapshot of current overlays and
  can include text, emoji, and PNG-backed image effects.
