package main

import "testing"

func TestRenderProbeCategoryEnabled(t *testing.T) {
	prev := renderProbeMode
	defer func() { renderProbeMode = prev }()

	renderProbeMode = ""
	if renderProbeCategoryEnabled("hud") {
		t.Fatal("expected disabled mode to reject category")
	}

	renderProbeMode = "all"
	if !renderProbeCategoryEnabled("hud") || !renderProbeCategoryEnabled("") {
		t.Fatal("expected all mode to enable every category")
	}

	renderProbeMode = "hud"
	if !renderProbeCategoryEnabled("hud") {
		t.Fatal("expected exact category match")
	}
	if !renderProbeCategoryEnabled("hud-top") {
		t.Fatal("expected prefix category match")
	}
	if renderProbeCategoryEnabled("input") {
		t.Fatal("expected unrelated category to stay disabled")
	}
}
