package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/vhs/lexer"
	"github.com/charmbracelet/vhs/parser"
	"github.com/charmbracelet/vhs/token"
	"github.com/go-rod/rod/lib/input"
	"github.com/mattn/go-runewidth"
)

// Execute executes a command on a running instance of vhs.
func Execute(c parser.Command, v *VHS) {
	if c.Type == token.SOURCE {
		ExecuteSourceTape(c, v)
	} else {
		CommandFuncs[c.Type](c, v)
	}

	if v.recording && v.Options.Test.Output != "" {
		v.SaveOutput()
	}
}

// CommandFunc is a function that executes a command on a running
// instance of vhs.
type CommandFunc func(c parser.Command, v *VHS)

// CommandFuncs maps command types to their executable functions.
var CommandFuncs = map[parser.CommandType]CommandFunc{
	token.BACKSPACE:  ExecuteKey(input.Backspace),
	token.DELETE:     ExecuteKey(input.Delete),
	token.INSERT:     ExecuteKey(input.Insert),
	token.DOWN:       ExecuteKey(input.ArrowDown),
	token.ENTER:      ExecuteKey(input.Enter),
	token.LEFT:       ExecuteKey(input.ArrowLeft),
	token.RIGHT:      ExecuteKey(input.ArrowRight),
	token.SPACE:      ExecuteKey(input.Space),
	token.UP:         ExecuteKey(input.ArrowUp),
	token.TAB:        ExecuteKey(input.Tab),
	token.ESCAPE:     ExecuteKey(input.Escape),
	token.PAGEUP:     ExecuteKey(input.PageUp),
	token.PAGEDOWN:   ExecuteKey(input.PageDown),
	token.HIDE:       ExecuteHide,
	token.REQUIRE:    ExecuteRequire,
	token.SHOW:       ExecuteShow,
	token.SET:        ExecuteSet,
	token.OUTPUT:     ExecuteOutput,
	token.SLEEP:      ExecuteSleep,
	token.TYPE:       ExecuteType,
	token.CTRL:       ExecuteCtrl,
	token.ALT:        ExecuteAlt,
	token.SHIFT:      ExecuteShift,
	token.ILLEGAL:    ExecuteNoop,
	token.SCREENSHOT: ExecuteScreenshot,
	token.COPY:       ExecuteCopy,
	token.PASTE:      ExecutePaste,
	token.ENV:        ExecuteEnv,
}

// ExecuteNoop is a no-op command that does nothing.
// Generally, this is used for Unknown commands when dealing with
// commands that are not recognized.
func ExecuteNoop(_ parser.Command, _ *VHS) {}

// ExecuteKey is a higher-order function that returns a CommandFunc to execute
// a key press for a given key. This is so that the logic for key pressing
// (since they are repeatable and delayable) can be re-used.
//
// i.e. ExecuteKey(input.ArrowDown) would return a CommandFunc that executes
// the ArrowDown key press.
func ExecuteKey(k input.Key) CommandFunc {
	return func(c parser.Command, v *VHS) {
		typingSpeed, err := time.ParseDuration(c.Options)
		if err != nil {
			typingSpeed = v.Options.TypingSpeed
		}
		repeat, err := strconv.Atoi(c.Args)
		if err != nil {
			repeat = 1
		}
		for i := 0; i < repeat; i++ {
			_ = v.Page.Keyboard.Type(k)
			time.Sleep(typingSpeed)
		}
	}
}

// ExecuteCtrl is a CommandFunc that presses the argument keys and/or modifiers
// with the ctrl key held down on the running instance of vhs.
func ExecuteCtrl(c parser.Command, v *VHS) {
	// Create key combination by holding ControlLeft
	action := v.Page.KeyActions().Press(input.ControlLeft)
	keys := strings.Split(c.Args, " ")

	for i, key := range keys {
		var inputKey *input.Key

		switch key {
		case "Shift":
			inputKey = &input.ShiftLeft
		case "Alt":
			inputKey = &input.AltLeft
		case "Enter":
			inputKey = &input.Enter
		case "Space":
			inputKey = &input.Space
		case "Backspace":
			inputKey = &input.Backspace
		default:
			r := rune(key[0])
			if k, ok := keymap[r]; ok {
				inputKey = &k
			}
		}

		// Press or hold key in case it's valid
		if inputKey != nil {
			if i != len(keys)-1 {
				action.Press(*inputKey)
			} else {
				// Other keys will remain pressed until the combination reaches the end
				action.Type(*inputKey)
			}
		}
	}

	action.MustDo()
}

