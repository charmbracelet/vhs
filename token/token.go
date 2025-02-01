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
	SHIFT     = "SHIFT"

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

	HIDE                  = "HIDE"
	OUTPUT                = "OUTPUT"
	REQUIRE               = "REQUIRE"
	SET                   = "SET"
	SHOW                  = "SHOW"
	SOURCE                = "SOURCE"
	TYPE                  = "TYPE"
	TYPE_VARIABLE         = "TYPE_VARIABLE"
	SCREENSHOT            = "SCREENSHOT"
	COPY                  = "COPY"
	PASTE                 = "PASTE"
	SHELL                 = "SHELL"
	ENV                   = "ENV"
	FONT_FAMILY           = "FONT_FAMILY" //nolint:revive
	FONT_SIZE             = "FONT_SIZE"   //nolint:revive
	FRAMERATE             = "FRAMERATE"
	PLAYBACK_SPEED        = "PLAYBACK_SPEED" //nolint:revive
	HEIGHT                = "HEIGHT"
	WIDTH                 = "WIDTH"
	LETTER_SPACING        = "LETTER_SPACING" //nolint:revive
	LINE_HEIGHT           = "LINE_HEIGHT"    //nolint:revive
	TYPING_SPEED          = "TYPING_SPEED"   //nolint:revive
	TYPING_SPEED_VARIABLE = "TYPING_SPEED_VARIABLE"
	PADDING               = "PADDING"
	THEME                 = "THEME"
	LOOP_OFFSET           = "LOOP_OFFSET"     //nolint:revive
	MARGIN_FILL           = "MARGIN_FILL"     //nolint:revive
	MARGIN                = "MARGIN"          //nolint:revive
	WINDOW_BAR            = "WINDOW_BAR"      //nolint:revive
	WINDOW_BAR_SIZE       = "WINDOW_BAR_SIZE" //nolint:revive
	BORDER_RADIUS         = "CORNER_RADIUS"   //nolint:revive
	WAIT                  = "WAIT"            //nolint:revive
	WAIT_TIMEOUT          = "WAIT_TIMEOUT"    //nolint:revive
	WAIT_PATTERN          = "WAIT_PATTERN"    //nolint:revive
	CURSOR_BLINK          = "CURSOR_BLINK"    //nolint:revive
)

// Keywords maps keyword strings to tokens.
var Keywords = map[string]Type{
	"em":                  EM,
	"px":                  PX,
	"ms":                  MILLISECONDS,
	"s":                   SECONDS,
	"m":                   MINUTES,
	"Set":                 SET,
	"Sleep":               SLEEP,
	"Type":                TYPE,
	"TypeVariable":        TYPE_VARIABLE,
	"Enter":               ENTER,
	"Space":               SPACE,
	"Backspace":           BACKSPACE,
	"Delete":              DELETE,
	"Insert":              INSERT,
	"Ctrl":                CTRL,
	"Alt":                 ALT,
	"Shift":               SHIFT,
	"Down":                DOWN,
	"Left":                LEFT,
	"Right":               RIGHT,
	"Up":                  UP,
	"PageUp":              PAGEUP,
	"PageDown":            PAGEDOWN,
	"Tab":                 TAB,
	"Escape":              ESCAPE,
	"End":                 END,
	"Hide":                HIDE,
	"Require":             REQUIRE,
	"Show":                SHOW,
	"Output":              OUTPUT,
	"Shell":               SHELL,
	"FontFamily":          FONT_FAMILY,
	"MarginFill":          MARGIN_FILL,
	"Margin":              MARGIN,
	"WindowBar":           WINDOW_BAR,
	"WindowBarSize":       WINDOW_BAR_SIZE,
	"BorderRadius":        BORDER_RADIUS,
	"FontSize":            FONT_SIZE,
	"Framerate":           FRAMERATE,
	"Height":              HEIGHT,
	"LetterSpacing":       LETTER_SPACING,
	"LineHeight":          LINE_HEIGHT,
	"PlaybackSpeed":       PLAYBACK_SPEED,
	"TypingSpeed":         TYPING_SPEED,
	"TypingSpeedVariable": TYPING_SPEED_VARIABLE,
	"Padding":             PADDING,
	"Theme":               THEME,
	"Width":               WIDTH,
	"LoopOffset":          LOOP_OFFSET,
	"WaitTimeout":         WAIT_TIMEOUT,
	"WaitPattern":         WAIT_PATTERN,
	"Wait":                WAIT,
	"Source":              SOURCE,
	"CursorBlink":         CURSOR_BLINK,
	"true":                BOOLEAN,
	"false":               BOOLEAN,
	"Screenshot":          SCREENSHOT,
	"Copy":                COPY,
	"Paste":               PASTE,
	"Env":                 ENV,
}

// IsSetting returns whether a token is a setting.
func IsSetting(t Type) bool {
	switch t {
	case SHELL, FONT_FAMILY, FONT_SIZE, LETTER_SPACING, LINE_HEIGHT,
		FRAMERATE, TYPING_SPEED, TYPING_SPEED_VARIABLE, THEME, PLAYBACK_SPEED, HEIGHT, WIDTH,
		PADDING, LOOP_OFFSET, MARGIN_FILL, MARGIN, WINDOW_BAR,
		WINDOW_BAR_SIZE, BORDER_RADIUS, CURSOR_BLINK, WAIT_TIMEOUT, WAIT_PATTERN:
		return true
	default:
		return false
	}
}

// IsCommand returns whether the string is a command.
func IsCommand(t Type) bool {
	switch t {
	case TYPE, SLEEP,
		UP, DOWN, RIGHT, LEFT, PAGEUP, PAGEDOWN,
		ENTER, BACKSPACE, DELETE, TAB,
		ESCAPE, HOME, INSERT, END, CTRL, SOURCE, SCREENSHOT, COPY, PASTE, WAIT:
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
