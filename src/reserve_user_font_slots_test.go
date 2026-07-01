package main

import "testing"

func TestReserveUserFontSlots(t *testing.T) {
	t.Parallel()

	m := &Motif{
		UserIniFile: NewIniFile(),
		Fnt:         map[int]*Fnt{3: newFnt()},
	}
	files, err := m.UserIniFile.NewSection("Files")
	if err != nil {
		t.Fatalf("NewSection: %v", err)
	}
	files.NewKey("font3", "select.fnt")
	files.NewKey("font7", "stage.fnt")
	files.NewKey("other", "ignored")

	reserveUserFontSlots(m)

	if _, ok := m.Fnt[3]; !ok {
		t.Fatal("existing occupied font slot should remain tracked")
	}
	if m.Fnt[7] != nil {
		t.Fatalf("font7 should be reserved as nil placeholder, got %#v", m.Fnt[7])
	}
	if _, ok := m.Fnt[8]; ok {
		t.Fatalf("unexpected font slot reservation: %#v", m.Fnt)
	}
}
