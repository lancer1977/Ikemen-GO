package main

import "testing"

func TestBuildFFFilterGraph(t *testing.T) {
	t.Parallel()

	if got := buildFFFilterGraph(0, 720, 1280, 720, SM_Fit, SF_Bilinear); got != "" {
		t.Fatalf("invalid input should return empty graph, got %q", got)
	}

	if got := buildFFFilterGraph(640, 360, 1280, 720, SM_None, SF_Bilinear); got != "" {
		t.Fatalf("SM_None should return empty graph, got %q", got)
	}

	if got := buildFFFilterGraph(640, 360, 1280, 720, SM_Stretch, SF_Neighbor); got != "scale=1280:720:flags=neighbor,format=rgba" {
		t.Fatalf("stretch graph = %q", got)
	}

	if got := buildFFFilterGraph(640, 360, 1280, 720, SM_Fit, SF_FastBilinear); got != "scale=1280:720:flags=fast_bilinear:force_original_aspect_ratio=decrease:force_divisible_by=2,pad=1280:720:(ow-iw)/2:(oh-ih)/2:color=black,format=rgba" {
		t.Fatalf("fit graph = %q", got)
	}

	if got := buildFFFilterGraph(640, 360, 1280, 720, SM_ZoomFill, SF_Lanczos); got != "scale=1280:720:flags=lanczos:force_original_aspect_ratio=increase:force_divisible_by=2,crop=1280:720:floor((iw-1280)/2):floor((ih-720)/2),format=rgba" {
		t.Fatalf("zoomfill graph = %q", got)
	}
}
