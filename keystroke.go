package main

import (
	"strconv"
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
	events    []KeyStrokeEvent
	startTime time.Time
}

// NewKeyStrokeEvents creates a new KeyStrokeEvents struct.
func NewKeyStrokeEvents() KeyStrokeEvents {
	return KeyStrokeEvents{
		events:    make([]KeyStrokeEvent, 0),
		startTime: time.Now(),
	}
}

// Push adds a new key press event to the collection.
func (k KeyStrokeEvents) Push(display string) {
	event := KeyStrokeEvent{Display: strconv.Quote(display), WhenMS: time.Now().Sub(k.startTime).Milliseconds()}
	k.events = append(k.events, event)
}

// Slice returns the underlying slice of key press events.
// NOTE: This is a reference.
func (k KeyStrokeEvents) Slice() []KeyStrokeEvent {
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
	KeyStrokeEvents KeyStrokeEvents
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
	KeyStrokeEvents KeyStrokeEvents
}

// Press is a wrapper around the rod.Keyboard#Press method.
func (k *Keyboard) Press(key input.Key) {
	k.KeyStrokeEvents.Push(string(inverseKeymap[key]))
	k.Keyboard.Press(key)
}

// Type is a wrapper around the rod.Keyboard#Type method.
func (k *Keyboard) Type(key input.Key) {
	k.KeyStrokeEvents.Push(string(inverseKeymap[key]))
	k.Keyboard.Type(key)
}

// Input is a wrapper around the rod.Keyboard#Input method.
func (k *Keyboard) Input(text string) {
	k.KeyStrokeEvents.Push(text)
	k.textAreaElem.Input(text)
}
