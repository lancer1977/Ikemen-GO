# Workflow Smoke Matrix

## Summary

This feature adds a repeatable startup and workflow smoke harness for Ikemen GO.
The goal is to validate the real executable against the launch branches that
have historically been fragile: plain startup, command-line option handling,
option mutation, debug toggles, window sizing, and optional gameplay startup
when fixture content is available.

## Current State

- [x] Define a workflow matrix for the common startup branches
- [x] Add a reusable smoke runner script
- [x] Add safe startup guards for fragile Lua option calls
- [ ] Wire the matrix into CI or a self-hosted runner
- [ ] Add a Windows/MSYS2 execution path for the matrix
- [x] Add a content-backed quick-vs fixture profile for fuller gameplay smoke
- [x] Add a stream-box deploy guard for the Windows rig runtime tree
- [x] Suppress blocking Windows crash UI during stream-box smoke runs
- [x] Add engine-level no-error-dialog mode for programmatic debug runs
- [x] Add per-character self-play compatibility sweep with JSON report output
- [x] Add stream-box character/def normalization command to remove
  space-punctuated path fragility
- [x] Add a small script toolbox on `stream-box` so routine checks do not need
  long command lines

## Coverage

The smoke runner currently exercises:

- plain startup
- `-windowed`
- `-nosound`
- `-jsonstdout`
- `-ailevel`
- `-speedtest`
- `-framerate`
- `-debug`
- `-togglelifebars`
- `-maxpowermode`
- `-setvolume`
- `-width` / `-height`

If `IKEMEN_WORKFLOW_FIXTURE_ROOT` points at a full content tree and both
`IKEMEN_WORKFLOW_PLAYER_DEF` is valid, the harness runs:

- quick-vs using the configured fixture players and default stage resolution
- optional stage-override quick-vs (same fixture with explicit `-s`, when
  `IKEMEN_WORKFLOW_STAGE_DEF` exists)
- quick-vs result-file checks (default and AI/palette variants)
- character self-play sweep via `--char-sweep` when enabled (default 30s and
  3-rounds)
- stream-box compatibility run override via `--time`, `--rounds`, `--p1-def`,
  and `--p2-def` for long-run validation scenarios

Quick-vs defaults in both local and stream-box smoke paths use AI level 8 for both
players by default. You can override:

- Local matrix: set `IKEMEN_WORKFLOW_P1_AI_LEVEL` and
  `IKEMEN_WORKFLOW_P2_AI_LEVEL` (e.g. `0` for no AI actuation).
- Stream-box smoke: pass `--ai-level` or set `IKEMEN_STREAM_BOX_AI_LEVEL`.

For the character sweep, each record includes:

- `character`: folder name used for discovery
- `path`: resolved character definition path used for launch
- `status`: `ok`, `error`, `timeout`, or `missing-result`
- `compatible`: boolean compatibility flag
- `elapsedSeconds`: elapsed launch and match duration in seconds
- `exitCode`: Ikemen process exit code
- `resultExists`: whether `-resultfile` was written

The default shared deploy fixture is:

- `/mnt/shared/Emu/ikemen`

## Notes

- The harness prefers `xvfb-run` when no display is available.
- The test tree is isolated in a temporary work directory so the repo checkout
  does not get dirtied by startup smoke runs.
- The runner fails on the first Lua panic or runtime exception string.
- `make smoke-stream-box` validates the deployed `C:\\mugen` tree, checks that
  `external/script/main.lua` contains the expected startup guards, launches a
  visible quick-vs smoke through Task Scheduler, and scans `Ikemen.log` for the
  same crash patterns.
- The stream-box smoke sets `HKCU\\Software\\Microsoft\\Windows\\Windows Error
  Reporting\\DontShowUI=1` for the test user and captures process stderr/stdout
  under `C:\\Apps\\ikemen-workflow-smoke.*.log`, so launch failures are reported
  back through the smoke command instead of blocking on a rig popup.
- The engine also supports `-noerrordialog` and
  `IKEMEN_SUPPRESS_ERROR_DIALOG=1`. Smoke/debug launchers should use both when
  they need to log a failure and proceed to the next test without manually
  clicking `OK`.
- The stream-box smoke run uses character-only quick-vs by default, with optional
  stage-override coverage when a stage definition is explicitly available. The
  current bridge path should prefer the per-match result file, with stdout left
  in place for local debugging and smoke triage.
- For candidate-library validation, prefer the stream-box launcher toolbox over
  repo-side code analysis so the workflow proves the real box-side behavior.
- Windows remains the required stream-rig target because HDMI capture and
  MilkDrop are Windows-dependent in the current setup. The durable automation
  path is tracked in
  [Windows Result Transport](./windows-result-transport.md).
- The operator script inventory is tracked in
  [Stream Box Toolbox](./stream-box-toolbox.md).
- The current per-character failure inventory is tracked in
  [Character Error Analysis](./character-error-analysis.md).
