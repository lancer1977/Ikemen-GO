from __future__ import annotations

import unittest

from scripts.stream_box_contracts import unknown_commands


class UnknownCommandsMissingTests(unittest.TestCase):
    def test_unknown_commands_returns_missing_cases(self) -> None:
        self.assertEqual(
            unknown_commands(["install", "check", "launch"], ["install"]),
            ["check", "launch"],
        )


if __name__ == "__main__":
    unittest.main()
