package main

import "testing"

func TestNewCommandKeyRemap(t *testing.T) {
	t.Parallel()

	r := NewCommandKeyRemap()
	if r == nil {
		t.Fatal("NewCommandKeyRemap returned nil")
	}
	if *r != (CommandKeyRemap{CK_a, CK_b, CK_c, CK_x, CK_y, CK_z, CK_s, CK_d, CK_w, CK_m}) {
		t.Fatalf("unexpected remap: %#v", r)
	}
}
