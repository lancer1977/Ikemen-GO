# Render Probe Screen Mapping - Pass 4 Swarm Plan

## Purpose

User asked to "swarm decorate" the remaining unexposed layer/draw types with unique identifiers, then visually report which codes appear. This is intended to be the final broad select-screen ownership pass.

## Run Mode

```bash
IKEMEN_RENDER_PROBES=swarm ./Ikemen_GO_Linux
```

`swarm` enables categories prefixed with `swarm-`, including:

- `swarm-select`
- `swarm-batch`

## Added Select/Lua Codes

Function and list construction:

- `SW01 updateDrawList enter`
- `SW02 updateDrawList return`

Portrait drawing:

- `SW10 portraits enter p1/p2`
- `SW11 face2 pre p1/p2`
- `SW12 face2 post p1/p2`
- `SW13 face pre p1/p2`
- `SW14 face post p1/p2`
- `SW15 icons pre p1/p2`
- `SW16 portraits exit p1/p2`

Main select loop blocks:

- `SW20 portraits block pre`
- `SW21 portraits block post`
- `SW30 batchDraw pre`
- `SW31 batchDraw post`
- `SW40 done cursors pre`
- `SW41 done cursors post`
- `SW50 team/select menu pre`
- `SW51 team/select menu post`
- `SW60 names pre`
- `SW61 names post`
- `SW70 complete block`
- `SW71 stage block pre`
- `SW72 record block pre`
- `SW80 timer pre`
- `SW81 timer post`
- `SW90 hook pre`
- `SW91 hook post`
- `SW92 bg1 pre`
- `SW93 bg1/stats post final` or `SW93 bg1 post/final`
- `SW99 before refresh`

Go batch draw layer callback:

- `SWB batch L<layer>`

## What To Report

For the select screen, report any visible `SW*` / `SWB*` codes and where they appear relative to:

- 3D-ish/perspective roster cells
- big portraits / secondary angled portraits
- active cursors
- names
- team menu
- top/background layer

If no `SW*` codes appear on select, the visible select surface is likely either clearing/covering all ordinary debug text or using a path outside the instrumented Lua select loop. The next move would be lower-level renderer or debug-text-over-present probing.

## Future Greenscreen Ask

User also wants a greenscreen background at some point. Keep this as a later probe/feature option: add a gated mode that clears the background to chroma green or swaps select/fight background layers for a green fill without affecting character/HUD rendering.
