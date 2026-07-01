package main

import "testing"

func TestEffectiveFightAspectKey(t *testing.T) {
	t.Run("stage", func(t *testing.T) {
		got, err := effectiveFightAspectKey(map[string]string{
			"Video.FightAspectWidth":  "-1",
			"Video.FightAspectHeight": "-1",
		})
		if err != nil {
			t.Fatalf("effectiveFightAspectKey returned error: %v", err)
		}
		if got != "stage" {
			t.Fatalf("effectiveFightAspectKey = %q, want stage", got)
		}
	})

	t.Run("resolution", func(t *testing.T) {
		got, err := effectiveFightAspectKey(map[string]string{
			"Video.FightAspectWidth":  "0",
			"Video.FightAspectHeight": "0",
			"Video.GameWidth":         "1280",
			"Video.GameHeight":        "720",
		})
		if err != nil {
			t.Fatalf("effectiveFightAspectKey returned error: %v", err)
		}
		if got != "resolution:16:9" {
			t.Fatalf("effectiveFightAspectKey = %q, want resolution:16:9", got)
		}
	})

	t.Run("custom", func(t *testing.T) {
		got, err := effectiveFightAspectKey(map[string]string{
			"Video.FightAspectWidth":  "16",
			"Video.FightAspectHeight": "9",
		})
		if err != nil {
			t.Fatalf("effectiveFightAspectKey returned error: %v", err)
		}
		if got != "custom:16:9" {
			t.Fatalf("effectiveFightAspectKey = %q, want custom:16:9", got)
		}
	})

	t.Run("missing", func(t *testing.T) {
		if _, err := effectiveFightAspectKey(map[string]string{
			"Video.FightAspectWidth": "16",
		}); err == nil {
			t.Fatalf("effectiveFightAspectKey succeeded with missing height")
		}
	})
}
