from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_main_case_commands


class ParseMainCommentLinesTests(unittest.TestCase):
    def test_parse_main_case_commands_ignores_comment_lines(self) -> None:
        text = """case "$1" in
# comment line
install)
# another comment
check)
*)
esac
"""
        self.assertEqual(parse_main_case_commands(text), {"install", "check"})
