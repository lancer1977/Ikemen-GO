from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_candidate_quality_options


class ParseCandidateFlagsTests(unittest.TestCase):
    def test_parse_candidate_quality_options_parses_flags(self) -> None:
        text = """candidate-quality)
  --dry-run
  --force
  --report-json
esac
"""
        self.assertEqual(parse_candidate_quality_options(text), {"dry-run", "force", "report-json"})


if __name__ == "__main__":
    unittest.main()
