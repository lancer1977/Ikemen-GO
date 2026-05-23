from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_usage_commands


class ParseUsageWhitespaceOnlyTests(unittest.TestCase):
    def test_parse_usage_commands_with_whitespace_only_text_returns_empty(self) -> None:
        self.assertEqual(parse_usage_commands("   "), [])
