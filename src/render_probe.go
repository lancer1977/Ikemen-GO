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

	r = Clamp(r, 0, 255)
	g = Clamp(g, 0, 255)
	b = Clamp(b, 0, 255)
	if w < 12 {
		w = 12
	}
	if h < 12 {
		h = 12
	}

	FillRect([4]int32{int32(x), int32(y), w, h}, uint32((r<<16)|(g<<8)|b), [2]int32{240, 0}, nil)
	sys.debugFont.SetColor(255-r/2, 255-g/2, 255-b/2, 255)
	sys.debugFont.fnt.Print(label, x+4, y+float32(h/2), sys.debugFont.xscl/sys.widthScale,
		sys.debugFont.yscl/sys.heightScale, 0, Rotation{0, 0, 0}, 0, 0, 0, 0, &sys.scrrect,
		sys.debugFont.palfx, sys.debugFont.frgba)
}
