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
		err = os.MkdirAll(filepath.Dir(v.Options.Test.Output), os.ModePerm)
		if err != nil {
			file, err = os.CreateTemp(os.TempDir(), "vhs-*.txt")
			return
		}
		file, err = os.Create(v.Options.Test.Output)
	})
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	// Get the current buffer.
	buf, err := v.Page.Eval("() => Array(term.rows).fill(0).map((e, i) => term.buffer.active.getLine(i).translateToString().trimEnd())")
	if err != nil {
		return fmt.Errorf("failed to get buffer: %w", err)
	}

	for _, line := range buf.Value.Arr() {
		str := line.Str()
		_, err = file.WriteString(str + "\n")
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
