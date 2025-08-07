package main

import (
	"crypto/md5" //nolint:gosec // MD5 is used for deduplication, not security
	"fmt"
	"html"
	"log"
	"strings"

	"github.com/go-rod/rod"
)

// Default colors used throughout SVG generation.
const (
	defaultBackgroundColor = "#171717"
	defaultForegroundColor = "#dddddd"
	defaultCursorColor     = "#dddddd"
	defaultBarColor        = "#2d2d2d"
	defaultMarginColor     = "#000000"

	// Style constants.
	nilValue                = "<nil>"
	nullValue               = "null"
	fontWeightBold          = "font-weight:bold;"
	fontStyleItalic         = "font-style:italic;"
	textDecorationUnderline = "text-decoration:underline;"
	svgDefaultFontFamily    = "monospace"
)

// Window control colors (macOS-style).
var windowControlColors = []string{"#ff5f58", "#ffbd2e", "#18c132"}

// SVGFrame represents a single frame in the SVG animation.
type SVGFrame struct {
	Lines      []string
	LineColors [][]CharStyle // Color/style info for each character on each line
	CursorX    int
	CursorY    int
	Timestamp  float64
	CharWidth  float64
	CharHeight float64
	CursorChar string // The cursor character (e.g., '█' for block)
}

// CharStyle represents the style of a character.
type CharStyle struct {
	FgColor   string
	BgColor   string
	Bold      bool
	Italic    bool
	Underline bool
}

// SVGConfig contains the full configuration for SVG generation.
type SVGConfig struct {
	Width         int
	Height        int
	FontSize      int
	FontFamily    string
	Theme         Theme
	Frames        []SVGFrame
	Duration      float64
	Style         *StyleOptions // Include all style options
	LineHeight    float64
	CursorBlink   bool
	PlaybackSpeed float64
	LoopOffset    float64
	OptimizeSize  bool // Enable size optimizations for smaller output
	Debug         bool // Enable debug logging
}

// TerminalState represents a unique terminal state for deduplication.
type TerminalState struct {
	Lines          []string
	LineColors     [][]CharStyle // Color/style info for each character on each line
	CursorX        int
	CursorY        int
	Hash           string
	IsCursorActive bool    // Whether cursor moved from previous state
	CursorIdleTime float64 // Time since last cursor movement in seconds
	CursorChar     string  // The cursor character (e.g., '█' for block)
}

// KeyframeStop represents a point in the animation timeline.
type KeyframeStop struct {
	Percentage float64
	StateIndex int
}

// PatternType represents the type of change pattern detected.
type PatternType int

const (
	PatternStatic PatternType = iota
	PatternTyping
	PatternBackspace
)

// FramePattern represents a detected pattern in frame sequences.
type FramePattern struct {
	Type       PatternType
	StartFrame int
	EndFrame   int
	StartTime  float64
	EndTime    float64
	
	// For typing patterns
	Line     int
	StartCol int
	Text     string
	
	// For backspace patterns
	DeletedText  string // Text that was deleted
	DeletedCount int    // Number of characters deleted
	
	// Store the initial and final states
	InitialState TerminalState
	FinalState   TerminalState
}

// SVGGenerator handles the generation of optimized animated SVG files.
type SVGGenerator struct {
	options             SVGConfig
	charWidth           float64
	charHeight          float64
	fontSize            float64
	states              []TerminalState // Unique terminal states
	stateMap            map[string]int  // Hash -> state index
	timeline            []KeyframeStop  // Animation timeline
	patterns            []FramePattern  // Detected patterns for optimization
	frameSpacing        float64         // Spacing between frames in SVG units
	prevCursorX         int             // Previous cursor X position for activity detection
	prevCursorY         int             // Previous cursor Y position for activity detection
	cursorIdleThreshold float64         // Time threshold before cursor starts blinking (seconds)
	// Class names (shorter when OptimizeSize is enabled)
	textClass         string
	cursorActiveClass string
	cursorIdleClass   string
}

// NewSVGGenerator creates a new SVG generator.
func NewSVGGenerator(opts SVGConfig) *SVGGenerator {
	// Get character dimensions from the first frame if available
	charWidth := float64(opts.FontSize) * 0.55 // fallback
	charHeight := float64(opts.FontSize) * 1.2 // fallback

	if len(opts.Frames) > 0 && opts.Frames[0].CharWidth > 0 {
		// Use actual dimensions from xterm.js
		charWidth = opts.Frames[0].CharWidth
		charHeight = opts.Frames[0].CharHeight
	}

	// Get style for calculating frame spacing
	style := opts.Style
	if style == nil {
		style = DefaultStyleOptions()
	}

	// Frame spacing should match the inner terminal width
	innerWidth := style.Width - (style.Padding * 2)

	// Set class names based on optimization settings
	textClass := "f"
	cursorActiveClass := "cursor-active"
	cursorIdleClass := "cursor-idle"
	if opts.OptimizeSize {
		textClass = "t"
		cursorActiveClass = "ca"
		cursorIdleClass = "ci"
	}

	return &SVGGenerator{
		options:             opts,
		charWidth:           charWidth,
		charHeight:          charHeight,
		stateMap:            make(map[string]int),
		frameSpacing:        float64(innerWidth), // Frame spacing matches inner terminal width
		prevCursorX:         -1,                  // Initialize to -1 to detect first frame
		prevCursorY:         -1,
		cursorIdleThreshold: 0.5, // Default 0.5 seconds before cursor starts blinking
		textClass:           textClass,
		cursorActiveClass:   cursorActiveClass,
		cursorIdleClass:     cursorIdleClass,
	}
}

// Generate creates the complete SVG animation.
func (g *SVGGenerator) Generate() string {
	if g.options.Debug {
		log.Printf("SVG Generator Debug is enabled")
	}

	// Use style options for dimensions
	style := g.options.Style
	if style == nil {
		style = DefaultStyleOptions()
	}

	// Process frames to extract unique states
	g.processFrames()

	// Calculate fontSize early so it's available for symbol generation
	g.fontSize = float64(g.options.FontSize)
	if g.fontSize <= 0 {
		g.fontSize = 20
	}

	var sb strings.Builder

	// Calculate total dimensions including margins
	totalWidth := style.Width
	totalHeight := style.Height
	if style.Margin > 0 {
		totalWidth += style.Margin * 2
		totalHeight += style.Margin * 2
	}

	// SVG root element
	sb.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`,
		totalWidth, totalHeight))
	g.writeNewline(&sb)

	// Add margin group if needed
	if style.Margin > 0 {
		marginColor := style.MarginFill
		if marginColor == "" {
			marginColor = defaultMarginColor
		}
		sb.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="%s"/>`,
			totalWidth, totalHeight, marginColor))
		g.writeNewline(&sb)
		sb.WriteString(fmt.Sprintf(`<g transform="translate(%d,%d)">`, style.Margin, style.Margin))
		g.writeNewline(&sb)
	}

	// Terminal window
	sb.WriteString(g.generateTerminalWindow())

	// Calculate inner terminal area
	barHeight := 0
	if style.WindowBar != "" {
		barHeight = style.WindowBarSize
	}

	padding := style.Padding
	innerX := padding
	innerY := barHeight + padding
	innerWidth := style.Width - (padding * 2)
	innerHeight := style.Height - barHeight - (padding * 2)

	// Inner terminal SVG with viewBox for animation
	// Calculate actual terminal content height
	maxLines := 0
	for _, state := range g.states {
		if len(state.Lines) > maxLines {
			maxLines = len(state.Lines)
		}
	}
	// viewBox width should match frame spacing (one frame width), height matches terminal
	viewBoxWidth := g.frameSpacing
	viewBoxHeight := float64(innerHeight)

	// Create inner SVG with viewBox that shows one frame at a time
	sb.WriteString(fmt.Sprintf(`<svg x="%d" y="%d" width="%d" height="%d" viewBox="0 0 %s %s">`,
		innerX, innerY, innerWidth, innerHeight, formatCoord(viewBoxWidth), formatCoord(viewBoxHeight)))
	g.writeNewline(&sb)

	// Add terminal background
	terminalBgColor := g.options.Theme.Background
	if terminalBgColor == "" {
		terminalBgColor = defaultMarginColor
	}
	sb.WriteString(fmt.Sprintf(`<rect width="%s" height="%s" fill="%s"/>`,
		formatCoord(viewBoxWidth), formatCoord(viewBoxHeight), terminalBgColor))
	g.writeNewline(&sb)

	// Add styles including CSS animation
	sb.WriteString(g.generateStyles())

	// Add defs section for reusable elements
	sb.WriteString("<defs>")
	g.writeNewline(&sb)
	sb.WriteString(g.generateCursorSymbols())
	sb.WriteString("</defs>")
	g.writeNewline(&sb)

	// Animation container without additional clipping (viewBox handles it)
	sb.WriteString(`<g class="animation-container">`)
	g.writeNewline(&sb)

	// Generate all unique states
	for i, state := range g.states {
		if g.options.Debug {
			// Count background colors in this state
			bgCount := 0
			for _, lineColors := range state.LineColors {
				for _, style := range lineColors {
					if style.BgColor != "" && style.BgColor != nilValue && style.BgColor != nullValue {
						bgCount++
					}
				}
			}
			if bgCount > 0 {
				log.Printf("Generating state %d with %d background colors", i, bgCount)
			}
		}
		sb.WriteString(g.generateState(i, &state))
	}

	sb.WriteString("</g>") // Close animation container
	g.writeNewline(&sb)
	sb.WriteString("</svg>") // Close inner SVG
	g.writeNewline(&sb)

	// Close margin group if opened
	if style.Margin > 0 {
		sb.WriteString("</g>")
		g.writeNewline(&sb)
	}

	sb.WriteString("</svg>")
	g.writeNewline(&sb)

	return sb.String()
}

