from __future__ import annotations

import unittest

from scripts.stream_box_contracts import unknown_commands


class UnknownCommandsIterableCaseTests(unittest.TestCase):
    def test_unknown_commands_with_iterables_preserves_order(self) -> None:
        self.assertEqual(unknown_commands(("check", "install"), ("install",)), ["check"])
