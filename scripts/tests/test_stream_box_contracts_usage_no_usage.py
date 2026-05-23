from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_usage_commands


class ParseUsageNoUsageTests(unittest.TestCase):
    def test_parse_usage_commands_returns_empty_when_no_usage(self) -> None:
        self.assertEqual(parse_usage_commands(""), [])


if __name__ == "__main__":
    unittest.main()
