package main

import (
	"reflect"
	"testing"
)

func TestSplitMusicParamValue(t *testing.T) {
	path, extras := splitMusicParamValue("  sound/title.ogg  , loop=1 , vol=80 ")
	if path != "sound/title.ogg" {
		t.Fatalf("splitMusicParamValue path = %q", path)
	}
	// splitMusicParamValue now strips commas before tokenising, so comma-separated
	// format is correctly parsed alongside space-separated format.
	if !reflect.DeepEqual(extras, []string{"loop=1", "vol=80"}) {
		t.Fatalf("splitMusicParamValue extras = %#v", extras)
	}

	path, extras = splitMusicParamValue("music/intro")
	if path != "music/intro" || extras != nil {
		t.Fatalf("splitMusicParamValue no-extension = %q %#v", path, extras)
	}

	path, extras = splitMusicParamValue("   ")
	if path != "" || extras != nil {
		t.Fatalf("splitMusicParamValue blank = %q %#v", path, extras)
	}
}
