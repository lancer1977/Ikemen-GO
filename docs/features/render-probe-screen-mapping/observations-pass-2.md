# Render Probe Screen Mapping - Observation Pass 2

## Context

Pass 2 build was deployed to the local runtime at `/home/lancer1977/code/ikemen-app` and run with:

```bash
IKEMEN_RENDER_PROBES=1 ./Ikemen_GO_Linux
```

Pass 2 changes included:

- A high-visibility `SELECT PATH HIT` sentinel in the local `start.f_selectScreen()` loop.
- Select probes moved to an upper-center vertical stack.
- Global `G*` probes shifted inward from the left edge.
- More granular FightScreen probes around health/life, face/portrait, power, and name draw sections.

Observation source: user visual report after pass 2.

## Select Screen Observation

User report:

> I dont think I saw anything on select

Interpretation:

- `SELECT PATH HIT` was not noticed on the select screen.
- Since `SELECT PATH HIT` is placed at the beginning of the local app's `start.f_selectScreen()` loop, the current observed path likely did not spend visible time in that loop.
- The most likely cause is that the observed flow is still Demo Mode / attract / preselected battle flow, or otherwise bypasses the interactive select screen loop.
- This does not yet prove select drawing cannot be probed; it means the next check must instrument a higher-level menu/mode entry path or force an interactive mode that calls `start.f_selectScreen()`.

Relevant local-app script evidence:

- `start.f_selectScreen()` exists in `/home/lancer1977/code/ikemen-app/external/script/start.lua`.
- It is called from several mode flows, including around the arcade/versus paths.
- The early guard returns before drawing if:

```lua
if (not main.selectMenu[1] and not main.selectMenu[2]) or selScreenEnd then
    return true
end
```

Next select-specific probe should therefore distinguish:

1. Did the mode flow call `start.f_selectScreen()` at all?
2. Did the function return early before the draw loop?
3. Did it enter the `while not selScreenEnd do` loop?

Recommended next select probes:

- `SELECT FUNC ENTER` immediately at the top of `start.f_selectScreen()`.
- `SELECT EARLY RETURN` immediately before the early return.
- Existing `SELECT PATH HIT` inside the while loop.
- Optional mode-level sentinels near each caller, such as arcade/versus selection entrypoints.

## Fight Screen Observation

User report:

> the in play screen had a ton of stuff

Interpretation:

- Pass 2 confirms the in-fight probe density is now high.
- The granular HUD probes are likely firing around the FightScreen subcomponent sections.
- This confirms that FightScreen is the right place to map healthbar/portrait/power/name draw order, but the current pass may be visually too noisy for easy manual reading.

Recommended next fight-specific adjustment:

- Keep the granular probes available, but add a narrower mode/toggle for one HUD family at a time:
  - health only
  - face/portrait only
  - power only
  - name only
- Alternatively, use a single environment variable mode, for example:
  - `IKEMEN_RENDER_PROBES=fight-health`
  - `IKEMEN_RENDER_PROBES=fight-face`
  - `IKEMEN_RENDER_PROBES=fight-power`
  - `IKEMEN_RENDER_PROBES=select`
  - `IKEMEN_RENDER_PROBES=all`

## Pass 2 Conclusions

1. Select probes still have no confirmed visual hit.
2. The current select sentinel is inside `start.f_selectScreen()`'s while loop, so absence points to either not calling that function, an early return, or not reaching/remaining in the loop visibly.
3. Fight-screen probes are very active and may now be too dense for practical manual observation.
4. The next improvement should add probe categories/modes so we can isolate one target at a time.
5. For the design question, the evidence still favors FightScreen for battle-time healthbar/roster-style UI and Lua/motif for screenpack/select-style UI, but select ownership needs a forced interactive-path proof.

## Next Probe Pass

- Add probe mode filtering instead of treating probes as all-or-nothing.
- Instrument select function entry and early return before adding more select draw-layer probes.
- Run with select-only probes first to avoid fight-screen noise:

```bash
IKEMEN_RENDER_PROBES=select ./Ikemen_GO_Linux
```

- Then run focused fight passes:

```bash
IKEMEN_RENDER_PROBES=fight-health ./Ikemen_GO_Linux
IKEMEN_RENDER_PROBES=fight-face ./Ikemen_GO_Linux
IKEMEN_RENDER_PROBES=fight-power ./Ikemen_GO_Linux
```
