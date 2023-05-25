package main

import "strings"

// TokenType represents a token's type.
type TokenType string

// Token represents a lexer token.
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// Tokens for the VHS language
const (
	AT      = "@"
	EQUAL   = "="
	PLUS    = "+"
	PERCENT = "%"
	SLASH   = "/"
	DOT     = "."
	DASH    = "-"

	EM           = "EM"
	MILLISECONDS = "MILLISECONDS"
	MINUTES      = "MINUTES"
	PX           = "PX"
	SECONDS      = "SECONDS"

	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"

	ALT       = "ALT"
	BACKSPACE = "BACKSPACE"
	CTRL      = "CTRL"
	DELETE    = "DELETE"
	END       = "END"
	ENTER     = "ENTER"
	ESCAPE    = "ESCAPE"
	HOME      = "HOME"
	INSERT    = "INSERT"
	PAGEDOWN  = "PAGEDOWN"
	PAGEUP    = "PAGEUP"
	SLEEP     = "SLEEP"
	SPACE     = "SPACE"
	TAB       = "TAB"

	COMMENT = "COMMENT"
	NUMBER  = "NUMBER"
	STRING  = "STRING"
	JSON    = "JSON"
	REGEX   = "REGEX"

	DOWN  = "DOWN"
	LEFT  = "LEFT"
	RIGHT = "RIGHT"
	UP    = "UP"

	HIDE    = "HIDE"
	OUTPUT  = "OUTPUT"
	REQUIRE = "REQUIRE"
	SET     = "SET"
	SHOW    = "SHOW"
	SOURCE  = "SOURCE"
	TYPE    = "TYPE"

	SHELL           = "SHELL"
	FONT_FAMILY     = "FONT_FAMILY" //nolint:revive
	FONT_SIZE       = "FONT_SIZE"   //nolint:revive
	FRAMERATE       = "FRAMERATE"
	PLAYBACK_SPEED  = "PLAYBACK_SPEED" //nolint:revive
	HEIGHT          = "HEIGHT"
	WIDTH           = "WIDTH"
	LETTER_SPACING  = "LETTER_SPACING" //nolint:revive
	LINE_HEIGHT     = "LINE_HEIGHT"    //nolint:revive
	TYPING_SPEED    = "TYPING_SPEED"   //nolint:revive
	PADDING         = "PADDING"
	THEME           = "THEME"
	LOOP_OFFSET     = "LOOP_OFFSET"     //nolint:revive
	MARGIN_FILL     = "MARGIN_FILL"     //nolint:revive
	MARGIN          = "MARGIN"          //nolint:revive
	WINDOW_BAR      = "WINDOW_BAR"      //nolint:revive
	WINDOW_BAR_SIZE = "WINDOW_BAR_SIZE" //nolint:revive
	BORDER_RADIUS   = "CORNER_RADIUS"   //nolint:revive
	WAIT            = "WAIT"            //nolint:revive
	WAIT_TIMEOUT    = "WAIT_TIMEOUT"    //nolint:revive
	WAIT_PATTERN    = "WAIT_PATTERN"    //nolint:revive
)

var keywords = map[string]TokenType{
	"em":            EM,
	"px":            PX,
	"ms":            MILLISECONDS,
	"s":             SECONDS,
	"m":             MINUTES,
	"Set":           SET,
	"Sleep":         SLEEP,
	"Type":          TYPE,
	"Enter":         ENTER,
	"Space":         SPACE,
	"Backspace":     BACKSPACE,
	"Ctrl":          CTRL,
	"Alt":           ALT,
	"Down":          DOWN,
	"Left":          LEFT,
	"Right":         RIGHT,
	"Up":            UP,
	"PageUp":        PAGEUP,
	"PageDown":      PAGEDOWN,
	"Tab":           TAB,
	"Escape":        ESCAPE,
	"End":           END,
	"Hide":          HIDE,
	"Require":       REQUIRE,
	"Show":          SHOW,
	"Output":        OUTPUT,
	"Shell":         SHELL,
	"FontFamily":    FONT_FAMILY,
	"MarginFill":    MARGIN_FILL,
	"Margin":        MARGIN,
	"WindowBar":     WINDOW_BAR,
	"WindowBarSize": WINDOW_BAR_SIZE,
	"BorderRadius":  BORDER_RADIUS,
	"FontSize":      FONT_SIZE,
	"Framerate":     FRAMERATE,
	"Height":        HEIGHT,
	"LetterSpacing": LETTER_SPACING,
	"LineHeight":    LINE_HEIGHT,
	"PlaybackSpeed": PLAYBACK_SPEED,
	"TypingSpeed":   TYPING_SPEED,
	"Padding":       PADDING,
	"Theme":         THEME,
	"Width":         WIDTH,
	"LoopOffset":    LOOP_OFFSET,
	"WaitTimeout":   WAIT_TIMEOUT,
	"WaitPattern":   WAIT_PATTERN,
	"Wait":          WAIT,
	"Source":        SOURCE,
}

// IsSetting returns whether a token is a setting.
func IsSetting(t TokenType) bool {
	switch t {
	case SHELL, FONT_FAMILY, FONT_SIZE, LETTER_SPACING, LINE_HEIGHT,
		FRAMERATE, TYPING_SPEED, THEME, PLAYBACK_SPEED, HEIGHT, WIDTH,
		PADDING, LOOP_OFFSET, MARGIN_FILL, MARGIN, WINDOW_BAR,
		WINDOW_BAR_SIZE, BORDER_RADIUS, WAIT_TIMEOUT, WAIT_PATTERN:
		return true
	default:
		return false
	}
}

// IsCommand returns whether the string is a command
func IsCommand(t TokenType) bool {
	switch t {
	case TYPE, SLEEP,
		UP, DOWN, RIGHT, LEFT, PAGEUP, PAGEDOWN,
		ENTER, BACKSPACE, DELETE, TAB,
		ESCAPE, HOME, INSERT, END, CTRL, SOURCE, WAIT:
		return true
	default:
		return false
	}
}

// String converts a token to it's human readable string format.
func (t TokenType) String() string {
	if IsCommand(t) || IsSetting(t) {
		return toCamel(string(t))
	}
	return string(t)
}

func toCamel(s string) string {
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
func LookupIdentifier(ident string) TokenType {
	if t, ok := keywords[ident]; ok {
		return t
	}
	return STRING
}
