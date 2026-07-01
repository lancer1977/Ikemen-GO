package main

import "testing"

func TestBytecodeStackPushPopAndClear(t *testing.T) {
	var bs BytecodeStack

	bs.Push(BytecodeInt(1))
	bs.PushI(2)
	bs.PushI64(3)
	bs.PushF(4.5)
	bs.PushB(true)

	if len(bs) != 5 {
		t.Fatalf("unexpected stack length: %d", len(bs))
	}
	if got := bs.Top(); got == nil || !got.ToB() {
		t.Fatalf("Top() should reference last pushed value, got %#v", got)
	}

	if got := bs.Pop(); got.ToB() != true || got.ToI() != 1 {
		t.Fatalf("Pop() = %#v, want last pushed bool", got)
	}
	if got := bs.Pop(); got.ToF() != 4.5 {
		t.Fatalf("second Pop() = %#v, want 4.5", got)
	}
	if got := bs.Pop(); got.ToI64() != 3 {
		t.Fatalf("third Pop() = %#v, want 3", got)
	}
	if got := bs.Pop(); got.ToI() != 2 {
		t.Fatalf("fourth Pop() = %#v, want 2", got)
	}
	if got := bs.Pop(); got.ToI() != 1 {
		t.Fatalf("fifth Pop() = %#v, want 1", got)
	}

	bs.Push(BytecodeInt(9))
	bs.Clear()
	if len(bs) != 0 {
		t.Fatalf("Clear() should empty stack, got len=%d", len(bs))
	}
}

func TestBytecodeStackEmptyTopAndPopPanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic from empty BytecodeStack access")
		}
	}()

	var bs BytecodeStack
	_ = bs.Top()
}

func TestBytecodeStackEmptyPopPanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic from empty BytecodeStack pop")
		}
	}()

	var bs BytecodeStack
	_ = bs.Pop()
}
