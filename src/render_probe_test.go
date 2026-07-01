package main

import (
	"os"
	"testing"
)

func TestParseRenderProbeMode(t *testing.T) {
	t.Parallel()

	prev, had := os.LookupEnv("IKEMEN_RENDER_PROBES")
	restore := func() {
		if had {
			_ = os.Setenv("IKEMEN_RENDER_PROBES", prev)
		} else {
			_ = os.Unsetenv("IKEMEN_RENDER_PROBES")
		}
	}
	defer restore()

	_ = os.Unsetenv("IKEMEN_RENDER_PROBES")
	if got := parseRenderProbeMode(); got != "" {
		t.Fatalf("expected empty mode for unset env, got %q", got)
	}

	if err := os.Setenv("IKEMEN_RENDER_PROBES", " yes "); err != nil {
		t.Fatalf("Setenv yes: %v", err)
	}
	if got := parseRenderProbeMode(); got != "all" {
		t.Fatalf("expected yes to map to all, got %q", got)
	}

	if err := os.Setenv("IKEMEN_RENDER_PROBES", "hud"); err != nil {
		t.Fatalf("Setenv hud: %v", err)
	}
	if got := parseRenderProbeMode(); got != "hud" {
		t.Fatalf("expected literal category, got %q", got)
	}
}

func TestGetViewport(t *testing.T) {
	t.Parallel()

	if got := getViewport(1280, 720, 640, 480); got != [4]float64{160, 0, 960, 720} {
		t.Fatalf("wide viewport = %#v", got)
	}
	if got := getViewport(640, 480, 1280, 720); got != [4]float64{0, 60, 640, 360} {
		t.Fatalf("tall viewport = %#v", got)
	}
	if got := getViewport(800, 600, 1600, 1200); got != [4]float64{0, 0, 800, 600} {
		t.Fatalf("same-aspect viewport = %#v", got)
	}
}
