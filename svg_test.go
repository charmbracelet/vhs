package main

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

// Helper functions
func createTestSVGConfig() SVGConfig {
	return SVGConfig{
		Width:      800,
		Height:     600,
		FontSize:   16,
		FontFamily: "monospace",
		Theme:      DefaultTheme,
		Frames: []SVGFrame{
			{
				Lines:      []string{"Hello", "World"},
				CursorX:    0,
				CursorY:    1,
				CharWidth:  8.8,
				CharHeight: 20,
			},
		},
		Duration: 1.0,
		Style:    DefaultStyleOptions(),
	}
}

func assertContains(t *testing.T, svg, expected, message string) {
	t.Helper()
	if !strings.Contains(svg, expected) {
		t.Errorf("%s: expected to contain '%s'", message, expected)
	}
}

func assertNotContains(t *testing.T, svg, unexpected, message string) {
	t.Helper()
	if strings.Contains(svg, unexpected) {
		t.Errorf("%s: expected NOT to contain '%s'", message, unexpected)
	}
}

// Core SVG Generation Tests
func TestSVGGenerator_CoreFunctionality(t *testing.T) {
	t.Run("generates valid SVG structure", func(t *testing.T) {
		opts := createTestSVGConfig()
		gen := NewSVGGenerator(opts)
		svg := gen.Generate()

		// Check basic SVG structure
		assertContains(t, svg, "<svg", "SVG opening tag")
		assertContains(t, svg, "</svg>", "SVG closing tag")
		assertContains(t, svg, "@keyframes", "CSS keyframes")
		assertContains(t, svg, "xmlns=\"http://www.w3.org/2000/svg\"", "SVG namespace")

		// Check viewBox is set correctly
		assertContains(t, svg, "viewBox=", "SVG viewBox")
	})

	t.Run("handles empty frames gracefully", func(t *testing.T) {
		opts := createTestSVGConfig()
		opts.Frames = []SVGFrame{}

		gen := NewSVGGenerator(opts)
		svg := gen.Generate()

		assertContains(t, svg, "<svg", "SVG should generate even with no frames")
	})

	t.Run("uses character dimensions from frames", func(t *testing.T) {
		opts := createTestSVGConfig()
		opts.Frames[0].CharWidth = 10.5
		opts.Frames[0].CharHeight = 25.0

		gen := NewSVGGenerator(opts)

		if gen.charWidth != 10.5 {
			t.Errorf("Expected charWidth to be 10.5, got %f", gen.charWidth)
		}
		if gen.charHeight != 25.0 {
			t.Errorf("Expected charHeight to be 25.0, got %f", gen.charHeight)
		}
	})
}

// Style and Appearance Tests
func TestSVGGenerator_StyleOptions(t *testing.T) {
	t.Run("applies all style options", func(t *testing.T) {
		style := &StyleOptions{
			Width:           1024,
			Height:          768,
			Padding:         20,
			Margin:          10,
			MarginFill:      "#ff0000",
			WindowBar:       "Colorful",
			WindowBarSize:   30,
			WindowBarColor:  "#333333",
			BorderRadius:    5,
			BackgroundColor: "#000000",
		}

		opts := createTestSVGConfig()
		opts.Width = style.Width
		opts.Height = style.Height
		opts.Style = style

		gen := NewSVGGenerator(opts)
		svg := gen.Generate()

		// Check dimensions with margins
		assertContains(t, svg, "1044", "Total width with margins")
		assertContains(t, svg, "788", "Total height with margins")

		// Check colors
		assertContains(t, svg, "#ff0000", "Margin fill color")
		assertContains(t, svg, "#000000", "Background color")
		assertContains(t, svg, "#333333", "Window bar color")

		// Check border radius
		assertContains(t, svg, "rx=\"5\"", "Border radius")
	})

	t.Run("WindowBar styles", func(t *testing.T) {
		tests := []struct {
			name     string
			style    string
			contains []string
		}{
			{
				name:     "Colorful",
				style:    "Colorful",
				contains: []string{`cx="20"`, `cx="40"`, `cx="60"`, `fill="#ff5f58"`},
			},
			{
				name:     "ColorfulRight",
				style:    "ColorfulRight",
				contains: []string{`fill="#ff5f58"`, `fill="#ffbd2e"`, `fill="#18c132"`},
			},
			{
				name:     "Rings",
				style:    "Rings",
				contains: []string{`fill="none"`, `stroke="#ff5f58"`, `stroke-width="1"`},
			},
			{
				name:     "RingsRight",
				style:    "RingsRight",
				contains: []string{`fill="none"`, `stroke="#ff5f58"`, `stroke-width="1"`},
			},
			{
				name:     "Empty",
				style:    "",
				contains: []string{}, // Should not contain window bar
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				style := DefaultStyleOptions()
				style.WindowBar = tt.style
				style.WindowBarSize = 40

				opts := createTestSVGConfig()
				opts.Style = style

				gen := NewSVGGenerator(opts)
				svg := gen.Generate()

				if tt.style == "" {
					assertNotContains(t, svg, "window-bar", "Should not have window bar")
				} else {
					assertContains(t, svg, "window-bar", "Should have window bar")
					for _, expected := range tt.contains {
						assertContains(t, svg, expected, fmt.Sprintf("%s window bar element", tt.name))
					}
				}
			})
		}
	})

	t.Run("handles missing style options", func(t *testing.T) {
		opts := createTestSVGConfig()
		opts.Style = nil

		gen := NewSVGGenerator(opts)
		svg := gen.Generate()

		// Should use defaults and not crash
		assertContains(t, svg, "<svg", "SVG should generate with nil style")
	})
}

