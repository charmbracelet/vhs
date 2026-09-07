//go:build integration

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-rod/rod"
)

// Run with go test -tags integration -run '^$' -bench BenchmarkSixelCapture -benchtime=3s.
// Reports throughput of canvas readback and PNG writes, without the frame timer.
func BenchmarkSixelCapture(b *testing.B) {
	for _, mode := range []string{"disabled", "blank", "image"} {
		b.Run(mode, func(b *testing.B) {
			v := New()
			b.Cleanup(func() { _ = v.Cleanup() })
			v.Options.Video.Style.Width = 1200
			v.Options.Video.Style.Height = 700
			v.Options.Video.Style.Padding = 20
			v.Options.Video.Sixel = mode != "disabled"
			if err := v.Start(context.Background()); err != nil {
				b.Fatal(err)
			}
			b.Cleanup(func() { _ = v.terminate() })
			if err := v.Page.Wait(rod.Eval("() => window.term !== undefined")); err != nil {
				b.Fatal(err)
			}
			v.Setup()
			if mode == "image" {
				v.Page.MustEval(`() => new Promise(resolve => term.write('\x1bPq#0;2;100;0;0#0!300~-!300~-!300~-!300~-!300~\x1b\\', resolve))`)
				v.Page.Timeout(5 * time.Second).MustWait(`() => !!document.querySelector('canvas.xterm-image-layer')`)
			}
			var images imageLayerCapture
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cursor, err := v.CursorCanvas.CanvasToImage("image/png", quality)
				if err != nil {
					b.Fatal(err)
				}
				text, err := v.TextCanvas.CanvasToImage("image/png", quality)
				if err != nil {
					b.Fatal(err)
				}
				frames := [][]byte{text, cursor}
				if mode != "disabled" {
					img, err := images.capture(v.Page)
					if err != nil {
						b.Fatal(err)
					}
					frames = append(frames, img)
				}
				for layer, frame := range frames {
					if err := os.WriteFile(filepath.Join(v.Options.Video.Input, fmt.Sprintf("layer-%d.png", layer)), frame, 0o600); err != nil {
						b.Fatal(err)
					}
				}
			}
			b.StopTimer()
			b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "frames/s")
		})
	}
}

// Measure the recorder with its actual 60 fps timer and per-frame file names.
func TestSixelCaptureRate(t *testing.T) {
	if os.Getenv("VHS_MEASURE_CAPTURE") == "" {
		t.Skip("set VHS_MEASURE_CAPTURE=1 for capture-rate measurements")
	}
	for _, mode := range []string{"disabled", "blank", "image"} {
		t.Run(mode, func(t *testing.T) {
			v := New()
			t.Cleanup(func() { _ = v.Cleanup() })
			v.Options.Video.Style.Width = 1200
			v.Options.Video.Style.Height = 700
			v.Options.Video.Style.Padding = 20
			v.Options.Video.Framerate = 60
			v.Options.Video.Sixel = mode != "disabled"
			requireNoErr(t, v.Start(context.Background()))
			requireNoErr(t, v.Page.Wait(rod.Eval("() => window.term !== undefined")))
			v.Setup()
			if mode == "image" {
				v.Page.MustEval(`() => new Promise(resolve => term.write('\x1bPq#0;2;100;0;0#0!300~-!300~-!300~-!300~-!300~\x1b\\', resolve))`)
				v.Page.Timeout(5 * time.Second).MustWait(`() => !!document.querySelector('canvas.xterm-image-layer')`)
			}
			const duration = 3 * time.Second
			ctx, cancel := context.WithTimeout(context.Background(), duration)
			defer cancel()
			for err := range v.Record(ctx) {
				t.Error(err)
			}
			t.Logf("1200x700, requested 60 fps: %d frames in %s, %.1f fps", v.totalFrames, duration, float64(v.totalFrames)/duration.Seconds())
		})
	}
}
