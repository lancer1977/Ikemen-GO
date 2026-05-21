#!/usr/bin/env python3
from __future__ import annotations

import argparse
import os
import shutil
import sys
from pathlib import Path


def parse_roster_entries(select_def: Path) -> list[tuple[str, str]]:
    entries: list[tuple[str, str]] = []
    in_characters = False

    for raw in select_def.read_text(encoding="utf-8", errors="ignore").splitlines():
        line = raw.strip()

        if not in_characters:
            if line.lower() == "[characters]":
                in_characters = True
            continue

        if line.startswith("["):
            break
        if not line or line.startswith(";"):
            continue
        if line.lower() == "randomselect":
            continue
        if line.startswith("slot =") or line == "}":
            continue

        name = line.split(",", 1)[0].strip()
        if name:
            entries.append((name, line))

    return entries


def find_insertion_index(lines: list[str]) -> int:
    in_characters = False
    last_content_idx: int | None = None

    for idx, raw in enumerate(lines):
        line = raw.strip()

        if not in_characters:
            if line.lower() == "[characters]":
                in_characters = True
            continue

        if line.startswith("["):
            return (last_content_idx + 1) if last_content_idx is not None else idx
        if not line or line.startswith(";"):
            continue
        last_content_idx = idx

    if not in_characters:
        raise RuntimeError("did not find a [Characters] section")

    return (last_content_idx + 1) if last_content_idx is not None else len(lines)


def update_select_def(select_def: Path, new_lines: list[str]) -> bool:
    if not new_lines:
        return False

    old_text = select_def.read_text(encoding="utf-8", errors="ignore")
    lines = old_text.splitlines()
    insert_at = find_insertion_index(lines)
    updated = lines[:insert_at] + new_lines + lines[insert_at:]
    new_text = "\n".join(updated) + "\n"

    if old_text == new_text:
        return False

    select_def.write_text(new_text, encoding="utf-8")
    return True


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Copy the next missing source characters into ~/apps/ikemen-dev in batches."
    )
    parser.add_argument("--dry-run", action="store_true", help="show what would change without writing files")
    parser.add_argument("--count", type=int, default=5, help="how many missing characters to add at once")
    args = parser.parse_args()

    source_root = Path(os.environ.get("IKEMEN_SOURCE_ROOT", str(Path.home() / "apps" / "ikemen-source"))).expanduser()
    dest_root = Path(os.environ.get("IKEMEN_LANCERO_TEST_ROOT", str(Path.home() / "apps" / "ikemen-dev"))).expanduser()

    source_select = source_root / "data" / "select.def"
    dest_select = dest_root / "data" / "select.def"
    source_chars = source_root / "chars"
    dest_chars = dest_root / "chars"

    if not source_select.is_file():
        print(f"ERROR: source select.def not found: {source_select}", file=sys.stderr)
        return 1
    if not dest_select.is_file():
        print(f"ERROR: destination select.def not found: {dest_select}", file=sys.stderr)
        return 1

    source_entries = parse_roster_entries(source_select)
    dest_names = {name for name, _ in parse_roster_entries(dest_select)}
    missing = [(name, line) for name, line in source_entries if name not in dest_names]
    selected = missing[: args.count]

    if args.dry_run:
        print("Dry run: no files will be changed.")

    if not selected:
        print("No missing source roster entries found.")
        return 0

    for name, line in selected:
        src_dir = source_chars / name
        dst_dir = dest_chars / name

        if not src_dir.is_dir():
            print(f"skip missing source char dir {name}")
            continue

        if dst_dir.exists():
            print(f"keep char dir {name}")
        else:
            print(f"copy char dir {name}")
            if not args.dry_run:
                shutil.copytree(src_dir, dst_dir, copy_function=shutil.copy2)

        print(f"add roster entry {name}")

    selected_lines = [line for name, line in selected if name not in dest_names]
    if selected_lines:
        if args.dry_run:
            print("update select.def")
        else:
            if update_select_def(dest_select, selected_lines):
                print("update select.def")
            else:
                print("select.def unchanged")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
