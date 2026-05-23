from __future__ import annotations

import unittest

from scripts.stream_box_contracts import unknown_commands


class UnknownCommandsUnhashableTests(unittest.TestCase):
    def test_unknown_commands_raises_with_unhashable_members(self) -> None:
        with self.assertRaises(TypeError):
            unknown_commands(["install", ["nested"]], ["install"])
