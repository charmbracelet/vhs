package parser

import (
	"os"
	"strings"
	"testing"

	"github.com/charmbracelet/vhs/lexer"
	"github.com/charmbracelet/vhs/token"
)

func TestParser(t *testing.T) {
	input := `
Set TypingSpeed 100ms
Set WaitTimeout 1m
Set WaitPattern /foo/
Type "echo 'Hello, World!'"
Enter
Backspace@0.1 5
Backspace@.1 5
Backspace@1 5
Backspace@100ms 5
Delete 2
Insert 2
Right 3
Left 3
Up@50ms
Down 2
Ctrl+C
Ctrl+L
Alt+.
Sleep 100ms
Sleep 3
Wait
Wait+Screen
Wait@100ms /foobar/`

	expected := []Command{
		{Type: token.SET, Options: "TypingSpeed", Args: "100ms"},
		{Type: token.SET, Options: "WaitTimeout", Args: "1m"},
		{Type: token.SET, Options: "WaitPattern", Args: "foo"},
		{Type: token.TYPE, Options: "", Args: "echo 'Hello, World!'"},
		{Type: token.ENTER, Options: "", Args: "1"},
		{Type: token.BACKSPACE, Options: "0.1s", Args: "5"},
		{Type: token.BACKSPACE, Options: ".1s", Args: "5"},
		{Type: token.BACKSPACE, Options: "1s", Args: "5"},
		{Type: token.BACKSPACE, Options: "100ms", Args: "5"},
		{Type: token.DELETE, Options: "", Args: "2"},
		{Type: token.INSERT, Options: "", Args: "2"},
		{Type: token.RIGHT, Options: "", Args: "3"},
		{Type: token.LEFT, Options: "", Args: "3"},
		{Type: token.UP, Options: "50ms", Args: "1"},
		{Type: token.DOWN, Options: "", Args: "2"},
		{Type: token.CTRL, Options: "", Args: "C"},
		{Type: token.CTRL, Options: "", Args: "L"},
		{Type: token.ALT, Options: "", Args: "."},
		{Type: token.SLEEP, Args: "100ms"},
		{Type: token.SLEEP, Args: "3s"},
		{Type: token.WAIT, Args: "Line"},
		{Type: token.WAIT, Args: "Screen"},
		{Type: token.WAIT, Options: "100ms", Args: "Line foobar"},
	}

	l := lexer.New(input)
	p := New(l)

	cmds := p.Parse()

	if len(cmds) != len(expected) {
		t.Fatalf("Expected %d commands, got %d; %v", len(expected), len(cmds), cmds)
	}

	for i, cmd := range cmds {
		if cmd.Type != expected[i].Type {
			t.Errorf("Expected command %d to be %s, got %s", i, expected[i].Type, cmd.Type)
		}
		if cmd.Args != expected[i].Args {
			t.Errorf("Expected command %d to have args %s, got %s", i, expected[i].Args, cmd.Args)
		}
		if cmd.Options != expected[i].Options {
			t.Errorf("Expected command %d to have options %s, got %s", i, expected[i].Options, cmd.Options)
		}
	}
}

func TestParserErrors(t *testing.T) {
	input := `
Type Enter
Type "echo 'Hello, World!'" Enter
Foo
Sleep Bar`

	l := lexer.New(input)
	p := New(l)

	_ = p.Parse()

	expectedErrors := []string{
		" 2:6  │ Type expects string",
		" 4:1  │ Invalid command: Foo",
		" 5:1  │ Expected time after Sleep",
		" 5:7  │ Invalid command: Bar",
	}

	if len(p.errors) != len(expectedErrors) {
		t.Fatalf("Expected %d errors, got %d", len(expectedErrors), len(p.errors))
	}

	for i, err := range p.errors {
		if err.String() != expectedErrors[i] {
			t.Errorf("Expected error %d to be [%s], got (%s)", i, expectedErrors[i], err)
		}
	}
}

