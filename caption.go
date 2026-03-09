package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

// CaptionOptions holds configuration for keystroke captioning.
type CaptionOptions struct {
	Enabled           bool
	Font              string
	FontSize          int
	KeyStyle          KeyStyle
	MaxKeysOnscreen   int
	InactivityTimerMs int
	Alignment         CaptionAlignment
	FontColor         string // #RRGGBB hex
	HighlightColor    string // #RRGGBB hex
	BoxColor          string // #RRGGBB hex
	BoxOpacity        float64
	BoxPadding        int
	MarginLeft        int
	MarginRight       int
	MarginVertical    int
}

// DefaultCaptionOptions returns caption options with sensible defaults.
func DefaultCaptionOptions() CaptionOptions {
	font := "monospace"
	switch runtime.GOOS {
	case "windows":
		font = "Consolas"
	case "darwin":
		font = "Menlo"
	}
	return CaptionOptions{
		Font:              font,
		FontSize:          22,
		KeyStyle:          KeyStyleIcon,
		MaxKeysOnscreen:   10,
		InactivityTimerMs: 1000,
		Alignment:         AlignBottomCenter,
		FontColor:         "#FFFFFF",
		HighlightColor:    "#FFCC66",
		BoxColor:          "#000000",
		BoxOpacity:        0.5,
		BoxPadding:        10,
		MarginLeft:        20,
		MarginRight:       20,
		MarginVertical:    20,
	}
}

// KeyStyle represents the rendering style for key captions.
type KeyStyle string

const (
	KeyStyleVim  KeyStyle = "vim"
	KeyStyleIcon KeyStyle = "icon"
)

// CaptionAlignment controls where captions appear on screen.
type CaptionAlignment string

const (
	AlignBottomLeft   CaptionAlignment = "bottom-left"
	AlignBottomCenter CaptionAlignment = "bottom-center"
	AlignBottomRight  CaptionAlignment = "bottom-right"
	AlignMiddleLeft   CaptionAlignment = "middle-left"
	AlignMiddleCenter CaptionAlignment = "middle-center"
	AlignMiddleRight  CaptionAlignment = "middle-right"
	AlignTopLeft      CaptionAlignment = "top-left"
	AlignTopCenter    CaptionAlignment = "top-center"
	AlignTopRight     CaptionAlignment = "top-right"
)

// CaptionAlignments maps alignment names to ASS numpad values (1–9).
var CaptionAlignments = map[CaptionAlignment]int{
	AlignBottomLeft:   1,
	AlignBottomCenter: 2,
	AlignBottomRight:  3,
	AlignMiddleLeft:   4,
	AlignMiddleCenter: 5,
	AlignMiddleRight:  6,
	AlignTopLeft:      7,
	AlignTopCenter:    8,
	AlignTopRight:     9,
}

// ASSValue returns the ASS numpad alignment integer (1–9).
func (a CaptionAlignment) ASSValue() int {
	if v, ok := CaptionAlignments[a]; ok {
		return v
	}
	return CaptionAlignments[AlignBottomRight]
}

// Normalize transforms a key name (e.g. "Ctrl+Shift+a") into styled display text.
func (s KeyStyle) Normalize(key string) string {
	parts := strings.Split(key, "+")
	overrides := vimOverrides
	if s == KeyStyleIcon {
		overrides = iconOverrides
	}
	for i, p := range parts {
		if o, ok := overrides[p]; ok {
			parts[i] = o
		}
	}
	if s == KeyStyleIcon {
		return strings.Join(parts, "")
	}
	if len(parts) == 1 && len([]rune(parts[0])) == 1 {
		return parts[0]
	}
	return "<" + strings.Join(parts, "-") + ">"
}

var vimOverrides = map[string]string{
	"Backspace": "BS",
	"Delete":    "Del",
	"Enter":     "CR",
	"Escape":    "Esc",
	"Ctrl":      "C",
	"Alt":       "A",
	"Shift":     "S",
}

