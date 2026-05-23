from __future__ import annotations

import unittest

from scripts.stream_box_contracts import unknown_commands


class UnknownCommandsIgnoreExtraTests(unittest.TestCase):
    def test_unknown_commands_ignores_extra_cases(self) -> None:
        self.assertEqual(unknown_commands(["install"], ["install", "check", "launch"]), [])


if __name__ == "__main__":
    unittest.main()
