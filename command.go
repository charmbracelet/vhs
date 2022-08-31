package dolly

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-rod/rod/lib/input"
)

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
)

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
}

func (c CommandType) String() string {
	return string(c)
}

type CommandFunc func(c Command, d *Dolly)

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
}

type Command struct {
	Type    CommandType
	Options string
	Args    string
}

func (c Command) String() string {
	if c.Options != "" {
		return fmt.Sprintf("%s %s %s", c.Type, c.Options, c.Args)
	}
	return fmt.Sprintf("%s %s", c.Type, c.Args)
}

func (c Command) Execute(d *Dolly) {
	CommandFuncs[c.Type](c, d)
}

func ExecuteKey(k input.Key) CommandFunc {
	return func(c Command, d *Dolly) {
		repeat, err := strconv.Atoi(c.Args)
		if err != nil {
			repeat = 1
		}
		delayMs, err := strconv.Atoi(c.Options)
		if err != nil {
			delayMs = 100
		}
		for i := 0; i < repeat; i++ {
			d.Page.Keyboard.Type(k)
			time.Sleep(time.Millisecond * time.Duration(delayMs))
		}
	}
}

func ExecuteSleep(c Command, d *Dolly) {
	dur, err := time.ParseDuration(c.Args)
	if err != nil {
		return
	}
	time.Sleep(dur)
}

func ExecuteType(c Command, d *Dolly) {
	for _, r := range c.Args {
		k, ok := keymap[r]
		if ok {
			d.Page.Keyboard.Type(k)
		} else {
			d.Page.MustElement("textarea").Input(string(r))
			d.Page.MustWaitIdle()
		}
		delayMs, err := strconv.Atoi(c.Options)
		if err != nil {
			delayMs = 100
		}
		time.Sleep(time.Millisecond * time.Duration(delayMs))
	}
}

var SetCommands = map[string]CommandFunc{
	"FontSize":   ApplyFontSize,
	"FontFamily": ApplyFontFamily,
	"Height":     ApplyHeight,
	"Width":      ApplyWidth,
	"LineHeight": ApplyLineHeight,
	"Theme":      ApplyTheme,
	"Padding":    ApplyPadding,
	"Framerate":  ApplyFramerate,
}

func ExecuteSet(c Command, d *Dolly) {
	SetCommands[c.Options](c, d)
}

func ApplyFontSize(c Command, d *Dolly) {
	fontSize, _ := strconv.Atoi(c.Args)
	d.Options.FontSize = fontSize
	d.Page.Eval(fmt.Sprintf("term.setOption('fontSize', '%d')", fontSize))
}

func ApplyFontFamily(c Command, d *Dolly) {
	d.Options.FontFamily = c.Args
	d.Page.Eval(fmt.Sprintf("term.setOption('fontFamily', '%s')", c.Args))
}

func ApplyHeight(c Command, d *Dolly) {
	d.Options.Height, _ = strconv.Atoi(c.Args)
}

func ApplyWidth(c Command, d *Dolly) {
	d.Options.Width, _ = strconv.Atoi(c.Args)
}

func ApplyLineHeight(c Command, d *Dolly) {
	lineHeight, _ := strconv.ParseFloat(c.Args, 64)
	d.Options.LineHeight = lineHeight
	d.Page.Eval(fmt.Sprintf("term.setOption('lineHeight', '%f')", lineHeight))
}

func ApplyTheme(c Command, d *Dolly) {
	err := json.Unmarshal([]byte(c.Args), &d.Options.Theme)
	if err != nil {
		d.Options.Theme = DefaultTheme
		return
	}
	d.Page.Eval(fmt.Sprintf("term.setOption('theme', %s)", c.Args))
}

func ApplyPadding(c Command, d *Dolly) {
	d.Options.Padding = c.Args
	d.Page.MustElement(".xterm").Eval(fmt.Sprintf(`this.style.padding = '%s'`, c.Args))
}

func ApplyFramerate(c Command, d *Dolly) {
	d.Options.Framerate, _ = strconv.ParseFloat(c.Args, 64)
}
