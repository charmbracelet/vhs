package dolly

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
	DOT          = "."
	PX           = "PX"
	EM           = "EM"
	EOF          = "EOF"
	ILLEGAL      = "ILLEGAL"
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
	SECONDS      = "SECONDS"
	MILLISECONDS = "MILLISECONDS"
	MINUTES      = "MINUTES"
)

var keywords = map[string]TokenType{
	"em":         EM,
	"px":         PX,
	"s":          SECONDS,
	"ms":         MILLISECONDS,
	"m":          MINUTES,
	"Set":        SET,
	"Sleep":      SLEEP,
	"Type":       TYPE,
	"Enter":      ENTER,
	"Backspace":  BACKSPACE,
	"Ctrl":       CTRL,
	"Down":       DOWN,
	"Left":       LEFT,
	"Right":      RIGHT,
	"Up":         UP,
	"FontFamily": SETTING,
	"FontSize":   SETTING,
	"Framerate":  SETTING,
	"Height":     SETTING,
	"LineHeight": SETTING,
	"Output":     SETTING,
	"Padding":    SETTING,
	"Theme":      SETTING,
	"Width":      SETTING,
}

// LookupIdentifier returns whether the identifier is a keyword.
// In `dolly`, there are no _actual_ identifiers, i.e. there are no variables.
// Instead, identifiers are simply strings (i.e. bare words).
func LookupIdentifier(ident string) TokenType {
	if t, ok := keywords[ident]; ok {
		return t
	}
	return STRING
}
