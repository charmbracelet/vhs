package vhs

// Lexer is a lexer that tokenizes the input.
type Lexer struct {
	ch      byte
	input   string
	pos     int
	nextPos int
	line    int
	column  int
}

// NewLexer returns a new lexer for tokenizing the input string.
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.readChar()
	return l
}

// readChar advances the lexer to the next character.
func (l *Lexer) readChar() {
	l.column += 1
	l.ch = l.peekChar()
	l.pos = l.nextPos
	l.nextPos += 1
}

// NextToken returns the next token in the input.
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	var tok = Token{Line: l.line, Column: l.column}

	switch l.ch {
	case 0:
		tok = l.newToken(EOF, l.ch)
	case '@':
		tok = l.newToken(AT, l.ch)
		l.readChar()
	case '=':
		tok = l.newToken(EQUAL, l.ch)
		l.readChar()
	case '%':
		tok = l.newToken(PERCENT, l.ch)
		l.readChar()
	case '+':
		tok = l.newToken(PLUS, l.ch)
		l.readChar()
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
		l.readChar()
	default:
		if isLetter(l.ch) || isDot(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdentifier(tok.Literal)
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = NUMBER
		} else {
			tok = l.newToken(ILLEGAL, l.ch)
			l.readChar()
		}
	}
	return tok
}

// newToken creates a new token with the given type and literal.
func (l *Lexer) newToken(tokenType TokenType, ch byte) Token {
	literal := string(ch)
	return Token{
		Type:    tokenType,
		Literal: literal,
		Line:    l.line,
		Column:  l.column,
	}
}

// readString reads a string from the input.
// "Foo" => Token(Foo)
func (l *Lexer) readString() string {
	pos := l.pos + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[pos:l.pos]
}

// readNumber reads a number from the input.
// 123 => Token(123)
func (l *Lexer) readNumber() string {
	pos := l.pos
	for isDigit(l.ch) || isDot(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

// readIdentifier reads an identifier from the input.
// Foo => Token(Foo)
func (l *Lexer) readIdentifier() string {
	pos := l.pos
	for isLetter(l.ch) || isDot(l.ch) || isDash(l.ch) || isSlash(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

// skipWhitespace skips whitespace characters.
// If it encounters a newline, it increments the line counter to keep track
// of the token's line number.
func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		if isNewLine(l.ch) {
			l.line += 1
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
