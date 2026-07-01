package main

import "testing"

func TestParseOverrideKeyRejectsEmptyFieldSegment(t *testing.T) {
	if _, _, _, ok := parseOverrideKey("p1.2."); ok {
		t.Fatal("parseOverrideKey should reject an empty field segment")
	}
}
