package main

import "testing"

func TestNewIniSection_CreatesWritableEmptyMap(t *testing.T) {
	is := NewIniSection()
	if is == nil || len(is) != 0 {
		t.Fatalf("NewIniSection() = %#v, want empty map", is)
	}
	is["value"] = "10"
	if is["value"] != "10" {
		t.Fatalf("NewIniSection map did not accept writes: %#v", is)
	}
}

func TestIniSectionParse_StripsCommentsIgnoresInvalidLinesAndKeepsFirstValue(t *testing.T) {
	is := NewIniSection()
	lines := []string{
		"  ; comment",
		"invalid line",
		"foo = 1 ; trailing comment",
		"foo = 2",
		"bar 3",
		"[Next]",
	}
	i := 0
	is.Parse(lines, &i)
	if is["foo"] != "1" {
		t.Fatalf("Parse kept foo = %q, want %q", is["foo"], "1")
	}
	if is["bar"] != "3" {
		t.Fatalf("Parse kept bar = %q, want %q", is["bar"], "3")
	}
	if i != 5 {
		t.Fatalf("Parse cursor = %d, want 5", i)
	}
}