// ExecuteAlt is a CommandFunc that presses the argument key with the alt key
// held down on the running instance of vhs.
func ExecuteAlt(c parser.Command, v *VHS) {
	_ = v.Page.Keyboard.Press(input.AltLeft)
	if k, ok := token.Keywords[c.Args]; ok {
		switch k {
		case token.ENTER:
			_ = v.Page.Keyboard.Type(input.Enter)
		case token.TAB:
			_ = v.Page.Keyboard.Type(input.Tab)
		}
	} else {
		for _, r := range c.Args {
			if k, ok := keymap[r]; ok {
				_ = v.Page.Keyboard.Type(k)
			}
		}
	}

	_ = v.Page.Keyboard.Release(input.AltLeft)
}

// ExecuteShift is a CommandFunc that presses the argument key with the shift
// key held down on the running instance of vhs.
func ExecuteShift(c parser.Command, v *VHS) {
	_ = v.Page.Keyboard.Press(input.ShiftLeft)
	if k, ok := token.Keywords[c.Args]; ok {
		switch k {
		case token.ENTER:
			_ = v.Page.Keyboard.Type(input.Enter)
		case token.TAB:
			_ = v.Page.Keyboard.Type(input.Tab)
		}
	} else {
		for _, r := range c.Args {
			if k, ok := keymap[r]; ok {
				_ = v.Page.Keyboard.Type(k)
			}
		}
	}

	_ = v.Page.Keyboard.Release(input.AltLeft)
}

// ExecuteHide is a CommandFunc that starts or stops the recording of the vhs.
func ExecuteHide(_ parser.Command, v *VHS) {
	v.PauseRecording()
}

// ExecuteRequire is a CommandFunc that checks if all the binaries mentioned in the
// Require command are present. If not, it exits with a non-zero error.
func ExecuteRequire(c parser.Command, v *VHS) {
	_, err := exec.LookPath(c.Args)
	if err != nil {
		v.Errors = append(v.Errors, err)
	}
}

// ExecuteShow is a CommandFunc that resumes the recording of the vhs.
func ExecuteShow(_ parser.Command, v *VHS) {
	v.ResumeRecording()
}

// ExecuteSleep sleeps for the desired time specified through the argument of
// the Sleep command.
func ExecuteSleep(c parser.Command, _ *VHS) {
	dur, err := time.ParseDuration(c.Args)
	if err != nil {
		return
	}
	time.Sleep(dur)
}

// ExecuteType types the argument string on the running instance of vhs.
func ExecuteType(c parser.Command, v *VHS) {
	typingSpeed, err := time.ParseDuration(c.Options)
	if err != nil {
		typingSpeed = v.Options.TypingSpeed
	}
	for _, r := range c.Args {
		k, ok := keymap[r]
		if ok {
			_ = v.Page.Keyboard.Type(k)
		} else {
			_ = v.Page.MustElement("textarea").Input(string(r))
			v.Page.MustWaitIdle()
		}
		time.Sleep(typingSpeed)
	}
}

// ExecuteOutput applies the output on the vhs videos.
func ExecuteOutput(c parser.Command, v *VHS) {
	switch c.Options {
	case ".mp4":
		v.Options.Video.Output.MP4 = c.Args
	case ".test", ".ascii", ".txt":
		v.Options.Test.Output = c.Args
	case ".png":
		v.Options.Video.Output.Frames = c.Args
	case ".webm":
		v.Options.Video.Output.WebM = c.Args
	default:
		v.Options.Video.Output.GIF = c.Args
	}
}

// ExecuteCopy copies text to the clipboard.
func ExecuteCopy(c parser.Command, _ *VHS) {
	_ = clipboard.WriteAll(c.Args)
}

// ExecuteEnv sets env with given key-value pair.
func ExecuteEnv(c parser.Command, v *VHS) {
	_ = os.Setenv(c.Options, c.Args)

}

// ExecutePaste pastes text from the clipboard.
func ExecutePaste(_ parser.Command, v *VHS) {
	clip, err := clipboard.ReadAll()
	if err != nil {
		return
	}
	for _, r := range clip {
		k, ok := keymap[r]
		if ok {
			_ = v.Page.Keyboard.Type(k)
		} else {
			_ = v.Page.MustElement("textarea").Input(string(r))
			v.Page.MustWaitIdle()
		}
	}
}

