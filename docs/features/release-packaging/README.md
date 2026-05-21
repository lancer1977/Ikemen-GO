# Release Packaging

## Summary

This feature covers local generation of desktop release archives for Ikemen GO.
The current packaging target is the shared output root at `/mnt/syn1/games/Ikemen`.

## Current State

- [x] Release workflow inspected for the canonical Windows and Linux archive layout
- [x] Local packaging script added for desktop installers
- [x] Host-aware deploy wrapper added for local publish
- [ ] Windows packaging validated on a Windows or MSYS2 host
- [x] Linux packaging validated end-to-end on this host
- [x] Local package includes `external/script/main.lua` and the rest of the runtime tree
- [x] Packaging flow wired into `make installers` and `make deploy-local`
- [x] Explicit `make deploy-local-linux` and `make deploy-local-windows` targets added
- [x] GitHub Actions workflow added for self-hosted local deploy
- [ ] Windows GitHub Actions runs on the self-hosted runner registered as `ikemen-windows`
- [x] Stream rig validation uses the `stream-box` SSH alias instead of raw IPs
- [x] The Windows test deploy root is `C:\\mugen` on the stream box

## Output Layout

- `Ikemen_GO-dev-linux.zip`
- `Ikemen_GO-dev-windows.zip`

Each archive mirrors the release workflow shape:

- desktop binary
- runtime libraries
- runtime assets
- README
- license notices

## Notes

- Linux packages also include `Ikemen_GO.desktop` and `Ikemen_GO.command`
- Windows packaging requires the MinGW toolchain that matches the existing release build
- Screenpack assets are pulled from the upstream screenpack repository when they are not already present locally
- `./scripts/deploy-local.sh` now dispatches to the correct platform target instead of trying both on every host
- Use `ssh stream-box` or `scripts/stream-box-ssh.sh --check` when validating
  the Windows stream rig from the repo or from deploy automation.
- Copy the packaged Ikemen runtime into `C:\\mugen` for test deploys.
