from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_usage_commands


class ParseUsageMultipleCommandsTests(unittest.TestCase):
    def test_parse_usage_commands_parses_multiple_commands(self) -> None:
        text = "Usage: scripts/stream-box/ikemen-box.sh <install|check|launch>\n"
        self.assertEqual(parse_usage_commands(text), ["install", "check", "launch"])


if __name__ == "__main__":
    unittest.main()
