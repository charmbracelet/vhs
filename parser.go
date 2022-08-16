// ./dolly.vhs
//
// Type@100 echo "Hi, there!"
// Left 3
// Right 2
// Enter
// Sleep 1s
//
package dolly

import (
	"fmt"
	"strings"
)

const commentPrefix = "#"
const optionsPrefix = "@"

type CommandType string

const (
	Backspace CommandType = "Backspace"
	Down      CommandType = "Down"
	Enter     CommandType = "Enter"
	Left      CommandType = "Left"
	Right     CommandType = "Right"
	Sleep     CommandType = "Sleep"
	Space     CommandType = "Space"
	Type      CommandType = "Type"
	Up        CommandType = "Up"
)

var allCommands = []CommandType{Type, Down, Enter, Left, Right, Sleep, Up}

type Command struct {
	Type    CommandType
	Options string
	Args    string
}

// Parse takes a string as input and returns the commands to be executed.
func Parse(s string) ([]Command, error) {
	var commands []Command

	lines := strings.Split(s, "\n")

	for _, line := range lines {
		if shouldSkip(line) {
			continue
		}

		for _, command := range allCommands {
			if strings.HasPrefix(line, string(command)) {
				options, args, err := parseArgs(command, line)
				if err != nil {
					return nil, err
				}
				commands = append(commands, Command{command, options, args})
				break
			}
		}
	}

	return commands, nil
}

func parseArgs(commandType CommandType, line string) (string, string, error) {
	rawArgs := line[len(commandType):]
	if !strings.HasPrefix(rawArgs, optionsPrefix) {
		return "", strings.TrimPrefix(rawArgs, " "), nil
	}

	var options, args string
	splitIndex := strings.Index(rawArgs, " ")

	if splitIndex < 0 || splitIndex == len(rawArgs)-1 {
		return "", "", fmt.Errorf("no arguments found for %s", commandType)
	}

	options = rawArgs[:splitIndex]

	if splitIndex >= len(rawArgs) {
		return options, "", nil
	}

	args = rawArgs[splitIndex+1:]
	return options, args, nil
}

func shouldSkip(line string) bool {
	return strings.HasPrefix(line, commentPrefix) || strings.TrimSpace(line) == ""
}
