package main

import (
	"encoding/json"
	"os"
	"path/filepath"
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

	l.LogKey("a")
	frame = 5 // 100ms at 50fps
	l.LogKey("b")

	// Verify correct number of events
	if len(l.events) != 2 {
		t.Errorf("Expected 2 events, got %d", len(l.events))
	}

	// Verify ordering of events
	if l.events[0].Ms >= l.events[1].Ms {
		t.Error("Expected first event to have earlier timestamp")
	}

	// Verify correct timestamps
	if l.events[0].Ms != 0 {
		t.Errorf("Expected first event at 0ms, got %d", l.events[0].Ms)
	}
	if l.events[1].Ms != 100 {
		t.Errorf("Expected second event at 100ms, got %d", l.events[1].Ms)
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

func TestKeyLogger_Save(t *testing.T) {
	l := NewKeyLogger()

	var frame int64
	l.Start(&frame, 50) // 50 fps

	l.LogKey("a")
	l.LogKey("b")
	l.LogKey("c")

	// Write events to JSON file
	tempFile := filepath.Join(t.TempDir(), "keylog.json")
	if err := l.Save(tempFile); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Validate JSON file contains what we expected
	data, err := os.ReadFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	var events []KeyEvent
	if err := json.Unmarshal(data, &events); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if len(events) != 3 {
		t.Errorf("Expected 3 events in file, got %d", len(events))
	}

	expectedKeys := []string{"a", "b", "c"}
	for i, expected := range expectedKeys {
		if events[i].Key != expected {
			t.Errorf("Event %d: expected key '%s', got '%s'", i, expected, events[i].Key)
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

			for _, f := range tt.frames {
				frame = f
				l.LogKey("a")
			}

			if len(l.events) != len(tt.wantMs) {
				t.Errorf("Expected %d events, got %d", len(tt.wantMs), len(l.events))
			}

			for i, want := range tt.wantMs {
				if l.events[i].Ms != want {
					t.Errorf("Event %d: expected %dms, got %dms", i, want, l.events[i].Ms)
				}
			}

		})
	}
}
