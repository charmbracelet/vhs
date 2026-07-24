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

// OverlayEvent represents a single text overlay at a specific point in the recording.
type OverlayEvent struct {
	StartMs    int64
	DurationMs int64
	Text       string
}

// OverlayOptions holds configuration for text overlays.
type OverlayOptions struct {
	Font           string
	FontSize       int
	Alignment      CaptionAlignment
	FontColor      string // #RRGGBB hex
	BoxColor       string // #RRGGBB hex
	BoxOpacity     float64
	BoxPadding     int
	MarginLeft     int
	MarginRight    int
	MarginVertical int
}

const defaultOverlayDurationMs = 3000

// DefaultOverlayOptions returns overlay options with sensible defaults.
func DefaultOverlayOptions() OverlayOptions {
	font := defaultOSFont()
	return OverlayOptions{
		Font:           font,
		FontSize:       22,
		Alignment:      AlignTopCenter,
		FontColor:      "#FFFFFF",
		BoxColor:       "#000000",
		BoxOpacity:     0.5,
		BoxPadding:     10,
		MarginLeft:     20,
		MarginRight:    20,
		MarginVertical: 20,
	}
}

// CaptionOptions holds configuration for keystroke captioning.
type CaptionOptions struct {
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
	font := defaultOSFont()
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

func defaultOSFont() string {
	font := "monospace"
	switch runtime.GOOS {
	case "windows":
		font = "Consolas"
	case "darwin":
		font = "Menlo"
	}
	return font
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
		if i+1 < n && int(events[i+1].StartMs)-int(current.StartMs) < inactivityTimerMs {
			showUntil = int(events[i+1].StartMs)
		} else {
			showUntil = int(current.StartMs) + inactivityTimerMs
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
	// Restore ASS newline sequences that were double-escaped
	s = strings.ReplaceAll(s, `\\N`, `\N`)
	s = strings.ReplaceAll(s, `\\n`, `\n`)
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
{{- if .HasCaption}}
Style: KeysBG,{{.Font}},{{.FontSize}},&H00FFFFFF,&H000000FF,{{.BoxColor}},&H00000000,1,0,0,0,100,100,0,0,3,{{.BoxPadding}},0,{{.Alignment}},{{.MarginLeft}},{{.MarginRight}},{{.MarginVertical}},1
Style: KeysFG,{{.Font}},{{.FontSize}},{{.FontColor}},&H000000FF,&H00000000,&H00000000,1,0,0,0,100,100,0,0,1,0,0,{{.Alignment}},{{.MarginLeft}},{{.MarginRight}},{{.MarginVertical}},1
{{- end}}
{{- if .HasOverlay}}
Style: OverlayBG,{{.OverlayFont}},{{.OverlayFontSize}},&H00FFFFFF,&H000000FF,{{.OverlayBoxColor}},&H00000000,1,0,0,0,100,100,0,0,3,{{.OverlayBoxPadding}},0,{{.OverlayAlignment}},{{.OverlayMarginLeft}},{{.OverlayMarginRight}},{{.OverlayMarginVertical}},1
Style: OverlayFG,{{.OverlayFont}},{{.OverlayFontSize}},{{.OverlayFontColor}},&H000000FF,&H00000000,&H00000000,1,0,0,0,100,100,0,0,1,0,0,{{.OverlayAlignment}},{{.OverlayMarginLeft}},{{.OverlayMarginRight}},{{.OverlayMarginVertical}},1
{{- end}}

[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text`

type ASSHeaderData struct {
	ResX, ResY            int
	Font                  string
	FontColor             string
	FontSize              int
	Alignment             int
	MarginLeft            int
	MarginRight           int
	MarginVertical        int
	BoxColor              string
	BoxPadding            int
	HasCaption            bool
	HasOverlay            bool
	OverlayFont           string
	OverlayFontColor      string
	OverlayFontSize       int
	OverlayAlignment      int
	OverlayMarginLeft     int
	OverlayMarginRight    int
	OverlayMarginVertical int
	OverlayBoxColor       string
	OverlayBoxPadding     int
}

var ASSHeader = template.Must(template.New("captionASS").Parse(ASSHeaderTemplate))

// GenerateCaptionFile creates an ASS subtitle file from key events and/or overlay events and returns its path.
func GenerateCaptionFile(events []KeyEvent, overlays []OverlayEvent, videoOpts VideoOptions, opts CaptionOptions, overlayOpts OverlayOptions) (string, error) {
	if len(events) == 0 && len(overlays) == 0 {
		return "", nil
	}

	width, height := calcTermDimensions(*videoOpts.Style)
	playbackSpeed := videoOpts.PlaybackSpeed
	tempDir := videoOpts.Input

	// Normalize caption events: apply key style and adjust for playback speed.
	normalized := make([]KeyEvent, len(events))
	for i, e := range events {
		normalized[i] = KeyEvent{
			StartMs: e.StartMs,
			Key:     opts.KeyStyle.Normalize(e.Key),
		}
	}

	// Normalize overlay events: copy so we can adjust timestamps in place.
	normalizedOverlays := make([]OverlayEvent, len(overlays))
	copy(normalizedOverlays, overlays)

	// Adjust all timestamps for playback speed.
	if playbackSpeed != 0 && playbackSpeed != 1.0 {
		for i := range normalized {
			normalized[i].StartMs = int64(float64(normalized[i].StartMs) / playbackSpeed)
		}
		for i := range normalizedOverlays {
			normalizedOverlays[i].StartMs = int64(float64(normalizedOverlays[i].StartMs) / playbackSpeed)
			normalizedOverlays[i].DurationMs = int64(float64(normalizedOverlays[i].DurationMs) / playbackSpeed)
		}
	}

	assPath := filepath.Join(tempDir, "captions.ass")
	f, err := os.Create(assPath)
	if err != nil {
		return "", fmt.Errorf("failed to create caption file: %w", err)
	}
	defer f.Close() //nolint:errcheck

	data := ASSHeaderData{
		ResX:                  width,
		ResY:                  height,
		HasCaption:            len(events) > 0,
		HasOverlay:            len(overlays) > 0,
		Font:                  opts.Font,
		FontSize:              opts.FontSize,
		Alignment:             opts.Alignment.ASSValue(),
		MarginLeft:            opts.MarginLeft,
		MarginRight:           opts.MarginRight,
		MarginVertical:        opts.MarginVertical,
		BoxPadding:            opts.BoxPadding,
		FontColor:             hexToASSBGR(opts.FontColor),
		BoxColor:              hexToASSBGR(opts.BoxColor),
		OverlayFont:           overlayOpts.Font,
		OverlayFontSize:       overlayOpts.FontSize,
		OverlayAlignment:      overlayOpts.Alignment.ASSValue(),
		OverlayMarginLeft:     overlayOpts.MarginLeft,
		OverlayMarginRight:    overlayOpts.MarginRight,
		OverlayMarginVertical: overlayOpts.MarginVertical,
		OverlayBoxPadding:     overlayOpts.BoxPadding,
		OverlayFontColor:      hexToASSBGR(overlayOpts.FontColor),
		OverlayBoxColor:       hexToASSBGR(overlayOpts.BoxColor),
	}

	if err := ASSHeader.Execute(f, data); err != nil {
		return "", fmt.Errorf("failed to write caption header: %w", err)
	}
	_, _ = fmt.Fprintln(f)

	alpha := opacityToASSAlpha(opts.BoxOpacity)
	highlightColor := hexToASSBGR(opts.HighlightColor)

	for _, w := range captionWindows(normalized, opts.InactivityTimerMs, max(opts.MaxKeysOnscreen, 1)) {
		start := msToASS(int(normalized[w.Last].StartMs))
		end := msToASS(w.ShowUntil)

		// Layer 0: invisible text with seamless opaque box (no inline overrides = no seams)
		plainText := captionWindowPlainText(normalized, w)
		bgLine := fmt.Sprintf(`Dialogue: 0,%s,%s,KeysBG,,0,0,0,,{\3a&H%s&\1a&HFF&}%s`,
			start, end, alpha, plainText)
		_, _ = fmt.Fprintln(f, bgLine)

		// Layer 1: visible colored text with no box
		coloredText := captionWindowColoredText(normalized, w, highlightColor)
		fgLine := fmt.Sprintf(`Dialogue: 1,%s,%s,KeysFG,,0,0,0,,%s`,
			start, end, coloredText)
		_, _ = fmt.Fprintln(f, fgLine)
	}

	overlayAlpha := opacityToASSAlpha(overlayOpts.BoxOpacity)

	for _, o := range normalizedOverlays {
		start := msToASS(int(o.StartMs))
		end := msToASS(int(o.StartMs + o.DurationMs))
		text := assEscape(o.Text)

		// Layer 2: invisible text with opaque box
		bgLine := fmt.Sprintf(`Dialogue: 2,%s,%s,OverlayBG,,0,0,0,,{\3a&H%s&\1a&HFF&}%s`,
			start, end, overlayAlpha, text)
		_, _ = fmt.Fprintln(f, bgLine)

		// Layer 3: visible text
		fgLine := fmt.Sprintf(`Dialogue: 3,%s,%s,OverlayFG,,0,0,0,,%s`,
			start, end, text)
		_, _ = fmt.Fprintln(f, fgLine)
	}

	return assPath, nil
}
