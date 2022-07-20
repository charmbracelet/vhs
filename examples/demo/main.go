package main

import (
	"time"

	"github.com/charmbracelet/frame"
)

func main() {
	f := frame.New(frame.WithFontSize(42))
	defer f.Cleanup()

	f.Type("echo 'Hello, Demo!'", 50*time.Millisecond)
	f.Enter()

	time.Sleep(time.Second)
}
