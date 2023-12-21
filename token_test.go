package main

import "testing"

func TestToCamel(t *testing.T) {
	// simple case
	res := toCamel("SIMPLE")
	if res != "Simple" {
		t.Errorf("expected Simple, got %s", res)
	}

	// multiple words
	res = toCamel("MULTIPLE_WORDS")
	if res != "MultipleWords" {
		t.Errorf("expected MultipleWords, got %s", res)
	}
}
