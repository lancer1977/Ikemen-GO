package main

import "testing"

func TestSplitLangPrefixRejectsNonTwoLetterPrefixes(t *testing.T) {
	lang, base, ok := splitLangPrefix("eng.select")
	if ok || lang != "" || base != "eng.select" {
		t.Fatalf("splitLangPrefix(eng.select) = %q %q %v, want no prefix", lang, base, ok)
	}

	lang, base, ok = splitLangPrefix("a.select")
	if ok || lang != "" || base != "a.select" {
		t.Fatalf("splitLangPrefix(a.select) = %q %q %v, want no prefix", lang, base, ok)
	}
}
