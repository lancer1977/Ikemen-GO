package main

import (
	"os"
	"testing"
)

func TestParseRenderProbeModeTreatsFalseyValuesAsDisabled(t *testing.T) {
	prev, had := os.LookupEnv("IKEMEN_RENDER_PROBES")
	defer func() {
		if had {
			_ = os.Setenv("IKEMEN_RENDER_PROBES", prev)
		} else {
			_ = os.Unsetenv("IKEMEN_RENDER_PROBES")
		}
	}()

	for _, v := range []string{"0", "false", "No", "off"} {
		if err := os.Setenv("IKEMEN_RENDER_PROBES", v); err != nil {
			t.Fatalf("Setenv(%q): %v", v, err)
		}
		if got := parseRenderProbeMode(); got != "" {
			t.Fatalf("parseRenderProbeMode(%q) = %q, want empty", v, got)
		}
	}
}
