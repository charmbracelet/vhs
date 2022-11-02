//go:generate make all
package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/agnivade/levenshtein"
)

var (
	//go:embed themes.json
	themesBts []byte

	//go:embed themes_custom.json
	customThemesBts []byte
)

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
	for _, bts := range [][]byte{themesBts, customThemesBts} {
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
	for _, bts := range [][]byte{themesBts, customThemesBts} {
		themes, err := parseThemes(bts)
		if err != nil {
			return Theme{}, err
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
		return Theme{}, err
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
	return Theme{}, ThemeNotFoundError{name, suggestions}
}

func parseThemes(bts []byte) ([]Theme, error) {
	var themes []Theme
	if err := json.Unmarshal(bts, &themes); err != nil {
		return nil, fmt.Errorf("could not load themes.json: %w", err)
	}
	return themes, nil
}
