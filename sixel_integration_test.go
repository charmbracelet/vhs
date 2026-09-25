//go:build integration

package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	imagegif "image/gif"
	"image/png"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func TestSixelRecording(t *testing.T) {
	dir := t.TempDir()
	// The image arrives after recording starts, leaves, then returns in green.
	tape := fmt.Sprintf(`Output "%s/frames/"
Output "%s/out.gif"
Set Shell bash
Set Width 800
Set Height 400
Set FontSize 18
Set Framerate 24
Set Sixel true
Set LoopOffset 50%%
Hide
Type "clear"
Enter
Show
Sleep 300ms
Screenshot "%s/blank.png"
Sleep 300ms
Hide
Type `+"`"+`printf '\033Pq#0;2;100;0;0#0!300~-!300~-!300~-!300~-!300~\033\\\n'`+"`"+`
Enter
Show
Sleep 300ms
Screenshot "%s/red.png"
Sleep 300ms
Hide
Type "clear"
Enter
Show
Sleep 300ms
Screenshot "%s/cleared.png"
Sleep 300ms
Hide
Type `+"`"+`printf '\033Pq#0;2;0;100;0#0!300~-!300~-!300~-!300~-!300~\033\\\n'`+"`"+`
Enter
Show
Sleep 300ms
Screenshot "%s/green.png"
Sleep 300ms
`, dir, dir, dir, dir, dir, dir)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if errs := Evaluate(ctx, tape, io.Discard); len(errs) != 0 {
		t.Fatal(errs)
	}
	for _, name := range []string{"blank", "red", "cleared", "green"} {
		data, err := os.ReadFile(filepath.Join(dir, name+".png"))
		requireNoErr(t, err)
		img, err := png.Decode(bytes.NewReader(data))
		requireNoErr(t, err)
		var red, green int
		for y := 0; y < img.Bounds().Dy(); y++ {
			for x := 0; x < img.Bounds().Dx(); x++ {
				c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
				if c.R > 220 && c.G < 30 && c.B < 30 {
					red++
				}
				if c.G > 220 && c.R < 30 && c.B < 30 {
					green++
				}
			}
		}
		if (name == "red") != (red > 8000) || (name == "green") != (green > 8000) {
			t.Fatalf("%s: red pixels=%d, green pixels=%d", name, red, green)
		}
	}
	frames, err := filepath.Glob(filepath.Join(dir, "frames", "frame-text-*.png"))
	requireNoErr(t, err)
	if len(frames) == 0 {
		t.Fatal("no exported frames")
	}
	for _, frame := range frames {
		for _, layer := range []string{"cursor", "image"} {
			_, err := os.Stat(strings.Replace(frame, "frame-text-", "frame-"+layer+"-", 1))
			requireNoErr(t, err)
		}
	}
	data, err := os.ReadFile(filepath.Join(dir, "out.gif"))
	requireNoErr(t, err)
	_, err = imagegif.DecodeAll(bytes.NewReader(data))
	requireNoErr(t, err)
}

// Run with go test -tags integration; requires a local Chrome and FFmpeg.
func TestImageLayerLifecycle(t *testing.T) {
	bin, found := launcher.LookPath()
	if !found {
		t.Fatal("integration test requires Chrome or Chromium")
	}
	u := launcher.New().Bin(bin).Leakless(false).NoSandbox(os.Getenv("VHS_NO_SANDBOX") != "").MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	t.Cleanup(func() { browser.MustClose() })
	page := browser.MustPage().Timeout(15 * time.Second)
	page.MustEval(`() => {
		const text = document.createElement('canvas');
		text.className = 'xterm-text-layer'; text.width = 40; text.height = 20;
		document.body.appendChild(text);
	}`)
	var capture imageLayerCapture
	check := func(want color.NRGBA, width, height int) {
		t.Helper()
		data, err := capture.capture(page)
		requireNoErr(t, err)
		img, err := png.Decode(bytes.NewReader(data))
		requireNoErr(t, err)
		if img.Bounds() != image.Rect(0, 0, width, height) {
			t.Fatalf("incorrect dimensions: %v", img.Bounds())
		}
		if got := color.NRGBAModel.Convert(img.At(width/2, height/2)).(color.NRGBA); got != want {
			t.Fatalf("pixel: %v; want %v", got, want)
		}
	}
	check(color.NRGBA{}, 40, 20)
	page.MustEval(`() => {
		const img = document.createElement('canvas');
		img.className = 'xterm-image-layer'; img.width = 20; img.height = 10;
		document.body.appendChild(img);
		const ctx = img.getContext('2d', {alpha: true, desynchronized: true});
		ctx.fillStyle = '#ff0000'; ctx.fillRect(0, 0, 20, 10);
	}`)
	check(color.NRGBA{R: 255, A: 255}, 40, 20)
	page.MustEval(`() => document.querySelector('.xterm-image-layer').remove()`)
	check(color.NRGBA{}, 40, 20)
	page.MustEval(`() => {
		const img = document.createElement('canvas');
		img.className = 'xterm-image-layer'; img.width = 40; img.height = 20;
		document.body.appendChild(img);
	}`)
	check(color.NRGBA{}, 40, 20)
	page.MustEval(`() => {
		const ctx = document.querySelector('.xterm-image-layer').getContext('2d');
		ctx.fillStyle = '#00ff00'; ctx.fillRect(0, 0, 40, 20);
	}`)
	check(color.NRGBA{G: 255, A: 255}, 40, 20)
	page.MustEval(`() => document.querySelector('.xterm-image-layer').width = 0`)
	check(color.NRGBA{}, 40, 20)
	page.MustEval(`() => document.querySelector('.xterm-text-layer').width = 80`)
	check(color.NRGBA{}, 80, 20)
	page.MustEval(`() => document.querySelector('.xterm-text-layer').remove()`)
	_, err := capture.capture(page)
	requireErr(t, err)
}

