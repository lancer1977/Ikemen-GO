from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_candidate_quality_options


class ParseCandidateShortOptionTests(unittest.TestCase):
    def test_parse_candidate_quality_options_ignores_short_option(self) -> None:
        text = """candidate-quality)
  -h
  --help
  --output
esac
"""
        self.assertEqual(parse_candidate_quality_options(text), {"help", "output"})


if __name__ == "__main__":
    unittest.main()
