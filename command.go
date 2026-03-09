package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/vhs/parser"
	"github.com/charmbracelet/vhs/token"
	"github.com/go-rod/rod/lib/input"
)

// Execute executes a command on a running instance of vhs.
func Execute(c parser.Command, v *VHS) error {
	err := CommandFuncs[c.Type](c, v)
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	if v.recording && v.Options.Test.Output != "" {
		err := v.SaveOutput()
		if err != nil {
			return fmt.Errorf("failed to save output: %w", err)
		}
	}

	return nil
}

// CommandFunc is a function that executes a command on a running
// instance of vhs.
type CommandFunc func(c parser.Command, v *VHS) error

// CommandFuncs maps command types to their executable functions.
var CommandFuncs = map[parser.CommandType]CommandFunc{
	token.BACKSPACE:  ExecuteKey(input.Backspace),
	token.DELETE:     ExecuteKey(input.Delete),
	token.INSERT:     ExecuteKey(input.Insert),
	token.DOWN:       ExecuteKey(input.ArrowDown),
	token.ENTER:      ExecuteKey(input.Enter),
	token.LEFT:       ExecuteKey(input.ArrowLeft),
	token.RIGHT:      ExecuteKey(input.ArrowRight),
	token.SPACE:      ExecuteKey(input.Space),
	token.UP:         ExecuteKey(input.ArrowUp),
	token.TAB:        ExecuteKey(input.Tab),
	token.ESCAPE:     ExecuteKey(input.Escape),
	token.PAGE_UP:    ExecuteKey(input.PageUp),
	token.PAGE_DOWN:  ExecuteKey(input.PageDown),
	token.HIDE:       ExecuteHide,
	token.REQUIRE:    ExecuteRequire,
	token.SHOW:       ExecuteShow,
	token.SET:        ExecuteSet,
	token.OUTPUT:     ExecuteOutput,
	token.SLEEP:      ExecuteSleep,
	token.TYPE:       ExecuteType,
	token.CTRL:       ExecuteCtrl,
	token.ALT:        ExecuteAlt,
	token.SHIFT:      ExecuteShift,
	token.ILLEGAL:    ExecuteNoop,
	token.SCREENSHOT: ExecuteScreenshot,
	token.COPY:       ExecuteCopy,
	token.PASTE:      ExecutePaste,
	token.ENV:        ExecuteEnv,
	token.AWAIT_PROMPT: ExecuteAwaitPrompt,
	token.WAIT:         ExecuteWait,
}

// ExecuteNoop is a no-op command that does nothing.
// Generally, this is used for Unknown commands when dealing with
// commands that are not recognized.
func ExecuteNoop(_ parser.Command, _ *VHS) error { return nil }

// ExecuteKey is a higher-order function that returns a CommandFunc to execute
// a key press for a given key. This is so that the logic for key pressing
// (since they are repeatable and delayable) can be re-used.
//
// i.e. ExecuteKey(input.ArrowDown) would return a CommandFunc that executes
// the ArrowDown key press.
func ExecuteKey(k input.Key) CommandFunc {
	return func(c parser.Command, v *VHS) error {
		typingSpeed, err := time.ParseDuration(c.Options)
		if err != nil {
			typingSpeed = v.Options.TypingSpeed
		}
		repeat, err := strconv.Atoi(c.Args)
		if err != nil {
			repeat = 1
		}
		for i := 0; i < repeat; i++ {
			err = v.Page.Keyboard.Type(k)
			if err != nil {
				return fmt.Errorf("failed to type key %c: %w", k, err)
			}
			time.Sleep(typingSpeed)
		}

		return nil
	}
}

// WaitTick is the amount of time to wait between checking for a match.
const WaitTick = 10 * time.Millisecond

