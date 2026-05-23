from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_candidate_quality_options


class ParseCandidateMixedCaseTests(unittest.TestCase):
    def test_parse_candidate_quality_options_accepts_mixed_case(self) -> None:
        text = """candidate-quality)
  --Dry-Run
  --OutputDir
esac
"""
        self.assertEqual(parse_candidate_quality_options(text), {"dry-run", "outputdir"})
