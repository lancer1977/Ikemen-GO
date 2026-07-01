package main

import (
	"reflect"
	"testing"
)

func TestIsMusicKey(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		in   string
		want bool
	}{
		{"music.title", true},
		{"stage.bgm", true},
		{"charparam.music", false},
		{"notmusic", false},
	} {
		if got := isMusicKey(tc.in); got != tc.want {
			t.Fatalf("isMusicKey(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestIsZeroValue(t *testing.T) {
	t.Parallel()

	type sample struct {
		A int
		B []string
		C struct {
			N int
		}
	}

	if !isZeroValue(reflect.ValueOf(sample{})) {
		t.Fatal("zero struct should be zero value")
	}

	if isZeroValue(reflect.ValueOf(sample{A: 1})) {
		t.Fatal("non-zero struct should not be zero value")
	}

	if !isZeroValue(reflect.ValueOf([]int(nil))) {
		t.Fatal("nil slice should be zero value")
	}
}
