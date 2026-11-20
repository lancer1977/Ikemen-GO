# patches/ and uncommitted work (written 2026-10-06)

Context for whoever picks up `~/code/Ikemen-GO` on the `mugen` host. Nothing
here is committed or pushed. Branch: `feature/auto-skip-intro` (no upstream).
The running engine binary is `~/apps/ikemen-dev/Ikemen_GO_Linux` (built
2026-10-01); this tree is its source.

## Uncommitted change in the working tree: `src/system.go`

Adds an `-autoskipintro <seconds>` command-line flag (+35/-4 lines):

- New `introFrames` counter on `SystemStateVars`, reset in `resetRound()`.
- New `autoSkipIntroFrames()` parses the flag (0 = off).
- `runIntroSkip()` now force-skips an intro that runs past the limit, even for
  characters that assert `RoundNotSkip` (the code comment cites long intros
  such as Vergil KA in unattended CPU-vs-CPU runs).

`.gitignore` also gained 3 lines. `.serena/` and `.ignore` are tooling files.

**Unverified:** that it builds, and whether any service passes
`-autoskipintro` (a grep of `~/apps/ikemen-dev/scripts` and the systemd user
units found no use of it). `ikemen-intro-skip.service` exists and may be a
separate mechanism.

## Patch files

These are `git apply`-able diffs against `src/`. They are dated by the day
they were saved. Sizes from `git apply --stat`:

| File | Touches | Size |
| --- | --- | --- |
| `2026-09-29-betting-countdown-hold.patch` | `fightscreen.go`, `system.go` | +57/-5 |
| `2026-09-30-hold-plus-setPalNo.patch` | `fightscreen.go`, `script.go`, `system.go` | +76/-5 |
| `2026-10-01-exact-fight-tier.patch` | `fightscreen.go`, `script.go`, `system.go` | +165/-5 |
| `2026-10-01-no-bet-countdown-env.patch` | `fightscreen.go` | +59/-5 |

The first three touch the same files with the same leading hunks and grow
in size, so they look like **successive iterations of one change, not
independent features**. That is inferred from the stats, not from reading
each patch end to end. The fourth is a smaller `fightscreen.go`-only variant.
Before applying any, diff them against each other and against current
`develop`; applying all four will conflict.

The HEAD commits (`afc6866f Require strict tier tournament pools`,
`752e5e18 Add endless random tournament OBS tracking`) may already contain
some of this work. Not checked.

## Suggested next steps

1. Decide which patch (if any) is the final form; archive the rest.
2. Build (`test.sh` is the only script at the repo root) and run the
   `-autoskipintro` change before committing.
3. Commit the `system.go` change on its own branch with an issue/PR; push
   only after that.
