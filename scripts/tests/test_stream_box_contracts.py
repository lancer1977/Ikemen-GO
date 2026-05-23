from __future__ import annotations

import unittest

from scripts.stream_box_contracts import (
    parse_candidate_quality_options,
    parse_main_case_commands,
    parse_usage_commands,
    unknown_commands,
)


class StreamBoxContractParserTests(unittest.TestCase):
    def test_parse_usage_commands_parses_single_command(self) -> None:
        text = "Usage: scripts/stream-box/ikemen-box.sh <install>\n"
        self.assertEqual(parse_usage_commands(text), ["install"])

    def test_parse_usage_commands_parses_multiple_commands(self) -> None:
        text = "Usage: scripts/stream-box/ikemen-box.sh <install|check|launch>\n"
        self.assertEqual(parse_usage_commands(text), ["install", "check", "launch"])

    def test_parse_usage_commands_preserves_order(self) -> None:
        text = "Usage: scripts/stream-box/ikemen-box.sh <zeta|alpha|omega>\n"
        self.assertEqual(parse_usage_commands(text), ["zeta", "alpha", "omega"])

    def test_parse_usage_commands_strips_spaces(self) -> None:
        text = "Usage: scripts/stream-box/ikemen-box.sh < install | check | launch >\n"
        self.assertEqual(parse_usage_commands(text), ["install", "check", "launch"])

    def test_parse_usage_commands_handles_crlf(self) -> None:
        text = "Usage: scripts/stream-box/ikemen-box.sh <install|check>\r\n"
        self.assertEqual(parse_usage_commands(text), ["install", "check"])

    def test_parse_usage_commands_returns_empty_when_missing(self) -> None:
        self.assertEqual(parse_usage_commands("Usage: other-command <noop>\n"), [])

    def test_parse_usage_commands_returns_empty_when_no_usage(self) -> None:
        self.assertEqual(parse_usage_commands(""), [])

    def test_parse_usage_commands_ignores_trailing_pipe(self) -> None:
        text = "Usage: scripts/stream-box/ikemen-box.sh <install|check|launch|>\n"
        self.assertEqual(parse_usage_commands(text), ["install", "check", "launch"])

    def test_parse_usage_commands_preserves_duplicates(self) -> None:
        text = (
            "Usage: scripts/stream-box/ikemen-box.sh <install|check|install|launch>\n"
        )
        self.assertEqual(
            parse_usage_commands(text),
            ["install", "check", "install", "launch"],
        )

    def test_parse_usage_commands_accepts_extra_text(self) -> None:
        text = (
            "Usage: scripts/stream-box/ikemen-box.sh <install|check>\n"
            "Extra details follow."
        )
        self.assertEqual(parse_usage_commands(text), ["install", "check"])

    def test_parse_main_case_commands_parses_simple_block(self) -> None:
        text = """case \"$1\" in
install)
check)
launch)
*)
esac
"""
        self.assertEqual(
            parse_main_case_commands(text),
            {"install", "check", "launch"},
        )

    def test_parse_main_case_commands_ignores_help_commands(self) -> None:
        text = """case \"$1\" in
install)
help)
-h)
--help)
*)
esac
"""
        self.assertEqual(parse_main_case_commands(text), {"install"})

    def test_parse_main_case_commands_handles_hyphen_tokens(self) -> None:
        text = """case \"$1\" in
normalize-chars)
normalize-chars-apply)
candidate-quality)
*)
esac
"""
        self.assertEqual(
            parse_main_case_commands(text),
            {"normalize-chars", "normalize-chars-apply", "candidate-quality"},
        )

    def test_parse_main_case_commands_empty_when_missing(self) -> None:
        self.assertEqual(parse_main_case_commands("echo none"), set())

    def test_parse_main_case_commands_filters_non_command_lines(self) -> None:
        text = """case \"$1\" in
install)
foo=$1
bad-token)
*)
esac
"""
        self.assertEqual(
            parse_main_case_commands(text),
            {"install", "bad-token"},
        )

    def test_parse_main_case_commands_windows_newlines(self) -> None:
        text = "case \"$1\" in\r\ninstall)\r\ncheck)\r\n*)\r\nesac\r\n"
        self.assertEqual(parse_main_case_commands(text), {"install", "check"})

    def test_parse_candidate_quality_options_parses_flags(self) -> None:
        text = """candidate-quality)
  --dry-run
  --force
  --report-json
esac
"""
        self.assertEqual(
            parse_candidate_quality_options(text),
            {"dry-run", "force", "report-json"},
        )

    def test_parse_candidate_quality_options_no_block(self) -> None:
        text = """case \"$1\" in
install)
*)
esac
"""
        self.assertEqual(parse_candidate_quality_options(text), set())

    def test_parse_candidate_quality_options_empty_block(self) -> None:
        text = "candidate-quality)\n"
        self.assertEqual(parse_candidate_quality_options(text), set())

    def test_parse_candidate_quality_options_ignores_short_option(self) -> None:
        text = """candidate-quality)
  -h
  --help
  --output
esac
"""
        self.assertEqual(parse_candidate_quality_options(text), {"help", "output"})

    def test_parse_candidate_quality_options_handles_duplicates(self) -> None:
        text = """candidate-quality)
  --dry-run
  --dry-run
  --out
esac
"""
        self.assertEqual(parse_candidate_quality_options(text), {"dry-run", "out"})

    def test_parse_candidate_quality_options_accepts_mixed_case(self) -> None:
        text = """candidate-quality)
  --Dry-Run
  --OutputDir
esac
"""
        self.assertEqual(parse_candidate_quality_options(text), {"dry-run", "outputdir"})

    def test_unknown_commands_returns_missing_cases(self) -> None:
        self.assertEqual(
            unknown_commands(["install", "check", "launch"], ["install"]),
            ["check", "launch"],
        )

    def test_unknown_commands_returns_empty_when_full_match(self) -> None:
        self.assertEqual(
            unknown_commands(["install", "check"], ["check", "install"]),
            [],
        )

    def test_unknown_commands_ignores_extra_cases(self) -> None:
        self.assertEqual(
            unknown_commands(["install"], ["install", "check", "launch"]),
            [],
        )

    def test_unknown_commands_sorted(self) -> None:
        self.assertEqual(
            unknown_commands(["zeta", "alpha", "gamma"], ["zeta"]),
            ["alpha", "gamma"],
        )

    def test_unknown_commands_with_empty_usage(self) -> None:
        self.assertEqual(unknown_commands([], ["install"]), [])

    def test_unknown_commands_with_empty_cases(self) -> None:
        self.assertEqual(unknown_commands(["install", "check"], []), ["check", "install"])

    def test_unknown_commands_supports_iterators(self) -> None:
        self.assertEqual(
            unknown_commands(("install", "check"), ("install",)),
            ["check"],
        )

    def test_unknown_commands_raises_with_unhashable_members(self) -> None:
        with self.assertRaises(TypeError):
            unknown_commands(["install", ["nested"]], ["install"])

    def test_parser_combo_parses_real_usage_and_cases(self) -> None:
        usage_text = (
            "Usage: scripts/stream-box/ikemen-box.sh "
            "<install|check|launch|launch-menu|status|tail|stop|repoint|bundle|"
            "normalize-chars|normalize-chars-apply|candidate-quality>\n"
        )
        case_text = """case \"$1\" in
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
            unknown_commands(
                parse_usage_commands(usage_text),
                parse_main_case_commands(case_text),
            ),
            [],
        )

    def test_parser_combo_has_unknown(self) -> None:
        usage_text = "Usage: scripts/stream-box/ikemen-box.sh <install|check|launch>\n"
        case_text = """case \"$1\" in
install)
check)
*)
esac
"""
        self.assertEqual(
            unknown_commands(
                parse_usage_commands(usage_text),
                parse_main_case_commands(case_text),
            ),
            ["launch"],
        )

    def test_parser_combo_detects_mismatch(self) -> None:
        usage_text = "Usage: scripts/stream-box/ikemen-box.sh <install|check>\n"
        case_text = "case \"$1\" in\ninstall)\n*)\nesac\n"
        self.assertEqual(
            unknown_commands(
                parse_usage_commands(usage_text),
                parse_main_case_commands(case_text),
            ),
            ["check"],
        )

    def test_parse_candidate_quality_options_with_comments(self) -> None:
        text = """candidate-quality)
# comment
  --dry-run # dry run
  --output-dir # output path
esac
"""
        self.assertEqual(
            parse_candidate_quality_options(text),
            {"dry-run", "output-dir"},
        )


if __name__ == "__main__":
    unittest.main()
