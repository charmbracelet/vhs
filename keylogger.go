package main

import (
	"encoding/json"
	"log"
	"os"
	"sync/atomic"
)

// KeyEvent represents a single key press and timing information.
type KeyEvent struct {
	// Ms is the number of milliseconds relative to the recording start.
	Ms int64 `json:"ms"`

	//  Key is the string representation of the key pressed.
	Key string `json:"key"`
}

// KeyLogger tracks key events during tape execution.
type KeyLogger struct {
	events    []KeyEvent
	paused    bool
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

// LogKey records the current key and time at which it occurred with respect
// to the current frame being captured. If key logging has not been started or
// if it has been paused, this does nothing.
func (l *KeyLogger) LogKey(key string) {
	if l.frame == nil || l.paused {
		return
	}

	frameNum := atomic.LoadInt64(l.frame)
	timeMs := frameNum * 1000 / int64(l.framerate)

	event := KeyEvent{
		Ms:  timeMs,
		Key: key,
	}
	l.events = append(l.events, event)
}

// Save writes the recorded key events to logFile as JSON. If logFile is an
// empty string or if there are no events, this does nothing.
func (l *KeyLogger) Save(logFile string) error {
	if logFile == "" || len(l.events) == 0 {
		return nil
	}

	log.Println(GrayStyle.Render("Saving keylog to " + logFile + "..."))
	data, err := json.MarshalIndent(l.events, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(logFile, data, 0o644)
}