func TestParseTapeFile(t *testing.T) {
	input, err := os.ReadFile("../examples/fixtures/all.tape")
	if err != nil {
		t.Fatal("could not read fixture file")
	}

	expected := []Command{
		{Type: token.OUTPUT, Options: ".gif", Args: "examples/fixtures/all.gif"},
		{Type: token.OUTPUT, Options: ".mp4", Args: "examples/fixtures/all.mp4"},
		{Type: token.OUTPUT, Options: ".webm", Args: "examples/fixtures/all.webm"},
		{Type: token.SET, Options: "Shell", Args: "fish"},
		{Type: token.SET, Options: "FontSize", Args: "22"},
		{Type: token.SET, Options: "FontFamily", Args: "DejaVu Sans Mono"},
		{Type: token.SET, Options: "Height", Args: "600"},
		{Type: token.SET, Options: "Width", Args: "1200"},
		{Type: token.SET, Options: "LetterSpacing", Args: "1"},
		{Type: token.SET, Options: "LineHeight", Args: "1.2"},
		{Type: token.SET, Options: "Theme", Args: "{ \"name\": \"Whimsy\", \"black\": \"#535178\", \"red\": \"#ef6487\", \"green\": \"#5eca89\", \"yellow\": \"#fdd877\", \"blue\": \"#65aef7\", \"purple\": \"#aa7ff0\", \"cyan\": \"#43c1be\", \"white\": \"#ffffff\", \"brightBlack\": \"#535178\", \"brightRed\": \"#ef6487\", \"brightGreen\": \"#5eca89\", \"brightYellow\": \"#fdd877\", \"brightBlue\": \"#65aef7\", \"brightPurple\": \"#aa7ff0\", \"brightCyan\": \"#43c1be\", \"brightWhite\": \"#ffffff\", \"background\": \"#29283b\", \"foreground\": \"#b3b0d6\", \"selectionBackground\": \"#3d3c58\", \"cursorColor\": \"#b3b0d6\" }"},
		{Type: token.SET, Options: "Theme", Args: "Catppuccin Mocha"},
		{Type: token.SET, Options: "Padding", Args: "50"},
		{Type: token.SET, Options: "Framerate", Args: "60"},
		{Type: token.SET, Options: "PlaybackSpeed", Args: "2"},
		{Type: token.SET, Options: "TypingSpeed", Args: ".1s"},
		{Type: token.SET, Options: "LoopOffset", Args: "60.4%"},
		{Type: token.SET, Options: "LoopOffset", Args: "20.99%"},
		{Type: token.SET, Options: "CursorBlink", Args: "false"},
		{Type: token.SLEEP, Options: "", Args: "1s"},
		{Type: token.SLEEP, Options: "", Args: "500ms"},
		{Type: token.SLEEP, Options: "", Args: ".5s"},
		{Type: token.SLEEP, Options: "", Args: "0.5s"},
		{Type: token.TYPE, Options: ".5s", Args: "All"},
		{Type: token.TYPE, Options: "500ms", Args: "All"},
		{Type: token.TYPE, Options: "", Args: "Double Quote"},
		{Type: token.TYPE, Options: "", Args: "\"Single\" Quote"},
		{Type: token.TYPE, Options: "", Args: `"Backtick" 'Quote'`},
		{Type: token.BACKSPACE, Options: "", Args: "1"},
		{Type: token.BACKSPACE, Options: "", Args: "2"},
		{Type: token.BACKSPACE, Options: "1s", Args: "3"},
		{Type: token.DELETE, Options: "", Args: "1"},
		{Type: token.DELETE, Options: "", Args: "2"},
		{Type: token.DELETE, Options: "1s", Args: "3"},
		{Type: token.INSERT, Options: "", Args: "1"},
		{Type: token.INSERT, Options: "", Args: "2"},
		{Type: token.INSERT, Options: "1s", Args: "3"},
		{Type: token.DOWN, Options: "", Args: "1"},
		{Type: token.DOWN, Options: "", Args: "2"},
		{Type: token.DOWN, Options: "1s", Args: "3"},
		{Type: token.PAGE_DOWN, Options: "", Args: "1"},
		{Type: token.PAGE_DOWN, Options: "", Args: "2"},
		{Type: token.PAGE_DOWN, Options: "1s", Args: "3"},
		{Type: token.ENTER, Options: "", Args: "1"},
		{Type: token.ENTER, Options: "", Args: "2"},
		{Type: token.ENTER, Options: "1s", Args: "3"},
		{Type: token.SPACE, Options: "", Args: "1"},
		{Type: token.SPACE, Options: "", Args: "2"},
		{Type: token.SPACE, Options: "1s", Args: "3"},
		{Type: token.TAB, Options: "", Args: "1"},
		{Type: token.TAB, Options: "", Args: "2"},
		{Type: token.TAB, Options: "1s", Args: "3"},
		{Type: token.LEFT, Options: "", Args: "1"},
		{Type: token.LEFT, Options: "", Args: "2"},
		{Type: token.LEFT, Options: "1s", Args: "3"},
		{Type: token.RIGHT, Options: "", Args: "1"},
		{Type: token.RIGHT, Options: "", Args: "2"},
		{Type: token.RIGHT, Options: "1s", Args: "3"},
		{Type: token.UP, Options: "", Args: "1"},
		{Type: token.UP, Options: "", Args: "2"},
		{Type: token.UP, Options: "1s", Args: "3"},
		{Type: token.PAGE_UP, Options: "", Args: "1"},
		{Type: token.PAGE_UP, Options: "", Args: "2"},
		{Type: token.PAGE_UP, Options: "1s", Args: "3"},
		{Type: token.DOWN, Options: "", Args: "1"},
		{Type: token.DOWN, Options: "", Args: "2"},
		{Type: token.DOWN, Options: "1s", Args: "3"},
		{Type: token.CTRL, Options: "", Args: "C"},
		{Type: token.CTRL, Options: "", Args: "L"},
		{Type: token.CTRL, Options: "", Args: "R"},
		{Type: token.ALT, Options: "", Args: "."},
		{Type: token.ALT, Options: "", Args: "L"},
		{Type: token.ALT, Options: "", Args: "i"},
		{Type: token.HIDE, Options: "", Args: ""},
		{Type: token.SHOW, Options: "", Args: ""},
	}

	l := lexer.New(string(input))
	p := New(l)

	cmds := p.Parse()

	if len(cmds) != len(expected) {
		t.Fatalf("Expected %d commands, got %d", len(expected), len(cmds))
	}

	for i, cmd := range cmds {
		if cmd.Type != expected[i].Type {
			t.Errorf("Expected command %d to be %s, got %s", i, expected[i].Type, cmd.Type)
		}
		if cmd.Args != expected[i].Args {
			t.Errorf("Expected command %d to have args %s, got %s", i, expected[i].Args, cmd.Args)
		}
		if cmd.Options != expected[i].Options {
			t.Errorf("Expected command %d to have options %s, got %s", i, expected[i].Options, cmd.Options)
		}
	}
}

