package main

import (
	"image"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

// emojiFonts returns the color emoji fonts on this machine: Apple's (sbix)
// and Noto's (CBDT), from the system or the Go module cache.
func emojiFonts(t *testing.T) []string {
	t.Helper()
	candidates := []string{
		"/System/Library/Fonts/Apple Color Emoji.ttc",
		"/usr/share/fonts/truetype/noto/NotoColorEmoji.ttf",
		"/usr/share/fonts/noto/NotoColorEmoji.ttf",
		"/usr/share/fonts/google-noto-emoji/NotoColorEmoji.ttf",
	}
	if out, err := exec.Command("go", "env", "GOMODCACHE").Output(); err == nil {
		matches, _ := filepath.Glob(filepath.Join(strings.TrimSpace(string(out)),
			"github.com/go-text/typesetting-utils@*/opentype/bitmap/NotoColorEmoji.ttf"))
		candidates = append(candidates, matches...)
	}
	var found []string
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			found = append(found, p)
		}
	}
	if len(found) == 0 {
		t.Skip("no color emoji font found")
	}
	return found
}

func TestColorGlyphs(t *testing.T) {
	for _, path := range emojiFonts(t) {
		t.Run(filepath.Base(path), func(t *testing.T) {
			f, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			coll, err := opentype.ParseCollectionReaderAt(f)
			if err != nil {
				t.Fatal(err)
			}
			ft, err := coll.Font(0)
			if err != nil {
				t.Fatal(err)
			}
			colors, err := readColorGlyphs(f, 0)
			if err != nil || colors == nil {
				t.Fatalf("expected color glyphs, got %v", err)
			}
			var buf sfnt.Buffer
			gid, err := ft.GlyphIndex(&buf, '😀')
			if err != nil || gid == 0 {
				t.Fatalf("no glyph for an emoji: %v", err)
			}
			for _, ppem := range []int{8, 22, 64, 1000} {
				img, got, err := colors.glyph(gid, ppem)
				if err != nil || img.Bounds().Dx() < 8 || got <= 0 {
					t.Fatalf("ppem %d: expected an image, got %v (%v)", ppem, img, err)
				}
			}
			if _, _, err := colors.glyph(0, 22); err == nil {
				t.Fatal("expected no image for the missing glyph")
			}
		})
	}
}

func TestReadColorGlyphsOfTextFont(t *testing.T) {
	opts := DefaultVHSOptions()
	r := newTestRenderer(t, opts)
	f := r.fonts.face('W', styleRegular)
	if f.color != nil {
		t.Fatal("expected a text font to have no color glyphs")
	}
}

func TestEmojiAreDrawnInColor(t *testing.T) {
	emojiFonts(t)
	term, r := newScreen(t, 6, 1)
	write(t, term, "😀")
	text, _ := render(t, term, r, term.LastInput())
	w, h := r.CellSize()
	colors := map[[3]uint8]bool{}
	for y := range h {
		for x := range 2 * w {
			c := text.RGBAAt(x, y)
			colors[[3]uint8{c.R, c.G, c.B}] = true
		}
	}
	if len(colors) < 20 {
		t.Fatalf("expected a color image, got %d colors", len(colors))
	}
}

func TestFontFallbacks(t *testing.T) {
	r := newTestRenderer(t, DefaultVHSOptions())
	fs := r.fonts
	if f := fs.face('日', styleRegular); f == nil {
		t.Skip("no CJK font installed")
	}
	if f := fs.face('😀', styleRegular); f != nil && f.color == nil {
		t.Error("expected emoji to come from a color font")
	}
	// Private use characters are only drawn from listed fonts, since they
	// mean something different in each font.
	if f := fs.face('\U000f0001', styleRegular); f != nil {
		t.Error("expected no face for an unlisted private use character")
	}
	// Bold falls back to the regular face, made bold.
	if f := fs.face('W', styleBold); f == nil {
		t.Error("expected a bold face")
	}
}