// ExecuteWait is a CommandFunc that waits for a regex match for the given amount of time.
func ExecuteWait(c parser.Command, v *VHS) error {
	scope, rxStr, ok := strings.Cut(c.Args, " ")
	rx := v.Options.WaitPattern
	if ok {
		// This is validated on parse so using MustCompile reduces noise.
		rx = regexp.MustCompile(rxStr)
	}

	timeout := v.Options.WaitTimeout
	if c.Options != "" {
		t, err := time.ParseDuration(c.Options)
		if err != nil {
			// Shouldn't be possible due to parse validation.
			return fmt.Errorf("failed to parse duration: %w", err)
		}
		timeout = t
	}

	checkT := time.NewTicker(WaitTick)
	defer checkT.Stop()
	timeoutT := time.NewTimer(timeout)
	defer timeoutT.Stop()

	for {
		var last string
		switch scope {
		case "Line":
			line, err := v.CurrentLine()
			if err != nil {
				return fmt.Errorf("failed to get current line: %w", err)
			}
			last = line

			if rx.MatchString(line) {
				return nil
			}
		case "Screen":
			lines, err := v.Buffer()
			if err != nil {
				return fmt.Errorf("failed to get buffer: %w", err)
			}
			last = strings.Join(lines, "\n")

			if rx.MatchString(last) {
				return nil
			}
		default:
			// Should be impossible due to parse validation, but we don't want to
			// hang if it does happen due to a bug.
			return fmt.Errorf("invalid scope %q", scope)
		}

		select {
		case <-checkT.C:
			continue
		case <-timeoutT.C:
			return fmt.Errorf("timeout waiting for %q to match %s; last value was: %s", c.Args, rx.String(), last)
		}
	}
}

// ExecuteAwaitPrompt waits for the shell to emit a new prompt marker.
// It detects prompt markers (OSC 7777) that are embedded in each shell's prompt
// configuration. Unlike Wait (which matches terminal content), AwaitPrompt
// detects when the shell has finished executing a command and is ready for input.
func ExecuteAwaitPrompt(c parser.Command, v *VHS) error {
	timeout := v.Options.WaitTimeout
	if c.Options != "" {
		t, err := time.ParseDuration(c.Options)
		if err != nil {
			return fmt.Errorf("failed to parse duration: %w", err)
		}
		timeout = t
	}

	// Record the current prompt count so we can detect the next one.
	baseline, err := v.PromptCount()
	if err != nil {
		return fmt.Errorf("failed to read prompt count: %w", err)
	}

	checkT := time.NewTicker(WaitTick)
	defer checkT.Stop()
	timeoutT := time.NewTimer(timeout)
	defer timeoutT.Stop()

	for {
		select {
		case <-checkT.C:
			current, err := v.PromptCount()
			if err != nil {
				return fmt.Errorf("failed to read prompt count: %w", err)
			}
			if current > baseline {
				return nil
			}
		case <-timeoutT.C:
			return fmt.Errorf("timeout waiting for shell prompt (waited %s)", timeout)
		}
	}
}

// ExecuteCtrl is a CommandFunc that presses the argument keys and/or modifiers
// with the ctrl key held down on the running instance of vhs.
func ExecuteCtrl(c parser.Command, v *VHS) error {
	// Create key combination by holding ControlLeft
	action := v.Page.KeyActions().Press(input.ControlLeft)
	keys := strings.Split(c.Args, " ")

	for i, key := range keys {
		var inputKey *input.Key

		switch key {
		case "Shift":
			inputKey = &input.ShiftLeft
		case "Alt":
			inputKey = &input.AltLeft
		case "Enter":
			inputKey = &input.Enter
		case "Space":
			inputKey = &input.Space
		case "Backspace":
			inputKey = &input.Backspace
		default:
			r := rune(key[0])
			if k, ok := keymap[r]; ok {
				inputKey = &k
			}
		}

		// Press or hold key in case it's valid
		if inputKey != nil {
			if i != len(keys)-1 {
				action.Press(*inputKey)
			} else {
				// Other keys will remain pressed until the combination reaches the end
				action.Type(*inputKey)
			}
		}
	}

	err := action.Do()
	if err != nil {
		return fmt.Errorf("failed to type key %s: %w", c.Args, err)
	}

	return nil
}

