package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/go-rod/rod/lib/input"
)

// CommandType is a type that represents a command.
type CommandType TokenType

// CommandTypes is a list of the available commands that can be executed.
var CommandTypes = []CommandType{ //nolint: deadcode
	BACKSPACE,
	CTRL,
	DOWN,
	ENTER,
	ESCAPE,
	ILLEGAL,
	LEFT,
	RIGHT,
	SET,
	OUTPUT,
	SLEEP,
	SPACE,
	HIDE,
	REQUIRE,
	SHOW,
	TAB,
	TYPE,
	UP,
}

// String returns the string representation of the command.
func (c CommandType) String() string {
	if len(c) < 1 {
		return ""
	}
	s := string(c)
	return string(s[0]) + strings.ToLower(s[1:])
}

// CommandFunc is a function that executes a command on a running
// instance of vhs.
type CommandFunc func(c Command, v *VHS)

// CommandFuncs maps command types to their executable functions.
var CommandFuncs = map[CommandType]CommandFunc{
	BACKSPACE: ExecuteKey(input.Backspace),
	DOWN:      ExecuteKey(input.ArrowDown),
	ENTER:     ExecuteKey(input.Enter),
	LEFT:      ExecuteKey(input.ArrowLeft),
	RIGHT:     ExecuteKey(input.ArrowRight),
	SPACE:     ExecuteKey(input.Space),
	UP:        ExecuteKey(input.ArrowUp),
	TAB:       ExecuteKey(input.Tab),
	ESCAPE:    ExecuteKey(input.Escape),
	HIDE:      ExecuteHide,
	REQUIRE:   ExecuteRequire,
	SHOW:      ExecuteShow,
	SET:       ExecuteSet,
	OUTPUT:    ExecuteOutput,
	SLEEP:     ExecuteSleep,
	TYPE:      ExecuteType,
	CTRL:      ExecuteCtrl,
	ILLEGAL:   ExecuteNoop,
}

// Command represents a command with options and arguments.
type Command struct {
	Type    CommandType
	Options string
	Args    string
}

// String returns the string representation of the command.
// This includes the options and arguments of the command.
func (c Command) String() string {
	if c.Options != "" {
		return fmt.Sprintf("%s %s %s", c.Type, c.Options, c.Args)
	}
	return fmt.Sprintf("%s %s", c.Type, c.Options)
}

// Execute executes a command on a running instance of vhs.
func (c Command) Execute(v *VHS) {
	CommandFuncs[c.Type](c, v)
	if v.recording && v.Options.Test.Output != "" {
		v.SaveOutput()
	}
}

// ExecuteNoop is a no-op command that does nothing.
// Generally, this is used for Unknown commands when dealing with
// commands that are not recognized.
func ExecuteNoop(c Command, v *VHS) {}

// ExecuteKey is a higher-order function that returns a CommandFunc to execute
// a key press for a given key. This is so that the logic for key pressing
// (since they are repeatable and delayable) can be re-used.
//
// i.e. ExecuteKey(input.ArrowDown) would return a CommandFunc that executes
// the ArrowDown key press.
func ExecuteKey(k input.Key) CommandFunc {
	return func(c Command, v *VHS) {
		typingSpeed, err := time.ParseDuration(c.Options)
		if err != nil {
			typingSpeed = v.Options.TypingSpeed
		}
		repeat, err := strconv.Atoi(c.Args)
		if err != nil {
			repeat = 1
		}
		for i := 0; i < repeat; i++ {
			_ = v.Page.Keyboard.Type(k)
			time.Sleep(typingSpeed)
		}
	}
}

// ExecuteCtrl is a CommandFunc that presses the argument key with the ctrl key
// held down on the running instance of vhs.
func ExecuteCtrl(c Command, v *VHS) {
	_ = v.Page.Keyboard.Press(input.ControlLeft)
	for _, r := range c.Args {
		if k, ok := keymap[r]; ok {
			_ = v.Page.Keyboard.Type(k)
		}
	}
	_ = v.Page.Keyboard.Release(input.ControlLeft)
}

