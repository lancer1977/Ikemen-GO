package main

import (
	"os"
	"testing"
)

func TestReplayHeaderIO(t *testing.T) {
	t.Parallel()

	tmp, err := os.CreateTemp("", "replay-header-*.rep")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()

	want := &ReplayHeader{
		FormatVersion:      7,
		SyncVersion:        12,
		Strict:             []SyncSetting{{Path: "A", Value: "1"}},
		Host:               []SyncSetting{{Path: "B", Value: "2"}},
		ContentFingerprint: "abc123",
	}

	if err := writeReplayHeader(tmp, want); err != nil {
		t.Fatalf("writeReplayHeader returned error: %v", err)
	}
	if got, err := readReplayHeader(tmp); err != nil {
		t.Fatalf("readReplayHeader returned error: %v", err)
	} else if got == nil {
		t.Fatal("readReplayHeader returned nil")
	} else {
		if got.FormatVersion != want.FormatVersion || got.SyncVersion != want.SyncVersion || got.ContentFingerprint != want.ContentFingerprint {
			t.Fatalf("header fields mismatch: got %#v want %#v", got, want)
		}
		if len(got.Strict) != 1 || got.Strict[0] != want.Strict[0] {
			t.Fatalf("Strict mismatch: %#v", got.Strict)
		}
		if len(got.Host) != 1 || got.Host[0] != want.Host[0] {
			t.Fatalf("Host mismatch: %#v", got.Host)
		}
	}
}

func TestWriteReplayHeaderNil(t *testing.T) {
	t.Parallel()

	tmp, err := os.CreateTemp("", "replay-header-nil-*.rep")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()

	if err := writeReplayHeader(tmp, nil); err != nil {
		t.Fatalf("writeReplayHeader(nil) returned error: %v", err)
	}
	info, err := tmp.Stat()
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if info.Size() != 0 {
		t.Fatalf("writeReplayHeader(nil) wrote %d bytes, want 0", info.Size())
	}
}

func TestReadReplayHeaderFallbacks(t *testing.T) {
	t.Parallel()

	tmp, err := os.CreateTemp("", "replay-header-fallback-*.rep")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()

	if _, err := tmp.Write([]byte("not-a-replay")); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if got, err := readReplayHeader(tmp); err != nil {
		t.Fatalf("readReplayHeader(bad magic) returned error: %v", err)
	} else if got != nil {
		t.Fatalf("readReplayHeader(bad magic) = %#v, want nil", got)
	}

	if err := tmp.Truncate(0); err != nil {
		t.Fatalf("Truncate: %v", err)
	}
	if _, err := tmp.Seek(0, 0); err != nil {
		t.Fatalf("Seek: %v", err)
	}
	if _, err := tmp.Write([]byte(replayMagic[:len(replayMagic)-1])); err != nil {
		t.Fatalf("Write short magic: %v", err)
	}
	if got, err := readReplayHeader(tmp); err != nil {
		t.Fatalf("readReplayHeader(short file) returned error: %v", err)
	} else if got != nil {
		t.Fatalf("readReplayHeader(short file) = %#v, want nil", got)
	}
}
