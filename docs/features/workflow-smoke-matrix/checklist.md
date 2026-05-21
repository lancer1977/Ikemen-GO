# Checklist

## Discovery

- [x] Identify the common startup and CLI branches that should be covered
- [x] Confirm the repo does not already have a general-purpose smoke harness

## Implementation

- [x] Add `safeGameOption`, `safeFightFramesPerCount`, and `safeModifyGameOption` guards in startup Lua
- [x] Remove the nested helper scope issue that blocked later startup branches
- [x] Add `scripts/smoke/ikemen-workflow-matrix.sh`
- [x] Make the harness use a temporary work root and optional fixture content
- [x] Support the deploy-tree motif path at `data/ikemen1/system.def`
- [x] Add `scripts/smoke/stream-box-ikemen-smoke.sh` for the Windows rig deploy
- [x] Suppress blocking Windows crash UI and redirect stream-box stdout/stderr to `C:\\Apps`
- [x] Add `-noerrordialog` / `IKEMEN_SUPPRESS_ERROR_DIALOG=1` for logged,
  non-blocking engine failures
- [x] Install a `stream-box` toolbox script pack for check/status/tail/launch
- [x] Add a compact stream-box diagnostic bundle command
- [x] Add a launcher-driven candidate-quality validation command
- [x] Harden the `C:\\Apps\\mugen` repoint command so it will not delete a
  real directory by default
- [x] Add stream-box character path normalization command for folder/`.def`
  names with dry-run and apply modes
- [x] Add character self-play sweep and compatibility JSON output in smoke matrix
- [x] Pass AI levels into character self-play sweep launches
- [ ] Add stream-box native character sweep reporting
- [ ] Add a Windows/MSYS2 execution wrapper if we want native Windows coverage
- [ ] Add CI wiring or a self-hosted runner hook

## Validation

- [ ] Run the startup matrix against the local Linux binary (blocked by current startup menu crash)
- [ ] Run the matrix with `/mnt/shared/Emu/ikemen` for quick-vs coverage
- [ ] Run `make smoke-stream-box`-equivalent checks against the stream-box rig
- [ ] Confirm the harness catches the known Lua startup panic class
- [x] Add a strict stream-box assertion for quick-vs result-file output
- [x] Replace repeated ad hoc launch commands with named toolbox scripts
- [x] Prefer the launcher toolbox for candidate-library validation over repo-side analysis
- [x] Capture a diagnostic bundle for the current quick-vs/menu failures
- [ ] Rebuild and run one known-failing case to confirm the error lands in logs
  without showing a modal dialog
- [ ] Run `scripts/smoke/ikemen-workflow-matrix.sh --char-sweep` against a
  representative fixture and review `ikemen-char-sweep-results.json`

## Follow-up

- [ ] Decide whether a package-install smoke should live in this same matrix or in a separate release-packaging feature
- [ ] Decide whether to add a failing-log snapshot fixture for known regressions
- [x] Implement the Windows result-file contract in [Windows Result Transport](./windows-result-transport.md)
- [ ] Keep named pipe / local callback transport as a future live-event goal
