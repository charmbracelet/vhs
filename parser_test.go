package main

import (
	"os"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	input := `
Set TypingSpeed 100ms
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
Sleep 3`

	expected := []Command{
		{Type: SET, Options: "TypingSpeed", Args: "100ms"},
		{Type: TYPE, Options: "", Args: "echo 'Hello, World!'"},
		{Type: ENTER, Options: "", Args: "1"},
		{Type: BACKSPACE, Options: "0.1s", Args: "5"},
		{Type: BACKSPACE, Options: ".1s", Args: "5"},
		{Type: BACKSPACE, Options: "1s", Args: "5"},
		{Type: BACKSPACE, Options: "100ms", Args: "5"},
		{Type: DELETE, Options: "", Args: "2"},
		{Type: INSERT, Options: "", Args: "2"},
		{Type: RIGHT, Options: "", Args: "3"},
		{Type: LEFT, Options: "", Args: "3"},
		{Type: UP, Options: "50ms", Args: "1"},
		{Type: DOWN, Options: "", Args: "2"},
		{Type: CTRL, Options: "", Args: "C"},
		{Type: CTRL, Options: "", Args: "L"},
		{Type: ALT, Options: "", Args: "."},
		{Type: SLEEP, Args: "100ms"},
		{Type: SLEEP, Args: "3s"},
	}

	l := NewLexer(input)
	p := NewParser(l)

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

func TestParserErrors(t *testing.T) {
	input := `
Type Enter
Type "echo 'Hello, World!'" Enter
Foo
Sleep Bar`

	l := NewLexer(input)
	p := NewParser(l)

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
	input, err := os.ReadFile("examples/fixtures/all.tape")
	if err != nil {
		t.Fatal("could not read fixture file")
	}

	expected := []Command{
		{Type: OUTPUT, Options: ".gif", Args: "examples/fixtures/all.gif"},
		{Type: OUTPUT, Options: ".mp4", Args: "examples/fixtures/all.mp4"},
		{Type: OUTPUT, Options: ".webm", Args: "examples/fixtures/all.webm"},
		{Type: SET, Options: "Shell", Args: "fish"},
		{Type: SET, Options: "FontSize", Args: "22"},
		{Type: SET, Options: "FontFamily", Args: "DejaVu Sans Mono"},
		{Type: SET, Options: "Height", Args: "600"},
		{Type: SET, Options: "Width", Args: "1200"},
		{Type: SET, Options: "LetterSpacing", Args: "1"},
		{Type: SET, Options: "LineHeight", Args: "1.2"},
		{Type: SET, Options: "Theme", Args: "{ \"name\": \"Whimsy\", \"black\": \"#535178\", \"red\": \"#ef6487\", \"green\": \"#5eca89\", \"yellow\": \"#fdd877\", \"blue\": \"#65aef7\", \"purple\": \"#aa7ff0\", \"cyan\": \"#43c1be\", \"white\": \"#ffffff\", \"brightBlack\": \"#535178\", \"brightRed\": \"#ef6487\", \"brightGreen\": \"#5eca89\", \"brightYellow\": \"#fdd877\", \"brightBlue\": \"#65aef7\", \"brightPurple\": \"#aa7ff0\", \"brightCyan\": \"#43c1be\", \"brightWhite\": \"#ffffff\", \"background\": \"#29283b\", \"foreground\": \"#b3b0d6\", \"selectionBackground\": \"#3d3c58\", \"cursorColor\": \"#b3b0d6\" }"},
		{Type: SET, Options: "Theme", Args: "Catppuccin Mocha"},
		{Type: SET, Options: "Padding", Args: "50"},
		{Type: SET, Options: "Framerate", Args: "60"},
		{Type: SET, Options: "PlaybackSpeed", Args: "2"},
		{Type: SET, Options: "TypingSpeed", Args: ".1s"},
		{Type: SET, Options: "LoopOffset", Args: "60.4%"},
		{Type: SET, Options: "LoopOffset", Args: "20.99%"},
		{Type: SET, Options: "CursorBlink", Args: "false"},
		{Type: SLEEP, Options: "", Args: "1s"},
		{Type: SLEEP, Options: "", Args: "500ms"},
		{Type: SLEEP, Options: "", Args: ".5s"},
		{Type: SLEEP, Options: "", Args: "0.5s"},
		{Type: TYPE, Options: ".5s", Args: "All"},
		{Type: TYPE, Options: "500ms", Args: "All"},
		{Type: TYPE, Options: "", Args: "Double Quote"},
		{Type: TYPE, Options: "", Args: "\"Single\" Quote"},
		{Type: TYPE, Options: "", Args: `"Backtick" 'Quote'`},
		{Type: BACKSPACE, Options: "", Args: "1"},
		{Type: BACKSPACE, Options: "", Args: "2"},
		{Type: BACKSPACE, Options: "1s", Args: "3"},
		{Type: DELETE, Options: "", Args: "1"},
		{Type: DELETE, Options: "", Args: "2"},
		{Type: DELETE, Options: "1s", Args: "3"},
		{Type: INSERT, Options: "", Args: "1"},
		{Type: INSERT, Options: "", Args: "2"},
		{Type: INSERT, Options: "1s", Args: "3"},
		{Type: DOWN, Options: "", Args: "1"},
		{Type: DOWN, Options: "", Args: "2"},
		{Type: DOWN, Options: "1s", Args: "3"},
		{Type: PAGEDOWN, Options: "", Args: "1"},
		{Type: PAGEDOWN, Options: "", Args: "2"},
		{Type: PAGEDOWN, Options: "1s", Args: "3"},
		{Type: ENTER, Options: "", Args: "1"},
		{Type: ENTER, Options: "", Args: "2"},
		{Type: ENTER, Options: "1s", Args: "3"},
		{Type: SPACE, Options: "", Args: "1"},
		{Type: SPACE, Options: "", Args: "2"},
		{Type: SPACE, Options: "1s", Args: "3"},
		{Type: TAB, Options: "", Args: "1"},
		{Type: TAB, Options: "", Args: "2"},
		{Type: TAB, Options: "1s", Args: "3"},
		{Type: LEFT, Options: "", Args: "1"},
		{Type: LEFT, Options: "", Args: "2"},
		{Type: LEFT, Options: "1s", Args: "3"},
		{Type: RIGHT, Options: "", Args: "1"},
		{Type: RIGHT, Options: "", Args: "2"},
		{Type: RIGHT, Options: "1s", Args: "3"},
		{Type: UP, Options: "", Args: "1"},
		{Type: UP, Options: "", Args: "2"},
		{Type: UP, Options: "1s", Args: "3"},
		{Type: PAGEUP, Options: "", Args: "1"},
		{Type: PAGEUP, Options: "", Args: "2"},
		{Type: PAGEUP, Options: "1s", Args: "3"},
		{Type: DOWN, Options: "", Args: "1"},
		{Type: DOWN, Options: "", Args: "2"},
		{Type: DOWN, Options: "1s", Args: "3"},
		{Type: CTRL, Options: "", Args: "C"},
		{Type: CTRL, Options: "", Args: "L"},
		{Type: CTRL, Options: "", Args: "R"},
		{Type: ALT, Options: "", Args: "."},
		{Type: ALT, Options: "", Args: "L"},
		{Type: ALT, Options: "", Args: "i"},
		{Type: HIDE, Options: "", Args: ""},
		{Type: SHOW, Options: "", Args: ""},
	}

	l := NewLexer(string(input))
	p := NewParser(l)

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
	t.Run("should parse with multiple modifiers", func(t *testing.T) {
		tape := "Ctrl+Shift+Alt+C"
		l := NewLexer(tape)
		p := NewParser(l)

		cmd := p.parseCtrl()

		expectedArgs := []string{"Shift", "Alt", "C"}
		args := strings.Split(cmd.Args, " ")

		if len(expectedArgs) != len(args) {
			t.Fatalf("Unable to parse args, expected args %d, got %d", len(expectedArgs), len(args))
		}

		for i, arg := range args {
			if expectedArgs[i] != arg {
				t.Errorf("Arg %d is wrong, expected %s, got %s", i, expectedArgs[i], arg)
			}
		}
	})

	t.Run("should parse with errors when using unknown modifier", func(t *testing.T) {
		tape := "Ctrl+AltRight"
		l := NewLexer(tape)
		p := NewParser(l)

		_ = p.parseCtrl()

		if len(p.errors) == 0 {
			t.Errorf("Expected to parse with errors but was success")
		}
	})

	t.Run("should parse with errors when using keyword as modifier", func(t *testing.T) {
		tape := "Ctrl+Backspace"
		l := NewLexer(tape)
		p := NewParser(l)

		_ = p.parseCtrl()

		if len(p.errors) == 0 {
			t.Errorf("Expected to parse with errors but was success")
		}
	})
}
