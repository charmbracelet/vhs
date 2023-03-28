package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/mattn/go-isatty"
	mcobra "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"github.com/spf13/cobra"
)

// We use this special character to represent the character that delimits the
// important symbols in the text. We replace these to code in markdown and bold
// in roff.
const specialChar = "%"

var (
	manDescription = `VHS let's you write terminal GIFs as code.
VHS reads .tape files and renders GIFs (videos).
A tape file is a script made up of commands describing what actions to perform in the render.

The following is a list of all possible commands in VHS:

* %Output% <path>.(gif|webm|mp4)
* %Require% <program>
* %Set% <setting> <value>
* %Sleep% <time>
* %Type% "<string>"
* %Ctrl%+<key>
* %Backspace% [repeat]
* %Down% [repeat]
* %Enter% [repeat]
* %Left% [repeat]
* %Right% [repeat]
* %Tab% [repeat]
* %Up% [repeat]
* %PageUp% [repeat]
* %PageDown% [repeat]
* %Hide%
* %Show%
`

	manOutput = `The Output command instructs VHS where to save the output of the recording.
File names with the extension %.gif%, %.webm%, %.mp4% will have the respective file types.
`

	manSettings = `The Set command allows VHS to adjust settings in the terminal, such as fonts, dimensions, and themes.

The following is a list of all possible setting commands in VHS:

* Set %Shell% <string>
* Set %FontSize% <number>
* Set %FontFamily% <string>
* Set %Height% <number>
* Set %Width% <number>
* Set %LetterSpacing% <float>
* Set %LineHeight% <float>
* Set %TypingSpeed% <time>
* Set %Theme% <json|string>
* Set %Padding% <number>
* Set %Framerate% <number>
* Set %PlaybackSpeed% <float>
`
	manBugs = "See GitHub Issues: <https://github.com/charmbracelet/vhs/issues>"

	manAuthor = "Charm <vt100@charm.sh>"
)

var manCmd = &cobra.Command{
	Use:     "manual",
	Aliases: []string{"man"},
	Short:   "Generate man pages",
	Args:    cobra.NoArgs,
	Hidden:  true,
	RunE: func(_ *cobra.Command, _ []string) error {
		if isatty.IsTerminal(os.Stdout.Fd()) {
			renderer, err := glamour.NewTermRenderer(
				glamour.WithStyles(GlamourTheme),
			)
			if err != nil {
				return err
			}
			v, err := renderer.Render(markdownManual())
			if err != nil {
				return err
			}
			fmt.Println(v)
			return nil
		}

		manPage, err := mcobra.NewManPage(1, rootCmd)
		if err != nil {
			return err
		}

		manPage = manPage.
			WithLongDescription(sanitizeSpecial(manDescription)).
			WithSection("Output", sanitizeSpecial(manOutput)).
			WithSection("Settings", sanitizeSpecial(manSettings)).
			WithSection("Bugs", sanitizeSpecial(manBugs)).
			WithSection("Author", sanitizeSpecial(manAuthor)).
			WithSection("Copyright", "(C) 2021-2022 Charmbracelet, Inc.\n"+
				"Released under MIT license.")

		fmt.Println(manPage.Build(roff.NewDocument()))
		return nil
	},
}

func markdownManual() string {
	return fmt.Sprint(
		"# MANUAL\n", sanitizeMarkdown(manDescription),
		"# OUTPUT\n", sanitizeMarkdown(manOutput),
		"# SETTING\n", sanitizeMarkdown(manSettings),
		"# BUGS\n", manBugs,
		"\n# AUTHOR\n", manAuthor,
	)
}

func sanitizeMarkdown(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(
		s, "<", "&lt;"), ">", "&gt;"), specialChar, "`")
}

func sanitizeSpecial(s string) string {
	return strings.ReplaceAll(s, specialChar, "")
}
