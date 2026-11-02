---
title: Mapped-Frontier Build Packet - Auto-deploy to dev from CI
status: draft
owner: Christopher Richmond
issue: https://github.com/lancer1977/Ikemen-GO/issues/45
created: 2026-11-02
updated: 2026-11-02
---

# Mapped-Frontier Build Packet: Auto-deploy to dev from CI

## Per-project merge policy

<!-- Mirror `.mfd/project.yaml` `merge_policy` for this repository. Keep this
section and that file in sync; do not let them drift. -->

- Required checks: deploy-dev
- Integration branch: develop
- Protected branch: main

- [x] Auto-deploy to dev: dev site deploys from CI/CD on push/merge to dev
  with no manual workflow_dispatch or approval; main/prod deploy needs an
  explicit decision from Christopher. Implemented by
  `.github/workflows/deploy-dev.yml`, refs #45.
