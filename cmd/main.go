package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/vhs"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: vhs < file.tape")
		os.Exit(1)
	}

	b, err := os.ReadFile(os.Args[1])
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