// processFrames deduplicates frames and builds timeline.
func (g *SVGGenerator) processFrames() {
	// First, detect patterns for optimization
	g.detectPatterns()
	
	// First pass: collect all unique states and track when they change
	lastStateIndex := -1
	lastCursorIdleTime := 0.0

	// Debug: Check for frames with background colors
	if g.options.Debug {
		for i, frame := range g.options.Frames {
			bgCount := 0
			for _, lineColors := range frame.LineColors {
				for _, style := range lineColors {
					if style.BgColor != "" && style.BgColor != nilValue && style.BgColor != nullValue {
						bgCount++
					}
				}
			}
			if bgCount > 0 {
				log.Printf("Frame %d has %d background colors", i, bgCount)
			}
		}
	}

	for i, frame := range g.options.Frames {
		// Create state from frame
		state := TerminalState{
			Lines:      frame.Lines,
			LineColors: frame.LineColors,
			CursorX:    frame.CursorX,
			CursorY:    frame.CursorY,
			CursorChar: frame.CursorChar,
		}

		// Detect cursor activity
		cursorMoved := false
		if g.prevCursorX != -1 && g.prevCursorY != -1 {
			// Check if cursor position changed
			if frame.CursorX != g.prevCursorX || frame.CursorY != g.prevCursorY {
				cursorMoved = true
			} else if i > 0 {
				// Check if text changed at cursor position (typing at same position)
				prevFrame := g.options.Frames[i-1]
				if frame.CursorY < len(frame.Lines) && frame.CursorY < len(prevFrame.Lines) {
					if frame.Lines[frame.CursorY] != prevFrame.Lines[frame.CursorY] {
						cursorMoved = true
					}
				}
			}
		}

		// Calculate cursor idle time
		if cursorMoved {
			state.IsCursorActive = true
			state.CursorIdleTime = 0.0
			lastCursorIdleTime = 0.0
		} else {
			// Calculate time difference from previous frame
			timeDiff := 0.0
			if i > 0 {
				timeDiff = frame.Timestamp - g.options.Frames[i-1].Timestamp
			}
			lastCursorIdleTime += timeDiff
			state.CursorIdleTime = lastCursorIdleTime
			state.IsCursorActive = lastCursorIdleTime < g.cursorIdleThreshold
		}

		// Update previous cursor position
		g.prevCursorX = frame.CursorX
		g.prevCursorY = frame.CursorY

		// Generate hash for deduplication
		hash := g.hashState(&state)
		state.Hash = hash

		// Check if we've seen this state before
		if idx, exists := g.stateMap[hash]; exists {
			// Reuse existing state - only add to timeline if state changed
			if idx != lastStateIndex {
				g.timeline = append(g.timeline, KeyframeStop{
					Percentage: float64(i) / float64(len(g.options.Frames)-1) * 100,
					StateIndex: idx,
				})
				lastStateIndex = idx
			}
		} else {
			// New unique state
			idx := len(g.states)
			g.states = append(g.states, state)
			g.stateMap[hash] = idx

			// Debug: Check if this state has background colors
			if g.options.Debug {
				bgCount := 0
				for _, lineColors := range state.LineColors {
					for _, style := range lineColors {
						if style.BgColor != "" && style.BgColor != nilValue && style.BgColor != nullValue {
							bgCount++
						}
					}
				}
				if bgCount > 0 {
					log.Printf("State %d (from frame %d) has %d background colors", idx, i, bgCount)
				}
			}

			g.timeline = append(g.timeline, KeyframeStop{
				Percentage: float64(i) / float64(len(g.options.Frames)-1) * 100,
				StateIndex: idx,
			})
			lastStateIndex = idx
		}
	}

	// Ensure we have the final frame at 100%
	if len(g.timeline) > 0 && g.timeline[len(g.timeline)-1].Percentage < 100.0 {
		// Add final keyframe at 100%
		g.timeline = append(g.timeline, KeyframeStop{
			Percentage: 100.0,
			StateIndex: lastStateIndex,
		})
	}

	// Log deduplication analysis
	if g.options.OptimizeSize {
		totalFrames := len(g.options.Frames)
		uniqueStates := len(g.states)
		duplicateFrames := totalFrames - uniqueStates
		deduplicationRate := float64(duplicateFrames) / float64(totalFrames) * 100

		if g.options.Debug {
			log.Printf("Frame deduplication analysis:")
			log.Printf("  Total frames: %d", totalFrames)
			log.Printf("  Unique states: %d", uniqueStates)
			log.Printf("  Duplicate frames: %d", duplicateFrames)
			log.Printf("  Deduplication rate: %.1f%%", deduplicationRate)
			log.Printf("  Timeline keyframes: %d", len(g.timeline))
		}

		// Analyze consecutive duplicate frames
		consecutiveDupes := 0
		maxConsecutive := 0
		currentConsecutive := 0
		prevHash := ""

		for i, frame := range g.options.Frames {
			state := TerminalState{
				Lines:      frame.Lines,
				LineColors: frame.LineColors,
				CursorX:    frame.CursorX,
				CursorY:    frame.CursorY,
				CursorChar: frame.CursorChar,
			}
			hash := g.hashState(&state)

			if i > 0 && hash == prevHash {
				currentConsecutive++
				consecutiveDupes++
				if currentConsecutive > maxConsecutive {
					maxConsecutive = currentConsecutive
				}
			} else {
				currentConsecutive = 0
			}
			prevHash = hash
		}

		if g.options.Debug {
			log.Printf("  Consecutive duplicate frames: %d", consecutiveDupes)
			log.Printf("  Max consecutive duplicates: %d", maxConsecutive)
		}
	}
}