// Settings maps the Set commands to their respective functions.
var Settings = map[string]CommandFunc{
	"FontFamily":    ExecuteSetFontFamily,
	"FontSize":      ExecuteSetFontSize,
	"Framerate":     ExecuteSetFramerate,
	"Height":        ExecuteSetHeight,
	"LetterSpacing": ExecuteSetLetterSpacing,
	"LineHeight":    ExecuteSetLineHeight,
	"PlaybackSpeed": ExecuteSetPlaybackSpeed,
	"Padding":       ExecuteSetPadding,
	"Theme":         ExecuteSetTheme,
	"TypingSpeed":   ExecuteSetTypingSpeed,
	"Width":         ExecuteSetWidth,
	"Shell":         ExecuteSetShell,
	"LoopOffset":    ExecuteLoopOffset,
	"MarginFill":    ExecuteSetMarginFill,
	"Margin":        ExecuteSetMargin,
	"WindowBar":     ExecuteSetWindowBar,
	"WindowBarSize": ExecuteSetWindowBarSize,
	"BorderRadius":  ExecuteSetBorderRadius,
	"CursorBlink":   ExecuteSetCursorBlink,
}

// ExecuteSet applies the settings on the running vhs specified by the
// option and argument pass to the command.
func ExecuteSet(c parser.Command, v *VHS) {
	Settings[c.Options](c, v)
}

// ExecuteSetFontSize applies the font size on the vhs.
func ExecuteSetFontSize(c parser.Command, v *VHS) {
	fontSize, _ := strconv.Atoi(c.Args)
	v.Options.FontSize = fontSize
	_, _ = v.Page.Eval(fmt.Sprintf("() => term.options.fontSize = %d", fontSize))

	// When changing the font size only the canvas dimensions change which are
	// scaled back during the render to fit the aspect ration and dimensions.
	//
	// We need to call term.fit to ensure that everything is resized properly.
	_, _ = v.Page.Eval("term.fit")
}

// ExecuteSetFontFamily applies the font family on the vhs.
func ExecuteSetFontFamily(c parser.Command, v *VHS) {
	v.Options.FontFamily = c.Args
	_, _ = v.Page.Eval(fmt.Sprintf("() => term.options.fontFamily = '%s'", withSymbolsFallback(c.Args)))
}

// ExecuteSetHeight applies the height on the vhs.
func ExecuteSetHeight(c parser.Command, v *VHS) {
	v.Options.Video.Style.Height, _ = strconv.Atoi(c.Args)
}

// ExecuteSetWidth applies the width on the vhs.
func ExecuteSetWidth(c parser.Command, v *VHS) {
	v.Options.Video.Style.Width, _ = strconv.Atoi(c.Args)
}

// ExecuteSetShell applies the shell on the vhs.
func ExecuteSetShell(c parser.Command, v *VHS) {
	if s, ok := Shells[c.Args]; ok {
		v.Options.Shell = s
	}
}

const (
	bitSize = 64
	base    = 10
)

// ExecuteSetLetterSpacing applies letter spacing (also known as tracking) on
// the vhs.
func ExecuteSetLetterSpacing(c parser.Command, v *VHS) {
	letterSpacing, _ := strconv.ParseFloat(c.Args, bitSize)
	v.Options.LetterSpacing = letterSpacing
	_, _ = v.Page.Eval(fmt.Sprintf("() => term.options.letterSpacing = %f", letterSpacing))
}

// ExecuteSetLineHeight applies the line height on the vhs.
func ExecuteSetLineHeight(c parser.Command, v *VHS) {
	lineHeight, _ := strconv.ParseFloat(c.Args, bitSize)
	v.Options.LineHeight = lineHeight
	_, _ = v.Page.Eval(fmt.Sprintf("() => term.options.lineHeight = %f", lineHeight))
}

// ExecuteSetTheme applies the theme on the vhs.
func ExecuteSetTheme(c parser.Command, v *VHS) {
	var err error
	v.Options.Theme, err = getTheme(c.Args)
	if err != nil {
		v.Errors = append(v.Errors, err)
		return
	}

	bts, _ := json.Marshal(v.Options.Theme)
	_, _ = v.Page.Eval(fmt.Sprintf("() => term.options.theme = %s", string(bts)))
	v.Options.Video.Style.BackgroundColor = v.Options.Theme.Background
	v.Options.Video.Style.WindowBarColor = v.Options.Theme.Background
}

// ExecuteSetTypingSpeed applies the default typing speed on the vhs.
func ExecuteSetTypingSpeed(c parser.Command, v *VHS) {
	typingSpeed, err := time.ParseDuration(c.Args)
	if err != nil {
		return
	}
	v.Options.TypingSpeed = typingSpeed
}

