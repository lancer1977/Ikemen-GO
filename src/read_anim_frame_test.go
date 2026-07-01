package main

import (
	"strings"
	"testing"
)

func TestReadAnimFrame_ParsesFlagsAlphaAndRejectsInvalidLines(t *testing.T) {
	if af, err := ReadAnimFrame("not a frame"); af != nil || err != nil {
		t.Fatalf("ReadAnimFrame(non-frame) = %#v, %v; want nil, nil", af, err)
	}

	if af, err := ReadAnimFrame("1, 2, 3, 4"); af != nil || err == nil {
		t.Fatalf("ReadAnimFrame(short) = %#v, %v; want nil, error", af, err)
	}

	af, err := ReadAnimFrame("1, 2, 3, 4, 5, HV, AS64d32, 1.5, 2.5, 3.5")
	if err != nil {
		t.Fatalf("ReadAnimFrame(valid) error = %v", err)
	}
	if af == nil {
		t.Fatal("ReadAnimFrame(valid) returned nil frame")
	}
	if af.Group != 1 || af.Number != 2 || af.Xoffset != -3 || af.Yoffset != -4 || af.Time != 5 {
		t.Fatalf("ReadAnimFrame frame basics = %#v", af)
	}
	if af.Hscale != -1 || af.Vscale != -1 {
		t.Fatalf("ReadAnimFrame scales = %v %v, want -1 -1", af.Hscale, af.Vscale)
	}
	if af.TransType != TT_add || af.SrcAlpha != 64 || af.DstAlpha != 32 {
		t.Fatalf("ReadAnimFrame alpha = %#v", af)
	}
	if af.Xscale != 1.5 || af.Yscale != 2.5 || af.Angle != 3.5 {
		t.Fatalf("ReadAnimFrame scale/angle = %#v", af)
	}

	if af, err := ReadAnimFrame("1, 2, 3, 4, 5, H, Z, 1, 1, 0"); af == nil || err == nil {
		t.Fatalf("ReadAnimFrame(invalid alpha) = %#v, %v; want frame and error", af, err)
	}
	if af, err := ReadAnimFrame("1, 2, 3, 4, 5, HX, A, 1, 1, 0"); af == nil || err == nil || !strings.Contains(err.Error(), "flip flag") {
		t.Fatalf("ReadAnimFrame(invalid flip) = %#v, %v; want error containing flip flag", af, err)
	}

	for _, tc := range []struct {
		name string
		line string
		want string
	}{
		{"xscale", "1, 2, 3, 4, 5, H, A, bad, 1, 0", "x-scale"},
		{"yscale", "1, 2, 3, 4, 5, H, A, 1, bad, 0", "y-scale"},
		{"angle", "1, 2, 3, 4, 5, H, A, 1, 1, bad", "angle"},
	} {
		af, err := ReadAnimFrame(tc.line)
		if af == nil || err == nil || err.Error() == "" || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("ReadAnimFrame(%s) = %#v, %v; want error containing %q", tc.name, af, err, tc.want)
		}
	}
}

func TestReadAnimFrame_PreservesDefaultsWhenOptionalFieldsAreMissing(t *testing.T) {
	af, err := ReadAnimFrame("1, 2, 3, 4, 5")
	if err != nil {
		t.Fatalf("ReadAnimFrame(defaults) error = %v", err)
	}
	if af == nil {
		t.Fatal("ReadAnimFrame(defaults) returned nil frame")
	}
	if af.Group != 1 || af.Number != 2 || af.Xoffset != 3 || af.Yoffset != 4 || af.Time != 5 {
		t.Fatalf("ReadAnimFrame(defaults) basics = %#v", af)
	}
	if af.Hscale != 1 || af.Vscale != 1 || af.TransType != TT_none || af.SrcAlpha != 255 || af.DstAlpha != 0 {
		t.Fatalf("ReadAnimFrame(defaults) alpha/flags = %#v", af)
	}
	if af.Xscale != 1 || af.Yscale != 1 || af.Angle != 0 {
		t.Fatalf("ReadAnimFrame(defaults) scale/angle = %#v", af)
	}
}
