# Lancero Helpers

This directory holds source-owned helper scripts for the local Lancero test
install.

- `sync-content-to-lancero-dev.sh`: copies `chars/`, `stages/`, `sound/`, and `music/` from `~/apps/ikemen-source` into `~/apps/ikemen-dev`.
- `add-source-char-to-lancero-dev.sh`: adds one source character to `~/apps/ikemen-dev/data/select.def` and writes `start.sh` and `start.cmd` into the dev install.
- `add-missing-source-chars-to-lancero-dev.sh`: adds the next 5 missing source roster entries to `~/apps/ikemen-dev/data/select.def` and copies any missing character folders.
- `backup-def-tree.sh`: writes a code-backed snapshot of the current character and stage def tree.
