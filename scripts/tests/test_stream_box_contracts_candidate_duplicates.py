from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_candidate_quality_options


class ParseCandidateDuplicatesTests(unittest.TestCase):
    def test_parse_candidate_quality_options_handles_duplicates(self) -> None:
        text = """candidate-quality)
  --dry-run
  --dry-run
  --out
esac
"""
        self.assertEqual(parse_candidate_quality_options(text), {"dry-run", "out"})


if __name__ == "__main__":
    unittest.main()