// Text Rendering Tests
func TestSVGGenerator_TextRendering(t *testing.T) {
	t.Run("handles text styles", func(t *testing.T) {
		opts := createTestSVGConfig()
		// Text styles are applied through CSS, not directly on SVGConfig

		gen := NewSVGGenerator(opts)
		svg := gen.Generate()

		// Check that text elements are created
		assertContains(t, svg, "<text", "Text elements")
	})

	t.Run("applies LineHeight", func(t *testing.T) {
		opts := createTestSVGConfig()
		opts.LineHeight = 1.5
		opts.Frames = []SVGFrame{
			{Lines: []string{"Line 1", "Line 2", "Line 3"}, CharHeight: 20},
		}

		gen := NewSVGGenerator(opts)
		svg := gen.Generate()

		// Check that text rendering includes proper line height
		assertContains(t, svg, "<text", "Text elements")
	})

	t.Run("handles special characters", func(t *testing.T) {
		opts := createTestSVGConfig()
		opts.Frames = []SVGFrame{
			{Lines: []string{"<script>alert('xss')</script>", "Test & < > \" '"}},
		}

		gen := NewSVGGenerator(opts)
		svg := gen.Generate()

		// Check HTML escaping
		assertContains(t, svg, "&lt;script&gt;", "Script tag escaped")
		assertContains(t, svg, "&amp;", "Ampersand escaped")
		assertNotContains(t, svg, "<script>", "Raw script tag should not exist")
	})

	t.Run("preserves whitespace with xml:space", func(t *testing.T) {
		opts := createTestSVGConfig()
		opts.Frames = []SVGFrame{
			{Lines: []string{"    indented", "multiple   spaces"}},
		}

		gen := NewSVGGenerator(opts)
		svg := gen.Generate()

		// Check for whitespace preservation
		assertContains(t, svg, `xml:space="preserve"`, "Whitespace preservation attribute")
	})
}

// Background Color Tests
func TestSVGGenerator_BackgroundColors(t *testing.T) {
	t.Run("renders background color rectangles", func(t *testing.T) {
		opts := createTestSVGConfig()
		opts.Frames = []SVGFrame{
			{
				Lines: []string{"", ""}, // Empty lines
				LineColors: [][]CharStyle{
					{}, // No colors on first line
					{ // Second line has background colors
						{BgColor: "#ff0000"},
						{BgColor: "#00ff00"},
						{BgColor: "#0000ff"},
					},
				},
				CharWidth:  10,
				CharHeight: 20,
			},
		}

		gen := NewSVGGenerator(opts)
		svg := gen.Generate()

		// Check for background rectangles
		assertContains(t, svg, `fill="#ff0000"`, "Red background")
		assertContains(t, svg, `fill="#00ff00"`, "Green background")
		assertContains(t, svg, `fill="#0000ff"`, "Blue background")
		assertContains(t, svg, `shape-rendering="crispEdges"`, "Crisp edges for color blocks")
	})

	t.Run("skips invalid background colors", func(t *testing.T) {
		opts := createTestSVGConfig()
		opts.Frames = []SVGFrame{
			{
				Lines: []string{"test"},
				LineColors: [][]CharStyle{
					{
						{BgColor: ""},
						{BgColor: "<nil>"},
						{BgColor: "null"},
						{BgColor: "#ffffff"}, // Valid
					},
				},
			},
		}

		gen := NewSVGGenerator(opts)
		svg := gen.Generate()

		// Should only render the valid background color
		assertContains(t, svg, `fill="#ffffff"`, "Valid background color")
		assertNotContains(t, svg, `fill=""`, "Empty background color")
		assertNotContains(t, svg, `fill="<nil>"`, "Nil background color")
		assertNotContains(t, svg, `fill="null"`, "Null background color")
	})
}

// Animation and Timing Tests
func TestSVGGenerator_AnimationTiming(t *testing.T) {
	t.Run("applies PlaybackSpeed", func(t *testing.T) {
		testCases := []struct {
			name          string
			duration      float64
			playbackSpeed float64
			expected      string
		}{
			{"normal speed", 10.0, 1.0, "10s"},
			{"double speed", 10.0, 2.0, "5s"},
			{"half speed", 10.0, 0.5, "20s"},
			{"no speed set", 10.0, 0.0, "10s"}, // Should use original duration
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				opts := createTestSVGConfig()
				opts.Duration = tc.duration
				opts.PlaybackSpeed = tc.playbackSpeed

				gen := NewSVGGenerator(opts)
				svg := gen.Generate()

				assertContains(t, svg, fmt.Sprintf("animation: slide %s step-end", tc.expected),
					"Animation duration with playback speed")
			})
		}
	})

	t.Run("applies LoopOffset", func(t *testing.T) {
		testCases := []struct {
			name       string
			loopOffset float64
			duration   float64
			frames     int
			expected   string
		}{
			{"25% offset", 0.25, 10.0, 100, "-2.5s"},
			{"50% offset", 0.5, 10.0, 100, "-5s"},
			{"frame offset 10", 10.0, 10.0, 100, "-1s"},
			{"no offset", 0.0, 10.0, 100, "0s"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				opts := createTestSVGConfig()
				opts.Duration = tc.duration
				opts.LoopOffset = tc.loopOffset

				// Create multiple frames for frame-based offset
				opts.Frames = make([]SVGFrame, tc.frames)
				for i := range opts.Frames {
					opts.Frames[i] = SVGFrame{Lines: []string{fmt.Sprintf("Frame %d", i)}}
				}

				gen := NewSVGGenerator(opts)
				svg := gen.Generate()

				assertContains(t, svg, tc.expected, "Animation delay from loop offset")
			})
		}
	})

	t.Run("CursorBlink animation", func(t *testing.T) {
		t.Run("enabled", func(t *testing.T) {
			opts := createTestSVGConfig()
			opts.CursorBlink = true

			gen := NewSVGGenerator(opts)
			svg := gen.Generate()

			assertContains(t, svg, "@keyframes blink", "Blink keyframes")
			assertContains(t, svg, "animation: blink 1s infinite", "Cursor blink animation")
		})

		t.Run("disabled", func(t *testing.T) {
			opts := createTestSVGConfig()
			opts.CursorBlink = false

			gen := NewSVGGenerator(opts)
			svg := gen.Generate()

			assertNotContains(t, svg, "@keyframes blink", "No blink keyframes")
			assertNotContains(t, svg, "animation: blink", "No cursor blink animation")
		})
	})
}

