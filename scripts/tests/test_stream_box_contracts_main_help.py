from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_main_case_commands


class ParseMainHelpCommandsTests(unittest.TestCase):
    def test_parse_main_case_commands_ignores_help_commands(self) -> None:
        text = """case "$1" in
install)
help)
-h)
--help)
*)
esac
"""
        self.assertEqual(parse_main_case_commands(text), {"install"})


if __name__ == "__main__":
    unittest.main()
