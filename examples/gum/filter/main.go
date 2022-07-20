package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/dolly"
)

const flavorsFile = "flavors.txt"
const favFile = "fav.txt"

var flavorsContent = []byte(strings.Join([]string{
	"Strawberry",
	"Banana",
	"Orange",
	"Cherry",
}, "\n"))

func main() {
	d := dolly.New(
		dolly.WithOutput("filter.gif"),
		dolly.WithFontSize(40),
	)
	defer d.Cleanup()

	os.WriteFile(flavorsFile, flavorsContent, 0644)
	defer os.Remove(flavorsFile)
	defer os.Remove(favFile)

	d.Type(fmt.Sprintf("cat %s | gum filter > %s", flavorsFile, favFile), dolly.WithSpeed(40))
	d.Enter()
	time.Sleep(time.Second)

	d.Type("nana", dolly.WithSpeed(250))
	time.Sleep(time.Second / 4)

	d.CtrlU()
	time.Sleep(time.Second / 2)

	d.Type("erry", dolly.WithSpeed(250))
	time.Sleep(time.Second / 2)

	d.CtrlU()
	time.Sleep(time.Second / 2)

	d.Type("↓↓↑", dolly.WithSpeed(500))

	d.Enter()

	time.Sleep(time.Second)

	d.Type("cat "+favFile, dolly.WithSpeed(30))
	d.Enter()

	time.Sleep(time.Second)
}
