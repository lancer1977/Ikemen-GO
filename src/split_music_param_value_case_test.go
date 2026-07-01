package main

import (
	"reflect"
	"testing"
)

func TestSplitMusicParamValueHandlesMixedCaseExtensions(t *testing.T) {
	path, extras := splitMusicParamValue("  Sound/Intro.Mp3   12 34 ")
	if path != "Sound/Intro.Mp3" {
		t.Fatalf("splitMusicParamValue path = %q, want %q", path, "Sound/Intro.Mp3")
	}
	if !reflect.DeepEqual(extras, []string{"12", "34"}) {
		t.Fatalf("splitMusicParamValue extras = %#v, want [12 34]", extras)
	}
}