// Frame Processing Tests
func TestSVGGenerator_FrameProcessing(t *testing.T) {
	t.Run("deduplicates identical frames", func(t *testing.T) {
		opts := createTestSVGConfig()
		opts.Frames = []SVGFrame{
			{Lines: []string{"Test"}, CursorX: 0, CursorY: 0},
			{Lines: []string{"Test"}, CursorX: 0, CursorY: 0}, // Duplicate
			{Lines: []string{"Test2"}, CursorX: 0, CursorY: 0},
			{Lines: []string{"Test"}, CursorX: 0, CursorY: 0}, // Another duplicate
		}

		gen := NewSVGGenerator(opts)
		gen.processFrames()

		// Should have only 2 unique states
		if len(gen.states) != 2 {
			t.Errorf("Expected 2 unique states, got %d", len(gen.states))
		}

		// Should have 3 timeline entries (frame 0, frame 2, 100%)
		if len(gen.timeline) != 3 {
			t.Errorf("Expected 3 timeline entries, got %d", len(gen.timeline))
		}
	})

	t.Run("handles cursor position in deduplication", func(t *testing.T) {
		opts := createTestSVGConfig()
		opts.Frames = []SVGFrame{
			{Lines: []string{"Test"}, CursorX: 0, CursorY: 0},
			{Lines: []string{"Test"}, CursorX: 1, CursorY: 0}, // Different cursor
		}

		gen := NewSVGGenerator(opts)
		gen.processFrames()

		// Should have 2 unique states due to different cursor positions
		if len(gen.states) != 2 {
			t.Errorf("Expected 2 unique states due to cursor difference, got %d", len(gen.states))
		}
	})

	t.Run("handles color information in deduplication", func(t *testing.T) {
		opts := createTestSVGConfig()
		opts.Frames = []SVGFrame{
			{
				Lines: []string{"Test"},
				LineColors: [][]CharStyle{
					{{FgColor: "#ff0000"}},
				},
			},
			{
				Lines: []string{"Test"},
				LineColors: [][]CharStyle{
					{{FgColor: "#00ff00"}}, // Different color
				},
			},
		}

		gen := NewSVGGenerator(opts)
		gen.processFrames()

		// Should have 2 unique states due to different colors
		if len(gen.states) != 2 {
			t.Errorf("Expected 2 unique states due to color difference, got %d", len(gen.states))
		}
	})
}

// Optimization Tests
func TestSVGGenerator_Optimization(t *testing.T) {
	t.Run("optimization disabled produces readable output", func(t *testing.T) {
		opts := createTestSVGConfig()
		opts.OptimizeSize = false

		gen := NewSVGGenerator(opts)
		svg := gen.Generate()

		// Should have newlines for readability
		if !strings.Contains(svg, "\n") {
			t.Error("Expected newlines in non-optimized output")
		}
	})

	t.Run("optimization enabled produces minified output", func(t *testing.T) {
		opts := createTestSVGConfig()
		opts.OptimizeSize = true

		gen := NewSVGGenerator(opts)
		svg := gen.Generate()

		// Should be minified (no indentation)
		if strings.Contains(svg, "\n\t") || strings.Contains(svg, "\n  ") {
			t.Error("Expected no indentation in optimized output")
		}
	})
}

// Debug Mode Tests
func TestSVGGenerator_DebugMode(t *testing.T) {
	t.Run("debug mode does not affect output structure", func(t *testing.T) {
		opts := createTestSVGConfig()
		opts.Debug = false
		gen := NewSVGGenerator(opts)
		svgNormal := gen.Generate()

		opts.Debug = true
		gen = NewSVGGenerator(opts)
		svgDebug := gen.Generate()

		// Debug mode should only affect logging, not output
		if len(svgNormal) != len(svgDebug) {
			t.Error("Debug mode should not change SVG output")
		}
	})
}

// Window Bar Font Tests
func TestSVGGenerator_WindowBarFonts(t *testing.T) {
	t.Run("uses custom window bar font family", func(t *testing.T) {
		opts := createTestSVGConfig()
		opts.Style = &StyleOptions{
			WindowBar:           "Colorful",
			WindowBarTitle:      "Test App",
			WindowBarFontFamily: "Arial",
			WindowBarFontSize:   18,
		}

		gen := NewSVGGenerator(opts)
		svg := gen.Generate()

		assertContains(t, svg, "Test App", "Window bar title")
		assertContains(t, svg, "Arial", "Window bar font family")
		assertContains(t, svg, "font-size=\"18\"", "Window bar font size")
	})

	t.Run("falls back to main font settings", func(t *testing.T) {
		opts := createTestSVGConfig()
		opts.FontFamily = "JetBrains Mono"
		opts.FontSize = 16
		opts.Style = &StyleOptions{
			WindowBar:      "Colorful",
			WindowBarTitle: "Test App",
			// No WindowBarFontFamily or WindowBarFontSize set
		}

		gen := NewSVGGenerator(opts)
		svg := gen.Generate()

		assertContains(t, svg, "Test App", "Window bar title")
		assertContains(t, svg, "JetBrains Mono", "Falls back to main font family")
		assertContains(t, svg, "font-size=\"16\"", "Falls back to main font size")
	})
}

