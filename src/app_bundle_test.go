package main

import (
	"runtime"
	"testing"
)

func TestIsRunningInsideAppBundle(t *testing.T) {
	if got := isRunningInsideAppBundle("/usr/local/bin/ikemen"); got {
		t.Fatal("expected non-app-bundle path to return false")
	}
	if got := isRunningInsideAppBundle("/Applications/Ikemen.app/Contents/MacOS/ikemen"); runtime.GOOS == "darwin" && !got {
		t.Fatal("expected macOS app bundle path to return true on darwin")
	}
	if got := isRunningInsideAppBundle("/Applications/Ikemen.app/Contents/MacOS/ikemen"); runtime.GOOS != "darwin" && got {
		t.Fatal("expected non-darwin runtime to return false")
	}
}
