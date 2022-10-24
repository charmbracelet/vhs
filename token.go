package vhs

// Type represents a token's type.
type TokenType string

// Token represents a lexer token.
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

//nolint:revive
const (
	AT             = "@"
	EQUAL          = "="
	PLUS           = "+"
	PERCENT        = "%"
	SLASH          = "/"
	DOT            = "."
	DASH           = "-"
	PX             = "PX"
	EM             = "EM"
	EOF            = "EOF"
	ILLEGAL        = "ILLEGAL"
	SPACE          = "SPACE"
	BACKSPACE      = "BACKSPACE"
	CTRL           = "CTRL"
	ENTER          = "ENTER"
	NUMBER         = "NUMBER"
	SET            = "SET"
	SLEEP          = "SLEEP"
	STRING         = "STRING"
	JSON           = "JSON"
	TYPE           = "TYPE"
	DOWN           = "DOWN"
	LEFT           = "LEFT"
	RIGHT          = "RIGHT"
	UP             = "UP"
	TAB            = "TAB"
	ESCAPE         = "ESCAPE"
	BEGIN          = "BEGIN"
	END            = "END"
	HIDE           = "HIDE"
	SHOW           = "SHOW"
	OUTPUT         = "OUTPUT"
	MILLISECONDS   = "MILLISECONDS"
	SECONDS        = "SECONDS"
	MINUTES        = "MINUTES"
	COMMENT        = "COMMENT"
	FONT_FAMILY    = "FONT_FAMILY"
	FONT_SIZE      = "FONT_SIZE"
	FRAMERATE      = "FRAMERATE"
	HEIGHT         = "HEIGHT"
	WIDTH          = "WIDTH"
	LETTER_SPACING = "LETTER_SPACING"
	LINE_HEIGHT    = "LINE_HEIGHT"
	TYPING_SPEED   = "TYPING_SPEED"
	PADDING        = "PADDING"
	THEME          = "THEME"
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
	"Down":          DOWN,
	"Left":          LEFT,
	"Right":         RIGHT,
	"Up":            UP,
	"Tab":           TAB,
	"Escape":        ESCAPE,
	"Begin":         BEGIN,
	"End":           END,
	"Hide":          HIDE,
	"Show":          SHOW,
	"Output":        OUTPUT,
	"FontFamily":    FONT_FAMILY,
	"FontSize":      FONT_SIZE,
	"Framerate":     FRAMERATE,
	"Height":        HEIGHT,
	"LetterSpacing": LETTER_SPACING,
	"LineHeight":    LINE_HEIGHT,
	"TypingSpeed":   TYPING_SPEED,
	"Padding":       PADDING,
	"Theme":         THEME,
	"Width":         WIDTH,
}

// IsSetting returns whether a token is a setting.
func IsSetting(t TokenType) bool {
	switch t {
	case FONT_FAMILY, FONT_SIZE, LETTER_SPACING, LINE_HEIGHT,
		FRAMERATE, TYPING_SPEED, THEME,
		HEIGHT, WIDTH, PADDING:
		return true
	default:
		return false
	}
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
