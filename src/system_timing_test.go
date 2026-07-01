package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"
)

func TestLoadTimeRoutesToShellAndConsole(t *testing.T) {
	oldSys := sys
	oldStdout := os.Stdout
	defer func() {
		sys = oldSys
		os.Stdout = oldStdout
	}()

	sys = oldSys
	sys.cfg.Debug.ConsoleRows = 4
	sys.consoleText = nil

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe(): %v", err)
	}
	os.Stdout = w

	sys.loadTime(time.Now().Add(-time.Millisecond), "load", true, true)
	w.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("reading stdout: %v", err)
	}
	r.Close()

	if got := buf.String(); !strings.Contains(got, "load; Load time:") {
		t.Fatalf("stdout output = %q, want load-time text", got)
	}
	if len(sys.consoleText) != 1 || !strings.Contains(sys.consoleText[0], "load; Load time:") {
		t.Fatalf("console output = %#v, want load-time text", sys.consoleText)
	}
}

func TestUpdateZScaleInterpolatesAndClamps(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys
	sys.stage = &Stage{}
	sys.stage.stageCamera.topz = 0
	sys.stage.stageCamera.botz = 100
	sys.stage.stageCamera.ztopscale = 2
	sys.stage.stageCamera.zbotscale = 0.5

	if got := sys.updateZScale(0, 1); got != 2 {
		t.Fatalf("updateZScale() at top = %v, want 2", got)
	}
	if got := sys.updateZScale(50, 1); got != 1.25 {
		t.Fatalf("updateZScale() mid-range = %v, want 1.25", got)
	}
	if got := sys.updateZScale(100, 1); got != 0.5 {
		t.Fatalf("updateZScale() at bottom = %v, want 0.5", got)
	}

	sys.stage.stageCamera.ztopscale = -1
	sys.stage.stageCamera.zbotscale = -2
	if got := sys.updateZScale(50, 1); got != 0 {
		t.Fatalf("updateZScale() should clamp negative scale to 0, got %v", got)
	}
}