// hashState generates a hash for a terminal state.
func (g *SVGGenerator) hashState(state *TerminalState) string {
	h := md5.New() //nolint:gosec // MD5 is used for deduplication, not security
	for i, line := range state.Lines {
		// Include the full line with trailing spaces in hash
		// This ensures states with different trailing spaces are treated as different
		h.Write([]byte(line))
		h.Write([]byte("\n"))

		// Include color information in hash
		if i < len(state.LineColors) {
			for _, style := range state.LineColors[i] {
				_, _ = fmt.Fprintf(h, "%s,%s,%t,%t,%t|",
					style.FgColor, style.BgColor, style.Bold, style.Italic, style.Underline)
			}
		}
		h.Write([]byte("\n"))
	}
	// Include cursor position and activity state in hash for accuracy
	// Only include active/idle status, not exact idle time to allow deduplication
	_, _ = fmt.Fprintf(h, "%d,%d,%v",
		state.CursorX, state.CursorY, state.IsCursorActive)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// detectPatterns analyzes frames to find typing and other patterns.
func (g *SVGGenerator) detectPatterns() {
	g.patterns = []FramePattern{}
	
	if len(g.options.Frames) < 2 {
		// Not enough frames to detect patterns
		return
	}
	
	i := 0
	for i < len(g.options.Frames) {
		// Try to detect typing pattern
		if pattern, consumed := g.detectTypingPattern(i); pattern != nil {
			g.patterns = append(g.patterns, *pattern)
			i += consumed
			continue
		}
		
		// Try to detect backspace pattern
		if pattern, consumed := g.detectBackspacePattern(i); pattern != nil {
			g.patterns = append(g.patterns, *pattern)
			i += consumed
			continue
		}
		
		// If no pattern detected, treat as static frame
		frame := g.options.Frames[i]
		g.patterns = append(g.patterns, FramePattern{
			Type:       PatternStatic,
			StartFrame: i,
			EndFrame:   i,
			StartTime:  frame.Timestamp,
			EndTime:    frame.Timestamp,
			FinalState: TerminalState{
				Lines:      frame.Lines,
				LineColors: frame.LineColors,
				CursorX:    frame.CursorX,
				CursorY:    frame.CursorY,
				CursorChar: frame.CursorChar,
			},
		})
		i++
	}
	
	if g.options.Debug {
		typingPatterns := 0
		typingFrames := 0
		backspacePatterns := 0
		backspaceFrames := 0
		for _, p := range g.patterns {
			switch p.Type {
			case PatternTyping:
				typingPatterns++
				typingFrames += p.EndFrame - p.StartFrame + 1
			case PatternBackspace:
				backspacePatterns++
				backspaceFrames += p.EndFrame - p.StartFrame + 1
			}
		}
		log.Printf("Pattern detection analysis:")
		log.Printf("  Total frames: %d", len(g.options.Frames))
		log.Printf("  Detected patterns: %d", len(g.patterns))
		log.Printf("  Typing patterns: %d (frames: %d)", typingPatterns, typingFrames)
		log.Printf("  Backspace patterns: %d (frames: %d)", backspacePatterns, backspaceFrames)
		totalOptimized := typingFrames + backspaceFrames
		log.Printf("  Total optimized frames: %d (%.1f%%)",
			totalOptimized, float64(totalOptimized)/float64(len(g.options.Frames))*100)
	}
}

// detectTypingPattern looks for consecutive frames where text is being typed on the same line.
func (g *SVGGenerator) detectTypingPattern(start int) (*FramePattern, int) {
	if start >= len(g.options.Frames)-1 {
		return nil, 0
	}
	
	firstFrame := g.options.Frames[start]
	line := firstFrame.CursorY
	startCol := firstFrame.CursorX
	
	// Track the typing sequence
	end := start + 1
	for end < len(g.options.Frames) {
		prev := g.options.Frames[end-1]
		curr := g.options.Frames[end]
		
		// Check if still typing on the same line
		if curr.CursorY != line {
			break
		}
		
		// Cursor should move forward (or stay for multi-byte chars)
		if curr.CursorX < prev.CursorX-1 { // Allow small backward movement for corrections
			break
		}
		
		// Check that only the cursor line changed
		if !g.isOnlyLineChanged(prev, curr, line) {
			break
		}
		
		// Line should grow (characters added)
		if line < len(prev.Lines) && line < len(curr.Lines) {
			prevLine := prev.Lines[line]
			currLine := curr.Lines[line]
			
			// Check if current line starts with previous line (typing appends)
			if !strings.HasPrefix(currLine, prevLine) {
				// If text got shorter, it's likely a backspace - break the pattern
				if len(currLine) < len(prevLine) {
					break
				}
				// If text changed but didn't grow from the previous, break
				break
			}
			
			// Check typing speed is reasonable (1-15 chars per frame is typical)
			charsChanged := abs(len(currLine) - len(prevLine))
			if charsChanged > 15 {
				break
			}
		} else {
			break
		}
		
		end++
	}
	
	// Need at least 3 frames to consider it a typing pattern
	framesInPattern := end - start
	if framesInPattern < 3 {
		return nil, 0
	}
	
	// Extract the typed text
	lastFrame := g.options.Frames[end-1]
	var typedText string
	
	if line < len(firstFrame.Lines) && line < len(lastFrame.Lines) {
		startLine := firstFrame.Lines[line]
		endLine := lastFrame.Lines[line]
		
		// Find the common prefix (unchanged part)
		commonPrefix := 0
		for i := 0; i < len(startLine) && i < len(endLine); i++ {
			if startLine[i] != endLine[i] {
				break
			}
			commonPrefix = i
		}
		
		// The typed text is what was added after the common prefix
		if len(endLine) > len(startLine) {
			typedText = endLine[len(startLine):]
		} else if commonPrefix < len(endLine) {
			// Handle case where text was modified, not just appended
			typedText = endLine[commonPrefix:]
		}
	}
	
	// Only create pattern if we actually typed something substantial
	if len(typedText) < 2 {
		return nil, 0
	}
	
	// Create initial and final states
	initialState := TerminalState{
		Lines:      firstFrame.Lines,
		LineColors: firstFrame.LineColors,
		CursorX:    firstFrame.CursorX,
		CursorY:    firstFrame.CursorY,
		CursorChar: firstFrame.CursorChar,
	}
	
	finalState := TerminalState{
		Lines:      lastFrame.Lines,
		LineColors: lastFrame.LineColors,
		CursorX:    lastFrame.CursorX,
		CursorY:    lastFrame.CursorY,
		CursorChar: lastFrame.CursorChar,
	}
	
	pattern := &FramePattern{
		Type:         PatternTyping,
		StartFrame:   start,
		EndFrame:     end - 1,
		StartTime:    firstFrame.Timestamp,
		EndTime:      lastFrame.Timestamp,
		Line:         line,
		StartCol:     startCol,
		Text:         typedText,
		InitialState: initialState,
		FinalState:   finalState,
	}
	
	if g.options.Debug {
		log.Printf("Detected typing pattern: frames %d-%d, line %d, text: %q (saved %d frames)",
			start, end-1, line, typedText, framesInPattern-2)
	}
	
	return pattern, framesInPattern
}

// detectBackspacePattern looks for consecutive frames where text is being deleted.
func (g *SVGGenerator) detectBackspacePattern(start int) (*FramePattern, int) {
	if start >= len(g.options.Frames)-1 {
		return nil, 0
	}
	
	firstFrame := g.options.Frames[start]
	line := firstFrame.CursorY
	
	// Track the backspace sequence
	end := start + 1
	totalDeleted := 0
	
	for end < len(g.options.Frames) {
		prev := g.options.Frames[end-1]
		curr := g.options.Frames[end]
		
		// Check if still on the same line
		if curr.CursorY != line {
			break
		}
		
		// Check that only the cursor line changed
		if !g.isOnlyLineChanged(prev, curr, line) {
			break
		}
		
		// Check if text is getting shorter (backspace pattern)
		if line < len(prev.Lines) && line < len(curr.Lines) {
			prevLine := prev.Lines[line]
			currLine := curr.Lines[line]
			
			// For backspace, current line should be shorter
			if len(currLine) >= len(prevLine) {
				break
			}
			
			// Check if it's a prefix (deleting from end)
			if !strings.HasPrefix(prevLine, currLine) {
				// Could be deletion in middle, but for now we'll break
				break
			}
			
			// Track how many characters were deleted
			deleted := len(prevLine) - len(currLine)
			totalDeleted += deleted
			
			// Don't group huge deletions (likely line clear, not backspace)
			if deleted > 10 {
				break
			}
		} else {
			break
		}
		
		end++
	}
	
	// Need at least 2 frames to consider it a backspace pattern
	framesInPattern := end - start
	if framesInPattern < 2 {
		return nil, 0
	}
	
	// Need to have deleted at least 2 characters to be worth optimizing
	if totalDeleted < 2 {
		return nil, 0
	}
	
	// Extract what was deleted
	lastFrame := g.options.Frames[end-1]
	var deletedText string
	
	if line < len(firstFrame.Lines) && line < len(lastFrame.Lines) {
		startLine := firstFrame.Lines[line]
		endLine := lastFrame.Lines[line]
		
		if strings.HasPrefix(startLine, endLine) {
			deletedText = startLine[len(endLine):]
		}
	}
	
	// Create states
	initialState := TerminalState{
		Lines:      firstFrame.Lines,
		LineColors: firstFrame.LineColors,
		CursorX:    firstFrame.CursorX,
		CursorY:    firstFrame.CursorY,
		CursorChar: firstFrame.CursorChar,
	}
	
	finalState := TerminalState{
		Lines:      lastFrame.Lines,
		LineColors: lastFrame.LineColors,
		CursorX:    lastFrame.CursorX,
		CursorY:    lastFrame.CursorY,
		CursorChar: lastFrame.CursorChar,
	}
	
	pattern := &FramePattern{
		Type:         PatternBackspace,
		StartFrame:   start,
		EndFrame:     end - 1,
		StartTime:    firstFrame.Timestamp,
		EndTime:      lastFrame.Timestamp,
		Line:         line,
		DeletedText:  deletedText,
		DeletedCount: totalDeleted,
		InitialState: initialState,
		FinalState:   finalState,
	}
	
	if g.options.Debug {
		log.Printf("Detected backspace pattern: frames %d-%d, line %d, deleted: %q (saved %d frames)",
			start, end-1, line, deletedText, framesInPattern-1)
	}
	
	return pattern, framesInPattern
}

// isOnlyLineChanged checks if only the specified line changed between frames.
func (g *SVGGenerator) isOnlyLineChanged(prev, curr SVGFrame, targetLine int) bool {
	// Check if number of lines changed significantly
	if abs(len(curr.Lines)-len(prev.Lines)) > 1 {
		return false
	}
	
	// Check each line
	maxLines := len(prev.Lines)
	if len(curr.Lines) < maxLines {
		maxLines = len(curr.Lines)
	}
	
	for i := 0; i < maxLines; i++ {
		if i != targetLine {
			// Other lines should remain unchanged
			if prev.Lines[i] != curr.Lines[i] {
				return false
			}
		}
	}
	
	return true
}

// abs returns the absolute value of an integer.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// generateTypingCSS generates CSS animation for a typing pattern.
func (g *SVGGenerator) generateTypingCSS(sb *strings.Builder, index int, pattern FramePattern) {
	// Calculate the width of the typed text
	textWidth := float64(len(pattern.Text)) * g.charWidth
	duration := pattern.EndTime - pattern.StartTime
	
	// Generate the keyframe animation
	sb.WriteString(fmt.Sprintf("@keyframes typing_%d {", index))
	g.writeNewline(sb)
	sb.WriteString("  from { width: 0; }")
	g.writeNewline(sb)
	sb.WriteString(fmt.Sprintf("  to { width: %spx; }", formatCoord(textWidth)))
	g.writeNewline(sb)
	sb.WriteString("}")
	g.writeNewline(sb)
	
	// Generate the class for this typing animation
	sb.WriteString(fmt.Sprintf(".typing_%d {", index))
	g.writeNewline(sb)
	sb.WriteString("  overflow: hidden;")
	g.writeNewline(sb)
	sb.WriteString("  white-space: nowrap;")
	g.writeNewline(sb)
	sb.WriteString("  display: inline-block;")
	g.writeNewline(sb)
	sb.WriteString(fmt.Sprintf("  animation: typing_%d %ss steps(%d, end) forwards;",
		index, formatDuration(duration), len(pattern.Text)))
	g.writeNewline(sb)
	sb.WriteString(fmt.Sprintf("  animation-delay: %ss;", formatDuration(pattern.StartTime)))
	g.writeNewline(sb)
	sb.WriteString("}")
	g.writeNewline(sb)
	g.writeNewline(sb)
}

// generateBackspaceCSS generates CSS animation for a backspace pattern.
func (g *SVGGenerator) generateBackspaceCSS(sb *strings.Builder, index int, pattern FramePattern) {
	// Calculate the width of the deleted text
	startWidth := float64(len(pattern.DeletedText)) * g.charWidth
	duration := pattern.EndTime - pattern.StartTime
	
	// Generate the keyframe animation (reverse of typing)
	sb.WriteString(fmt.Sprintf("@keyframes backspace_%d {", index))
	g.writeNewline(sb)
	sb.WriteString(fmt.Sprintf("  from { width: %spx; }", formatCoord(startWidth)))
	g.writeNewline(sb)
	sb.WriteString("  to { width: 0; }")
	g.writeNewline(sb)
	sb.WriteString("}")
	g.writeNewline(sb)
	
	// Generate the class for this backspace animation
	sb.WriteString(fmt.Sprintf(".backspace_%d {", index))
	g.writeNewline(sb)
	sb.WriteString("  overflow: hidden;")
	g.writeNewline(sb)
	sb.WriteString("  white-space: nowrap;")
	g.writeNewline(sb)
	sb.WriteString("  display: inline-block;")
	g.writeNewline(sb)
	sb.WriteString(fmt.Sprintf("  animation: backspace_%d %ss steps(%d, end) forwards;",
		index, formatDuration(duration), pattern.DeletedCount))
	g.writeNewline(sb)
	sb.WriteString(fmt.Sprintf("  animation-delay: %ss;", formatDuration(pattern.StartTime)))
	g.writeNewline(sb)
	sb.WriteString("}")
	g.writeNewline(sb)
	g.writeNewline(sb)
}

// generateStyles creates the CSS styles and animations.
func (g *SVGGenerator) generateStyles() string {
	var sb strings.Builder

	sb.WriteString("<style>")
	g.writeNewline(&sb)
	
	// Generate typing animations for detected patterns
	for i, pattern := range g.patterns {
		switch pattern.Type {
		case PatternTyping:
			g.generateTypingCSS(&sb, i, pattern)
		case PatternBackspace:
			g.generateBackspaceCSS(&sb, i, pattern)
		}
	}

	// Generate keyframes
	sb.WriteString("@keyframes slide {")
	g.writeNewline(&sb)

	// Build optimized keyframes from timeline
	keyframeCount := len(g.timeline)
	for _, stop := range g.timeline {
		offset := -float64(stop.StateIndex) * g.frameSpacing
		sb.WriteString(fmt.Sprintf("  %s%% { transform: translateX(%spx); }",
			formatPercentage(stop.Percentage, keyframeCount), formatCoord(offset)))
		g.writeNewline(&sb)
	}

	sb.WriteString("}")
	g.writeNewline(&sb)
	g.writeNewline(&sb)

	// Animation container style
	sb.WriteString(".animation-container {")
	g.writeNewline(&sb)
	// Apply playback speed to animation duration
	animationDuration := g.options.Duration
	if g.options.PlaybackSpeed > 0 {
		animationDuration = g.options.Duration / g.options.PlaybackSpeed
	}

	// Calculate animation delay based on LoopOffset
	animationDelay := 0.0
	if g.options.LoopOffset > 0 {
		// LoopOffset can be a percentage (0-100) or frame number
		if g.options.LoopOffset <= 1.0 {
			// Treat as percentage
			animationDelay = -animationDuration * g.options.LoopOffset
		} else {
			// Treat as frame number
			animationDelay = -(g.options.LoopOffset / float64(len(g.options.Frames))) * animationDuration
		}
	}

	// Use step-end timing to ensure frames change instantly
	sb.WriteString(fmt.Sprintf("  animation: slide %ss step-end %ss infinite;", formatDuration(animationDuration), formatDuration(animationDelay)))
	g.writeNewline(&sb)
	sb.WriteString("}")
	g.writeNewline(&sb)
	g.writeNewline(&sb)

	// Terminal styles
	theme := g.options.Theme

	// Text styles using classes for deduplication
	// Use shorter class names when optimization is enabled
	textClass := "f"
	cursorActiveClass := "cursor-active"
	cursorIdleClass := "cursor-idle"
	if g.options.OptimizeSize {
		textClass = "t"
		cursorActiveClass = "ca"
		cursorIdleClass = "ci"
	}

	// Ensure font-family is properly formatted for SVG
	fontFamily := g.options.FontFamily
	if fontFamily == "" {
		fontFamily = svgDefaultFontFamily
	}
	// Ensure foreground color is set
	foregroundColor := theme.Foreground
	if foregroundColor == "" {
		foregroundColor = defaultForegroundColor
	}
	// Use a simpler font stack for better compatibility
	textStyle := fmt.Sprintf("fill: %s; font-family: %s, monospace; font-size: %spx;",
		foregroundColor, fontFamily, formatCoord(g.fontSize))
	// Don't apply letter-spacing in SVG as it causes cursor misalignment
	// The character positions from xterm.js already account for the terminal's letter spacing
	sb.WriteString(fmt.Sprintf(".%s { %s }", textClass, textStyle))
	g.writeNewline(&sb)

	// Add ANSI color classes - use a map to avoid duplicates
	colorClasses := map[string]string{
		"black":   theme.Black,
		"red":     theme.Red,
		"green":   theme.Green,
		"yellow":  theme.Yellow,
		"blue":    theme.Blue,
		"magenta": theme.Magenta,
		"cyan":    theme.Cyan,
		"white":   theme.White,
	}

	if g.options.OptimizeSize {
		// Use single-letter class names for colors
		shortColorClasses := map[string]string{
			"k": theme.Black, // blacK
			"r": theme.Red,
			"g": theme.Green,
			"y": theme.Yellow,
			"b": theme.Blue,
			"m": theme.Magenta,
			"c": theme.Cyan,
			"w": theme.White,
		}
		for name, color := range shortColorClasses {
			sb.WriteString(fmt.Sprintf(".%s { fill: %s; }", name, color))
			g.writeNewline(&sb)
		}
		// Add prompt color class if we detect it's used frequently
		if theme.BrightBlue != "" {
			sb.WriteString(fmt.Sprintf(".p { fill: %s; }", theme.BrightBlue)) // prompt color
			g.writeNewline(&sb)
		}
	} else {
		for name, color := range colorClasses {
			sb.WriteString(fmt.Sprintf(".%s { fill: %s; }", name, color))
			g.writeNewline(&sb)
		}
	}

	// Cursor styles - for inline cursor with background
	// Note: SVG doesn't support background property on tspan, we'll need to use a different approach
	// We'll render a rect behind the cursor character
	// Active cursor is always visible
	sb.WriteString(fmt.Sprintf(".%s { }", cursorActiveClass))
	g.writeNewline(&sb)
	// Idle cursor blinks
	if g.options.CursorBlink {
		sb.WriteString("@keyframes blink { 0%, 49% { opacity: 1; } 50%, 100% { opacity: 0; } }")
		g.writeNewline(&sb)
		sb.WriteString(fmt.Sprintf(".%s { animation: blink 1s infinite; }", cursorIdleClass))
		g.writeNewline(&sb)
	}

	sb.WriteString("</style>")
	g.writeNewline(&sb)

	return sb.String()
}

// generateState creates a group for a single terminal state.
func (g *SVGGenerator) generateState(index int, state *TerminalState) string {
	var sb strings.Builder

	// Position this state in the animation sequence
	xOffset := float64(index) * g.frameSpacing

	sb.WriteString(fmt.Sprintf(`<g transform="translate(%s,0)">`, formatCoord(xOffset)))
	g.writeNewline(&sb)

	// Debug specific state with background colors
	if g.options.Debug && index == 19 {
		log.Printf("=== Generating state 19 with background colors ===")
		log.Printf("State has %d lines", len(state.Lines))
		for y, line := range state.Lines {
			if y < len(state.LineColors) && len(state.LineColors[y]) > 0 {
				bgCount := 0
				for x, style := range state.LineColors[y] {
					if style.BgColor != "" && style.BgColor != nilValue && style.BgColor != nullValue {
						bgCount++
						if bgCount == 1 {
							log.Printf("Line %d: first bg at x=%d color=%s", y, x, style.BgColor)
						}
					}
				}
				if bgCount > 0 {
					log.Printf("Line %d has %d bg colors, text=%q", y, bgCount, line)
				}
			}
		}
	}

	// Render lines with optimization
	for y, line := range state.Lines {
		// Check if cursor is on this line
		isCursorLine := y == state.CursorY

		// Check if we have color information for this line
		hasColors := y < len(state.LineColors) && len(state.LineColors[y]) > 0

		// Check if this line has background colors even if text is empty/spaces
		hasBackgroundColors := false
		if hasColors {
			for _, style := range state.LineColors[y] {
				if style.BgColor != "" && style.BgColor != nilValue && style.BgColor != nullValue {
					hasBackgroundColors = true
					if g.options.Debug && (y == 12 || y == 13) {
						log.Printf("Line %d has background color at hasBackgroundColors check", y)
					}
					break
				}
			}
		}

		// Skip empty lines ONLY if they don't have cursor or background colors
		if line == "" && y != state.CursorY && !hasBackgroundColors {
			continue
		}

		// Debug log lines with potential background colors
		if g.options.Debug && y < len(state.LineColors) {
			bgCount := 0
			for _, style := range state.LineColors[y] {
				if style.BgColor != "" && style.BgColor != nilValue && style.BgColor != nullValue {
					bgCount++
				}
			}
			if bgCount > 0 {
				log.Printf("Line %d has %d background colors, text: %q", y, bgCount, line)
			}
		}

		// Ultra-optimized rendering: use tspan for efficient text grouping
		// Unified rendering approach for both colored and non-colored text
		// This ensures consistent text selection and layout
		// Render if line has content, is cursor line, or has background colors
		if strings.TrimSpace(line) != "" || isCursorLine || hasBackgroundColors {
			// Render with colors using natural text flow
			lineHeight := g.options.LineHeight
			if lineHeight <= 0 {
				lineHeight = 1.0
			}
			yPos := float64(y)*g.charHeight*lineHeight + g.charHeight*0.8

			// First, render any background rectangles if we have color data
			if hasColors && y < len(state.LineColors) {
				if g.options.Debug {
					log.Printf("Processing line %d with %d color entries, line length: %d", y, len(state.LineColors[y]), len(line))
				}
				// Make sure we check all color entries, not just up to line length
				// This is important for lines that are all spaces with background colors
				maxX := len(state.LineColors[y])
				if len(line) > maxX {
					maxX = len(line)
				}
				bgFound := false
				for x := 0; x < maxX && x < len(state.LineColors[y]); x++ {
					style := state.LineColors[y][x]
					// Only render background if present and not empty
					if style.BgColor != "" && style.BgColor != nilValue && style.BgColor != nullValue {
						bgFound = true
						if g.options.Debug {
							log.Printf("Rendering background rect at (%d,%d) with color %s", x, y, style.BgColor)
						}
						charX := float64(x) * g.charWidth
						sb.WriteString(fmt.Sprintf(`<rect x="%s" y="%s" width="%s" height="%s" fill="%s" shape-rendering="crispEdges"/>`,
							formatCoord(charX), formatCoord(float64(y)*g.charHeight*lineHeight), formatCoord(g.charWidth), formatCoord(g.charHeight), style.BgColor))
						g.writeNewline(&sb)
					}
				}
				if g.options.Debug && !bgFound && hasBackgroundColors {
					// Debug: print first few color entries to see what's happening
					log.Printf("Line %d claims to have background colors but none found. First few entries:", y)
					for i := 0; i < 5 && i < len(state.LineColors[y]); i++ {
						log.Printf("  [%d]: fg=%q bg=%q", i, state.LineColors[y][i].FgColor, state.LineColors[y][i].BgColor)
					}
				}
			}

			// Note: cursor background will be rendered inline with text to ensure proper alignment

			// Convert line to runes to handle UTF-8 properly
			runes := []rune(line)

			// If cursor is beyond line end, pad with spaces
			if isCursorLine && state.CursorX > len(runes) {
				// Pad the line to reach cursor position
				padding := state.CursorX - len(runes)
				for i := 0; i < padding; i++ {
					runes = append(runes, ' ')
				}
			}

			// For inline cursor positioning, we need to render in segments
			if isCursorLine && state.CursorChar != "" {
				// Split the line into two parts: before cursor and after cursor
				var beforeCursor, afterCursor string

				if state.CursorX < len(runes) {
					beforeCursor = string(runes[:state.CursorX])
					if state.CursorX+1 < len(runes) {
						afterCursor = string(runes[state.CursorX+1:])
					}
				} else {
					beforeCursor = string(runes)
				}

				// Render all text in a single text element with inline cursor
				// Add xml:space="preserve" to preserve whitespace
				sb.WriteString(fmt.Sprintf(`<text y="%s" xml:space="preserve">`, formatCoord(yPos)))

				// Render text before cursor with proper styling
				if beforeCursor != "" {
					// Use a temporary builder for the before-cursor segment
					var tempSb strings.Builder
					g.renderTextSegment(&tempSb, beforeCursor, y, 0, state.CursorX, hasColors, state.LineColors)
					sb.WriteString(tempSb.String())
				} else {
					// No text before cursor, start at x="0"
					sb.WriteString(`<tspan x="0"></tspan>`)
				}

				// Render cursor as inline element with background
				cursorClass := g.cursorActiveClass
				if !state.IsCursorActive {
					cursorClass = g.cursorIdleClass
				}

				// Get cursor color (cursor is rendered as a block with foreground color)
				cursorBgColor := g.options.Theme.Foreground
				if cursorBgColor == "" {
					cursorBgColor = defaultCursorColor
				}

				// Render cursor inline
				// For a true inline solution, we'll render the cursor as a colored block character
				if state.CursorChar != "" && state.CursorChar != " " {
					// Use the cursor character from xterm.js (usually █)
					sb.WriteString(fmt.Sprintf(`<tspan class="%s %s" style="fill:%s;">%s</tspan>`,
						g.textClass, cursorClass, cursorBgColor, html.EscapeString(state.CursorChar)))
				} else {
					// Fallback to block character
					sb.WriteString(fmt.Sprintf(`<tspan class="%s %s" style="fill:%s;">█</tspan>`,
						g.textClass, cursorClass, cursorBgColor))
				}

				// Render text after cursor
				if afterCursor != "" {
					// Create a modified renderTextSegment call that doesn't add x="0"
					afterRunes := []rune(afterCursor)
					for i := 0; i < len(afterRunes); {
						// Similar logic to renderTextSegment but without x positioning
						charPos := state.CursorX + 1 + i

						// Get initial style
						var styleStr string
						var colorClass string
						if hasColors && y < len(state.LineColors) && charPos < len(state.LineColors[y]) {
							style := state.LineColors[y][charPos]
							if style.FgColor != "" && style.FgColor != nilValue {
								colorClass = g.getColorClass(style.FgColor)
								if colorClass == "" {
									styleStr = fmt.Sprintf("fill:%s;", style.FgColor)
								}
							}
							if style.Bold {
								styleStr += fontWeightBold
							}
							if style.Italic {
								styleStr += fontStyleItalic
							}
							if style.Underline {
								styleStr += textDecorationUnderline
							}
						}

						// Collect characters with same style
						segmentText := string(afterRunes[i])
						i++

						for i < len(afterRunes) {
							nextCharPos := state.CursorX + 1 + i
							var nextStyleStr string
							var nextColorClass string
							if hasColors && y < len(state.LineColors) && nextCharPos < len(state.LineColors[y]) {
								nextStyle := state.LineColors[y][nextCharPos]
								if nextStyle.FgColor != "" && nextStyle.FgColor != nilValue {
									nextColorClass = g.getColorClass(nextStyle.FgColor)
									if nextColorClass == "" {
										nextStyleStr = fmt.Sprintf("fill:%s;", nextStyle.FgColor)
									}
								}
								if nextStyle.Bold {
									nextStyleStr += fontWeightBold
								}
								if nextStyle.Italic {
									nextStyleStr += fontStyleItalic
								}
								if nextStyle.Underline {
									nextStyleStr += textDecorationUnderline
								}
							}

							if styleStr != nextStyleStr || colorClass != nextColorClass {
								break
							}

							segmentText += string(afterRunes[i])
							i++
						}

						// Render the segment without x positioning
						classes := g.textClass
						if colorClass != "" {
							classes += " " + colorClass
						}

						if styleStr != "" {
							sb.WriteString(fmt.Sprintf(`<tspan class="%s" style="%s">%s</tspan>`, classes, styleStr, html.EscapeString(segmentText)))
						} else {
							sb.WriteString(fmt.Sprintf(`<tspan class="%s">%s</tspan>`, classes, html.EscapeString(segmentText)))
						}
					}
				}

				sb.WriteString("</text>")
				g.writeNewline(&sb)
			} else {
				// No cursor on this line, render normally
				// Add xml:space="preserve" to preserve whitespace
				sb.WriteString(fmt.Sprintf(`<text y="%s" xml:space="preserve">`, formatCoord(yPos)))
				g.renderTextSegment(&sb, string(runes), y, 0, len(runes), hasColors, state.LineColors)
				sb.WriteString("</text>")
				g.writeNewline(&sb)
			}
		}
	}

	sb.WriteString("</g>")
	g.writeNewline(&sb)

	return sb.String()
}

// getColorClass returns the appropriate CSS class for a color.
func (g *SVGGenerator) getColorClass(color string) string {
	if !g.options.OptimizeSize {
		return ""
	}

	theme := g.options.Theme
	switch color {
	case theme.Black:
		return "k"
	case theme.Red:
		return "r"
	case theme.Green:
		return "g"
	case theme.Yellow:
		return "y"
	case theme.Blue:
		return "b"
	case theme.Magenta:
		return "m"
	case theme.Cyan:
		return "c"
	case theme.White:
		return "w"
	case theme.BrightBlue, "#5a56e0":
		return "p" // prompt color
	default:
		return ""
	}
}

// renderTextSegment renders a segment of text with appropriate styling.
func (g *SVGGenerator) renderTextSegment(sb *strings.Builder, text string, lineIndex, startChar, _ int, hasColors bool, lineColors [][]CharStyle) {
	runes := []rune(text)
	x := 0

	for x < len(runes) {
		// Group consecutive characters with the same style
		startX := x
		var style CharStyle
		var styleStr string

		// Get style if we have color data for this position
		charPos := startChar + x
		var colorClass string
		if hasColors && lineIndex < len(lineColors) && charPos < len(lineColors[lineIndex]) {
			style = lineColors[lineIndex][charPos]
			// Check if we can use a color class instead of inline style
			if style.FgColor != "" && style.FgColor != nilValue {
				colorClass = g.getColorClass(style.FgColor)
				if colorClass == "" {
					// No matching class, use inline style
					styleStr += fmt.Sprintf("fill:%s;", style.FgColor)
				}
			}
			if style.Bold {
				styleStr += fontWeightBold
			}
			if style.Italic {
				styleStr += fontStyleItalic
			}
			if style.Underline {
				styleStr += textDecorationUnderline
			}
		}

		// Collect characters with same style
		segmentText := string(runes[x])
		x++

		for x < len(runes) {
			nextCharPos := startChar + x
			var nextStyleStr string

			var nextColorClass string
			if hasColors && lineIndex < len(lineColors) && nextCharPos < len(lineColors[lineIndex]) {
				nextStyle := lineColors[lineIndex][nextCharPos]
				if nextStyle.FgColor != "" && nextStyle.FgColor != nilValue {
					nextColorClass = g.getColorClass(nextStyle.FgColor)
					if nextColorClass == "" {
						nextStyleStr += fmt.Sprintf("fill:%s;", nextStyle.FgColor)
					}
				}
				if nextStyle.Bold {
					nextStyleStr += fontWeightBold
				}
				if nextStyle.Italic {
					nextStyleStr += fontStyleItalic
				}
				if nextStyle.Underline {
					nextStyleStr += textDecorationUnderline
				}
			}

			// If style changes, break
			if styleStr != nextStyleStr || colorClass != nextColorClass {
				break
			}

			segmentText += string(runes[x])
			x++
		}

		// Render the segment
		classes := g.textClass
		if colorClass != "" {
			classes += " " + colorClass
		}

		if x == 1 && startX == 0 {
			// First segment, position at x=0
			if styleStr != "" {
				// x="0" is needed for tspan to reset x position
				fmt.Fprintf(sb, `<tspan x="0" class="%s" style="%s">%s</tspan>`, classes, styleStr, html.EscapeString(segmentText))
			} else {
				// x="0" is needed for tspan to reset x position
				fmt.Fprintf(sb, `<tspan x="0" class="%s">%s</tspan>`, classes, html.EscapeString(segmentText))
			}
		} else {
			// Let text flow naturally
			if styleStr != "" {
				fmt.Fprintf(sb, `<tspan class="%s" style="%s">%s</tspan>`, classes, styleStr, html.EscapeString(segmentText))
			} else {
				fmt.Fprintf(sb, `<tspan class="%s">%s</tspan>`, classes, html.EscapeString(segmentText))
			}
		}
	}
}

// writeNewline conditionally writes a newline based on optimization settings.
func (g *SVGGenerator) writeNewline(sb *strings.Builder) {
	if !g.options.OptimizeSize {
		sb.WriteString("\n")
	}
}

// formatCoord formats a coordinate value with minimal decimal places.
func formatCoord(val float64) string {
	// If it's a whole number, don't include decimals
	if val == float64(int(val)) {
		return fmt.Sprintf("%d", int(val))
	}
	// Otherwise use 1 decimal place and remove trailing zeros
	formatted := fmt.Sprintf("%.1f", val)
	// Remove trailing .0 if present
	if strings.HasSuffix(formatted, ".0") {
		return formatted[:len(formatted)-2]
	}
	return formatted
}

// formatPercentage formats a percentage value with appropriate precision
// to avoid keyframe collisions in large animations.
// The precision is dynamically calculated based on the number of keyframes.
func formatPercentage(val float64, keyframeCount int) string {
	// For whole numbers, keep minimal format
	if val == float64(int(val)) {
		return fmt.Sprintf("%d", int(val))
	}
	
	// Dynamically determine precision based on keyframe count
	// This ensures we have enough precision to avoid collisions
	// while keeping the output as compact as possible
	var precision int
	switch {
	case keyframeCount < 100:
		precision = 1  // Up to 100 unique values
	case keyframeCount < 1000:
		precision = 2  // Up to 1,000 unique values
	case keyframeCount < 10000:
		precision = 3  // Up to 10,000 unique values
	case keyframeCount < 100000:
		precision = 4  // Up to 100,000 unique values
	default:
		precision = 5  // Up to 1,000,000 unique values
	}
	
	// Format with calculated precision
	formatStr := fmt.Sprintf("%%.%df", precision)
	formatted := fmt.Sprintf(formatStr, val)
	
	// Remove trailing zeros but keep at least 1 decimal for consistency
	formatted = strings.TrimRight(formatted, "0")
	formatted = strings.TrimSuffix(formatted, ".")
	
	return formatted
}

// formatDuration formats a duration value with minimal decimal places.
func formatDuration(seconds float64) string {
	// If it's a whole number, don't include decimals
	if seconds == float64(int(seconds)) {
		return fmt.Sprintf("%d", int(seconds))
	}
	// Otherwise use 2 decimal places and remove trailing zeros
	formatted := fmt.Sprintf("%.2f", seconds)
	// Remove trailing zeros
	formatted = strings.TrimRight(formatted, "0")
	// Remove trailing decimal point if no decimals remain
	if strings.HasSuffix(formatted, ".") {
		return formatted[:len(formatted)-1]
	}
	return formatted
}

// generateCursorSymbols creates reusable cursor symbols.
func (g *SVGGenerator) generateCursorSymbols() string {
	// We're now using inline cursor rendering, so no symbols needed
	return ""
}

// generateTerminalWindow creates the terminal window chrome.
func (g *SVGGenerator) generateTerminalWindow() string {
	var sb strings.Builder

	style := g.options.Style
	if style == nil {
		style = DefaultStyleOptions()
	}

	// Window background with configurable rounded corners
	borderRadius := style.BorderRadius
	if borderRadius < 0 {
		borderRadius = 0
	}

	// Use background color from style or theme
	bgColor := style.BackgroundColor
	if bgColor == "" {
		bgColor = defaultBarColor
	}

	sb.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" rx="%d" fill="%s"/>`,
		g.options.Width, g.options.Height, borderRadius, bgColor))
	g.writeNewline(&sb)

	// Window bar if enabled
	if style.WindowBar != "" {
		sb.WriteString(g.generateWindowBar())
	}

	return sb.String()
}

