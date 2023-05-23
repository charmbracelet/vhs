package main

import (
	"testing"
)

type parseScreenshotTest struct {
	tape   string
	errors []string
}

func (st *parseScreenshotTest) run(t *testing.T) {
	l := NewLexer(st.tape)
	p := NewParser(l)

	_ = p.Parse()

	if len(p.errors) != len(st.errors) {
		t.Fatalf("Expected errors: %d, errors: %d", len(st.errors), len(p.errors))
	}

	for i := range st.errors {
		err := p.errors[i].Msg
		expected := st.errors[i]

		if err != expected {
			t.Errorf("Expected error: %s, actual error %s", expected, err)
		}
	}
}

func TestParseScreeenshot(t *testing.T) {
	t.Run("should return error when screenshot extension is NOT (.png)", func(t *testing.T) {
		test := &parseScreenshotTest{
			tape:   "Screenshot step_one_screenshot.jpg",
			errors: []string{"Expected file with .png extension"},
		}

		test.run(t)
	})

	t.Run("should return error when screenshot path is missing", func(t *testing.T) {
		test := &parseScreenshotTest{
			tape:   "Screenshot",
			errors: []string{"Expected path after Screenshot"},
		}

		test.run(t)
	})
}
