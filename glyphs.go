package main

import (
	"image"
	"image/color"
)

// boxArms describes the box drawing characters U+2500 to U+257F as the weight
// of their up, right, down, and left arms: 0 none, 1 light, 2 heavy. "----"
// marks characters left to the font (dashes, doubles, diagonals).
//
// ponytail: double lines are left to the font, so they may not connect
// across cells when LineHeight > 1. Add a third weight if anyone needs them.
const boxArms = "" +
	"0101020210102020--------------------------------" + // 2500-250B
	"0110021001200220001100120021002211001200210022001001100220012002" + // 250C-251B
	"1110121021101120212022101220222010111012201110212021201210222022" + // 251C-252B
	"0111011202110212012101220221022211011102120112022101210222012202" + // 252C-253B
	"1111111212111212211111212121211222111122122122121222212222212222" + // 253C-254B
	"----------------" + // 254C-254F
	"--------------------------------------------------------------------------------------------------------------------" + // 2550-256C
	"0110001110011100" + // 256D-2570
	"------------" + // 2571-2573
	"000110000100001000022000020000200201102001022010" // 2574-257F

// drawCustomGlyph draws box drawing and block element characters as shapes
// that fill the whole cell, so they connect across cells at any line height
// like xterm.js' custom glyphs. It reports false for any other rune.
func (r *Renderer) drawCustomGlyph(dst *image.RGBA, px, py int, ch rune, c color.RGBA) bool {
	switch {
	case ch >= 0x2500 && ch <= 0x257f:
		return r.drawBox(dst, px, py, ch, c)
	case ch >= 0x2580 && ch <= 0x259f:
		r.drawBlock(dst, px, py, ch, c)
		return true
	}
	return false
}

func (r *Renderer) drawBox(dst *image.RGBA, px, py int, ch rune, c color.RGBA) bool {
	i := int(ch-0x2500) * 4 //nolint:mnd
	arms := boxArms[i : i+4]
	if arms[0] == '-' {
		return false
	}

	light := r.lineW
	thickness := func(a byte) int { return int(a-'0') * light }
	tm := 0
	for j := range 4 {
		tm = max(tm, thickness(arms[j]))
	}

	w, h := r.cellW, r.cellH
	cx, cy := px+w/2, py+h/2
	lo := func(t int) int { return t / 2 } //nolint:mnd
	if t := thickness(arms[0]); t > 0 {
		x0 := cx - lo(t)
		fill(dst, image.Rect(x0, py, x0+t, cy-lo(tm)+tm), c)
	}
	if t := thickness(arms[1]); t > 0 {
		y0 := cy - lo(t)
		fill(dst, image.Rect(cx-lo(tm), y0, px+w, y0+t), c)
	}
	if t := thickness(arms[2]); t > 0 {
		x0 := cx - lo(t)
		fill(dst, image.Rect(x0, cy-lo(tm), x0+t, py+h), c)
	}
	if t := thickness(arms[3]); t > 0 {
		y0 := cy - lo(t)
		fill(dst, image.Rect(px, y0, cx-lo(tm)+tm, y0+t), c)
	}
	return true
}

//nolint:mnd
func (r *Renderer) drawBlock(dst *image.RGBA, px, py int, ch rune, c color.RGBA) {
	w, h := r.cellW, r.cellH
	eighthW := func(n int) int { return (w*n + 4) / 8 }
	eighthH := func(n int) int { return (h*n + 4) / 8 }
	rect := func(x0, y0, x1, y1 int) { fill(dst, image.Rect(px+x0, py+y0, px+x1, py+y1), c) }

	switch {
	case ch == 0x2580: // upper half
		rect(0, 0, w, h/2)
	case ch >= 0x2581 && ch <= 0x2588: // lower eighths, full block
		rect(0, h-eighthH(int(ch-0x2580)), w, h)
	case ch >= 0x2589 && ch <= 0x258f: // left eighths
		rect(0, 0, eighthW(int(0x2590-ch)), h)
	case ch == 0x2590: // right half
		rect(w/2, 0, w, h)
	case ch >= 0x2591 && ch <= 0x2593: // shades
		shade := c
		shade.A = uint8(int(ch-0x2590) * 64)
		shade.R = uint8(int(c.R) * int(shade.A) / 255) //nolint:gosec
		shade.G = uint8(int(c.G) * int(shade.A) / 255) //nolint:gosec
		shade.B = uint8(int(c.B) * int(shade.A) / 255) //nolint:gosec
		fill(dst, image.Rect(px, py, px+w, py+h), shade)
	case ch == 0x2594: // upper eighth
		rect(0, 0, w, eighthH(1))
	case ch == 0x2595: // right eighth
		rect(w-eighthW(1), 0, w, h)
	default: // quadrants
		const ul, ur, ll, lr = 1, 2, 4, 8
		quads := [...]int{
			ll, lr, ul, ul | ll | lr, ul | lr, ul | ur | ll, ul | ur | lr, ur, ur | ll, ur | ll | lr,
		}[ch-0x2596]
		if quads&ul != 0 {
			rect(0, 0, w/2, h/2)
		}
		if quads&ur != 0 {
			rect(w/2, 0, w, h/2)
		}
		if quads&ll != 0 {
			rect(0, h/2, w/2, h)
		}
		if quads&lr != 0 {
			rect(w/2, h/2, w, h)
		}
	}
}
