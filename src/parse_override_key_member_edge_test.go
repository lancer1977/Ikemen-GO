package main

import "testing"

func TestParseOverrideKeyRejectsNonNumericMember(t *testing.T) {
	if _, _, _, ok := parseOverrideKey("p1.two.life"); ok {
		t.Fatal("parseOverrideKey should reject a non-numeric member")
	}
}