var iconOverrides = map[string]string{
	"Backspace": "⌫",
	"Delete":    "⌦",
	"Ctrl":      "^",
	"Alt":       "⌥",
	"Shift":     "⇧",
	"Down":      "↓",
	"PageDown":  "PgDn",
	"Up":        "↑",
	"PageUp":    "PgUp",
	"Left":      "←",
	"Right":     "→",
	"Space":     "␣",
	"Enter":     "↵",
	"Escape":    "Esc",
	"Tab":       "⇥",
}

// CaptionWindow describes the sliding key display at the moment each key is pressed.
type CaptionWindow struct {
	First, Last int
	ShowUntil   int
	IsTruncated bool
}

// captionWindows returns one CaptionWindow per event. Windows are computed
// such that keys remain on screen until inactivityTimerMs has been exceeded.
// Further, there will be no more than maxKeys shown on the screen at any
// point in time.
func captionWindows(events []KeyEvent, inactivityTimerMs, maxKeys int) []CaptionWindow {
	n := len(events)
	wins := make([]CaptionWindow, 0, n)
	sessionStart := 0

	for i := range n {
		current := events[i]
		windowStart := max(sessionStart, i-maxKeys+1)
		isTruncated := windowStart > sessionStart

		var showUntil int
		if i+1 < n && int(events[i+1].Ms)-int(current.Ms) < inactivityTimerMs {
			showUntil = int(events[i+1].Ms)
		} else {
			showUntil = int(current.Ms) + inactivityTimerMs
			sessionStart = i + 1
		}

		wins = append(wins, CaptionWindow{
			First:       windowStart,
			Last:        i,
			ShowUntil:   showUntil,
			IsTruncated: isTruncated,
		})
	}
	return wins
}

// ASS helpers

func msToASS(ms int) string {
	h := ms / 3600000
	ms %= 3600000
	m := ms / 60000
	ms %= 60000
	s := ms / 1000
	ms %= 1000
	cs := ms / 10
	return fmt.Sprintf("%d:%02d:%02d.%02d", h, m, s, cs)
}

func assEscape(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, "{", `\{`)
	s = strings.ReplaceAll(s, "}", `\}`)
	return s
}

func opacityToASSAlpha(opacity float64) string {
	alpha := int(math.Round((1 - opacity) * 255))
	return fmt.Sprintf("%02X", alpha)
}

func hexToASSBGR(hex string) string {
	hex = strings.TrimPrefix(hex, "#")
	r, g, b := hex[0:2], hex[2:4], hex[4:6]
	return fmt.Sprintf("&H%s%s%s&", b, g, r)
}

// captionWindowPlainText builds plain text with no inline overrides.
// Used for the background box layer so the opaque box renders without seams.
func captionWindowPlainText(events []KeyEvent, w CaptionWindow) string {
	var parts []string
	if w.IsTruncated {
		parts = append(parts, "…")
	}
	for i := w.First; i <= w.Last; i++ {
		parts = append(parts, assEscape(events[i].Key))
	}
	return strings.Join(parts, " ")
}

// captionWindowColoredText builds ASS text with the last key highlighted.
// Used for the visible text layer (no opaque box, so color overrides cause no seams).
func captionWindowColoredText(events []KeyEvent, w CaptionWindow, highlightColor string) string {
	var parts []string
	if w.IsTruncated {
		parts = append(parts, "…")
	}
	for i := w.First; i < w.Last; i++ {
		parts = append(parts, assEscape(events[i].Key))
	}
	lastKey := assEscape(events[w.Last].Key)
	parts = append(parts, fmt.Sprintf(`{\c%s}%s`, highlightColor, lastKey))
	return strings.Join(parts, " ")
}

