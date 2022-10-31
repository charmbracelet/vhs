package main

import (
	"html/template"
	"strings"
)

// ExecuteTemplate executes the template with the given tape and returns the output.
func ExecuteTemplate(tape string) (string, error) {
	parsed, err := template.
		New("Tape").
		Funcs(template.FuncMap{
			"add": func(a, b int) int {
				return a + b
			},
			"repeat": func(s string, n int) string {
				return strings.Repeat(s, n)
			},
		}).
		Parse(tape)

	if err != nil {
		return "", err
	}

	var output strings.Builder

	err = parsed.Execute(&output, nil)

	if err != nil {
		return "", err
	}

	return output.String(), nil
}
