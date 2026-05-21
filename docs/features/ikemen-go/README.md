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

## Current State

- [x] Document the existing CLI launch surface
- [x] Preserve the current quick-vs and debug stdout path
- [x] Add canonical `-resultfile` output for bridge consumption
- [x] Keep `-jsonlog` and `-jsonstdout` as developer-friendly lanes
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
