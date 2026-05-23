from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_main_case_commands


class ParseMainStarOnlyTests(unittest.TestCase):
    def test_parse_main_case_commands_empty_block_with_star_only(self) -> None:
        text = """case "$1" in
*)
esac
"""
        self.assertEqual(parse_main_case_commands(text), set())
