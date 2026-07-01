package main

import (
	"math"
	"testing"
)

func TestBytecodeValueConversionsAndSetters(t *testing.T) {
	none := bvNone()
	if !none.IsNone() || none.IsUndefined() {
		t.Fatalf("bvNone() state = %#v", none)
	}

	undef := BytecodeUndefined()
	if !undef.IsUndefined() || undef.ToF() != 0 || undef.ToI() != 0 || undef.ToI64() != 0 || undef.ToB() {
		t.Fatalf("BytecodeUndefined() conversions = %#v", undef)
	}
	if _, ok := undef.ToAny().(UndefinedFormatter); !ok {
		t.Fatalf("BytecodeUndefined().ToAny() should return UndefinedFormatter, got %#v", undef.ToAny())
	}

	if got := BytecodeFloat(3.5); got.ToF() != 3.5 || got.ToI() != 3 || got.ToI64() != 3 || !got.ToB() {
		t.Fatalf("BytecodeFloat conversions = %#v", got)
	}
	if got := BytecodeFloat(float32(math.NaN())); !got.IsUndefined() {
		t.Fatalf("BytecodeFloat(NaN) should become undefined, got %#v", got)
	}

	if got := BytecodeInt(7); got.ToF() != 7 || got.ToI() != 7 || got.ToI64() != 7 || !got.ToB() {
		t.Fatalf("BytecodeInt conversions = %#v", got)
	}
	if got := BytecodeInt64(9); got.ToF() != 9 || got.ToI() != 9 || got.ToI64() != 9 || !got.ToB() {
		t.Fatalf("BytecodeInt64 conversions = %#v", got)
	}
	if !BytecodeBool(true).ToB() || BytecodeBool(false).ToB() {
		t.Fatal("BytecodeBool() conversions failed")
	}

	var v BytecodeValue
	v.SetF(4.5)
	if v.IsUndefined() || v.ToF() != 4.5 || v.ToI() != 4 {
		t.Fatalf("SetF() = %#v", v)
	}
	v.SetI(12)
	if v.ToI() != 12 || v.ToI64() != 12 || !v.ToB() {
		t.Fatalf("SetI() = %#v", v)
	}
	v.SetI64(14)
	if v.ToI64() != 14 || v.ToI() != 14 {
		t.Fatalf("SetI64() = %#v", v)
	}
	v.SetB(true)
	if !v.ToB() || v.ToI() != 1 {
		t.Fatalf("SetB(true) = %#v", v)
	}
	v.SetB(false)
	if v.ToB() {
		t.Fatalf("SetB(false) = %#v", v)
	}
}
