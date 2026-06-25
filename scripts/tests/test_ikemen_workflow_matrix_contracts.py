from __future__ import annotations

import unittest
from pathlib import Path


class IkemenWorkflowMatrixContractTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        cls.script_path = (
            Path(__file__).resolve().parents[1] / "smoke" / "ikemen-workflow-matrix.sh"
        )
        cls.script_text = cls.script_path.read_text(encoding="utf-8")

    def test_usage_exposes_the_extended_matrix_flags(self) -> None:
        self.assertIn("--char-sweep-watchdog|--no-char-sweep-watchdog", self.script_text)
        self.assertIn("--char-sweep-long-intro-timeout-multiplier N", self.script_text)
        self.assertIn("--case NAME", self.script_text)
        self.assertIn("--fixture-root PATH", self.script_text)
        self.assertIn("--char-dir PATH", self.script_text)

    def test_case_parser_recognizes_the_extended_matrix_flags(self) -> None:
        self.assertIn("--char-sweep-watchdog)", self.script_text)
        self.assertIn("--no-char-sweep-watchdog)", self.script_text)
        self.assertIn("--char-sweep-long-intro-timeout-multiplier)", self.script_text)
        self.assertIn("--fixture-root)", self.script_text)
        self.assertIn("--case)", self.script_text)
        self.assertIn("--char-dir)", self.script_text)

    def test_run_case_fails_timeout_and_nonzero_exit_codes(self) -> None:
        self.assertIn('FAIL: $name hit timeout (${IKEMEN_WORKFLOW_TIMEOUT}s)', self.script_text)
        self.assertIn('FAIL: $name exited with status $run_status', self.script_text)
        self.assertNotIn('OK: $name exited with status $run_status but logged no crash text', self.script_text)

    def test_result_case_keeps_match_end_state_in_summary(self) -> None:
        self.assertIn("fightEnded=$fight_ended", self.script_text)
        self.assertIn("liveMatchOver=$live_match_over", self.script_text)

    def test_bootstrap_requires_the_default_window_icon_assets(self) -> None:
        self.assertIn("ensure_window_icon_assets", self.script_text)
        self.assertIn("external/icons/IkemenCylia_256.png", self.script_text)
        self.assertIn("The default config expects the Cylia icon set under external/icons/.", self.script_text)


class StreamBoxToolboxContractTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        cls.toolbox_path = (
            Path(__file__).resolve().parents[1] / "stream-box" / "ikemen-box.sh"
        )
        cls.toolbox_text = cls.toolbox_path.read_text(encoding="utf-8")

    def test_check_deploy_requires_the_default_window_icon_assets(self) -> None:
        self.assertIn('external\\icons\\IkemenCylia_256.png', self.toolbox_text)
        self.assertIn('external\\icons\\IkemenCylia_96.png', self.toolbox_text)
        self.assertIn('external\\icons\\IkemenCylia_48.png', self.toolbox_text)
        self.assertIn('Missing required icon asset:', self.toolbox_text)


if __name__ == "__main__":
    unittest.main()