func TestSlant(t *testing.T) {
	mask := image.NewAlpha(image.Rect(0, 0, 4, 10))
	for i := range mask.Pix {
		mask.Pix[i] = 0xff
	}
	g := slant(&glyphMask{mask: mask, rect: image.Rect(1, 5, 5, 15)}, 15)
	if g.rect.Dx() <= 4 || g.rect.Dy() != 10 {
		t.Fatalf("expected a wider glyph of the same height, got %v", g.rect)
	}
	// The top row leans right, the bottom one stays on the baseline.
	if top, bottom := g.mask.Pix[g.mask.Stride-1], g.mask.Pix[(g.rect.Dy()-1)*g.mask.Stride]; top != 0xff || bottom != 0xff {
		t.Fatalf("expected the top row moved right and the bottom one in place")
	}
}

// coverage returns how much of the rectangle of the mask is covered.
func coverage(m *image.Alpha, r image.Rectangle) float64 {
	var sum float64
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			sum += float64(m.AlphaAt(x, y).A) / 0xff
		}
	}
	return sum / float64(r.Dx()*r.Dy())
}

func TestCustomGlyphs(t *testing.T) {
	r := newTestRenderer(t, DefaultVHSOptions())
	w, h := r.CellSize()
	glyph := func(ch rune) *image.Alpha {
		t.Helper()
		if !isCustomGlyph(ch) {
			t.Fatalf("%q is not a custom glyph", ch)
		}
		return r.customGlyph(ch)
	}
	row := func(m *image.Alpha, y int) float64 { return coverage(m, image.Rect(0, y, w, y+1)) }
	col := func(m *image.Alpha, x int) float64 { return coverage(m, image.Rect(x, 0, x+1, h)) }

	// Lines run edge to edge through the middle of the cell, 1 or 3 pixels
	// wide and crisp.
	if m := glyph('─'); row(m, h/2) != 1 || coverage(m, m.Rect) != float64(w)/float64(w*h) {
		t.Errorf("expected a 1 pixel line across the middle, got %.2f", coverage(m, m.Rect))
	}
	if m := glyph('━'); coverage(m, m.Rect) != 3*float64(w)/float64(w*h) {
		t.Errorf("expected a 3 pixel heavy line, got %.2f", coverage(m, m.Rect))
	}
	if m := glyph('│'); col(m, w/2) != 1 {
		t.Error("expected a line down the middle")
	}
	// A corner reaches the bottom and the right, not the top or left.
	if m := glyph('┌'); m.AlphaAt(w/2, h-1).A == 0 || m.AlphaAt(w-1, h/2).A == 0 || row(m, 0) != 0 || col(m, 0) != 0 {
		t.Error("expected a corner going down and right")
	}
	// A rounded corner has nothing where a square one has its corner.
	if m := glyph('╭'); m.AlphaAt(w/2, h/2).A != 0 || m.AlphaAt(w/2, h-1).A == 0 || m.AlphaAt(w-1, h/2).A == 0 {
		t.Error("expected a rounded corner")
	}
	// Double lines are two lines.
	if m := glyph('═'); coverage(m, m.Rect) != 2*float64(w)/float64(w*h) {
		t.Errorf("expected two lines, got %.2f", coverage(m, m.Rect))
	}
	// Dashes leave gaps.
	if m := glyph('┄'); row(m, h/2) >= 0.9 || row(m, h/2) <= 0.3 {
		t.Errorf("expected a dashed line, got %.2f", row(m, h/2))
	}
	// Shades are patterns of their density.
	for ch, want := range map[rune]float64{'░': 0.125, '▒': 0.25, '▓': 0.75} {
		if got := coverage(glyph(ch), image.Rect(0, 0, w/4*4, h/4*4)); got != want {
			t.Errorf("%q: expected %.3f coverage, got %.3f", ch, want, got)
		}
	}
	// Blocks fill whole pixels.
	if m := glyph('▄'); coverage(m, image.Rect(0, h/2+1, w, h)) != 1 || coverage(m, image.Rect(0, 0, w, h/2-1)) != 0 {
		t.Error("expected the lower half filled")
	}
	// Powerline arrows fill the left of the cell and point right.
	if m := glyph('\ue0b0'); col(m, 0) < 0.95 || m.AlphaAt(w-1, 0).A != 0 {
		t.Error("expected a right-pointing arrow")
	}
	if isCustomGlyph('a') || isCustomGlyph('日') {
		t.Error("expected text not to be a custom glyph")
	}
}
