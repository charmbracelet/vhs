package main

import (
	"os"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
Output examples/out.gif
Set FontSize 42
Set Padding 5
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
Sleep 2`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{OUTPUT, "Output"},
		{STRING, "examples/out.gif"},
		{SET, "Set"},
		{FONT_SIZE, "FontSize"},
		{NUMBER, "42"},
		{SET, "Set"},
		{PADDING, "Padding"},
		{NUMBER, "5"},
		{TYPE, "Type"},
		{STRING, "echo 'Hello, world!'"},
		{ENTER, "Enter"},
		{TYPE, "Type"},
		{AT, "@"},
		{NUMBER, ".1"},
		{STRING, "echo 'Hello, world!'"},
		{LEFT, "Left"},
		{NUMBER, "3"},
		{SLEEP, "Sleep"},
		{NUMBER, "1"},
		{RIGHT, "Right"},
		{AT, "@"},
		{NUMBER, "100"},
		{MILLISECONDS, "ms"},
		{NUMBER, "3"},
		{SLEEP, "Sleep"},
		{NUMBER, "500"},
		{MILLISECONDS, "ms"},
		{CTRL, "Ctrl"},
		{PLUS, "+"},
		{STRING, "C"},
		{ENTER, "Enter"},
		{SLEEP, "Sleep"},
		{NUMBER, ".1"},
		{SLEEP, "Sleep"},
		{NUMBER, "100"},
		{MILLISECONDS, "ms"},
		{SLEEP, "Sleep"},
		{NUMBER, "2"},
	}

	l := NewLexer(input)

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
	input, err := os.ReadFile("examples/fixtures/all.tape")
	if err != nil {
		t.Fatal("could not read all.tape file")
	}

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{COMMENT, " All Commands"},
		{COMMENT, " Output:"},
		{OUTPUT, "Output"},
		{STRING, "examples/fixtures/all.gif"},
		{OUTPUT, "Output"},
		{STRING, "examples/fixtures/all.mp4"},
		{OUTPUT, "Output"},
		{STRING, "examples/fixtures/all.webm"},
		{COMMENT, " Settings:"},
		{SET, "Set"},
		{FONT_SIZE, "FontSize"},
		{NUMBER, "22"},
		{SET, "Set"},
		{FONT_FAMILY, "FontFamily"},
		{STRING, "DejaVu Sans Mono"},
		{SET, "Set"},
		{HEIGHT, "Height"},
		{NUMBER, "600"},
		{SET, "Set"},
		{WIDTH, "Width"},
		{NUMBER, "1200"},
		{SET, "Set"},
		{LETTER_SPACING, "LetterSpacing"},
		{NUMBER, "1"},
		{SET, "Set"},
		{LINE_HEIGHT, "LineHeight"},
		{NUMBER, "1.2"},
		{SET, "Set"},
		{THEME, "Theme"},
		{JSON, "{ \"name\": \"Whimsy\", \"black\": \"#535178\", \"red\": \"#ef6487\", \"green\": \"#5eca89\", \"yellow\": \"#fdd877\", \"blue\": \"#65aef7\", \"purple\": \"#aa7ff0\", \"cyan\": \"#43c1be\", \"white\": \"#ffffff\", \"brightBlack\": \"#535178\", \"brightRed\": \"#ef6487\", \"brightGreen\": \"#5eca89\", \"brightYellow\": \"#fdd877\", \"brightBlue\": \"#65aef7\", \"brightPurple\": \"#aa7ff0\", \"brightCyan\": \"#43c1be\", \"brightWhite\": \"#ffffff\", \"background\": \"#29283b\", \"foreground\": \"#b3b0d6\", \"selectionBackground\": \"#3d3c58\", \"cursorColor\": \"#b3b0d6\" }"},
		{SET, "Set"},
		{PADDING, "Padding"},
		{NUMBER, "50"},
		{SET, "Set"},
		{FRAMERATE, "Framerate"},
		{NUMBER, "60"},
		{SET, "Set"},
		{PLAYBACK_SPEED, "PlaybackSpeed"},
		{NUMBER, "2"},
		{SET, "Set"},
		{TYPING_SPEED, "TypingSpeed"},
		{NUMBER, ".1"},
		{SET, "Set"},
		{LOOP_OFFSET, "LoopOffset"},
		{NUMBER, "60"},
		{SET, "Set"},
		{LOOP_OFFSET, "LoopOffset"},
		{NUMBER, "20"},
		{PERCENT, "%"},
		{COMMENT, " Sleep:"},
		{SLEEP, "Sleep"},
		{NUMBER, "1"},
		{SLEEP, "Sleep"},
		{NUMBER, "500"},
		{MILLISECONDS, "ms"},
		{SLEEP, "Sleep"},
		{NUMBER, ".5"},
		{SLEEP, "Sleep"},
		{NUMBER, "0.5"},
		{COMMENT, " Type:"},
		{TYPE, "Type"},
		{STRING, "All"},
		{TYPE, "Type"},
		{AT, "@"},
		{NUMBER, ".5"},
		{STRING, "All"},
		{TYPE, "Type"},
		{AT, "@"},
		{NUMBER, "500"},
		{MILLISECONDS, "ms"},
		{STRING, "All"},
		{COMMENT, " Keys:"},
		{BACKSPACE, "Backspace"},
		{BACKSPACE, "Backspace"},
		{NUMBER, "2"},
		{BACKSPACE, "Backspace"},
		{AT, "@"},
		{NUMBER, "1"},
		{NUMBER, "3"},
		{DOWN, "Down"},
		{DOWN, "Down"},
		{NUMBER, "2"},
		{DOWN, "Down"},
		{AT, "@"},
		{NUMBER, "1"},
		{NUMBER, "3"},
		{ENTER, "Enter"},
		{ENTER, "Enter"},
		{NUMBER, "2"},
		{ENTER, "Enter"},
		{AT, "@"},
		{NUMBER, "1"},
		{NUMBER, "3"},
		{SPACE, "Space"},
		{SPACE, "Space"},
		{NUMBER, "2"},
		{SPACE, "Space"},
		{AT, "@"},
		{NUMBER, "1"},
		{NUMBER, "3"},
		{TAB, "Tab"},
		{TAB, "Tab"},
		{NUMBER, "2"},
		{TAB, "Tab"},
		{AT, "@"},
		{NUMBER, "1"},
		{NUMBER, "3"},
		{LEFT, "Left"},
		{LEFT, "Left"},
		{NUMBER, "2"},
		{LEFT, "Left"},
		{AT, "@"},
		{NUMBER, "1"},
		{NUMBER, "3"},
		{RIGHT, "Right"},
		{RIGHT, "Right"},
		{NUMBER, "2"},
		{RIGHT, "Right"},
		{AT, "@"},
		{NUMBER, "1"},
		{NUMBER, "3"},
		{UP, "Up"},
		{UP, "Up"},
		{NUMBER, "2"},
		{UP, "Up"},
		{AT, "@"},
		{NUMBER, "1"},
		{NUMBER, "3"},
		{DOWN, "Down"},
		{DOWN, "Down"},
		{NUMBER, "2"},
		{DOWN, "Down"},
		{AT, "@"},
		{NUMBER, "1"},
		{NUMBER, "3"},
		{COMMENT, " Control:"},
		{CTRL, "Ctrl"},
		{PLUS, "+"},
		{STRING, "C"},
		{CTRL, "Ctrl"},
		{PLUS, "+"},
		{STRING, "L"},
		{CTRL, "Ctrl"},
		{PLUS, "+"},
		{STRING, "R"},
		{COMMENT, " Display:"},
		{HIDE, "Hide"},
		{SHOW, "Show"},
	}

	l := NewLexer(string(input))

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
