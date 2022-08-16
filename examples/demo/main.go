package main

import (
	"os"

	"github.com/charmbracelet/dolly"
)

func main() {
	b, err := os.ReadFile("dolly.vhs")
	if err != nil {
		panic(err)
	}
	cmds, err := dolly.Parse(string(b))
	if err != nil {
		panic(err)
	}
	dolly.Run(cmds)
}
