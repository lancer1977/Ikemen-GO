# Mugen Character Manifests

## Summary

This feature defines a standardized markdown manifest for each exported Mugen
or Ikemen character folder. The manifest is generated from the primary `.def`
file in the folder and lives beside the character files as `manifest.md`.

The goal is to make character packs easier to scan by humans, scripts, and AI
tools without needing to inspect every asset manually.

## Current State

- [x] Manifest rule defined
- [x] Generator script added and migrated to `Api.Ikemen`
- [x] Export pipeline can invoke manifest generation
- [ ] Backfill manifests for the full shared character library
- [ ] Decide whether stages should get a separate manifest rule

## Manifest Contract

- `schema: mugen-character-manifest/v1`
- `manifest.md` is written inside the character folder
- the primary `.def` file is used as the source of truth
- the manifest includes:
  - folder identity
  - display name
  - internal name
  - author
  - version date
  - Mugen version
  - palette defaults when present
  - referenced file list from the `.def`
  - top-level file inventory
  - missing file references when applicable

## Output Location

- `/mnt/syn1/games/Ikemen/test-assets/chars/<character-folder>/manifest.md`

## Canonical Tooling

- Generator: `Api.Ikemen/scripts/asset-normalization/generate-mugen-character-manifests.py`
- Default source root: `/mnt/syn1/games/Mugen/Characters`
- Use `--dry-run` to inspect planned writes without mutating character folders.

## Notes

- The manifest is intentionally markdown first, with YAML-style front matter for
  machine consumption.
- Missing references are preserved instead of being hidden.
- The output is designed to be easy to index from the shared Ikemen asset root.
- `identity_base` and `identity_slug` normalize pack labels out of the folder
  name so AI and scripts can group the same character across variants.
