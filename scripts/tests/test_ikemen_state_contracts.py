from __future__ import annotations

import unittest
from pathlib import Path


class IkemenStateContractTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        repo_root = Path(__file__).resolve().parents[2]
        cls.compiler_text = (repo_root / "src" / "compiler.go").read_text(encoding="utf-8")
        cls.live_snapshot_text = (repo_root / "src" / "live_snapshot.go").read_text(encoding="utf-8")
        cls.stats_text = (repo_root / "src" / "stats.go").read_text(encoding="utf-8")
        cls.script_text = (repo_root / "src" / "script.go").read_text(encoding="utf-8")

    def test_zss_animtype_parser_accepts_med_shorthand(self) -> None:
        self.assertIn('case "med":', self.compiler_text)
        self.assertIn('ra = RA_Medium', self.compiler_text)
        self.assertIn('ZSS: accept the canonical words plus common shorthand used by older content.', self.compiler_text)

    def test_live_and_final_stats_expose_match_end_state(self) -> None:
        self.assertIn('MatchOver         bool       `json:"matchOver"`', self.live_snapshot_text)
        self.assertIn('Ended      bool     `json:"ended"`', self.stats_text)
        self.assertIn('m.Ended = sys.matchOver()', self.stats_text)
        self.assertIn('MatchOver:         sys.matchOver(),', self.script_text)

    def test_main_help_exposes_the_result_and_snapshot_flags(self) -> None:
        main_text = (Path(__file__).resolve().parents[2] / "src" / "main.go").read_text(encoding="utf-8")
        self.assertIn("-resultfile <jsonfile>  Writes the final match JSON to <jsonfile> (canonical bridge contract)", main_text)
        self.assertIn("-livedatafile <jsonfile> Writes live match snapshots to <jsonfile> during the fight", main_text)
        self.assertIn("-noerrordialog          Logs errors without showing a blocking desktop dialog", main_text)


if __name__ == "__main__":
    unittest.main()
