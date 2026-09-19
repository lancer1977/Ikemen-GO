package main

import "testing"

// TestFontKeyBackslashNormalized verifies that fontKey normalizes backslashes to forward slashes
// for cross-platform consistency. Windows-style paths (e.g., font\select.fnt) are normalized
// to forward slashes (font/select.fnt) regardless of the build OS, ensuring correct deduplication
// on Unix systems when config files contain Windows-authored paths.
func TestFontKeyBackslashNormalized(t *testing.T) {
	t.Parallel()

	if got := fontKey(`font\select.fnt`, 12); got != `font/select.fnt|12` {
		t.Fatalf("fontKey = %q, want font/select.fnt|12", got)
	}
}

// TestRegisterFontIndexBackslashKey verifies that registerFontIndex uses normalized keys.
// When a Windows-style path (font\select.fnt) is registered, the key is normalized to
// forward slashes (font/select.fnt|12). This ensures that lookups for the same font
// with different path conventions (e.g., font/select.fnt vs font\select.fnt) will find
// the registered index.
func TestRegisterFontIndexBackslashKey(t *testing.T) {
	t.Parallel()

	indexByKey := map[string]int{}
	registerFontIndex(indexByKey, `font\select.fnt`, 12, 3)
	// The key is normalized to forward slashes
	if got := indexByKey[`font/select.fnt|12`]; got != 3 {
		t.Fatalf("registerFontIndex = %d, want 3", got)
	}

	registerFontIndex(nil, "font.fnt", 12, 4)
	registerFontIndex(indexByKey, "", 12, 4)
	if len(indexByKey) != 1 {
		t.Fatalf("expected nil/empty inputs to be ignored, got %#v", indexByKey)
	}
}