// ExecuteAlt is a CommandFunc that presses the argument key with the alt key
// held down on the running instance of vhs.
func ExecuteAlt(c parser.Command, v *VHS) error {
	err := v.Page.Keyboard.Press(input.AltLeft)
	if err != nil {
		return fmt.Errorf("failed to press Alt key: %w", err)
	}
	if k, ok := token.Keywords[c.Args]; ok { //nolint:nestif
		switch k {
		case token.ENTER:
			err = v.Page.Keyboard.Type(input.Enter)
			if err != nil {
				return fmt.Errorf("failed to type Enter key: %w", err)
			}
		case token.TAB:
			err := v.Page.Keyboard.Type(input.Tab)
			if err != nil {
				return fmt.Errorf("failed to type Tab key: %w", err)
			}
		}
	} else {
		for _, r := range c.Args {
			if k, ok := keymap[r]; ok {
				err = v.Page.Keyboard.Type(k)
				if err != nil {
					return fmt.Errorf("failed to type key %c: %w", r, err)
				}
			}
		}
	}

	err = v.Page.Keyboard.Release(input.AltLeft)
	if err != nil {
		return fmt.Errorf("failed to release Alt key: %w", err)
	}

	return nil
}

// ExecuteShift is a CommandFunc that presses the argument key with the shift
// key held down on the running instance of vhs.
func ExecuteShift(c parser.Command, v *VHS) error {
	err := v.Page.Keyboard.Press(input.ShiftLeft)
	if err != nil {
		return fmt.Errorf("failed to press Shift key: %w", err)
	}

	if k, ok := token.Keywords[c.Args]; ok { //nolint:nestif
		switch k {
		case token.ENTER:
			err = v.Page.Keyboard.Type(input.Enter)
			if err != nil {
				return fmt.Errorf("failed to type Enter key: %w", err)
			}
		case token.TAB:
			err = v.Page.Keyboard.Type(input.Tab)
			if err != nil {
				return fmt.Errorf("failed to type Tab key: %w", err)
			}
		}
	} else {
		for _, r := range c.Args {
			if k, ok := keymap[r]; ok {
				err = v.Page.Keyboard.Type(k)
				if err != nil {
					return fmt.Errorf("failed to type key %c: %w", r, err)
				}
			}
		}
	}

	err = v.Page.Keyboard.Release(input.ShiftLeft)
	if err != nil {
		return fmt.Errorf("failed to release Shift key: %w", err)
	}

	return nil
}

// ExecuteHide is a CommandFunc that starts or stops the recording of the vhs.
func ExecuteHide(_ parser.Command, v *VHS) error {
	v.PauseRecording()
	return nil
}

// ExecuteRequire is a CommandFunc that checks if all the binaries mentioned in the
// Require command are present. If not, it exits with a non-zero error.
func ExecuteRequire(c parser.Command, _ *VHS) error {
	_, err := exec.LookPath(c.Args)
	return err //nolint:wrapcheck
}

// ExecuteShow is a CommandFunc that resumes the recording of the vhs.
func ExecuteShow(_ parser.Command, v *VHS) error {
	v.ResumeRecording()
	return nil
}

// ExecuteSleep sleeps for the desired time specified through the argument of
// the Sleep command.
func ExecuteSleep(c parser.Command, _ *VHS) error {
	dur, err := time.ParseDuration(c.Args)
	if err != nil {
		return fmt.Errorf("failed to parse duration: %w", err)
	}
	time.Sleep(dur)
	return nil
}

