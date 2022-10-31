//go:generate make all
package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
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

// sortedThemeNames returns the names of the themes, sorted.
func sortedThemeNames() []string {
	var keys []string
	for _, bts := range [][]byte{themesBts, customThemesBts} {
		for _, theme := range parseThemes(bts) {
			keys = append(keys, theme.Name)
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		return strings.ToLower(keys[i]) < strings.ToLower(keys[j])
	})
	return keys
}

const distance = 2

// findTheme return the given theme, if it exists.
func findTheme(name string) (Theme, []string, bool) {
	for _, bts := range [][]byte{themesBts, customThemesBts} {
		for _, theme := range parseThemes(bts) {
			if theme.Name == name {
				return theme, nil, true
			}
		}
	}

	// not found, lets find similar themes!
	keys := sortedThemeNames()

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
	return Theme{}, suggestions, false
}

func parseThemes(bts []byte) []Theme {
	var themes []Theme
	if err := json.Unmarshal(bts, &themes); err != nil {
		fmt.Fprintf(os.Stderr, "could not load themes.json: %s\n", err)
		os.Exit(1)
	}
	return themes
}
