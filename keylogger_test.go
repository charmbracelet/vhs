package main

import (
	"testing"
)

func TestKeyLogger_Basic(t *testing.T) {
	l := NewKeyLogger()

	l.LogKey("a")
	if len(l.events) != 0 {
		t.Error("Expected no events when recorder has not been started")
	}

	var frame int64
	l.Start(&frame, 50) // 50 fps
	l.Enable()

	l.LogKey("a")
	frame = 5 // 100ms at 50fps
	l.LogKey("b")

	// Verify correct number of events
	if len(l.events) != 2 {
		t.Errorf("Expected 2 events, got %d", len(l.events))
	}

	// Verify ordering of events
	if l.events[0].StartMs >= l.events[1].StartMs {
		t.Error("Expected first event to have earlier timestamp")
	}

	// Verify correct timestamps
	if l.events[0].StartMs != 0 {
		t.Errorf("Expected first event at 0ms, got %d", l.events[0].StartMs)
	}
	if l.events[1].StartMs != 100 {
		t.Errorf("Expected second event at 100ms, got %d", l.events[1].StartMs)
	}

	// Verify correct keys
	if l.events[0].Key != "a" {
		t.Errorf("Expected first event as 'a', got '%s'", l.events[0].Key)
	}
	if l.events[1].Key != "b" {
		t.Errorf("Expected second event as 'b', got '%s'", l.events[1].Key)
	}

}

func TestKeyLogger_PauseResume(t *testing.T) {
	l := NewKeyLogger()

	var frame int64
	l.Start(&frame, 50) // 50 fps
	l.Enable()

	l.LogKey("a")
	l.Pause()
	l.LogKey("b") // this should not be logged
	l.Resume()
	l.LogKey("c")

	// Verify correct number of events
	if len(l.events) != 2 {
		t.Errorf("Expected 2 events, got %d", len(l.events))
	}

	// Verify correct keys
	if l.events[0].Key != "a" {
		t.Errorf("Expected first event as 'a', got '%s'", l.events[0].Key)
	}
	if l.events[1].Key != "c" {
		t.Errorf("Expected second event as 'c', got '%s'", l.events[1].Key)
	}
}

func TestKeyLogger_Events(t *testing.T) {
	l := NewKeyLogger()

	var frame int64
	l.Start(&frame, 50) // 50 fps
	l.Enable()

	l.LogKey("a")
	l.LogKey("b")
	l.LogKey("c")

	events := l.Events()
	if len(events) != 3 {
		t.Errorf("Expected 3 events, got %d", len(events))
	}

	expectedKeys := []string{"a", "b", "c"}
	for i, expected := range expectedKeys {
		if events[i].Key != expected {
			t.Errorf("Event %d: expected key '%s', got '%s'", i, expected, events[i].Key)
		}
	}
}

func TestKeyLogger_EnableDisable(t *testing.T) {
	l := NewKeyLogger()

	var frame int64
	l.Start(&frame, 50)

	// Not enabled by default — keys should be ignored
	l.LogKey("a")
	if len(l.events) != 0 {
		t.Errorf("Expected 0 events when disabled, got %d", len(l.events))
	}

	l.Enable()
	l.LogKey("b")
	if len(l.events) != 1 {
		t.Errorf("Expected 1 event after Enable, got %d", len(l.events))
	}

	l.Disable()
	l.LogKey("c")
	if len(l.events) != 1 {
		t.Errorf("Expected 1 event after Disable, got %d", len(l.events))
	}

	l.Enable()
	l.LogKey("d")
	if len(l.events) != 2 {
		t.Errorf("Expected 2 events after re-Enable, got %d", len(l.events))
	}

	if l.events[0].Key != "b" || l.events[1].Key != "d" {
		t.Errorf("Expected keys 'b' and 'd', got '%s' and '%s'", l.events[0].Key, l.events[1].Key)
	}
}

func TestKeyLogger_PauseAndEnabledInteraction(t *testing.T) {
	l := NewKeyLogger()

	var frame int64
	l.Start(&frame, 50)
	l.Enable()

	l.LogKey("a") // logged
	l.Pause()
	l.LogKey("b") // paused — not logged
	l.Resume()
	l.LogKey("c") // logged (enabled + not paused)

	// Disable while not paused
	l.Disable()
	l.LogKey("d") // disabled — not logged

	// Pause while disabled, then enable
	l.Pause()
	l.Enable()
	l.LogKey("e") // paused — not logged even though enabled

	l.Resume()
	l.LogKey("f") // enabled + not paused — logged

	if len(l.events) != 3 {
		t.Errorf("Expected 3 events, got %d", len(l.events))
	}
	expectedKeys := []string{"a", "c", "f"}
	for i, expected := range expectedKeys {
		if l.events[i].Key != expected {
			t.Errorf("Event %d: expected key '%s', got '%s'", i, expected, l.events[i].Key)
		}
	}
}

func TestKeyLogger_FrameAlignment(t *testing.T) {
	tests := []struct {
		name      string
		framerate int
		frames    []int64
		wantMs    []int64
	}{
		{
			name:      "50fps",
			framerate: 50,
			frames:    []int64{0, 25, 50, 100},
			wantMs:    []int64{0, 500, 1000, 2000},
		},
		{
			name:      "30fps",
			framerate: 30,
			frames:    []int64{0, 15, 30, 60},
			wantMs:    []int64{0, 500, 1000, 2000},
		},
		{
			name:      "60fps",
			framerate: 60,
			frames:    []int64{0, 30, 60, 120},
			wantMs:    []int64{0, 500, 1000, 2000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewKeyLogger()
			var frame int64
			l.Start(&frame, tt.framerate)
			l.Enable()

			for _, f := range tt.frames {
				frame = f
				l.LogKey("a")
			}

			if len(l.events) != len(tt.wantMs) {
				t.Errorf("Expected %d events, got %d", len(tt.wantMs), len(l.events))
			}

			for i, want := range tt.wantMs {
				if l.events[i].StartMs != want {
					t.Errorf("Event %d: expected %dms, got %dms", i, want, l.events[i].StartMs)
				}
			}

		})
	}
}