// generateWindowBar creates the window bar based on style.
func (g *SVGGenerator) generateWindowBar() string {
	var sb strings.Builder

	style := g.options.Style
	if style == nil {
		style = DefaultStyleOptions()
	}

	barSize := style.WindowBarSize
	barColor := style.WindowBarColor
	if barColor == "" {
		barColor = defaultBarColor
	}

	// Bar background
	borderRadius := style.BorderRadius
	if borderRadius < 0 {
		borderRadius = 0
	}

	sb.WriteString(`<g id="window-bar">`)
	g.writeNewline(&sb)

	// Bar background with rounded top corners
	sb.WriteString(fmt.Sprintf(`<path d="M %d,0 L %d,0 Q %d,0 %d,%d L %d,%d L 0,%d L 0,%d Q 0,0 %d,0 Z" fill="%s"/>`,
		borderRadius, g.options.Width-borderRadius, g.options.Width, g.options.Width, borderRadius,
		g.options.Width, barSize, barSize, borderRadius, borderRadius, barColor))
	g.writeNewline(&sb)

	// Window controls based on style
	switch style.WindowBar {
	case "Colorful":
		// Colorful circles on the left (macOS-style)
		for i, color := range windowControlColors {
			x := 20 + i*20
			sb.WriteString(fmt.Sprintf(`<circle cx="%d" cy="%d" r="6" fill="%s"/>`, x, barSize/2, color))
			g.writeNewline(&sb)
		}
	case "ColorfulRight":
		// Colorful circles on the right
		for i, color := range windowControlColors {
			x := g.options.Width - 80 + i*20
			sb.WriteString(fmt.Sprintf(`<circle cx="%d" cy="%d" r="6" fill="%s"/>`, x, barSize/2, color))
			g.writeNewline(&sb)
		}
	case "Rings":
		// Ring circles on the left
		for i, color := range windowControlColors {
			x := 20 + i*20
			sb.WriteString(fmt.Sprintf(`<circle cx="%d" cy="%d" r="6" fill="none" stroke="%s" stroke-width="1"/>`, x, barSize/2, color))
			g.writeNewline(&sb)
		}
	case "RingsRight":
		// Ring circles on the right
		for i, color := range windowControlColors {
			x := g.options.Width - 80 + i*20
			sb.WriteString(fmt.Sprintf(`<circle cx="%d" cy="%d" r="6" fill="none" stroke="%s" stroke-width="1"/>`, x, barSize/2, color))
			g.writeNewline(&sb)
		}
	}

	// Title text (if provided)
	if style.WindowBarTitle != "" {
		// Get the appropriate font family with fallbacks
		fontFamily := getWindowBarFontFamily(style, g.options.FontFamily)
		// Get the appropriate font size with fallback
		fontSize := style.WindowBarFontSize
		if fontSize == 0 {
			fontSize = g.options.FontSize
		}

		// Calculate vertical position with proper padding for different font sizes
		// For SVG text, the y position is the baseline
		// We need to ensure there's enough top padding, especially for larger fonts

		// Calculate the text baseline position
		// Rule: ensure adequate top padding
		minTopPadding := int(float64(fontSize) * fontPaddingRatio)
		if minTopPadding < minFontPadding {
			minTopPadding = minFontPadding
		}

		// The baseline should be positioned considering:
		// - Font ascent calculation
		// - We want the text centered but with adequate top padding
		fontAscent := int(float64(fontSize) * fontAscentRatio)

		// Calculate baseline position
		yPos := minTopPadding + fontAscent

		// But also try to center in the bar if there's room
		centerBaseline := (barSize + fontAscent) / 2
		if centerBaseline > yPos && (centerBaseline+int(float64(fontSize)*fontPaddingRatio)) <= barSize {
			yPos = centerBaseline
		}

		// Add text with padding constraints
		// The text will be centered but constrained to avoid overlapping with window controls
		// Window controls occupy roughly 80px on each side
		centerX := g.options.Width / 2

		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" font-family="%s" font-size="%d" fill="#cccccc">`,
			centerX, yPos, fontFamily, fontSize))
		sb.WriteString(html.EscapeString(style.WindowBarTitle))
		sb.WriteString(`</text>`)
		g.writeNewline(&sb)
	}

	sb.WriteString("</g>")
	g.writeNewline(&sb)

	return sb.String()
}

