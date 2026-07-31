package main

import (
	"strings"
	"testing"
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

	for _, want := range []string{"loop=-1,fps=60[loopbar]", ",fps=60[bg]"} {
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