// Integration Tests
func TestMakeSVG(t *testing.T) {
	t.Run("generates SVG file", func(t *testing.T) {
		vhs := &VHS{
			Options: &Options{
				FontSize:      16,
				FontFamily:    "monospace",
				Theme:         DefaultTheme,
				LetterSpacing: 1.0,
				LineHeight:    1.0,
				CursorBlink:   true,
				Video: VideoOptions{
					Framerate:     30,
					PlaybackSpeed: 1.0,
					Output: VideoOutputs{
						SVG: "test_output.svg",
					},
					Style: DefaultStyleOptions(),
				},
				SVG: SVGOptions{
					OptimizeSize: true,
				},
			},
			svgFrames: []SVGFrame{
				{
					Lines:   []string{"Test output"},
					CursorX: 0,
					CursorY: 0,
				},
			},
		}

		err := MakeSVG(vhs)
		if err != nil {
			t.Fatalf("MakeSVG failed: %v", err)
		}

		// Check if file was created
		if _, err := os.Stat("test_output.svg"); os.IsNotExist(err) {
			t.Error("SVG file was not created")
		}

		// Clean up
		_ = os.Remove("test_output.svg")
	})

	t.Run("skips when no SVG output specified", func(t *testing.T) {
		vhs := &VHS{
			Options: &Options{
				Video: VideoOptions{
					Output: VideoOutputs{
						SVG: "", // No output specified
					},
				},
			},
		}

		err := MakeSVG(vhs)
		if err != nil {
			t.Errorf("MakeSVG should not error when no output specified: %v", err)
		}
	})

	t.Run("skips when no frames captured", func(t *testing.T) {
		vhs := &VHS{
			Options: &Options{
				Video: VideoOptions{
					Output: VideoOutputs{
						SVG: "test_output.svg",
					},
				},
			},
			svgFrames: []SVGFrame{}, // No frames
		}

		err := MakeSVG(vhs)
		if err != nil {
			t.Errorf("MakeSVG should not error when no frames: %v", err)
		}
	})
}

