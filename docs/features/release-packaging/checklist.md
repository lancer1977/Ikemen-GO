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
- [x] Restore the tracked desktop and Android build toolchain deleted by `0c649765`
- [x] Disable matrix fail-fast so Android failures do not cancel desktop evidence
- [x] Remove Homebrew tap-forcing environment overrides
- [x] Gate manual publication behind an explicit `publishRelease` input
- [x] Sanitize branch names before using them in archive paths

## Validation

- [x] Run the Linux packaging flow end-to-end
- [ ] Run the Windows packaging flow on a Windows/MSYS2 host
- [x] Confirm the Linux archive appears in the writable local output directory used for validation
- [x] Confirm the final Linux archive appears in the shared output directory
- [x] Fail fast on Windows packaging from the wrong host
- [ ] Validate the new GitHub Actions workflow on a self-hosted runner
- [x] Validate the stream-box test deploy root at `C:\\mugen`

## Follow-up

- [x] Prove the restored Android artifact lane on GitHub Actions
- [x] Prove Linux and macOS desktop artifact lanes in the same diagnostic run

- [x] Decide whether `make installers` should become the primary entry point
- [x] Decide whether the script should fail hard on missing screenpack assets or keep producing a partial package
