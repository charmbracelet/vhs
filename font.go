package main

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	"path/filepath"
	"strings"
	"sync"

	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const (
	// Font family constants.
	monospaceFont = "monospace"
	
	// Font positioning constants.
	fontAscentRatio  = 0.8 // Font ascent is typically 80% of font size
	fontPaddingRatio = 0.2 // Top padding is 20% of font size
	minFontPadding   = 4   // Minimum padding in pixels

	// Window bar layout constants.
	windowControlPadding = 100 // Horizontal padding for window controls
)

var (
	fontLoader     *FontLoader
	fontLoaderOnce sync.Once
)

// getFontLoader returns the singleton font loader.
func getFontLoader() *FontLoader {
	fontLoaderOnce.Do(func() {
		fontLoader = NewFontLoader()
	})
	return fontLoader
}

// getDefaultFont returns the basic font face to use for text rendering.
func getDefaultFont() font.Face {
	// Using basicfont.Face7x13 as it's built-in and always available
	// This provides good readability at small sizes
	return basicfont.Face7x13
}

// getWindowBarFont returns the font to use for window bar titles.
func getWindowBarFont(fontFamily string, fontSize float64) font.Face {
	loader := getFontLoader()

	// Try to load the requested font
	face, err := loader.LoadFont(fontFamily, fontSize)
	if err != nil {
		// Fall back to a common monospace font
		return loader.GetFallbackFont(fontSize)
	}

	return face
}

// measureText returns the width of the given text in pixels.
func measureText(face font.Face, text string) int {
	d := &font.Drawer{Face: face}
	return d.MeasureString(text).Round()
}

// drawCenteredText draws text centered horizontally at the given y position.
func drawCenteredText(img draw.Image, face font.Face, text string, y int, color image.Image) {
	if strings.TrimSpace(text) == "" {
		return
	}

	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X

	// Reserve space for window controls and padding
	horizontalPadding := windowControlPadding
	maxTextWidth := width - (2 * horizontalPadding)

	// Truncate text if it's too wide
	displayText := text
	textWidth := measureText(face, displayText)

	// If text is too wide, truncate with ellipsis
	if textWidth > maxTextWidth && len(displayText) > 3 {
		ellipsis := "..."
		ellipsisWidth := measureText(face, ellipsis)
		targetWidth := maxTextWidth - ellipsisWidth

		// Binary search for the right truncation point
		for len(displayText) > 0 && measureText(face, displayText) > targetWidth {
			// Remove characters from the end
			runes := []rune(displayText)
			if len(runes) > 1 {
				displayText = string(runes[:len(runes)-1])
			} else {
				break
			}
		}

		if len(displayText) > 0 {
			displayText += ellipsis
		}
		textWidth = measureText(face, displayText)
	}

	// Calculate centered X position
	x := bounds.Min.X + (width-textWidth)/2

	// Use antialiased text rendering if possible
	if _, ok := face.(*basicfont.Face); ok {
		// Basic font doesn't support antialiasing, use simple drawer
		d := &font.Drawer{
			Dst:  img,
			Src:  color,
			Face: face,
			Dot:  fixed.P(x, y),
		}
		d.DrawString(displayText)
	} else {
		// Use antialiased rendering for TrueType fonts
		drawAntialiasedText(img, face, displayText, x, y, color)
	}
}

// drawAntialiasedText draws antialiased text at the given position.
func drawAntialiasedText(dst draw.Image, face font.Face, text string, x, y int, src image.Image) {
	// Create a temporary image for the text with alpha channel
	metrics := face.Metrics()
	textWidth := measureText(face, text)
	textHeight := (metrics.Ascent + metrics.Descent).Ceil()

	// Create temporary image with alpha channel for smooth rendering
	tmp := image.NewRGBA(image.Rect(0, 0, textWidth, textHeight))

	// Draw text to temporary image
	d := &font.Drawer{
		Dst:  tmp,
		Src:  src,
		Face: face,
		Dot:  fixed.P(0, metrics.Ascent.Round()),
	}
	d.DrawString(text)

	// Draw the temporary image to the destination with bilinear interpolation
	dstRect := image.Rect(x, y-metrics.Ascent.Round(), x+textWidth, y+metrics.Descent.Round())
	xdraw.BiLinear.Scale(dst, dstRect, tmp, tmp.Bounds(), xdraw.Over, nil)
}


