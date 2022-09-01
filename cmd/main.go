package main

import (
	"fmt"
	"io"
	"log"
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

	d := vhs.New()
	defer d.Cleanup()

	l := vhs.NewLexer(string(b))
	p := vhs.NewParser(l)

	cmds := p.Parse()
	errs := p.Errors()
	if len(errs) != 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		os.Exit(1)
	}

	var offset int

	for i, cmd := range cmds {
		if cmd.Type == vhs.Set {
			log.Printf("Setting %s to %s", cmd.Options, cmd.Args)
			cmd.Execute(&d)
		} else {
			offset = i
			break
		}
	}

	d.Start()

	for _, cmd := range cmds[offset:] {
		log.Println(cmd)
		cmd.Execute(&d)
	}

}
