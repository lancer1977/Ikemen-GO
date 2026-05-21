#!/usr/bin/env python3
from __future__ import annotations

import argparse
import datetime as dt
import os
from pathlib import Path


def parse_roster_entries(select_def: Path) -> list[str]:
    entries: list[str] = []
    in_characters = False

    for raw in select_def.read_text(encoding="utf-8", errors="ignore").splitlines():
        line = raw.rstrip()
        stripped = line.strip()

        if not in_characters:
            if stripped.lower() == "[characters]":
                in_characters = True
            continue

        if stripped.startswith("["):
            break
        if not stripped or stripped.startswith(";"):
            continue
        if stripped.lower() == "randomselect":
            continue
        if stripped.startswith("slot =") or stripped == "}":
            continue

        entries.append(stripped)

    return entries


def list_char_defs(chars_root: Path) -> list[str]:
    out: list[str] = []
    if not chars_root.is_dir():
        return out

    for char_dir in sorted((p for p in chars_root.iterdir() if p.is_dir()), key=lambda p: p.name.lower()):
        for def_file in sorted(char_dir.rglob("*.def")):
            out.append(str(def_file.relative_to(chars_root.parent)).replace(os.sep, "/"))
    return out


def list_stage_defs(stages_root: Path) -> list[str]:
    out: list[str] = []
    if not stages_root.is_dir():
        return out
    for def_file in sorted(stages_root.rglob("*.def")):
        out.append(str(def_file.relative_to(stages_root.parent)).replace(os.sep, "/"))
    return out


def main() -> int:
    parser = argparse.ArgumentParser(description="Write a code-backed backup of the current def tree.")
    parser.add_argument(
        "--source-root",
        default=os.environ.get("IKEMEN_SOURCE_ROOT", str(Path.home() / "apps" / "ikemen-source")),
        help="source tree root",
    )
    parser.add_argument(
        "--target-root",
        default=os.environ.get("IKEMEN_LANCERO_TEST_ROOT", str(Path.home() / "apps" / "ikemen-dev")),
        help="live install root to snapshot",
    )
    parser.add_argument(
        "--output",
        default=None,
        help="backup file path (defaults to build/lancero/def-tree.backup.md)",
    )
    args = parser.parse_args()

    source_root = Path(args.source_root).expanduser()
    target_root = Path(args.target_root).expanduser()
    output = Path(args.output).expanduser() if args.output else Path(__file__).with_name("def-tree.backup.md")

    select_def = target_root / "data" / "select.def"
    chars_root = target_root / "chars"
    stages_root = target_root / "stages"

    if not select_def.is_file():
        raise SystemExit(f"ERROR: select.def not found: {select_def}")

    lines: list[str] = []
    lines.append("# Def Tree Backup")
    lines.append("")
    lines.append(f"Generated: {dt.datetime.now(dt.timezone.utc).isoformat()}")
    lines.append(f"Source root: {source_root}")
    lines.append(f"Target root: {target_root}")
    lines.append("")
    lines.append("## select.def [Characters]")
    lines.append("")
    for entry in parse_roster_entries(select_def):
        lines.append(f"- {entry}")
    lines.append("")
    lines.append("## Character Defs")
    lines.append("")
    for rel in list_char_defs(chars_root):
        lines.append(f"- {rel}")
    lines.append("")
    lines.append("## Stage Defs")
    lines.append("")
    for rel in list_stage_defs(stages_root):
        lines.append(f"- {rel}")
    lines.append("")

    output.parent.mkdir(parents=True, exist_ok=True)
    output.write_text("\n".join(lines), encoding="utf-8")
    print(output)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
