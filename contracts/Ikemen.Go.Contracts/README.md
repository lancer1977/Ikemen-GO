# Ikemen.Go.Contracts

Strong C# models for the existing Ikemen GO result and live snapshot JSON
contract.

The current engine payload is unversioned JSON. This library names that shape
`IkemenGo.MatchState.V1` so bridge and tooling consumers can compile against a
stable model while the engine keeps emitting the existing fields.

## Models

- `GameStatsSnapshotV1`: final JSON emitted by `-resultfile`, `-jsonlog`, and
  `-jsonstdout`.
- `GameLiveSnapshotV1`: live JSON emitted by `-livedatafile`.
- `StatsLogV1`, `StatsMatchV1`, `StatsRoundV1`, and
  `StatsFighterStateV1`: shared nested contract objects.

The authoritative field list is documented in
`docs/features/ikemen-go/result-contract-v1.md`.