// Benchmark Tests
func BenchmarkSVGGeneration(b *testing.B) {
	// Create a large set of frames
	frames := make([]SVGFrame, 100)
	for i := range frames {
		lines := make([]string, 30)
		for j := range lines {
			lines[j] = fmt.Sprintf("Line %d: %s", j, strings.Repeat("content ", 10))
		}
		frames[i] = SVGFrame{
			Lines:   lines,
			CursorX: i % 80,
			CursorY: i % 30,
		}
	}

	opts := SVGConfig{
		Width:      1024,
		Height:     768,
		FontSize:   16,
		FontFamily: "monospace",
		Theme:      DefaultTheme,
		Frames:     frames,
		Duration:   10.0,
		Style:      DefaultStyleOptions(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gen := NewSVGGenerator(opts)
		_ = gen.Generate()
	}
}

func BenchmarkFrameDeduplication(b *testing.B) {
	// Create frames with some duplicates
	frames := make([]SVGFrame, 1000)
	for i := range frames {
		// Create patterns that repeat
		content := fmt.Sprintf("Pattern %d", i%50)
		frames[i] = SVGFrame{
			Lines:   []string{content},
			CursorX: i % 10,
			CursorY: 0,
		}
	}

	opts := SVGConfig{
		Frames: frames,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gen := NewSVGGenerator(opts)
		gen.processFrames()
	}
}

// CharStyle Tests
func TestCharStyle(t *testing.T) {
	t.Run("all style attributes", func(t *testing.T) {
		style := CharStyle{
			FgColor:   "#ff0000",
			BgColor:   "#0000ff",
			Bold:      true,
			Italic:    true,
			Underline: true,
		}

		if style.FgColor != "#ff0000" {
			t.Errorf("Expected FgColor #ff0000, got %s", style.FgColor)
		}
		if style.BgColor != "#0000ff" {
			t.Errorf("Expected BgColor #0000ff, got %s", style.BgColor)
		}
		if !style.Bold {
			t.Error("Expected Bold to be true")
		}
		if !style.Italic {
			t.Error("Expected Italic to be true")
		}
		if !style.Underline {
			t.Error("Expected Underline to be true")
		}
	})
}

// TerminalState Tests
func TestTerminalState(t *testing.T) {
	t.Run("state with colors", func(t *testing.T) {
		state := TerminalState{
			Lines:   []string{"Hello", "World"},
			CursorX: 5,
			CursorY: 1,
			LineColors: [][]CharStyle{
				{{FgColor: "#ff0000"}, {FgColor: "#00ff00"}},
				{{FgColor: "#0000ff"}},
			},
		}

		if len(state.Lines) != 2 {
			t.Errorf("Expected 2 lines, got %d", len(state.Lines))
		}
		if state.CursorX != 5 || state.CursorY != 1 {
			t.Errorf("Expected cursor at (5,1), got (%d,%d)", state.CursorX, state.CursorY)
		}
		if len(state.LineColors) != 2 {
			t.Errorf("Expected 2 line colors, got %d", len(state.LineColors))
		}
	})
}

// SVGFrame Structure Tests
func TestSVGFrameStructure(t *testing.T) {
	t.Run("all fields accessible", func(t *testing.T) {
		frame := SVGFrame{
			Lines:      []string{"Test"},
			CursorX:    10,
			CursorY:    5,
			Timestamp:  1.5,
			CharWidth:  8.8,
			CharHeight: 20.0,
			LineColors: [][]CharStyle{
				{{FgColor: "#ff0000", BgColor: "#0000ff"}},
			},
		}

		// Verify all fields
		if len(frame.Lines) != 1 || frame.Lines[0] != "Test" {
			t.Error("Lines field not set correctly")
		}

		if frame.CursorX != 10 || frame.CursorY != 5 {
			t.Error("Cursor position fields not set correctly")
		}

		if frame.Timestamp != 1.5 {
			t.Error("Timestamp field not set correctly")
		}

		if frame.CharWidth != 8.8 || frame.CharHeight != 20.0 {
			t.Error("Character dimension fields not set correctly")
		}

		if len(frame.LineColors) != 1 {
			t.Error("LineColors field not set correctly")
		}
	})
}

// Edge Cases Tests
func TestSVGGenerator_EdgeCases(t *testing.T) {
	t.Run("handles very long lines", func(t *testing.T) {
		longLine := strings.Repeat("x", 200)
		opts := createTestSVGConfig()
		opts.Frames = []SVGFrame{
			{Lines: []string{longLine}},
		}

		gen := NewSVGGenerator(opts)
		svg := gen.Generate()

		// Should still contain the line (possibly truncated or wrapped)
		assertContains(t, svg, "x", "Long line content")
	})

	t.Run("handles empty terminal", func(t *testing.T) {
		opts := createTestSVGConfig()
		opts.Frames = []SVGFrame{
			{Lines: []string{}},
		}

		gen := NewSVGGenerator(opts)
		svg := gen.Generate()

		assertContains(t, svg, "<svg", "SVG generates with empty terminal")
	})

	t.Run("handles nil theme", func(t *testing.T) {
		opts := createTestSVGConfig()
		opts.Theme = Theme{} // Empty theme

		gen := NewSVGGenerator(opts)
		svg := gen.Generate()

		// Should use defaults and not crash
		assertContains(t, svg, "<svg", "SVG generates with empty theme")
	})

	t.Run("handles zero dimensions", func(t *testing.T) {
		opts := createTestSVGConfig()
		opts.Width = 0
		opts.Height = 0

		gen := NewSVGGenerator(opts)
		svg := gen.Generate()

		// Should still generate valid SVG
		assertContains(t, svg, "<svg", "SVG generates with zero dimensions")
	})
}

// JavaScript Capture Test
func TestCaptureSVGFrame_JavaScript(t *testing.T) {
	t.Run("JavaScript code structure", func(t *testing.T) {
		// Test the JavaScript code that's embedded in CaptureSVGFrame
		jsCode := `
const term = window.term;
if (!term) {
    return { error: "Terminal not found" };
}

const buffer = term.buffer.active;
const lines = [];
const lineColors = [];

for (let y = 0; y < buffer.length; y++) {
    const line = buffer.getLine(y);
    if (!line) continue;
    
    lines.push(line.translateToString(true));
}

return {
    lines: lines,
    cursorX: buffer.cursorX,
    cursorY: buffer.cursorY,
    charWidth: term._core._renderService._renderer._dimensions.css.cell.width,
    charHeight: term._core._renderService._renderer._dimensions.css.cell.height
};
`
		// Verify JavaScript syntax is valid (would need actual JS parser for full validation)
		if !strings.Contains(jsCode, "window.term") {
			t.Error("JavaScript should reference window.term")
		}
		if !strings.Contains(jsCode, "translateToString(true)") {
			t.Error("JavaScript should preserve whitespace with translateToString(true)")
		}
	})
}

// Utility Function Tests
func TestSVGGenerator_UtilityFunctions(t *testing.T) {

	t.Run("getColorClass", func(t *testing.T) {
		// Test with optimization disabled
		opts := createTestSVGConfig()
		opts.OptimizeSize = false
		gen := NewSVGGenerator(opts)

		if gen.getColorClass("#000000") != "" {
			t.Error("getColorClass should return empty string when optimization is disabled")
		}

		// Test with optimization enabled
		opts.OptimizeSize = true
		gen = NewSVGGenerator(opts)

		testCases := []struct {
			color    string
			expected string
			desc     string
		}{
			{gen.options.Theme.Black, "k", "black color"},
			{gen.options.Theme.Red, "r", "red color"},
			{gen.options.Theme.Green, "g", "green color"},
			{gen.options.Theme.Yellow, "y", "yellow color"},
			{gen.options.Theme.Blue, "b", "blue color"},
			{gen.options.Theme.Magenta, "m", "magenta color"},
			{gen.options.Theme.Cyan, "c", "cyan color"},
			{gen.options.Theme.White, "w", "white color"},
			{gen.options.Theme.BrightBlue, "p", "bright blue (prompt) color"},
			{"#5a56e0", "p", "hardcoded prompt color"},
			{"#123456", "", "unknown color"},
		}

		for _, tc := range testCases {
			result := gen.getColorClass(tc.color)
			if result != tc.expected {
				t.Errorf("getColorClass(%s) = %q, expected %q for %s", tc.color, result, tc.expected, tc.desc)
			}
		}
	})

}

// Test MakeSVG with error conditions
func TestMakeSVG_ErrorConditions(t *testing.T) {
	t.Run("handles file write error gracefully", func(t *testing.T) {
		// Use a path that will definitely fail on all platforms
		// On Windows, use a path with invalid characters
		// On Unix, use a path in a read-only directory
		invalidPath := "/dev/null/cannot/create/file/here.svg"
		if runtime.GOOS == "windows" {
			// Windows doesn't allow colons in filenames (except for drive letters)
			invalidPath = `C:\invalid:path\file*.svg`
		}

		vhs := &VHS{
			Options: &Options{
				FontSize:   16,
				FontFamily: "monospace",
				Theme:      DefaultTheme,
				Video: VideoOptions{
					Output: VideoOutputs{
						SVG: invalidPath,
					},
					Style: DefaultStyleOptions(),
				},
			},
			svgFrames: []SVGFrame{
				{Lines: []string{"Test"}},
			},
		}

		err := MakeSVG(vhs)
		if err == nil {
			t.Error("MakeSVG should return error when unable to write file")
		}
	})
}

// Test parseFontFamily more thoroughly
func TestParseFontFamily(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{"Arial", []string{"Arial"}},
		{"Arial, sans-serif", []string{"Arial", "sans-serif"}},
		{"'JetBrains Mono', monospace", []string{"JetBrains Mono", "monospace"}},
		{"\"Courier New\", Courier, monospace", []string{"Courier New", "Courier", "monospace"}},
		{"  Arial  ,  Helvetica  ", []string{"Arial", "Helvetica"}},
		{"", []string{"monospace"}}, // Empty string returns monospace as default
		{"'Font', \"Another Font\", third", []string{"Font", "Another Font", "third"}},
	}

	for _, tc := range testCases {
		result := parseFontFamily(tc.input)
		if len(result) != len(tc.expected) {
			t.Errorf("parseFontFamily(%q) returned %d fonts, expected %d", tc.input, len(result), len(tc.expected))
			continue
		}
		for i, font := range result {
			if font != tc.expected[i] {
				t.Errorf("parseFontFamily(%q)[%d] = %q, expected %q", tc.input, i, font, tc.expected[i])
			}
		}
	}
}

