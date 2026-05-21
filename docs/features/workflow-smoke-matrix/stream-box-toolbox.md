# Stream Box Toolbox

## Summary

This is the small script pack I want available on `stream-box` so we stop
passing long command lines by hand every time we want to check or launch the
Ikemen runtime.

The idea is to keep the machine-side commands boring:

- `check`
- `launch-quickvs`
- `launch-menu`
- `status`
- `tail`
- `stop`
- `repoint`
- `bundle`
- `normalize-chars`
- `candidate-quality`

The repo-local helper should install those scripts into a single tools folder
and then call them by name.

## Current State

- [x] `stream-box` is the canonical Windows rig
- [x] `C:\\Apps\\mugen` is the live Ikemen root on that box
- [x] The repo already has a reusable `scripts/smoke/stream-box-ikemen-smoke.sh`
- [x] Install a compact PowerShell toolbox on the box itself
- [x] Replace repeated ad hoc `ssh` / `schtasks` / `scp` command lines with a
  named script command

## Proposed Scripts

- `ikemen-check.ps1`
  - Verifies the live deploy tree, the executable, and the core runtime assets.
- `ikemen-launch-quickvs.ps1`
  - Launches the fixed quick-vs smoke with the current `C:\\Apps\\mugen` tree.
  - Returns a failing exit code when redirected stderr contains panic/fatal
    runtime text.
  - Treats a still-running process as healthy when no panic/fatal text is found.
- `ikemen-launch-menu.ps1`
  - Starts the game without fight args so menu/bootstrap behavior can be compared
    to the quick-vs lane.
  - Returns a failing exit code when redirected stderr contains panic/fatal
    runtime text.
  - Treats a still-running process as healthy when no panic/fatal text is found.
- `ikemen-status.ps1`
  - Prints process status plus the last known exit state for the smoke lane.
- `ikemen-tail.ps1`
  - Tails `Ikemen.log` and the smoke stderr/stdout logs.
- `ikemen-stop.ps1`
  - Stops any stale `Ikemen_GO` process before a new launch.
- `ikemen-repoint.ps1`
  - Replaces the `C:\\Apps\\mugen` junction with a new build root.
  - Refuses to delete a non-link directory unless `-ForceDirectory` is used.
- `ikemen-bundle.ps1`
  - Captures the current log files, root link state, process state, display
    details, and recent bridge result files into `C:\\Apps\\ikemen-diagnostics`.
- `ikemen-normalize-chars.ps1`
  - Creates a collision-safe normalization plan for top-level `chars` folders
    and all `.def` filenames.
  - Supports dry-run and apply modes.
  - Writes a TSV mapping report under `C:\\Apps\\ikemen-char-normalization.*.tsv`.
- `ikemen-candidate-quality.ps1`
  - Runs the full candidate-quality validation loop from the launcher side.
  - Inventories `T:\\chars-candidates`, normalizes into a timestamped run,
    writes triage, duplicate review, and promotion-gate reports, and dry-runs
    promotion for the first ready candidate when one exists.
  - Accepts a source-root override so the workflow can point at the real
    character tree when `T:` is not mounted, and resolves the live promotion
    target on the box.
  - Keeps the validation flow on the box instead of relying on repo-side
    analysis.

## Follow-Up Scripts

- `ikemen-result-check.ps1`
  - Reads a future Windows result file once the result-file contract lands.

## Validation

- [x] Install the toolbox on `stream-box`
- [x] Run `ikemen-check.ps1`
- [x] Run `ikemen-launch-quickvs.ps1`
- [x] Run `ikemen-launch-menu.ps1`
- [x] Run `ikemen-status.ps1`
- [ ] Run `ikemen-tail.ps1`
- [x] Run `ikemen-repoint.ps1`
- [x] Run `ikemen-bundle.ps1`
- [x] Run `ikemen-normalize-chars.ps1` dry-run
- [x] Run `ikemen-normalize-chars.ps1 -Apply`
- [ ] Run `ikemen-candidate-quality.ps1`
- [ ] Confirm stale-process cleanup works before a second launch

## Notes

- The toolbox should keep the awkward launch arguments inside the script, not in
  the operator’s shell history.
- The repo-local smoke script can remain the source of truth for the exact args,
  but the box-side scripts should be the daily operator entrypoints for
  validation, especially when the workflow needs to prove real launcher
  behavior rather than just repo-side analysis.
- After any launch failure, collect a bundle before changing roots or rebuilding:
  `bash scripts/stream-box/ikemen-box.sh bundle failure-label`.
