package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// The margin fill and window bar sources default to 25fps and are the *main*
// input of their overlay, so they must be pinned to the target framerate or
// they silently override `Set Framerate`.
func TestSecondaryStreamsUseFramerate(t *testing.T) {
	opts := DefaultVideoOptions()
	opts.Framerate = 60
	opts.Style = DefaultStyleOptions()
	opts.Style.WindowBar = "Colorful"
	opts.Style.MarginFill = "#000000"

	got := strings.Join(NewVideoFilterBuilder(&opts).
		WithWindowBar(2).
		WithMarginFill(3).
		Build(), " ")

	for _, want := range []string{"loop=loop=-1:size=1:start=0,fps=60[loopbar]", ",fps=60[bg]"} {
		if !strings.Contains(got, want) {
			t.Errorf("filter_complex missing %q\ngot: %s", want, got)
		}
	}
}

// Screenshots have no framerate; an fps filter there would be `fps=0`.
func TestScreenshotOmitsFramerate(t *testing.T) {
	style := DefaultStyleOptions()
	style.WindowBar = "Colorful"

	got := strings.Join(NewScreenshotFilterComplexBuilder(style).
		WithWindowBar(2).
		WithMarginFill(3).
		Build(), " ")

	if strings.Contains(got, "fps=") {
		t.Errorf("screenshot filter should not set fps, got: %s", got)
	}
}

func TestWindowBarEncodedOutput(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg unavailable")
	}
	if _, err := exec.LookPath("ffprobe"); err != nil {
		t.Skip("ffprobe unavailable")
	}
	for _, rate := range []int{10, 50} {
		t.Run(fmt.Sprint(rate), func(t *testing.T) {
			dir := t.TempDir()
			bar := filepath.Join(dir, "bar.png")
			output := filepath.Join(dir, "out.mp4")
			run := func(name string, args ...string) []byte {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				out, err := exec.CommandContext(ctx, name, args...).CombinedOutput()
				if err != nil {
					t.Fatalf("%s: %v: %s", name, err, out)
				}
				return out
			}
			run("ffmpeg", "-v", "error", "-f", "lavfi", "-i", "color=s=160x120", "-frames:v", "1", bar)
			opts := DefaultVideoOptions()
			opts.Framerate = rate
			opts.Style = DefaultStyleOptions()
			opts.Style.WindowBar = "Colorful"
			opts.Style.WindowBarSize = 20
			opts.Style.Width = 160
			opts.Style.Height = 140
			opts.Style.MarginFill = "#000000"
			fb := &FilterComplexBuilder{filterComplex: &strings.Builder{}, style: opts.Style, framerate: rate, prevStageName: "0:v"}
			fb.WithWindowBar(1).WithMarginFill(2)
			graph := strings.TrimPrefix(fb.filterComplex.String(), ";")
			run("ffmpeg", "-v", "error", "-f", "lavfi", "-i", fmt.Sprintf("testsrc2=size=160x100:rate=%d:duration=2", rate), "-i", bar, "-f", "lavfi", "-i", "color=s=160x140", "-filter_complex", graph, "-map", "[withbg]", "-c:v", "libx264", "-pix_fmt", "yuv420p", output)
			var probe struct {
				Streams []struct {
					Frames   string `json:"nb_frames"`
					Duration string `json:"duration"`
				}
			}
			if err := json.Unmarshal(run("ffprobe", "-v", "error", "-show_entries", "stream=nb_frames,duration", "-of", "json", output), &probe); err != nil {
				t.Fatal(err)
			}
			if len(probe.Streams) != 1 {
				t.Fatalf("expected one video stream, got %d", len(probe.Streams))
			}
			if got := probe.Streams[0]; got.Frames != fmt.Sprint(2*rate) || got.Duration != "2.000000" {
				t.Fatalf("unexpected encoded output: %+v", got)
			}
		})
	}
}
