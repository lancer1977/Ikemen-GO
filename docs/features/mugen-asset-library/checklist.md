# Checklist

## Discovery

- [x] Confirm current Ikemen asset root
- [x] Inspect the source Mugen library layout
- [x] Identify existing deploy-aligned test-assets content

## Implementation

- [x] Add a reusable export script
- [x] Export top-level character directories
- [x] Export character-split run directories
- [x] Export stage files from all stage trees
- [x] Add root convenience links for `chars` and `stages`

## Validation

- [ ] Run the exporter against the full Mugen tree
- [ ] Spot-check that missing chars and stages now exist in the Ikemen asset store
- [ ] Confirm deploy-linux still sees the same asset layout
- [ ] Generate standardized character manifests for the shared library

## Follow-up

- [ ] Decide whether the exporter should also generate manifests by origin tag
- [ ] Decide whether to auto-normalize folder names during export
