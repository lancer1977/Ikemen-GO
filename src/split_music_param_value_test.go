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
	// DEFECT: splitMusicParamValue tokenises with strings.Fields(), which splits on
	// whitespace only, so commas survive as standalone tokens. expandMusicKV then
	// consumes those tokens positionally as volume/loopstart/loopend, and Atoi(",")
	// is 0 -- a comma-separated music line plays silently instead of at the default
	// volume of 100. Tracked as lancer1977/Ikemen-GO#16.
	if !reflect.DeepEqual(extras, []string{",", "loop=1", ",", "vol=80"}) {
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
