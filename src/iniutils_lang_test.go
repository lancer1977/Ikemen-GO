package main

import (
	"testing"

	"gopkg.in/ini.v1"
)

func TestSelectedLanguage_NormalizesExplicitConfiguration(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys.cfg.Config.Language = "  PT-BR "
	if got := SelectedLanguage(); got != "pt" {
		t.Fatalf("SelectedLanguage() = %q, want pt", got)
	}

	sys.cfg.Config.Language = ""
	if got := SelectedLanguage(); got != "en" {
		t.Fatalf("SelectedLanguage() with empty config = %q, want en", got)
	}
}

func TestResolveLangSectionNameAndNames_PreferSpecificAndBaseSections(t *testing.T) {
	f := ini.Empty()
	_, _ = f.NewSection("select")
	_, _ = f.NewSection("pt.select")

	if got := ResolveLangSectionName(f, "select", "pt"); got != "pt.select" {
		t.Fatalf("ResolveLangSectionName() = %q, want pt.select", got)
	}
	if got := ResolveLangSectionName(f, "select", "fr"); got != "select" {
		t.Fatalf("ResolveLangSectionName() fallback = %q, want select", got)
	}
	if got := ResolveLangSectionName(f, "missing", "fr"); got != "missing" {
		t.Fatalf("ResolveLangSectionName() missing = %q, want missing", got)
	}

	if got := ResolveLangSectionNames(f, "select", "pt"); len(got) != 2 || got[0] != "select" || got[1] != "pt.select" {
		t.Fatalf("ResolveLangSectionNames() = %#v, want [select pt.select]", got)
	}
	if got := ResolveLangSectionNames(f, "missing", "pt"); len(got) != 1 || got[0] != "missing" {
		t.Fatalf("ResolveLangSectionNames() missing = %#v, want [missing]", got)
	}
}
