package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/dolly"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: dolly <file.vhs>")
		os.Exit(1)
	}

	b, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cmds, errs := dolly.Parse(string(b))
	if len(errs) != 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		os.Exit(1)
	}
	dolly.Run(cmds)
}