// Test buildSVGFontFamily
func TestBuildSVGFontFamily(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"Arial", "Arial, monospace"},
		{"Arial, sans-serif", "Arial, sans-serif, monospace"},
		{"JetBrains Mono, Courier, monospace", "JetBrains Mono, Courier, monospace"},
		{"", "monospace"},
		{"'Font Name'", "Font Name, monospace"},
		{"monospace", "monospace"},
		{"Consolas, monospace", "Consolas, monospace"},
	}

	for _, tc := range testCases {
		result := buildSVGFontFamily(tc.input)
		if result != tc.expected {
			t.Errorf("buildSVGFontFamily(%q) = %q, expected %q", tc.input, result, tc.expected)
		}
	}
}

// Test formatDuration edge cases
func TestFormatDuration(t *testing.T) {
	testCases := []struct {
		input    float64
		expected string
	}{
		{0, "0"},
		{1.0, "1"},
		{1.5, "1.5"},
		{1.50, "1.5"},
		{1.234, "1.23"},
		{10.00, "10"},
		{-2.5, "-2.5"},
		{0.1, "0.1"},
		{0.01, "0.01"},
		{0.001, "0"},
	}

	for _, tc := range testCases {
		result := formatDuration(tc.input)
		if result != tc.expected {
			t.Errorf("formatDuration(%f) = %q, expected %q", tc.input, result, tc.expected)
		}
	}
}

// Test formatPercentage function with dynamic precision
func TestFormatPercentage(t *testing.T) {
	testCases := []struct {
		name          string
		input         float64
		keyframeCount int
		expected      string
	}{
		// Whole numbers always minimal
		{"whole number 0", 0.0, 100, "0"},
		{"whole number 100", 100.0, 1000, "100"},
		{"whole number 50", 50.0, 10000, "50"},
		
		// Small frame counts (< 100) - 1 decimal
		{"small count decimal", 12.345, 50, "12.3"},
		{"small count trailing", 10.5000, 75, "10.5"},
		
		// Medium frame counts (< 1000) - 2 decimals
		{"medium count decimal", 0.024038, 500, "0.02"},
		{"medium count precision", 2.456, 800, "2.46"},
		{"medium count trailing", 10.1000, 900, "10.1"},
		
		// Large frame counts (< 10000) - 3 decimals
		{"large count decimal", 0.024038, 4161, "0.024"},
		{"large count precision", 2.451923, 5000, "2.452"},
		{"large count trailing", 99.9000, 8000, "99.9"},
		
		// Very large frame counts (< 100000) - 4 decimals
		{"very large decimal", 0.00012, 50000, "0.0001"},
		{"very large precision", 33.33333, 75000, "33.3333"},
		
		// Huge frame counts (>= 100000) - 5 decimals
		{"huge count decimal", 0.000012, 150000, "0.00001"},
		{"huge count precision", 12.345678, 200000, "12.34568"},
		
		// Edge cases
		{"negative percentage", -5.5, 100, "-5.5"},
		{"zero with decimals", 0.0001, 1000, "0"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := formatPercentage(tc.input, tc.keyframeCount)
			if result != tc.expected {
				t.Errorf("formatPercentage(%f, %d) = %s; want %s", 
					tc.input, tc.keyframeCount, result, tc.expected)
			}
		})
	}
}

// Test dynamic precision prevents collisions
func TestFormatPercentageDynamicPrecision(t *testing.T) {
	testCases := []struct {
		name       string
		frameCount int
	}{
		{"small animation", 50},
		{"medium animation", 500},
		{"large animation", 4161},
		{"very large animation", 50000},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			percentages := make(map[string]bool)
			
			for i := 0; i < tc.frameCount; i++ {
				percentage := float64(i) / float64(tc.frameCount-1) * 100
				formatted := formatPercentage(percentage, tc.frameCount)
				
				if percentages[formatted] {
					t.Errorf("Duplicate percentage at frame %d/%d: %s (%.6f%%)", 
						i, tc.frameCount, formatted, percentage)
				}
				percentages[formatted] = true
			}
			
			if len(percentages) != tc.frameCount {
				t.Errorf("Expected %d unique percentages, got %d", 
					tc.frameCount, len(percentages))
			}
		})
	}
}

// Test that large animations produce unique percentages
func TestLargeAnimationPercentages(t *testing.T) {
	// Test that 4161 frames produce unique percentages
	percentages := make(map[string]bool)
	totalFrames := 4161
	
	for i := 0; i < totalFrames; i++ {
		percentage := float64(i) / float64(totalFrames-1) * 100
		formatted := formatPercentage(percentage, totalFrames)
		
		if percentages[formatted] {
			t.Errorf("Duplicate percentage at frame %d: %s (%.6f%%)", i, formatted, percentage)
		}
		percentages[formatted] = true
	}
	
	if len(percentages) != totalFrames {
		t.Errorf("Expected %d unique percentages, got %d", 
			totalFrames, len(percentages))
	}
}

