from __future__ import annotations

import re
import unittest
from pathlib import Path

import yaml

# GitHub Actions normalizes the "on:" key to True (YAML 1.1 boolean) when parsed by PyYAML.
ON_KEY = True

# Actions maintained by the "actions" GitHub org (first-party, github-script included)
# are exempt from the SHA-pin requirement below only when referenced by a full commit SHA
# anyway -- this repo still requires them pinned, but the allow-list documents intent.
FIRST_PARTY_OWNER = "actions"

FULL_SHA_RE = re.compile(r"^[0-9a-f]{40}$")


class PullRequestTargetWorkflowAuthorityTests(unittest.TestCase):
    """
    Static contract for .github/workflows/conventional-title.yml and
    conventional-label.yml (GitHub issue #9): both run on pull_request_target,
    which grants base-repo token context including secrets. They must never
    check out or execute pull-request source, must not hold broad
    (write-all) authority, and any third-party action they use must be
    pinned to an immutable commit SHA.
    """

    WORKFLOW_NAMES = ("conventional-title.yml", "conventional-label.yml")

    @classmethod
    def setUpClass(cls) -> None:
        repo_root = Path(__file__).resolve().parents[2]
        workflows_dir = repo_root / ".github" / "workflows"
        cls.workflows: dict[str, dict] = {}
        cls.raw_text: dict[str, str] = {}
        for name in cls.WORKFLOW_NAMES:
            path = workflows_dir / name
            text = path.read_text(encoding="utf-8")
            cls.raw_text[name] = text
            cls.workflows[name] = yaml.safe_load(text)

    def _all_steps(self, workflow: dict):
        for job in workflow.get("jobs", {}).values():
            for step in job.get("steps", []) or []:
                yield step

    def test_workflows_still_trigger_on_pull_request_target(self) -> None:
        # This is deliberate: pull_request_target is required for the
        # metadata-mutation (comment/label) purpose. It must not be swapped
        # for a plain pull_request trigger, which would lack write context.
        for name, workflow in self.workflows.items():
            with self.subTest(workflow=name):
                triggers = workflow[ON_KEY]
                self.assertIn("pull_request_target", triggers)

    def test_no_write_all_permissions(self) -> None:
        for name, workflow in self.workflows.items():
            with self.subTest(workflow=name):
                permissions = workflow.get("permissions")
                self.assertNotEqual(permissions, "write-all")
                self.assertIsInstance(
                    permissions,
                    dict,
                    f"{name}: permissions must be an explicit minimal mapping, not a bare string",
                )
                self.assertNotIn("write-all", permissions.values())

    def test_permissions_are_minimal_pull_requests_write_only(self) -> None:
        for name, workflow in self.workflows.items():
            with self.subTest(workflow=name):
                permissions = workflow.get("permissions", {})
                self.assertEqual(
                    permissions,
                    {"pull-requests": "write"},
                    f"{name}: expected the minimum scope (pull-requests: write) and nothing broader",
                )

    def test_no_checkout_step_of_any_kind(self) -> None:
        # Neither workflow needs repo contents at all -- they only read PR
        # title/body/number via `context.payload.pull_request` and post a
        # review or label. A checkout step (of PR head or otherwise) would
        # reintroduce the ability to execute pull-request-controlled code
        # under base-repo token/secret context.
        for name, workflow in self.workflows.items():
            with self.subTest(workflow=name):
                for step in self._all_steps(workflow):
                    uses = step.get("uses", "")
                    self.assertNotIn(
                        "actions/checkout",
                        uses,
                        f"{name}: pull_request_target workflow must never check out PR source",
                    )

    def test_no_reference_to_pull_request_head_anywhere(self) -> None:
        for name, text in self.raw_text.items():
            with self.subTest(workflow=name):
                self.assertNotIn("event.pull_request.head", text)
                self.assertNotIn("pull_request.head.sha", text)
                self.assertNotIn("pull_request.head.ref", text)

    def test_third_party_actions_are_pinned_to_a_full_commit_sha(self) -> None:
        for name, workflow in self.workflows.items():
            with self.subTest(workflow=name):
                for step in self._all_steps(workflow):
                    uses = step.get("uses")
                    if not uses:
                        continue
                    self.assertIn(
                        "@",
                        uses,
                        f"{name}: {uses!r} must be pinned to a ref",
                    )
                    action_ref, _, pinned = uses.partition("@")
                    self.assertTrue(
                        FULL_SHA_RE.match(pinned),
                        f"{name}: {uses!r} must be pinned to a full 40-character commit SHA, not a floating tag",
                    )

    def test_only_first_party_github_script_action_is_used(self) -> None:
        # The prior implementation depended on morrisoncole/pr-lint-action and
        # bcoe/conventional-release-labels, two unpinned third-party actions
        # running with base-repo secrets. They have been replaced outright
        # with actions/github-script, removing the third-party supply-chain
        # surface entirely rather than just pinning it.
        for name, workflow in self.workflows.items():
            with self.subTest(workflow=name):
                action_owners = set()
                for step in self._all_steps(workflow):
                    uses = step.get("uses")
                    if not uses:
                        continue
                    action_ref = uses.split("@", 1)[0]
                    owner = action_ref.split("/", 1)[0]
                    action_owners.add(owner)
                self.assertEqual(
                    action_owners,
                    {FIRST_PARTY_OWNER},
                    f"{name}: only actions/* (github-script) steps are expected here",
                )

    def test_failure_path_is_explicit_not_silent(self) -> None:
        # Each script step must call core.setFailed on an unexpected error
        # instead of allowing the step -- and therefore the check -- to
        # report green while doing nothing (fail-open).
        for name, text in self.raw_text.items():
            with self.subTest(workflow=name):
                self.assertIn("core.setFailed", text)
                self.assertIn("catch (error)", text)


if __name__ == "__main__":
    unittest.main()
