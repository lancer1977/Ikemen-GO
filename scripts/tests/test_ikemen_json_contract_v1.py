from __future__ import annotations

import re
import unittest
from pathlib import Path


class IkemenJsonContractV1Tests(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        cls.repo_root = Path(__file__).resolve().parents[2]
        cls.stats_text = (cls.repo_root / "src" / "stats.go").read_text(encoding="utf-8")
        cls.live_text = (cls.repo_root / "src" / "live_snapshot.go").read_text(encoding="utf-8")
        cls.script_text = (cls.repo_root / "src" / "script.go").read_text(encoding="utf-8")
        cls.doc_text = (
            cls.repo_root / "docs" / "features" / "ikemen-go" / "result-contract-v1.md"
        ).read_text(encoding="utf-8")
        cls.contract_dir = cls.repo_root / "contracts" / "Ikemen.Go.Contracts"
        cls.contract_text = "\n".join(
            path.read_text(encoding="utf-8")
            for path in sorted(cls.contract_dir.glob("*.cs"))
        )

    def assert_contract_fields(self, fields: set[str]) -> None:
        for field in sorted(fields):
            with self.subTest(field=field):
                self.assertIn(f'[JsonPropertyName("{field}")]', self.contract_text)
                self.assertIn(f"`{field}`", self.doc_text)

    def test_csharp_contract_covers_final_result_json_tags(self) -> None:
        tags = set(re.findall(r'json:"([^",]+)"', self.stats_text))
        tags.update({"statsLog", "continueFlg", "persistRoundCount", "matchOver"})

        self.assert_contract_fields(tags)

    def test_csharp_contract_covers_live_snapshot_json_tags(self) -> None:
        tags = set(re.findall(r'json:"([^",]+)"', self.live_text))

        self.assert_contract_fields(tags)

    def test_v1_contract_documents_current_non_event_surface(self) -> None:
        self.assertIn("V1 is not an append-only event stream.", self.doc_text)
        self.assertIn("V1 does not emit discrete damage", self.doc_text)
        self.assertIn("GameStatsSnapshotV1", self.contract_text)
        self.assertIn("GameLiveSnapshotV1", self.contract_text)
        self.assertIn('public const string Name = "IkemenGo.MatchState.V1";', self.contract_text)

    def test_engine_json_export_still_uses_game_stats_snapshot(self) -> None:
        self.assertIn("GameStatsSnapshot{", self.script_text)
        self.assertIn("StatsLog:          sys.statsLog,", self.script_text)
        self.assertIn("MatchOver:         sys.matchOver(),", self.script_text)


if __name__ == "__main__":
    unittest.main()
