package main

import (
	"fmt"

	mcobra "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"github.com/spf13/cobra"
)

const (
	manDescription = `VHS let's you write terminal GIFs as code.
VHS reads .tape files and renders GIFs (videos).
A tape file is a script made up of commands describing what actions to perform in the render.

The following is a list of all possible commands in VHS:

• Output <path>.(gif|webm|mp4)
• Set <setting> <value>
• Sleep <time>
• Type "<string>"
• Ctrl+<key>
• Backspace [repeat]
• Down [repeat]
• Enter [repeat]
• Left [repeat]
• Right [repeat]
• Tab [repeat]
• Up [repeat]
• Hide
• Show
`

	manOutput = `The Output command instructs VHS where to save the output of the recording.
File names with the extension .gif, .webm, .mp4 will have the respective file types.
`

	manSettings = `The Set command allows VHS to adjust settings in the terminal, such as fonts, dimensions, and themes.

The following is a list of all possible setting commands in VHS:

• Set FontSize <number>
• Set FontFamily <string>
• Set Height <number>
• Set Width <number>
• Set LetterSpacing <float>
• Set LineHeight <float>
• Set TypingSpeed <time>
• Set Theme <json>
• Set Padding <number>
• Set Framerate <number>
`

	manBugs = "See GitHub Issues: <https://github.com/charmbracelet/vhs/issues>"

	manAuthor = "Charm <vt100@charm.sh>"
)

var (
	manCmd = &cobra.Command{
		Use:    "man",
		Short:  "Generate man pages",
		Args:   cobra.NoArgs,
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			manPage, err := mcobra.NewManPage(1, rootCmd)
			if err != nil {
				return err
			}

			manPage = manPage.WithLongDescription(manDescription).
				WithSection("Output", manOutput).
				WithSection("Settings", manSettings).
				WithSection("Bugs", manBugs).
				WithSection("Author", manAuthor).
				WithSection("Copyright", "(C) 2021-2022 Charmbracelet, Inc.\n"+
					"Released under MIT license.")
			fmt.Println(manPage.Build(roff.NewDocument()))
			return nil
		},
	}
)
