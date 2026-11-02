package main

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"strconv"
	"strings"

	"golang.org/x/image/vector"
)

// Box drawing characters, block elements, and powerline symbols are drawn as
// shapes filling the cell rather than with a font, like xterm.js does with
// its customGlyphs option: lines join up between cells at any line height, and
// they don't need a font that has them. The shapes are xterm.js' own.

// blockRect is a rectangle of a block element, in eighths of the cell.
type blockRect struct{ x, y, w, h int }

// blockRects are the block elements and legacy computing symbols xterm.js
// draws, U+2580 to U+259F and some of U+1FB70 to U+1FB97.
var blockRects = map[rune][]blockRect{
	'▀': {{0, 0, 8, 4}}, '▁': {{0, 7, 8, 1}}, '▂': {{0, 6, 8, 2}}, '▃': {{0, 5, 8, 3}},
	'▄': {{0, 4, 8, 4}}, '▅': {{0, 3, 8, 5}}, '▆': {{0, 2, 8, 6}}, '▇': {{0, 1, 8, 7}},
	'█': {{0, 0, 8, 8}}, '▉': {{0, 0, 7, 8}}, '▊': {{0, 0, 6, 8}}, '▋': {{0, 0, 5, 8}},
	'▌': {{0, 0, 4, 8}}, '▍': {{0, 0, 3, 8}}, '▎': {{0, 0, 2, 8}}, '▏': {{0, 0, 1, 8}},
	'▐': {{4, 0, 4, 8}}, '▔': {{0, 0, 8, 1}}, '▕': {{7, 0, 1, 8}},
	'▖': {{0, 4, 4, 4}}, '▗': {{4, 4, 4, 4}}, '▘': {{0, 0, 4, 4}},
	'▙': {{0, 0, 4, 8}, {0, 4, 8, 4}}, '▚': {{0, 0, 4, 4}, {4, 4, 4, 4}},
	'▛': {{0, 0, 4, 8}, {4, 0, 4, 4}}, '▜': {{0, 0, 8, 4}, {4, 0, 4, 8}},
	'▝': {{4, 0, 4, 4}}, '▞': {{4, 0, 4, 4}, {0, 4, 4, 4}}, '▟': {{4, 0, 4, 8}, {0, 4, 8, 4}},
	'🭰': {{1, 0, 1, 8}}, '🭱': {{2, 0, 1, 8}}, '🭲': {{3, 0, 1, 8}}, '🭳': {{4, 0, 1, 8}},
	'🭴': {{5, 0, 1, 8}}, '🭵': {{6, 0, 1, 8}}, '🭶': {{0, 1, 8, 1}}, '🭷': {{0, 2, 8, 1}},
	'🭸': {{0, 3, 8, 1}}, '🭹': {{0, 4, 8, 1}}, '🭺': {{0, 5, 8, 1}}, '🭻': {{0, 6, 8, 1}},
	'🭼': {{0, 0, 1, 8}, {0, 7, 8, 1}}, '🭽': {{0, 0, 1, 8}, {0, 0, 8, 1}},
	'🭾': {{7, 0, 1, 8}, {0, 0, 8, 1}}, '🭿': {{7, 0, 1, 8}, {0, 7, 8, 1}},
	'🮀': {{0, 0, 8, 1}, {0, 7, 8, 1}}, '🮁': {{0, 0, 8, 1}, {0, 2, 8, 1}, {0, 4, 8, 1}, {0, 7, 8, 1}},
	'🮂': {{0, 0, 8, 2}}, '🮃': {{0, 0, 8, 3}}, '🮄': {{0, 0, 8, 5}}, '🮅': {{0, 0, 8, 6}},
	'🮆': {{0, 0, 8, 7}}, '🮇': {{6, 0, 2, 8}}, '🮈': {{5, 0, 3, 8}}, '🮉': {{3, 0, 5, 8}},
	'🮊': {{2, 0, 6, 8}}, '🮋': {{1, 0, 7, 8}},
	'🮕': {{0, 0, 2, 2}, {4, 0, 2, 2}, {2, 2, 2, 2}, {6, 2, 2, 2}, {0, 4, 2, 2}, {4, 4, 2, 2}, {2, 6, 2, 2}, {6, 6, 2, 2}},
	'🮖': {{2, 0, 2, 2}, {6, 0, 2, 2}, {0, 2, 2, 2}, {4, 2, 2, 2}, {2, 4, 2, 2}, {6, 4, 2, 2}, {0, 6, 2, 2}, {4, 6, 2, 2}},
	'🮗': {{0, 2, 8, 2}, {0, 6, 8, 2}},
}

