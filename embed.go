package vhs

import (
	_ "embed"
	"fmt"

	"github.com/charmbracelet/glamour"
)

//go:embed demo.tape
var DemoTape []byte

//go:embed help.txt
var HelpText []byte

func Help() {
	fmt.Println(string(HelpText))
}

//go:embed manual.md
var ManualText []byte

func Manual() {
	renderer, err := glamour.NewTermRenderer(glamour.WithAutoStyle())
	if err != nil {
		fmt.Println(string(ManualText))
	}
	man, err := renderer.RenderBytes(ManualText)
	if err != nil {
		fmt.Println(string(ManualText))
	}
	fmt.Println(string(man))
}
