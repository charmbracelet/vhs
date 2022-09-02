package vhs

import "testing"

func TestParser(t *testing.T) {
	input := `
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
		{Type: Type, Options: "100ms", Args: "echo 'Hello, World!'"},
		{Type: Enter, Options: "100ms", Args: "1"},
		{Type: Backspace, Options: "100ms", Args: "5"},
		{Type: Backspace, Options: "100ms", Args: "5"},
		{Type: Backspace, Options: "1s", Args: "5"},
		{Type: Right, Options: "100ms", Args: "3"},
		{Type: Left, Options: "100ms", Args: "3"},
		{Type: Up, Options: "50ms", Args: "1"},
		{Type: Down, Options: "100ms", Args: "2"},
		{Type: Ctrl, Options: "", Args: "C"},
		{Type: Ctrl, Options: "", Args: "L"},
		{Type: Sleep, Args: "100ms"},
		{Type: Sleep, Args: "3s"},
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
Type
Type "echo 'Hello, World!'" Enter
Foo
Sleep Bar`

	l := NewLexer(input)
	p := NewParser(l)

	_ = p.Parse()

	expectedErrors := []string{
		" 2:1  │ Type expects string",
		" 4:1  │ Invalid command: Foo",
		" 5:1  │ Expected time, got Bar",
		" 5:1  │ Invalid command: Bar",
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
