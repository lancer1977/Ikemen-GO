package main

import (
	"testing"
	"time"
)

func TestLiveTimestampUTC_UsesRFC3339Nano(t *testing.T) {
	ts := liveTimestampUTC()
	if _, err := time.Parse(time.RFC3339Nano, ts); err != nil {
		t.Fatalf("expected RFC3339Nano timestamp, got %q: %v", ts, err)
	}
}
