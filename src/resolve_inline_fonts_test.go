package main

import "testing"

func TestResolveInlineFontsNilInput(t *testing.T) {
	t.Parallel()

	called := false
	resolveInlineFonts(nil, "base", nil, nil, func(query string, value interface{}) error {
		called = true
		return nil
	})
	if called {
		t.Fatal("nil ini file should not invoke callback")
	}
}

func TestResolveInlineFontsSkipsNonFontKeys(t *testing.T) {
	t.Parallel()

	f := NewIniFile()
	sec, err := f.NewSection("General")
	if err != nil {
		t.Fatalf("NewSection: %v", err)
	}
	sec.NewKey("title", "value")

	called := false
	resolveInlineFonts(f, "base", nil, nil, func(query string, value interface{}) error {
		called = true
		return nil
	})
	if called {
		t.Fatal("non-font keys should be skipped")
	}
}

func TestResolveInlineFontsSkipsFilesSection(t *testing.T) {
	t.Parallel()

	f := NewIniFile()
	sec, err := f.NewSection("Files")
	if err != nil {
		t.Fatalf("NewSection: %v", err)
	}
	sec.NewKey("font", "select.fnt")

	called := false
	resolveInlineFonts(f, "base", nil, nil, func(query string, value interface{}) error {
		called = true
		return nil
	})
	if called {
		t.Fatal("[Files] section should be skipped")
	}
}

func TestResolveInlineFontsSkipsMusicSection(t *testing.T) {
	t.Parallel()

	f := NewIniFile()
	sec, err := f.NewSection("Music")
	if err != nil {
		t.Fatalf("NewSection: %v", err)
	}
	sec.NewKey("font", "select.fnt")

	called := false
	resolveInlineFonts(f, "base", nil, nil, func(query string, value interface{}) error {
		called = true
		return nil
	})
	if called {
		t.Fatal("[Music] section should be skipped")
	}
}

func TestResolveInlineFontsSkipsDefaultSection(t *testing.T) {
	t.Parallel()

	f := NewIniFile()
	sec, err := f.NewSection("")
	if err != nil {
		t.Fatalf("NewSection: %v", err)
	}
	sec.NewKey("font", "select.fnt")

	called := false
	resolveInlineFonts(f, "base", nil, nil, func(query string, value interface{}) error {
		called = true
		return nil
	})
	if called {
		t.Fatal("default section should be skipped")
	}
}

func TestResolveInlineFontsSkipsEmptyFontValue(t *testing.T) {
	t.Parallel()

	f := NewIniFile()
	sec, err := f.NewSection("General")
	if err != nil {
		t.Fatalf("NewSection: %v", err)
	}
	sec.NewKey("font", "   ")

	called := false
	resolveInlineFonts(f, "base", nil, nil, func(query string, value interface{}) error {
		called = true
		return nil
	})
	if called {
		t.Fatal("empty font value should be skipped")
	}
}

func TestResolveInlineFontsSkipsAlreadyIndexedFont(t *testing.T) {
	t.Parallel()

	f := NewIniFile()
	sec, err := f.NewSection("General")
	if err != nil {
		t.Fatalf("NewSection: %v", err)
	}
	sec.NewKey("font", "3,0,0")

	called := false
	resolveInlineFonts(f, "base", nil, nil, func(query string, value interface{}) error {
		called = true
		return nil
	})
	if called {
		t.Fatal("already-indexed font should be skipped")
	}
}

func TestResolveInlineFontsUsesHeightFromEighthField(t *testing.T) {
	t.Parallel()

	f := NewIniFile()
	sec, err := f.NewSection("General")
	if err != nil {
		t.Fatalf("NewSection: %v", err)
	}
	sec.NewKey("font", "font.fnt,0,0,0,0,0,0,42")

	indexByKey := map[string]int{fontKey("font.fnt", 42): 7}
	var gotQuery string
	var gotValue interface{}
	resolveInlineFonts(f, "base", nil, indexByKey, func(query string, value interface{}) error {
		gotQuery = query
		gotValue = value
		return nil
	})

	if gotQuery != "General.font" {
		t.Fatalf("query = %q, want %q", gotQuery, "General.font")
	}
	if gotValue != "7,0,0,0,0,0,0,42" {
		t.Fatalf("value = %v, want %q", gotValue, "7,0,0,0,0,0,0,42")
	}
}

func TestEnsureFontIndexReturnsPreseededIndex(t *testing.T) {
	t.Parallel()

	indexByKey := map[string]int{fontKey("font.fnt", 42): 7}
	if got := ensureFontIndex(nil, indexByKey, "base", "font.fnt", 42); got != 7 {
		t.Fatalf("ensureFontIndex = %d, want 7", got)
	}
}