// ExecuteType types the argument string on the running instance of vhs.
func ExecuteType(c parser.Command, v *VHS) error {
	typingSpeed := v.Options.TypingSpeed
	if c.Options != "" {
		var err error
		typingSpeed, err = time.ParseDuration(c.Options)
		if err != nil {
			return fmt.Errorf("failed to parse typing speed: %w", err)
		}
	}
	for _, r := range c.Args {
		k, ok := keymap[r]
		if ok {
			err := v.Page.Keyboard.Type(k)
			if err != nil {
				return fmt.Errorf("failed to type key %c: %w", r, err)
			}
		} else {
			err := v.Page.MustElement("textarea").Input(string(r))
			if err != nil {
				return fmt.Errorf("failed to input text: %w", err)
			}

			v.Page.MustWaitIdle()
		}
		time.Sleep(typingSpeed)
	}

	return nil
}

// ExecuteOutput applies the output on the vhs videos.
func ExecuteOutput(c parser.Command, v *VHS) error {
	switch c.Options {
	case ".mp4":
		v.Options.Video.Output.MP4 = c.Args
	case ".test", ".ascii", ".txt":
		v.Options.Test.Output = c.Args
	case ".png":
		v.Options.Video.Output.Frames = c.Args
	case ".webm":
		v.Options.Video.Output.WebM = c.Args
	default:
		v.Options.Video.Output.GIF = c.Args
	}

	return nil
}

// ExecuteCopy copies text to the clipboard.
func ExecuteCopy(c parser.Command, _ *VHS) error {
	return clipboard.WriteAll(c.Args) //nolint:wrapcheck
}

// ExecuteEnv sets env with given key-value pair.
func ExecuteEnv(c parser.Command, _ *VHS) error {
	return os.Setenv(c.Options, c.Args) //nolint:wrapcheck
}

// ExecutePaste pastes text from the clipboard.
func ExecutePaste(_ parser.Command, v *VHS) error {
	clip, err := clipboard.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read clipboard: %w", err)
	}
	for _, r := range clip {
		k, ok := keymap[r]
		if ok {
			err = v.Page.Keyboard.Type(k)
			if err != nil {
				return fmt.Errorf("failed to type key %c: %w", r, err)
			}
		} else {
			err = v.Page.MustElement("textarea").Input(string(r))
			if err != nil {
				return fmt.Errorf("failed to input text: %w", err)
			}
			v.Page.MustWaitIdle()
		}
	}

	return nil
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
	"Shell":         ExecuteSetShell,
	"LoopOffset":    ExecuteLoopOffset,
	"MarginFill":    ExecuteSetMarginFill,
	"Margin":        ExecuteSetMargin,
	"WindowBar":     ExecuteSetWindowBar,
	"WindowBarSize": ExecuteSetWindowBarSize,
	"BorderRadius":  ExecuteSetBorderRadius,
	"WaitPattern":   ExecuteSetWaitPattern,
	"WaitTimeout":   ExecuteSetWaitTimeout,
	"CursorBlink":   ExecuteSetCursorBlink,
}

// ExecuteSet applies the settings on the running vhs specified by the
// option and argument pass to the command.
func ExecuteSet(c parser.Command, v *VHS) error {
	return Settings[c.Options](c, v)
}

// ExecuteSetFontSize applies the font size on the vhs.
func ExecuteSetFontSize(c parser.Command, v *VHS) error {
	fontSize, err := strconv.Atoi(c.Args)
	if err != nil {
		return fmt.Errorf("failed to parse font size: %w", err)
	}
	v.Options.FontSize = fontSize
	_, err = v.Page.Eval(fmt.Sprintf("() => term.options.fontSize = %d", fontSize))
	if err != nil {
		return fmt.Errorf("failed to set font size: %w", err)
	}

	// When changing the font size only the canvas dimensions change which are
	// scaled back during the render to fit the aspect ration and dimensions.
	//
	// We need to call term.fit to ensure that everything is resized properly.
	_, err = v.Page.Eval("term.fit")
	if err != nil {
		return fmt.Errorf("failed to fit terminal: %w", err)
	}

	return nil
}

// ExecuteSetFontFamily applies the font family on the vhs.
func ExecuteSetFontFamily(c parser.Command, v *VHS) error {
	v.Options.FontFamily = c.Args
	_, err := v.Page.Eval(fmt.Sprintf("() => term.options.fontFamily = '%s'", withSymbolsFallback(c.Args)))
	if err != nil {
		return fmt.Errorf("failed to set font family: %w", err)
	}

	return nil
}

