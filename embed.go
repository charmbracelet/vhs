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

func PrintHelp() {
	fmt.Println(string(HelpText))
}

//go:embed vhs.1.md
var ManualText []byte

func PrintManual() {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStyles(GlamourTheme),
	)
	if err != nil {
		fmt.Println(string(ManualText))
	}
	man, err := renderer.RenderBytes(ManualText)
	if err != nil {
		fmt.Println(string(ManualText))
	}
	fmt.Println(string(man))
}
