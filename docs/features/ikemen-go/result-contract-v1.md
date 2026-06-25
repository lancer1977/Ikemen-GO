# Ikemen GO Result Contract V1

## Summary

This document defines the current V1 JSON contract emitted by Ikemen GO for
automation and bridge consumers.

V1 is intentionally a documentation and model contract for the existing output.
It does not add fields to the engine payload. The payload is currently
unversioned JSON; consumers should treat the shapes below as
`IkemenGo.MatchState.V1`.

## Transports

- `-resultfile <jsonfile>` writes the final match JSON after the fight ends.
  This is the canonical durable bridge surface.
- `-jsonlog <jsonfile>` writes the same final match JSON. When omitted and
  `-nojsonlog` is not set, Ikemen GO writes `save/last-match.json`.
- `-jsonstdout` writes the same final match JSON to stdout.
- `-livedatafile <jsonfile>` writes a live snapshot JSON during the fight and
  refreshes it once more after the final result is available.
- `-overlayfile <jsonfile>` writes a live overlay snapshot contract that lets a
  sidecar request text, emoji, and PNG-backed image effects during a fight.

The result and live files are written atomically. Consumers should read the
whole file as a JSON document, not tail it as an append-only event log.

## Result Snapshot

The final result payload is represented by
`Ikemen.Go.Contracts.V1.GameStatsSnapshotV1`.

```json
{
  "statsLog": {
    "matches": []
  },
  "continueFlg": false,
  "persistRoundCount": 0,
  "matchOver": true
}
```

### Fields

- `statsLog`: stats container for the current game session.
- `matches`: ordered match records inside `statsLog`.
- `statsLog.matches`: dotted path to the ordered match records accumulated by
  the current game session.
- `continueFlg`: current continue flag from the engine.
- `persistRoundCount`: persisted round counter from the engine.
- `matchOver`: true when the current fight has reached the engine match-over
  condition.

## Live Snapshot

The live payload is represented by
`Ikemen.Go.Contracts.V1.GameLiveSnapshotV1`.

```json
{
  "statsLog": {
    "matches": []
  },
  "continueFlg": false,
  "persistRoundCount": 0,
  "matchOver": false,
  "frameCounter": 0,
  "matchTime": 0,
  "curRoundTime": 0,
  "currentRound": {
    "index": 1,
    "timer": 0,
    "score": [0, 0],
    "fighters": [[], []]
  }
}
```

### Additional Live Fields

- `frameCounter`: current engine frame counter.
- `matchTime`: accumulated match time in ticks.
- `curRoundTime`: current round timer in ticks.
- `currentRound`: pollable in-flight round snapshot.

## Match

`Ikemen.Go.Contracts.V1.StatsMatchV1`

- `matchTime`: total match time in ticks.
- `roundTime`: configured round time in ticks.
- `winSide`: winning side as currently emitted by the engine.
- `ended`: true once the fight has fully ended.
- `lastRound`: final round index played, 1-based.
- `draws`: drawn round count.
- `wins`: two-element array, `[P1Wins, P2Wins]`.
- `teamModes`: two-element array, engine team mode per side.
- `totalScore`: two-element array, cumulative score per side.
- `rounds`: ordered round records.

## Round

`Ikemen.Go.Contracts.V1.StatsRoundV1`

- `index`: 1-based round number.
- `timer`: round timer value in ticks.
- `score`: two-element array, `[P1Score, P2Score]`.
- `fighters`: two-element array of fighter lists, `[P1Side, P2Side]`.

## Fighter

`Ikemen.Go.Contracts.V1.StatsFighterStateV1`

- `name`: character name.
- `id`: runtime character ID.
- `memberNo`: 0-based team member index.
- `selectNo`: select screen index.
- `aiLevel`: CPU AI level, where `0` means human.
- `palNo`: palette number.
- `life`: current life at snapshot time.
- `lifeMax`: maximum life.
- `power`: current power at snapshot time.
- `powerMax`: maximum power.
- `winQuote`: selected win quote, or `-1` when unused.
- `win`: true when this fighter's side won the round.
- `winKO`: true when the side won by KO.
- `winTime`: true when the side won by time-out.
- `winPerfect`: true when the side won perfectly.
- `winSpecial`: true when the side won with a special.
- `winHyper`: true when the side won with a hyper.
- `drawGame`: true when the round was declared a draw.
- `ko`: true when this fighter was KO'd.
- `overKO`: true when this fighter is in over-KO state.

## Explicit Non-Goals In V1

- V1 is not an append-only event stream.
- V1 does not emit discrete damage, hit, combo, match-start, or round-start
  events.
- V1 does not include full select-screen character lists or stage lists.
- V1 does not include selected stage metadata in the final result payload.
- V1 does not replace the overlay snapshot contract with the result snapshot;
  live overlays are a separate file contract.

Those are candidates for a later V2 or for a separate live event transport.
