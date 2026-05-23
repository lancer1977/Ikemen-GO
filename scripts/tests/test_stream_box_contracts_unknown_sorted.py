from __future__ import annotations

import unittest

from scripts.stream_box_contracts import unknown_commands


class UnknownCommandsSortedTests(unittest.TestCase):
    def test_unknown_commands_sorted(self) -> None:
        self.assertEqual(
            unknown_commands(["zeta", "alpha", "gamma"], ["zeta"]),
            ["alpha", "gamma"],
        )


if __name__ == "__main__":
    unittest.main()
