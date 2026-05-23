from __future__ import annotations

import unittest

from scripts.stream_box_contracts import unknown_commands


class UnknownCommandsFullMatchTests(unittest.TestCase):
    def test_unknown_commands_returns_empty_when_full_match(self) -> None:
        self.assertEqual(
            unknown_commands(["install", "check"], ["check", "install"]),
            [],
        )


if __name__ == "__main__":
    unittest.main()
