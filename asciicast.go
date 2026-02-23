package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// AsciicastHeader is the header of an asciicast v2 file.
type AsciicastHeader struct {
	Version   int               `json:"version"`
	Width     int               `json:"width"`
	Height    int               `json:"height"`
	Timestamp int64             `json:"timestamp,omitempty"`
	Env       map[string]string `json:"env,omitempty"`
}

// AsciicastEvent represents a single event in an asciicast v2 recording.
// Format: [time, event-type, event-data]
type AsciicastEvent struct {
	Time float64
	Type string
	Data string
}

// MarshalJSON implements custom JSON marshaling for AsciicastEvent
// to produce the [time, type, data] array format.
func (e AsciicastEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{e.Time, e.Type, e.Data})
}

const cast = ".cast"

// GenerateAsciicast creates an asciicast v2 file from terminal buffer snapshots
// captured during recording. It writes the asciicast v2 format (JSON Lines)
// with a header line followed by event lines.
func GenerateAsciicast(v *VHS) error {
	path := v.Options.Video.Output.Asciicast
	if path == "" {
		return nil
	}

	log.Println(GrayStyle.Render("Creating " + path + "..."))

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return fmt.Errorf("failed to create directory for asciicast: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create asciicast file: %w", err)
	}
	defer f.Close()

	// Calculate terminal dimensions (cols x rows) from the xterm viewport.
	// We use rough estimates based on font size and pixel dimensions.
	style := v.Options.Video.Style
	padding := style.Padding
	margin := 0
	if style.MarginFill != "" {
		margin = style.Margin
	}
	bar := 0
	if style.WindowBar != "" {
		bar = style.WindowBarSize
	}

	pixelWidth := style.Width - double(padding) - double(margin)
	pixelHeight := style.Height - double(padding) - double(margin) - bar

	// Approximate character cell size: width ~= fontSize * 0.6, height ~= fontSize * 1.2
	charWidth := int(float64(v.Options.FontSize) * 0.6)  //nolint:mnd
	charHeight := int(float64(v.Options.FontSize) * 1.2) //nolint:mnd
	if charWidth <= 0 {
		charWidth = 1
	}
	if charHeight <= 0 {
		charHeight = 1
	}

	cols := pixelWidth / charWidth
	rows := pixelHeight / charHeight
	if cols <= 0 {
		cols = 80 //nolint:mnd
	}
	if rows <= 0 {
		rows = 24 //nolint:mnd
	}

	// Write the header line.
	header := AsciicastHeader{
		Version:   2,
		Width:     cols,
		Height:    rows,
		Timestamp: time.Now().Unix(),
		Env: map[string]string{
			"TERM":  "xterm-256color",
			"SHELL": "/bin/bash",
		},
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return fmt.Errorf("failed to marshal asciicast header: %w", err)
	}

	if _, err := fmt.Fprintln(f, string(headerJSON)); err != nil {
		return fmt.Errorf("failed to write asciicast header: %w", err)
	}

	// Write the buffered frames as asciicast output events.
	// Each frame becomes an "o" (output) event at the appropriate timestamp.
	framerate := v.Options.Video.Framerate
	if framerate <= 0 {
		framerate = defaultFramerate
	}

	var prev string
	for i, frame := range v.AsciicastFrames {
		content := strings.Join(frame, "\r\n")
		// Only emit events when content changes to reduce file size.
		if content == prev {
			continue
		}
		prev = content

		elapsed := float64(i) / float64(framerate)
		event := AsciicastEvent{
			Time: elapsed,
			Type: "o",
			Data: content + "\r\n",
		}

		eventJSON, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("failed to marshal asciicast event: %w", err)
		}

		if _, err := fmt.Fprintln(f, string(eventJSON)); err != nil {
			return fmt.Errorf("failed to write asciicast event: %w", err)
		}
	}

	return nil
}
