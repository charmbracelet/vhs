package dolly

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMissingArguments = errors.New("missing arguments")
	ErrUnknownCommand   = errors.New("unknown command")
	ErrUnknownOptions   = errors.New("unknown options")
)

const (
	commentPrefix = "#"
	optionsPrefix = "@"
)

// Parse takes a string as input and returns the commands to be executed.
func Parse(s string) ([]Command, error) {
	var commands []Command

	lines := strings.Split(s, "\n")

	for _, line := range lines {
		if shouldSkip(line) {
			continue
		}

		valid := false
		for commandType, command := range Commands {
			if strings.HasPrefix(line, command) {
				options, args, err := parseArgs(command, line)
				if err != nil {
					return nil, err
				}
				commands = append(commands, Command{commandType, options, args})
				valid = true
				break
			}
		}
		if !valid {
			return nil, fmt.Errorf("%s: %s", ErrUnknownCommand, line)
		}
	}

	return commands, nil
}

func parseArgs(command string, line string) (string, string, error) {
	rawArgs := strings.TrimPrefix(line[len(command):], " ")

	if command == Commands[Set] {
		splitIndex := strings.Index(rawArgs, " ")
		if splitIndex == -1 {
			return "", "", fmt.Errorf("%s: %s", ErrMissingArguments, line)
		}

		options := rawArgs[:splitIndex]
		args := rawArgs[splitIndex+1:]
		_, ok := SetCommands[options]
		if !ok {
			return "", "", fmt.Errorf("%s: %s", ErrUnknownOptions, line)
		}

		return options, args, nil
	}

	if !strings.HasPrefix(rawArgs, optionsPrefix) {
		if command == Commands[Type] && rawArgs == "" {
			return "", "", fmt.Errorf("%s: %s", ErrMissingArguments, line)
		}
		return "", rawArgs, nil
	}

	var options, args string
	splitIndex := strings.Index(rawArgs, " ")

	if splitIndex < 0 || splitIndex == len(rawArgs)-1 {
		return "", "", fmt.Errorf("%s: %s", ErrMissingArguments, line)
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
