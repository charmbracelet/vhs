package dolly

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

type Command struct {
	Type    CommandType
	Options string
	Args    string
}
