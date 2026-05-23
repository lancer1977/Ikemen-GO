from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_usage_commands


class ParseUsageExtraCharactersTests(unittest.TestCase):
    def test_parse_usage_commands_accepts_extra_whitespace_and_text(self) -> None:
        text = (
            "Usage: scripts/stream-box/ikemen-box.sh <install|check>\n"
            "Extra details follow."
        )
        self.assertEqual(parse_usage_commands(text), ["install", "check"])