func TestSixelFFmpegPixels(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Fatal("integration test requires FFmpeg")
	}
	for _, enabled := range []bool{false, true} {
		for _, decorated := range []bool{false, true} {
			t.Run(fmt.Sprintf("sixel=%t/decorated=%t", enabled, decorated), func(t *testing.T) {
				dir := t.TempDir()
				style := &StyleOptions{Width: 80, Height: 60, Padding: 2, BackgroundColor: "#000000"}
				x, y := 2, 2
				if decorated {
					style.Margin = 10
					style.MarginFill = "#000000"
					style.WindowBar = "Colorful"
					style.WindowBarSize = 30
					style.WindowBarColor = "#000000"
					style.BorderRadius = 4
					style.Width += 20
					style.Height += 50
					x += 10
					y += 40
				}
				// Leave a trailing frame for the fps filter's end-of-stream rounding.
				for frame := 1; frame <= 3; frame++ {
					for _, format := range []string{textFrameFormat, cursorFrameFormat, imageFrameFormat} {
						img := image.NewNRGBA(image.Rect(0, 0, 76, 56))
						switch format {
						case textFrameFormat:
							draw.Draw(img, img.Bounds(), image.NewUniform(color.NRGBA{B: 255, A: 255}), image.Point{}, draw.Src)
						case cursorFrameFormat:
							draw.Draw(img, image.Rect(30, 20, 40, 30), image.NewUniform(color.NRGBA{G: 255, A: 255}), image.Point{}, draw.Src)
						case imageFrameFormat:
							if frame == 1 {
								draw.Draw(img, image.Rect(10, 10, 50, 40), image.NewUniform(color.NRGBA{R: 255, A: 255}), image.Point{}, draw.Src)
							}
						}
						var buf bytes.Buffer
						requireNoErr(t, png.Encode(&buf, img))
						requireNoErr(t, os.WriteFile(filepath.Join(dir, fmt.Sprintf(format, frame)), buf.Bytes(), 0o600))
					}
				}
				check := func(img image.Image, frame int) {
					t.Helper()
					if img.Bounds().Dx() != style.Width || img.Bounds().Dy() != style.Height {
						t.Fatalf("output size: %v", img.Bounds())
					}
					wantImage := color.NRGBA{B: 255, A: 255}
					if enabled && frame == 1 {
						wantImage = color.NRGBA{R: 255, A: 255}
					}
					for _, p := range []struct {
						x, y int
						want color.NRGBA
					}{
						{5, 5, color.NRGBA{B: 255, A: 255}},
						{20, 20, wantImage},
						{35, 25, color.NRGBA{G: 255, A: 255}},
					} {
						got := color.NRGBAModel.Convert(img.At(x+p.x, y+p.y)).(color.NRGBA)
						// FFmpeg's overlay filter converts through YUV by default.
						for channel, value := range []uint8{got.R, got.G, got.B, got.A} {
							want := []uint8{p.want.R, p.want.G, p.want.B, p.want.A}[channel]
							if diff := int(value) - int(want); diff < -5 || diff > 5 {
								t.Fatalf("frame %d pixel (%d,%d): %v; want approximately %v", frame, p.x, p.y, got, p.want)
							}
						}
					}
				}
				run := func(args []string) {
					t.Helper()
					ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
					defer cancel()
					if out, err := exec.CommandContext(ctx, "ffmpeg", args...).CombinedOutput(); err != nil {
						t.Fatalf("ffmpeg: %v\n%s", err, out)
					}
				}
				shot := NewScreenshotOptions(dir, style)
				shot.sixel = enabled
				for frame := 1; frame <= 2; frame++ {
					path := filepath.Join(dir, "out.png")
					run(shot.buildFFopts(path, filepath.Join(dir, fmt.Sprintf(textFrameFormat, frame)), filepath.Join(dir, fmt.Sprintf(cursorFrameFormat, frame)), filepath.Join(dir, fmt.Sprintf(imageFrameFormat, frame))))
					data, err := os.ReadFile(path)
					requireNoErr(t, err)
					img, err := png.Decode(bytes.NewReader(data))
					requireNoErr(t, err)
					check(img, frame)
				}
				path := filepath.Join(dir, "out.gif")
				run(buildFFopts(VideoOptions{Input: dir, Style: style, Framerate: 25, StartingFrame: 1, PlaybackSpeed: 1, Sixel: enabled}, path))
				data, err := os.ReadFile(path)
				requireNoErr(t, err)
				animation, err := imagegif.DecodeAll(bytes.NewReader(data))
				requireNoErr(t, err)
				if len(animation.Image) < 2 {
					t.Fatal("expected both video frames")
				}
				// GIF encoders may store only changed rectangles in later frames.
				composite := image.NewNRGBA(image.Rect(0, 0, style.Width, style.Height))
				for i, frame := range animation.Image {
					draw.Draw(composite, frame.Bounds(), frame, frame.Bounds().Min, draw.Over)
					check(composite, i+1)
				}
			})
		}
	}
}
