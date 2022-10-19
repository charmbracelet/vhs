package main

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/charmbracelet/vhs"
	"github.com/charmbracelet/vhs/style"
	version "github.com/hashicorp/go-version"
)

var ttydMinVersion = version.Must(version.NewVersion("1.7.2"))

// main runs the VHS command line interface and handles the argument parsing.
func main() {
	ensurDependencies()

	var err error
	var command string

	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "help", "--help", "-h":
		vhs.PrintHelp()
	case "version", "--version", "-v":
		PrintVersion()
	case "man", "manual":
		vhs.PrintManual()
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

// Run runs a given tape file and generates its outputs.
func Run(args []string) error {
	if len(args) < 1 && !hasStdin() {
		vhs.PrintHelp()
		return errors.New("no input provided")
	}

	if len(args) > 1 {
		return errors.New("expects 1 file")
	}

	if hasStdin() {
		in, _ := io.ReadAll(os.Stdin)
		return vhs.Evaluate(string(in), os.Stdout, "")
	}

	f, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}
	fmt.Println(style.File.Render("File: " + args[0]))
	err = vhs.Evaluate(string(f), os.Stdout, "")
	if err != nil {
		return err
	}
	return nil
}

// hasStdin returns whether stdin has been piped in.
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

// Parse takes a glob file path and parses all the files to ensure they are
// valid without running them.
//
// This allows CI to ensure all tape files are valid with the given version
// of VHS.
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

// New creates a new tape file with example tape file contents and
// documentation. Contents are copied from by demo.tape.
func New(args []string) error {
	if len(args) != 1 {
		return errors.New("new expects a file name to create")
	}

	fileName := strings.TrimSuffix(args[0], extension) + extension

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}

	f.Write(vhs.DemoTape)

	fmt.Println("Created " + fileName)

	return nil
}

// Version stores the build version of VHS at the time of packages through -ldflags
//   go build -ldflags "-s -w -X=main.Version=$(VERSION)" cmd/vhs/vhs.go -o vhs
var Version string

// PrintVersion prints the version of VHS.
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

var versionRegex = regexp.MustCompile(`\d+\.\d+\.\d+`)

// getVersion returns the parsed version of a program
func getVersion(program string) *version.Version {
	cmd := exec.Command(program, "--version")
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	programVersion, _ := version.NewVersion(versionRegex.FindString(string(out)))
	return programVersion
}

// ensureDependencies ensures that all dependencies are correctly installed
// and versioned before continuing
func ensurDependencies() {
	var shouldExit bool

	_, ffmpegErr := exec.LookPath("ffmpeg")
	if ffmpegErr != nil {
		fmt.Println("ffmpeg is not installed.")
		fmt.Println("Install it from: http://ffmpeg.org")
		shouldExit = true
	}
	_, ttydErr := exec.LookPath("ttyd")
	if ttydErr != nil {
		fmt.Println("ttyd is not installed.")
		fmt.Println("Install it from: https://github.com/tsl0922/ttyd")
		shouldExit = true
	}
	_, bashErr := exec.LookPath("bash")
	if bashErr != nil {
		fmt.Println("bash is not installed.")
		shouldExit = true
	}

	if shouldExit {
		os.Exit(1)
	}

	ttydVersion := getVersion("ttyd")
	if ttydVersion == nil || ttydVersion.LessThan(ttydMinVersion) {
		fmt.Printf("ttyd version (%s) is out of date, VHS requires %s\n", ttydVersion, ttydMinVersion)
		fmt.Println("Install the latest version from: https://github.com/tsl0922/ttyd")
		os.Exit(1)
	}
}
