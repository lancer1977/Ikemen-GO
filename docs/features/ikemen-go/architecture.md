# Ikemen Launch Architecture

## Runtime Flow

1. ChannelCheevos (or another launcher) builds a temporary run ID.
2. CLI args are passed to Ikemen GO through:
   - `-p<n>` family for players and overrides
   - `-tmode<n>`, `-time`, `-rounds`, `-s` for match configuration
   - `-resultfile <path>` for durable bridge consumption
   - `-livedatafile <path>` for live snapshots during an active match
3. `src/main.go` parses and forwards flags unchanged to `external/script/main.lua`.
4. `main.lua` runs startup initialization in safe mode:
   - normal startup path first
   - quick-vs path when both `-p1` and `-p2` are supplied without `-loadmotif`
5. On match completion, stats are serialized once with `getGameStatsJson()` and
   written to:
   - JSON file via `-jsonlog` (default `save/last-match.json`)
   - canonical result file via `-resultfile`
   - optional stdout JSON via `-jsonstdout`

## Failure Surfaces

- Startup null-safe wrappers (`safeGameOption`, `safeCommandLineValue`, `safeCall`,
  protected file IO) reduce fatal panics from missing assets.
- Smoke scripts treat common runtime crash signatures as hard failures:
  - `attempt to call a non-function object`
  - `runtime error`
  - `panic`
  - `fatal error`
  - `segmentation fault`
  - `stack traceback`

## Data Contracts

- `-resultfile` is the primary transport for automation/bridge workflows.
- `-jsonstdout` remains available for operator diagnostics and local workflow
  smoke.
- `-livedatafile` streams atomic match snapshots while the fight is active.
- Result output is expected to be JSON containing `statsLog`.

## Operating Notes

- `scripts/stream-box` scripts launch with `-jsonstdout` plus a per-run `-resultfile`
  path.
- Debug smoke helpers keep process stdio in `C:\\Apps` and suppress Windows crash
  popups.
