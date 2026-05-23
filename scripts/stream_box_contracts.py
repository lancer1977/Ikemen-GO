from __future__ import annotations

import re
from typing import Iterable, List, Set

_USAGE_RE = re.compile(r"Usage:\s*scripts/stream-box/ikemen-box\.sh\s*<([^>]+)>")
_MAIN_CASE_RE = re.compile(r"^\s*([a-zA-Z0-9-]+)\)", re.M)
_CQ_OPTION_RE = re.compile(r"^\s*--([a-zA-Z0-9-]+)", re.M)


def parse_usage_commands(script_text: str) -> List[str]:
    match = _USAGE_RE.search(script_text)
    if not match:
        return []
    inner = match.group(1)
    return [token.strip() for token in inner.split("|") if token.strip()]


def parse_main_case_commands(script_text: str) -> Set[str]:
    start_match = re.search(r"case\s+\"(?:\\)?\$1\"\s+in", script_text)
    if not start_match:
        start_match = re.search(r"case\s+'(?:\\)?\$1'\s+in", script_text)
    if not start_match:
        return set()
    start = start_match.start()
    body = script_text[start:]
    lines = body.splitlines()
    commands = set()
    for line in lines:
        m = _MAIN_CASE_RE.match(line)
        if not m:
            continue
        cmd = m.group(1)
        if cmd in {"*", "-h", "--help", "help"}:
            continue
        commands.add(cmd)
    return commands


def parse_candidate_quality_options(script_text: str) -> Set[str]:
    start = script_text.find("candidate-quality)")
    if start == -1:
        return set()
    block = script_text[start:]
    options: Set[str] = set()
    for m in _CQ_OPTION_RE.finditer(block):
        options.add(m.group(1).lower())
    return options


def unknown_commands(usage: Iterable[str], cases: Iterable[str]) -> List[str]:
    usage_set = set(usage)
    case_set = set(cases)
    return sorted(usage_set - case_set)
