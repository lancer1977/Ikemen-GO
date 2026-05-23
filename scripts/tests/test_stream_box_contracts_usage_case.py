from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_usage_commands


class ParseUsageCaseTests(unittest.TestCase):
    def test_parse_usage_commands_is_case_sensitive(self) -> None:
        text = "Usage: scripts/stream-box/ikemen-box.sh <Install|Check>\n"
        self.assertEqual(parse_usage_commands(text), ["Install", "Check"])
