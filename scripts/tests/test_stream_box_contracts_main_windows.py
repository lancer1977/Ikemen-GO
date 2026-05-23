from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_main_case_commands


class ParseMainWindowsNewlinesTests(unittest.TestCase):
    def test_parse_main_case_commands_windows_newlines(self) -> None:
        text = "case \"$1\" in\r\ninstall)\r\ncheck)\r\n*)\r\nesac\r\n"
        self.assertEqual(parse_main_case_commands(text), {"install", "check"})


if __name__ == "__main__":
    unittest.main()
