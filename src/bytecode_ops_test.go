package main

import "testing"

func TestBytecodeExpUnaryAndArithmeticOperators(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys
	sys.loader.state = LS_Complete

	v := BytecodeInt(5)
	(BytecodeExp{}).neg(&v)
	if v.ToI() != -5 {
		t.Fatalf("neg() on int = %#v, want -5", v)
	}

	v = BytecodeFloat(2.5)
	(BytecodeExp{}).neg(&v)
	if v.ToF() != -2.5 {
		t.Fatalf("neg() on float = %#v, want -2.5", v)
	}

	v = BytecodeInt(5)
	(BytecodeExp{}).not(&v)
	if v.ToI() != ^int32(5) {
		t.Fatalf("not() = %#v, want bitwise inverse", v)
	}

	v = BytecodeBool(true)
	(BytecodeExp{}).blnot(&v)
	if v.ToB() {
		t.Fatalf("blnot() on true should be false, got %#v", v)
	}

	a := BytecodeInt(6)
	(BytecodeExp{}).mul(&a, BytecodeInt(7))
	if a.ToI() != 42 {
		t.Fatalf("mul() = %#v, want 42", a)
	}

	a = BytecodeFloat(6.5)
	(BytecodeExp{}).mul(&a, BytecodeInt(2))
	if a.ToF() != 13 {
		t.Fatalf("mul() float branch = %#v, want 13", a)
	}

	a = BytecodeInt(12)
	(BytecodeExp{}).div(&a, BytecodeInt(3))
	if a.ToI() != 4 {
		t.Fatalf("div() integer branch = %#v, want 4", a)
	}

	a = BytecodeFloat(9)
	(BytecodeExp{}).div(&a, BytecodeFloat(2))
	if a.ToF() != 4.5 {
		t.Fatalf("div() float branch = %#v, want 4.5", a)
	}

	a = BytecodeInt(10)
	(BytecodeExp{}).mod(&a, BytecodeInt(4))
	if a.ToI() != 2 {
		t.Fatalf("mod() = %#v, want 2", a)
	}

	a = BytecodeInt(1)
	(BytecodeExp{}).add(&a, BytecodeInt(2))
	if a.ToI() != 3 {
		t.Fatalf("add() = %#v, want 3", a)
	}

	a = BytecodeFloat(5.5)
	(BytecodeExp{}).sub(&a, BytecodeInt(1))
	if a.ToF() != 4.5 {
		t.Fatalf("sub() = %#v, want 4.5", a)
	}
}

func TestBytecodeExpDivisionAndModulusByZeroSetUndefined(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys
	sys.loader.state = LS_Complete

	a := BytecodeInt(3)
	(BytecodeExp{}).div(&a, BytecodeInt(0))
	if !a.IsUndefined() {
		t.Fatalf("div() by zero should yield undefined, got %#v", a)
	}

	a = BytecodeInt(3)
	(BytecodeExp{}).mod(&a, BytecodeInt(0))
	if !a.IsUndefined() {
		t.Fatalf("mod() by zero should yield undefined, got %#v", a)
	}
}

func TestBytecodeExpPowerHandlesFloatAndVersionedIntegerBranches(t *testing.T) {
	oldSys := sys
	defer func() { sys = oldSys }()

	sys = oldSys
	sys.loader.state = LS_Complete

	base := BytecodeInt(2)
	(BytecodeExp{}).pow(&base, BytecodeInt(3), 0)
	if base.ToI() != 8 {
		t.Fatalf("pow() integer branch = %#v, want 8", base)
	}

	base = BytecodeFloat(2)
	(BytecodeExp{}).pow(&base, BytecodeInt(3), 0)
	if base.ToF() != 8 {
		t.Fatalf("pow() float branch = %#v, want 8", base)
	}

	base = BytecodeInt(2)
	sys.cgi[0].ikemenver = [2]int{0, 0}
	sys.cgi[0].mugenver = [2]int{0, 0}
	(BytecodeExp{}).pow(&base, BytecodeInt(3), 0)
	if base.ToI() != 8 {
		t.Fatalf("pow() old-version integer branch = %#v, want 8", base)
	}
}