func TestParseCtrl(t *testing.T) {
	tests := []struct {
		name     string
		tape     string
		wantArgs []string
		wantErr  bool
	}{
		{
			name:     "should parse with multiple modifiers",
			tape:     "Ctrl+Shift+Alt+C",
			wantArgs: []string{"Shift", "Alt", "C"},
			wantErr:  false,
		},
		{
			name:    "should not parse with out of order modifiers",
			tape:    "Ctrl+Shift+C+Alt",
			wantErr: true,
		},
		{
			name:    "should not parse with out of order modifiers",
			tape:    "Ctrl+Shift+C+Alt+C",
			wantErr: true,
		},
		{
			tape:    "Ctrl+Alt+Right",
			wantErr: true,
		},
		{
			name:     "Ctrl+Backspace",
			tape:     "Ctrl+Backspace",
			wantArgs: []string{"Backspace"},
			wantErr:  false,
		},
		{
			name:     "Ctrl+Space",
			tape:     "Ctrl+Space",
			wantArgs: []string{"Space"},
			wantErr:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l := lexer.New(tc.tape)
			p := New(l)

			cmd := p.parseCtrl()
			if tc.wantErr {
				if len(p.errors) == 0 {
					t.Errorf("Expected to parse with errors but was success")
				}
				return
			}

			if len(p.errors) > 0 {
				t.Errorf("Expected to parse with no errors but was failure")
			}

			args := strings.Split(cmd.Args, " ")
			if len(tc.wantArgs) != len(args) {
				t.Fatalf("Unable to parse args, expected args %d, got %d", len(tc.wantArgs), len(args))
			}

			for i, arg := range args {
				if tc.wantArgs[i] != arg {
					t.Errorf("Arg %d is wrong, expected %s, got %s", i, tc.wantArgs[i], arg)
				}
			}
		})
	}
}

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

	l := lexer.New(st.tape)
	p := New(l)

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

type parseScreenshotTest struct {
	tape   string
	errors []string
}

func (st *parseScreenshotTest) run(t *testing.T) {
	l := lexer.New(st.tape)
	p := New(l)

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
