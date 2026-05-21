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

- [ ] Generate manifests for the full exported character library
- [ ] Spot-check a few large and oddly named packs
- [ ] Verify the manifests are readable by downstream tooling

## Follow-up

- [ ] Decide whether stage folders need the same manifest rule
- [ ] Decide whether to create a root index file for all manifests
- [ ] Decide whether to add palette/stat extraction beyond `.def` parsing
