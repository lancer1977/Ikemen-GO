from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_usage_commands


class ParseUsageCommentaryTests(unittest.TestCase):
    def test_parse_usage_commands_ignores_non_usage_commentary(self) -> None:
        self.assertEqual(parse_usage_commands("# Usage: nothing"), [])
