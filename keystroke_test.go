package main

import (
	"testing"
)

// checkKeyStrokeEvents checks if the key stroke events are as expected.
func checkKeyStrokeEvents(t *testing.T, events *KeyStrokeEvents, expected ...string) {
	for i, event := range events.Slice() {
		if actual := event.Display; expected[i] != actual {
			t.Fatalf("expected event display %q, got %q", expected, actual)
		}
	}
}

func defaultKeyStrokeEvents() *KeyStrokeEvents {
	events := NewKeyStrokeEvents(DefaultMaxDisplaySize)
	events.enabled = true

	return events
}
