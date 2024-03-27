package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// TestOptions is the set of options for the testing functionality.
type TestOptions struct {
	Output string
	Golden string
}

// DefaultTestOptions returns the default set of options for the testing functionality.
func DefaultTestOptions() TestOptions {
	return TestOptions{
		Output: "out.test",
	}
}

// Alternatively, `var separator = strings.Repeat("─", 80)`.
const separator = "────────────────────────────────────────────────────────────────────────────────"

var (
	once sync.Once
	file *os.File
)

// SaveOutput saves the current buffer to the output file.
func (v *VHS) SaveOutput() {
	// Create output file (once)
	once.Do(func() {
		err := os.MkdirAll(filepath.Dir(v.Options.Test.Output), os.ModePerm)
		if err != nil {
			file, _ = os.CreateTemp(os.TempDir(), "vhs-*.txt")
			return
		}
		file, _ = os.Create(v.Options.Test.Output)
	})

	lines, err := v.Buffer()
	if err != nil {
		return
	}

	for _, line := range lines {
		_, _ = file.WriteString(line + "\n")
	}

	_, _ = file.WriteString(separator + "\n")
}

// Buffer returns the current buffer.
func (v *VHS) Buffer() ([]string, error) {
	// Get the current buffer.
	buf, err := v.Page.Eval("() => Array(term.rows).fill(0).map((e, i) => term.buffer.active.getLine(term.buffer.active.viewportY+i).translateToString().trimEnd())")
	if err != nil {
		return nil, fmt.Errorf("read buffer: %w", err)
	}

	var lines []string
	for _, line := range buf.Value.Arr() {
		lines = append(lines, line.Str())
	}

	return lines, nil
}

// CurrentLine returns the current line from the buffer.
func (v *VHS) CurrentLine() (string, error) {
	buf, err := v.Page.Eval("() => term.buffer.active.getLine(term.buffer.active.cursorY+term.buffer.active.viewportY).translateToString().trimEnd()")
	if err != nil {
		return "", fmt.Errorf("read curent line from buffer: %w", err)
	}

	return buf.Value.Str(), nil
}
