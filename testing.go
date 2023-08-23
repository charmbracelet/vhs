package main

import (
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

	// Get the current buffer.
	buf, err := v.mainTerm.Page.Eval("() => Array(term.rows).fill(0).map((e, i) => term.buffer.active.getLine(i).translateToString().trimEnd())")
	if err != nil {
		return
	}

	for _, line := range buf.Value.Arr() {
		str := line.Str()
		_, _ = file.WriteString(str + "\n")
	}

	_, _ = file.WriteString(separator + "\n")
}
