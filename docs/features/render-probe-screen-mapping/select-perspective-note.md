# Select Screen Perspective Note

User observed that the current select screen has a newer 3D-ish look and still shows no select probe labels in `IKEMEN_RENDER_PROBES=select` mode.

## Evidence

The local runtime screenpack at `/home/lancer1977/code/ikemen-app/data/ikemen1/system.def` uses perspective-projected select elements:

- `[Select Info]` defines the roster grid.
- `cell.*-0` through `cell.*-9` set `projection = perspective`.
- Those same rows also set `xangle`, `focallength`, `spacing`, `scale`, and `offset`.
- `p1.face2` and `p2.face2` also use `projection = perspective` and `yangle` for angled secondary portraits.

So the 3D-ish look is not necessarily a wholly separate renderer. It is likely the regular Lua/select screenpack path using Ikemen's perspective projection support for animations and batched cell drawing.

## Current implication

The absence of select probes is no longer best explained as "new renderer" first. More likely candidates are:

1. The visible select screen is drawn by the same Lua/screenpack path, but the current probe placement/order is being covered or cleared by projected/layered select elements.
2. The observed screen is a different select-like menu path than `start.f_selectScreen()`.
3. The select function is reached but our immediate function-entry probes are drawn before a later clear/background pass and never survive to the presented frame.

## Next useful probe

Add a probe after the final visible select draw calls, especially after projected cells/faces and any fade draw, or add instrumentation inside `batchDraw` / animation projection draw to prove the perspective cell path itself.
