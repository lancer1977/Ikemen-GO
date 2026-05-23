from __future__ import annotations

import unittest

from scripts.stream_box_contracts import parse_usage_commands, parse_main_case_commands, unknown_commands


class UnknownCommandsComboRealUsageTests(unittest.TestCase):
    def test_unknown_commands_parser_combo_parses_real_usage_and_cases(self) -> None:
        usage_text = (
            "Usage: scripts/stream-box/ikemen-box.sh "
            "<install|check|launch|launch-menu|status|tail|stop|repoint|bundle|"
            "normalize-chars|normalize-chars-apply|candidate-quality>\n"
        )
        case_text = """case "$1" in
install)
check)
launch)
launch-menu)
status)
tail)
stop)
repoint)
bundle)
normalize-chars)
normalize-chars-apply)
candidate-quality)
*)
esac
"""
        self.assertEqual(
            unknown_commands(parse_usage_commands(usage_text), parse_main_case_commands(case_text)),
            [],
        )