// CaptureSVGFrame captures the current terminal state and returns an SVGFrame.
func CaptureSVGFrame(page *rod.Page, counter int, framerate int) (*SVGFrame, error) {
	// Get cursor position and exact character positions from xterm.js
	termInfo, err := page.Eval(`() => {
		const term = window.term;
		if (!term) {
			console.error('term is not available');
			return null;
		}
		const buffer = term.buffer.active;
		const cursorX = buffer.cursorX;
		// Cursor Y is relative to the viewport (0 = top of visible area)
		const cursorY = buffer.cursorY;
		
		// Debug logging
		const cursorLine = buffer.getLine(cursorY + buffer.viewportY);
		if (cursorLine) {
			const lineText = cursorLine.translateToString(true);
			console.log('Cursor Debug - xterm.js cursor position:', cursorX, 'on line:', JSON.stringify(lineText));
			console.log('Cursor Debug - Line length:', lineText.length, 'chars');
			
			// More detailed debugging
			if (cursorX < lineText.length) {
				console.log('Cursor Debug - Character at cursor position (' + cursorX + '):', JSON.stringify(lineText[cursorX]));
				console.log('Cursor Debug - Text before cursor:', JSON.stringify(lineText.substring(0, cursorX)));
				console.log('Cursor Debug - Text including cursor:', JSON.stringify(lineText.substring(0, cursorX + 1)));
			} else {
				console.log('Cursor Debug - Cursor is at end of line (position ' + cursorX + ')');
			}
			
			// Check what xterm.js thinks about cursor positioning
			console.log('Cursor Debug - buffer.cursorX:', buffer.cursorX);
			console.log('Cursor Debug - Is cursor at line end?', cursorX >= lineText.length);
		}
		
		// Get character dimensions
		let charWidth = 0;
		let charHeight = 0;
		
		// Get dimensions from the rendered canvas
		// This is the most reliable source as it represents the actual rendered output
		const textCanvas = document.querySelector('canvas.xterm-text-layer');
		const cols = term.cols;
		const rows = term.rows;
		charWidth = textCanvas.width / cols;
		charHeight = textCanvas.height / rows;
		
		// Get cursor character from buffer
		let cursorChar = '█'; // Default block cursor
		
		// Helper function to convert xterm.js color to hex
		function xtermColorToHex(color) {
			if (!color) return null;
			// Handle RGB colors
			if (color.mode === 'rgb') {
				const r = color.r.toString(16).padStart(2, '0');
				const g = color.g.toString(16).padStart(2, '0');
				const b = color.b.toString(16).padStart(2, '0');
				return '#' + r + g + b;
			}
			// Handle palette colors (0-255)
			if (typeof color === 'number') {
				// Use xterm.js's color palette
				const palette = term.options.theme;
				if (color < 16 && palette) {
					// Basic 16 colors
					const colorNames = ['black', 'red', 'green', 'yellow', 'blue', 'magenta', 'cyan', 'white',
									   'brightBlack', 'brightRed', 'brightGreen', 'brightYellow', 
									   'brightBlue', 'brightMagenta', 'brightCyan', 'brightWhite'];
					return palette[colorNames[color]] || null;
				}
				// For extended colors, we'd need the full palette
				return null;
			}
			return null;
		}
		
		// Get color information for all visible lines
		const lineColors = [];
		const activeBuffer = term.buffer.active;
		console.log('term exists:', !!term, 'buffer exists:', !!term.buffer, 'active exists:', !!activeBuffer);
		const viewportStart = activeBuffer ? activeBuffer.viewportY : 0;
		const viewportEnd = viewportStart + term.rows;
		
		console.log('Capturing colors for viewport:', viewportStart, 'to', viewportEnd, 'buffer length:', activeBuffer ? activeBuffer.length : 'no buffer');
		
		let cellCount = 0;
		for (let y = viewportStart; y < viewportEnd && activeBuffer && y < activeBuffer.length; y++) {
			const line = activeBuffer.getLine(y);
			const lineColorData = [];
			
			if (line) {
				// Get the full line including trailing spaces
				// translateToString(true) preserves trailing whitespace
				const lineText = line.translateToString(true);
				// Use term.cols to ensure we capture all columns, not just non-empty ones
				for (let x = 0; x < term.cols; x++) {
					const cell = line.getCell(x);
					if (cell) {
						cellCount++;
						const chars = cell.getChars();
						let fgColor = null;
						let bgColor = null;
						
						
						// Check if cell has foreground color
						if (cell.isFgRGB()) {
							const fg = cell.getFgColor();
							fgColor = '#' + ((fg >> 16) & 0xff).toString(16).padStart(2, '0') +
									 ((fg >> 8) & 0xff).toString(16).padStart(2, '0') +
									 (fg & 0xff).toString(16).padStart(2, '0');
						} else if (cell.isFgPalette()) {
							// Handle palette colors
							const paletteIndex = cell.getFgColor();
							if (paletteIndex >= 0 && paletteIndex < 16) {
								const colorNames = ['black', 'red', 'green', 'yellow', 'blue', 'magenta', 'cyan', 'white',
												   'brightBlack', 'brightRed', 'brightGreen', 'brightYellow', 
												   'brightBlue', 'brightMagenta', 'brightCyan', 'brightWhite'];
								const palette = term.options.theme;
								if (palette && palette[colorNames[paletteIndex]]) {
									fgColor = palette[colorNames[paletteIndex]];
								}
							}
						}
						
						// Check if cell has background color
						if (cell.isBgRGB()) {
							const bg = cell.getBgColor();
							bgColor = '#' + ((bg >> 16) & 0xff).toString(16).padStart(2, '0') +
									((bg >> 8) & 0xff).toString(16).padStart(2, '0') +
									(bg & 0xff).toString(16).padStart(2, '0');
							console.log('Found RGB background color:', bgColor, 'at', x, y);
						} else if (cell.isBgPalette()) {
							// Handle palette colors for background
							const paletteIndex = cell.getBgColor();
							if (paletteIndex >= 0 && paletteIndex < 16) {
								const colorNames = ['black', 'red', 'green', 'yellow', 'blue', 'magenta', 'cyan', 'white',
												   'brightBlack', 'brightRed', 'brightGreen', 'brightYellow', 
												   'brightBlue', 'brightMagenta', 'brightCyan', 'brightWhite'];
								const palette = term.options.theme;
								const colorName = colorNames[paletteIndex];
								console.log('Background - Palette index:', paletteIndex, 'colorName:', colorName, 'theme exists:', !!palette);
								if (palette && palette[colorName]) {
									bgColor = palette[colorName];
									console.log('Found palette background color:', bgColor, 'for', colorName, 'at', x, y);
								} else {
									// Use default ANSI colors if theme doesn't have them
									const defaultColors = ['#000000', '#cc0000', '#4e9a06', '#c4a000', '#3465a4', '#75507b', '#06989a', '#d3d7cf',
														   '#555753', '#ef2929', '#8ae234', '#fce94f', '#729fcf', '#ad7fa8', '#34e2e2', '#eeeeec'];
									bgColor = defaultColors[paletteIndex];
									console.log('Using default color:', bgColor, 'for palette index:', paletteIndex);
								}
							}
						}
						
						
						lineColorData.push({
							char: chars || ' ',
							fgColor: fgColor === null ? '' : fgColor,
							bgColor: bgColor === null ? '' : bgColor,
							bold: cell.isBold() === 1,
							italic: cell.isItalic() === 1,
							underline: cell.isUnderline() === 1
						});
					}
				}
			}
			lineColors.push(lineColorData);
		}
		
		// Note: We no longer need to calculate character positions since
		// we're using text-push positioning in SVG
		
		
		return {
			cursorX: cursorX,
			cursorY: cursorY,
			charWidth: charWidth,
			charHeight: charHeight,
			lineColors: lineColors,
			cursorChar: cursorChar
		};
	}`)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate terminal info: %w", err)
	}

	// Get buffer content as lines
	bufferResult, err := page.Eval(`() => {
		const term = window.term;
		const buffer = term.buffer.active;
		const lines = [];
		
		// Get all visible lines
		const viewportStart = buffer.viewportY;
		const viewportEnd = viewportStart + term.rows;
		
		for (let y = viewportStart; y < viewportEnd && y < buffer.length; y++) {
			const line = buffer.getLine(y);
			if (line) {
				// translateToString(true) preserves trailing whitespace
				lines.push(line.translateToString(true).trimEnd());
			} else {
				lines.push('');
			}
		}
		
		return lines;
	}`)
	if err != nil {
		return nil, fmt.Errorf("failed to get buffer content: %w", err)
	}

	// Convert buffer result to string slice
	buffer := []string{}
	bufferArray := bufferResult.Value.Arr()
	for _, line := range bufferArray {
		buffer = append(buffer, line.Str())
	}

	// Parse the terminal info
	cursorX := termInfo.Value.Get("cursorX").Int()
	cursorY := termInfo.Value.Get("cursorY").Int()
	charWidth := termInfo.Value.Get("charWidth").Num()
	charHeight := termInfo.Value.Get("charHeight").Num()
	cursorChar := termInfo.Value.Get("cursorChar").Str()

	// Parse line colors
	lineColors := [][]CharStyle{}
	lineColorsJSON := termInfo.Value.Get("lineColors")
	if !lineColorsJSON.Nil() {
		lines := lineColorsJSON.Arr()
		for _, line := range lines {
			lineStyles := []CharStyle{}
			chars := line.Arr()
			for _, charData := range chars {
				// Handle nil/null values properly
				fgColor := ""
				bgColor := ""
				fgColorVal := charData.Get("fgColor")
				bgColorVal := charData.Get("bgColor")
				if !fgColorVal.Nil() && fgColorVal.Str() != nilValue {
					fgColor = fgColorVal.Str()
				}
				if !bgColorVal.Nil() && bgColorVal.Str() != nilValue {
					bgColor = bgColorVal.Str()
				}

				style := CharStyle{
					FgColor:   fgColor,
					BgColor:   bgColor,
					Bold:      charData.Get("bold").Bool(),
					Italic:    charData.Get("italic").Bool(),
					Underline: charData.Get("underline").Bool(),
				}
				lineStyles = append(lineStyles, style)
			}
			lineColors = append(lineColors, lineStyles)
		}
	}

	svgFrame := &SVGFrame{
		Lines:      buffer,
		LineColors: lineColors,
		CursorX:    cursorX,
		CursorY:    cursorY,
		CharWidth:  charWidth,
		CharHeight: charHeight,
		Timestamp:  float64(counter) / float64(framerate),
		CursorChar: cursorChar,
	}

	return svgFrame, nil
}