// shadePatterns are the pixel patterns the shade characters are filled with.
var shadePatterns = map[rune][][]uint8{
	'░': {{1, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 1, 0}, {0, 0, 0, 0}},
	'▒': {{1, 0}, {0, 0}, {0, 1}, {0, 0}},
	'▓': {{0, 1}, {1, 1}, {1, 0}, {1, 1}},
}

// powerlineGlyph is a powerline symbol: a path in fractions of the cell,
// filled or stroked, with padding on the sides in twenty-fourths of the font
// size.
type powerlineGlyph struct {
	path        string
	stroke      bool
	left, right float64
}

// powerlineGlyphs are the powerline symbols xterm.js draws.
var powerlineGlyphs = map[rune]powerlineGlyph{
	0xe0b0: {path: "M0,0 L1,.5 L0,1", right: 2},
	0xe0b1: {path: "M-1,-.5 L1,.5 L-1,1.5", stroke: true, left: 1, right: 1},
	0xe0b2: {path: "M1,0 L0,.5 L1,1", left: 2},
	0xe0b3: {path: "M2,-.5 L0,.5 L2,1.5", stroke: true, left: 1, right: 1},
	0xe0b4: {path: "M0,0 L0,1 C0.552,1,1,0.776,1,.5 C1,0.224,0.552,0,0,0", right: 1},
	0xe0b5: {path: "M.2,1 C.422,1,.8,.826,.78,.5 C.8,.174,0.422,0,.2,0", stroke: true, right: 1},
	0xe0b6: {path: "M1,0 L1,1 C0.448,1,0,0.776,0,.5 C0,0.224,0.448,0,1,0", left: 1},
	0xe0b7: {path: "M.8,1 C0.578,1,0.2,.826,.22,.5 C0.2,0.174,0.578,0,0.8,0", stroke: true, left: 1},
	0xe0b8: {path: "M-.5,-.5 L1.5,1.5 L-.5,1.5"},
	0xe0b9: {path: "M-.5,-.5 L1.5,1.5", stroke: true, left: 1, right: 1},
	0xe0ba: {path: "M1.5,-.5 L-.5,1.5 L1.5,1.5"},
	0xe0bb: {path: "M1.5,-.5 L-.5,1.5", stroke: true, left: 1, right: 1},
	0xe0bc: {path: "M1.5,-.5 L-.5,1.5 L-.5,-.5"},
	0xe0bd: {path: "M1.5,-.5 L-.5,1.5", stroke: true, left: 1, right: 1},
	0xe0be: {path: "M-.5,-.5 L1.5,1.5 L1.5,-.5"},
	0xe0bf: {path: "M-.5,-.5 L1.5,1.5", stroke: true, left: 1, right: 1},
}

// drawCustomGlyph draws box drawing characters, block elements, and powerline
// symbols as shapes that fill the cell. It reports false for any other rune.
func (r *Renderer) drawCustomGlyph(dst *image.RGBA, px, py int, ch rune, c color.RGBA) bool {
	mask, ok := r.custom[ch]
	if !ok {
		if !isCustomGlyph(ch) {
			return false
		}
		mask = r.customGlyph(ch)
		r.custom[ch] = mask
	}
	draw.DrawMask(dst, mask.Rect.Add(image.Pt(px, py)), image.NewUniform(c), image.Point{}, mask, image.Point{}, draw.Over)
	return true
}

func isCustomGlyph(ch rune) bool {
	_, block := blockRects[ch]
	_, shade := shadePatterns[ch]
	_, box := boxStrokes[ch]
	_, powerline := powerlineGlyphs[ch]
	return block || shade || box || powerline
}

// customGlyph draws the coverage of a custom glyph in a cell.
func (r *Renderer) customGlyph(ch rune) *image.Alpha {
	w, h := r.cellW, r.cellH
	mask := image.NewAlpha(image.Rect(0, 0, w, h))

	if rects, ok := blockRects[ch]; ok {
		// Rounded to whole pixels, so neighboring blocks meet without seams.
		eighth := func(n, size int) int { return (size*n + 4) / 8 } //nolint:mnd
		for _, b := range rects {
			rect := image.Rect(eighth(b.x, w), eighth(b.y, h), eighth(b.x+b.w, w), eighth(b.y+b.h, h))
			draw.Draw(mask, rect, image.Opaque, image.Point{}, draw.Src)
		}
		return mask
	}

	if pattern, ok := shadePatterns[ch]; ok {
		for y := range h {
			row := pattern[y%len(pattern)]
			for x := range w {
				if row[x%len(row)] != 0 {
					mask.Pix[y*mask.Stride+x] = 0xff
				}
			}
		}
		return mask
	}

	z := vector.NewRasterizer(w, h)
	if strokes, ok := boxStrokes[ch]; ok {
		strokeBox(z, strokes, w, h)
	} else if g, ok := powerlineGlyphs[ch]; ok {
		unit := r.fonts.size / 12                    //nolint:mnd
		left, right := g.left*unit/2, g.right*unit/2 //nolint:mnd
		paths := parseGlyphPath(g.path,
			func(x float64) float64 { return x*(float64(w)-left-right) + left },
			func(y float64) float64 { return y * float64(h) },
		)
		if g.stroke {
			strokePath(z, paths, unit)
		} else {
			fillPath(z, paths)
		}
	}
	z.Draw(mask, mask.Rect, image.Opaque, image.Point{})
	return mask
}

