from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_main_case_commands


class ParseMainIgnoreBlankCommandsTests(unittest.TestCase):
    def test_parse_main_case_commands_ignores_blank_command_lines(self) -> None:
        text = """case "$1" in

install)

check)
*)  
esac
"""
        self.assertEqual(parse_main_case_commands(text), {"install", "check"})
