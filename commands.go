package dolly

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod/lib/input"
)

type CommandType int

const (
	Backspace CommandType = iota
	Down
	Enter
	Left
	Right
	Space
	Up

	Type
	Set
	Sleep
)

var Commands = map[CommandType]string{
	Backspace: "Backspace",
	Down:      "Down",
	Enter:     "Enter",
	Left:      "Left",
	Right:     "Right",
	Space:     "Space",
	Up:        "Up",

	Set:   "Set",
	Sleep: "Sleep",
	Type:  "Type",
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

	Set:   ExecuteSet,
	Sleep: ExecuteSleep,
	Type:  ExecuteType,
}

type Command struct {
	Type    CommandType
	Options string
	Args    string
}

func (c Command) Execute(d *Dolly) {
	CommandFuncs[c.Type](c, d)
}

func (c Command) String() string {
	return Commands[c.Type] + " " + c.Args
}

func ExecuteKey(k input.Key) CommandFunc {
	return func(c Command, d *Dolly) {
		num, err := strconv.Atoi(c.Args)
		if err != nil {
			num = 1
		}
		for i := 0; i < num; i++ {
			d.Page.Keyboard.Type(k)
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
		time.Sleep(time.Millisecond * 100)
	}
}

func ExecuteSet(c Command, d *Dolly) {
	if strings.HasPrefix(c.Args, "FontSize") {
		d.Options.TTY.FontSize, _ = strconv.Atoi(c.Args[len("FontSize "):])
	} else if strings.HasPrefix(c.Args, "FontFamily") {
		d.Options.TTY.FontFamily = c.Args[len("FontFamily "):]
	} else if strings.HasPrefix(c.Args, "Width") {
		d.Options.Width, _ = strconv.Atoi(c.Args[(len("Width ")):])
	} else if strings.HasPrefix(c.Args, "Height") {
		d.Options.Height, _ = strconv.Atoi(c.Args[(len("Height ")):])
	} else if strings.HasPrefix(c.Args, "LineHeight") {
		d.Options.TTY.LineHeight, _ = strconv.ParseFloat(c.Args[(len("LineHeight ")):], 64)
	} else if strings.HasPrefix(c.Args, "Padding") {
		d.Options.Padding = c.Args[(len("Padding ")):]
	} else if strings.HasPrefix(c.Args, "Framerate") {
		d.Options.Framerate, _ = strconv.ParseFloat(c.Args[(len("Framerate ")):], 64)
	}
}
