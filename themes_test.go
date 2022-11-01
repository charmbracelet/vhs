package main

import "testing"

func TestFindAllThemes(t *testing.T) {
	themes := sortedThemeNames()
	expect := 295
	if l := len(themes); l != expect {
		t.Errorf("expected to load %d themes, got %d", expect, l)
	}
}

func TestFindTheme(t *testing.T) {
	theme, suggestions, ok := findTheme("caTppuccin ltt")
	if ok {
		t.Fatal("expected to not be found:", theme)
	}
	if len(suggestions) != 1 {
		t.Fatal("expected 1 suggestions, got:", suggestions)
	}
	if sg := suggestions[0]; sg != "Catppuccin Latte" {
		t.Fatal("wrong suggestion:", suggestions[0])
	}
}
