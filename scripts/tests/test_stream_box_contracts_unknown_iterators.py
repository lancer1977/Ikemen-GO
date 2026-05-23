from __future__ import annotations

import unittest

from scripts.stream_box_contracts import unknown_commands


class UnknownCommandsIteratorsTests(unittest.TestCase):
    def test_unknown_commands_supports_iterators(self) -> None:
        self.assertEqual(unknown_commands(("install", "check"), ("install",)), ["check"])


if __name__ == "__main__":
    unittest.main()
