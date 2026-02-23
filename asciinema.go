package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// AsciinemaHeader is the header for an asciinema v2 recording.
type AsciinemaHeader struct {
	Version   int               `json:"version"`
	Width     int               `json:"width"`
	Height    int               `json:"height"`
	Timestamp int64             `json:"timestamp"`
	Title     string            `json:"title,omitempty"`
	Env       map[string]string `json:"env,omitempty"`
}

// AsciinemaEvent represents a single event in an asciinema recording.
// Format: [time, "o", data] for output events.
type AsciinemaEvent struct {
	Time float64
	Type string // "o" for output, "i" for input
	Data string
}

// MarshalJSON implements custom JSON marshaling for AsciinemaEvent.
func (e AsciinemaEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{e.Time, e.Type, e.Data}) //nolint:wrapcheck
}

// AsciinemaRecorder captures terminal output for asciinema format.
type AsciinemaRecorder struct {
	mu        sync.Mutex
	events    []AsciinemaEvent
	startTime time.Time
	lastState string
	output    string
	width     int
	height    int
	title     string
}

// NewAsciinemaRecorder creates a new asciinema recorder.
func NewAsciinemaRecorder(output string, width, height int) *AsciinemaRecorder {
	return &AsciinemaRecorder{
		output:    output,
		width:     width,
		height:    height,
		startTime: time.Now(),
		events:    make([]AsciinemaEvent, 0),
	}
}

// SetTitle sets the recording title.
func (r *AsciinemaRecorder) SetTitle(title string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.title = title
}

// CaptureFrame captures the current terminal state from xterm.js buffer.
// Only records if the state has changed from the previous capture.
//
// Note: Since VHS renders via xterm.js in a browser, we don't have direct
// PTY access. This captures screen snapshots and emits ANSI sequences to
// recreate the display. The result is playable but may differ from native
// asciinema recordings.
func (r *AsciinemaRecorder) CaptureFrame(lines []string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Build current state for change detection
	currentState := strings.Join(lines, "\n")

	// Only record if changed
	if currentState == r.lastState {
		return
	}

	// Calculate time offset in seconds
	timeOffset := time.Since(r.startTime).Seconds()

	// Emit ANSI clear screen + content
	var output strings.Builder
	output.WriteString("\x1b[2J\x1b[H")

	for i, line := range lines {
		output.WriteString(line)
		if i < len(lines)-1 {
			output.WriteString("\n")
		}
	}

	r.events = append(r.events, AsciinemaEvent{
		Time: timeOffset,
		Type: "o",
		Data: output.String(),
	})

	r.lastState = currentState
}

// Save writes the asciinema recording to the output file.
func (r *AsciinemaRecorder) Save() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.output == "" {
		return nil
	}

	// Ensure directory exists
	if dir := filepath.Dir(r.output); dir != "." {
		if err := os.MkdirAll(dir, 0o750); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	file, err := os.Create(r.output)
	if err != nil {
		return fmt.Errorf("failed to create asciinema file: %w", err)
	}
	defer file.Close() //nolint:errcheck

	// Write header
	header := AsciinemaHeader{
		Version:   2,
		Width:     r.width,
		Height:    r.height,
		Timestamp: r.startTime.Unix(),
		Title:     r.title,
		Env: map[string]string{
			"SHELL": "/bin/bash",
			"TERM":  "xterm-256color",
		},
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return fmt.Errorf("failed to marshal header: %w", err)
	}
	_, _ = fmt.Fprintln(file, string(headerJSON))

	// Write events
	for _, event := range r.events {
		eventJSON, err := json.Marshal(event)
		if err != nil {
			continue
		}
		_, _ = fmt.Fprintln(file, string(eventJSON))
	}

	return nil
}

// EventCount returns the number of recorded events.
func (r *AsciinemaRecorder) EventCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.events)
}
