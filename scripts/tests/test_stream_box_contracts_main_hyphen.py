from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_main_case_commands


class ParseMainHyphenTokensTests(unittest.TestCase):
    def test_parse_main_case_commands_handles_hyphen_tokens(self) -> None:
        text = """case "$1" in
normalize-chars)
normalize-chars-apply)
candidate-quality)
*)
esac
"""
        self.assertEqual(parse_main_case_commands(text), {"normalize-chars", "normalize-chars-apply", "candidate-quality"})


if __name__ == "__main__":
    unittest.main()
