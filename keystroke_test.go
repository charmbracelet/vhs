package main

import (
	"testing"

	"github.com/go-rod/rod/lib/input"
)

// checkKeyStrokeEvents checks if the key stroke events are as expected.
func checkKeyStrokeEvents(t *testing.T, events *KeyStrokeEvents, expected ...string) {
	for i, event := range events.events {
		if actual := event.Display; expected[i] != actual {
			t.Fatalf("expected event display %q, got %q", expected, actual)
		}
	}
}

func defaultKeyStrokeEvents() *KeyStrokeEvents {
	events := NewKeyStrokeEvents(DefaultMaxDisplaySize)
	events.Enable()

	return events
}

func TestKeyStrokeEventsRemembersKeyStrokes(t *testing.T) {
	events := defaultKeyStrokeEvents()
	events.Push("a")
	events.Push("b")
	events.Push("c")
	checkKeyStrokeEvents(t, events, "a", "ab", "abc")
}

func TestKeyStrokeEventsHonorsMaxDisplaySize(t *testing.T) {
	events := defaultKeyStrokeEvents()
	events.maxDisplaySize = 2

	events.Push("a")
	events.Push("b")
	events.Push("c")

	// NOTE: It should not be "ab", but "bc" at the end -- we should be acting
	// like a ring buffer.
	checkKeyStrokeEvents(t, events, "a", "ab", "bc")
}

func TestKeyStrokeEventsShowsNothingIfDisabled(t *testing.T) {
	events := defaultKeyStrokeEvents()
	events.Disable()

	events.Push("a")
	events.Push("b")
	events.Push("c")

	checkKeyStrokeEvents(t, events)
}

func TestKeyStrokeEventsKeyToDisplay(t *testing.T) {
	cases := []struct {
		name     string
		key      input.Key
		expected string
	}{
		{
			name:     "letter",
			key:      input.KeyA,
			expected: "a",
		},
		{
			name:     "number",
			key:      input.Digit1,
			expected: "1",
		},
		{
			name:     "symbol no override",
			key:      input.Minus,
			expected: "-",
		},
		{
			name:     "shifted key",
			key:      shift(input.Minus),
			expected: "_",
		},
		{
			name:     "symbol override",
			key:      input.Backspace,
			expected: "⌫",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if actual := keyToDisplay(tc.key); tc.expected != actual {
				t.Fatalf("expected display %q, got %q", tc.expected, actual)
			}
		})
	}
}
