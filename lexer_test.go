package vhs

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
Set FontSize 42
Set Padding 5em
Type "echo 'Hello, world!'"
Enter
Type@100ms "echo 'Hello, world!'"
Left 3
Sleep 1s
Right@100ms 3
Sleep 500ms
Ctrl+C
Enter
Sleep 2s`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{SET, "Set"},
		{SETTING, "FontSize"},
		{NUMBER, "42"},
		{SET, "Set"},
		{SETTING, "Padding"},
		{NUMBER, "5"},
		{EM, "em"},
		{TYPE, "Type"},
		{STRING, "echo 'Hello, world!'"},
		{ENTER, "Enter"},
		{TYPE, "Type"},
		{AT, "@"},
		{NUMBER, "100"},
		{MILLISECONDS, "ms"},
		{STRING, "echo 'Hello, world!'"},
		{LEFT, "Left"},
		{NUMBER, "3"},
		{SLEEP, "Sleep"},
		{NUMBER, "1"},
		{SECONDS, "s"},
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
		{NUMBER, "2"},
		{SECONDS, "s"},
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
