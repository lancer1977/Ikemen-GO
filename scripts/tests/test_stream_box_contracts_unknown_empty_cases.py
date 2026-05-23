from __future__ import annotations

import unittest

from scripts.stream_box_contracts import unknown_commands


class UnknownCommandsEmptyCasesTests(unittest.TestCase):
    def test_unknown_commands_with_empty_cases(self) -> None:
        self.assertEqual(unknown_commands(["install", "check"], []), ["check", "install"])
