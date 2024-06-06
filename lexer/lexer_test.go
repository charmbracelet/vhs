package lexer

import (
	"os"
	"testing"

	"github.com/charmbracelet/vhs/token"
)

func TestNextToken(t *testing.T) {
	input := `
Output examples/out.gif
Set FontSize 42
Set Padding 5
Set CursorBlink false
Type "echo 'Hello, world!'"
Enter
Type@.1 "echo 'Hello, world!'"
Left 3
Sleep 1
Right@100ms 3
Sleep 500ms
Ctrl+C
Enter
Sleep .1
Sleep 100ms
Sleep 2
Wait+Screen@1m /foobar/`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.OUTPUT, "Output"},
		{token.STRING, "examples/out.gif"},
		{token.SET, "Set"},
		{token.FONT_SIZE, "FontSize"},
		{token.NUMBER, "42"},
		{token.SET, "Set"},
		{token.PADDING, "Padding"},
		{token.NUMBER, "5"},
		{token.SET, "Set"},
		{token.CURSOR_BLINK, "CursorBlink"},
		{token.BOOLEAN, "false"},
		{token.TYPE, "Type"},
		{token.STRING, "echo 'Hello, world!'"},
		{token.ENTER, "Enter"},
		{token.TYPE, "Type"},
		{token.AT, "@"},
		{token.NUMBER, ".1"},
		{token.STRING, "echo 'Hello, world!'"},
		{token.LEFT, "Left"},
		{token.NUMBER, "3"},
		{token.SLEEP, "Sleep"},
		{token.NUMBER, "1"},
		{token.RIGHT, "Right"},
		{token.AT, "@"},
		{token.NUMBER, "100"},
		{token.MILLISECONDS, "ms"},
		{token.NUMBER, "3"},
		{token.SLEEP, "Sleep"},
		{token.NUMBER, "500"},
		{token.MILLISECONDS, "ms"},
		{token.CTRL, "Ctrl"},
		{token.PLUS, "+"},
		{token.STRING, "C"},
		{token.ENTER, "Enter"},
		{token.SLEEP, "Sleep"},
		{token.NUMBER, ".1"},
		{token.SLEEP, "Sleep"},
		{token.NUMBER, "100"},
		{token.MILLISECONDS, "ms"},
		{token.SLEEP, "Sleep"},
		{token.NUMBER, "2"},
		{token.WAIT, "Wait"},
		{token.PLUS, "+"},
		{token.STRING, "Screen"},
		{token.AT, "@"},
		{token.NUMBER, "1"},
		{token.MINUTES, "m"},
		{token.REGEX, "foobar"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLexTapeFile(t *testing.T) {
	input, err := os.ReadFile("../examples/fixtures/all.tape")
	if err != nil {
		t.Fatal("could not read all.tape file")
	}

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.COMMENT, " All Commands"},
		{token.COMMENT, " Output:"},
		{token.OUTPUT, "Output"},
		{token.STRING, "examples/fixtures/all.gif"},
		{token.OUTPUT, "Output"},
		{token.STRING, "examples/fixtures/all.mp4"},
		{token.OUTPUT, "Output"},
		{token.STRING, "examples/fixtures/all.webm"},
		{token.COMMENT, " Settings:"},
		{token.SET, "Set"},
		{token.SHELL, "Shell"},
		{token.STRING, "fish"},
		{token.SET, "Set"},
		{token.FONT_SIZE, "FontSize"},
		{token.NUMBER, "22"},
		{token.SET, "Set"},
		{token.FONT_FAMILY, "FontFamily"},
		{token.STRING, "DejaVu Sans Mono"},
		{token.SET, "Set"},
		{token.HEIGHT, "Height"},
		{token.NUMBER, "600"},
		{token.SET, "Set"},
		{token.WIDTH, "Width"},
		{token.NUMBER, "1200"},
		{token.SET, "Set"},
		{token.LETTER_SPACING, "LetterSpacing"},
		{token.NUMBER, "1"},
		{token.SET, "Set"},
		{token.LINE_HEIGHT, "LineHeight"},
		{token.NUMBER, "1.2"},
		{token.SET, "Set"},
		{token.THEME, "Theme"},
		{token.JSON, "{ \"name\": \"Whimsy\", \"black\": \"#535178\", \"red\": \"#ef6487\", \"green\": \"#5eca89\", \"yellow\": \"#fdd877\", \"blue\": \"#65aef7\", \"purple\": \"#aa7ff0\", \"cyan\": \"#43c1be\", \"white\": \"#ffffff\", \"brightBlack\": \"#535178\", \"brightRed\": \"#ef6487\", \"brightGreen\": \"#5eca89\", \"brightYellow\": \"#fdd877\", \"brightBlue\": \"#65aef7\", \"brightPurple\": \"#aa7ff0\", \"brightCyan\": \"#43c1be\", \"brightWhite\": \"#ffffff\", \"background\": \"#29283b\", \"foreground\": \"#b3b0d6\", \"selectionBackground\": \"#3d3c58\", \"cursorColor\": \"#b3b0d6\" }"},
		{token.SET, "Set"},
		{token.THEME, "Theme"},
		{token.STRING, "Catppuccin Mocha"},
		{token.SET, "Set"},
		{token.PADDING, "Padding"},
		{token.NUMBER, "50"},
		{token.SET, "Set"},
		{token.FRAMERATE, "Framerate"},
		{token.NUMBER, "60"},
		{token.SET, "Set"},
		{token.PLAYBACK_SPEED, "PlaybackSpeed"},
		{token.NUMBER, "2"},
		{token.SET, "Set"},
		{token.TYPING_SPEED, "TypingSpeed"},
		{token.NUMBER, ".1"},
		{token.SET, "Set"},
		{token.LOOP_OFFSET, "LoopOffset"},
		{token.NUMBER, "60.4"},
		{token.SET, "Set"},
		{token.LOOP_OFFSET, "LoopOffset"},
		{token.NUMBER, "20.99"},
		{token.PERCENT, "%"},
		{token.SET, "Set"},
		{token.CURSOR_BLINK, "CursorBlink"},
		{token.BOOLEAN, "false"},
		{token.COMMENT, " Sleep:"},
		{token.SLEEP, "Sleep"},
		{token.NUMBER, "1"},
		{token.SLEEP, "Sleep"},
		{token.NUMBER, "500"},
		{token.MILLISECONDS, "ms"},
		{token.SLEEP, "Sleep"},
		{token.NUMBER, ".5"},
		{token.SLEEP, "Sleep"},
		{token.NUMBER, "0.5"},
		{token.COMMENT, " Type:"},
		{token.TYPE, "Type"},
		{token.AT, "@"},
		{token.NUMBER, ".5"},
		{token.STRING, "All"},
		{token.TYPE, "Type"},
		{token.AT, "@"},
		{token.NUMBER, "500"},
		{token.MILLISECONDS, "ms"},
		{token.STRING, "All"},
		{token.TYPE, "Type"},
		{token.STRING, "Double Quote"},
		{token.TYPE, "Type"},
		{token.STRING, "\"Single\" Quote"},
		{token.TYPE, "Type"},
		{token.STRING, `"Backtick" 'Quote'`},
		{token.COMMENT, " Keys:"},
		{token.BACKSPACE, "Backspace"},
		{token.BACKSPACE, "Backspace"},
		{token.NUMBER, "2"},
		{token.BACKSPACE, "Backspace"},
		{token.AT, "@"},
		{token.NUMBER, "1"},
		{token.NUMBER, "3"},
		{token.DELETE, "Delete"},
		{token.DELETE, "Delete"},
		{token.NUMBER, "2"},
		{token.DELETE, "Delete"},
		{token.AT, "@"},
		{token.NUMBER, "1"},
		{token.NUMBER, "3"},
		{token.INSERT, "Insert"},
		{token.INSERT, "Insert"},
		{token.NUMBER, "2"},
		{token.INSERT, "Insert"},
		{token.AT, "@"},
		{token.NUMBER, "1"},
		{token.NUMBER, "3"},
		{token.DOWN, "Down"},
		{token.DOWN, "Down"},
		{token.NUMBER, "2"},
		{token.DOWN, "Down"},
		{token.AT, "@"},
		{token.NUMBER, "1"},
		{token.NUMBER, "3"},
		{token.PAGEDOWN, "PageDown"},
		{token.PAGEDOWN, "PageDown"},
		{token.NUMBER, "2"},
		{token.PAGEDOWN, "PageDown"},
		{token.AT, "@"},
		{token.NUMBER, "1"},
		{token.NUMBER, "3"},
		{token.ENTER, "Enter"},
		{token.ENTER, "Enter"},
		{token.NUMBER, "2"},
		{token.ENTER, "Enter"},
		{token.AT, "@"},
		{token.NUMBER, "1"},
		{token.NUMBER, "3"},
		{token.SPACE, "Space"},
		{token.SPACE, "Space"},
		{token.NUMBER, "2"},
		{token.SPACE, "Space"},
		{token.AT, "@"},
		{token.NUMBER, "1"},
		{token.NUMBER, "3"},
		{token.TAB, "Tab"},
		{token.TAB, "Tab"},
		{token.NUMBER, "2"},
		{token.TAB, "Tab"},
		{token.AT, "@"},
		{token.NUMBER, "1"},
		{token.NUMBER, "3"},
		{token.LEFT, "Left"},
		{token.LEFT, "Left"},
		{token.NUMBER, "2"},
		{token.LEFT, "Left"},
		{token.AT, "@"},
		{token.NUMBER, "1"},
		{token.NUMBER, "3"},
		{token.RIGHT, "Right"},
		{token.RIGHT, "Right"},
		{token.NUMBER, "2"},
		{token.RIGHT, "Right"},
		{token.AT, "@"},
		{token.NUMBER, "1"},
		{token.NUMBER, "3"},
		{token.UP, "Up"},
		{token.UP, "Up"},
		{token.NUMBER, "2"},
		{token.UP, "Up"},
		{token.AT, "@"},
		{token.NUMBER, "1"},
		{token.NUMBER, "3"},
		{token.PAGEUP, "PageUp"},
		{token.PAGEUP, "PageUp"},
		{token.NUMBER, "2"},
		{token.PAGEUP, "PageUp"},
		{token.AT, "@"},
		{token.NUMBER, "1"},
		{token.NUMBER, "3"},
		{token.DOWN, "Down"},
		{token.DOWN, "Down"},
		{token.NUMBER, "2"},
		{token.DOWN, "Down"},
		{token.AT, "@"},
		{token.NUMBER, "1"},
		{token.NUMBER, "3"},
		{token.COMMENT, " Control:"},
		{token.CTRL, "Ctrl"},
		{token.PLUS, "+"},
		{token.STRING, "C"},
		{token.CTRL, "Ctrl"},
		{token.PLUS, "+"},
		{token.STRING, "L"},
		{token.CTRL, "Ctrl"},
		{token.PLUS, "+"},
		{token.STRING, "R"},
		{token.COMMENT, " Alt:"},
		{token.ALT, "Alt"},
		{token.PLUS, "+"},
		{token.STRING, "."},
		{token.ALT, "Alt"},
		{token.PLUS, "+"},
		{token.STRING, "L"},
		{token.ALT, "Alt"},
		{token.PLUS, "+"},
		{token.STRING, "i"},
		{token.COMMENT, " Display:"},
		{token.HIDE, "Hide"},
		{token.SHOW, "Show"},
	}

	l := New(string(input))

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