// ExecuteHide is a CommandFunc that starts or stops the recording of the vhs.
func ExecuteHide(c Command, v *VHS) {
	v.PauseRecording()
}

// ExecuteRequire is a CommandFunc that checks if all the binaries mentioned in the
// Require command are present. If not, it exits with a non-zero error.
func ExecuteRequire(c Command, v *VHS) {
	_, err := exec.LookPath(c.Args)
	if err != nil {
		v.Errors = append(v.Errors, err)
	}
}

// ExecuteShow is a CommandFunc that resumes the recording of the vhs.
func ExecuteShow(c Command, v *VHS) {
	v.ResumeRecording()
}

// ExecuteSleep sleeps for the desired time specified through the argument of
// the Sleep command.
func ExecuteSleep(c Command, v *VHS) {
	dur, err := time.ParseDuration(c.Args)
	if err != nil {
		return
	}
	time.Sleep(dur)
}

// ExecuteType types the argument string on the running instance of vhs.
func ExecuteType(c Command, v *VHS) {
	typingSpeed, err := time.ParseDuration(c.Options)
	if err != nil {
		typingSpeed = v.Options.TypingSpeed
	}
	for _, r := range c.Args {
		k, ok := keymap[r]
		if ok {
			_ = v.Page.Keyboard.Type(k)
		} else {
			_ = v.Page.MustElement("textarea").Input(string(r))
			v.Page.MustWaitIdle()
		}
		time.Sleep(typingSpeed)
	}
}

// ExecuteOutput applies the output on the vhs videos.
func ExecuteOutput(c Command, v *VHS) {
	switch c.Options {
	case ".mp4":
		v.Options.Video.Output.MP4 = c.Args
	case ".test", ".ascii", ".txt":
		v.Options.Test.Output = c.Args
	case ".png":
		v.Options.Video.Input = c.Args
		v.Options.Video.CleanupFrames = false
	case ".webm":
		v.Options.Video.Output.WebM = c.Args
	default:
		v.Options.Video.Output.GIF = c.Args
	}
}

// Settings maps the Set commands to their respective functions.
var Settings = map[string]CommandFunc{
	"FontFamily":    ExecuteSetFontFamily,
	"FontSize":      ExecuteSetFontSize,
	"Framerate":     ExecuteSetFramerate,
	"Height":        ExecuteSetHeight,
	"LetterSpacing": ExecuteSetLetterSpacing,
	"LineHeight":    ExecuteSetLineHeight,
	"PlaybackSpeed": ExecuteSetPlaybackSpeed,
	"Padding":       ExecuteSetPadding,
	"Theme":         ExecuteSetTheme,
	"TypingSpeed":   ExecuteSetTypingSpeed,
	"Width":         ExecuteSetWidth,
	"LoopOffset":    ExecuteLoopOffset,
}

// ExecuteSet applies the settings on the running vhs specified by the
// option and argument pass to the command.
func ExecuteSet(c Command, v *VHS) {
	Settings[c.Options](c, v)
}

// ExecuteSetFontSize applies the font size on the vhs.
func ExecuteSetFontSize(c Command, v *VHS) {
	fontSize, _ := strconv.Atoi(c.Args)
	v.Options.FontSize = fontSize
	_, _ = v.Page.Eval(fmt.Sprintf("() => term.options.fontSize = %d", fontSize))

	// When changing the font size only the canvas dimensions change which are
	// scaled back during the render to fit the aspect ration and dimensions.
	//
	// We need to call term.fit to ensure that everything is resized properly.
	_, _ = v.Page.Eval("term.fit")
}

// ExecuteSetFontFamily applies the font family on the vhs.
func ExecuteSetFontFamily(c Command, v *VHS) {
	v.Options.FontFamily = c.Args
	_, _ = v.Page.Eval(fmt.Sprintf("() => term.options.fontFamily = '%s'", c.Args))
}

