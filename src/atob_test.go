package main

import "testing"

func TestAtob_ParsesCommonInputsAndRejectsInvalid(t *testing.T) {
	for _, s := range []string{"true", "TRUE", "1"} {
		if !Atob(s) {
			t.Fatalf("Atob(%q) = false, want true", s)
		}
	}
	for _, s := range []string{"false", "FALSE", "0", "not-a-bool"} {
		if Atob(s) {
			t.Fatalf("Atob(%q) = true, want false", s)
		}
	}
}
