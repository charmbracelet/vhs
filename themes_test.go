package main

import (
	"errors"
	"reflect"
	"sort"
	"testing"
)

func TestFindAllThemes(t *testing.T) {
	themes, err := sortedThemeNames()
	if err != nil {
		t.Fatal(err)
	}
	expect := 348
	if l := len(themes); l != expect {
		t.Errorf("expected to load %d themes, got %d", expect, l)
	}
}

func TestFindTheme(t *testing.T) {
	tests := []struct {
		tname string
		theme string
		err   error
	}{
		{
			tname: "exact match",
			theme: "Catppuccin Latte",
			err:   nil,
		},
		{
			tname: "match found",
			theme: "caTppuccin ltt",
			err:   ThemeNotFoundError{"caTppuccin ltt", []string{"Catppuccin Latte"}},
		},
		{
			tname: "no match found",
			theme: "stArf1sh",
			err:   ThemeNotFoundError{"stArf1sh", []string{}},
		},
		{
			tname: "single char",
			theme: "s",
			err:   ThemeNotFoundError{"s", []string{}},
		},
		{
			tname: "empty string",
			theme: "",
			err:   ThemeNotFoundError{"", []string{}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.tname, func(t *testing.T) {
			_, err := findTheme(tc.theme)
			if tc.err != nil {
				if err == nil {
					t.Fatal("expected an error:", tc.err)
				}
				// check we got the right error
				var perr ThemeNotFoundError
				if !errors.As(err, &perr) {
					t.Fatal(err)
				}
				gotErr := err.(ThemeNotFoundError)
				wantErr := tc.err.(ThemeNotFoundError)
				// check suggestions
				sort.Strings(gotErr.Suggestions)
				sort.Strings(wantErr.Suggestions)
				if !reflect.DeepEqual(gotErr.Suggestions, wantErr.Suggestions) {
					t.Fatalf("got != want. got: %v, want: %v", err, tc.err)
				}
				// check names
				if !reflect.DeepEqual(gotErr.Theme, wantErr.Theme) {
					t.Fatalf("got != want. got: %v, want: %v", err, tc.err)
				}
			}
			if err != nil && tc.err == nil {
				t.Fatal("unexpected error:", err)
			}
		})
	}
}
