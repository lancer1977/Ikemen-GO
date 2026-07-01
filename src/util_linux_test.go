//go:build linux

package main

import "testing"

func TestOsPreferredLanguage(t *testing.T) {
	t.Setenv("LC_ALL", "")
	t.Setenv("LC_MESSAGES", "")
	t.Setenv("LANG", "")
	if got := osPreferredLanguage(); got != "" {
		t.Fatalf("osPreferredLanguage() with no env = %q, want empty", got)
	}

	t.Setenv("LANG", "en_US.UTF-8")
	if got := osPreferredLanguage(); got != "en_US.UTF-8" {
		t.Fatalf("osPreferredLanguage() with LANG = %q, want en_US.UTF-8", got)
	}

	t.Setenv("LC_MESSAGES", "pt_BR.UTF-8")
	if got := osPreferredLanguage(); got != "pt_BR.UTF-8" {
		t.Fatalf("osPreferredLanguage() should prefer LC_MESSAGES, got %q", got)
	}

	t.Setenv("LC_ALL", "ja_JP.UTF-8")
	if got := osPreferredLanguage(); got != "ja_JP.UTF-8" {
		t.Fatalf("osPreferredLanguage() should prefer LC_ALL, got %q", got)
	}
}
