// Package vhs keys.go maps tape keys and characters to the key presses sent
// to the terminal.
//
// Type Hello, world!
//
// Each character is sent as typed. Characters that are keys of their own,
// like a newline or an arrow, are sent as that key so the terminal encodes
// them according to its current modes.
package main

import (
	"strings"
	"unicode"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/vt"
)

// Named keys used by tape commands.
var (
	keyBackspace = uv.KeyPressEvent{Code: vt.KeyBackspace}
	keyDelete    = uv.KeyPressEvent{Code: vt.KeyDelete}
	keyInsert    = uv.KeyPressEvent{Code: vt.KeyInsert}
	keyDown      = uv.KeyPressEvent{Code: vt.KeyDown}
	keyEnter     = uv.KeyPressEvent{Code: vt.KeyEnter}
	keyLeft      = uv.KeyPressEvent{Code: vt.KeyLeft}
	keyRight     = uv.KeyPressEvent{Code: vt.KeyRight}
	keySpace     = uv.KeyPressEvent{Code: vt.KeySpace, Text: " "}
	keyUp        = uv.KeyPressEvent{Code: vt.KeyUp}
	keyTab       = uv.KeyPressEvent{Code: vt.KeyTab}
	keyEscape    = uv.KeyPressEvent{Code: vt.KeyEscape}
	keyPgUp      = uv.KeyPressEvent{Code: vt.KeyPgUp}
	keyPgDown    = uv.KeyPressEvent{Code: vt.KeyPgDown}
)

// runeKeys are the characters that are sent as a key of their own.
var runeKeys = map[rune]uv.KeyPressEvent{
	'\b':   keyBackspace,
	'\n':   keyEnter,
	'\r':   keyEnter,
	'\t':   keyTab,
	'\x1b': keyEscape,
	'←':    keyLeft,
	'↑':    keyUp,
	'→':    keyRight,
	'↓':    keyDown,
}

// shiftedRunes maps the keys of a US keyboard to what they type with Shift
// held.
var shiftedRunes = map[rune]rune{
	'1': '!', '2': '@', '3': '#', '4': '$', '5': '%',
	'6': '^', '7': '&', '8': '*', '9': '(', '0': ')',
	'-': '_', '=': '+', '[': '{', ']': '}', '\\': '|',
	';': ':', '\'': '"', ',': '<', '.': '>', '/': '?', '`': '~',
}

// shiftRune returns what r types with Shift held.
func shiftRune(r rune) rune {
	if s, ok := shiftedRunes[r]; ok {
		return s
	}
	return unicode.ToUpper(r)
}

// runeKey returns the key press that types r with the given modifiers.
func runeKey(r rune, mod uv.KeyMod) uv.KeyPressEvent {
	if k, ok := runeKeys[r]; ok {
		k.Mod = mod
		return k
	}
	if mod&uv.ModShift != 0 {
		r = shiftRune(r)
	}
	return uv.KeyPressEvent{Code: r, Text: string(r), Mod: mod}
}

// ctrlKey returns the key press for the arguments of a Ctrl command, e.g.
// "Shift Alt C" or "Left".
func ctrlKey(args string) (uv.KeyPressEvent, bool) {
	k := uv.KeyPressEvent{Mod: uv.ModCtrl}
	fields := strings.Fields(args)
	for i, arg := range fields {
		if i < len(fields)-1 {
			switch arg {
			case "Shift":
				k.Mod |= uv.ModShift
			case "Alt":
				k.Mod |= uv.ModAlt
			}
			continue
		}

		switch arg {
		case "Shift", "Alt":
			// A modifier on its own sends nothing.
			return k, false
		case "Enter":
			k.Code = vt.KeyEnter
		case "Space":
			k.Code = vt.KeySpace
		case "Backspace":
			k.Code = vt.KeyBackspace
		case "Left":
			k.Code = vt.KeyLeft
		case "Right":
			k.Code = vt.KeyRight
		case "Up":
			k.Code = vt.KeyUp
		case "Down":
			k.Code = vt.KeyDown
		default:
			r := []rune(arg)[0]
			switch {
			case unicode.IsUpper(r):
				r = unicode.ToLower(r)
				k.Mod |= uv.ModShift
			case r == '@':
				r = vt.KeySpace // Ctrl+@ is NUL, like Ctrl+Space.
			case r == '-':
				r = '_' // Ctrl+- is the same control character as Ctrl+_.
			}
			k.Code = r
		}
	}
	return k, k.Code != 0
}
