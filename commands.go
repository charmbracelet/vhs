package dolly

import (
	"fmt"
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
		fontSize, _ := strconv.Atoi(c.Args[len("FontSize "):])
		d.Options.TTY.FontSize = fontSize
		d.Page.Eval(fmt.Sprintf("term.setOption('fontSize', '%d')", fontSize))
	} else if strings.HasPrefix(c.Args, "FontFamily") {
		fontFamily := c.Args[len("FontFamily "):]
		d.Options.TTY.FontFamily = fontFamily
		d.Page.Eval(fmt.Sprintf("term.setOption('fontFamily', '%s')", fontFamily))
	} else if strings.HasPrefix(c.Args, "Width") {
		d.Options.Width, _ = strconv.Atoi(c.Args[(len("Width ")):])
	} else if strings.HasPrefix(c.Args, "Height") {
		d.Options.Height, _ = strconv.Atoi(c.Args[(len("Height ")):])
	} else if strings.HasPrefix(c.Args, "LineHeight") {
		lineHeight, _ := strconv.ParseFloat(c.Args[(len("LineHeight ")):], 64)
		d.Options.TTY.LineHeight = lineHeight
		d.Page.Eval(fmt.Sprintf("term.setOption('lineHeight', '%f')", lineHeight))
	} else if strings.HasPrefix(c.Args, "Theme") {
		theme := c.Args[len("Theme "):]
		d.Page.Eval(fmt.Sprintf("term.setOption('theme', '%s')", theme))
	} else if strings.HasPrefix(c.Args, "Padding") {
		padding := c.Args[(len("Padding ")):]
		d.Page.MustElement(".xterm").Eval(fmt.Sprintf(`this.style.padding = '%s'`, padding))
		d.Options.Padding = padding
	} else if strings.HasPrefix(c.Args, "Framerate") {
		d.Options.Framerate, _ = strconv.ParseFloat(c.Args[(len("Framerate ")):], 64)
	}
}
