from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_candidate_quality_options


class ParseCandidateEmptyBlockTests(unittest.TestCase):
    def test_parse_candidate_quality_options_empty_block(self) -> None:
        text = "candidate-quality)\n"
        self.assertEqual(parse_candidate_quality_options(text), set())


if __name__ == "__main__":
    unittest.main()
