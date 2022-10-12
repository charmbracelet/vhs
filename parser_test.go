package vhs

import "testing"

func TestParser(t *testing.T) {
	input := `
Set TypingSpeed 100ms
Type "echo 'Hello, World!'" Enter
Backspace@100 5
Backspace@100ms 5
Backspace@1s 5
Right 3 Left 3
Up@50 Down 2
Ctrl+C
Ctrl+L
Sleep 100
Sleep 3s`

	expected := []Command{
		{Type: SET, Options: "TypingSpeed", Args: "100ms"},
		{Type: TYPE, Options: "", Args: "echo 'Hello, World!'"},
		{Type: ENTER, Options: "", Args: "1"},
		{Type: BACKSPACE, Options: "100ms", Args: "5"},
		{Type: BACKSPACE, Options: "100ms", Args: "5"},
		{Type: BACKSPACE, Options: "1s", Args: "5"},
		{Type: RIGHT, Options: "", Args: "3"},
		{Type: LEFT, Options: "", Args: "3"},
		{Type: UP, Options: "50ms", Args: "1"},
		{Type: DOWN, Options: "", Args: "2"},
		{Type: CTRL, Options: "", Args: "C"},
		{Type: CTRL, Options: "", Args: "L"},
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
