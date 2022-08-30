package lexer

import (
	"testing"

	"github.com/charmbracelet/dolly/token"
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
Enter
Sleep 2s`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.SET, "Set"},
		{token.SETTING, "FontSize"},
		{token.NUMBER, "42"},
		{token.SET, "Set"},
		{token.SETTING, "Padding"},
		{token.NUMBER, "5"},
		{token.EM, "em"},
		{token.TYPE, "Type"},
		{token.STRING, "echo 'Hello, world!'"},
		{token.ENTER, "Enter"},
		{token.TYPE, "Type"},
		{token.AT, "@"},
		{token.NUMBER, "100"},
		{token.MILLISECONDS, "ms"},
		{token.STRING, "echo 'Hello, world!'"},
		{token.LEFT, "Left"},
		{token.NUMBER, "3"},
		{token.SLEEP, "Sleep"},
		{token.NUMBER, "1"},
		{token.SECONDS, "s"},
		{token.RIGHT, "Right"},
		{token.AT, "@"},
		{token.NUMBER, "100"},
		{token.MILLISECONDS, "ms"},
		{token.NUMBER, "3"},
		{token.SLEEP, "Sleep"},
		{token.NUMBER, "500"},
		{token.MILLISECONDS, "ms"},
		{token.ENTER, "Enter"},
		{token.SLEEP, "Sleep"},
		{token.NUMBER, "2"},
		{token.SECONDS, "s"},
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
