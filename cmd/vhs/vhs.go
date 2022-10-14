package main

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/vhs"
)

//go:embed demo.tape
var demoTape []byte

func main() {
	ensureInstalled("ffmpeg", "ttyd", "bash")

	var err error

	var command string
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "new":
		err = New(os.Args[2:])
	case "parse":
		err = Parse(os.Args[2:])
	default:
		err = Run(os.Args[1:])
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Run(args []string) error {
	var b []byte
	var err error
	if len(args) >= 1 {
		b, err = os.ReadFile(args[0])
	} else {
		b, err = io.ReadAll(os.Stdin)
	}
	if err != nil {
		return err
	}
	return vhs.Evaluate(string(b), os.Stdout, "")
}

var errorStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("8")).
	Foreground(lipgloss.Color("1")).
	Padding(0, 1).
	Width(80)

func Parse(args []string) error {
	if len(args) < 1 {
		return errors.New("parse expects at least one file")
	}

	passing := true

	for _, file := range args {
		b, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		l := vhs.NewLexer(string(b))
		p := vhs.NewParser(l)

		_ = p.Parse()
		errs := p.Errors()

		if len(errs) != 0 {
			lines := strings.Split(string(b), "\n")
			fmt.Println(errorStyle.Render(file))
			for _, err := range errs {
				fmt.Print(vhs.LineNumber(err.Token.Line))
				fmt.Println(lines[err.Token.Line-1])
				fmt.Print(strings.Repeat(" ", err.Token.Column+5))
				fmt.Println(vhs.Underline(len(err.Token.Literal)), err.Msg)
				fmt.Println()
			}
			passing = false
		}
	}

	if !passing {
		return errors.New("invalid tape file(s)")
	}

	return nil
}

const extension = ".tape"

func New(args []string) error {
	if len(args) != 1 {
		return errors.New("new expects a file name to create")
	}

	fileName := strings.TrimSuffix(args[0], extension) + extension

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}

	f.Write(demoTape)

	fmt.Println("Created " + fileName)

	return nil
}

func ensureInstalled(programs ...string) {
	var missing bool

	for _, p := range programs {
		_, err := exec.LookPath(p)
		if err != nil {
			fmt.Printf("%s is not installed.\n", p)
			missing = true
		}
	}

	if missing {
		fmt.Println("Required programs are missing, install and add them to your PATH.")
		os.Exit(1)
	}
}
