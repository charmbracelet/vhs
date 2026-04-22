package main

import "testing"

func TestScreenshot(t *testing.T) {
	t.Run("makeScreenshot should add screenshot to map and disable capture", func(t *testing.T) {
		path := "sample.png"
		targetFrame := 60
		opts := ScreenshotOptions{
			nextScreenshotPath: path,
			frameCapture:       true,
			screenshots:        make(map[string]int),
		}

		opts.makeScreenshot(targetFrame)

		frame, ok := opts.screenshots[path]
		if !ok {
			t.Errorf("Unable to create screenshot: %s", path)
		}

		if frame != targetFrame {
			t.Errorf("Unable to create screenshot: %s", path)
		}

		if opts.frameCapture {
			t.Error("frameCapture should be false after invoking makeScreenshot")
		}
	})

	t.Run("enableFrameCapture", func(t *testing.T) {
		path := "sample.png"

		opts := ScreenshotOptions{
			frameCapture: false,
		}

		opts.enableFrameCapture(path)

		if !opts.frameCapture {
			t.Error("frameCapture should be true after invoking enableFrameCapture")
		}

		if opts.nextScreenshotPath != path {
			t.Errorf("nextScreenshotPath: %s, expected: %s", opts.nextScreenshotPath, path)
		}
	})
}
