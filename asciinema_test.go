package main

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"
)

func TestAsciinemaRecorder_CaptureFrame(t *testing.T) {
	r := NewAsciinemaRecorder("", 80, 24)

	// Capture first frame
	r.CaptureFrame([]string{"$ echo hello", "hello", "$"})
	if r.EventCount() != 1 {
		t.Errorf("expected 1 event, got %d", r.EventCount())
	}

	// Same content should not create new event
	r.CaptureFrame([]string{"$ echo hello", "hello", "$"})
	if r.EventCount() != 1 {
		t.Errorf("expected 1 event (no change), got %d", r.EventCount())
	}

	// Different content should create new event
	r.CaptureFrame([]string{"$ echo world", "world", "$"})
	if r.EventCount() != 2 {
		t.Errorf("expected 2 events, got %d", r.EventCount())
	}
}

func TestAsciinemaRecorder_Save(t *testing.T) {
	tmpFile := t.TempDir() + "/test.cast"

	r := NewAsciinemaRecorder(tmpFile, 80, 24)
	r.SetTitle("Test Recording")

	// Capture some frames
	r.CaptureFrame([]string{"$ echo hello", "hello"})
	time.Sleep(100 * time.Millisecond)
	r.CaptureFrame([]string{"$ echo world", "world"})

	// Save
	err := r.Save()
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Read and verify
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) < 3 {
		t.Fatalf("expected at least 3 lines (header + 2 events), got %d", len(lines))
	}

	// Verify header
	var header AsciinemaHeader
	if err := json.Unmarshal([]byte(lines[0]), &header); err != nil {
		t.Fatalf("Failed to parse header: %v", err)
	}

	if header.Version != 2 {
		t.Errorf("expected version 2, got %d", header.Version)
	}
	if header.Width != 80 {
		t.Errorf("expected width 80, got %d", header.Width)
	}
	if header.Height != 24 {
		t.Errorf("expected height 24, got %d", header.Height)
	}
	if header.Title != "Test Recording" {
		t.Errorf("expected title 'Test Recording', got %q", header.Title)
	}

	// Verify events are parseable
	for i, line := range lines[1:] {
		var event []interface{}
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			t.Errorf("Failed to parse event %d: %v", i, err)
		}
		if len(event) != 3 {
			t.Errorf("Event %d: expected 3 elements, got %d", i, len(event))
		}
	}
}

func TestAsciinemaEvent_MarshalJSON(t *testing.T) {
	event := AsciinemaEvent{
		Time: 1.5,
		Type: "o",
		Data: "hello\n",
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	expected := `[1.5,"o","hello\n"]`
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}
