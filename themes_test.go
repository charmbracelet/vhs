package main

import "testing"

func TestFindAllThemes(t *testing.T) {
	themes := sortedThemeNames()
	expect := 295
	if l := len(themes); l != expect {
		t.Errorf("expected to load %d themes, got %d", expect, l)
	}
}
