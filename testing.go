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
func (v *VHS) SaveOutput() error {
	var err error
	// Create output file (once)
	once.Do(func() {
		err = os.MkdirAll(filepath.Dir(v.Options.Test.Output), 0o750)
		if err != nil {
			file, err = os.CreateTemp(os.TempDir(), "vhs-*.txt")
			return
		}
		file, err = os.Create(v.Options.Test.Output)
	})
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	lines, err := v.Buffer()
	if err != nil {
		return fmt.Errorf("failed to get buffer: %w", err)
	}

	for _, line := range lines {
		_, err = file.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("failed to write buffer to file: %w", err)
		}
	}

	_, err = file.WriteString(separator + "\n")
	if err != nil {
		return fmt.Errorf("failed to write separator to file: %w", err)
	}

	return nil
}

// Buffer returns the current buffer.
func (v *VHS) Buffer() ([]string, error) {
	// Get the current buffer.
	buf, err := v.Page.Eval("() => Array(term.rows).fill(0).map((e, i) => term.buffer.active.getLine(i + term.buffer.active.viewportY).translateToString().trimEnd())")
	if err != nil {
		return nil, fmt.Errorf("read buffer: %w", err)
	}

	arr := buf.Value.Arr()
	lines := make([]string, 0, len(arr))
	for _, line := range arr {
		lines = append(lines, line.Str())
	}

	return lines, nil
}

// CurrentLine returns the current line from the buffer.
func (v *VHS) CurrentLine() (string, error) {
	buf, err := v.Page.Eval("() => term.buffer.active.getLine(term.buffer.active.cursorY+term.buffer.active.viewportY).translateToString().trimEnd()")
	if err != nil {
		return "", fmt.Errorf("read current line from buffer: %w", err)
	}

	return buf.Value.Str(), nil
}