// parseFontFamily parses a font family string and returns a list of individual fonts.
func parseFontFamily(fontFamily string) []string {
	if fontFamily == "" {
		return []string{svgDefaultFontFamily}
	}

	// Split by comma and clean up each font name
	fonts := strings.Split(fontFamily, ",")
	result := make([]string, 0, len(fonts))

	for _, font := range fonts {
		font = strings.TrimSpace(font)
		// Remove quotes if present
		font = strings.Trim(font, "\"'")
		if font != "" {
			result = append(result, font)
		}
	}

	return result
}

// buildSVGFontFamily creates a properly formatted font-family string for SVG
// with appropriate fallbacks.
func buildSVGFontFamily(fontFamily string) string {
	fonts := parseFontFamily(fontFamily)

	// Build a new list without quotes - SVG will handle the attribute quoting
	fontList := make([]string, 0, len(fonts)+1)
	hasMonospace := false

	for _, font := range fonts {
		// Check if this is a generic font family
		if font == svgDefaultFontFamily || font == "ui-monospace" || font == "sans-serif" || font == "serif" {
			fontList = append(fontList, font)
			if font == svgDefaultFontFamily {
				hasMonospace = true
			}
		} else {
			// Add font names as-is, SVG attribute will be quoted
			fontList = append(fontList, font)
		}
	}

	// Always add monospace as fallback if not already present
	if !hasMonospace {
		fontList = append(fontList, svgDefaultFontFamily)
	}

	return strings.Join(fontList, ", ")
}

// getWindowBarFontFamily returns the font family to use for window bar titles
// with proper fallback handling.
func getWindowBarFontFamily(style *StyleOptions, defaultFontFamily string) string {
	// Use WindowBarFontFamily if specified
	if style.WindowBarFontFamily != "" {
		return buildSVGFontFamily(style.WindowBarFontFamily)
	}

	// Fall back to FontFamily from style
	if style.FontFamily != "" {
		return buildSVGFontFamily(style.FontFamily)
	}

	// Fall back to default font family
	if defaultFontFamily != "" {
		return buildSVGFontFamily(defaultFontFamily)
	}

	// Ultimate fallback
	return svgDefaultFontFamily
}
