from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_usage_commands


class ParseUsageCrlfTests(unittest.TestCase):
    def test_parse_usage_commands_handles_crlf(self) -> None:
        text = "Usage: scripts/stream-box/ikemen-box.sh <install|check>\r\n"
        self.assertEqual(parse_usage_commands(text), ["install", "check"])


if __name__ == "__main__":
    unittest.main()
