package main

import (
	"strings"
	"testing"
)

func TestInputToTape(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name: "ctrl key combinations",
			input: `echo "Hello,.
BACKSPACE
LEFT
LEFT
RIGHT
RIGHT
 world"
ENTER
ENTER
ENTER
ls
ENTER
ENTER
BACKSPACE
CTRL+C
CTRL+C
CTRL+C
CTRL+W
CTRL+A
CTRL+E
SLEEP
SLEEP
ALT+.
SLEEP
exit
`,
			want: `Type 'echo "Hello,.'
Backspace
Left 2
Right 2
Type ' world"'
Enter 3
Type "ls"
Enter 2
Backspace
Ctrl+C
Ctrl+C
Ctrl+C
Ctrl+W
Ctrl+A
Ctrl+E
Sleep 1s
Alt+.
Sleep 500ms
`,
		},
		{
			name: "PageUp, PageDown #559",
			input: `echo "Hello,.
PAGE_UP
PAGE_UP
PAGE_UP
PAGE_UP
PAGE_UP
PAGE_UP
PAGE_UP
PAGE_UP
PAGE_DOWN
PAGE_DOWN
PAGE_DOWN
PAGE_DOWN
exit
`,
			want: `Type 'echo "Hello,.'
PageUp 8
PageDown 4
`,
		},
	}
	for _, tc := range tests {
		got := inputToTape(tc.input)
		if tc.want != got {
			t.Fatalf("want:\n%s\ngot:\n%s\n", tc.want, got)
		}
	}
}

func TestInputToTapeLongSleep(t *testing.T) {
	input := strings.Repeat("SLEEP\n", 121) + "exit"
	want := "Sleep 60.5s\n"
	got := inputToTape(input)
	if want != got {
		t.Fatalf("want:\n%s\ngot:\n%s\n", want, got)
	}
}

func TestInputToTape_RepeatedSleepAfterExit(t *testing.T) {
	input := "SLEEP\nexit\nSLEEP\nSLEEP"
	want := "Sleep 500ms\nType \"exit\"\nSleep 1s\n"
	got := inputToTape(input)
	if want != got {
		t.Fatalf("want:\n%s\ngot:\n%s\n", want, got)
	}
}
