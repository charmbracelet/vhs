package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/charmbracelet/vhs"
)

func main() {
	ensureInstalled("ffmpeg", "ttyd")

	var b []byte
	var err error

	if len(os.Args) > 1 {
		b, err = os.ReadFile(os.Args[1])
	} else {
		b, err = io.ReadAll(os.Stdin)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = vhs.Evaluate(string(b), os.Stdout, "")
	if err != nil {
		os.Exit(1)
	}
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
