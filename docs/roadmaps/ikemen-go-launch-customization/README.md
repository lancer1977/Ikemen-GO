# Ikemen GO Launch Customization Roadmap

## Summary

This roadmap tracks the next engine-level launch work that sits beyond the
current CLI contract.

The immediate result-file path is now in place. The remaining work is about
making higher-level launch orchestration easier to express without growing the
flag surface indefinitely.

## Scope

- structured launch payloads
- live match-event transport
- automation-only exit modes
- future bridge/runtime negotiation

## Current State

- [x] The repo has a documented CLI launch contract
- [x] Match results can be written to a unique result file
- [x] Local stdout JSON remains available for debugging
- [ ] A structured launch payload has been designed
- [ ] Live event transport has been prototyped

