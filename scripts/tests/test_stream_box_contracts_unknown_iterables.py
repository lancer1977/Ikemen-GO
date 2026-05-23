from __future__ import annotations

import unittest

from scripts.stream_box_contracts import unknown_commands


class UnknownCommandsIterablesTests(unittest.TestCase):
    def test_unknown_commands_supports_iterable_inputs(self) -> None:
        self.assertEqual(
            unknown_commands({"install", "check", "launch"}, {"install", "status"}),
            ["check", "launch"],
        )
