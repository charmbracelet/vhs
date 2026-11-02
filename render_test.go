package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/vt"
)

func newTestRenderer(tb testing.TB, opts Options) *Renderer {
	tb.Helper()
	r, err := NewRenderer(&opts)
	if err != nil {
		tb.Skipf("no fonts: %v", err)
	}
	return r
}

// newScreen returns a terminal without a shell, whose output the test
// writes, and a renderer sized to it.
func newScreen(tb testing.TB, cols, rows int) (*Terminal, *Renderer) {
	tb.Helper()
	opts := DefaultVHSOptions()
	r := newTestRenderer(tb, opts)
	term := &Terminal{emu: vt.NewSafeEmulator(cols, rows), lastInput: time.Now()}
	term.SetTheme(opts.Theme)
	w, h := r.CellSize()
	r.SetSize(cols*w, rows*h, cols, rows)
	return term, r
}

func write(tb testing.TB, term *Terminal, s string) {
	tb.Helper()
	if _, err := term.emu.Write([]byte(s)); err != nil {
		tb.Fatal(err)
	}
}

// render renders a frame even if nothing changed since the last one.
func render(tb testing.TB, term *Terminal, r *Renderer, now time.Time) (text, cursor *image.RGBA) {
	tb.Helper()
	r.drawn = false
	text, cursor = r.Frame(term, now)
	if text == nil || cursor == nil {
		tb.Fatal("expected a frame")
	}
	return text, cursor
}

// cellCenter returns the color at the center of a cell.
func cellCenter(img *image.RGBA, r *Renderer, col, row int) color.RGBA {
	w, h := r.CellSize()
	return img.RGBAAt(col*w+w/2, row*h+h/2)
}

func TestRendererCellMetrics(t *testing.T) {
	opts := DefaultVHSOptions()
	base := newTestRenderer(t, opts)
	if base.cellW <= 0 || base.cellH <= 0 {
		t.Fatalf("invalid cell size %dx%d", base.cellW, base.cellH)
	}

	opts.LetterSpacing = 3
	if r := newTestRenderer(t, opts); r.cellW != base.cellW+2 {
		t.Errorf("letter spacing 3: expected cell width %d, got %d", base.cellW+2, r.cellW)
	}

	opts = DefaultVHSOptions()
	opts.LineHeight = 2
	if r := newTestRenderer(t, opts); r.cellH != 2*base.cellH {
		t.Errorf("line height 2: expected cell height %d, got %d", 2*base.cellH, r.cellH)
	}
}

func TestRendererMonospaceFallback(t *testing.T) {
	// A family that isn't installed must fall back to a monospace font, not
	// to the symbol fonts that come next.
	opts := DefaultVHSOptions()
	opts.FontFamily = "Font That Does Not Exist"
	r := newTestRenderer(t, opts)
	wide, _ := r.fonts.face('W', styleRegular).face.GlyphAdvance('W')
	narrow, _ := r.fonts.face('i', styleRegular).face.GlyphAdvance('i')
	if wide != narrow {
		t.Fatalf("expected a monospace font, 'W' is %v wide and 'i' %v", wide, narrow)
	}
}

func TestFontFamilies(t *testing.T) {
	tests := []struct {
		in             string
		want, fallback string
	}{
		{"DejaVu Sans Mono", "DejaVu Sans Mono,monospace", "Apple Symbols"},
		{"Hack, 'Fira Code' ,Apple Symbols", "Hack,Fira Code,monospace", "Apple Symbols,Apple Symbols"},
		{"Hack,monospace,Noto Sans", "Hack,monospace", "Noto Sans,Apple Symbols"},
		{"Hack,Monospace", "Hack,Monospace", "Apple Symbols"},
	}
	for _, tt := range tests {
		families, fallback := fontFamilies(tt.in)
		if got := strings.Join(families, ","); got != tt.want {
			t.Errorf("%q: expected families %q, got %q", tt.in, tt.want, got)
		}
		if got := strings.Join(fallback, ","); got != tt.fallback {
			t.Errorf("%q: expected fallback %q, got %q", tt.in, tt.fallback, got)
		}
	}
}

