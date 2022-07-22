package dolly

import (
	"math/rand"
	"time"

	"github.com/go-rod/rod/lib/input"
)

// TypeOptions are the possible typing options.
type TypeOptions struct {
	Speed    float64
	Variance float64
}

// DefaultTypeOptions returns the default typing options.
func DefaultTypeOptions() TypeOptions {
	return TypeOptions{
		Speed:    75,
		Variance: 0.1,
	}
}

// TypeOption is a typing option.
type TypeOption func(*TypeOptions)

// WithSpeed sets the typing speed.
func WithSpeed(speed float64) TypeOption {
	return func(o *TypeOptions) { o.Speed = speed }
}

// WithVariance sets the typing speed variance.
func WithVariance(variance float64) TypeOption {
	return func(o *TypeOptions) { o.Variance = variance }
}

// Type types the given string onto the page at the given speed. The delay is
// the time between each key press.
func (d Dolly) Type(str string, opts ...TypeOption) {
	options := DefaultTypeOptions()
	for _, opt := range opts {
		opt(&options)
	}

	for _, r := range str {
		k, ok := keymap[r]
		if ok {
			d.Page.Keyboard.Type(k)
		} else {
			d.Page.MustElement("textarea").Input(string(r))
			d.Page.MustWaitIdle()
		}

		r := (rand.Float64() - 0.5)
		v := r * (options.Variance * options.Speed)
		time.Sleep(time.Millisecond * time.Duration(v+options.Speed))
	}
}

// Enter is a helper function that press the enter key.
func (d Dolly) Enter() { d.Page.Keyboard.Type(input.Enter) }

// Execute executes a command in the terminal without showing output and clears
// the screen.
func (d Dolly) Execute(cmd string) {
	d.Type(cmd, WithSpeed(0))
	d.Enter()
	d.Clear()
}

// WithCtrl presses a key with the ctrl key held down.
func (d Dolly) WithCtrl(k input.Key) {
	d.Page.Keyboard.Press(input.ControlLeft)
	d.Page.Keyboard.Type(k)
	d.Page.Keyboard.Release(input.ControlLeft)
}

// Clear is a helper function that clears the screen.
// Must be currently on shell to work (not inside input / program)
func (d Dolly) Clear() { d.WithCtrl(shift(input.KeyL)) }

// CtrlU is a helper function that presses the ctrl-u key.
func (d Dolly) CtrlU() { d.WithCtrl(input.KeyU) }

// CtrlC is a helper function that presses the ctrl-c key.
func (d Dolly) CtrlC() { d.WithCtrl(shift(input.KeyC)) }

func shift(k input.Key) input.Key {
	k, _ = k.Shift()
	return k
}

var keymap = map[rune]input.Key{
	' ':    input.Space,
	'!':    shift(input.Digit1),
	'"':    shift(input.Quote),
	'#':    shift(input.Digit3),
	'$':    shift(input.Digit4),
	'%':    shift(input.Digit5),
	'&':    shift(input.Digit7),
	'(':    shift(input.Digit9),
	')':    shift(input.Digit0),
	'*':    shift(input.Digit8),
	'+':    shift(input.Equal),
	',':    input.Comma,
	'-':    input.Minus,
	'.':    input.Period,
	'/':    input.Slash,
	'0':    input.Digit0,
	'1':    input.Digit1,
	'2':    input.Digit2,
	'3':    input.Digit3,
	'4':    input.Digit4,
	'5':    input.Digit5,
	'6':    input.Digit6,
	'7':    input.Digit7,
	'8':    input.Digit8,
	'9':    input.Digit9,
	':':    shift(input.Semicolon),
	';':    input.Semicolon,
	'<':    shift(input.Comma),
	'=':    input.Equal,
	'>':    shift(input.Period),
	'?':    shift(input.Slash),
	'@':    shift(input.Digit2),
	'A':    shift(input.KeyA),
	'B':    shift(input.KeyB),
	'C':    shift(input.KeyC),
	'D':    shift(input.KeyD),
	'E':    shift(input.KeyE),
	'F':    shift(input.KeyF),
	'G':    shift(input.KeyG),
	'H':    shift(input.KeyH),
	'I':    shift(input.KeyI),
	'J':    shift(input.KeyJ),
	'K':    shift(input.KeyK),
	'L':    shift(input.KeyL),
	'M':    shift(input.KeyM),
	'N':    shift(input.KeyN),
	'O':    shift(input.KeyO),
	'P':    shift(input.KeyP),
	'Q':    shift(input.KeyQ),
	'R':    shift(input.KeyR),
	'S':    shift(input.KeyS),
	'T':    shift(input.KeyT),
	'U':    shift(input.KeyU),
	'V':    shift(input.KeyV),
	'W':    shift(input.KeyW),
	'X':    shift(input.KeyX),
	'Y':    shift(input.KeyY),
	'Z':    shift(input.KeyZ),
	'[':    input.BracketLeft,
	'\'':   input.Quote,
	'\\':   input.Backslash,
	'\b':   input.Backspace,
	'\n':   input.Enter,
	'\r':   input.Enter,
	'\t':   input.Tab,
	'\x1b': input.Escape,
	']':    input.BracketRight,
	'^':    shift(input.Digit6),
	'_':    shift(input.Minus),
	'`':    input.Backquote,
	'a':    input.KeyA,
	'b':    input.KeyB,
	'c':    input.KeyC,
	'd':    input.KeyD,
	'e':    input.KeyE,
	'f':    input.KeyF,
	'g':    input.KeyG,
	'h':    input.KeyH,
	'i':    input.KeyI,
	'j':    input.KeyJ,
	'k':    input.KeyK,
	'l':    input.KeyL,
	'm':    input.KeyM,
	'n':    input.KeyN,
	'o':    input.KeyO,
	'p':    input.KeyP,
	'q':    input.KeyQ,
	'r':    input.KeyR,
	's':    input.KeyS,
	't':    input.KeyT,
	'u':    input.KeyU,
	'v':    input.KeyV,
	'w':    input.KeyW,
	'x':    input.KeyX,
	'y':    input.KeyY,
	'z':    input.KeyZ,
	'{':    shift(input.BracketLeft),
	'|':    shift(input.Backslash),
	'}':    shift(input.BracketRight),
	'~':    shift(input.Backquote),
	'←':    input.ArrowLeft,
	'↑':    input.ArrowUp,
	'→':    input.ArrowRight,
	'↓':    input.ArrowDown,
}
