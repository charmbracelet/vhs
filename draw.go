package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
)

type circle struct {
	p image.Point
	r int
}

const white = 0xFF
const black = 0x17

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circle) Bounds() image.Rectangle {
	return image.Rect(
		c.p.X-c.r,
		c.p.Y-c.r,
		c.p.X+c.r,
		c.p.Y+c.r,
	)
}

const halfPixel = 0.5
const doublingFactor = 2

func double(i int) int { return i * doublingFactor }
func half(i int) int   { return i / doublingFactor }

func (c *circle) At(x, y int) color.Color {
	// Prepare points for circle calculations.
	// We subtract 1 from the radius to leave space for
	// antialiased pixels.
	xx := float64(x-c.p.X) + halfPixel
	yy := float64(y-c.p.Y) + halfPixel
	rr := float64(c.r) - 1

	// The distance from this pixel to the closest point
	// in the circle.
	dist := math.Sqrt(xx*xx+yy*yy) - rr

	if dist < 0 {
		// This pixel is inside the circle
		return color.Alpha{white}
	} else if dist <= 1 {
		// This pixel is partly inside the circle
		// and needs antialiasing
		return color.Alpha{
			uint8((1 - dist) * white),
		}
	}

	// This pixel is outside the circle
	// and should be fully transparent
	return color.Alpha{0x00}
}

type rect struct {
	pa image.Point
	pb image.Point
}

func (r *rect) ColorModel() color.Model {
	return color.AlphaModel
}

func (r *rect) Bounds() image.Rectangle {
	return image.Rect(r.pa.X, r.pa.Y, r.pb.X, r.pb.Y)
}

func (r *rect) At(x, y int) color.Color {
	if (x >= r.pa.X) &&
		(x < r.pb.X) &&
		(y >= r.pa.Y) &&
		(y < r.pb.Y) {
		return color.Alpha{white}
	}
	return color.Alpha{0x00}
}

type roundedrect struct {
	pa     image.Point
	pb     image.Point
	radius int
}

func (r *roundedrect) ColorModel() color.Model {
	return color.AlphaModel
}

func (r *roundedrect) Bounds() image.Rectangle {
	return image.Rect(r.pa.X, r.pa.Y, r.pb.X, r.pb.Y)
}

func (r *roundedrect) At(x, y int) color.Color {
	// Top-left corner
	if (x >= r.pa.X) &&
		(x < r.pa.X+r.radius) &&
		(y >= r.pa.Y) &&
		(y < r.pa.Y+r.radius) {
		c := circle{
			image.Point{
				r.radius,
				r.radius,
			},
			// Add one to corner radius so that
			// fully-opaque pixels match the rectangle.
			// The outermost pixels of a circle are
			// always antialiased and thus transparent.
			r.radius + 1,
		}
		return c.At(x, y)
	}

	// Top-right corner
	if (x >= r.pb.X-r.radius) &&
		(x < r.pb.X) &&
		(y >= r.pa.Y) &&
		(y < r.pa.Y+r.radius) {
		c := circle{
			image.Point{
				r.pb.X - r.radius,
				r.radius,
			},
			r.radius + 1,
		}
		return c.At(x, y)
	}

	// Bottom-left corner
	if (x >= r.pa.X) &&
		(x < r.pa.X+r.radius) &&
		(y >= r.pb.Y-r.radius) &&
		(y < r.pb.Y) {
		c := circle{
			image.Point{
				r.radius,
				r.pb.Y - r.radius,
			},
			r.radius + 1,
		}
		return c.At(x, y)
	}

	// Bottom-right corner
	if (x >= r.pb.X-r.radius) &&
		(x < r.pb.X) &&
		(y >= r.pb.Y-r.radius) &&
		(y < r.pb.Y) {
		c := circle{
			image.Point{
				r.pb.X - r.radius,
				r.pb.Y - r.radius,
			},
			r.radius + 1,
		}
		return c.At(x, y)
	}

	return color.Alpha{white}
}

// Make a mask to round a terminal's corners
func MakeBorderRadiusMask(width, height, radius int, targetpng string) {
	img := image.NewGray(
		image.Rectangle{
			image.Point{0, 0},
			image.Point{width, height},
		},
	)

	// Fill image with black
	draw.DrawMask(
		img, img.Bounds(), &image.Uniform{color.Gray{0x00}}, image.Point{0, 0},
		&rect{image.Point{0, 0}, image.Point{width, height}},
		image.Point{0, 0}, draw.Src,
	)

	// Put mask in white on top
	draw.DrawMask(
		img, img.Bounds(), &image.Uniform{color.Gray{white}}, image.Point{0, 0},
		&roundedrect{image.Point{0, 0}, image.Point{width, height}, radius},
		image.Point{0, 0}, draw.Over,
	)

	f, err := os.Create(targetpng)
	if err != nil {
		fmt.Println(ErrorStyle.Render("Could not draw Border Mask: unable to save file."))
	} else {
		err = png.Encode(f, img)
	}

	if err != nil {
		fmt.Println(ErrorStyle.Render("Could not draw Border Mask: encoding failed."))
	}
}

// Make a window bar and save it to a file
func MakeWindowBar(termWidth, termHeight int, opts StyleOptions, file string) {
	var err error
	switch opts.WindowBar {
	case "Colorful":
		err = makeColorfulBar(termWidth, termHeight, false, opts, file)
	case "ColorfulRight":
		err = makeColorfulBar(termWidth, termHeight, true, opts, file)
	case "Rings":
		err = makeRingBar(termWidth, termHeight, false, opts, file)
	case "RingsRight":
		err = makeRingBar(termWidth, termHeight, true, opts, file)
	}

	if err != nil {
		fmt.Println(ErrorStyle.Render("Couldn't draw Bar: encoding failed"))
	}
}

