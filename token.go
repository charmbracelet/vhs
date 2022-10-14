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
	DASH         = "-"
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
	JSON         = "JSON"
	TYPE         = "TYPE"
	DOWN         = "DOWN"
	LEFT         = "LEFT"
	RIGHT        = "RIGHT"
	UP           = "UP"
	TAB          = "TAB"
	ESCAPE       = "ESCAPE"
	BEGIN        = "BEGIN"
	END          = "END"
	HIDE         = "HIDE"
	SHOW         = "SHOW"
	OUTPUT       = "OUTPUT"
	MILLISECONDS = "MILLISECONDS"
	SECONDS      = "SECONDS"
	MINUTES      = "MINUTES"
	COMMENT      = "COMMENT"
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
	"FontFamily":    SETTING,
	"FontSize":      SETTING,
	"Framerate":     SETTING,
	"Height":        SETTING,
	"LetterSpacing": SETTING,
	"LineHeight":    SETTING,
	"TypingSpeed":   SETTING,
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