// ExecuteSetHeight applies the height on the vhs.
func ExecuteSetHeight(c parser.Command, v *VHS) error {
	height, err := strconv.Atoi(c.Args)
	if err != nil {
		return fmt.Errorf("failed to parse height: %w", err)
	}
	v.Options.Video.Style.Height = height

	return nil
}

// ExecuteSetWidth applies the width on the vhs.
func ExecuteSetWidth(c parser.Command, v *VHS) error {
	width, err := strconv.Atoi(c.Args)
	if err != nil {
		return fmt.Errorf("failed to parse width: %w", err)
	}
	v.Options.Video.Style.Width = width

	return nil
}

// ExecuteSetShell applies the shell on the vhs.
func ExecuteSetShell(c parser.Command, v *VHS) error {
	s, ok := Shells[c.Args]
	if !ok {
		return fmt.Errorf("invalid shell %s", c.Args)
	}

	v.Options.Shell = s
	return nil
}

const (
	bitSize = 64
	base    = 10
)

// ExecuteSetLetterSpacing applies letter spacing (also known as tracking) on
// the vhs.
func ExecuteSetLetterSpacing(c parser.Command, v *VHS) error {
	letterSpacing, err := strconv.ParseFloat(c.Args, bitSize)
	if err != nil {
		return fmt.Errorf("failed to parse letter spacing: %w", err)
	}

	v.Options.LetterSpacing = letterSpacing
	_, err = v.Page.Eval(fmt.Sprintf("() => term.options.letterSpacing = %f", letterSpacing))
	if err != nil {
		return fmt.Errorf("failed to set letter spacing: %w", err)
	}

	return nil
}

// ExecuteSetLineHeight applies the line height on the vhs.
func ExecuteSetLineHeight(c parser.Command, v *VHS) error {
	lineHeight, err := strconv.ParseFloat(c.Args, bitSize)
	if err != nil {
		return fmt.Errorf("failed to parse line height: %w", err)
	}

	v.Options.LineHeight = lineHeight
	_, err = v.Page.Eval(fmt.Sprintf("() => term.options.lineHeight = %f", lineHeight))
	if err != nil {
		return fmt.Errorf("failed to set line height: %w", err)
	}

	return nil
}

// ExecuteSetTheme applies the theme on the vhs.
func ExecuteSetTheme(c parser.Command, v *VHS) error {
	var err error
	v.Options.Theme, err = getTheme(c.Args)
	if err != nil {
		return err
	}

	bts, err := json.Marshal(v.Options.Theme)
	if err != nil {
		return fmt.Errorf("failed to marshal theme: %w", err)
	}

	_, err = v.Page.Eval(fmt.Sprintf("() => term.options.theme = %s", string(bts)))
	if err != nil {
		return fmt.Errorf("failed to set theme: %w", err)
	}

	v.Options.Video.Style.BackgroundColor = v.Options.Theme.Background
	v.Options.Video.Style.WindowBarColor = v.Options.Theme.Background

	return nil
}

// ExecuteSetTypingSpeed applies the default typing speed on the vhs.
func ExecuteSetTypingSpeed(c parser.Command, v *VHS) error {
	typingSpeed, err := time.ParseDuration(c.Args)
	if err != nil {
		return fmt.Errorf("failed to parse typing speed: %w", err)
	}

	v.Options.TypingSpeed = typingSpeed
	return nil
}

// ExecuteSetWaitTimeout applies the default wait timeout on the vhs.
func ExecuteSetWaitTimeout(c parser.Command, v *VHS) error {
	waitTimeout, err := time.ParseDuration(c.Args)
	if err != nil {
		return fmt.Errorf("failed to parse wait timeout: %w", err)
	}
	v.Options.WaitTimeout = waitTimeout
	return nil
}

