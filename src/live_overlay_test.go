package main

import (
	"testing"
)

func TestParseLiveOverlayFile_DefaultSchemaAndEffects(t *testing.T) {
	file, err := parseLiveOverlayFile([]byte(`{
		"sourceSequence": 7,
		"overlays": [
			{
				"id": "cmd-text-1",
				"kind": "text",
				"text": "GG!",
				"x": 160,
				"y": 48,
				"durationMs": 1200
			},
			{
				"id": "cmd-image-1",
				"kind": "image",
				"assetPath": "assets/emotes/hype.png",
				"x": 80,
				"y": 90,
				"durationMs": 900
			}
		]
	}`))
	if err != nil {
		t.Fatalf("parseLiveOverlayFile returned error: %v", err)
	}
	if file.Schema != "live-lancero/overlay-state/v1" {
		t.Fatalf("unexpected schema: %q", file.Schema)
	}
	if file.SourceSequence != 7 {
		t.Fatalf("unexpected source sequence: %d", file.SourceSequence)
	}
	if len(file.Overlays) != 2 {
		t.Fatalf("unexpected overlay count: %d", len(file.Overlays))
	}
	if file.Overlays[0].Kind != "text" || file.Overlays[0].Text != "GG!" {
		t.Fatalf("unexpected first overlay: %#v", file.Overlays[0])
	}
	if file.Overlays[1].Kind != "image" || file.Overlays[1].AssetPath != "assets/emotes/hype.png" {
		t.Fatalf("unexpected second overlay: %#v", file.Overlays[1])
	}
}

func TestEffectDurationFrames_DefaultsToSixtyFps(t *testing.T) {
	prev := sys.cfg.Video.Framerate
	sys.cfg.Video.Framerate = 0
	t.Cleanup(func() {
		sys.cfg.Video.Framerate = prev
	})

	if got := effectDurationFrames(1000); got != 60 {
		t.Fatalf("expected 60 frames for 1000ms, got %d", got)
	}
	if got := effectDurationFrames(0); got != 72 {
		t.Fatalf("expected default duration of 72 frames, got %d", got)
	}
}