// getTextYPositionForFont calculates the Y position for centered text with a specific font.
func getTextYPositionForFont(barHeight int, face font.Face, fontSize int) int {
	// Get font metrics
	metrics := face.Metrics()
	ascent := metrics.Ascent.Round()
	descent := metrics.Descent.Round()

	// Match SVG logic exactly:
	// SVG uses the actual font size and calculates:
	// - fontAscent = fontSize * 0.8
	// - padding = max(4, fontSize * 0.2)
	// - baseline = padding + fontAscent

	// Calculate font ascent using SVG's approach
	fontAscent := int(float64(fontSize) * fontAscentRatio)

	// Calculate minimum top padding
	minPadding := int(float64(fontSize) * fontPaddingRatio)
	if minPadding < minFontPadding {
		minPadding = minFontPadding
	}

	// Use the same calculation as SVG: baseline = padding + fontAscent
	baseline := minPadding + fontAscent

	// For very small bars, ensure we fit
	if baseline > barHeight-2 {
		// If we can't fit with proper padding, just center as best we can
		baseline = barHeight - descent - 2
		if baseline < ascent {
			baseline = ascent // Absolute minimum
		}
	}

	return baseline
}

// FontLoader manages loading and caching of fonts from the filesystem.
// It provides a singleton interface for efficient font loading with caching support.
type FontLoader struct {
	cache map[string]font.Face
}

// NewFontLoader creates a new font loader.
func NewFontLoader() *FontLoader {
	return &FontLoader{
		cache: make(map[string]font.Face),
	}
}

// LoadFont attempts to load a font face with the given font family and size.
func (fl *FontLoader) LoadFont(fontFamily string, fontSize float64) (font.Face, error) {
	// Create cache key
	cacheKey := fmt.Sprintf("%s_%.2f", fontFamily, fontSize)

	// Check cache first
	if face, ok := fl.cache[cacheKey]; ok {
		return face, nil
	}

	// Try to load the font
	face, err := fl.loadFontFromFamily(fontFamily, fontSize)
	if err != nil {
		return nil, err
	}

	// Cache the loaded font
	fl.cache[cacheKey] = face
	return face, nil
}

// loadFontFromFamily attempts to load a font from the font family string.
// It supports CSS-style font family lists (e.g., "Monaco, Courier, monospace")
// and will try each font in order until one is found.
func (fl *FontLoader) loadFontFromFamily(fontFamily string, fontSize float64) (font.Face, error) {
	// Split the font family string into individual fonts
	fonts := strings.Split(fontFamily, ",")

	for _, fontName := range fonts {
		fontName = strings.TrimSpace(fontName)

		// Skip generic font families
		if fontName == monospaceFont || fontName == "ui-monospace" {
			continue
		}

		// Try to load this font
		face, err := fl.loadSingleFont(fontName, fontSize)
		if err == nil {
			return face, nil
		}
	}

	// If no font could be loaded, return error
	return nil, fmt.Errorf("could not load any font from family '%s': none of the specified fonts were found", fontFamily)
}

// loadSingleFont attempts to load a single font by name.
// It searches common font directories across different platforms.
func (fl *FontLoader) loadSingleFont(fontName string, fontSize float64) (font.Face, error) {
	// Map of font names to potential file paths
	fontPaths := fl.getFontPaths(fontName)

	for _, path := range fontPaths {
		face, err := fl.loadFontFromFile(path, fontSize)
		if err == nil {
			return face, nil
		}
	}

	return nil, fmt.Errorf("could not find font '%s' in system font directories", fontName)
}

