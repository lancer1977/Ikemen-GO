package main

import (
	"math"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestLuaNumberToJsonFormatsNumericRangesAndRejectsInvalidValues(t *testing.T) {
	if got, err := luaNumberToJson(lua.LNumber(12)); err != nil || got != int64(12) {
		t.Fatalf("luaNumberToJson(int-like) = (%T) %#v, %v; want int64(12), nil", got, got, err)
	}
	if got, err := luaNumberToJson(lua.LNumber(12.5)); err != nil || got != float64(12.5) {
		t.Fatalf("luaNumberToJson(float) = (%T) %#v, %v; want float64(12.5), nil", got, got, err)
	}
	// Use uint64(3) << 62 = 3*2^62, which is:
	// - Well outside int64 max (9.223e18 < 13.835e18)
	// - Representable exactly in float64 (multiple of a power of 2)
	testVal := uint64(3) << 62
	if got, err := luaNumberToJson(lua.LNumber(float64(testVal))); err != nil || got != testVal {
		t.Fatalf("luaNumberToJson(uint64-range) = (%T) %#v, %v; want uint64(3<<62), nil", got, got, err)
	}
	if _, err := luaNumberToJson(lua.LNumber(math.NaN())); err == nil {
		t.Fatal("luaNumberToJson should reject NaN")
	}
}
