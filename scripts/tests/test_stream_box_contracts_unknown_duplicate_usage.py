from __future__ import annotations

import unittest

from scripts.stream_box_contracts import unknown_commands


class UnknownCommandsDuplicateUsageTests(unittest.TestCase):
    def test_unknown_commands_deduplicates_with_duplicate_usage_inputs(self) -> None:
        self.assertEqual(
            unknown_commands(["install", "install", "check"], ["install"]),
            ["check"],
        )