func TestFrameColors(t *testing.T) {
	term, r := newScreen(t, 10, 2)
	palette := DefaultTheme.palette()
	fg, bg, _, _ := DefaultTheme.colors()

	write(t, term, "\x1b[41m \x1b[0m")          // 0: red background
	write(t, term, "\x1b[48;5;196m \x1b[0m")    // 1: 256 color background
	write(t, term, "\x1b[48;2;1;2;3m \x1b[0m")  // 2: true color background
	write(t, term, "\x1b[7m \x1b[0m")           // 3: reverse
	write(t, term, "\x1b[31m\u2588\x1b[0m")     // 4: red full block
	write(t, term, "\x1b[1;31m\u2588\x1b[0m")   // 5: bold red is bright red
	write(t, term, "\x1b[2;31m\u2588\x1b[0m")   // 6: faint red
	write(t, term, "\x1b[38;5;3m\u2588\x1b[0m") // 7: indexed theme color
	write(t, term, "\x1b[8;31m\u2588\x1b[0m")   // 8: concealed
	text, _ := render(t, term, r, time.Now())

	tests := []struct {
		col  int
		want color.Color
	}{
		{0, palette[1]},
		{1, ansi.IndexedColor(196)},
		{2, color.RGBA{1, 2, 3, 0xff}},
		{3, fg},
		{4, palette[1]},
		{5, palette[9]},
		{6, blend(toRGBA(palette[1]), toRGBA(bg))},
		{7, palette[3]},
		{8, bg},
		{9, bg},
	}
	for _, tt := range tests {
		if got, want := cellCenter(text, r, tt.col, 0), toRGBA(tt.want); got != want {
			t.Errorf("cell %d: expected %v, got %v", tt.col, want, got)
		}
	}
}

func TestFrameBoxDrawingFillsCell(t *testing.T) {
	term, r := newScreen(t, 2, 1)
	write(t, term, "\u2500\u2502") // horizontal, vertical
	text, _ := render(t, term, r, time.Now())
	w, h := r.CellSize()
	fg, _, _, _ := DefaultTheme.colors()
	want := toRGBA(fg)

	// The lines run from edge to edge, so they connect to their neighbors.
	for _, p := range []image.Point{{0, h / 2}, {w - 1, h / 2}, {w + w/2, 0}, {w + w/2, h - 1}} {
		if got := text.RGBAAt(p.X, p.Y); got != want {
			t.Errorf("pixel %v: expected %v, got %v", p, want, got)
		}
	}
}

func TestFrameCursor(t *testing.T) {
	term, r := newScreen(t, 4, 2)
	_, _, cursorColor, _ := DefaultTheme.colors()
	start := term.LastInput()
	transparent := color.RGBA{}

	_, cursor := render(t, term, r, start.Add(100*time.Millisecond))
	if got := cellCenter(cursor, r, 0, 0); got != toRGBA(cursorColor) {
		t.Errorf("expected the cursor at 0,0, got %v", got)
	}
	if got := cellCenter(cursor, r, 1, 0); got != transparent {
		t.Errorf("expected nothing but the cursor on the cursor layer, got %v", got)
	}

	// Blinking: off in the second interval after input.
	_, cursor = render(t, term, r, start.Add(cursorBlinkInterval+100*time.Millisecond))
	if got := cellCenter(cursor, r, 0, 0); got != transparent {
		t.Errorf("expected the cursor to blink off, got %v", got)
	}

	// A steady cursor doesn't blink.
	write(t, term, "\x1b[2 q")
	_, cursor = render(t, term, r, start.Add(cursorBlinkInterval+100*time.Millisecond))
	if got := cellCenter(cursor, r, 0, 0); got != toRGBA(cursorColor) {
		t.Errorf("expected a steady cursor, got %v", got)
	}

	write(t, term, "\x1b[?25l")
	_, cursor = render(t, term, r, start)
	if got := cellCenter(cursor, r, 0, 0); got != transparent {
		t.Errorf("expected a hidden cursor, got %v", got)
	}
}

