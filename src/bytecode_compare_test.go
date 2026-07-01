package main

import "testing"

func TestBytecodeExpComparisonHelpersHandleIntAndFloatBranches(t *testing.T) {
	compare := BytecodeExp{}

	v := BytecodeInt(7)
	compare.gt(&v, BytecodeInt(6))
	if !v.ToB() {
		t.Fatalf("gt() int branch = %#v, want true", v)
	}
	compare.ge(&v, BytecodeInt(7))
	if !v.ToB() {
		t.Fatalf("ge() int branch = %#v, want true", v)
	}
	compare.lt(&v, BytecodeInt(8))
	if !v.ToB() {
		t.Fatalf("lt() int branch = %#v, want true", v)
	}
	compare.le(&v, BytecodeInt(7))
	if !v.ToB() {
		t.Fatalf("le() int branch = %#v, want true", v)
	}
	compare.eq(&v, BytecodeInt(7))
	if !v.ToB() {
		t.Fatalf("eq() int branch = %#v, want true", v)
	}
	compare.ne(&v, BytecodeInt(8))
	if !v.ToB() {
		t.Fatalf("ne() int branch = %#v, want true", v)
	}

	v = BytecodeFloat(3.5)
	compare.gt(&v, BytecodeInt(3))
	if !v.ToB() {
		t.Fatalf("gt() float branch = %#v, want true", v)
	}
	compare.le(&v, BytecodeFloat(3.5))
	if !v.ToB() {
		t.Fatalf("le() float branch = %#v, want true", v)
	}
	compare.eq(&v, BytecodeFloat(3.5))
	if !v.ToB() {
		t.Fatalf("eq() float branch = %#v, want true", v)
	}
	compare.ne(&v, BytecodeFloat(4.5))
	if !v.ToB() {
		t.Fatalf("ne() float branch = %#v, want true", v)
	}
}

func TestBytecodeExpComparisonHelpersWithUndefinedRemainBooleanized(t *testing.T) {
	compare := BytecodeExp{}

	v := BytecodeUndefined()
	compare.eq(&v, BytecodeUndefined())
	if !v.ToB() {
		t.Fatalf("eq() with undefined values should still booleanize, got %#v", v)
	}

	v = BytecodeUndefined()
	compare.ne(&v, BytecodeInt(1))
	if !v.ToB() {
		t.Fatalf("ne() with undefined lhs should booleanize, got %#v", v)
	}
}
