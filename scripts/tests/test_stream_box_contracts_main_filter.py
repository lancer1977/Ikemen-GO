from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_main_case_commands


class ParseMainFilterTests(unittest.TestCase):
    def test_parse_main_case_commands_filters_non_command_lines(self) -> None:
        text = """case "$1" in
install)
foo=$1
bad-token)
*)
esac
"""
        self.assertEqual(parse_main_case_commands(text), {"install", "bad-token"})


if __name__ == "__main__":
    unittest.main()
