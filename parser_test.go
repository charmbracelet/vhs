package dolly

import "testing"

func TestAllCommands(t *testing.T) {
	if len(allCommands) != 9 {
		t.Errorf("unexpected number of commands: %d", len(allCommands))
	}
}

func TestParse(t *testing.T) {
	input := `
Type@100 echo 'Hi, there!'
Left 3
Right 2

# Comment
Enter
Sleep 1s`

	expected := []Command{
		{Type: Type, Options: "@100", Args: "echo 'Hi, there!'"},
		{Type: Left, Args: "3"},
		{Type: Right, Args: "2"},
		{Type: Enter},
		{Type: Sleep, Args: "1s"},
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

func TestTypeOptions(t *testing.T) {
	input := `
Type@1 foo
# Trailing whitespace is significant
Type bar `

	expected := []Command{
		{Type: Type, Options: "@1", Args: "foo"},
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
		{Type: Space, Options: "@1000", Args: "5"},
		{Type: Space, Options: "", Args: "5"},
		{Type: Backspace, Options: "@100", Args: "10"},
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
