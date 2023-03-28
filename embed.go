package main

import (
	_ "embed"
)

// DemoTape contains the demo.tape file contents to place in a new tape.
//go:embed examples/demo.tape
var DemoTape []byte
