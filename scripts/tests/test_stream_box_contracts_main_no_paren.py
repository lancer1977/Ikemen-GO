from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_main_case_commands


class ParseMainNoParenTests(unittest.TestCase):
    def test_parse_main_case_commands_ignores_entries_without_trailing_paren(self) -> None:
        text = """case "$1" in
install
check)
*)
esac
"""
        self.assertEqual(parse_main_case_commands(text), {"check"})