// strokeBox strokes the lines of a box drawing character in a w by h cell.
// Lines are 1 or 3 pixels wide, and their ends are moved to the middle of a
// pixel so that they are crisp, as xterm.js does.
func strokeBox(z *vector.Rasterizer, strokes []boxStroke, w, h int) {
	fw, fh := float64(w), float64(h)
	snap := func(v, size float64) float64 {
		if v *= size; v != 0 {
			v = min(max(math.Floor(v+1)-0.5, 0), size)
		}
		return v
	}
	for _, s := range strokes {
		path := s.path
		if s.pathFn != nil {
			path = s.pathFn(0.15, 0.15/fh*fw) //nolint:mnd
		}
		strokePath(z, parseGlyphPath(path,
			func(x float64) float64 { return snap(x, fw) },
			func(y float64) float64 { return snap(y, fh) },
		), float64(s.weight))
	}
}

// point is a point of a glyph path, in pixels.
type point struct{ x, y float64 }

// parseGlyphPath parses a path of xterm.js glyph definitions, made of moves,
// lines, and cubic curves ("M0,.5 L1,.5 C.5,1,.5,.5,1,.5"), into lines
// joining points, mapping coordinates with fx and fy.
func parseGlyphPath(path string, fx, fy func(float64) float64) [][]point {
	var paths [][]point
	for _, ins := range strings.Fields(path) {
		var p []point
		var nums []float64
		for i, s := range strings.Split(ins[1:], ",") {
			v, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return paths
			}
			if i%2 == 0 {
				nums = append(nums, fx(v))
			} else {
				nums = append(nums, fy(v))
			}
		}
		for i := 0; i+1 < len(nums); i += 2 {
			p = append(p, point{nums[i], nums[i+1]})
		}
		switch {
		case ins[0] == 'M' && len(p) == 1:
			paths = append(paths, p)
		case len(paths) == 0:
		case ins[0] == 'L' && len(p) == 1:
			paths[len(paths)-1] = append(paths[len(paths)-1], p[0])
		case ins[0] == 'C' && len(p) == 3:
			cur := paths[len(paths)-1]
			p0 := cur[len(cur)-1]
			const steps = 16
			for i := 1; i <= steps; i++ {
				t := float64(i) / steps
				u := 1 - t
				cur = append(cur, point{
					u*u*u*p0.x + 3*u*u*t*p[0].x + 3*u*t*t*p[1].x + t*t*t*p[2].x,
					u*u*u*p0.y + 3*u*u*t*p[0].y + 3*u*t*t*p[1].y + t*t*t*p[2].y,
				})
			}
			paths[len(paths)-1] = cur
		}
	}
	return paths
}

// strokePath strokes lines of the given width, with square ends where lines
// join and flat ones where they end, like a canvas with miter joins and butt
// caps draws right angles.
func strokePath(z *vector.Rasterizer, paths [][]point, width float64) {
	for _, pts := range paths {
		for i := 0; i+1 < len(pts); i++ {
			a, b := pts[i], pts[i+1]
			dx, dy := b.x-a.x, b.y-a.y
			l := math.Hypot(dx, dy)
			if l == 0 {
				continue
			}
			ux, uy := dx/l, dy/l
			nx, ny := -uy*width/2, ux*width/2 //nolint:mnd
			if i > 0 {
				a.x, a.y = a.x-ux*width/2, a.y-uy*width/2 //nolint:mnd
			}
			if i+2 < len(pts) {
				b.x, b.y = b.x+ux*width/2, b.y+uy*width/2 //nolint:mnd
			}
			z.MoveTo(float32(a.x+nx), float32(a.y+ny))
			z.LineTo(float32(b.x+nx), float32(b.y+ny))
			z.LineTo(float32(b.x-nx), float32(b.y-ny))
			z.LineTo(float32(a.x-nx), float32(a.y-ny))
			z.ClosePath()
		}
	}
}

// fillPath fills the shapes the lines enclose.
func fillPath(z *vector.Rasterizer, paths [][]point) {
	for _, pts := range paths {
		z.MoveTo(float32(pts[0].x), float32(pts[0].y))
		for _, p := range pts[1:] {
			z.LineTo(float32(p.x), float32(p.y))
		}
		z.ClosePath()
	}
}
