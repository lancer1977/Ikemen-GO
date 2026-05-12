# Checklist

## Discovery

- [x] Read the existing build and release workflow
- [x] Confirm the release artifact structure for desktop zips
- [x] Check whether the runtime screenpack assets are present in the local checkout

## Implementation

- [x] Add a local packaging script
- [x] Keep output rooted at `/mnt/syn1/games/Ikemen` by default
- [x] Preserve the release workflow archive names for Linux and Windows
- [x] Include runtime assets and license notices in the staged package
- [x] Add a host-aware deploy wrapper for local publish
- [x] Add direct Linux and Windows deploy wrappers
- [x] Wire local publish into `make installers` and `make deploy-local`
- [x] Add explicit `make deploy-local-linux` and `make deploy-local-windows` targets
- [x] Align local staging with the release workflow runtime tree
- [x] Add GitHub Actions workflow for self-hosted deploy

## Validation

- [x] Run the Linux packaging flow end-to-end
- [ ] Run the Windows packaging flow on a Windows/MSYS2 host
- [x] Confirm the Linux archive appears in the writable local output directory used for validation
- [x] Confirm the final Linux archive appears in the shared output directory
- [x] Fail fast on Windows packaging from the wrong host
- [ ] Validate the new GitHub Actions workflow on a self-hosted runner

## Follow-up

- [ ] Decide whether `make installers` should become the primary entry point
- [ ] Decide whether the script should fail hard on missing screenpack assets or keep producing a partial package
