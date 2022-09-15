package vhs

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

const frameSeparator = "---"

var (
	once sync.Once
	file *os.File
)

// SaveOutput saves the current buffer to the output file.
func (v *VHS) SaveOutput() {

	// Create output file (once)
	once.Do(func() {
		err := os.MkdirAll(filepath.Dir(v.Options.Test.Output), 0770)
		if err != nil {
			file, _ = os.CreateTemp(os.TempDir(), "vhs-*.txt")
			return
		}
		file, _ = os.Create(v.Options.Test.Output)
	})

	// Get the current normal buffer.
	buf, err := v.Page.Eval("() => Array(term.rows).fill(0).map((e, i) => term.buffer.normal.getLine(i).translateToString().trimEnd())")
	if err != nil {
		return
	}

	for _, line := range buf.Value.Arr() {
		str := line.Str()
		if str == "" {
			continue
		}
		_, _ = file.WriteString(str + "\n")
	}

	_, _ = file.WriteString(frameSeparator + "\n")

	// Get the current alternative buffer.
	alt, err := v.Page.Eval("() => Array(term.rows).fill(0).map((e, i) => term.buffer.alternate.getLine(i).translateToString().trimEnd())")
	if err != nil {
		return
	}

	for _, line := range alt.Value.Arr() {
		_, _ = file.WriteString(line.Str() + "\n")
	}

	_, _ = file.WriteString(frameSeparator + "\n")
}