const barToDotRatio = 6
const barToDotBorderRatio = 5

func makeColorfulBar(termWidth int, termHeight int, isRight bool, opts StyleOptions, targetpng string) error {
	// Radius of dots
	dotRad := opts.WindowBarSize / barToDotRatio
	dotDia := double(dotRad)
	// Space between dots and edge
	dotGap := half(opts.WindowBarSize - dotDia)
	// Space between dot centers
	dotSpace := dotDia + opts.WindowBarSize/barToDotRatio

	// Dimensions of bar image
	width := termWidth
	height := termHeight + opts.WindowBarSize

	img := image.NewRGBA(
		image.Rectangle{
			image.Point{0, 0},
			image.Point{width, height},
		},
	)

	bg, _ := parseHexColor(opts.WindowBarColor)
	dotA := color.RGBA{white, 0x4F, 0x4D, white}
	dotB := color.RGBA{0xFE, 0xBB, 0x00, white}
	dotC := color.RGBA{0x00, 0xCC, 0x1D, white}

	var pta, ptb, ptc image.Point
	if isRight {
		pta = image.Point{termWidth - (dotGap + dotRad), dotRad + dotGap}
		ptb = image.Point{termWidth - (dotGap + dotRad + dotSpace), dotRad + dotGap}
		ptc = image.Point{termWidth - (dotGap + dotRad + 2*dotSpace), dotRad + dotGap}
	} else {
		pta = image.Point{dotGap + dotRad, dotRad + dotGap}
		ptb = image.Point{dotGap + dotRad + dotSpace, dotRad + dotGap}
		ptc = image.Point{dotGap + dotRad + 2*dotSpace, dotRad + dotGap}
	}

	draw.DrawMask(
		img, img.Bounds(), &image.Uniform{bg}, image.Point{0, 0},
		&rect{image.Point{0, 0}, image.Point{width, height}},
		image.Point{0, 0}, draw.Src,
	)

	draw.DrawMask(
		img,
		img.Bounds(),
		&image.Uniform{dotA},
		image.Point{0, 0},
		&circle{pta, dotRad},
		image.Point{0, 0},
		draw.Over,
	)

	draw.DrawMask(
		img,
		img.Bounds(),
		&image.Uniform{dotB},
		image.Point{0, 0},
		&circle{ptb, dotRad},
		image.Point{0, 0},
		draw.Over,
	)

	draw.DrawMask(
		img,
		img.Bounds(),
		&image.Uniform{dotC},
		image.Point{0, 0},
		&circle{ptc, dotRad},
		image.Point{0, 0},
		draw.Over,
	)

	f, err := os.Create(targetpng)
	if err != nil {
		fmt.Println(ErrorStyle.Render("Couldn't draw colorful bar: unable to save file."))
	} else {
		err = png.Encode(f, img)
	}
	return err
}

func makeRingBar(termWidth int, termHeight int, isRight bool, opts StyleOptions, targetpng string) error {
	// Radius of dots
	outerRad := opts.WindowBarSize / barToDotBorderRatio
	outerDia := double(outerRad)
	innerRad := double(outerDia) / barToDotBorderRatio
	// Space between dots and edge
	ringGap := half(opts.WindowBarSize - outerDia)
	// Space between dot centers
	ringSpace := outerDia + opts.WindowBarSize/barToDotRatio

	// Dimensions of bar image
	width := termWidth
	height := termHeight + opts.WindowBarSize

	img := image.NewRGBA(
		image.Rectangle{
			image.Point{0, 0},
			image.Point{width, height},
		},
	)

	bg, _ := parseHexColor(opts.WindowBarColor)
	ring := color.RGBA{0x33, 0x33, 0x33, white}

	draw.DrawMask(
		img, img.Bounds(), &image.Uniform{bg}, image.Point{0, 0},
		&rect{image.Point{0, 0}, image.Point{width, height}},
		image.Point{0, 0}, draw.Src,
	)

	for i := 0; i <= 2; i++ {
		var pt image.Point
		if isRight {
			pt = image.Point{
				termWidth - (ringGap + outerRad + i*ringSpace),
				outerRad + ringGap,
			}
		} else {
			pt = image.Point{
				ringGap + outerRad + i*ringSpace,
				outerRad + ringGap,
			}
		}

		draw.DrawMask(
			img,
			img.Bounds(),
			&image.Uniform{ring},
			image.Point{0, 0},
			&circle{pt, outerRad},
			image.Point{0, 0},
			draw.Over,
		)

		draw.DrawMask(
			img,
			img.Bounds(),
			&image.Uniform{bg},
			image.Point{0, 0},
			&circle{pt, innerRad},
			image.Point{0, 0},
			draw.Over,
		)
	}

	f, err := os.Create(targetpng)
	if err != nil {
		fmt.Println(ErrorStyle.Render("Couldn't draw ring bar: unable to save file."))
	} else {
		err = png.Encode(f, img)
	}
	return err
}

func parseHexColor(s string) (c color.RGBA, err error) {
	c.R, c.G, c.B, c.A = black, black, black, white
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 6:
		_, err = fmt.Sscanf(s, "%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	case 3:
		_, err = fmt.Sscanf(s, "%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("%s color of invalid length", s)
	}
	return
}
