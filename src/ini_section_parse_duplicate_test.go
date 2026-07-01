package main

import "testing"

func TestIniSectionParse_IgnoresDuplicateKeysAfterFirstAssignment(t *testing.T) {
	is := NewIniSection()
	lines := []string{
		"foo = 1",
		"foo = 2",
		"[Next]",
	}
	i := 0
	is.Parse(lines, &i)
	if is["foo"] != "1" {
		t.Fatalf("Parse duplicate key = %q, want %q", is["foo"], "1")
	}
	if i != 2 {
		t.Fatalf("Parse cursor = %d, want 2", i)
	}
}
