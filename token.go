package dolly

// Type represents a token's type.
type TokenType string

// Token represents a lexer token.
type Token struct {
	Type    TokenType
	Literal string
}

const (
	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"
	STRING  = "STRING"
	NUMBER  = "NUMBER"
	IDENT   = "IDENT"
	SETTING = "SETTING"

	AT      = "@"
	EQUAL   = "="
	PERCENT = "%"

	PX           = "PX"
	EM           = "EM"
	SECONDS      = "SECONDS"
	MILLISECONDS = "MILLISECONDS"
	MINUTES      = "MINUTES"

	ENTER     = "ENTER"
	SET       = "SET"
	SLEEP     = "SLEEP"
	TYPE      = "TYPE"
	BACKSPACE = "BACKSPACE"

	DOWN  = "DOWN"
	LEFT  = "LEFT"
	RIGHT = "RIGHT"
	UP    = "UP"
)

var keywords = map[string]TokenType{
	// Commands
	"Set":       SET,
	"Sleep":     SLEEP,
	"Type":      TYPE,
	"Enter":     ENTER,
	"Backspace": BACKSPACE,
	"Down":      DOWN,
	"Left":      LEFT,
	"Right":     RIGHT,
	"Up":        UP,

	// Units
	"em": EM,
	"px": PX,
	"%":  PERCENT,
	"s":  SECONDS,
	"ms": MILLISECONDS,
	"m":  MINUTES,

	// Settings
	"FontFamily": SETTING,
	"FontSize":   SETTING,
	"Framerate":  SETTING,
	"Height":     SETTING,
	"LineHeight": SETTING,
	"Padding":    SETTING,
	"Theme":      SETTING,
	"Width":      SETTING,
}

func LookupIdentifier(ident string) TokenType {
	if t, ok := keywords[ident]; ok {
		return t
	}
	return STRING
}
