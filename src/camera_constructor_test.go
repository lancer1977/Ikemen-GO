package main

import "testing"

func TestNewCameraInitializesDefaults(t *testing.T) {
	c := newCamera()
	if c == nil {
		t.Fatal("newCamera returned nil")
	}
	if c.View != Fighting_View {
		t.Fatalf("newCamera View = %v, want %v", c.View, Fighting_View)
	}
	if c.LegacyZoomMin != 5.0/8.0 || c.LegacyZoomMax != 1 {
		t.Fatalf("unexpected legacy zoom defaults: %#v", c)
	}
}
