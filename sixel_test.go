package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/charmbracelet/vhs/lexer"
	"github.com/charmbracelet/vhs/parser"
)

func TestSixelSetting(t *testing.T) {
	v := VHS{Options: &Options{}}
	if v.Options.Video.Sixel || v.Options.Screenshot.sixel {
		t.Fatal("sixel must default to disabled")
	}
	for _, value := range []string{"true", "false"} {
		p := parser.New(lexer.New("Set Sixel " + value))
		commands := p.Parse()
		if len(p.Errors()) != 0 || len(commands) != 1 {
			t.Fatalf("parse %s: %v", value, p.Errors())
		}
		requireNoErr(t, Execute(commands[0], &v))
		if v.Options.Video.Sixel != (value == "true") || v.Options.Screenshot.sixel != (value == "true") {
			t.Fatalf("setting %s not applied to both outputs", value)
		}
	}
	for _, value := range []string{"Enabled", "1", ""} {
		p := parser.New(lexer.New("Set Sixel " + value))
		p.Parse()
		if len(p.Errors()) == 0 {
			t.Fatalf("accepted invalid boolean %q", value)
		}
	}
}

func TestSixelInputs(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		t.Run(fmt.Sprint(enabled), func(t *testing.T) {
			dir := t.TempDir()
			style := DefaultStyleOptions()
			style.WindowBar = "Colorful"
			style.BorderRadius = 8
			opts := VideoOptions{Input: dir, Style: style, Framerate: 60, StartingFrame: 7, PlaybackSpeed: 1, Sixel: enabled}
			shot := NewScreenshotOptions(dir, style)
			shot.sixel = enabled
			for _, args := range [][]string{
				buildFFopts(opts, "out.gif"),
				shot.buildFFopts("out.png", filepath.Join(dir, textFrameFormat), filepath.Join(dir, cursorFrameFormat), filepath.Join(dir, imageFrameFormat)),
			} {
				var inputs []string
				var filter string
				for i, arg := range args {
					if arg == "-i" {
						inputs = append(inputs, args[i+1])
					}
					if arg == "-filter_complex" {
						filter = args[i+1]
					}
				}
				want := []string{filepath.Join(dir, textFrameFormat), filepath.Join(dir, cursorFrameFormat)}
				margin := 2
				if enabled {
					want = append(want, filepath.Join(dir, imageFrameFormat))
					margin++
					if !strings.Contains(filter, "[0][2]overlay[withimg];[withimg][1]overlay[merged]") {
						t.Fatalf("incorrect image overlay: %s", filter)
					}
				} else if !strings.Contains(filter, "[0][1]overlay[merged]") || strings.Contains(filter, "withimg") {
					t.Fatalf("default overlay changed: %s", filter)
				}
				want = append(want, "color="+style.MarginFill+":s=1200x600", filepath.Join(dir, "bar.png"), filepath.Join(dir, "mask.png"))
				if !reflect.DeepEqual(inputs, want) {
					t.Fatalf("inputs: %v; want %v", inputs, want)
				}
				for _, stage := range []string{fmt.Sprintf("[%d]scale=", margin), fmt.Sprintf("[%d]loop=", margin+1), fmt.Sprintf("[%d]", margin+2)} {
					if !strings.Contains(filter, stage) {
						t.Fatalf("missing decoration %s in %s", stage, filter)
					}
				}
			}
		})
	}
}

func TestSixelLoopOffset(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		t.Run(fmt.Sprint(enabled), func(t *testing.T) {
			dir := t.TempDir()
			formats := []string{textFrameFormat, cursorFrameFormat}
			if enabled {
				formats = append(formats, imageFrameFormat)
			}
			for frame := 1; frame <= 4; frame++ {
				for _, format := range formats {
					name := fmt.Sprintf(format, frame)
					requireNoErr(t, os.WriteFile(filepath.Join(dir, name), []byte(name), 0o600))
				}
			}
			v := VHS{totalFrames: 4, Options: &Options{
				LoopOffset: 50,
				Video:      VideoOptions{Input: dir, StartingFrame: 1, Sixel: enabled},
				Screenshot: ScreenshotOptions{
					input: dir, style: DefaultStyleOptions(), sixel: enabled,
					screenshots: map[string]int{"early.png": 1, "late.png": 4},
				},
			}}
			requireNoErr(t, v.ApplyLoopOffset())
			if v.Options.Video.StartingFrame != 3 || v.Options.Screenshot.screenshots["early.png"] != 5 || v.Options.Screenshot.screenshots["late.png"] != 4 {
				t.Fatal("incorrect frame references after loop offset")
			}
			for index, original := range []int{3, 4, 1, 2} {
				for _, format := range formats {
					data, err := os.ReadFile(filepath.Join(dir, fmt.Sprintf(format, index+3)))
					requireNoErr(t, err)
					if string(data) != fmt.Sprintf(format, original) {
						t.Fatalf("incorrect frame: %s", data)
					}
				}
			}
			for _, cmd := range MakeScreenshots(context.Background(), v.Options.Screenshot) {
				for i, arg := range cmd.Args {
					if arg == "-i" && strings.Contains(cmd.Args[i+1], "frame-") {
						_, err := os.Stat(cmd.Args[i+1])
						requireNoErr(t, err)
					}
				}
			}
		})
	}
}
