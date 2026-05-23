from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_main_case_commands


class ParseMainSimpleTests(unittest.TestCase):
    def test_parse_main_case_commands_parses_simple_block(self) -> None:
        text = """case "$1" in
install)
check)
launch)
*)
esac
"""
        self.assertEqual(parse_main_case_commands(text), {"install", "check", "launch"})


if __name__ == "__main__":
    unittest.main()
