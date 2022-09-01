package dolly

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
// instance of dolly.
type CommandFunc func(c Command, d *Dolly)

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

// Execute executes a command on a running instance of dolly.
func (c Command) Execute(d *Dolly) {
	CommandFuncs[c.Type](c, d)
}

// ExecuteNoop is a no-op command that does nothing.
// Generally, this is used for Unknown commands when dealing with
// commands that are not recognized.
func ExecuteNoop(c Command, d *Dolly) {}

// ExecuteKey is a higher-order function that returns a CommandFunc to execute
// a key press for a given key. This is so that the logic for key pressing
// (since they are repeatable and delayable) can be re-used.
//
// i.e. ExecuteKey(input.ArrowDown) would return a CommandFunc that executes
// the ArrowDown key press.
func ExecuteKey(k input.Key) CommandFunc {
	return func(c Command, d *Dolly) {
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
// held down on the running instance of dolly.
func ExecuteCtrl(c Command, d *Dolly) {
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
func ExecuteSleep(c Command, d *Dolly) {
	dur, err := time.ParseDuration(c.Args)
	if err != nil {
		return
	}
	time.Sleep(dur)
}

// ExecuteType types the argument string on the running instance of dolly.
func ExecuteType(c Command, d *Dolly) {
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
	"FontSize":   ApplyFontSize,
	"FontFamily": ApplyFontFamily,
	"Height":     ApplyHeight,
	"Width":      ApplyWidth,
	"LineHeight": ApplyLineHeight,
	"Theme":      ApplyTheme,
	"Padding":    ApplyPadding,
	"Framerate":  ApplyFramerate,
	"Output":     ApplyOutput,
}

// ExecuteSet applies the settings on the running dolly specified by the
// option and argument pass to the command.
func ExecuteSet(c Command, d *Dolly) {
	Settings[c.Options](c, d)
}

// ApplyFontSize applies the font size on the dolly.
func ApplyFontSize(c Command, d *Dolly) {
	fontSize, _ := strconv.Atoi(c.Args)
	d.Options.FontSize = fontSize
	_, _ = d.Page.Eval(fmt.Sprintf("term.setOption('fontSize', '%d')", fontSize))
}

// ApplyFontFamily applies the font family on the dolly.
func ApplyFontFamily(c Command, d *Dolly) {
	d.Options.FontFamily = c.Args
	_, _ = d.Page.Eval(fmt.Sprintf("term.setOption('fontFamily', '%s')", c.Args))
}

// ApplyHeight applies the height on the dolly.
func ApplyHeight(c Command, d *Dolly) {
	d.Options.Height, _ = strconv.Atoi(c.Args)
}

// ApplyWidth applies the width on the dolly.
func ApplyWidth(c Command, d *Dolly) {
	d.Options.Width, _ = strconv.Atoi(c.Args)
	d.Options.GIF.Width, _ = strconv.Atoi(c.Args)
}

// ApplyLineHeight applies the line height on the dolly.
func ApplyLineHeight(c Command, d *Dolly) {
	lineHeight, _ := strconv.ParseFloat(c.Args, 64)
	d.Options.LineHeight = lineHeight
	_, _ = d.Page.Eval(fmt.Sprintf("term.setOption('lineHeight', '%f')", lineHeight))
}

// ApplyTheme applies the theme on the dolly.
func ApplyTheme(c Command, d *Dolly) {
	err := json.Unmarshal([]byte(c.Args), &d.Options.Theme)
	if err != nil {
		d.Options.Theme = DefaultTheme
		return
	}
	_, _ = d.Page.Eval(fmt.Sprintf("term.setOption('theme', %s)", c.Args))
}

// ApplyPadding applies the padding on the dolly.
func ApplyPadding(c Command, d *Dolly) {
	d.Options.Padding = c.Args
	_, _ = d.Page.MustElement(".xterm").Eval(fmt.Sprintf(`this.style.padding = '%s'`, c.Args))
}

// ApplyFramerate applies the framerate on the dolly.
func ApplyFramerate(c Command, d *Dolly) {
	d.Options.Framerate, _ = strconv.ParseFloat(c.Args, 64)
}

// ApplyOutput applies the output on the dolly GIF.
func ApplyOutput(c Command, d *Dolly) {
	d.Options.GIF.Output = c.Args
}
