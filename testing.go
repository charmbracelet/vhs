package main

import (
	"fmt"
	"os"
	"path/filepath"
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

// SaveOutput saves the current buffer to the output file. The file is created
// on the first call of each run, and the later calls append to it.
func (v *VHS) SaveOutput() error {
	if v.outputFile == nil {
		file, err := v.createOutputFile()
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		v.outputFile = file
	}
	file := v.outputFile
	var err error

	for _, line := range v.Buffer() {
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

// Buffer returns the text of the visible terminal rows.
func (v *VHS) Buffer() []string {
	return v.term.Lines()
}

// CurrentLine returns the line the cursor is on.
func (v *VHS) CurrentLine() string {
	return v.term.CurrentLine()
}

// createOutputFile creates the file the test output is saved to, or a
// temporary file if its directory can't be created.
func (v *VHS) createOutputFile() (*os.File, error) {
	var file *os.File
	var err error
	if err = os.MkdirAll(filepath.Dir(v.Options.Test.Output), 0o750); err != nil {
		file, err = os.CreateTemp(os.TempDir(), "vhs-*.txt")
	} else {
		file, err = os.Create(v.Options.Test.Output)
	}
	if err != nil {
		return nil, fmt.Errorf("could not create %s: %w", v.Options.Test.Output, err)
	}
	return file, nil
}
