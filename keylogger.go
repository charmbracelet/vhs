package main

import (
	"sync/atomic"
)

// KeyEvent represents a single key press and timing information.
type KeyEvent struct {
	// StartMs is the number of milliseconds relative to the recording start.
	StartMs int64 `json:"ms"`

	// Key is the string representation of the key pressed.
	Key string `json:"key"`
}

// KeyLogger tracks key events during tape execution.
type KeyLogger struct {
	events    []KeyEvent
	paused    bool
	enabled   bool
	frame     *int64 // pointer to shared frame counter
	framerate int
}

// NewKeyLogger creates a new KeyLogger.
func NewKeyLogger() *KeyLogger {
	return &KeyLogger{
		events: make([]KeyEvent, 0),
	}
}

// Start begins key recording.
//
// The first parameter is a pointer to the shared frame counter that is being
// incremented asynchronously by VHS as captures are being made. This is used
// with the framerate parameter to compute when a key event was made.
func (l *KeyLogger) Start(frame *int64, framerate int) {
	l.frame = frame
	l.framerate = framerate
}

// Pause suspends key logging until Resume is called.
func (l *KeyLogger) Pause() {
	l.paused = true
}

// Resume enables key logging after Pause is called.
func (l *KeyLogger) Resume() {
	l.paused = false
}

// Enable turns on caption key logging.
func (l *KeyLogger) Enable() {
	l.enabled = true
}

// Disable turns off caption key logging.
func (l *KeyLogger) Disable() {
	l.enabled = false
}

// LogKey records the current key and time at which it occurred with respect
// to the current frame being captured. If key logging has not been started or
// if it has been paused, this does nothing.
func (l *KeyLogger) LogKey(key string) {
	if l.frame == nil || l.paused || !l.enabled {
		return
	}

	frameNum := atomic.LoadInt64(l.frame)
	timeMs := frameNum * 1000 / int64(l.framerate)

	event := KeyEvent{
		StartMs: timeMs,
		Key:     key,
	}
	l.events = append(l.events, event)
}

// Events returns the recorded key events.
func (l *KeyLogger) Events() []KeyEvent {
	return l.events
}
