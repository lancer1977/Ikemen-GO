package main

import (
	"strings"
	"testing"
)

func TestCountFmtVerbsIgnoresEscapedPercents(t *testing.T) {
	t.Parallel()

	if got := countFmtVerbs("score %% %s %03d"); got != 2 {
		t.Fatalf("countFmtVerbs = %d, want 2", got)
	}
	if got := countFmtVerbs("literal only %%"); got != 0 {
		t.Fatalf("countFmtVerbs = %d, want 0", got)
	}
}

func TestPreprocessINIContentMovesInfoboxTextIntoInfoBox(t *testing.T) {
	t.Parallel()

	in := "[InfoBox]\n[Infobox Text]\nline 1\nline 2 %s %s\n"
	got := preprocessINIContent(in)
	if strings.Contains(got, "[Infobox Text]") {
		t.Fatalf("preprocessINIContent kept source section: %q", got)
	}
	wantLine := "\ttext.text = line 1\\nline 2 " + Version + " " + BuildTime + "\n\n"
	if !strings.Contains(got, wantLine) {
		t.Fatalf("preprocessINIContent missing rewritten text line: %q", got)
	}
}

func TestMotifHiscoreNameHelpers(t *testing.T) {
	t.Parallel()

	mo := &Motif{
		HiscoreInfo: HiscoreInfoProperties{
			Glyphs: []string{"A", ">", "C"},
			Item: HiscoreItemProperties{
				Name: HiscoreItemNameProperties{
					Uppercase: true,
					Text: map[string]string{
						"default": "%s",
					},
				},
			},
		},
	}

	if got := initialsWidth("name %4s"); got != 4 {
		t.Fatalf("initialsWidth = %d, want 4", got)
	}
	if got := initialsWidth("name %s"); got != 3 {
		t.Fatalf("initialsWidth default = %d, want 3", got)
	}
	if got := currentGlyph(mo, []int{1, 3}); got != "C" {
		t.Fatalf("currentGlyph = %q, want %q", got, "C")
	}
	if got := currentGlyph(mo, []int{0}); got != "" {
		t.Fatalf("currentGlyph invalid index = %q, want blank", got)
	}
	if got := buildNameFromLetters(mo, []int{1, 2, 3, 4}); got != "A C" {
		t.Fatalf("buildNameFromLetters = %q, want %q", got, "A C")
	}

	row := &rankingRow{
		nameData:        &TextSprite{},
		nameDataActive:  &TextSprite{},
		nameDataActive2: &TextSprite{},
	}
	updateRowNameFromLetters(mo, row, []int{1, 2, 3})
	if row.name != "A C" {
		t.Fatalf("row.name = %q, want %q", row.name, "A C")
	}
	if row.nameData.text != "A C" || row.nameDataActive.text != "A C" || row.nameDataActive2.text != "A C" {
		t.Fatalf("row text sprites not updated consistently: %+v %+v %+v", row.nameData, row.nameDataActive, row.nameDataActive2)
	}
}

func TestMatchDef(t *testing.T) {
	t.Parallel()

	if !matchDef(`chars\\kfm_zss\\char.def`, "kfm_zss") {
		t.Fatal("expected directory name match")
	}
	if !matchDef(`chars\\kfm\\char.def`, "char") {
		t.Fatal("expected basename match")
	}
	if matchDef(`chars\\kfm\\char.def`, "") {
		t.Fatal("expected empty key to be rejected")
	}
}
