package main

import (
	"testing"

	"github.com/veandco/go-sdl2/sdl"
)

func TestInputLookupHelpers(t *testing.T) {
	initLUTs()

	if got := StringToKey("ESCAPE"); got != sdl.K_ESCAPE {
		t.Fatalf("StringToKey(ESCAPE) = %v", got)
	}
	if got := StringToKey("missing"); got != sdl.K_UNKNOWN {
		t.Fatalf("StringToKey(missing) = %v, want unknown", got)
	}
	if got := KeyToString(sdl.K_RETURN); got != "RETURN" {
		t.Fatalf("KeyToString(RETURN) = %q", got)
	}
	if got := KeyToString(sdl.K_UNKNOWN); got != "" {
		t.Fatalf("KeyToString(unknown) = %q, want blank", got)
	}

	got := NewModifierKey(true, true, true)
	want := sdl.KMOD_CTRL | sdl.KMOD_ALT | sdl.KMOD_SHIFT
	if got != want {
		t.Fatalf("NewModifierKey = %v, want %v", got, want)
	}
}
