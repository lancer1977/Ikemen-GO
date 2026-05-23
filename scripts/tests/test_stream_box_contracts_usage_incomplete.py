from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_usage_commands


class ParseUsageIncompleteTests(unittest.TestCase):
    def test_parse_usage_commands_with_incomplete_angle_block(self) -> None:
        self.assertEqual(parse_usage_commands("Usage: scripts/stream-box/ikemen-box.sh <install|check\n"), [])
