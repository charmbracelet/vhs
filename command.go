package vhs

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-rod/rod/lib/input"
)

// CommandType is a type that represents a command.
type CommandType string

const (
	Backspace CommandType = "Backspace"
	Down      CommandType = "Down"
	Enter     CommandType = "Enter"
	Left      CommandType = "Left"
	Right     CommandType = "Right"
	Space     CommandType = "Space"
	Up        CommandType = "Up"
	Type      CommandType = "Type"
	Set       CommandType = "Set"
	Sleep     CommandType = "Sleep"
	Ctrl      CommandType = "Ctrl"
	Unknown   CommandType = "Unknown"
)

// CommandTypes is a list of the available commands that can be executed.
var CommandTypes = []CommandType{
	Backspace,
	Down,
	Enter,
	Left,
	Right,
	Space,
	Up,
	Type,
	Set,
	Sleep,
	Ctrl,
	Unknown,
}

// String returns the string representation of the command.
func (c CommandType) String() string {
	return string(c)
}

// CommandFunc is a function that executes a command on a running
// instance of vhs.
type CommandFunc func(c Command, d *VHS)

// CommandFuncs maps command types to their executable functions.
var CommandFuncs = map[CommandType]CommandFunc{
	Backspace: ExecuteKey(input.Backspace),
	Down:      ExecuteKey(input.ArrowDown),
	Enter:     ExecuteKey(input.Enter),
	Left:      ExecuteKey(input.ArrowLeft),
	Right:     ExecuteKey(input.ArrowRight),
	Space:     ExecuteKey(input.Space),
	Up:        ExecuteKey(input.ArrowUp),
	Set:       ExecuteSet,
	Sleep:     ExecuteSleep,
	Type:      ExecuteType,
	Ctrl:      ExecuteCtrl,
	Unknown:   ExecuteNoop,
}

// Command represents a command with options and arguments.
type Command struct {
	Type    CommandType
	Options string
	Args    string
}

// String returns the string representation of the command.
// This includes the options and arguments of the command.
func (c Command) String() string {
	if c.Options != "" {
		return fmt.Sprintf("%s %s %s", c.Type, c.Options, c.Args)
	}
	return fmt.Sprintf("%s %s", c.Type, c.Args)
}

// Execute executes a command on a running instance of vhs.
func (c Command) Execute(d *VHS) {
	CommandFuncs[c.Type](c, d)
}

// ExecuteNoop is a no-op command that does nothing.
// Generally, this is used for Unknown commands when dealing with
// commands that are not recognized.
func ExecuteNoop(c Command, d *VHS) {}

// ExecuteKey is a higher-order function that returns a CommandFunc to execute
// a key press for a given key. This is so that the logic for key pressing
// (since they are repeatable and delayable) can be re-used.
//
// i.e. ExecuteKey(input.ArrowDown) would return a CommandFunc that executes
// the ArrowDown key press.
func ExecuteKey(k input.Key) CommandFunc {
	return func(c Command, d *VHS) {
		repeat, err := strconv.Atoi(c.Args)
		if err != nil {
			repeat = 1
		}
		delay, err := time.ParseDuration(c.Options)
		if err != nil {
			delay = time.Millisecond * 100
		}
		for i := 0; i < repeat; i++ {
			_ = d.Page.Keyboard.Type(k)
			time.Sleep(delay)
		}
	}
}

// ExecuteCtrl is a CommandFunc that presses the argument key with the ctrl key
// held down on the running instance of vhs.
func ExecuteCtrl(c Command, d *VHS) {
	_ = d.Page.Keyboard.Press(input.ControlLeft)
	for _, r := range c.Args {
		if k, ok := keymap[r]; ok {
			_ = d.Page.Keyboard.Type(k)
		}
	}
	_ = d.Page.Keyboard.Release(input.ControlLeft)
}

// ExecuteSleep sleeps for the desired time specified through the argument of
// the Sleep command.
func ExecuteSleep(c Command, d *VHS) {
	dur, err := time.ParseDuration(c.Args)
	if err != nil {
		return
	}
	time.Sleep(dur)
}

