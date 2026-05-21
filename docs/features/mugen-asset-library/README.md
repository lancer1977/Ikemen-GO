# Mugen Asset Library

## Summary

This feature tracks the shared asset store used for Ikemen GO validation and
test deploys. The canonical storage root is:

`/mnt/syn1/games/Ikemen/test-assets`

The top-level `chars` and `stages` paths under `/mnt/syn1/games/Ikemen` are
kept as convenience links to the same asset store.

## Current State

- [x] Shared asset root identified
- [x] Export script added for Mugen character and stage gaps
- [x] Root convenience links created for `chars` and `stages`
- [ ] Validate the full export against the current Mugen library
- [ ] Decide whether to split the library into per-pack manifests
- [x] Decide whether to add metadata files for unique character identity and origin tags

The manifest follow-up now lives in the dedicated
[Mugen Character Manifests](../mugen-character-manifests/README.md) feature.
Canonical asset analysis and manifest tooling now lives in `Api.Ikemen/scripts/asset-normalization/`.

## Source Roots

- `/mnt/syn1/games/Mugen/Characters`
- `/mnt/syn1/games/Mugen/Characters-split/runs`
- `/mnt/syn1/games/Mugen/**/stages`

## Notes

- The exporter preserves existing files and only fills gaps.
- The stage export is filename-based so shared stage packs do not overwrite
  existing entries.
- The character export copies each top-level character directory into the
  shared asset store without clobbering existing contents.
- The exporter can also call
  `Api.Ikemen/scripts/asset-normalization/generate-mugen-character-manifests.py`
  for standardized per-character `manifest.md` files.
