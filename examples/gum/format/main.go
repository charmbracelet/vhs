package main

import (
	"os"
	"time"

	"github.com/charmbracelet/dolly"
)

const markdownFile = "readme.md"
const markdownContent = `# Gum Formats
- Markdown
- Code
- Template
- Emoji`

const codeFile = "main.go"
const codeContent = `package main

import "fmt"

func main() {
	fmt.Println("Charm™ Gum!")
}`

const templateFile = "tmpl.txt"

func main() {
	d := dolly.New(dolly.WithOutput("format.gif"), dolly.WithFontSize(24), dolly.WithLineHeight(1))
	defer d.Cleanup()

	os.WriteFile(markdownFile, []byte(markdownContent), 0644)
	os.WriteFile(codeFile, []byte(codeContent), 0644)
	defer os.Remove(markdownFile)
	defer os.Remove(codeFile)
	defer os.Remove(templateFile)

	d.Type("gum format < "+markdownFile, dolly.WithSpeed(50))
	d.Enter()
	time.Sleep(3 * time.Second)

	d.Clear()
	time.Sleep(time.Second)

	d.Type("gum format < main.go -t code", dolly.WithSpeed(50))
	d.Enter()
	time.Sleep(3 * time.Second)

	d.Clear()
	time.Sleep(time.Second)

	d.Type(`echo '{{ Bold "Tasty" }} {{ Color "99" "0" " Gum! " }}' > ` + templateFile)
	d.Enter()
	time.Sleep(time.Second / 2)
	d.Type(`cat ` + templateFile + ` | gum format -t template`)
	d.Enter()
	time.Sleep(2 * time.Second)

	d.Clear()
	time.Sleep(time.Second)

	d.Type(`echo 'I :heart: Bubble Gum :candy:' | gum format -t emoji`, dolly.WithSpeed(50))
	d.Enter()
	time.Sleep(2 * time.Second)
}
