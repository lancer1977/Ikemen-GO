package main

import "testing"

func TestFontKey(t *testing.T) {
	t.Parallel()

	if got := fontKey(`font\select.fnt`, 12); got != "font/select.fnt|12" {
		t.Fatalf("fontKey = %q", got)
	}
}

func TestRegisterFontIndex(t *testing.T) {
	t.Parallel()

	indexByKey := map[string]int{}
	registerFontIndex(indexByKey, `font\select.fnt`, 12, 3)
	if got := indexByKey["font/select.fnt|12"]; got != 3 {
		t.Fatalf("registerFontIndex = %d, want 3", got)
	}

	registerFontIndex(nil, "font.fnt", 12, 4)
	registerFontIndex(indexByKey, "", 12, 4)
	if len(indexByKey) != 1 {
		t.Fatalf("expected nil/empty inputs to be ignored, got %#v", indexByKey)
	}
}
