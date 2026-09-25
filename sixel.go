package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"

	"github.com/go-rod/rod"
)

// The image addon creates and removes its canvas as images enter and leave the
// screen. Query and read it in one evaluation so no detached element is retained.
// Image canvases may use CSS pixels while the text canvas uses device pixels.
const captureImageLayer = `() => {
	const text = document.querySelector('canvas.xterm-text-layer');
	const width = text.width, height = text.height;
	let image = document.querySelector('canvas.xterm-image-layer');
	if (!image || !image.width || !image.height) return {width, height, data: ''};
	if (image.width !== width || image.height !== height) {
		const scaled = document.createElement('canvas');
		scaled.width = width;
		scaled.height = height;
		scaled.getContext('2d').drawImage(image, 0, 0, width, height);
		image = scaled;
	}
	return {width, height, data: image.toDataURL('image/png').split(',')[1]};
}`

type imageLayerCapture struct {
	width, height int
	transparent   []byte
}

func (c *imageLayerCapture) capture(page *rod.Page) ([]byte, error) {
	result, err := page.Eval(captureImageLayer)
	if err != nil {
		return nil, fmt.Errorf("capture image layer: %w", err)
	}
	if data := result.Value.Get("data").Str(); data != "" {
		frame, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			return nil, fmt.Errorf("decode image layer: %w", err)
		}
		return frame, nil
	}
	width := result.Value.Get("width").Int()
	height := result.Value.Get("height").Int()
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid image frame dimensions: %dx%d", width, height)
	}
	if c.transparent == nil || c.width != width || c.height != height {
		var buf bytes.Buffer
		if err := png.Encode(&buf, image.NewNRGBA(image.Rect(0, 0, width, height))); err != nil {
			return nil, fmt.Errorf("encode transparent image frame: %w", err)
		}
		c.width, c.height = width, height
		c.transparent = buf.Bytes()
	}
	return c.transparent, nil
}
