from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_usage_commands


class ParseUsageDuplicatesTests(unittest.TestCase):
    def test_parse_usage_commands_preserves_duplicates(self) -> None:
        text = (
            "Usage: scripts/stream-box/ikemen-box.sh <install|check|install|launch>\n"
        )
        self.assertEqual(parse_usage_commands(text), ["install", "check", "install", "launch"])


if __name__ == "__main__":
    unittest.main()
