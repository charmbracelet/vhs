package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <command>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "input":
		gumInput()
	case "write":
		gumWrite()
	case "filter":
		gumFilter()
	case "choose":
	case "spin":
	case "style":
	case "join":
	case "format":
	}
}
