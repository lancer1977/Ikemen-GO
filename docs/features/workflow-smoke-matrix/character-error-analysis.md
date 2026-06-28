# Character Error Analysis

## Summary

This note tracks the current per-character sweep failures from the
`/mnt/shared/games/ikemen2` fixture set.

## Current Findings

- [x] Rebuilt `Ikemen_GO_Linux` with `-noerrordialog` support.
- [x] Confirmed the sweep can run without blocking on a modal dialog.
- [x] Preserved rerun logs under `/tmp/ikemen-char-sweep-logs-postfix2`.
- [x] Added AI level arguments to character self-play sweep launches.
- [x] Backfilled missing core `data/*.zss` files into fixture work roots before
  launch (`scripts/smoke/ikemen-workflow-matrix.sh`).
- [x] Added metadata enrichment and folder-name normalization in sweep JSON
  (`schema: ikemen-char-sweep/v2`).
- [x] Added character metadata overrides file:
  `docs/features/workflow-smoke-matrix/character-metadata-overrides.tsv`.
- [x] Documented `brolyz2` as `long-intro` in metadata tags/notes.
- [x] Normalized `stream-box` character folder names and `.def` filenames to
  remove whitespace-delimited launch path failures.
- [x] Investigate one remaining deterministic parser/content failure
  (`kfmZ/kfm.zss:657`).
- [x] Investigate one remaining runtime timeout (`BrolyZ2`).
- [x] Update the ZSS `animtype` parser to accept `Med` shorthand for legacy
  content in source.
- [x] Add quick-vs watchdog reporting so a timeout with no result file is
  tracked as a harness failure instead of a character compatibility verdict.
- [x] Extend sweep timeouts for `long-intro` characters so `BrolyZ2` can use a
  longer run window without weakening the default matrix.

## Prior Throwing Characters

The completed sweep report in `ikemen-char-sweep-results.json` captured 15
throwing characters. All failed with `exitCode: 1` and no result file.

Common signature:

```text
Panic: external/script/main.lua:1083
```

That line is the `game()` call in the command-line quick-vs path. Most entries
only preserved the selected character `.def` path because the temporary detailed
logs from that run were not kept.

| Character | Selected path | Captured error hint |
| --- | --- | --- |
| BrolyZ2 | `chars/BrolyZ2/BrolyZ2.def` | `main.lua:1083 -> BrolyZ2.def` |
| Chun_LI_JJ | `chars/Chun_LI_JJ/Chun_LI_JJ.def` | `main.lua:1083 -> Chun_LI_JJ.def` |
| Future_Gohan | `chars/Future_Gohan/Future_Gohan.def` | `main.lua:1083 -> Future_Gohan.def` |
| GokuZ2_1.0-ai | `chars/GokuZ2_1.0-ai/GokuZ2_1.0-ai.def` | `main.lua:1083 -> GokuZ2_1.0-ai.def` |
| Iori-KOF98 | `chars/Iori-KOF98/Iori-KOF98.def` | `main.lua:1083 -> Iori-KOF98.def` |
| PiccoloZ2 | `chars/PiccoloZ2/PiccoloZ2.def` | `main.lua:1083 -> PiccoloZ2.def` |
| Rei_AI-Patch | `chars/Rei_AI-Patch/Rei_AI-Patch.def` | `main.lua:1083 -> Rei_AI-Patch.def` |
| ShinSmoke Ryu | `chars/ShinSmoke Ryu/ShinSmoke Ryu.def` | `main.lua:1083 -> ShinSmoke Ryu.def` |
| Spiderman | `chars/Spiderman/Spiderman.def` | `main.lua:1083 -> Spiderman.def` |
| Thor_AvX | `chars/Thor_AvX/Thor_AvX.def` | `main.lua:1083 -> Thor_AvX.def` |
| kfm | `chars/kfm/kfm.def` | `main.lua:1083 -> kfm.def` |
| kfm720 | `chars/kfm720/kfm720.def` | `main.lua:1083 -> kfm720.def` |
| kfmZ | `chars/kfmZ/kfmZ.def` | `main.lua:1083 -> kfm.zss:657` |
| kfm_zaxis | `chars/kfm_zaxis/kfm_zaxis.def` | `main.lua:1083 -> kfm_zaxis.def` |
| kfm_zss | `chars/kfm_zss/kfm_zss.def` | `main.lua:1083 -> kfm_zss.def` |

## Latest Rebuilt Sweep

The post-fix sweep report is:

- `/tmp/ikemen-char-sweep-results-postfix2.json`

Latest status:

- `totalCharacters`: 17
- `testedCharacters`: 15
- `skippedCharacters`: 2
- `compatibleCount`: 13
- `incompatibleCount`: 2

Remaining incompatible characters:

- `BrolyZ2`: classified as a watchdog timeout rather than a parser/content failure. The latest rebuilt sweep hit the 35s window without a completed match/result file, while nearby peers finished normally.
- `kfmZ`: parser/content failure at `kfm.zss:657` with `animtype: Invalid animtype: Med` in the currently deployed binary. The source tree already accepts `Med`; it still needs to be rebuilt and rerun in the live fixture.

This confirms the earlier broad failure class (`open data/demo.zss`) was a
fixture bootstrap issue, not a character-wide engine regression.

## Follow-up

- [x] Preserve per-character stderr/stdout logs with
  `IKEMEN_WORKFLOW_CHAR_SWEEP_LOG_DIR`.
- [x] Add a quick-vs watchdog mode that treats a still-running process with no
  result file as a harness failure, not a character compatibility verdict.
- [x] Reproduce prior throwing characters with full logs.
- [ ] Rebuild and rerun the sweep to confirm `kfmZ/kfm.zss:657` clears after
  the `Med` compatibility update.
- [ ] Rerun the sweep to verify `BrolyZ2` completes with the longer
  `long-intro` timeout window and still writes `-resultfile`.
