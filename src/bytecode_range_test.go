package main

import "testing"

func TestBytecodeExpRangeCheckHandlesInclusiveExclusiveAndTypeBranches(t *testing.T) {
	ops := BytecodeExp{}

	v := BytecodeInt(5)
	ops.rangeCheck(&v, BytecodeInt(5), BytecodeInt(5), OC_range_ii)
	if !v.ToB() {
		t.Fatalf("rangeCheck() inclusive int branch = %#v, want true", v)
	}

	v = BytecodeInt(5)
	ops.rangeCheck(&v, BytecodeInt(5), BytecodeInt(6), OC_range_ie)
	if !v.ToB() {
		t.Fatalf("rangeCheck() lower-inclusive int branch = %#v, want true", v)
	}

	v = BytecodeInt(5)
	ops.rangeCheck(&v, BytecodeInt(4), BytecodeInt(5), OC_range_ei)
	if !v.ToB() {
		t.Fatalf("rangeCheck() upper-inclusive int branch = %#v, want true", v)
	}

	v = BytecodeFloat(5.5)
	ops.rangeCheck(&v, BytecodeFloat(5.0), BytecodeFloat(6.0), OC_range_ee)
	if !v.ToB() {
		t.Fatalf("rangeCheck() exclusive float branch = %#v, want true", v)
	}

	v = BytecodeFloat(5.0)
	ops.rangeCheck(&v, BytecodeFloat(5.0), BytecodeFloat(6.0), OC_range_ee)
	if v.ToB() {
		t.Fatalf("rangeCheck() exclusive float branch should fail on lower bound, got %#v", v)
	}
}

func TestBytecodeExpRangeCheckPropagatesUndefined(t *testing.T) {
	ops := BytecodeExp{}

	v := BytecodeUndefined()
	ops.rangeCheck(&v, BytecodeInt(0), BytecodeInt(1), OC_range_ii)
	if !v.IsUndefined() {
		t.Fatalf("rangeCheck() with undefined value should remain undefined, got %#v", v)
	}
}
