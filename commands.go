package dolly

import (
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
	Sleep
	Space
	Type
	Up
)

var Commands = map[CommandType]string{
	Backspace: "Backspace",
	Down:      "Down",
	Enter:     "Enter",
	Left:      "Left",
	Right:     "Right",
	Sleep:     "Sleep",
	Space:     "Space",
	Type:      "Type",
	Up:        "Up",
}

type CommandFunc func(c Command, d *Dolly)

var CommandFuncs = map[CommandType]CommandFunc{
	Backspace: ExecuteBackspace,
	Down:      ExecuteDown,
	Enter:     ExecuteEnter,
	Left:      ExecuteLeft,
	Right:     ExecuteRight,
	Sleep:     ExecuteSleep,
	Space:     ExecuteSpace,
	Type:      ExecuteType,
	Up:        ExecuteUp,
}

type Command struct {
	Type    CommandType
	Options string
	Args    string
}

func (c Command) Execute(d *Dolly) {
	CommandFuncs[c.Type](c, d)
}

func ExecuteBackspace(c Command, d *Dolly) {
	d.Page.Keyboard.Type(input.Backspace)
}

func ExecuteDown(c Command, d *Dolly) {
	d.Page.Keyboard.Type(input.ArrowDown)
}

func ExecuteEnter(c Command, d *Dolly) {
	d.Page.Keyboard.Type(input.Enter)
}

func ExecuteLeft(c Command, d *Dolly) {
	d.Page.Keyboard.Type(input.ArrowLeft)
}

func ExecuteRight(c Command, d *Dolly) {
	d.Page.Keyboard.Type(input.ArrowRight)
}

func ExecuteUp(c Command, d *Dolly) {
	d.Page.Keyboard.Type(input.ArrowUp)
}

func ExecuteSpace(c Command, d *Dolly) {
	d.Page.Keyboard.Type(input.Space)
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
