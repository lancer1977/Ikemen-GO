from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_usage_commands


class ParseUsageNoPrefixTests(unittest.TestCase):
    def test_parse_usage_commands_without_usage_prefix_returns_empty(self) -> None:
        self.assertEqual(parse_usage_commands("scripts/stream-box/ikemen-box.sh <install|check>\n"), [])
