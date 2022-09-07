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

const (
	AT           = "@"
	EQUAL        = "="
	PLUS         = "+"
	PERCENT      = "%"
	SLASH        = "/"
	DOT          = "."
	PX           = "PX"
	EM           = "EM"
	EOF          = "EOF"
	ILLEGAL      = "ILLEGAL"
	SPACE        = "SPACE"
	BACKSPACE    = "BACKSPACE"
	CTRL         = "CTRL"
	ENTER        = "ENTER"
	NUMBER       = "NUMBER"
	SET          = "SET"
	SETTING      = "SETTING"
	SLEEP        = "SLEEP"
	STRING       = "STRING"
	TYPE         = "TYPE"
	DOWN         = "DOWN"
	LEFT         = "LEFT"
	RIGHT        = "RIGHT"
	UP           = "UP"
	TAB          = "TAB"
	ESCAPE       = "ESCAPE"
	SECONDS      = "SECONDS"
	MILLISECONDS = "MILLISECONDS"
	MINUTES      = "MINUTES"
)

var keywords = map[string]TokenType{
	"em":            EM,
	"px":            PX,
	"s":             SECONDS,
	"ms":            MILLISECONDS,
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
	"FontFamily":    SETTING,
	"FontSize":      SETTING,
	"Framerate":     SETTING,
	"Height":        SETTING,
	"LetterSpacing": SETTING,
	"LineHeight":    SETTING,
	"Output":        SETTING,
	"Padding":       SETTING,
	"Theme":         SETTING,
	"Width":         SETTING,
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
