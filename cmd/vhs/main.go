package main

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/vhs"
)

func main() {
	var b []byte
	var err error

	if len(os.Args) > 1 {
		b, err = os.ReadFile(os.Args[1])
	} else {
		b, err = io.ReadAll(os.Stdin)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = vhs.Evaluate(string(b), os.Stdout)
	if err != nil {
		os.Exit(1)
	}
}
