package dolly

import (
	"fmt"
	"strings"
)

const commentPrefix = "#"
const optionsPrefix = "@"

// Parse takes a string as input and returns the commands to be executed.
func Parse(s string) ([]Command, error) {
	var commands []Command

	lines := strings.Split(s, "\n")

	for _, line := range lines {
		if shouldSkip(line) {
			continue
		}

		for commandType, command := range Commands {
			if strings.HasPrefix(line, command) {
				options, args, err := parseArgs(command, line)
				if err != nil {
					return nil, err
				}
				commands = append(commands, Command{commandType, options, args})
				break
			}
		}
	}

	return commands, nil
}

func parseArgs(command string, line string) (string, string, error) {
	rawArgs := line[len(command):]
	if !strings.HasPrefix(rawArgs, optionsPrefix) {
		return "", strings.TrimPrefix(rawArgs, " "), nil
	}

	var options, args string
	splitIndex := strings.Index(rawArgs, " ")

	if splitIndex < 0 || splitIndex == len(rawArgs)-1 {
		return "", "", fmt.Errorf("no arguments found for %s", command)
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
