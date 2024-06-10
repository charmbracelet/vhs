package main

import (
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

// KeyStrokeEvent represents a key press event for the purposes of keystroke
// overlay.
type KeyStrokeEvent struct {
	// Display generally includes the current key stroke sequence.
	Display string
	// WhenMS is the time in milliseconds when the key was pressed starting
	// from the beginning of the recording.
	WhenMS int64
}

// KeyStrokeEvents is a collection of key press events that you can push to.
type KeyStrokeEvents struct {
	enabled        bool
	display        string
	events         []KeyStrokeEvent
	once           sync.Once
	startTime      time.Time
	maxDisplaySize int
}

// NewKeyStrokeEvents creates a new KeyStrokeEvents struct.
func NewKeyStrokeEvents() *KeyStrokeEvents {
	return &KeyStrokeEvents{
		display: "",
		events:  make([]KeyStrokeEvent, 0),
		// NOTE: This is actually setting the startTime too early. It
		// takes a while (in computer time) to get to the point where
		// we start recording. Therefore, we actually set this another
		// time on the first push. Without this, the final overlay
		// would be slightly desynced by a 20-40 ms, which is
		// noticeable to the human eye.
		startTime:      time.Now(),
		maxDisplaySize: 20,
	}
}

// keypressSymbolOverrides maps certain input keys to their corresponding
// keypress string or symbol. These override the default rune for the
// corresponding input key to improve the visuals or readability of the keypress
// overlay. A good example of this improvement can be seen in things like Enter
// (newline). The description string and symbol are embedded into an inner map,
// which can be indexed into based on whether special symbols are requested or
// not.
// TODO: I think we can ignore the non-special symbol case and just use the
// symbols since we can rely on font fallback behavior.
var keypressSymbolOverrides = map[input.Key]map[bool]string{
	input.Backspace: {
		true:  "\\\\\\\\b",
		false: "⌫",
	},
	input.Delete: {
		true:  "\\\\\\\\d",
		false: "␡",
	},
	input.ControlLeft: {
		true:  "<CTRL>+",
		false: "C-",
	},
	input.ControlRight: {
		true:  "<CTRL>+",
		false: "C-",
	},
	input.AltLeft: {
		true:  "<ALT>+",
		false: "⎇-",
	},
	input.AltRight: {
		true:  "<ALT>+",
		false: "⎇-",
	},
	input.ArrowDown: {
		true:  "<DOWN>",
		false: "↓",
	},
	input.PageDown: {
		true:  "<PAGEDOWN>",
		false: "⤓",
	},
	input.ArrowUp: {
		true:  "<UP>",
		false: "↑",
	},
	input.PageUp: {
		true:  "<PAGEUP>",
		false: "⤒",
	},
	input.ArrowLeft: {
		true:  "<LEFT>",
		false: "←",
	},
	input.ArrowRight: {
		true:  "<RIGHT>",
		false: "→",
	},
	input.Space: {
		true:  "<SPACE>",
		false: "␣",
	},
	input.Enter: {
		true:  "<ENTER>",
		false: "⏎",
	},
	input.Escape: {
		true:  "<ESCAPE>",
		false: "⎋",
	},
	input.Tab: {
		true:  "<TAB>",
		false: "⇥",
	},
}

func keyToDisplay(key input.Key) string {
	if override, ok := keypressSymbolOverrides[key]; ok {
		if symbol, ok := override[false]; ok {
			return symbol
		}
	}
	return string(inverseKeymap[key])
}

// Enable enables key press event recording.
func (k *KeyStrokeEvents) Enable() {
	k.enabled = true
}

// Disable disables key press event recording.
func (k *KeyStrokeEvents) Disable() {
	k.enabled = false
}

// Push adds a new key press event to the collection.
func (k *KeyStrokeEvents) Push(display string) {
	k.once.Do(func() {
		k.startTime = time.Now()
	})

	// If we're not enabled, we don't want to do anything.
	// But note that we still want to update the start time -- this is because
	// we need to know the global start time if we want to render any subsequent
	// events correctly, and the keystroke overlay may be re-enabled later in
	// the recording.
	if !k.enabled {
		return
	}

	k.display += display
	// Keep k.display @ 20 max.
	// Anymore than that is probably overkill, and we don't want to run into
	// issues where the overlay text is longer than the video width itself.
	if len(k.display) > k.maxDisplaySize {
		// We need to be cognizant of unicode -- we can't just slice off a byte,
		// we have to slice off a _rune_. The conversion back-and-forth may be a
		// bit inefficient, but k.display will always be tiny thanks to
		// k.maxDisplaySize.
		k.display = string([]rune(k.display)[1:])
	}
	event := KeyStrokeEvent{Display: k.display, WhenMS: time.Now().Sub(k.startTime).Milliseconds()}
	k.events = append(k.events, event)
}

// Slice returns the underlying slice of key press events.
// NOTE: This is a reference.
func (k *KeyStrokeEvents) Slice() []KeyStrokeEvent {
	return k.events
}

// Page is a wrapper around the rod.Page object.
// It's primary purpose is to decorate the rod.Page struct such that we can
// record keypress events during the recording for keypress overlays. We prefer
// decorating so that we that minimize the possibility of future bugs around
// forgetting to log key presses, since all input is done through rod.Page (and
// technically rod.Page.MustElement() + rod.Page.Keyboard).
type Page struct {
	*rod.Page
	Keyboard        Keyboard
	KeyStrokeEvents *KeyStrokeEvents
}

// NewPage creates a new wrapper Page object.
func NewPage(page *rod.Page) *Page {
	keyStrokeEvents := NewKeyStrokeEvents()
	return &Page{Page: page, KeyStrokeEvents: keyStrokeEvents, Keyboard: Keyboard{page.Keyboard, page.MustElement("textarea"), keyStrokeEvents}}
}

// MustSetViewport is a wrapper around the rod.Page#MustSetViewport method.
func (p *Page) MustSetViewport(width, height int, deviceScaleFactor float64, mobile bool) *Page {
	p.Page.MustSetViewport(width, height, deviceScaleFactor, mobile)
	return p
}

// MustWait is a wrapper around the rod.Page#MustWait method.
func (p *Page) MustWait(js string) *Page {
	p.Page.MustWait(js)
	return p
}

// Keyboard is a wrapper around the rod.Keyboard object.
type Keyboard struct {
	*rod.Keyboard
	textAreaElem    *rod.Element
	KeyStrokeEvents *KeyStrokeEvents
}

// Press is a wrapper around the rod.Keyboard#Press method.
func (k *Keyboard) Press(key input.Key) {
	k.KeyStrokeEvents.Push(keyToDisplay(key))
	k.Keyboard.Press(key)
}

// Type is a wrapper around the rod.Keyboard#Type method.
func (k *Keyboard) Type(key input.Key) {
	k.KeyStrokeEvents.Push(keyToDisplay(key))
	k.Keyboard.Type(key)
}

// Input is a wrapper around the rod.Page#MustElement("textarea")#Input method.
func (k *Keyboard) Input(text string) {
	k.KeyStrokeEvents.Push(text)
	k.textAreaElem.Input(text)
}