// ExecuteSetPadding applies the padding on the vhs.
func ExecuteSetPadding(c parser.Command, v *VHS) {
	v.Options.Video.Style.Padding, _ = strconv.Atoi(c.Args)
}

// ExecuteSetFramerate applies the framerate on the vhs.
func ExecuteSetFramerate(c parser.Command, v *VHS) {
	framerate, err := strconv.ParseInt(c.Args, base, 0)
	if err != nil {
		return
	}
	v.Options.Video.Framerate = int(framerate)
}

// ExecuteSetPlaybackSpeed applies the playback speed option on the vhs.
func ExecuteSetPlaybackSpeed(c parser.Command, v *VHS) {
	playbackSpeed, err := strconv.ParseFloat(c.Args, bitSize)
	if err != nil {
		return
	}
	v.Options.Video.PlaybackSpeed = playbackSpeed
}

// ExecuteLoopOffset applies the loop offset option on the vhs.
func ExecuteLoopOffset(c parser.Command, v *VHS) {
	loopOffset, err := strconv.ParseFloat(strings.TrimRight(c.Args, "%"), bitSize)
	if err != nil {
		return
	}
	v.Options.LoopOffset = loopOffset
}

// ExecuteSetMarginFill sets vhs margin fill
func ExecuteSetMarginFill(c parser.Command, v *VHS) {
	v.Options.Video.Style.MarginFill = c.Args
}

// ExecuteSetMargin sets vhs margin size
func ExecuteSetMargin(c parser.Command, v *VHS) {
	v.Options.Video.Style.Margin, _ = strconv.Atoi(c.Args)
}

// ExecuteSetWindowBar sets window bar type
func ExecuteSetWindowBar(c parser.Command, v *VHS) {
	v.Options.Video.Style.WindowBar = c.Args
}

// ExecuteSetWindowBar sets window bar size
func ExecuteSetWindowBarSize(c parser.Command, v *VHS) {
	v.Options.Video.Style.WindowBarSize, _ = strconv.Atoi(c.Args)
}

// ExecuteSetBorderRadius sets corner radius
func ExecuteSetBorderRadius(c parser.Command, v *VHS) {
	v.Options.Video.Style.BorderRadius, _ = strconv.Atoi(c.Args)
}

// ExecuteSetCursorBlink sets cursor blinking
func ExecuteSetCursorBlink(c parser.Command, v *VHS) {
	var err error
	v.Options.CursorBlink, err = strconv.ParseBool(c.Args)
	if err != nil {
		return
	}
}

const sourceDisplayMaxLength = 10

// ExecuteSourceTape is a CommandFunc that executes all commands of source tape.
func ExecuteSourceTape(c parser.Command, v *VHS) {
	tapePath := c.Args
	var out io.Writer = os.Stdout
	if quietFlag {
		out = io.Discard
	}

	// read tape file
	tape, err := os.ReadFile(tapePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	l := lexer.New(string(tape))
	p := parser.New(l)

	cmds := p.Parse()

	errs := []error{}
	for _, parsedErr := range p.Errors() {
		errs = append(errs, parsedErr)
	}

	if len(errs) != 0 {
		fmt.Fprintln(out, ErrorStyle.Render(fmt.Sprintf("tape %s has errors", tapePath)))
		printErrors(out, tapePath, errs)
		return
	}

	displayPath := runewidth.Truncate(strings.TrimSuffix(tapePath, extension), sourceDisplayMaxLength, "…")

	// Run all commands from the sourced tape file.
	for _, cmd := range cmds {
		// Output have to be avoid in order to not overwrite output of the original tape.
		if cmd.Type == token.SOURCE ||
			cmd.Type == token.OUTPUT {
			continue
		}
		fmt.Fprintf(out, "%s %s\n", GrayStyle.Render(displayPath+":"), Highlight(cmd, false))
		CommandFuncs[cmd.Type](cmd, v)
	}
}

// ExecuteScreenshot is a CommandFunc that indicates a new screenshot must be taken.
func ExecuteScreenshot(c parser.Command, v *VHS) {
	v.ScreenshotNextFrame(c.Args)
}

func getTheme(s string) (Theme, error) {
	if strings.TrimSpace(s) == "" {
		return DefaultTheme, nil
	}
	switch s[0] {
	case '{':
		return getJSONTheme(s)
	default:
		return findTheme(s)
	}
}

func getJSONTheme(s string) (Theme, error) {
	var t Theme
	if err := json.Unmarshal([]byte(s), &t); err != nil {
		return DefaultTheme, fmt.Errorf("invalid `Set Theme %q: %w`", s, err)
	}
	return t, nil
}
