# Windows Result Transport

## Summary

The stream rig must stay Windows-native because the HDMI capture path and
MilkDrop workflow depend on Windows. Linux remains useful for fast local and CI
validation, but it cannot replace the stream-box visual deployment target.

The deploy contract should therefore make Windows result capture durable instead
of depending on GUI-process stdout.

## Current Decision

- [x] Keep `stream-box` as the primary visible Windows validation target
- [x] Keep `C:\\mugen` as the Ikemen deploy root
- [x] Suppress blocking Windows crash UI during smoke runs
- [x] Capture stderr/stdout under `C:\\Apps`
- [x] Track bridge-side OS-aware strategy selection in
  `cc-desktopbridge/docs/roadmaps/ikemen-result-transport-roadmap.md`
- [x] Add a durable fight-result file contract
- [ ] Treat named pipe or local callback transport as a future live-event goal

## Near-Term Contract

Add a result-file argument to the Ikemen fork:

```text
-resultfile C:\Apps\ikemen-results\<runId>.json
```

The expected flow:

- [ ] `cc-desktopbridge` creates a run id
- [x] bridge passes `-resultfile` when launching Ikemen
- [x] Ikemen writes the final fight result atomically
- [ ] bridge waits for process exit or file creation
- [ ] bridge parses winner, rounds, fighters, and exit metadata
- [ ] smoke fails on missing file, invalid JSON, crash log, or timeout

This should become the primary Windows automation contract. Console stdout can
remain useful for Linux and local developer smoke, but it is too fragile as the
only Windows bridge surface.

## Future Goal

Add a live result/event channel after the file contract is stable.

Candidate transports:

- [ ] Named pipe from Ikemen to `cc-desktopbridge`
- [ ] Localhost HTTP callback to the bridge
- [ ] WebSocket or SignalR relay if we later need round-by-round state

Use this only when we need live match events, not just final result capture. The
first version should focus on reliable final result delivery.

## Validation

- [x] Linux smoke still emits useful JSON for fast local testing
- [x] Windows smoke writes a result file for quick-vs
- [x] stream-box smoke proves no crash popup and no stale process
- [ ] bridge smoke proves it can read and forward the result
- [ ] result file cleanup avoids stale run reuse

## Cross-Repo Ownership

- `Ikemen-GO` owns the CLI transports and result schema emission.
- `cc-desktopbridge` owns OS detection, strategy selection, launch arguments,
  result collection, and forwarding to the rest of the ChannelCheevos stack.
