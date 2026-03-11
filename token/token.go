// Package token provides the token types and structures for the VHS Tape
// language.
package token

import (
	"strings"
)

// Type represents a token's type.
type Type string

// Token represents a lexer token.
type Token struct {
	Type    Type
	Literal string
	Line    int
	Column  int
}

// Tokens for the VHS language.
const (
	AT        = "@"
	EQUAL     = "="
	PLUS      = "+"
	PERCENT   = "%"
	SLASH     = "/"
	BACKSLASH = "\\"
	DOT       = "."
	DASH      = "-"

	MINUS         = "-"
	RIGHT_BRACKET = "]"
	LEFT_BRACKET  = "["
	CARET         = "^"

	EM           = "EM"
	MILLISECONDS = "MILLISECONDS"
	MINUTES      = "MINUTES"
	PX           = "PX"
	SECONDS      = "SECONDS"

	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"

	ALT         = "ALT"
	BACKSPACE   = "BACKSPACE"
	CTRL        = "CTRL"
	DELETE      = "DELETE"
	END         = "END"
	ENTER       = "ENTER"
	ESCAPE      = "ESCAPE"
	HOME        = "HOME"
	INSERT      = "INSERT"
	PAGE_DOWN   = "PAGE_DOWN"
	PAGE_UP     = "PAGE_UP"
	SCROLL_DOWN = "SCROLL_DOWN"
	SCROLL_UP   = "SCROLL_UP"
	SLEEP       = "SLEEP"
	SPACE       = "SPACE"
	TAB         = "TAB"
	SHIFT       = "SHIFT"

	COMMENT = "COMMENT"
	NUMBER  = "NUMBER"
	STRING  = "STRING"
	JSON    = "JSON"
	REGEX   = "REGEX"
	BOOLEAN = "BOOLEAN"

	DOWN  = "DOWN"
	LEFT  = "LEFT"
	RIGHT = "RIGHT"
	UP    = "UP"

	HIDE            = "HIDE"
	OUTPUT          = "OUTPUT"
	REQUIRE         = "REQUIRE"
	SET             = "SET"
	SHOW            = "SHOW"
	SOURCE          = "SOURCE"
	TYPE            = "TYPE"
	SCREENSHOT      = "SCREENSHOT"
	COPY            = "COPY"
	PASTE           = "PASTE"
	SHELL           = "SHELL"
	ENV             = "ENV"
	FONT_FAMILY     = "FONT_FAMILY"
	FONT_SIZE       = "FONT_SIZE"
	FRAMERATE       = "FRAMERATE"
	PLAYBACK_SPEED  = "PLAYBACK_SPEED"
	HEIGHT          = "HEIGHT"
	WIDTH           = "WIDTH"
	LETTER_SPACING  = "LETTER_SPACING"
	LINE_HEIGHT     = "LINE_HEIGHT"
	TYPING_SPEED    = "TYPING_SPEED"
	PADDING         = "PADDING"
	THEME           = "THEME"
	LOOP_OFFSET     = "LOOP_OFFSET"
	MARGIN_FILL     = "MARGIN_FILL"
	MARGIN          = "MARGIN"
	WINDOW_BAR      = "WINDOW_BAR"
	WINDOW_BAR_SIZE = "WINDOW_BAR_SIZE"
	BORDER_RADIUS   = "CORNER_RADIUS"
	WAIT            = "WAIT"
	WAIT_TIMEOUT    = "WAIT_TIMEOUT"
	WAIT_PATTERN    = "WAIT_PATTERN"
	CURSOR_BLINK    = "CURSOR_BLINK"

	CAPTION_ON               = "CAPTION_ON"
	CAPTION_OFF              = "CAPTION_OFF"
	CAPTION_FONT             = "CAPTION_FONT"
	CAPTION_FONT_SIZE        = "CAPTION_FONT_SIZE"
	CAPTION_MAX_KEYS         = "CAPTION_MAX_KEYS"
	CAPTION_INACTIVITY_TIMER = "CAPTION_INACTIVITY_TIMER"
	CAPTION_HIGHLIGHT_COLOR  = "CAPTION_HIGHLIGHT_COLOR"
	CAPTION_FONT_COLOR       = "CAPTION_FONT_COLOR"
	CAPTION_BOX_COLOR        = "CAPTION_BOX_COLOR"
	CAPTION_BOX_OPACITY      = "CAPTION_BOX_OPACITY"
	CAPTION_KEY_STYLE        = "CAPTION_KEY_STYLE"
	CAPTION_MARGIN_LEFT      = "CAPTION_MARGIN_LEFT"
	CAPTION_MARGIN_RIGHT     = "CAPTION_MARGIN_RIGHT"
	CAPTION_MARGIN_VERTICAL  = "CAPTION_MARGIN_VERTICAL"
	CAPTION_ALIGNMENT        = "CAPTION_ALIGNMENT"
	CAPTION_BOX_PADDING      = "CAPTION_BOX_PADDING"

	OVERLAY                 = "OVERLAY"
	OVERLAY_FONT            = "OVERLAY_FONT"
	OVERLAY_FONT_SIZE       = "OVERLAY_FONT_SIZE"
	OVERLAY_FONT_COLOR      = "OVERLAY_FONT_COLOR"
	OVERLAY_BOX_COLOR       = "OVERLAY_BOX_COLOR"
	OVERLAY_BOX_OPACITY     = "OVERLAY_BOX_OPACITY"
	OVERLAY_BOX_PADDING     = "OVERLAY_BOX_PADDING"
	OVERLAY_ALIGNMENT       = "OVERLAY_ALIGNMENT"
	OVERLAY_MARGIN_LEFT     = "OVERLAY_MARGIN_LEFT"
	OVERLAY_MARGIN_RIGHT    = "OVERLAY_MARGIN_RIGHT"
	OVERLAY_MARGIN_VERTICAL = "OVERLAY_MARGIN_VERTICAL"
)