func TestFrameUnchanged(t *testing.T) {
	term, r := newScreen(t, 4, 2)
	r.blink = false
	text, cursor := r.Frame(term, time.Now())
	if text == nil || cursor == nil || text.Rect != image.Rect(0, 0, r.width, r.height) {
		t.Fatalf("expected a %dx%d frame", r.width, r.height)
	}
	if text, cursor := r.Frame(term, time.Now()); text != nil || cursor != nil {
		t.Fatal("expected no new layers when nothing changed")
	}
	// Output changes both layers: the text, and the cursor moving.
	write(t, term, "x")
	if text, cursor := r.Frame(term, time.Now()); text == nil || cursor == nil {
		t.Fatal("expected new layers after output")
	}
	// Moving the cursor back over the same text changes the cursor only.
	write(t, term, "\b")
	if text, cursor := r.Frame(term, time.Now()); text != nil || cursor == nil {
		t.Fatal("expected only a new cursor layer after moving the cursor")
	}
	// Color changes change the text only.
	write(t, term, "\x1b[31mx\x1b[0m\b")
	if text, cursor := r.Frame(term, time.Now()); text == nil || cursor != nil {
		t.Fatal("expected only a new text layer after recoloring")
	}
}

func TestTerminalScroll(t *testing.T) {
	term, r := newScreen(t, 10, 10)
	for i := 1; i <= 30; i++ {
		write(t, term, fmt.Sprintf("line %d\r\n", i))
	}
	first := func() string { return term.Lines()[0] }
	if got := first(); got != "line 22" {
		t.Fatalf("expected the bottom of the output, got %q", got)
	}

	term.Scroll(-5)
	if got := first(); got != "line 17" {
		t.Fatalf("expected to scroll up 5 lines, got %q", got)
	}
	// New output doesn't move the viewport.
	write(t, term, "line 31\r\n")
	if got := first(); got != "line 17" {
		t.Fatalf("expected the viewport to stay put, got %q", got)
	}
	// The cursor isn't drawn while scrolled up.
	if _, cursor := render(t, term, r, term.LastInput()); !isTransparent(cursor) {
		t.Error("expected no cursor while scrolled up")
	}

	term.Scroll(-100)
	if got := first(); got != "line 1" {
		t.Fatalf("expected to stop at the top, got %q", got)
	}
	term.Scroll(100)
	if got := first(); got != "line 23" {
		t.Fatalf("expected to stop at the bottom, got %q", got)
	}

	term.Scroll(-3)
	term.userInput()
	if term.ViewOffset() != 0 {
		t.Fatal("expected input to scroll back to the bottom")
	}
}

func isTransparent(img *image.RGBA) bool {
	for _, b := range img.Pix {
		if b != 0 {
			return false
		}
	}
	return true
}

func BenchmarkFrame(b *testing.B) {
	term, r := newShellTerminal(b)
	term.SendText("ls -la /usr/bin | head -15; echo DONE\r")
	waitForLine(b, term, "DONE")
	b.ReportAllocs()
	var i int64
	for b.Loop() {
		i++
		r.drawn = false
		text, cur := r.Frame(term, time.Unix(i, 0))
		r.ReleaseText(text)
		r.ReleaseCursor(cur)
	}
}

func BenchmarkEncode(b *testing.B) {
	term, r := newShellTerminal(b)
	term.SendText("ls -la /usr/bin | head -15; echo DONE\r")
	waitForLine(b, term, "DONE")
	text, _ := r.Frame(term, time.Now())
	b.ReportAllocs()
	for b.Loop() {
		var buf bytes.Buffer
		if err := pngEncoder.Encode(&buf, text); err != nil {
			b.Fatal(err)
		}
	}
}

func newShellTerminal(tb testing.TB) (*Terminal, *Renderer) {
	tb.Helper()
	opts := DefaultVHSOptions()
	r := newTestRenderer(tb, opts)
	term, err := StartTerminal(context.Background(), Shell{Command: []string{"sh"}, Env: []string{"PS1=> "}}, 60, 20)
	if err != nil {
		tb.Fatal(err)
	}
	tb.Cleanup(func() { _ = term.Close() })
	term.SetTheme(opts.Theme)
	w, h := r.CellSize()
	r.SetSize(60*w, 20*h, 60, 20)
	return term, r
}

// waitForLine waits until a line of the screen is want.
func waitForLine(tb testing.TB, term *Terminal, want string) {
	tb.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if slices.Contains(term.Lines(), want) {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	tb.Fatalf("timed out waiting for %q, screen:\n%s", want, strings.Join(term.Lines(), "\n"))
}
