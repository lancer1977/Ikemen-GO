from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_usage_commands


class ParseUsageTrailingPipeTests(unittest.TestCase):
    def test_parse_usage_commands_ignores_trailing_pipe(self) -> None:
        text = "Usage: scripts/stream-box/ikemen-box.sh <install|check|launch|>\n"
        self.assertEqual(parse_usage_commands(text), ["install", "check", "launch"])


if __name__ == "__main__":
    unittest.main()
