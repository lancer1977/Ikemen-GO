from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_usage_commands, parse_main_case_commands, unknown_commands


class UnknownCommandsComboHasUnknownTests(unittest.TestCase):
    def test_unknown_commands_parser_combo_has_unknown(self) -> None:
        usage_text = "Usage: scripts/stream-box/ikemen-box.sh <install|check|launch>\n"
        case_text = """case "$1" in
install)
check)
*)
esac
"""
        self.assertEqual(
            unknown_commands(parse_usage_commands(usage_text), parse_main_case_commands(case_text)),
            ["launch"],
        )
