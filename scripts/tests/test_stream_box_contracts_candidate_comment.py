from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_candidate_quality_options


class ParseCandidateCommentsTests(unittest.TestCase):
    def test_parse_candidate_quality_options_with_comments(self) -> None:
        text = """candidate-quality)
# comment
  --dry-run # dry run
  --output-dir # output path
esac
"""
        self.assertEqual(parse_candidate_quality_options(text), {"dry-run", "output-dir"})
