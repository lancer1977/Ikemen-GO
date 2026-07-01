package main

import "testing"

func TestSoundChannelsSetSize(t *testing.T) {
	channels := SoundChannels{
		{playerID: 3, channelNo: 4, group: 5, number: 6, timeStamp: 7},
		{playerID: 8, channelNo: 9, group: 10, number: 11, timeStamp: 12},
	}

	channels.SetSize(4)
	if len(channels) != 4 {
		t.Fatalf("SetSize grow len = %d, want 4", len(channels))
	}
	if channels[0].playerID != 3 || channels[1].playerID != 8 {
		t.Fatalf("SetSize should keep existing channels: %#v", channels[:2])
	}
	if channels[2].playerID != -1 || channels[2].channelNo != -1 || channels[3].group != -1 || channels[3].number != -1 {
		t.Fatalf("SetSize should initialize new channels: %#v", channels[2:])
	}

	channels[0].playerID = 13
	channels[1].playerID = 14
	channels[2].playerID = 15
	channels.SetSize(2)
	if len(channels) != 2 {
		t.Fatalf("SetSize shrink len = %d, want 2", len(channels))
	}
	if channels[0].playerID != 13 || channels[1].playerID != 14 {
		t.Fatalf("SetSize should keep first entries when shrinking: %#v", channels)
	}
}