// Keywords maps keyword strings to tokens.
var Keywords = map[string]Type{
	"em":                     EM,
	"px":                     PX,
	"ms":                     MILLISECONDS,
	"s":                      SECONDS,
	"m":                      MINUTES,
	"Set":                    SET,
	"Sleep":                  SLEEP,
	"Type":                   TYPE,
	"Enter":                  ENTER,
	"Space":                  SPACE,
	"Backspace":              BACKSPACE,
	"Delete":                 DELETE,
	"Insert":                 INSERT,
	"Ctrl":                   CTRL,
	"Alt":                    ALT,
	"Shift":                  SHIFT,
	"Down":                   DOWN,
	"Left":                   LEFT,
	"Right":                  RIGHT,
	"Up":                     UP,
	"PageUp":                 PAGE_UP,
	"PageDown":               PAGE_DOWN,
	"ScrollUp":               SCROLL_UP,
	"ScrollDown":             SCROLL_DOWN,
	"Tab":                    TAB,
	"Escape":                 ESCAPE,
	"End":                    END,
	"Hide":                   HIDE,
	"Require":                REQUIRE,
	"Show":                   SHOW,
	"Output":                 OUTPUT,
	"Shell":                  SHELL,
	"FontFamily":             FONT_FAMILY,
	"MarginFill":             MARGIN_FILL,
	"Margin":                 MARGIN,
	"WindowBar":              WINDOW_BAR,
	"WindowBarSize":          WINDOW_BAR_SIZE,
	"BorderRadius":           BORDER_RADIUS,
	"FontSize":               FONT_SIZE,
	"Framerate":              FRAMERATE,
	"Height":                 HEIGHT,
	"LetterSpacing":          LETTER_SPACING,
	"LineHeight":             LINE_HEIGHT,
	"PlaybackSpeed":          PLAYBACK_SPEED,
	"TypingSpeed":            TYPING_SPEED,
	"Padding":                PADDING,
	"Theme":                  THEME,
	"Width":                  WIDTH,
	"LoopOffset":             LOOP_OFFSET,
	"WaitTimeout":            WAIT_TIMEOUT,
	"WaitPattern":            WAIT_PATTERN,
	"Wait":                   WAIT,
	"Source":                 SOURCE,
	"CursorBlink":            CURSOR_BLINK,
	"CaptionOn":              CAPTION_ON,
	"CaptionOff":             CAPTION_OFF,
	"CaptionFont":            CAPTION_FONT,
	"CaptionFontSize":        CAPTION_FONT_SIZE,
	"CaptionMaxKeys":         CAPTION_MAX_KEYS,
	"CaptionInactivityTimer": CAPTION_INACTIVITY_TIMER,
	"CaptionHighlightColor":  CAPTION_HIGHLIGHT_COLOR,
	"CaptionFontColor":       CAPTION_FONT_COLOR,
	"CaptionBoxColor":        CAPTION_BOX_COLOR,
	"CaptionBoxOpacity":      CAPTION_BOX_OPACITY,
	"CaptionKeyStyle":        CAPTION_KEY_STYLE,
	"CaptionMarginLeft":      CAPTION_MARGIN_LEFT,
	"CaptionMarginRight":     CAPTION_MARGIN_RIGHT,
	"CaptionMarginVertical":  CAPTION_MARGIN_VERTICAL,
	"CaptionAlignment":       CAPTION_ALIGNMENT,
	"CaptionBoxPadding":      CAPTION_BOX_PADDING,
	"Overlay":                OVERLAY,
	"OverlayFont":            OVERLAY_FONT,
	"OverlayFontSize":        OVERLAY_FONT_SIZE,
	"OverlayFontColor":       OVERLAY_FONT_COLOR,
	"OverlayBoxColor":        OVERLAY_BOX_COLOR,
	"OverlayBoxOpacity":      OVERLAY_BOX_OPACITY,
	"OverlayBoxPadding":      OVERLAY_BOX_PADDING,
	"OverlayAlignment":       OVERLAY_ALIGNMENT,
	"OverlayMarginLeft":      OVERLAY_MARGIN_LEFT,
	"OverlayMarginRight":     OVERLAY_MARGIN_RIGHT,
	"OverlayMarginVertical":  OVERLAY_MARGIN_VERTICAL,
	"true":                   BOOLEAN,
	"false":                  BOOLEAN,
	"Screenshot":             SCREENSHOT,
	"Copy":                   COPY,
	"Paste":                  PASTE,
	"Env":                    ENV,
}

