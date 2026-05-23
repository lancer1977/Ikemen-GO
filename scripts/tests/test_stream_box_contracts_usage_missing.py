from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_usage_commands


class ParseUsageMissingTests(unittest.TestCase):
    def test_parse_usage_commands_returns_empty_when_missing(self) -> None:
        self.assertEqual(parse_usage_commands("Usage: other-command <noop>\n"), [])


if __name__ == "__main__":
    unittest.main()