// ExecuteType types the argument string on the running instance of vhs.
func ExecuteType(c Command, d *VHS) {
	for _, r := range c.Args {
		k, ok := keymap[r]
		if ok {
			_ = d.Page.Keyboard.Type(k)
		} else {
			_ = d.Page.MustElement("textarea").Input(string(r))
			d.Page.MustWaitIdle()
		}
		delayMs, err := strconv.Atoi(c.Options)
		if err != nil {
			delayMs = 100
		}
		time.Sleep(time.Millisecond * time.Duration(delayMs))
	}
}

// Settings maps the Set commands to their respective functions.
var Settings = map[string]CommandFunc{
	"FontSize":      ApplyFontSize,
	"FontFamily":    ApplyFontFamily,
	"Height":        ApplyHeight,
	"Width":         ApplyWidth,
	"LetterSpacing": ApplyLetterSpacing,
	"LineHeight":    ApplyLineHeight,
	"Theme":         ApplyTheme,
	"Padding":       ApplyPadding,
	"Framerate":     ApplyFramerate,
	"Output":        ApplyOutput,
}

// ExecuteSet applies the settings on the running vhs specified by the
// option and argument pass to the command.
func ExecuteSet(c Command, d *VHS) {
	Settings[c.Options](c, d)
}

// ApplyFontSize applies the font size on the vhs.
func ApplyFontSize(c Command, d *VHS) {
	fontSize, _ := strconv.Atoi(c.Args)
	d.Options.FontSize = fontSize
	_, _ = d.Page.Eval(fmt.Sprintf("term.setOption('fontSize', '%d')", fontSize))
}

// ApplyFontFamily applies the font family on the vhs.
func ApplyFontFamily(c Command, d *VHS) {
	d.Options.FontFamily = c.Args
	_, _ = d.Page.Eval(fmt.Sprintf("term.setOption('fontFamily', '%s')", c.Args))
}

// ApplyHeight applies the height on the vhs.
func ApplyHeight(c Command, d *VHS) {
	d.Options.Height, _ = strconv.Atoi(c.Args)
}

// ApplyWidth applies the width on the vhs.
func ApplyWidth(c Command, d *VHS) {
	d.Options.Width, _ = strconv.Atoi(c.Args)
	d.Options.GIF.Width, _ = strconv.Atoi(c.Args)
}

// ApplyLetterSpacing applies letter tracking (also known as tracking) on the
// vhs.
func ApplyLetterSpacing(c Command, d *VHS) {
	letterSpacing, _ := strconv.ParseFloat(c.Args, 64)
	d.Options.LetterSpacing = letterSpacing
	_, _ = d.Page.Eval(fmt.Sprintf("term.setOption('letterSpacing', '%f')", letterSpacing))
}

// ApplyLineHeight applies the line height on the vhs.
func ApplyLineHeight(c Command, d *VHS) {
	lineHeight, _ := strconv.ParseFloat(c.Args, 64)
	d.Options.LineHeight = lineHeight
	_, _ = d.Page.Eval(fmt.Sprintf("term.setOption('lineHeight', '%f')", lineHeight))
}

// ApplyTheme applies the theme on the vhs.
func ApplyTheme(c Command, d *VHS) {
	err := json.Unmarshal([]byte(c.Args), &d.Options.Theme)
	if err != nil {
		d.Options.Theme = DefaultTheme
		return
	}
	_, _ = d.Page.Eval(fmt.Sprintf("term.setOption('theme', %s)", c.Args))
}

// ApplyPadding applies the padding on the vhs.
func ApplyPadding(c Command, d *VHS) {
	d.Options.Padding = c.Args
	_, _ = d.Page.MustElement(".xterm").Eval(fmt.Sprintf(`this.style.padding = '%s'`, c.Args))
}

// ApplyFramerate applies the framerate on the vhs.
func ApplyFramerate(c Command, d *VHS) {
	d.Options.Framerate, _ = strconv.ParseFloat(c.Args, 64)
}

// ApplyOutput applies the output on the vhs GIF.
func ApplyOutput(c Command, d *VHS) {
	d.Options.GIF.Output = c.Args
}
