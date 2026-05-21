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

## Validation

- [ ] Run the local workflow smoke matrix with the result-file case
- [ ] Run the stream-box smoke path against the live Windows rig
- [ ] Confirm the result file contains the expected JSON payload
- [ ] Confirm the bridge can keep the stdout lane enabled without changing the
  result-file contract

## Follow-up

- [ ] Decide whether the flag surface should be replaced by a structured launch
  payload
- [ ] Decide whether live round-event transport needs a named pipe or callback
  lane
- [ ] Decide whether match-result files should carry additional metadata beyond
  the JSON stats snapshot
