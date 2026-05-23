import unittest

from scripts.stream_box_contracts import parse_usage_commands


class StreamBoxContractParserEmptyBracketsTests(unittest.TestCase):
    def test_parse_usage_commands_empty_brackets(self) -> None:
        text = "Usage: scripts/stream-box/ikemen-box.sh <>\n"
        self.assertEqual(parse_usage_commands(text), [])
