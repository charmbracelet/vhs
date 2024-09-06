// Package main theme.go contains the information about a terminal theme.
// It stores the 16 base colors as well as the background and foreground colors
// of the terminal theme.
//
// It can be changed through the Set command.
//
// Set Theme {"background": "#171717"}
// Set Theme "Catppuccin Mocha"
//
//go:generate make all
package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/glamour/ansi"

	_ "embed"

	"github.com/agnivade/levenshtein"
)

// Theme is a terminal theme for xterm.js
// It is used for marshalling between the xterm.js readable json format and a
// valid go struct.
// https://xtermjs.org/docs/api/terminal/interfaces/itheme/
type Theme struct {
	Name          string `json:"name"`
	Background    string `json:"background"`
	Foreground    string `json:"foreground"`
	Selection     string `json:"selection"`
	Cursor        string `json:"cursor"`
	CursorAccent  string `json:"cursorAccent"`
	Black         string `json:"black"`
	BrightBlack   string `json:"brightBlack"`
	Red           string `json:"red"`
	BrightRed     string `json:"brightRed"`
	Green         string `json:"green"`
	BrightGreen   string `json:"brightGreen"`
	Yellow        string `json:"yellow"`
	BrightYellow  string `json:"brightYellow"`
	Blue          string `json:"blue"`
	BrightBlue    string `json:"brightBlue"`
	Magenta       string `json:"magenta"`
	BrightMagenta string `json:"brightMagenta"`
	Cyan          string `json:"cyan"`
	BrightCyan    string `json:"brightCyan"`
	White         string `json:"white"`
	BrightWhite   string `json:"brightWhite"`
}

func (t Theme) String() string {
	ts, err := json.Marshal(t)
	if err != nil {
		dts, _ := json.Marshal(DefaultTheme)
		return string(dts)
	}
	return string(ts)
}

// DefaultTheme is the default theme to use for recording demos and
// screenshots.
//
// Taken from https://github.com/meowgorithm/dotfiles.
var DefaultTheme = Theme{
	Background:    Background,
	Foreground:    Foreground,
	Cursor:        Foreground,
	CursorAccent:  Background,
	Black:         Black,
	BrightBlack:   BrightBlack,
	Red:           Red,
	BrightRed:     BrightRed,
	Green:         Green,
	BrightGreen:   BrightGreen,
	Yellow:        Yellow,
	BrightYellow:  BrightYellow,
	Blue:          Blue,
	BrightBlue:    BrightBlue,
	Magenta:       Magenta,
	BrightMagenta: BrightMagenta,
	Cyan:          Cyan,
	BrightCyan:    BrightCyan,
	White:         White,
	BrightWhite:   BrightWhite,
}

const margin = 2

// GlamourTheme is the theme for printing out the manual page.
// $ vhs man
var GlamourTheme = ansi.StyleConfig{
	Document: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			BlockPrefix: "\n",
			BlockSuffix: "\n",
		},
		Margin: uintPtr(margin),
	},
	Heading: ansi.StyleBlock{
		StylePrimitive: ansi.StylePrimitive{
			BlockSuffix: "\n",
			Color:       stringPtr("99"),
			Bold:        boolPtr(true),
		},
	},
	Item:     ansi.StylePrimitive{Prefix: "· "},
	Emph:     ansi.StylePrimitive{Color: stringPtr(BrightBlack)},
	Strong:   ansi.StylePrimitive{Bold: boolPtr(true)},
	Link:     ansi.StylePrimitive{Color: stringPtr("42"), Underline: boolPtr(true)},
	LinkText: ansi.StylePrimitive{Color: stringPtr("207")},
	Code:     ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Color: stringPtr("204")}},
}

func boolPtr(b bool) *bool       { return &b }
func stringPtr(s string) *string { return &s }
func uintPtr(u uint) *uint       { return &u }

//go:embed themes.json
var themesBts []byte

// ThemeNotFoundError is returned when a requested theme is not found.
type ThemeNotFoundError struct {
	Theme       string
	Suggestions []string
}

func (e ThemeNotFoundError) Error() string {
	if len(e.Suggestions) == 0 {
		return fmt.Sprintf("invalid `Set Theme %q`: theme does not exist", e.Theme)
	}

	return fmt.Sprintf("invalid `Set Theme %q`: did you mean %q",
		e.Theme,
		strings.Join(e.Suggestions, ", "),
	)
}

// sortedThemeNames returns the names of the themes, sorted.
func sortedThemeNames() ([]string, error) {
	var keys []string
	for _, bts := range [][]byte{themesBts} {
		themes, err := parseThemes(bts)
		if err != nil {
			return nil, err
		}

		for _, theme := range themes {
			keys = append(keys, theme.Name)
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		return strings.ToLower(keys[i]) < strings.ToLower(keys[j])
	})
	return keys, nil
}

const distance = 2

// findTheme return the given theme, if it exists.
func findTheme(name string) (Theme, error) {
	for _, bts := range [][]byte{themesBts} {
		themes, err := parseThemes(bts)
		if err != nil {
			return DefaultTheme, err
		}

		for _, theme := range themes {
			if theme.Name == name {
				return theme, nil
			}
		}
	}

	// not found, lets find similar themes!
	keys, err := sortedThemeNames()
	if err != nil {
		return DefaultTheme, err
	}

	suggestions := []string{}
	lname := strings.ToLower(name)
	for _, theme := range keys {
		ltheme := strings.ToLower(theme)
		levenshteinDistance := levenshtein.ComputeDistance(lname, ltheme)
		suggestByLevenshtein := levenshteinDistance <= distance
		suggestByPrefix := strings.HasPrefix(lname, ltheme)
		if suggestByLevenshtein || suggestByPrefix {
			suggestions = append(suggestions, theme)
		}
	}
	return DefaultTheme, ThemeNotFoundError{name, suggestions}
}

func parseThemes(bts []byte) ([]Theme, error) {
	var themes []Theme
	if err := json.Unmarshal(bts, &themes); err != nil {
		return nil, fmt.Errorf("could not load themes.json: %w", err)
	}
	return themes, nil
}
