# Checklist

## Discovery

- [x] Map the current Ikemen GO launch CLI surface
- [x] Separate bridge-owned policy from engine-owned execution
- [x] Confirm the repo already emits match stats as JSON

## Implementation

- [x] Add `-resultfile` to the documented launch contract
- [x] Emit result JSON to a per-match file path
- [x] Keep stdout JSON available for local debugging
- [x] Update smoke paths to exercise the result-file lane
- [x] Update stream-box operator scripts to surface result files
- [x] Document request-to-flag mapping for bridge launch intent
- [x] Document the current result/live JSON shape as V1
- [x] Add a C# model/interface library for V1 contract consumers

## Validation

- The local repository can prove the `-resultfile` and `matchOver` contracts
  here; live Windows rig validation still depends on `cc-desktopbridge`.

- [x] Run the local workflow smoke matrix with the result-file case
- [ ] Run the stream-box smoke path against the live Windows rig
- [x] Confirm the result file contains the expected JSON payload
- [x] Confirm the stdout lane can stay enabled alongside `-resultfile` without
  changing the result-file contract
- [x] Confirm a known-failing launch writes the panic to stderr with
  `-noerrordialog` active
- [x] Build the C# V1 contract library
- [x] Add local contract tests that compare Go JSON tags with C# model
  `JsonPropertyName` attributes and the V1 docs
- [x] Add `make contracts` as the repo-native V1 contract validation hook
- [x] Prove the live overlay lane in a real match loop: load
  `-overlayfile`, refresh live overlays, and draw text/image entries during a
  running fight.

## Follow-up

- [ ] Decide whether the flag surface should be replaced by a structured launch
  payload
- [ ] Decide whether live round-event transport needs a named pipe or callback
  lane
- [ ] Decide whether match-result files should carry additional metadata beyond
  the JSON stats snapshot