const ASSHeaderTemplate = `[Script Info]
ScriptType: v4.00+
PlayResX: {{.ResX}}
PlayResY: {{.ResY}}
ScaledBorderAndShadow: yes

[V4+ Styles]
Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding
Style: KeysBG,{{.Font}},{{.FontSize}},&H00FFFFFF,&H000000FF,{{.BoxColor}},&H00000000,1,0,0,0,100,100,0,0,3,{{.BoxPadding}},0,{{.Alignment}},{{.MarginLeft}},{{.MarginRight}},{{.MarginVertical}},1
Style: KeysFG,{{.Font}},{{.FontSize}},{{.FontColor}},&H000000FF,&H00000000,&H00000000,1,0,0,0,100,100,0,0,1,0,0,{{.Alignment}},{{.MarginLeft}},{{.MarginRight}},{{.MarginVertical}},1

[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text`

type ASSHeaderData struct {
	ResX, ResY     int
	Font           string
	FontColor      string
	FontSize       int
	Alignment      int
	MarginLeft     int
	MarginRight    int
	MarginVertical int
	BoxColor       string
	BoxPadding     int
}

var ASSHeader = template.Must(template.New("captionASS").Parse(ASSHeaderTemplate))

// GenerateCaptionFile creates an ASS subtitle file from key events and returns its path.
func GenerateCaptionFile(events []KeyEvent, videoOpts VideoOptions, opts CaptionOptions) (string, error) {
	if len(events) == 0 {
		return "", nil
	}

	width, height := calcTermDimensions(*videoOpts.Style)
	playbackSpeed := videoOpts.PlaybackSpeed
	tempDir := videoOpts.Input

	normalized := make([]KeyEvent, len(events))
	for i, e := range events {
		normalized[i] = KeyEvent{
			Ms:  e.Ms,
			Key: opts.KeyStyle.Normalize(e.Key),
		}
	}

	// Adjust timestamps for playback speed
	if playbackSpeed != 0 && playbackSpeed != 1.0 {
		for i := range normalized {
			normalized[i].Ms = int64(float64(normalized[i].Ms) / playbackSpeed)
		}
	}

	assPath := filepath.Join(tempDir, "captions.ass")
	f, err := os.Create(assPath)
	if err != nil {
		return "", fmt.Errorf("failed to create caption file: %w", err)
	}
	defer f.Close()

	data := ASSHeaderData{
		ResX:           width,
		ResY:           height,
		Font:           opts.Font,
		FontSize:       opts.FontSize,
		Alignment:      opts.Alignment.ASSValue(),
		MarginLeft:     opts.MarginLeft,
		MarginRight:    opts.MarginRight,
		MarginVertical: opts.MarginVertical,
		BoxPadding:     opts.BoxPadding,
		FontColor:      hexToASSBGR(opts.FontColor),
		BoxColor:       hexToASSBGR(opts.BoxColor),
	}
	if err := ASSHeader.Execute(f, data); err != nil {
		return "", fmt.Errorf("failed to write caption header: %w", err)
	}
	fmt.Fprintln(f)

	alpha := opacityToASSAlpha(opts.BoxOpacity)
	highlightColor := hexToASSBGR(opts.HighlightColor)

	for _, w := range captionWindows(normalized, opts.InactivityTimerMs, max(opts.MaxKeysOnscreen, 1)) {
		start := msToASS(int(normalized[w.Last].Ms))
		end := msToASS(w.ShowUntil)

		// Layer 0: invisible text with seamless opaque box (no inline overrides = no seams)
		plainText := captionWindowPlainText(normalized, w)
		bgLine := fmt.Sprintf(`Dialogue: 0,%s,%s,KeysBG,,0,0,0,,{\3a&H%s&\1a&HFF&}%s`,
			start, end, alpha, plainText)
		fmt.Fprintln(f, bgLine)

		// Layer 1: visible colored text with no box
		coloredText := captionWindowColoredText(normalized, w, highlightColor)
		fgLine := fmt.Sprintf(`Dialogue: 1,%s,%s,KeysFG,,0,0,0,,%s`,
			start, end, coloredText)
		fmt.Fprintln(f, fgLine)
	}

	return assPath, nil
}
