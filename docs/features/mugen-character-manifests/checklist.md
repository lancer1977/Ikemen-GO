# Checklist

## Discovery

- [x] Identify the shared Ikemen asset root
- [x] Sample real character `.def` files
- [x] Decide on a per-folder `manifest.md` output

## Rule Definition

- [x] Define the manifest schema
- [x] Define the primary `.def` selection heuristic
- [x] Define how missing references are recorded

## Implementation

- [x] Add a manifest generator script
- [x] Add an export-hook invocation
- [x] Update the feature docs

## Validation

- [x] Confirm full-library manifest ownership moved to `Api.Ikemen` / `Ikemen.Core`
- [ ] Generate manifests for the full exported character library
- [ ] Spot-check a few large and oddly named packs
- [ ] Verify the manifests are readable by downstream tooling

## Follow-up

- [x] Decide whether stage folders need the same manifest rule
- [x] Decide whether to create a root index file for all manifests
- [x] Decide whether to add palette/stat extraction beyond `.def` parsing

## Decision Notes

- Root indexing is not added; consumers should walk the per-folder manifests instead.
- Stage manifests remain a separate generator concern outside the character rule.
- Palette/stat extraction stays `.def`-only for this feature.
