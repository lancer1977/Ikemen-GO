package main

import "testing"

func TestCommandListAccessors(t *testing.T) {
	cl := NewCommandList(nil)
	cl.Add(Command{name: "fire", curbuftime: 0, buffer_shared: true})
	cl.Add(Command{name: "fire", curbuftime: 3, buffer_shared: true})
	cl.Add(Command{name: "jump", curbuftime: 1, buffer_shared: false})

	if got := cl.Get("missing"); got != nil {
		t.Fatalf("Get(missing) = %#v, want nil", got)
	}
	if got := cl.Get("fire"); len(got) != 2 {
		t.Fatalf("Get(fire) len = %d, want 2", len(got))
	}
	if !cl.GetState("fire") {
		t.Fatal("expected GetState(fire) to be true")
	}
	if cl.GetState("jump") == false {
		t.Fatal("expected GetState(jump) to be true")
	}
	if cl.Assert("missing", 9) {
		t.Fatal("expected Assert(missing) to be false")
	}
	if !cl.Assert("fire", 9) {
		t.Fatal("expected Assert(fire) to be true")
	}
	if cl.Commands[0][0].curbuftime != 9 || cl.Commands[0][1].curbuftime != 9 {
		t.Fatalf("Assert did not update all matching commands: %#v", cl.Commands[0])
	}
	cl.ClearName("fire")
	if cl.Commands[0][0].curbuftime != 0 || cl.Commands[0][1].curbuftime != 0 {
		t.Fatalf("ClearName did not clear shared commands: %#v", cl.Commands[0])
	}
	if cl.Commands[1][0].curbuftime != 1 {
		t.Fatalf("ClearName should not affect other commands: %#v", cl.Commands[1])
	}
}
