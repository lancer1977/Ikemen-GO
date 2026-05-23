from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_usage_commands


class ParseUsagePathWithAngleTests(unittest.TestCase):
    def test_parse_usage_commands_ignores_nested_angle_brackets(self) -> None:
        self.assertEqual(parse_usage_commands("Usage: scripts/stream-box/ikemen-box.sh <install <nested>>\n"), ["install <nested"])
