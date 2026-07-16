# Workflow Matrix

The workflow smoke matrix exists in two forms:

- [`scripts/smoke/ikemen-workflow-matrix.sh`](../../scripts/smoke/ikemen-workflow-matrix.sh)
  for local Linux and MSYS2 Bash runs
- [`.github/workflows/workflow-smoke-matrix.yml`](../../.github/workflows/workflow-smoke-matrix.yml)
  for self-hosted Linux and Windows runners

## Runner layout

- Linux lane: `self-hosted`, `linux`, `ikemen-linux`
- Windows lane: `self-hosted`, `windows`, `ikemen-windows`
- Formatting validation lane: `self-hosted`, `linux`, `x64`, `pr-validation`,
  `ikemen-go`. This repository-scoped runner is for trusted same-repository
  validation only; fork pull requests are not routed to it.

## Formatting validation

Run the same non-mutating check used by `lint-code-style` locally:

```bash
./scripts/check-gofmt.sh
```

If it reports files, format those files with `gofmt -w` and rerun the check.
The workflow has read-only repository permissions and does not commit or push
changes.

## Local usage

Run the matrix directly from a checkout:

```bash
./scripts/smoke/ikemen-workflow-matrix.sh
```

To point the matrix at a mounted content tree, pass `--fixture-root PATH`
instead of exporting `IKEMEN_WORKFLOW_FIXTURE_ROOT`.

To run one lane at a time, pass `--case NAME` with a lane such as
`quickvs-resultfile` or `windowed-jsonstdout`.

For custom fixture layouts, `--char-dir PATH` overrides the default `chars`
subdirectory, and `--char-sweep-watchdog` / `--no-char-sweep-watchdog` controls
whether long character runs are reported as watchdog failures.

On Windows with MSYS2 Bash available, use:

```cmd
scripts\smoke\ikemen-workflow-matrix.cmd
```

## Workflow usage

Dispatch the workflow from GitHub Actions and choose:

- `both` to run both lanes
- `linux` to run only the Linux self-hosted runner
- `windows` to run only the Windows self-hosted runner

The workflow also supports toggling the per-character sweep and overriding the
fixture root when a different runtime tree is mounted on the runners.