// ExecuteSetWaitPattern applies the default wait pattern on the vhs.
func ExecuteSetWaitPattern(c parser.Command, v *VHS) error {
	rx, err := regexp.Compile(c.Args)
	if err != nil {
		return fmt.Errorf("failed to compile regexp: %w", err)
	}
	v.Options.WaitPattern = rx
	return nil
}

// ExecuteSetPadding applies the padding on the vhs.
func ExecuteSetPadding(c parser.Command, v *VHS) error {
	padding, err := strconv.Atoi(c.Args)
	if err != nil {
		return fmt.Errorf("failed to parse padding: %w", err)
	}

	v.Options.Video.Style.Padding = padding
	return nil
}

// ExecuteSetFramerate applies the framerate on the vhs.
func ExecuteSetFramerate(c parser.Command, v *VHS) error {
	framerate, err := strconv.ParseInt(c.Args, base, 0)
	if err != nil {
		return fmt.Errorf("failed to parse framerate: %w", err)
	}

	v.Options.Video.Framerate = int(framerate)
	return nil
}

// ExecuteSetPlaybackSpeed applies the playback speed option on the vhs.
func ExecuteSetPlaybackSpeed(c parser.Command, v *VHS) error {
	playbackSpeed, err := strconv.ParseFloat(c.Args, bitSize)
	if err != nil {
		return fmt.Errorf("failed to parse playback speed: %w", err)
	}

	v.Options.Video.PlaybackSpeed = playbackSpeed
	return nil
}

// ExecuteLoopOffset applies the loop offset option on the vhs.
func ExecuteLoopOffset(c parser.Command, v *VHS) error {
	loopOffset, err := strconv.ParseFloat(strings.TrimRight(c.Args, "%"), bitSize)
	if err != nil {
		return fmt.Errorf("failed to parse loop offset: %w", err)
	}

	v.Options.LoopOffset = loopOffset
	return nil
}

// ExecuteSetMarginFill sets vhs margin fill.
func ExecuteSetMarginFill(c parser.Command, v *VHS) error {
	v.Options.Video.Style.MarginFill = c.Args
	return nil
}

// ExecuteSetMargin sets vhs margin size.
func ExecuteSetMargin(c parser.Command, v *VHS) error {
	margin, err := strconv.Atoi(c.Args)
	if err != nil {
		return fmt.Errorf("failed to parse margin: %w", err)
	}

	v.Options.Video.Style.Margin = margin
	return nil
}

// ExecuteSetWindowBar sets window bar type.
func ExecuteSetWindowBar(c parser.Command, v *VHS) error {
	v.Options.Video.Style.WindowBar = c.Args
	return nil
}

// ExecuteSetWindowBarSize sets window bar size.
func ExecuteSetWindowBarSize(c parser.Command, v *VHS) error {
	windowBarSize, err := strconv.Atoi(c.Args)
	if err != nil {
		return fmt.Errorf("failed to parse window bar size: %w", err)
	}

	v.Options.Video.Style.WindowBarSize = windowBarSize
	return nil
}

// ExecuteSetBorderRadius sets corner radius.
func ExecuteSetBorderRadius(c parser.Command, v *VHS) error {
	borderRadius, err := strconv.Atoi(c.Args)
	if err != nil {
		return fmt.Errorf("failed to parse border radius: %w", err)
	}

	v.Options.Video.Style.BorderRadius = borderRadius
	return nil
}

// ExecuteSetCursorBlink sets cursor blinking.
func ExecuteSetCursorBlink(c parser.Command, v *VHS) error {
	var err error
	v.Options.CursorBlink, err = strconv.ParseBool(c.Args)
	if err != nil {
		return fmt.Errorf("failed to parse cursor blink: %w", err)
	}

	return nil
}

// ExecuteScreenshot is a CommandFunc that indicates a new screenshot must be taken.
func ExecuteScreenshot(c parser.Command, v *VHS) error {
	v.ScreenshotNextFrame(c.Args)
	return nil
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
