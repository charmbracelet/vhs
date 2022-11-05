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
	Events []Event
}

func (a Asciicast) ToTape() string {
	settings := make([]string, 0)
	if shell, ok := a.Meta.Env["SHELL"]; ok && shell != "bash" {
		settings = append(settings, fmt.Sprintf("Set Shell %v", shell))
	}

	input := make([]string, 0)
	for _, e := range a.Events {
		input = append(input, e.Content)
	}
	return strings.Join(settings, "\n") + "\n" + inputToTape(strings.Join(input, ""))

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

type Event struct {
	Timestamp float64
	Type      string
	Content   string
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

		fmt.Println(cast.ToTape())

		return nil
	},
}