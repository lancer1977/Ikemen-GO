from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_usage_commands


class ParseUsageOrderTests(unittest.TestCase):
    def test_parse_usage_commands_preserves_order(self) -> None:
        text = "Usage: scripts/stream-box/ikemen-box.sh <zeta|alpha|omega>\n"
        self.assertEqual(parse_usage_commands(text), ["zeta", "alpha", "omega"])


if __name__ == "__main__":
    unittest.main()
