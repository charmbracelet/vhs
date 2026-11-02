// Code generated from the custom glyphs of xterm.js 5.x (MIT license), the
// version VHS used to render with. DO NOT EDIT.

package main

import "fmt"

// boxStroke is a path xterm.js strokes to draw a box drawing character, in
// fractions of the cell, with its weight: 1 for light lines, 3 for heavy ones.
// Double lines and arcs depend on the cell's shape, through the offsets of
// double lines from the middle: xp of the width, and yp of the height.
type boxStroke struct {
	weight int
	path   string
	pathFn func(xp, yp float64) string
}

// boxStrokes are the box drawing characters, U+2500 to U+257F.
var boxStrokes = map[rune][]boxStroke{
	'─': {{weight: 1, path: "M0,.5 L1,.5"}},
	'━': {{weight: 3, path: "M0,.5 L1,.5"}},
	'│': {{weight: 1, path: "M.5,0 L.5,1"}},
	'┃': {{weight: 3, path: "M.5,0 L.5,1"}},
	'┌': {{weight: 1, path: "M0.5,1 L.5,.5 L1,.5"}},
	'┏': {{weight: 3, path: "M0.5,1 L.5,.5 L1,.5"}},
	'┐': {{weight: 1, path: "M0,.5 L.5,.5 L.5,1"}},
	'┓': {{weight: 3, path: "M0,.5 L.5,.5 L.5,1"}},
	'└': {{weight: 1, path: "M.5,0 L.5,.5 L1,.5"}},
	'┗': {{weight: 3, path: "M.5,0 L.5,.5 L1,.5"}},
	'┘': {{weight: 1, path: "M.5,0 L.5,.5 L0,.5"}},
	'┛': {{weight: 3, path: "M.5,0 L.5,.5 L0,.5"}},
	'├': {{weight: 1, path: "M.5,0 L.5,1 M.5,.5 L1,.5"}},
	'┣': {{weight: 3, path: "M.5,0 L.5,1 M.5,.5 L1,.5"}},
	'┤': {{weight: 1, path: "M.5,0 L.5,1 M.5,.5 L0,.5"}},
	'┫': {{weight: 3, path: "M.5,0 L.5,1 M.5,.5 L0,.5"}},
	'┬': {{weight: 1, path: "M0,.5 L1,.5 M.5,.5 L.5,1"}},
	'┳': {{weight: 3, path: "M0,.5 L1,.5 M.5,.5 L.5,1"}},
	'┴': {{weight: 1, path: "M0,.5 L1,.5 M.5,.5 L.5,0"}},
	'┻': {{weight: 3, path: "M0,.5 L1,.5 M.5,.5 L.5,0"}},
	'┼': {{weight: 1, path: "M0,.5 L1,.5 M.5,0 L.5,1"}},
	'╋': {{weight: 3, path: "M0,.5 L1,.5 M.5,0 L.5,1"}},
	'╴': {{weight: 1, path: "M.5,.5 L0,.5"}},
	'╸': {{weight: 3, path: "M.5,.5 L0,.5"}},
	'╵': {{weight: 1, path: "M.5,.5 L.5,0"}},
	'╹': {{weight: 3, path: "M.5,.5 L.5,0"}},
	'╶': {{weight: 1, path: "M.5,.5 L1,.5"}},
	'╺': {{weight: 3, path: "M.5,.5 L1,.5"}},
	'╷': {{weight: 1, path: "M.5,.5 L.5,1"}},
	'╻': {{weight: 3, path: "M.5,.5 L.5,1"}},
	'═': {{weight: 1, pathFn: func(xp, yp float64) string { return fmt.Sprintf("M0,%g L1,%g M0,%g L1,%g", .5-yp, .5-yp, .5+yp, .5+yp) }}},
	'║': {{weight: 1, pathFn: func(xp, yp float64) string { return fmt.Sprintf("M%g,0 L%g,1 M%g,0 L%g,1", .5-xp, .5-xp, .5+xp, .5+xp) }}},
	'╒': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M.5,1 L.5,%g L1,%g M.5,%g L1,%g", .5-yp, .5-yp, .5+yp, .5+yp)
	}}},
	'╓': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M%g,1 L%g,.5 L1,.5 M%g,.5 L%g,1", .5-xp, .5-xp, .5+xp, .5+xp)
	}}},
	'╔': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M1,%g L%g,%g L%g,1 M1,%g L%g,%g L%g,1", .5-yp, .5-xp, .5-yp, .5-xp, .5+yp, .5+xp, .5+yp, .5+xp)
	}}},
	'╕': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M0,%g L.5,%g L.5,1 M0,%g L.5,%g", .5-yp, .5-yp, .5+yp, .5+yp)
	}}},
	'╖': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M%g,1 L%g,.5 L0,.5 M%g,.5 L%g,1", .5+xp, .5+xp, .5-xp, .5-xp)
	}}},
	'╗': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M0,%g L%g,%g L%g,1 M0,%g L%g,%g L%g,1", .5+yp, .5-xp, .5+yp, .5-xp, .5-yp, .5+xp, .5-yp, .5+xp)
	}}},
	'╘': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M.5,0 L.5,%g L1,%g M.5,%g L1,%g", .5+yp, .5+yp, .5-yp, .5-yp)
	}}},
	'╙': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M1,.5 L%g,.5 L%g,0 M%g,.5 L%g,0", .5-xp, .5-xp, .5+xp, .5+xp)
	}}},
	'╚': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M1,%g L%g,%g L%g,0 M1,%g L%g,%g L%g,0", .5-yp, .5+xp, .5-yp, .5+xp, .5+yp, .5-xp, .5+yp, .5-xp)
	}}},
	'╛': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M0,%g L.5,%g L.5,0 M0,%g L.5,%g", .5+yp, .5+yp, .5-yp, .5-yp)
	}}},
	'╜': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M0,.5 L%g,.5 L%g,0 M%g,.5 L%g,0", .5+xp, .5+xp, .5-xp, .5-xp)
	}}},
	'╝': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M0,%g L%g,%g L%g,0 M0,%g L%g,%g L%g,0", .5-yp, .5-xp, .5-yp, .5-xp, .5+yp, .5+xp, .5+yp, .5+xp)
	}}},
	'╞': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M.5,0 L.5,1 M.5,%g L1,%g M.5,%g L1,%g", .5-yp, .5-yp, .5+yp, .5+yp)
	}}},
	'╟': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M%g,0 L%g,1 M%g,0 L%g,1 M%g,.5 L1,.5", .5-xp, .5-xp, .5+xp, .5+xp, .5+xp)
	}}},
	'╠': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M%g,0 L%g,1 M1,%g L%g,%g L%g,1 M1,%g L%g,%g L%g,0", .5-xp, .5-xp, .5+yp, .5+xp, .5+yp, .5+xp, .5-yp, .5+xp, .5-yp, .5+xp)
	}}},
	'╡': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M.5,0 L.5,1 M0,%g L.5,%g M0,%g L.5,%g", .5-yp, .5-yp, .5+yp, .5+yp)
	}}},
	'╢': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M0,.5 L%g,.5 M%g,0 L%g,1 M%g,0 L%g,1", .5-xp, .5-xp, .5-xp, .5+xp, .5+xp)
	}}},
	'╣': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M%g,0 L%g,1 M0,%g L%g,%g L%g,1 M0,%g L%g,%g L%g,0", .5+xp, .5+xp, .5+yp, .5-xp, .5+yp, .5-xp, .5-yp, .5-xp, .5-yp, .5-xp)
	}}},
	'╤': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M0,%g L1,%g M0,%g L1,%g M.5,%g L.5,1", .5-yp, .5-yp, .5+yp, .5+yp, .5+yp)
	}}},
	'╥': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M0,.5 L1,.5 M%g,.5 L%g,1 M%g,.5 L%g,1", .5-xp, .5-xp, .5+xp, .5+xp)
	}}},
	'╦': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M0,%g L1,%g M0,%g L%g,%g L%g,1 M1,%g L%g,%g L%g,1", .5-yp, .5-yp, .5+yp, .5-xp, .5+yp, .5-xp, .5+yp, .5+xp, .5+yp, .5+xp)
	}}},
	'╧': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M.5,0 L.5,%g M0,%g L1,%g M0,%g L1,%g", .5-yp, .5-yp, .5-yp, .5+yp, .5+yp)
	}}},
	'╨': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M0,.5 L1,.5 M%g,.5 L%g,0 M%g,.5 L%g,0", .5-xp, .5-xp, .5+xp, .5+xp)
	}}},
	'╩': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M0,%g L1,%g M0,%g L%g,%g L%g,0 M1,%g L%g,%g L%g,0", .5+yp, .5+yp, .5-yp, .5-xp, .5-yp, .5-xp, .5-yp, .5+xp, .5-yp, .5+xp)
	}}},
	'╪': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M.5,0 L.5,1 M0,%g L1,%g M0,%g L1,%g", .5-yp, .5-yp, .5+yp, .5+yp)
	}}},
	'╫': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M0,.5 L1,.5 M%g,0 L%g,1 M%g,0 L%g,1", .5-xp, .5-xp, .5+xp, .5+xp)
	}}},
	'╬': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M0,%g L%g,%g L%g,1 M1,%g L%g,%g L%g,1 M0,%g L%g,%g L%g,0 M1,%g L%g,%g L%g,0", .5+yp, .5-xp, .5+yp, .5-xp, .5+yp, .5+xp, .5+yp, .5+xp, .5-yp, .5-xp, .5-yp, .5-xp, .5-yp, .5+xp, .5-yp, .5+xp)
	}}},
	'╱': {{weight: 1, path: "M1,0 L0,1"}},
	'╲': {{weight: 1, path: "M0,0 L1,1"}},
	'╳': {{weight: 1, path: "M1,0 L0,1 M0,0 L1,1"}},
	'╼': {{weight: 1, path: "M.5,.5 L0,.5"}, {weight: 3, path: "M.5,.5 L1,.5"}},
	'╽': {{weight: 1, path: "M.5,.5 L.5,0"}, {weight: 3, path: "M.5,.5 L.5,1"}},
	'╾': {{weight: 1, path: "M.5,.5 L1,.5"}, {weight: 3, path: "M.5,.5 L0,.5"}},
	'╿': {{weight: 1, path: "M.5,.5 L.5,1"}, {weight: 3, path: "M.5,.5 L.5,0"}},
	'┍': {{weight: 1, path: "M.5,.5 L.5,1"}, {weight: 3, path: "M.5,.5 L1,.5"}},
	'┎': {{weight: 1, path: "M.5,.5 L1,.5"}, {weight: 3, path: "M.5,.5 L.5,1"}},
	'┑': {{weight: 1, path: "M.5,.5 L.5,1"}, {weight: 3, path: "M.5,.5 L0,.5"}},
	'┒': {{weight: 1, path: "M.5,.5 L0,.5"}, {weight: 3, path: "M.5,.5 L.5,1"}},
	'┕': {{weight: 1, path: "M.5,.5 L.5,0"}, {weight: 3, path: "M.5,.5 L1,.5"}},
	'┖': {{weight: 1, path: "M.5,.5 L1,.5"}, {weight: 3, path: "M.5,.5 L.5,0"}},
	'┙': {{weight: 1, path: "M.5,.5 L.5,0"}, {weight: 3, path: "M.5,.5 L0,.5"}},
	'┚': {{weight: 1, path: "M.5,.5 L0,.5"}, {weight: 3, path: "M.5,.5 L.5,0"}},
	'┝': {{weight: 1, path: "M.5,0 L.5,1"}, {weight: 3, path: "M.5,.5 L1,.5"}},
	'┞': {{weight: 1, path: "M0.5,1 L.5,.5 L1,.5"}, {weight: 3, path: "M.5,.5 L.5,0"}},
	'┟': {{weight: 1, path: "M.5,0 L.5,.5 L1,.5"}, {weight: 3, path: "M.5,.5 L.5,1"}},
	'┠': {{weight: 1, path: "M.5,.5 L1,.5"}, {weight: 3, path: "M.5,0 L.5,1"}},
	'┡': {{weight: 1, path: "M.5,.5 L.5,1"}, {weight: 3, path: "M.5,0 L.5,.5 L1,.5"}},
	'┢': {{weight: 1, path: "M.5,.5 L.5,0"}, {weight: 3, path: "M0.5,1 L.5,.5 L1,.5"}},
	'┥': {{weight: 1, path: "M.5,0 L.5,1"}, {weight: 3, path: "M.5,.5 L0,.5"}},
	'┦': {{weight: 1, path: "M0,.5 L.5,.5 L.5,1"}, {weight: 3, path: "M.5,.5 L.5,0"}},
	'┧': {{weight: 1, path: "M.5,0 L.5,.5 L0,.5"}, {weight: 3, path: "M.5,.5 L.5,1"}},
	'┨': {{weight: 1, path: "M.5,.5 L0,.5"}, {weight: 3, path: "M.5,0 L.5,1"}},
	'┩': {{weight: 1, path: "M.5,.5 L.5,1"}, {weight: 3, path: "M.5,0 L.5,.5 L0,.5"}},
	'┪': {{weight: 1, path: "M.5,.5 L.5,0"}, {weight: 3, path: "M0,.5 L.5,.5 L.5,1"}},
	'┭': {{weight: 1, path: "M0.5,1 L.5,.5 L1,.5"}, {weight: 3, path: "M.5,.5 L0,.5"}},
	'┮': {{weight: 1, path: "M0,.5 L.5,.5 L.5,1"}, {weight: 3, path: "M.5,.5 L1,.5"}},
	'┯': {{weight: 1, path: "M.5,.5 L.5,1"}, {weight: 3, path: "M0,.5 L1,.5"}},
	'┰': {{weight: 1, path: "M0,.5 L1,.5"}, {weight: 3, path: "M.5,.5 L.5,1"}},
	'┱': {{weight: 1, path: "M.5,.5 L1,.5"}, {weight: 3, path: "M0,.5 L.5,.5 L.5,1"}},
	'┲': {{weight: 1, path: "M.5,.5 L0,.5"}, {weight: 3, path: "M0.5,1 L.5,.5 L1,.5"}},
	'┵': {{weight: 1, path: "M.5,0 L.5,.5 L1,.5"}, {weight: 3, path: "M.5,.5 L0,.5"}},
	'┶': {{weight: 1, path: "M.5,0 L.5,.5 L0,.5"}, {weight: 3, path: "M.5,.5 L1,.5"}},
	'┷': {{weight: 1, path: "M.5,.5 L.5,0"}, {weight: 3, path: "M0,.5 L1,.5"}},
	'┸': {{weight: 1, path: "M0,.5 L1,.5"}, {weight: 3, path: "M.5,.5 L.5,0"}},
	'┹': {{weight: 1, path: "M.5,.5 L1,.5"}, {weight: 3, path: "M.5,0 L.5,.5 L0,.5"}},
	'┺': {{weight: 1, path: "M.5,.5 L0,.5"}, {weight: 3, path: "M.5,0 L.5,.5 L1,.5"}},
	'┽': {{weight: 1, path: "M.5,0 L.5,1 M.5,.5 L1,.5"}, {weight: 3, path: "M.5,.5 L0,.5"}},
	'┾': {{weight: 1, path: "M.5,0 L.5,1 M.5,.5 L0,.5"}, {weight: 3, path: "M.5,.5 L1,.5"}},
	'┿': {{weight: 1, path: "M.5,0 L.5,1"}, {weight: 3, path: "M0,.5 L1,.5"}},
	'╀': {{weight: 1, path: "M0,.5 L1,.5 M.5,.5 L.5,1"}, {weight: 3, path: "M.5,.5 L.5,0"}},
	'╁': {{weight: 1, path: "M.5,.5 L.5,0 M0,.5 L1,.5"}, {weight: 3, path: "M.5,.5 L.5,1"}},
	'╂': {{weight: 1, path: "M0,.5 L1,.5"}, {weight: 3, path: "M.5,0 L.5,1"}},
	'╃': {{weight: 1, path: "M0.5,1 L.5,.5 L1,.5"}, {weight: 3, path: "M.5,0 L.5,.5 L0,.5"}},
	'╄': {{weight: 1, path: "M0,.5 L.5,.5 L.5,1"}, {weight: 3, path: "M.5,0 L.5,.5 L1,.5"}},
	'╅': {{weight: 1, path: "M.5,0 L.5,.5 L1,.5"}, {weight: 3, path: "M0,.5 L.5,.5 L.5,1"}},
	'╆': {{weight: 1, path: "M.5,0 L.5,.5 L0,.5"}, {weight: 3, path: "M0.5,1 L.5,.5 L1,.5"}},
	'╇': {{weight: 1, path: "M.5,.5 L.5,1"}, {weight: 3, path: "M.5,.5 L.5,0 M0,.5 L1,.5"}},
	'╈': {{weight: 1, path: "M.5,.5 L.5,0"}, {weight: 3, path: "M0,.5 L1,.5 M.5,.5 L.5,1"}},
	'╉': {{weight: 1, path: "M.5,.5 L1,.5"}, {weight: 3, path: "M.5,0 L.5,1 M.5,.5 L0,.5"}},
	'╊': {{weight: 1, path: "M.5,.5 L0,.5"}, {weight: 3, path: "M.5,0 L.5,1 M.5,.5 L1,.5"}},
	'╌': {{weight: 1, path: "M.1,.5 L.4,.5 M.6,.5 L.9,.5"}},
	'╍': {{weight: 3, path: "M.1,.5 L.4,.5 M.6,.5 L.9,.5"}},
	'┄': {{weight: 1, path: "M.0667,.5 L.2667,.5 M.4,.5 L.6,.5 M.7333,.5 L.9333,.5"}},
	'┅': {{weight: 3, path: "M.0667,.5 L.2667,.5 M.4,.5 L.6,.5 M.7333,.5 L.9333,.5"}},
	'┈': {{weight: 1, path: "M.05,.5 L.2,.5 M.3,.5 L.45,.5 M.55,.5 L.7,.5 M.8,.5 L.95,.5"}},
	'┉': {{weight: 3, path: "M.05,.5 L.2,.5 M.3,.5 L.45,.5 M.55,.5 L.7,.5 M.8,.5 L.95,.5"}},
	'╎': {{weight: 1, path: "M.5,.1 L.5,.4 M.5,.6 L.5,.9"}},
	'╏': {{weight: 3, path: "M.5,.1 L.5,.4 M.5,.6 L.5,.9"}},
	'┆': {{weight: 1, path: "M.5,.0667 L.5,.2667 M.5,.4 L.5,.6 M.5,.7333 L.5,.9333"}},
	'┇': {{weight: 3, path: "M.5,.0667 L.5,.2667 M.5,.4 L.5,.6 M.5,.7333 L.5,.9333"}},
	'┊': {{weight: 1, path: "M.5,.05 L.5,.2 M.5,.3 L.5,.45 L.5,.55 M.5,.7 L.5,.95"}},
	'┋': {{weight: 3, path: "M.5,.05 L.5,.2 M.5,.3 L.5,.45 L.5,.55 M.5,.7 L.5,.95"}},
	'╭': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M.5,1 L.5,%g C.5,%g,.5,.5,1,.5", .5+yp/.15*.5, .5+yp/.15*.5)
	}}},
	'╮': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M.5,1 L.5,%g C.5,%g,.5,.5,0,.5", .5+yp/.15*.5, .5+yp/.15*.5)
	}}},
	'╯': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M.5,0 L.5,%g C.5,%g,.5,.5,0,.5", .5-yp/.15*.5, .5-yp/.15*.5)
	}}},
	'╰': {{weight: 1, pathFn: func(xp, yp float64) string {
		return fmt.Sprintf("M.5,0 L.5,%g C.5,%g,.5,.5,1,.5", .5-yp/.15*.5, .5-yp/.15*.5)
	}}},
}
