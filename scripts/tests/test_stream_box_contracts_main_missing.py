from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_main_case_commands


class ParseMainMissingTests(unittest.TestCase):
    def test_parse_main_case_commands_empty_when_missing(self) -> None:
        self.assertEqual(parse_main_case_commands("echo none"), set())


if __name__ == "__main__":
    unittest.main()
