from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_usage_commands


class ParseUsageEmptyTokenTests(unittest.TestCase):
    def test_parse_usage_commands_ignores_empty_token_inside(self) -> None:
        text = "Usage: scripts/stream-box/ikemen-box.sh <install||check>\n"
        self.assertEqual(parse_usage_commands(text), ["install", "check"])
