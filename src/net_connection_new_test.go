package main

import "testing"

func TestNewNetConnection(t *testing.T) {
	t.Parallel()

	nc := NewNetConnection()
	if nc == nil {
		t.Fatal("NewNetConnection returned nil")
	}
	if nc.st != NS_Stop {
		t.Fatalf("st = %v, want NS_Stop", nc.st)
	}
	if nc.sendEnd == nil || nc.recvEnd == nil || nc.closing == nil {
		t.Fatal("expected channels to be initialized")
	}
	select {
	case <-nc.sendEnd:
	default:
		t.Fatal("expected sendEnd to be pre-seeded")
	}
	select {
	case <-nc.recvEnd:
	default:
		t.Fatal("expected recvEnd to be pre-seeded")
	}
	if nc.buf[0].InputReader == nil {
		t.Fatal("expected net buffers to be initialized")
	}
}
