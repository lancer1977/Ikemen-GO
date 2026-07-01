package main

import "testing"

func TestFindActiveSffSearchOrder(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	charSff := &Sff{filename: "shared.sff"}
	stageSff := &Sff{filename: "shared.sff"}
	ffxSff := &Sff{filename: "shared.sff"}

	sys.cgi[0].sff = charSff
	sys.stage = &Stage{sff: stageSff}
	sys.ffx = map[string]*FightFx{
		"fx": &FightFx{sff: ffxSff},
	}

	if got := findActiveSff("shared.sff"); got != charSff {
		t.Fatalf("expected character SFF to win lookup, got %#v", got)
	}
	sys.cgi[0].sff = nil
	if got := findActiveSff("shared.sff"); got != stageSff {
		t.Fatalf("expected stage SFF to win when no char match, got %#v", got)
	}
	sys.stage = nil
	if got := findActiveSff("shared.sff"); got != ffxSff {
		t.Fatalf("expected common FX SFF to win when no char/stage match, got %#v", got)
	}
	if got := findActiveSff("missing.sff"); got != nil {
		t.Fatalf("expected missing SFF lookup to return nil, got %#v", got)
	}
}