// IsSetting returns whether a token is a setting.
func IsSetting(t Type) bool {
	switch t {
	case SHELL, FONT_FAMILY, FONT_SIZE, LETTER_SPACING, LINE_HEIGHT,
		FRAMERATE, TYPING_SPEED, THEME, PLAYBACK_SPEED, HEIGHT, WIDTH,
		PADDING, LOOP_OFFSET, MARGIN_FILL, MARGIN, WINDOW_BAR,
		WINDOW_BAR_SIZE, BORDER_RADIUS, CURSOR_BLINK, WAIT_TIMEOUT, WAIT_PATTERN,
		CAPTION_FONT, CAPTION_FONT_SIZE, CAPTION_MAX_KEYS, CAPTION_INACTIVITY_TIMER,
		CAPTION_HIGHLIGHT_COLOR, CAPTION_FONT_COLOR, CAPTION_BOX_COLOR,
		CAPTION_BOX_OPACITY, CAPTION_KEY_STYLE,
		CAPTION_MARGIN_LEFT, CAPTION_MARGIN_RIGHT, CAPTION_MARGIN_VERTICAL,
		CAPTION_ALIGNMENT, CAPTION_BOX_PADDING,
		OVERLAY_FONT, OVERLAY_FONT_SIZE, OVERLAY_FONT_COLOR,
		OVERLAY_BOX_COLOR, OVERLAY_BOX_OPACITY, OVERLAY_BOX_PADDING,
		OVERLAY_ALIGNMENT, OVERLAY_MARGIN_LEFT, OVERLAY_MARGIN_RIGHT,
		OVERLAY_MARGIN_VERTICAL:
		return true
	default:
		return false
	}
}

// IsCommand returns whether the string is a command.
func IsCommand(t Type) bool {
	switch t {
	case TYPE, SLEEP,
		UP, DOWN, RIGHT, LEFT, PAGE_UP, PAGE_DOWN, SCROLL_UP, SCROLL_DOWN,
		ENTER, BACKSPACE, DELETE, TAB,
		ESCAPE, HOME, INSERT, END, CTRL, SOURCE, SCREENSHOT, COPY, PASTE, WAIT,
		OVERLAY, CAPTION_ON, CAPTION_OFF:
		return true
	default:
		return false
	}
}

// IsModifier returns whether the token is a modifier.
func IsModifier(t Type) bool {
	return t == ALT || t == SHIFT
}

// String converts a token to it's human readable string format.
func (t Type) String() string {
	if IsCommand(t) || IsSetting(t) {
		return ToCamel(string(t))
	}
	return string(t)
}

// ToCamel converts a snake_case string to CamelCase.
func ToCamel(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		p = strings.ToUpper(p[:1]) + strings.ToLower(p[1:])
		parts[i] = p
	}

	return strings.Join(parts, "")
}

// LookupIdentifier returns whether the identifier is a keyword.
// In `vhs`, there are no _actual_ identifiers, i.e. there are no variables.
// Instead, identifiers are simply strings (i.e. bare words).
func LookupIdentifier(ident string) Type {
	if t, ok := Keywords[ident]; ok {
		return t
	}
	return STRING
}
