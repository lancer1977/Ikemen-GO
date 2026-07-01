package main

import (
	"testing"

	"gopkg.in/ini.v1"
)

func TestPickLangSection(t *testing.T) {
	t.Parallel()

	oldLang := sys.cfg.Config.Language
	sys.cfg.Config.Language = "pt"
	defer func() { sys.cfg.Config.Language = oldLang }()

	f := ini.Empty()
	base, _ := f.NewSection("select")
	lang, _ := f.NewSection("pt.select")
	_, _ = base.NewKey("title", "base")
	_, _ = lang.NewKey("title", "pt")

	if got := pickLangSection(f, "select"); got == nil || got.Name() != "pt.select" {
		t.Fatalf("pickLangSection() = %#v", got)
	}
	merged := pickLangSectionMerged(f, "select")
	if merged == nil || merged.Name() != "select" {
		t.Fatalf("pickLangSectionMerged() = %#v", merged)
	}
	if got := merged.Key("title").String(); got != "pt" {
		t.Fatalf("merged title = %q, want pt", got)
	}
}
