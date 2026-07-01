package main

import "testing"

func TestMusicAppendParamsAndPinSelection_HandleMusicFieldsAndSelection(t *testing.T) {
	m := Music{}
	m.AppendParams([]string{
		"stage.bgmusic = intro.ogg",
		"stage.bgmloop = 2",
		"stage.bgmvolume = 80",
		"stage.extra = ignored",
		"title.music = title.ogg",
	})

	if len(m["stage"]) != 1 || m["stage"][0].bgmusic != "intro.ogg" || m["stage"][0].bgmloop != 2 || m["stage"][0].bgmvolume != 80 {
		t.Fatalf("AppendParams stage entry = %#v", m["stage"])
	}
	if len(m["title"]) != 1 || m["title"][0].bgmusic != "title.ogg" {
		t.Fatalf("AppendParams title entry = %#v", m["title"])
	}
	if _, ok := m["stage_extra"]; ok {
		t.Fatal("AppendParams should ignore non-music fields")
	}

	first := &bgMusic{bgmusic: "a"}
	second := &bgMusic{bgmusic: "b"}
	m["battle"] = []*bgMusic{first, second}
	m.pinSelection("battle.music")
	if len(m["battle"]) != 1 || (!first.selected && !second.selected) {
		t.Fatalf("pinSelection should keep one selected candidate: %#v", m["battle"])
	}

	m.ClearSelection()
	if first.selected || second.selected {
		t.Fatal("ClearSelection should clear selection flags")
	}
}
