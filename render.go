package main

import (
	"hash/maphash"
	"image"
	"image/color"
	"image/draw"
	"math"
	"sync"
	"time"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/vt"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// cursorBlinkInterval is how long the cursor stays on, then off, while
// blinking, matching xterm.js.
const cursorBlinkInterval = 600 * time.Millisecond

// Renderer rasterizes the terminal into frames: an opaque text layer and a
// transparent cursor layer, the same two layers xterm.js used to draw.
type Renderer struct {
	fonts *fontSet

	// Cell size in pixels, where text sits in a cell, and the thickness of
	// drawn lines.
	cellW, cellH      int
	glyphX, baseline  int
	strikeY, lineW    int
	width, height     int
	blink             bool
	cursor, accent    color.Color
	glyphs            map[glyphKey]*glyphMask
	seed              maphash.Seed
	lastHash          uint64
	hasLast           bool
	textPool, curPool sync.Pool
}

type glyphKey struct {
	r    rune
	face font.Face
}

// glyphMask is a rasterized glyph and where it goes relative to the top-left
// of its cell.
type glyphMask struct {
	mask *image.Alpha
	rect image.Rectangle
}

// NewRenderer loads the fonts and measures the cells for the given options.
func NewRenderer(opts *Options) (*Renderer, error) {
	fonts, err := newFontSet(opts.FontFamily, opts.FontSize)
	if err != nil {
		return nil, err
	}
	face := fonts.face('W', styleRegular).face
	adv, _ := face.GlyphAdvance('W')
	m := face.Metrics()

	// Cell metrics follow xterm.js: the character width is floored and the
	// letter spacing rounded, and line height scales the character height
	// with the text centered vertically.
	spacing := int(math.Round(opts.LetterSpacing))
	charH := (m.Ascent + m.Descent).Ceil()
	r := &Renderer{
		fonts:  fonts,
		cellW:  max(1, adv.Floor()+spacing),
		cellH:  max(1, int(math.Floor(float64(charH)*opts.LineHeight))),
		glyphX: max(0, spacing/2),                                  //nolint:mnd
		lineW:  max(1, int(math.Round(float64(opts.FontSize)/15))), //nolint:mnd
		blink:  opts.CursorBlink,
		glyphs: map[glyphKey]*glyphMask{},
		seed:   maphash.MakeSeed(),
	}
	top := 0
	if r.cellH != charH {
		top = int(math.Round(float64(r.cellH-charH) / 2)) //nolint:mnd
	}
	r.baseline = top + m.Ascent.Round()
	xHeight := m.XHeight.Round()
	if xHeight <= 0 {
		xHeight = m.Ascent.Round() * 2 / 3 //nolint:mnd
	}
	r.strikeY = r.baseline - xHeight/2 //nolint:mnd
	_, _, r.cursor, r.accent = opts.Theme.colors()
	return r, nil
}

// CellSize returns the size of a terminal cell in pixels.
func (r *Renderer) CellSize() (w, h int) { return r.cellW, r.cellH }

// SetSize sets the size of the frames in pixels.
func (r *Renderer) SetSize(w, h int) {
	r.width, r.height = w, h
}

// Release returns frame images to be reused once they have been encoded.
func (r *Renderer) Release(text, cursor *image.RGBA) {
	if text != nil {
		r.textPool.Put(text)
	}
	if cursor != nil {
		r.curPool.Put(cursor)
	}
}

func (r *Renderer) newImage(pool *sync.Pool) *image.RGBA {
	if img, ok := pool.Get().(*image.RGBA); ok && img.Rect.Dx() == r.width && img.Rect.Dy() == r.height {
		return img
	}
	return image.NewRGBA(image.Rect(0, 0, r.width, r.height))
}

// cellDraw is a cell resolved to the colors it is drawn with.
type cellDraw struct {
	x, y    int
	cell    *uv.Cell
	fg, bg  color.RGBA
	ul      color.RGBA
	style   fontStyle
	conceal bool
}

// cursorState is what the cursor layer depends on.
type cursorState struct {
	visible bool
	pos     uv.Position
	style   vt.CursorStyle
}

// Frame renders the terminal. When nothing visible changed since the last
// frame it returns changed == false and no images.
func (r *Renderer) Frame(t *Terminal, now time.Time) (text, cursor *image.RGBA, changed bool) {
	buf := t.Snapshot()
	offset := t.ViewOffset()
	defFg, defBg := t.Colors()
	fg0, bg0 := toRGBA(defFg), toRGBA(defBg)

	var palette [256]color.RGBA
	var paletteSet [256]bool
	indexed := func(i int) color.RGBA {
		if !paletteSet[i] {
			palette[i] = toRGBA(t.emu.IndexedColor(i))
			paletteSet[i] = true
		}
		return palette[i]
	}
	resolve := func(c color.Color, def color.RGBA, bold bool) color.RGBA {
		switch c := c.(type) {
		case nil:
			return def
		case ansi.BasicColor:
			if bold && c < 8 { //nolint:mnd
				c += 8
			}
			return indexed(int(c))
		case ansi.IndexedColor:
			return indexed(int(c))
		default:
			return toRGBA(c)
		}
	}

	var h maphash.Hash
	h.SetSeed(r.seed)
	writeColor := func(c color.RGBA) { _, _ = h.Write([]byte{c.R, c.G, c.B, c.A}) }
	writeInt := func(i int) { _, _ = h.Write([]byte{byte(i), byte(i >> 8), byte(i >> 16)}) } //nolint:gosec,mnd
	writeColor(bg0)

	cells := make([]cellDraw, 0, buf.Width()*buf.Height())
	for y := range buf.Height() {
		for x := 0; x < buf.Width(); {
			c := buf.CellAt(x, y)
			if c == nil || c.IsZero() {
				x++
				continue
			}
			st := c.Style
			bold := st.Attrs&uv.AttrBold != 0
			d := cellDraw{
				x: x, y: y, cell: c,
				fg:      resolve(st.Fg, fg0, bold),
				bg:      resolve(st.Bg, bg0, false),
				conceal: st.Attrs&uv.AttrConceal != 0,
			}
			if st.Attrs&uv.AttrReverse != 0 {
				d.fg, d.bg = d.bg, d.fg
			}
			if st.Attrs&uv.AttrFaint != 0 {
				d.fg = blend(d.fg, d.bg)
			}
			d.ul = d.fg
			if st.UnderlineColor != nil {
				d.ul = resolve(st.UnderlineColor, fg0, false)
			}
			if bold {
				d.style |= styleBold
			}
			if st.Attrs&uv.AttrItalic != 0 {
				d.style |= styleItalic
			}
			cells = append(cells, d)

			_, _ = h.WriteString(c.Content)
			writeInt(x)
			writeInt(c.Width)
			_, _ = h.Write([]byte{st.Attrs, st.Underline})
			writeColor(d.fg)
			writeColor(d.bg)
			writeColor(d.ul)
			x += max(c.Width, 1)
		}
		writeInt(-1)
	}

	cur := r.cursorState(t, buf, offset, now)
	if cur.visible {
		writeInt(cur.pos.X)
		writeInt(cur.pos.Y)
		writeInt(int(cur.style))
	}

	sum := h.Sum64()
	if r.hasLast && sum == r.lastHash {
		return nil, nil, false
	}
	r.lastHash, r.hasLast = sum, true

	text = r.newImage(&r.textPool)
	draw.Draw(text, text.Rect, image.NewUniform(bg0), image.Point{}, draw.Src)
	for _, d := range cells {
		if d.bg != bg0 {
			draw.Draw(text, r.cellRect(d.x, d.y, d.cell.Width), image.NewUniform(d.bg), image.Point{}, draw.Src)
		}
	}
	for _, d := range cells {
		r.drawCell(text, d)
	}

	cursor = r.newImage(&r.curPool)
	clear(cursor.Pix)
	if cur.visible {
		r.drawCursor(cursor, buf, cur)
	}
	return text, cursor, true
}

// cursorState reports where and how the cursor is drawn, if at all.
func (r *Renderer) cursorState(t *Terminal, buf uv.ScreenBuffer, offset int, now time.Time) cursorState {
	style, blink := t.emu.CursorStyle()
	st := cursorState{
		visible: !t.emu.CursorHidden() && offset == 0,
		pos:     t.emu.CursorPosition(),
		style:   style,
	}
	if st.visible && r.blink && blink {
		st.visible = (now.Sub(t.LastInput())/cursorBlinkInterval)%2 == 0 //nolint:mnd
	}
	st.pos.X = min(max(st.pos.X, 0), buf.Width()-1)
	st.pos.Y = min(max(st.pos.Y, 0), buf.Height()-1)
	return st
}

func (r *Renderer) cellRect(x, y, width int) image.Rectangle {
	px, py := x*r.cellW, y*r.cellH
	return image.Rect(px, py, px+max(width, 1)*r.cellW, py+r.cellH)
}

func (r *Renderer) drawCell(dst *image.RGBA, d cellDraw) {
	c := d.cell
	px, py := d.x*r.cellW, d.y*r.cellH
	if !d.conceal && c.Content != "" && c.Content != " " {
		r.drawText(dst, px, py, c.Content, d.style, d.fg)
	}

	width := max(c.Width, 1) * r.cellW
	if c.Style.Underline != uv.UnderlineNone {
		r.drawUnderline(dst, px, py, width, c.Style.Underline, d.ul)
	}
	if c.Style.Attrs&uv.AttrStrikethrough != 0 {
		fill(dst, image.Rect(px, py+r.strikeY, px+width, py+r.strikeY+r.lineW), d.fg)
	}
}

// drawText draws a grapheme cluster with its top-left at px, py.
func (r *Renderer) drawText(dst *image.RGBA, px, py int, s string, style fontStyle, fg color.RGBA) {
	runes := []rune(s)
	if len(runes) == 1 && r.drawCustomGlyph(dst, px, py, runes[0], fg) {
		return
	}
	src := image.NewUniform(fg)
	for _, ch := range runes {
		gf := r.fonts.face(ch, style)
		if gf == nil {
			continue
		}
		g := r.glyph(ch, gf.face)
		if g == nil {
			continue
		}
		dr := g.rect.Add(image.Pt(px+r.glyphX, py))
		draw.DrawMask(dst, dr, src, image.Point{}, g.mask, image.Point{}, draw.Over)
		if gf.fakeBold {
			draw.DrawMask(dst, dr.Add(image.Pt(1, 0)), src, image.Point{}, g.mask, image.Point{}, draw.Over)
		}
	}
}

// glyph returns the rasterized glyph, cached since the rasterizer reuses its
// mask buffer between calls.
func (r *Renderer) glyph(ch rune, face font.Face) *glyphMask {
	key := glyphKey{ch, face}
	if g, ok := r.glyphs[key]; ok {
		return g
	}
	var g *glyphMask
	dr, mask, maskp, _, ok := face.Glyph(fixed.P(0, r.baseline), ch)
	if ok && !dr.Empty() {
		a := image.NewAlpha(image.Rect(0, 0, dr.Dx(), dr.Dy()))
		draw.Draw(a, a.Rect, mask, maskp, draw.Src)
		g = &glyphMask{mask: a, rect: dr}
	}
	r.glyphs[key] = g
	return g
}

func (r *Renderer) drawUnderline(dst *image.RGBA, px, py, width int, style uv.Underline, c color.RGBA) {
	lw := r.lineW
	uy := py + min(r.cellH-lw, r.baseline+lw)
	switch style {
	case uv.UnderlineDouble:
		top := max(py+r.baseline+1, uy-lw)
		fill(dst, image.Rect(px, top, px+width, top+lw), c)
		bottom := min(py+r.cellH-lw, top+2*lw)
		fill(dst, image.Rect(px, bottom, px+width, bottom+lw), c)
	case uv.UnderlineCurly:
		amp := float64(max(lw, 1))
		for i := range width {
			y := uy + int(math.Round(amp*math.Sin(2*math.Pi*float64(i)/float64(r.cellW))))
			fill(dst, image.Rect(px+i, y, px+i+1, y+lw), c)
		}
	case uv.UnderlineDotted:
		for i := 0; i < width; i += 2 * lw {
			fill(dst, image.Rect(px+i, uy, px+min(i+lw, width), uy+lw), c)
		}
	case uv.UnderlineDashed:
		dash := max(2, r.cellW/3) //nolint:mnd
		for i := 0; i < width; i += 2 * dash {
			fill(dst, image.Rect(px+i, uy, px+min(i+dash, width), uy+lw), c)
		}
	default:
		fill(dst, image.Rect(px, uy, px+width, uy+lw), c)
	}
}

func (r *Renderer) drawCursor(dst *image.RGBA, buf uv.ScreenBuffer, cur cursorState) {
	cell := buf.CellAt(cur.pos.X, cur.pos.Y)
	width := 1
	if cell != nil && cell.Width > 1 {
		width = cell.Width
	}
	rect := r.cellRect(cur.pos.X, cur.pos.Y, width)
	cc := toRGBA(r.cursor)
	switch cur.style {
	case vt.CursorUnderline:
		fill(dst, image.Rect(rect.Min.X, rect.Max.Y-r.lineW, rect.Max.X, rect.Max.Y), cc)
	case vt.CursorBar:
		fill(dst, image.Rect(rect.Min.X, rect.Min.Y, rect.Min.X+r.lineW, rect.Max.Y), cc)
	default:
		fill(dst, rect, cc)
		if cell != nil && cell.Content != "" && cell.Content != " " {
			var style fontStyle
			if cell.Style.Attrs&uv.AttrBold != 0 {
				style |= styleBold
			}
			if cell.Style.Attrs&uv.AttrItalic != 0 {
				style |= styleItalic
			}
			r.drawText(dst, rect.Min.X, rect.Min.Y, cell.Content, style, toRGBA(r.accent))
		}
	}
}

func fill(dst *image.RGBA, rect image.Rectangle, c color.RGBA) {
	draw.Draw(dst, rect, image.NewUniform(c), image.Point{}, draw.Over)
}

// toRGBA converts any color to an opaque-premultiplied RGBA.
func toRGBA(c color.Color) color.RGBA {
	if c == nil {
		return color.RGBA{A: 0xff}
	}
	if rgba, ok := c.(color.RGBA); ok {
		return rgba
	}
	return color.RGBAModel.Convert(c).(color.RGBA) //nolint:forcetypeassert
}

// blend mixes fg halfway into bg, used for faint text.
func blend(fg, bg color.RGBA) color.RGBA {
	return color.RGBA{
		R: uint8((uint16(fg.R) + uint16(bg.R)) / 2), //nolint:gosec,mnd
		G: uint8((uint16(fg.G) + uint16(bg.G)) / 2), //nolint:gosec,mnd
		B: uint8((uint16(fg.B) + uint16(bg.B)) / 2), //nolint:gosec,mnd
		A: 0xff,
	}
}
