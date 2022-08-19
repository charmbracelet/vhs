package dolly

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	input := `
Type@100 echo 'Hi, there!'
Left 3
Right 2

# Comment
Enter
Sleep 1s`

	expected := []Command{
		{Type: Type, Options: "100", Args: "echo 'Hi, there!'"},
		{Type: Left, Args: "3"},
		{Type: Right, Args: "2"},
		{Type: Enter},
		{Type: Sleep, Args: "1s"},
	}

	commands, errs := Parse(input)
	if len(errs) != 0 {
		for _, err := range errs {
			t.Log(err)
		}
		t.Fail()
	}
	if len(commands) != 5 {
		t.Errorf("expected 5 commands, got %d", len(commands))
	}

	for i, command := range commands {
		if command != expected[i] {
			t.Errorf("expected %v, got %v", expected[i], command)
		}
	}
}

func TestTypeOptions(t *testing.T) {
	input := `
Type@1 foo
# Trailing whitespace is significant
Type bar `

	expected := []Command{
		{Type: Type, Options: "1", Args: "foo"},
		{Type: Type, Args: "bar" + " "},
	}

	commands, err := Parse(input)
	if err != nil {
		t.Error(err)
	}

	if len(commands) != 2 {
		t.Errorf("expected 2 commands, got %d", len(commands))
	}

	for i, command := range commands {
		if command != expected[i] {
			t.Errorf("expected %v, got %v", expected[i], command)
		}
	}
}

func TestSpaceCommand(t *testing.T) {
	input := `
Space@1000 5
Space 5
Backspace@100 10
`

	expected := []Command{
		{Type: Space, Options: "1000", Args: "5"},
		{Type: Space, Options: "", Args: "5"},
		{Type: Backspace, Options: "100", Args: "10"},
	}

	commands, err := Parse(input)
	if err != nil {
		t.Error(err)
	}

	if len(commands) != 3 {
		t.Errorf("expected 3 commands, got %d", len(commands))
	}

	for i, command := range commands {
		if command != expected[i] {
			t.Errorf("expected %v, got %v", expected[i], command)
		}
	}
}

func TestSleepCommand(t *testing.T) {
	input := `
Sleep 1s
Sleep 100ms
`

	expected := []Command{
		{Type: Sleep, Args: "1s"},
		{Type: Sleep, Args: "100ms"},
	}

	commands, err := Parse(input)
	if err != nil {
		t.Error(err)
	}

	if len(commands) != 2 {
		t.Errorf("expected 2 commands, got %d", len(commands))
	}

	for i, command := range commands {
		if command != expected[i] {
			t.Errorf("expected %v, got %v", expected[i], command)
		}
	}
}

func TestSetCommand(t *testing.T) {
	input := `
Set FontFamily 32
Set Padding 5em
Set FontSize 15

Type Foo
Type Bar
`
	expected := []Command{
		{Type: Set, Options: "FontFamily", Args: "32"},
		{Type: Set, Options: "Padding", Args: "5em"},
		{Type: Set, Options: "FontSize", Args: "15"},
		{Type: Type, Args: "Foo"},
		{Type: Type, Args: "Bar"},
	}

	commands, err := Parse(input)
	if err != nil {
		t.Error(err)
	}

	if len(commands) != 5 {
		t.Errorf("expected 5 commands, got %d", len(commands))
	}

	for i, command := range commands {
		if command != expected[i] {
			t.Errorf("expected %v, got %v", expected[i], command)
		}
	}
}

func TestSetMultipleCommand(t *testing.T) {
	input := `
Set FontFamily 32
Set Padding 5em
Set FontSize 15

Type Foo
Type Bar

Set Padding 10em
Set FontSize 30
`
	expected := []Command{
		{Type: Set, Options: "FontFamily", Args: "32"},
		{Type: Set, Options: "Padding", Args: "5em"},
		{Type: Set, Options: "FontSize", Args: "15"},
		{Type: Type, Args: "Foo"},
		{Type: Type, Args: "Bar"},
		{Type: Set, Options: "Padding", Args: "10em"},
		{Type: Set, Options: "FontSize", Args: "30"},
	}

	commands, err := Parse(input)
	if err != nil {
		t.Error(err)
	}

	if len(commands) != 7 {
		t.Errorf("expected 7 commands, got %d", len(commands))
	}

	for i, command := range commands {
		if command != expected[i] {
			t.Errorf("expected %v, got %v", expected[i], command)
		}
	}
}

func TestInvalidString(t *testing.T) {
	tests := []struct {
		input string
		errs  []error
	}{
		{input: "Set FontFamily", errs: []error{fmt.Errorf("%s\n1 | Set FontFamily", ErrMissingArguments)}},
		{input: "Foo", errs: []error{fmt.Errorf("%s\n1 | Foo", ErrUnknownCommand)}},
		{input: "Set Foo Bar", errs: []error{fmt.Errorf("%s\n1 | Set Foo Bar", ErrUnknownOptions)}},
		{input: "Type", errs: []error{fmt.Errorf("%s\n1 | Type", ErrMissingArguments)}},
		{input: "Set FontFamily Monospace\nFoobar\nType", errs: []error{fmt.Errorf("%s\n2 | Foobar", ErrUnknownCommand), fmt.Errorf("%s\n3 | Type", ErrMissingArguments)}},
	}

	for _, test := range tests {
		_, errs := Parse(test.input)
		if len(errs) != len(test.errs) {
			t.Errorf("expected %d errors, got %d", len(test.errs), len(errs))
		}

		for i, err := range errs {
			if err.Error() != test.errs[i].Error() {
				t.Logf("Expected:\n%s\n", test.errs[i])
				t.Logf("Got:\n%s\n", err)
				t.Fail()
			}
		}
	}
}
