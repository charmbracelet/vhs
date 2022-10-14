package main

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"

	"github.com/charmbracelet/vhs"
	"github.com/charmbracelet/vhs/style"
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
	case "help", "--help", "-h":
		Help()
	case "version", "--version", "-v":
		PrintVersion()
	case "man", "manual":
		Manual()
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

//go:embed help.txt
var help []byte

func Help() {
	fmt.Println(string(help))
}

//go:embed manual.txt
var manual []byte

func Manual() {
	fmt.Println(string(manual))
}

var Version string

func PrintVersion() {
	if Version == "" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Sum != "" {
			Version = info.Main.Version
		} else {
			Version = "unknown (built from source)"
		}
	}
	version := fmt.Sprintf("VHS version %s", Version)
	fmt.Println(version)
}

func Run(args []string) error {
	var b []byte
	var err error
	if len(args) >= 1 {
		b, err = os.ReadFile(args[0])
		if err != nil {
			Help()
		}
	} else {
		if !hasStdin() {
			Help()
			return nil
		}
		b, err = io.ReadAll(os.Stdin)
	}
	if err != nil {
		return err
	}
	return vhs.Evaluate(string(b), os.Stdout, "")
}

func hasStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	if stat.Mode()&os.ModeNamedPipe == 0 && stat.Size() == 0 {
		return false
	}
	return true
}

func Parse(args []string) error {
	if len(args) < 1 {
		return errors.New("parse expects at least one file")
	}

	valid := true

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
			fmt.Println(style.ErrorFile.Render(file))
			for _, err := range errs {
				fmt.Print(vhs.LineNumber(err.Token.Line))
				fmt.Println(lines[err.Token.Line-1])
				fmt.Print(strings.Repeat(" ", err.Token.Column+5))
				fmt.Println(vhs.Underline(len(err.Token.Literal)), err.Msg)
				fmt.Println()
			}
			valid = false
		}
	}

	if !valid {
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
