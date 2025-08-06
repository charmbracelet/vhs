package main

import (
	"image"
	"strings"
	"testing"

	"golang.org/x/image/font/basicfont"
)

// Test singleton font loader
func TestGetFontLoader(t *testing.T) {
	loader1 := getFontLoader()
	loader2 := getFontLoader()
	
	if loader1 != loader2 {
		t.Error("getFontLoader should return the same singleton instance")
	}
	
	if loader1 == nil {
		t.Error("getFontLoader should not return nil")
	}
}

// Test default font
func TestGetDefaultFont(t *testing.T) {
	face := getDefaultFont()
	
	if face == nil {
		t.Error("getDefaultFont should not return nil")
	}
	
	// Should return basicfont.Face7x13
	if face != basicfont.Face7x13 {
		t.Error("getDefaultFont should return basicfont.Face7x13")
	}
}

// Test window bar font
func TestGetWindowBarFont(t *testing.T) {
	testCases := []struct {
		name       string
		fontFamily string
		fontSize   float64
		shouldWork bool
	}{
		{"valid monospace font", "monospace", 14.0, true},
		{"empty font family", "", 14.0, true}, // Should fall back
		{"invalid font", "ThisFontDoesNotExist123456", 14.0, true}, // Should fall back
		{"common font Monaco", "Monaco", 14.0, true}, // May or may not exist depending on OS
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			face := getWindowBarFont(tc.fontFamily, tc.fontSize)
			if face == nil {
				t.Error("getWindowBarFont should always return a font face (with fallback)")
			}
		})
	}
}

// Test measureText
func TestMeasureText(t *testing.T) {
	face := basicfont.Face7x13
	
	testCases := []struct {
		text     string
		expected int
	}{
		{"", 0},
		{"a", 7}, // basicfont.Face7x13 has 7 pixel wide characters
		{"hello", 35}, // 5 * 7
		{"Hello World", 77}, // 11 * 7
	}
	
	for _, tc := range testCases {
		width := measureText(face, tc.text)
		if width != tc.expected {
			t.Errorf("measureText(%q) = %d, expected %d", tc.text, width, tc.expected)
		}
	}
}

// Test drawCenteredText
func TestDrawCenteredText(t *testing.T) {
	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 200, 50))
	face := basicfont.Face7x13
	color := image.Black
	
	// Test empty text (should not panic)
	drawCenteredText(img, face, "", 25, color)
	
	// Test whitespace only (should not draw)
	drawCenteredText(img, face, "   ", 25, color)
	
	// Test normal text
	drawCenteredText(img, face, "Test", 25, color)
	
	// Test text that's too long
	longText := strings.Repeat("A", 50)
	drawCenteredText(img, face, longText, 25, color)
	
	// Verify no panics occurred
	t.Log("drawCenteredText completed without panics")
}


// Test text Y position with specific font
func TestGetTextYPositionForFont(t *testing.T) {
	face := basicfont.Face7x13
	
	testCases := []struct {
		barHeight int
		fontSize  int
	}{
		{30, 16},
		{40, 20},
		{20, 12},
		{10, 8}, // Very small bar
	}
	
	for _, tc := range testCases {
		y := getTextYPositionForFont(tc.barHeight, face, tc.fontSize)
		
		// Should return a reasonable value - note that baseline can be slightly outside 
		// the bar height for very small bars due to font metrics
		if y < 0 {
			t.Errorf("getTextYPositionForFont(%d, _, %d) = %d, should not be negative",
				tc.barHeight, tc.fontSize, y)
		}
	}
}

// Test FontLoader creation
func TestNewFontLoader(t *testing.T) {
	loader := NewFontLoader()
	
	if loader == nil {
		t.Fatal("NewFontLoader should not return nil")
	}
	
	if loader.cache == nil {
		t.Error("FontLoader cache should be initialized")
	}
}

// Test FontLoader.LoadFont
func TestFontLoader_LoadFont(t *testing.T) {
	loader := NewFontLoader()
	
	// Test loading a generic font
	_, err := loader.LoadFont("monospace", 14.0)
	if err == nil {
		t.Error("LoadFont should return error for generic font family 'monospace'")
	}
	
	// Test caching - load same font twice
	loader.LoadFont("InvalidFont123", 14.0)
	loader.LoadFont("InvalidFont123", 14.0) // Should use cache
	
	// Different size should create different cache entry
	loader.LoadFont("InvalidFont123", 16.0)
}

// Test FontLoader.GetFallbackFont
func TestFontLoader_GetFallbackFont(t *testing.T) {
	loader := NewFontLoader()
	
	face := loader.GetFallbackFont(14.0)
	if face == nil {
		t.Error("GetFallbackFont should always return a font face")
	}
}

// Test getFontPaths
func TestFontLoader_GetFontPaths(t *testing.T) {
	loader := NewFontLoader()
	
	testCases := []string{
		"Monaco",
		"JetBrains Mono",
		"Arial",
		"Unknown Font",
	}
	
	for _, fontName := range testCases {
		paths := loader.getFontPaths(fontName)
		
		if len(paths) == 0 {
			t.Errorf("getFontPaths(%q) returned no paths", fontName)
		}
		
		// Check that paths include various extensions
		hasttf := false
		hasotf := false
		for _, path := range paths {
			if strings.HasSuffix(path, ".ttf") {
				hasttf = true
			}
			if strings.HasSuffix(path, ".otf") {
				hasotf = true
			}
		}
		
		if !hasttf {
			t.Errorf("getFontPaths(%q) should include .ttf paths", fontName)
		}
		if !hasotf {
			t.Errorf("getFontPaths(%q) should include .otf paths", fontName)
		}
	}
}

// Test loadFontFromFile with invalid path
func TestFontLoader_LoadFontFromFile(t *testing.T) {
	loader := NewFontLoader()
	
	// Test non-existent file
	_, err := loader.loadFontFromFile("/this/does/not/exist.ttf", 14.0)
	if err == nil {
		t.Error("loadFontFromFile should return error for non-existent file")
	}
	
	// Test invalid font data
	// We can't easily test this without creating actual files
}

// Test drawAntialiasedText basic functionality
func TestDrawAntialiasedText(t *testing.T) {
	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 200, 50))
	face := basicfont.Face7x13
	color := image.Black
	
	// Should not panic
	drawAntialiasedText(img, face, "Test", 10, 25, color)
	
	t.Log("drawAntialiasedText completed without panics")
}

// Benchmark font loading
func BenchmarkFontLoader_LoadFont(b *testing.B) {
	loader := NewFontLoader()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Try to load a font that might exist
		loader.LoadFont("Monaco", 14.0)
	}
}

// Benchmark text measurement
func BenchmarkMeasureText(b *testing.B) {
	face := basicfont.Face7x13
	text := "The quick brown fox jumps over the lazy dog"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		measureText(face, text)
	}
}