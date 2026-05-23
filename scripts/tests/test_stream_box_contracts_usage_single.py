from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_usage_commands


class ParseUsageSingleCommandTests(unittest.TestCase):
    def test_parse_usage_commands_parses_single_command(self) -> None:
        text = "Usage: scripts/stream-box/ikemen-box.sh <install>\n"
        self.assertEqual(parse_usage_commands(text), ["install"])


if __name__ == "__main__":
    unittest.main()
