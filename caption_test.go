package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func keyEvents(ms ...int64) []KeyEvent {
	out := make([]KeyEvent, len(ms))
	for i, t := range ms {
		out[i] = KeyEvent{StartMs: t, Key: "a"}
	}
	return out
}

func TestVimNormalize(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"a", "a"},
		{"Enter", "<CR>"},
		{"Backspace", "<BS>"},
		{"Escape", "<Esc>"},
		{"Ctrl+c", "<C-c>"},
		{"Alt+f", "<A-f>"},
		{"Shift+Tab", "<S-Tab>"},
		{"Space", "<Space>"},
		{"Delete", "<Del>"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := KeyStyleVim.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("KeyStyleVim.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestIconNormalize(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"a", "a"},
		{"Enter", "↵"},
		{"Escape", "Esc"},
		{"PageUp", "PgUp"},
		{"PageDown", "PgDn"},
		{"Backspace", "⌫"},
		{"Space", "␣"},
		{"Ctrl+c", "^c"},
		{"Alt+d", "⌥d"},
		{"Shift+Tab", "⇧⇥"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := KeyStyleIcon.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("KeyStyleIcon.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestCaptionWindows(t *testing.T) {
	tests := []struct {
		name              string
		events            []KeyEvent
		inactivityTimerMs int
		maxKeys           int
		want              []CaptionWindow
	}{
		{
			name:              "empty events",
			events:            []KeyEvent{},
			inactivityTimerMs: 1000,
			maxKeys:           3,
			want:              []CaptionWindow{},
		},
		{
			name:              "single event",
			events:            keyEvents(0),
			inactivityTimerMs: 1000,
			maxKeys:           3,
			want: []CaptionWindow{
				{First: 0, Last: 0, ShowUntil: 1000, IsTruncated: false},
			},
		},
		{
			name:              "two events in same session",
			events:            keyEvents(0, 500),
			inactivityTimerMs: 1000,
			maxKeys:           3,
			want: []CaptionWindow{
				{First: 0, Last: 0, ShowUntil: 500, IsTruncated: false},
				{First: 0, Last: 1, ShowUntil: 1500, IsTruncated: false},
			},
		},
		{
			name:              "gap exactly equals inactivity timer starts new session",
			events:            keyEvents(0, 1000),
			inactivityTimerMs: 1000,
			maxKeys:           3,
			want: []CaptionWindow{
				{First: 0, Last: 0, ShowUntil: 1000, IsTruncated: false},
				{First: 1, Last: 1, ShowUntil: 2000, IsTruncated: false},
			},
		},
		{
			name:              "gap one over inactivity timer starts new session",
			events:            keyEvents(0, 1001),
			inactivityTimerMs: 1000,
			maxKeys:           3,
			want: []CaptionWindow{
				{First: 0, Last: 0, ShowUntil: 1000, IsTruncated: false},
				{First: 1, Last: 1, ShowUntil: 2001, IsTruncated: false},
			},
		},
		{
			name:              "sliding window truncates when maxKeys exceeded",
			events:            keyEvents(0, 100, 200, 300),
			inactivityTimerMs: 1000,
			maxKeys:           3,
			want: []CaptionWindow{
				{First: 0, Last: 0, ShowUntil: 100, IsTruncated: false},
				{First: 0, Last: 1, ShowUntil: 200, IsTruncated: false},
				{First: 0, Last: 2, ShowUntil: 300, IsTruncated: false},
				{First: 1, Last: 3, ShowUntil: 1300, IsTruncated: true},
			},
		},
		{
			name:              "session reset clears truncation",
			events:            keyEvents(0, 100, 200, 2000, 2100),
			inactivityTimerMs: 1000,
			maxKeys:           2,
			want: []CaptionWindow{
				{First: 0, Last: 0, ShowUntil: 100, IsTruncated: false},
				{First: 0, Last: 1, ShowUntil: 200, IsTruncated: false},
				{First: 1, Last: 2, ShowUntil: 1200, IsTruncated: true},
				{First: 3, Last: 3, ShowUntil: 2100, IsTruncated: false},
				{First: 3, Last: 4, ShowUntil: 3100, IsTruncated: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := captionWindows(tt.events, tt.inactivityTimerMs, tt.maxKeys)
			if len(got) == 0 && len(tt.want) == 0 {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got  %v\nwant %v", got, tt.want)
			}
		})
	}
}

func TestMsToASS(t *testing.T) {
	tests := []struct {
		ms   int
		want string
	}{
		{0, "0:00:00.00"},
		{1000, "0:00:01.00"},
		{61000, "0:01:01.00"},
		{3661050, "1:01:01.05"},
		{500, "0:00:00.50"},
		{10, "0:00:00.01"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := msToASS(tt.ms)
			if got != tt.want {
				t.Errorf("msToASS(%d) = %q, want %q", tt.ms, got, tt.want)
			}
		})
	}
}

func TestAssEscape(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`hello`, `hello`},
		{`a\b`, `a\\b`},
		{`{tag}`, `\{tag\}`},
		{`a\{b}`, `a\\\{b\}`},
		{`a\Nb`, `a\Nb`},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := assEscape(tt.input)
			if got != tt.want {
				t.Errorf("assEscape(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestOpacityToASSAlpha(t *testing.T) {
	tests := []struct {
		opacity float64
		want    string
	}{
		{1.0, "00"},
		{0.0, "FF"},
		{0.5, "80"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := opacityToASSAlpha(tt.opacity)
			if got != tt.want {
				t.Errorf("opacityToASSAlpha(%f) = %q, want %q", tt.opacity, got, tt.want)
			}
		})
	}
}

func TestHexToASSBGR(t *testing.T) {
	tests := []struct {
		hex  string
		want string
	}{
		{"#FFCC66", "&H66CCFF&"},
		{"#FF0000", "&H0000FF&"},
		{"#00FF00", "&H00FF00&"},
	}
	for _, tt := range tests {
		t.Run(tt.hex, func(t *testing.T) {
			got := hexToASSBGR(tt.hex)
			if got != tt.want {
				t.Errorf("hexToASSBGR(%q) = %q, want %q", tt.hex, got, tt.want)
			}
		})
	}
}

func TestCaptionAlignment(t *testing.T) {
	tests := []struct {
		alignment CaptionAlignment
		want      int
	}{
		{AlignBottomLeft, 1},
		{AlignBottomCenter, 2},
		{AlignBottomRight, 3},
		{AlignMiddleLeft, 4},
		{AlignMiddleCenter, 5},
		{AlignMiddleRight, 6},
		{AlignTopLeft, 7},
		{AlignTopCenter, 8},
		{AlignTopRight, 9},
		{"unknown", 3}, // default
	}
	for _, tt := range tests {
		t.Run(string(tt.alignment), func(t *testing.T) {
			got := tt.alignment.ASSValue()
			if got != tt.want {
				t.Errorf("CaptionAlignment(%q).ASSValue() = %d, want %d", tt.alignment, got, tt.want)
			}
		})
	}
}

func TestGenerateCaptionFile(t *testing.T) {
	tmpDir := t.TempDir()

	events := []KeyEvent{
		{StartMs: 0, Key: "h"},
		{StartMs: 100, Key: "i"},
		{StartMs: 200, Key: "Enter"},
	}

	videoOpts := VideoOptions{PlaybackSpeed: 1.0, Input: tmpDir, Style: &StyleOptions{Width: 800, Height: 600}}
	path, err := GenerateCaptionFile(events, nil, videoOpts, DefaultCaptionOptions(), DefaultOverlayOptions())
	if err != nil {
		t.Fatalf("GenerateCaptionFile failed: %v", err)
	}

	if filepath.Dir(path) != tmpDir {
		t.Errorf("expected file in %s, got %s", tmpDir, path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read caption file: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, "[Script Info]") {
		t.Error("expected ASS header in output")
	}
	if !strings.Contains(content, "PlayResX: 800") {
		t.Error("expected PlayResX: 800")
	}
	if !strings.Contains(content, "Dialogue:") {
		t.Error("expected Dialogue lines in output")
	}
	expectedFont := DefaultCaptionOptions().Font
	if !strings.Contains(content, expectedFont) {
		t.Errorf("expected default font %q in output", expectedFont)
	}
	// Default alignment is bottom-center (ASS numpad 2)
	if !strings.Contains(content, ",2,20,20,20,") {
		t.Error("expected default alignment value 2 in style lines")
	}
}

func TestGenerateCaptionFileDefaultFont(t *testing.T) {
	tmpDir := t.TempDir()

	events := []KeyEvent{{StartMs: 0, Key: "a"}}

	// When CaptionFont is empty, should default to the OS-appropriate font
	videoOpts := VideoOptions{PlaybackSpeed: 1.0, Input: tmpDir, Style: &StyleOptions{Width: 800, Height: 600}}
	path, err := GenerateCaptionFile(events, nil, videoOpts, DefaultCaptionOptions(), DefaultOverlayOptions())
	if err != nil {
		t.Fatalf("GenerateCaptionFile failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read caption file: %v", err)
	}
	content := string(data)

	expectedFont := DefaultCaptionOptions().Font
	if !strings.Contains(content, expectedFont+",") {
		t.Errorf("expected default font %q in ASS output", expectedFont)
	}
}

func TestGenerateCaptionFileExplicitFontOverride(t *testing.T) {
	tmpDir := t.TempDir()

	events := []KeyEvent{{StartMs: 0, Key: "a"}}
	opts := DefaultCaptionOptions()
	opts.Font = "JetBrainsMonoNL NFM"

	videoOpts := VideoOptions{PlaybackSpeed: 1.0, Input: tmpDir, Style: &StyleOptions{Width: 800, Height: 600}}
	path, err := GenerateCaptionFile(events, nil, videoOpts, opts, DefaultOverlayOptions())
	if err != nil {
		t.Fatalf("GenerateCaptionFile failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read caption file: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, "JetBrainsMonoNL NFM") {
		t.Error("expected explicit CaptionFont in ASS output")
	}
	defaultFont := DefaultCaptionOptions().Font
	if strings.Contains(content, defaultFont) {
		t.Errorf("default font %q should not appear when CaptionFont is set explicitly", defaultFont)
	}
}

func TestGenerateCaptionFilePlaybackSpeed(t *testing.T) {
	tmpDir := t.TempDir()

	events := []KeyEvent{
		{StartMs: 0, Key: "a"},
		{StartMs: 2000, Key: "b"},
	}

	videoOpts := VideoOptions{PlaybackSpeed: 2.0, Input: tmpDir, Style: &StyleOptions{Width: 800, Height: 600}}
	path, err := GenerateCaptionFile(events, nil, videoOpts, DefaultCaptionOptions(), DefaultOverlayOptions())
	if err != nil {
		t.Fatalf("GenerateCaptionFile failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read caption file: %v", err)
	}
	content := string(data)

	// At 2x speed, 2000ms becomes 1000ms = 0:00:01.00
	if !strings.Contains(content, "0:00:01.00") {
		t.Error("expected timestamp adjusted for 2x playback speed")
	}
}

func TestCaptionFontSizeConsistentAcrossHeights(t *testing.T) {
	fontSize := 22
	var firstContent string
	for _, height := range []int{384, 600, 768, 1080} {
		t.Run(fmt.Sprintf("height_%d", height), func(t *testing.T) {
			tmpDir := t.TempDir()
			events := []KeyEvent{{StartMs: 0, Key: "a"}}
			opts := DefaultCaptionOptions()
			opts.FontSize = fontSize

			videoOpts := VideoOptions{
				PlaybackSpeed: 1.0,
				Input:         tmpDir,
				Style:         &StyleOptions{Width: 800, Height: height},
			}
			path, err := GenerateCaptionFile(events, nil, videoOpts, opts, DefaultOverlayOptions())
			if err != nil {
				t.Fatalf("GenerateCaptionFile failed: %v", err)
			}
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("failed to read caption file: %v", err)
			}
			content := string(data)
			// Extract the font size from the style line
			for _, line := range strings.Split(content, "\n") {
				if strings.HasPrefix(line, "Style: KeysFG,") {
					if firstContent == "" {
						firstContent = line
					} else if line != firstContent {
						t.Errorf("font size differs across heights:\n  first:   %s\n  current: %s", firstContent, line)
					}
					break
				}
			}
		})
	}
}

func TestGenerateCaptionFileEmpty(t *testing.T) {
	tmpDir := t.TempDir()

	videoOpts := VideoOptions{PlaybackSpeed: 1.0, Input: tmpDir, Style: &StyleOptions{Width: 800, Height: 600}}
	path, err := GenerateCaptionFile([]KeyEvent{}, nil, videoOpts, DefaultCaptionOptions(), DefaultOverlayOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "" {
		t.Errorf("expected empty path for no events, got %q", path)
	}
}
