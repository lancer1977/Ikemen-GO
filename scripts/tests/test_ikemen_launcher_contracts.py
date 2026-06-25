from __future__ import annotations

import unittest
from pathlib import Path


class IkemenLocalKfmSmokeContractTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        cls.script_text = (
            Path(__file__).resolve().parents[1]
            / "smoke"
            / "ikemen-local-kfm-smoke.sh"
        ).read_text(encoding="utf-8")

    def test_local_kfm_smoke_requires_the_default_window_icon_assets(self) -> None:
        self.assertIn("ensure_icon_assets", self.script_text)
        self.assertIn("external/icons/IkemenCylia_256.png", self.script_text)
        self.assertIn("The default config expects the Cylia icon set under external/icons/.", self.script_text)


class IkemenVisualSnapshotContractTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        cls.script_text = (
            Path(__file__).resolve().parents[1]
            / "visual"
            / "ikemen-control-snapshots.py"
        ).read_text(encoding="utf-8")

    def test_visual_helper_requires_the_default_window_icon_assets(self) -> None:
        self.assertIn("require_icon_assets", self.script_text)
        self.assertIn("IkemenCylia_256.png", self.script_text)
        self.assertIn("The default config expects the Cylia icon set under external/icons/.", self.script_text)


if __name__ == "__main__":
    unittest.main()
