# Phases

## Phase 1: Stabilize the current contract

- [x] Keep the existing CLI flags documented and supported
- [x] Add a per-match result file path for bridge integration
- [x] Keep stdout JSON for local developer smoke
- [ ] Validate the Windows stream-box path after the engine change

## Phase 2: Reduce launch surface pressure

- [ ] Define a structured launch payload schema
- [ ] Decide which current flags should become payload fields
- [ ] Decide how the bridge should version launch requests
- [ ] Keep the CLI path working as the compatibility fallback

## Phase 3: Add live event transport

- [ ] Evaluate named pipe delivery for round events
- [ ] Evaluate localhost callback delivery for round events
- [ ] Keep the final result file as the fallback when live transport fails
- [ ] Capture the bridge-side retry and timeout behavior

## Phase 4: Automation exit modes

- [ ] Add an automation-friendly way to exit after a result is written
- [ ] Keep the normal interactive path unchanged
- [ ] Verify the mode does not interfere with local debug workflows

