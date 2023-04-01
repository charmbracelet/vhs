package main

import (
	"os"
	"testing"
)

type parseSourceTest struct {
	tape      string
	srcTape   string
	errors    []string
	writeFile bool
}

func (st *parseSourceTest) run(t *testing.T) {
	if st.writeFile {
		err := os.WriteFile("source.tape", []byte(st.srcTape), os.ModePerm)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}

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

	os.Remove("source.tape")
}

func TestParseSource(t *testing.T) {
	t.Run("should not return errors when tape exist and is NOT empty", func(t *testing.T) {
		test := &parseSourceTest{
			tape:      "Source source.tape",
			srcTape:   `Type "echo 'Welcome to VHS!'"`,
			writeFile: true,
		}

		test.run(t)
	})

	t.Run("should return errors when tape NOT found", func(t *testing.T) {
		test := &parseSourceTest{
			tape:      "Source source.tape",
			errors:    []string{"File source.tape not found"},
			writeFile: false,
		}

		test.run(t)
	})

	t.Run("should return error when tape extension is NOT (.tape)", func(t *testing.T) {
		test := &parseSourceTest{
			tape:      "Source source.vhs",
			errors:    []string{"Expected file with .tape extension"},
			writeFile: true,
		}

		test.run(t)
	})

	t.Run("should return error when Source command does NOT have tape path", func(t *testing.T) {
		test := &parseSourceTest{
			tape:      "Source",
			errors:    []string{"Expected path after Source"},
			writeFile: true,
		}

		test.run(t)
	})

	t.Run("should return error when find nested Source commands", func(t *testing.T) {
		test := &parseSourceTest{
			tape: "Source source.tape",
			srcTape: `Type "echo 'Welcome to VHS!'"
	Source magic.tape
	Type "goodbye"
	`,
			errors:    []string{"Nested Source detected"},
			writeFile: true,
		}

		test.run(t)
	})
}
