package main

import (
	"os"
	"strings"
)

var renderProbeMode = parseRenderProbeMode()

func parseRenderProbeMode() string {
	v := strings.TrimSpace(strings.ToLower(os.Getenv("IKEMEN_RENDER_PROBES")))
	switch v {
	case "", "0", "false", "no", "off":
		return ""
	case "1", "true", "yes", "on":
		return "all"
	default:
		return v
	}
}

func renderProbeCategoryEnabled(category string) bool {
	if renderProbeMode == "" {
		return false
	}
	category = strings.TrimSpace(strings.ToLower(category))
	if category == "" {
		category = "all"
	}
	if renderProbeMode == "all" || renderProbeMode == category {
		return true
	}
	if strings.HasPrefix(category, renderProbeMode+"-") || strings.HasPrefix(renderProbeMode, category+"-") {
		return true
	}
	return false
}

func drawRenderProbe(label string, x, y float32, r, g, b int32) {
	drawRenderProbeMode("all", label, x, y, r, g, b)
}

func drawRenderProbeMode(category, label string, x, y float32, r, g, b int32) {
	if !renderProbeCategoryEnabled(category) || sys.debugFont == nil || sys.debugFont.fnt == nil || sys.frameSkip {
		return
	}

	r = Clamp(r, 0, 255)
	g = Clamp(g, 0, 255)
	b = Clamp(b, 0, 255)
	alpha := [2]int32{220, 0}
	w := int32(8*len(label) + 8)
	if w < 32 {
		w = 32
	}

	FillRect([4]int32{int32(x) - 2, int32(y) - 10, w, 14}, uint32((r<<16)|(g<<8)|b), alpha, nil)
	sys.debugFont.SetColor(255-r/2, 255-g/2, 255-b/2, 255)
	sys.debugFont.fnt.Print(label, x, y, sys.debugFont.xscl/sys.widthScale,
		sys.debugFont.yscl/sys.heightScale, 0, Rotation{0, 0, 0}, 0, 0, 0, 0, &sys.scrrect,
		sys.debugFont.palfx, sys.debugFont.frgba)
}

func drawRenderProbeBlockMode(category, label string, x, y float32, w, h int32, r, g, b int32) {
	if !renderProbeCategoryEnabled(category) || sys.debugFont == nil || sys.debugFont.fnt == nil || sys.frameSkip {
		return
	}

	drawRenderProbeBlock(label, x, y, w, h, r, g, b, [2]int32{240, 0})
}

func drawRenderProbeScreenBlockMode(category, label string, x, y float32, w, h int32, r, g, b int32) {
	if !renderProbeCategoryEnabled(category) || sys.debugFont == nil || sys.debugFont.fnt == nil || sys.frameSkip {
		return
	}

	// Queue in the top Lua layer so screenpack/bg/text draws flushed by refresh()
	// do not overwrite the probe. This is for top/bottom HUD seam probes.
	labelLocal := label
	xLocal := x
	yLocal := y
	wLocal := w
	hLocal := h
	rLocal := r
	gLocal := g
	bLocal := b
	sys.luaQueueLayerDraw(2, func() {
		maxW := sys.scrrect[2]
		maxH := sys.scrrect[3]
		xi := int32(xLocal)
		yi := int32(yLocal)
		if xi < 0 {
			xi = maxW + xi
		}
		if yi < 0 {
			yi = maxH + yi
		}
		xi = Clamp(xi, 0, maxW)
		yi = Clamp(yi, 0, maxH)
		if wLocal <= 0 {
			wLocal = maxW - xi
		}
		if hLocal <= 0 {
			hLocal = maxH - yi
		}
		if xi+wLocal > maxW {
			wLocal = maxW - xi
		}
		if yi+hLocal > maxH {
			hLocal = maxH - yi
		}
		if wLocal < 12 || hLocal < 12 {
			return
		}

		drawRenderProbeBlock(labelLocal, float32(xi), float32(yi), wLocal, hLocal, rLocal, gLocal, bLocal, [2]int32{255, 0})
	})
}

func drawRenderProbeBlock(label string, x, y float32, w, h int32, r, g, b int32, alpha [2]int32) {
	r = Clamp(r, 0, 255)
	g = Clamp(g, 0, 255)
	b = Clamp(b, 0, 255)

	FillRect([4]int32{int32(x), int32(y), w, h}, uint32((r<<16)|(g<<8)|b), alpha, nil)
	sys.debugFont.SetColor(255-r/2, 255-g/2, 255-b/2, 255)
	sys.debugFont.fnt.Print(label, x+4, y+float32(h/2), sys.debugFont.xscl/sys.widthScale,
		sys.debugFont.yscl/sys.heightScale, 0, Rotation{0, 0, 0}, 0, 0, 0, 0, &sys.scrrect,
		sys.debugFont.palfx, sys.debugFont.frgba)
}
