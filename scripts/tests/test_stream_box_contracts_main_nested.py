from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_main_case_commands


class ParseMainNestedTests(unittest.TestCase):
    def test_parse_main_case_commands_ignores_nested_case_tokens(self) -> None:
        text = """case "$1" in
install)
  case "$2" in
  sub)
*)
  esac
*)
esac
"""
        self.assertEqual(parse_main_case_commands(text), {"install", "sub"})


if __name__ == "__main__":
    unittest.main()
