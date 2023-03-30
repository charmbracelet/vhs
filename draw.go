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

func (c *circle) At(x, y int) color.Color {
	// Prepare points for circle calculations.
	// We subtract 1 from the radius to leave space for
	// antialiased pixels.
	xx := float64(x-c.p.X) + 0.5
	yy := float64(y-c.p.Y) + 0.5
	rr := float64(c.r) - 1

	// The distance from this pixel to the closest point
	// in the circle.
	dist := math.Sqrt(xx*xx+yy*yy) - rr

	if dist < 0 {
		// This pixel is inside the circle
		return color.Alpha{0xFF}
	} else if dist <= 1 {
		// This pixel is partly inside the circle
		// and needs antialiasing
		return color.Alpha{
			uint8((1 - dist) * 0xFF),
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
		return color.Alpha{0xFF}
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

	return color.Alpha{0xFF}
}

// Make a mask to round a terminal's corners
func MakeCornerMask(width, height, radius int, targetpng string) {
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
		img, img.Bounds(), &image.Uniform{color.Gray{0xFF}}, image.Point{0, 0},
		&roundedrect{image.Point{0, 0}, image.Point{width, height}, radius},
		image.Point{0, 0}, draw.Over,
	)

	f, err := os.Create(targetpng)
	if err != nil {
		fmt.Println(ErrorStyle.Render("Couldn't draw CornerMask: unable to save file."))
	} else {
		err = png.Encode(f, img)
	}

	if err != nil {
		fmt.Println(ErrorStyle.Render("Couldn't draw CornerMask: encoding failed."))
	}
}

// Check if a given windowbar type is valid
func CheckBar(windowbar string) bool {
	switch windowbar {
	case
		"",
		"Colorful",
		"ColorfulRight",
		"Rings",
		"RingsRight":
		return true
	}
	return false
}

// Make a window bar and save it to a file
func MakeBar(termWidth, termHeight int, opts VideoOptions, targetpng string) {
	var err error
	switch opts.WindowBar {
	case "Colorful":
		err = makeColorfulBar(termWidth, termHeight, false, opts, targetpng)
	case "ColorfulRight":
		err = makeColorfulBar(termWidth, termHeight, true, opts, targetpng)
	case "Rings":
		err = makeRingBar(termWidth, termHeight, false, opts, targetpng)
	case "RingsRight":
		err = makeRingBar(termWidth, termHeight, true, opts, targetpng)
	}

	if err != nil {
		fmt.Println(ErrorStyle.Render("Couldn't draw Bar: encoding failed"))
	}
}

func makeColorfulBar(termWidth int, termHeight int, isRight bool, opts VideoOptions, targetpng string) error {
	// Radius of dots
	dotrad := opts.WindowBarSize / 6
	// Space between dots and edge
	dotgap := (opts.WindowBarSize - 2*dotrad) / 2
	// Space between dot centers
	dotspace := (2 * dotrad) + opts.WindowBarSize/6

	// Dimensions of bar image
	width := termWidth
	height := termHeight + opts.WindowBarSize

	img := image.NewRGBA(
		image.Rectangle{
			image.Point{0, 0},
			image.Point{width, height},
		},
	)

	bg := color.RGBA{0x17, 0x17, 0x17, 0xFF}
	dota := color.RGBA{0xFF, 0x4F, 0x4D, 0xFF}
	dotb := color.RGBA{0xFE, 0xBB, 0x00, 0xFF}
	dotc := color.RGBA{0x00, 0xCC, 0x1D, 0xFF}

	var pta, ptb, ptc image.Point
	if isRight {
		pta = image.Point{termWidth - (dotgap + dotrad), dotrad + dotgap}
		ptb = image.Point{termWidth - (dotgap + dotrad + dotspace), dotrad + dotgap}
		ptc = image.Point{termWidth - (dotgap + dotrad + 2*dotspace), dotrad + dotgap}
	} else {
		pta = image.Point{dotgap + dotrad, dotrad + dotgap}
		ptb = image.Point{dotgap + dotrad + dotspace, dotrad + dotgap}
		ptc = image.Point{dotgap + dotrad + 2*dotspace, dotrad + dotgap}
	}

	draw.DrawMask(
		img, img.Bounds(), &image.Uniform{bg}, image.Point{0, 0},
		&rect{image.Point{0, 0}, image.Point{width, height}},
		image.Point{0, 0}, draw.Src,
	)

	draw.DrawMask(
		img,
		img.Bounds(),
		&image.Uniform{dota},
		image.Point{0, 0},
		&circle{pta, dotrad},
		image.Point{0, 0},
		draw.Over,
	)

	draw.DrawMask(
		img,
		img.Bounds(),
		&image.Uniform{dotb},
		image.Point{0, 0},
		&circle{ptb, dotrad},
		image.Point{0, 0},
		draw.Over,
	)

	draw.DrawMask(
		img,
		img.Bounds(),
		&image.Uniform{dotc},
		image.Point{0, 0},
		&circle{ptc, dotrad},
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

func makeRingBar(termWidth int, termHeight int, isRight bool, opts VideoOptions, targetpng string) error {
	// Radius of dots
	outerrad := opts.WindowBarSize / 5
	innerrad := (4 * outerrad) / 5
	// Space between dots and edge
	ringgap := (opts.WindowBarSize - 2*outerrad) / 2
	// Space between dot centers
	ringspace := (2 * outerrad) + opts.WindowBarSize/6

	// Dimensions of bar image
	width := termWidth
	height := termHeight + opts.WindowBarSize

	img := image.NewRGBA(
		image.Rectangle{
			image.Point{0, 0},
			image.Point{width, height},
		},
	)

	bg := color.RGBA{0x11, 0x11, 0x11, 0xFF}
	ring := color.RGBA{0x33, 0x33, 0x33, 0xFF}

	draw.DrawMask(
		img, img.Bounds(), &image.Uniform{bg}, image.Point{0, 0},
		&rect{image.Point{0, 0}, image.Point{width, height}},
		image.Point{0, 0}, draw.Src,
	)

	for i := 0; i <= 2; i++ {
		var pt image.Point
		if isRight {
			pt = image.Point{
				termWidth - (ringgap + outerrad + i*ringspace),
				outerrad + ringgap,
			}
		} else {
			pt = image.Point{
				ringgap + outerrad + i*ringspace,
				outerrad + ringgap,
			}
		}

		draw.DrawMask(
			img,
			img.Bounds(),
			&image.Uniform{ring},
			image.Point{0, 0},
			&circle{pt, outerrad},
			image.Point{0, 0},
			draw.Over,
		)

		draw.DrawMask(
			img,
			img.Bounds(),
			&image.Uniform{bg},
			image.Point{0, 0},
			&circle{pt, innerrad},
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
