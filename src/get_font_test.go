package main

import "testing"

func TestGetFont(t *testing.T) {
	t.Parallel()

	var f map[int]*Fnt
	if got := getFont(f, 1); got != nil {
		t.Fatalf("getFont(nil, 1) = %#v, want nil", got)
	}
	f = map[int]*Fnt{2: {}}
	if got := getFont(f, 2); got == nil {
		t.Fatal("getFont should return existing font")
	}
	if got := getFont(f, -1); got != nil {
		t.Fatalf("getFont should reject negative index, got %#v", got)
	}
	if got := getFont(f, 3); got != nil {
		t.Fatalf("getFont missing key = %#v, want nil", got)
	}
}