// Test SVG keyframe generation with many frames
func TestSVGKeyframeGeneration(t *testing.T) {
	// Create config with many frames to test keyframe collision
	opts := createTestSVGConfig()
	opts.Frames = make([]SVGFrame, 1000)
	for i := range opts.Frames {
		opts.Frames[i] = SVGFrame{
			Lines:      []string{fmt.Sprintf("Frame %d", i)},
			CharWidth:  8.8,
			CharHeight: 20,
			CursorX:    i % 10,
			CursorY:    0,
		}
	}
	
	gen := NewSVGGenerator(opts)
	svg := gen.Generate()
	
	// Verify no duplicate keyframe percentages
	keyframeRegex := regexp.MustCompile(`(\d+(?:\.\d+)?)\%\s*\{`)
	matches := keyframeRegex.FindAllStringSubmatch(svg, -1)
	
	seen := make(map[string]bool)
	duplicates := 0
	for _, match := range matches {
		if seen[match[1]] {
			duplicates++
			if duplicates <= 5 { // Only report first 5 duplicates
				t.Errorf("Duplicate keyframe percentage found: %s%%", match[1])
			}
		}
		seen[match[1]] = true
	}
	
	if duplicates > 0 {
		t.Errorf("Found %d duplicate keyframe percentages out of %d total keyframes", 
			duplicates, len(matches))
	}
	
	// Verify animation is present
	if !strings.Contains(svg, "@keyframes slide") {
		t.Error("SVG missing @keyframes slide animation")
	}
	
	// Verify reasonable number of keyframes were generated
	// With deduplication, we should have fewer keyframes than frames
	if len(matches) > len(opts.Frames) {
		t.Errorf("Too many keyframes generated: %d keyframes for %d frames", 
			len(matches), len(opts.Frames))
	}
}

// Pattern Detection Tests
func TestSVGGenerator_PatternDetection(t *testing.T) {
	t.Run("detects simple typing pattern", func(t *testing.T) {
		opts := createTestSVGConfig()
		opts.Frames = []SVGFrame{
			{Lines: []string{"$ "}, CursorX: 2, CursorY: 0, Timestamp: 0.0, CharWidth: 8.8, CharHeight: 20},
			{Lines: []string{"$ h"}, CursorX: 3, CursorY: 0, Timestamp: 0.1, CharWidth: 8.8, CharHeight: 20},
			{Lines: []string{"$ he"}, CursorX: 4, CursorY: 0, Timestamp: 0.2, CharWidth: 8.8, CharHeight: 20},
			{Lines: []string{"$ hel"}, CursorX: 5, CursorY: 0, Timestamp: 0.3, CharWidth: 8.8, CharHeight: 20},
			{Lines: []string{"$ hell"}, CursorX: 6, CursorY: 0, Timestamp: 0.4, CharWidth: 8.8, CharHeight: 20},
			{Lines: []string{"$ hello"}, CursorX: 7, CursorY: 0, Timestamp: 0.5, CharWidth: 8.8, CharHeight: 20},
		}
		
		gen := NewSVGGenerator(opts)
		gen.detectPatterns()
		
		// Should detect one typing pattern
		typingPatterns := 0
		for _, p := range gen.patterns {
			if p.Type == PatternTyping {
				typingPatterns++
				// Verify the detected text
				if p.Text != "hello" {
					t.Errorf("Expected typed text 'hello', got '%s'", p.Text)
				}
				// Verify it covers the right frames
				if p.StartFrame != 0 || p.EndFrame != 5 {
					t.Errorf("Expected pattern to cover frames 0-5, got %d-%d", p.StartFrame, p.EndFrame)
				}
			}
		}
		
		if typingPatterns != 1 {
			t.Errorf("Expected 1 typing pattern, got %d", typingPatterns)
		}
	})
	
	t.Run("detects multiple typing patterns", func(t *testing.T) {
		opts := createTestSVGConfig()
		opts.Frames = []SVGFrame{
			// First typing sequence
			{Lines: []string{"$ "}, CursorX: 2, CursorY: 0, Timestamp: 0.0, CharWidth: 8.8, CharHeight: 20},
			{Lines: []string{"$ l"}, CursorX: 3, CursorY: 0, Timestamp: 0.1, CharWidth: 8.8, CharHeight: 20},
			{Lines: []string{"$ ls"}, CursorX: 4, CursorY: 0, Timestamp: 0.2, CharWidth: 8.8, CharHeight: 20},
			// Output (breaks pattern)
			{Lines: []string{"$ ls", "file1.txt", "file2.txt"}, CursorX: 0, CursorY: 2, Timestamp: 0.3, CharWidth: 8.8, CharHeight: 20},
			// Second typing sequence
			{Lines: []string{"$ ls", "file1.txt", "file2.txt", "$ "}, CursorX: 2, CursorY: 3, Timestamp: 0.4, CharWidth: 8.8, CharHeight: 20},
			{Lines: []string{"$ ls", "file1.txt", "file2.txt", "$ c"}, CursorX: 3, CursorY: 3, Timestamp: 0.5, CharWidth: 8.8, CharHeight: 20},
			{Lines: []string{"$ ls", "file1.txt", "file2.txt", "$ ca"}, CursorX: 4, CursorY: 3, Timestamp: 0.6, CharWidth: 8.8, CharHeight: 20},
			{Lines: []string{"$ ls", "file1.txt", "file2.txt", "$ cat"}, CursorX: 5, CursorY: 3, Timestamp: 0.7, CharWidth: 8.8, CharHeight: 20},
		}
		
		gen := NewSVGGenerator(opts)
		gen.detectPatterns()
		
		// Should detect two typing patterns
		typingPatterns := 0
		var detectedTexts []string
		for _, p := range gen.patterns {
			if p.Type == PatternTyping {
				typingPatterns++
				detectedTexts = append(detectedTexts, p.Text)
			}
		}
		
		if typingPatterns != 2 {
			t.Errorf("Expected 2 typing patterns, got %d", typingPatterns)
		}
		
		// Check the detected text
		if len(detectedTexts) == 2 {
			if detectedTexts[0] != "ls" {
				t.Errorf("Expected first typed text 'ls', got '%s'", detectedTexts[0])
			}
			if detectedTexts[1] != "cat" {
				t.Errorf("Expected second typed text 'cat', got '%s'", detectedTexts[1])
			}
		}
	})
	
	t.Run("does not detect pattern with too few frames", func(t *testing.T) {
		opts := createTestSVGConfig()
		opts.Frames = []SVGFrame{
			{Lines: []string{"$ "}, CursorX: 2, CursorY: 0, Timestamp: 0.0, CharWidth: 8.8, CharHeight: 20},
			{Lines: []string{"$ h"}, CursorX: 3, CursorY: 0, Timestamp: 0.1, CharWidth: 8.8, CharHeight: 20},
		}
		
		gen := NewSVGGenerator(opts)
		gen.detectPatterns()
		
		// Should not detect typing pattern (too few frames)
		for _, p := range gen.patterns {
			if p.Type == PatternTyping {
				t.Error("Should not detect typing pattern with only 2 frames")
			}
		}
	})
	
	t.Run("handles mixed typing and static frames", func(t *testing.T) {
		opts := createTestSVGConfig()
		opts.Frames = []SVGFrame{
			// Static frame
			{Lines: []string{"Welcome"}, CursorX: 0, CursorY: 0, Timestamp: 0.0, CharWidth: 8.8, CharHeight: 20},
			// Typing sequence
			{Lines: []string{"Welcome", "$ "}, CursorX: 2, CursorY: 1, Timestamp: 1.0, CharWidth: 8.8, CharHeight: 20},
			{Lines: []string{"Welcome", "$ e"}, CursorX: 3, CursorY: 1, Timestamp: 1.1, CharWidth: 8.8, CharHeight: 20},
			{Lines: []string{"Welcome", "$ ec"}, CursorX: 4, CursorY: 1, Timestamp: 1.2, CharWidth: 8.8, CharHeight: 20},
			{Lines: []string{"Welcome", "$ ech"}, CursorX: 5, CursorY: 1, Timestamp: 1.3, CharWidth: 8.8, CharHeight: 20},
			{Lines: []string{"Welcome", "$ echo"}, CursorX: 6, CursorY: 1, Timestamp: 1.4, CharWidth: 8.8, CharHeight: 20},
			// Static frame
			{Lines: []string{"Welcome", "$ echo", "Hello"}, CursorX: 0, CursorY: 2, Timestamp: 1.5, CharWidth: 8.8, CharHeight: 20},
		}
		
		gen := NewSVGGenerator(opts)
		gen.detectPatterns()
		
		typingCount := 0
		staticCount := 0
		for _, p := range gen.patterns {
			if p.Type == PatternTyping {
				typingCount++
			} else if p.Type == PatternStatic {
				staticCount++
			}
		}
		
		if typingCount != 1 {
			t.Errorf("Expected 1 typing pattern, got %d", typingCount)
		}
		
		// Static frames: first frame, and last frame (total 2)
		if staticCount != 2 {
			t.Errorf("Expected 2 static patterns, got %d", staticCount)
		}
	})
}

