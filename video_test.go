package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	imagegif "image/gif"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRenderGIF(t *testing.T) {
	if os.Getenv("VHS_TEST_FFMPEG") == "" {
		t.Skip("set VHS_TEST_FFMPEG=1 to run this test")
	}
	for _, tc := range []struct {
		style         StyleOptions
		frames, delay int
	}{
		{StyleOptions{Width: 320, Height: 200, BackgroundColor: "#000000"}, 9, 2},
		{StyleOptions{Width: 320, Height: 200, BackgroundColor: "#000000", WindowBar: "Colorful", WindowBarSize: 30}, 6, 4},
		{StyleOptions{Width: 320, Height: 200, BackgroundColor: "#000000", BorderRadius: 10}, 9, 2},
		{StyleOptions{Width: 320, Height: 200, BackgroundColor: "#000000", Margin: 20, MarginFill: "#ffffff", WindowBar: "Colorful", WindowBarSize: 30, BorderRadius: 10}, 5, 4},
	} {
		style := tc.style
		t.Run(fmt.Sprintf("bar=%s/radius=%d/margin=%d", style.WindowBar, style.BorderRadius, style.Margin), func(t *testing.T) {
			testRenderGIF(t, &style, 10, tc.frames, tc.delay)
		})
	}
}

func testRenderGIF(t *testing.T, style *StyleOptions, frames, wantFrames, delay int) {
	t.Helper()
	v := New()
	t.Cleanup(func() { _ = v.Cleanup() })
	v.totalFrames = frames
	v.Options.Video.Style = style
	v.Options.Video.Output.GIF = filepath.Join(t.TempDir(), "output.gif")
	bounds := image.Rect(0, 0, style.Width, style.Height)
	colors := []color.RGBA{{R: 200, A: 255}, {B: 200, A: 255}}
	for _, name := range []string{textFrameFormat, cursorFrameFormat} {
		first := ""
		for i := 1; i <= frames; i++ {
			path := filepath.Join(v.Options.Video.Input, fmt.Sprintf(name, i))
			if i == 1 || i == frames/2+1 {
				img := image.NewRGBA(bounds)
				if name == textFrameFormat {
					draw.Draw(img, bounds, image.NewUniform(colors[(i-1)/(frames/2)]), image.Point{}, draw.Src)
				}
				file, err := os.Create(path)
				requireNoErr(t, err)
				requireNoErr(t, png.Encode(file, img))
				requireNoErr(t, file.Close())
				first = path
			} else {
				requireNoErr(t, os.Link(first, path))
			}
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := v.Render(ctx); err != nil {
		t.Fatalf("render failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(v.Options.Video.Input, "palette.png")); !os.IsNotExist(err) {
		t.Fatalf("palette was not removed: %v", err)
	}
	file, err := os.Open(v.Options.Video.Output.GIF)
	requireNoErr(t, err)
	defer file.Close()
	animation, err := imagegif.DecodeAll(file)
	requireNoErr(t, err)
	if len(animation.Image) != wantFrames {
		t.Fatalf("expected %d frames, got %d", wantFrames, len(animation.Image))
	}
	duration := 0
	for _, delay := range animation.Delay {
		duration += delay
	}
	if duration != wantFrames*delay {
		t.Fatalf("expected duration %d, got %d", wantFrames*delay, duration)
	}
	last := wantFrames - 1
	if style.WindowBar != "" {
		last--
	}
	checkFrames := []int{0, last}
	for i, frame := range checkFrames {
		img := animation.Image[frame]
		x, y := style.Width/2, style.Height/2
		want := colors[i]
		got := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
		if abs(int(got.R)-int(want.R)) > 1 || abs(int(got.G)-int(want.G)) > 1 || abs(int(got.B)-int(want.B)) > 1 || got.A != want.A {
			t.Fatalf("frame %d: expected color %v at (%d, %d), got %v", frame, want, x, y, got)
		}
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
