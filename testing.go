package vhs

// import (
// 	"fmt"
// )

// // TestOptions is the set of options for the testing functionality.
// type TestOptions struct {
// 	Output string
// 	Golden string
// }

// // DefaultTestOptions returns the default set of options for the testing functionality.
// func DefaultTestOptions() TestOptions {
// 	return TestOptions{
// 		Output: "out.test",
// 	}
// }

// // SaveOutput saves the current buffer to the output file.
// func (v *VHS) SaveOutput() {
// 	o, err := v.Page.Eval("() => Array(term.rows).fill(0).map((e, i) => term.buffer.normal.getLine(i).translateToString())")
// 	if err != nil {
// 		return
// 	}

// 	fmt.Println("Saving frame to ", testFile.Name())
// 	for _, line := range o.Value.Arr() {
// 		_, _ = testFile.WriteString(line.Str() + "\n")
// 	}
// 	_, _ = testFile.WriteString("--------------------------------------------------------------------------------\n")
// }