// getFontPaths returns potential file paths for a font name.
// It searches system font directories on macOS, Linux, and Windows,
// and handles common font naming variations.
func (fl *FontLoader) getFontPaths(fontName string) []string {
	var paths []string

	// Common font directories on different platforms
	var fontDirs []string

	// macOS font directories
	fontDirs = append(fontDirs,
		"/System/Library/Fonts",
		"/Library/Fonts",
		filepath.Join(os.Getenv("HOME"), "Library/Fonts"),
	)

	// Linux font directories
	fontDirs = append(fontDirs,
		"/usr/share/fonts",
		"/usr/local/share/fonts",
		filepath.Join(os.Getenv("HOME"), ".fonts"),
		filepath.Join(os.Getenv("HOME"), ".local/share/fonts"),
	)

	// Windows font directories
	fontDirs = append(fontDirs,
		"C:\\Windows\\Fonts",
		filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "Windows", "Fonts"),
	)

	// Clean up font name for file matching
	cleanName := strings.ReplaceAll(fontName, " ", "")

	// Common font file extensions
	extensions := []string{".ttf", ".otf", ".ttc"}

	// Special cases for known fonts
	fontFileMap := map[string][]string{
		"JetBrains Mono":   {"JetBrainsMono-Regular.ttf", "JetBrainsMono.ttf"},
		"DejaVu Sans Mono": {"DejaVuSansMono.ttf", "DejaVu Sans Mono.ttf"},
		"Menlo":            {"Menlo.ttc", "Menlo-Regular.ttf"},
		"Monaco":           {"Monaco.ttf"},
		"Courier":          {"Courier.ttc", "Courier New.ttf"},
		"Consolas":         {"consola.ttf", "Consolas.ttf"},
		"Inconsolata":      {"Inconsolata-Regular.ttf", "Inconsolata.ttf"},
		"Roboto Mono":      {"RobotoMono-Regular.ttf", "Roboto Mono.ttf"},
		"Hack":             {"Hack-Regular.ttf", "Hack.ttf"},
	}

	// Check special cases first
	if fileNames, ok := fontFileMap[fontName]; ok {
		for _, dir := range fontDirs {
			for _, fileName := range fileNames {
				paths = append(paths, filepath.Join(dir, fileName))
			}
		}
	}

	// Generic search
	for _, dir := range fontDirs {
		for _, ext := range extensions {
			// Try exact match
			paths = append(paths, filepath.Join(dir, fontName+ext))
			// Try without spaces
			paths = append(paths, filepath.Join(dir, cleanName+ext))
			// Try with -Regular suffix
			paths = append(paths, filepath.Join(dir, fontName+"-Regular"+ext))
			paths = append(paths, filepath.Join(dir, cleanName+"-Regular"+ext))
		}
	}

	return paths
}

// loadFontFromFile loads a font from a file.
// It supports both TrueType (.ttf) and OpenType (.otf) fonts,
// as well as TrueType Collection (.ttc) files.
func (fl *FontLoader) loadFontFromFile(path string, fontSize float64) (font.Face, error) {
	// Check if file exists
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("stat font file: %w", err)
	}

	// Read font file
	fontData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read font file: %w", err)
	}

	// Parse font - try single font first
	tt, err := opentype.Parse(fontData)
	if err != nil {
		// Try parsing as a collection (TTC file)
		ttc, err2 := opentype.ParseCollection(fontData)
		if err2 != nil {
			return nil, fmt.Errorf("failed to parse font file '%s': %v (tried as TrueType/OpenType)", path, err)
		}
		// Use first font in collection
		if ttc.NumFonts() > 0 {
			tt, err = ttc.Font(0)
			if err != nil {
				return nil, fmt.Errorf("failed to get font from collection '%s': %v", path, err)
			}
		} else {
			return nil, fmt.Errorf("font collection '%s' is empty", path)
		}
	}

	// Create font face
	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("create font face: %w", err)
	}

	return face, nil
}

// GetFallbackFont returns a fallback font face.
// It tries common monospace fonts in order of preference,
// and falls back to the basic built-in font if none are found.
func (fl *FontLoader) GetFallbackFont(fontSize float64) font.Face {
	// Try to load a common monospace font
	commonFonts := []string{"Monaco", "Menlo", "Courier", "DejaVu Sans Mono"}

	for _, fontName := range commonFonts {
		face, err := fl.loadSingleFont(fontName, fontSize)
		if err == nil {
			return face
		}
	}

	// If all else fails, return the basic font
	return getDefaultFont()
}
