import json
import re
import subprocess
import tempfile
import unittest
from pathlib import Path


class DeployIkemenDevContractsTests(unittest.TestCase):
    def test_install_helper_covers_check_deps_package_hints(self) -> None:
        repo_root = Path(__file__).resolve().parents[2]
        deploy_text = (repo_root / "scripts" / "deploy-ikemen-dev.sh").read_text(encoding="utf-8")
        helper_text = (repo_root / "scripts" / "install-linux-build-deps.sh").read_text(encoding="utf-8")

        hinted_packages = set(re.findall(r'check_(?:command|pkg_config) [^"\n]+ "([^"]+)"', deploy_text))
        helper_packages = set(re.findall(r"^\s+([A-Za-z0-9+_.-]+)\s*\\?$", helper_text, flags=re.MULTILINE))

        for package in hinted_packages - {"bash"}:
            with self.subTest(package=package):
                self.assertIn(package, helper_packages)

    def test_check_deps_install_hint_matches_reported_missing_dependencies(self) -> None:
        repo_root = Path(__file__).resolve().parents[2]
        with tempfile.TemporaryDirectory(prefix="ikemen-dev-check-") as tmp:
            app_root = Path(tmp) / "ikemen-dev"

            result = subprocess.run(
                [
                    "bash",
                    str(repo_root / "scripts" / "deploy-ikemen-dev.sh"),
                    "--check-deps",
                    "--app-root",
                    str(app_root),
                ],
                cwd=repo_root,
                text=True,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                check=True,
            )

            missing_packages = {
                parts[4]
                for line in result.stdout.splitlines()
                if (parts := line.split("|")) and len(parts) == 5
                and parts[0] == "dependency"
                and parts[1] in {"command", "pkg-config"}
                and parts[3] == "missing"
            }
            install_line = next(
                line for line in result.stdout.splitlines()
                if line.startswith("preflight|installMissing|")
            )
            if missing_packages:
                self.assertTrue(install_line.startswith("preflight|installMissing|sudo apt-get install -y "))
                for package in missing_packages:
                    self.assertIn(package, install_line)
            else:
                self.assertEqual(install_line, "preflight|installMissing|none")

    def test_check_deps_reports_missing_manifest_for_empty_app_root(self) -> None:
        repo_root = Path(__file__).resolve().parents[2]
        with tempfile.TemporaryDirectory(prefix="ikemen-dev-check-") as tmp:
            app_root = Path(tmp) / "ikemen-dev"

            result = subprocess.run(
                [
                    "bash",
                    str(repo_root / "scripts" / "deploy-ikemen-dev.sh"),
                    "--check-deps",
                    "--app-root",
                    str(app_root),
                ],
                cwd=repo_root,
                text=True,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                check=True,
            )

            self.assertIn(f"preflight|appRoot|{app_root}", result.stdout)
            self.assertIn(f"dependency|file|deployed-binary|missing|{app_root / 'Ikemen_GO_Linux'}", result.stdout)
            self.assertIn(f"dependency|file|deployed-manifest|missing|{app_root / 'deploy-manifest.json'}", result.stdout)
            self.assertIn(f"preflight|deployManifestHash|missing|{app_root / 'deploy-manifest.json'}", result.stdout)

    def test_check_deps_reports_manifest_hash_mismatch(self) -> None:
        repo_root = Path(__file__).resolve().parents[2]
        with tempfile.TemporaryDirectory(prefix="ikemen-dev-check-") as tmp:
            app_root = Path(tmp) / "ikemen-dev"
            app_root.mkdir(parents=True)
            binary = app_root / "Ikemen_GO_Linux"
            binary.write_bytes(b"new-binary")
            binary.chmod(0o755)
            (app_root / "deploy-manifest.json").write_text(
                json.dumps(
                    {
                        "schema": "ikemen-go/dev-deploy-manifest/v1",
                        "binary": {"sha256": "old-hash"},
                    }
                )
                + "\n",
                encoding="utf-8",
            )

            result = subprocess.run(
                [
                    "bash",
                    str(repo_root / "scripts" / "deploy-ikemen-dev.sh"),
                    "--check-deps",
                    "--app-root",
                    str(app_root),
                ],
                cwd=repo_root,
                text=True,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                check=True,
            )

            self.assertIn("preflight|deployManifestHash|mismatch|expected=old-hash|actual=", result.stdout)

    def test_no_build_deploy_manifest_records_final_binary_path(self) -> None:
        repo_root = Path(__file__).resolve().parents[2]
        if not (repo_root / "Ikemen_GO_Linux").is_file():
            self.skipTest("Ikemen_GO_Linux is required for --no-build deploy contract")

        with tempfile.TemporaryDirectory(prefix="ikemen-dev-deploy-") as tmp:
            root = Path(tmp)
            app_root = root / "ikemen-dev"
            screenpack = root / "screenpack"
            (screenpack / "chars" / "kfm").mkdir(parents=True)
            (screenpack / "stages").mkdir(parents=True)
            (screenpack / "data" / "ikemen1").mkdir(parents=True)
            (screenpack / "chars" / "kfm" / "kfm.def").write_text("[Info]\nname = kfm\n", encoding="utf-8")
            (screenpack / "stages" / "kfm.def").write_text("[Info]\nname = kfm\n", encoding="utf-8")
            (screenpack / "data" / "ikemen1" / "system.def").write_text("[Info]\nname = ikemen1\n", encoding="utf-8")

            subprocess.run(
                [
                    "bash",
                    str(repo_root / "scripts" / "deploy-ikemen-dev.sh"),
                    "--no-build",
                    "--app-root",
                    str(app_root),
                    "--screenpack-dir",
                    str(screenpack),
                ],
                cwd=repo_root,
                text=True,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                check=True,
            )

            manifest = json.loads((app_root / "deploy-manifest.json").read_text(encoding="utf-8"))
            self.assertEqual(manifest["appRoot"], str(app_root))
            self.assertEqual(manifest["binary"]["path"], str(app_root / "Ikemen_GO_Linux"))
            self.assertNotIn(".tmp.", manifest["binary"]["path"])


if __name__ == "__main__":
    unittest.main()