// ExecuteSetHeight applies the height on the vhs.
func ExecuteSetHeight(c Command, v *VHS) {
	v.Options.Video.Height, _ = strconv.Atoi(c.Args)
}

// ExecuteSetWidth applies the width on the vhs.
func ExecuteSetWidth(c Command, v *VHS) {
	v.Options.Video.Width, _ = strconv.Atoi(c.Args)
}

const (
	bitSize = 64
	base    = 10
)

// ExecuteSetLetterSpacing applies letter spacing (also known as tracking) on the
// vhs.
func ExecuteSetLetterSpacing(c Command, v *VHS) {
	letterSpacing, _ := strconv.ParseFloat(c.Args, bitSize)
	v.Options.LetterSpacing = letterSpacing
	_, _ = v.Page.Eval(fmt.Sprintf("() => term.options.letterSpacing = %f", letterSpacing))
}

// ExecuteSetLineHeight applies the line height on the vhs.
func ExecuteSetLineHeight(c Command, v *VHS) {
	lineHeight, _ := strconv.ParseFloat(c.Args, bitSize)
	v.Options.LineHeight = lineHeight
	_, _ = v.Page.Eval(fmt.Sprintf("() => term.options.lineHeight = %f", lineHeight))
}

// ExecuteSetTheme applies the theme on the vhs.
func ExecuteSetTheme(c Command, v *VHS) {
	var err error
	v.Options.Theme, err = getTheme(c.Args)
	if err != nil {
		v.Errors = append(v.Errors, err)
		return
	}

	bts, _ := json.Marshal(v.Options.Theme)
	_, _ = v.Page.Eval(fmt.Sprintf("() => term.options.theme = %s", string(bts)))
	v.Options.Video.BackgroundColor = v.Options.Theme.Background
}

// ExecuteSetTypingSpeed applies the default typing speed on the vhs.
func ExecuteSetTypingSpeed(c Command, v *VHS) {
	typingSpeed, err := time.ParseDuration(c.Args)
	if err != nil {
		return
	}
	v.Options.TypingSpeed = typingSpeed
}

// ExecuteSetPadding applies the padding on the vhs.
func ExecuteSetPadding(c Command, v *VHS) {
	v.Options.Video.Padding, _ = strconv.Atoi(c.Args)
}

// ExecuteSetFramerate applies the framerate on the vhs.
func ExecuteSetFramerate(c Command, v *VHS) {
	framerate, err := strconv.ParseInt(c.Args, base, 0)
	if err != nil {
		return
	}
	v.Options.Video.Framerate = int(framerate)
}

// ExecuteSetPlaybackSpeed applies the playback speed option on the vhs.
func ExecuteSetPlaybackSpeed(c Command, v *VHS) {
	playbackSpeed, err := strconv.ParseFloat(c.Args, bitSize)
	if err != nil {
		return
	}
	v.Options.Video.PlaybackSpeed = playbackSpeed
}

// ExecuteLoopOffset applies the loop offset option on the vhs.
func ExecuteLoopOffset(c Command, v *VHS) {
	loopOffset, err := strconv.ParseFloat(strings.TrimRight(c.Args, "%"), bitSize)
	if err != nil {
		return
	}
	v.Options.LoopOffset = loopOffset
}

func getTheme(s string) (Theme, error) {
	if strings.TrimSpace(s) == "" {
		return DefaultTheme, nil
	}
	switch s[0] {
	case '{':
		return getJSONTheme(s)
	default:
		return findTheme(s)
	}
}

func getJSONTheme(s string) (Theme, error) {
	var t Theme
	if err := json.Unmarshal([]byte(s), &t); err != nil {
		return DefaultTheme, fmt.Errorf("invalid `Set Theme %q: %w`", s, err)
	}
	return t, nil
}
