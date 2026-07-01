package main

import "testing"

func TestNewFightScreenScore(t *testing.T) {
	t.Parallel()

	sc := newFightScreenScore()
	if sc == nil {
		t.Fatal("newFightScreenScore returned nil")
	}
	if sc.enabled == nil {
		t.Fatal("newFightScreenScore should allocate enabled map")
	}
	if sc.separator != [2]string{"", "."} {
		t.Fatalf("separator = %#v, want [\"\" \".\"]", sc.separator)
	}
	if sc.active {
		t.Fatal("newFightScreenScore should start inactive")
	}
}
