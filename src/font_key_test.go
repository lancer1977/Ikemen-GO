package main

import "testing"

// TestFontKeyBackslashNotNormalized documents the real behavior of fontKey on Unix systems.
// DEFECT: filepath.ToSlash only converts OS-specific separators, not backslashes.
// On Windows, this would normalize backslashes to forward slashes. On Unix, backslashes
// are not path separators, so they are left as-is. This breaks cross-platform deduplication
// when config files contain Windows-style paths on Unix systems. Production code should
// manually replace both backslashes and forward slashes with a canonical separator.
func TestFontKeyBackslashNotNormalized(t *testing.T) {
	t.Parallel()

	if got := fontKey(`font\select.fnt`, 12); got != `font\select.fnt|12` {
		t.Fatalf("fontKey = %q", got)
	}
}

// TestRegisterFontIndexBackslashKey documents the real behavior when registering font indices.
// DEFECT: Cascading from fontKey's backslash handling, the key will contain backslashes
// on Unix systems, not normalized forward slashes. This means lookups that expect
// normalized keys (like from forward-slash paths) will fail to find the registered index.
func TestRegisterFontIndexBackslashKey(t *testing.T) {
	t.Parallel()

	indexByKey := map[string]int{}
	registerFontIndex(indexByKey, `font\select.fnt`, 12, 3)
	// On Unix, the key is stored with backslash, not normalized to forward slash
	if got := indexByKey[`font\select.fnt|12`]; got != 3 {
		t.Fatalf("registerFontIndex = %d, want 3", got)
	}

	registerFontIndex(nil, "font.fnt", 12, 4)
	registerFontIndex(indexByKey, "", 12, 4)
	if len(indexByKey) != 1 {
		t.Fatalf("expected nil/empty inputs to be ignored, got %#v", indexByKey)
	}
}
