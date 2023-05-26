package main

import (
	"testing"
)

func TestFindAllThemes(t *testing.T) {
	themes, err := sortedThemeNames()
	if err != nil {
		t.Fatal(err)
	}
	expect := 295
	if l := len(themes); l != expect {
		t.Errorf("expected to load %d themes, got %d", expect, l)
	}
}

func TestFindTheme(t *testing.T) {
	_, err := findTheme("Catppuccin Latte")
	if err != nil {
		t.Error(err)
	}

	theme, err := findTheme("caTppuccin ltt")
	te, ok := err.(ThemeNotFoundError)
	if !ok {
		t.Fatal("expected to not be found:", theme)
	}
	if len(te.Suggestions) != 1 {
		t.Fatal("expected 1 suggestion, got:", te.Suggestions)
	}
	if sg := te.Suggestions[0]; sg != "Catppuccin Latte" {
		t.Fatal("wrong suggestion:", te.Suggestions[0])
	}
}
