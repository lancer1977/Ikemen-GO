from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_candidate_quality_options


class ParseCandidatePrefixTests(unittest.TestCase):
    def test_parse_candidate_quality_options_ignores_prefixed_dash_dash(self) -> None:
        text = """candidate-quality)
  ---help
  ----force
esac
"""
        self.assertEqual(parse_candidate_quality_options(text), {"-help", "--force"})
