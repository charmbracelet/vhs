package lexer

import "github.com/charmbracelet/vhs/token"

// Lexer is a lexer that tokenizes the input.
type Lexer struct {
	ch      byte
	input   string
	pos     int
	nextPos int
	line    int
	column  int
}

// New returns a new lexer for tokenizing the input string.
func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.readChar()
	return l
}

// readChar advances the lexer to the next character.
func (l *Lexer) readChar() {
	l.column++
	l.ch = l.peekChar()
	l.pos = l.nextPos
	l.nextPos++
}

// NextToken returns the next token in the input.
func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	tok := token.Token{Line: l.line, Column: l.column}

	switch l.ch {
	case 0:
		tok = l.newToken(token.EOF, l.ch)
	case '@':
		tok = l.newToken(token.AT, l.ch)
		l.readChar()
	case '=':
		tok = l.newToken(token.EQUAL, l.ch)
		l.readChar()
	case '%':
		tok = l.newToken(token.PERCENT, l.ch)
		l.readChar()
	case '#':
		tok.Type = token.COMMENT
		tok.Literal = l.readComment()
	case '+':
		tok = l.newToken(token.PLUS, l.ch)
		l.readChar()
	case '{':
		tok.Type = token.JSON
		tok.Literal = "{" + l.readJSON() + "}"
		l.readChar()
	case '`':
		tok.Type = token.STRING
		tok.Literal = l.readString('`')
		l.readChar()
	case '\'':
		tok.Type = token.STRING
		tok.Literal = l.readString('\'')
		l.readChar()
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString('"')
		l.readChar()
	default:
		if isDigit(l.ch) || (isDot(l.ch) && isDigit(l.peekChar())) {
			tok.Literal = l.readNumber()
			tok.Type = token.NUMBER
		} else if isLetter(l.ch) || isDot(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Literal)
		} else {
			tok = l.newToken(token.ILLEGAL, l.ch)
			l.readChar()
		}
	}
	return tok
}

// newToken creates a new token with the given type and literal.
func (l *Lexer) newToken(tokenType token.Type, ch byte) token.Token {
	literal := string(ch)
	return token.Token{
		Type:    tokenType,
		Literal: literal,
		Line:    l.line,
		Column:  l.column,
	}
}

// readComment reads a comment.
// // Foo => Token(Foo).
func (l *Lexer) readComment() string {
	pos := l.pos + 1
	for {
		l.readChar()
		if isNewLine(l.ch) || l.ch == 0 {
			break
		}
	}
	// The current character is a newline.
	// skipWhitespace() will handle this for us and increment the line counter.
	return l.input[pos:l.pos]
}

// readString reads a string from the input.
// "Foo" => Token(Foo).
func (l *Lexer) readString(endChar byte) string {
	pos := l.pos + 1
	for {
		l.readChar()
		if l.ch == endChar || l.ch == 0 || isNewLine(l.ch) {
			break
		}
	}
	return l.input[pos:l.pos]
}

// readJSON reads a JSON object from the input.
// {"foo": "bar"} => Token({"foo": "bar"}).
func (l *Lexer) readJSON() string {
	pos := l.pos + 1
	for {
		l.readChar()
		if l.ch == '}' || l.ch == 0 {
			break
		}
	}
	return l.input[pos:l.pos]
}

// readNumber reads a number from the input.
// 123 => Token(123).
func (l *Lexer) readNumber() string {
	pos := l.pos
	for isDigit(l.ch) || isDot(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

// readIdentifier reads an identifier from the input.
// Foo => Token(Foo).
func (l *Lexer) readIdentifier() string {
	pos := l.pos
	for isLetter(l.ch) || isDot(l.ch) || isDash(l.ch) || isUnderscore(l.ch) || isSlash(l.ch) || isPercent(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

// skipWhitespace skips whitespace characters.
// If it encounters a newline, it increments the line counter to keep track
// of the token's line number.
func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		// Note: we don't use isNewline since we don't want to double count \r\n on
		// windows and increment the l.line.
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		l.readChar()
	}
}

// isDot returns whether a character is a dot.
func isDot(ch byte) bool {
	return ch == '.'
}

// isDash returns whether a character is a dash.
func isDash(ch byte) bool {
	return ch == '-'
}

// isUnderscore returns whether a character is an underscore.
func isUnderscore(ch byte) bool {
	return ch == '_'
}

// isPercent returns whether a character is a percent.
func isPercent(ch byte) bool {
	return ch == '%'
}

// isSlash returns whether a character is a slash.
func isSlash(ch byte) bool {
	return ch == '/'
}

// isLetter returns whether a character is a letter.
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

// isDigit returns whether a character is a digit.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// isWhitespace returns whether a character is a whitespace.
func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

// isNewLine returns whether a character is a newline.
//
// Note: in windows a single new line is \r\n so using isNewline is not
// recommended for counting the number of new lines.
func isNewLine(ch byte) bool {
	return ch == '\n' || ch == '\r'
}

// peekChar returns the next character in the input without advancing the lexer.
func (l *Lexer) peekChar() byte {
	if l.nextPos >= len(l.input) {
		return 0
	}
	return l.input[l.nextPos]
}
