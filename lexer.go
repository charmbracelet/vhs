package dolly

type Lexer struct {
	ch      byte
	input   string
	pos     int
	nextPos int
	line    int
	column  int
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 1}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	l.ch = l.peekChar()
	l.pos = l.nextPos
	l.nextPos += 1
}

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
	case '+':
		tok = l.newToken(PLUS, l.ch)
		l.readChar()
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
		l.readChar()
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdentifier(tok.Literal)
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = NUMBER
		} else {
			tok = l.newToken(ILLEGAL, l.ch)
		}
	}
	return tok
}

func (l *Lexer) newToken(tokenType TokenType, ch byte) Token {
	return Token{
		Type:    tokenType,
		Literal: string(ch),
		Line:    l.line,
		Column:  l.column,
	}
}

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

func (l *Lexer) readNumber() string {
	pos := l.pos
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readIdentifier() string {
	pos := l.pos
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		if isNewLine(l.ch) {
			l.line += 1
			l.column = 1
		}
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isNewLine(ch byte) bool {
	return ch == '\n' || ch == '\r'
}

func (l *Lexer) peekChar() byte {
	if l.nextPos >= len(l.input) {
		return 0
	}
	return l.input[l.nextPos]
}
