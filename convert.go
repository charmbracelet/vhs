package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
)

type Asciicast struct {
	Meta struct {
		Version   int
		Width     int
		Height    int
		Timestamp int // TODO timestamp format
		Env       map[string]string
	}
	Events EventStream
}

func (a Asciicast) ToTape() string {
	tape := make([]string, 0)
	if shell, ok := a.Meta.Env["SHELL"]; ok && shell != "bash" {
		tape = append(tape, fmt.Sprintf("Set Shell %v", shell))
	}
	tape = append(tape, a.Events.toTapeCommands()...)
	return strings.Join(tape, "\n")
}

func ReadFile(path string) (*Asciicast, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	asciicast := &Asciicast{
		Events: make([]Event, 0),
	}

	// header
	scanner.Scan()
	if err := json.Unmarshal(scanner.Bytes(), &asciicast.Meta); err != nil {
		return nil, err
	}

	// events
	for scanner.Scan() {
		var e Event
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			return nil, err
		} else if e.Type == "i" {
			asciicast.Events = append(asciicast.Events, e)
		}
	}
	return asciicast, nil

}

type EventStream []Event

func (es EventStream) traverse(f func(events []Event)) {
	buffer := make([]Event, 0)
	for _, event := range []Event(es) {
		if len(buffer) != 0 && (event.IsKey() || buffer[0].IsKey()) && event.Key() != buffer[0].Key() {
			f(buffer)
			buffer = make([]Event, 0) // TODO better way to clear it
		}
		buffer = append(buffer, event)
	}
	if len(buffer) != 0 {
		f(buffer)
	}
}

func (es EventStream) toTapeCommands() []string {
	commands := make([]string, 0)
	es.traverse(func(events []Event) {
		if events[0].IsKey() {
			if length := len(events); length > 0 {
				commands = append(commands, fmt.Sprintf("%v %v", events[0].Key(), length))
			} else {
				commands = append(commands, events[0].Key())
			}
		} else {
			s := make([]string, 0, len(events))
			for _, e := range events {
				s = append(s, e.Content)
			}
			commands = append(commands, fmt.Sprintf(`Type "%v"`, strings.ReplaceAll(strings.Join(s, ""), `"`, `\""`)))
		}
	})
	return commands
}

type Event struct {
	Timestamp float64
	Type      string
	Content   string
}

func (e Event) IsKey() bool {
	return e.Content != e.Key()
}

func (e Event) Key() string {
	switch e.Content {
	case "\x7f":
		return "Backspace"
	case "\r":
		return "Enter"
	case "\t":
		return "Tab"
	default:
		return e.Content
	}
}

func (e *Event) UnmarshalJSON(data []byte) error {
	var raw []interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	} else if len(raw) != 3 {
		return errors.New("invalid amount of elements")
	}

	if value := reflect.ValueOf(raw[0]); value.Kind() == reflect.Float64 {
		e.Timestamp = value.Float()
	} else {
		return errors.New("invalid format for Timestamp")
	}

	if value := reflect.ValueOf(raw[1]); value.Kind() == reflect.String {
		e.Type = value.String()
	} else {
		return errors.New("invalid format for Type")
	}

	if value := reflect.ValueOf(raw[2]); value.Kind() == reflect.String {
		e.Content = value.String()
	} else {
		return errors.New("invalid format for Content")
	}
	return nil
}

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "TODO",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cast, err := ReadFile(args[0])
		if err != nil {
			return err
		}
		// TODO fmt.Printf("%#v\n", cast)

		fmt.Println(cast.ToTape())
		return nil
	},
}