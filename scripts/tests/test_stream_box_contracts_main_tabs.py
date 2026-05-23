from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_main_case_commands


class ParseMainTabsTests(unittest.TestCase):
    def test_parse_main_case_commands_handles_tabs(self) -> None:
        text = """case "$1" in
\tinstall)\t
\tcheck)\t
*)
esac
"""
        self.assertEqual(parse_main_case_commands(text), {"install", "check"})