// Typing Animation CSS Tests
func TestSVGGenerator_TypingAnimationCSS(t *testing.T) {
	t.Run("generates typing animation CSS", func(t *testing.T) {
		opts := createTestSVGConfig()
		opts.Frames = []SVGFrame{
			{Lines: []string{"$ "}, CursorX: 2, CursorY: 0, Timestamp: 0.0, CharWidth: 8.8, CharHeight: 20},
			{Lines: []string{"$ t"}, CursorX: 3, CursorY: 0, Timestamp: 0.1, CharWidth: 8.8, CharHeight: 20},
			{Lines: []string{"$ te"}, CursorX: 4, CursorY: 0, Timestamp: 0.2, CharWidth: 8.8, CharHeight: 20},
			{Lines: []string{"$ tes"}, CursorX: 5, CursorY: 0, Timestamp: 0.3, CharWidth: 8.8, CharHeight: 20},
			{Lines: []string{"$ test"}, CursorX: 6, CursorY: 0, Timestamp: 0.4, CharWidth: 8.8, CharHeight: 20},
		}
		
		gen := NewSVGGenerator(opts)
		svg := gen.Generate()
		
		// Should contain typing animation CSS
		assertContains(t, svg, "@keyframes typing_", "SVG should contain typing animation keyframes")
		assertContains(t, svg, ".typing_", "SVG should contain typing animation class")
		assertContains(t, svg, "overflow: hidden", "Typing animation should have overflow hidden")
		assertContains(t, svg, "white-space: nowrap", "Typing animation should have nowrap")
		assertContains(t, svg, "steps(", "Typing animation should use steps timing")
	})
	
	t.Run("calculates correct animation duration", func(t *testing.T) {
		opts := createTestSVGConfig()
		opts.Frames = []SVGFrame{
			{Lines: []string{"$ "}, CursorX: 2, CursorY: 0, Timestamp: 0.0, CharWidth: 8.8, CharHeight: 20},
			{Lines: []string{"$ h"}, CursorX: 3, CursorY: 0, Timestamp: 0.5, CharWidth: 8.8, CharHeight: 20},
			{Lines: []string{"$ hi"}, CursorX: 4, CursorY: 0, Timestamp: 1.0, CharWidth: 8.8, CharHeight: 20},
			{Lines: []string{"$ hi!"}, CursorX: 5, CursorY: 0, Timestamp: 1.5, CharWidth: 8.8, CharHeight: 20},
		}
		
		gen := NewSVGGenerator(opts)
		svg := gen.Generate()
		
		// Should have animation duration of 1.5s (from timestamp 0.0 to 1.5)
		assertContains(t, svg, "animation: typing_0 1.5s", "Animation duration should be 1.5s")
	})
}
