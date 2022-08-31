package dolly

import "fmt"

// Parser is the structure that manages the parsing of tokens.
type Parser struct {
	l      *Lexer
	errors []string
	cur    Token
	peek   Token
}

// NewParser returns a new Parser.
func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	// Read two tokens, so cur and peek are both set.
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Parse() []Command {
	cmds := []Command{}

	for p.cur.Type != EOF {
		cmds = append(cmds, p.parseCommand())
		p.nextToken()
	}

	return cmds
}

func (p *Parser) parseCommand() Command {
	switch p.cur.Type {
	case BACKSPACE:
		return p.parseBackspace()
	case ENTER:
		return p.parseEnter()
	case SET:
		return p.parseSet()
	case SLEEP:
		return p.parseSleep()
	case TYPE:
		return p.parseType()
	case DOWN:
		return p.parseDown()
	case LEFT:
		return p.parseLeft()
	case RIGHT:
		return p.parseRight()
	case UP:
		return p.parseUp()
	default:
		p.errors = append(p.errors, fmt.Sprintf("unknown command: %s", p.cur.Literal))
		return Command{}
	}
}

// parseBackspace parses a backspace command.
// A backspace command takes an optional typing speed and optional count.
//
// Backspace[@<time>] [count]
//
func (p *Parser) parseBackspace() Command {
	return Command{
		Type: BACKSPACE,
	}
}

// parseEnter parses an enter command.
// An enter command takes an optional typing speed and optional count.
//
// Enter[@<time>] [count]
//
func (p *Parser) parseEnter() Command {
	return Command{
		Type: ENTER,
	}
}

// parseSet parses a set command.
// A set command takes a setting name and a value.
//
// Set <setting> <value>
//
func (p *Parser) parseSet() Command {
	return Command{
		Type: SET,
	}
}

// parseSleep parses a sleep command.
// A sleep command takes a time for how long to sleep.
//
// Sleep <time>
//
func (p *Parser) parseSleep() Command {
	return Command{
		Type: SLEEP,
	}
}

// parseType parses a type command.
// A type command takes a string to type.
//
// Type "string"
//
func (p *Parser) parseType() Command {
	return Command{
		Type: TYPE,
	}
}

// parseDown parses a down command.
// A down command takes an optional typing speed and optional count.
//
// Down[@<time>] [count]
//
func (p *Parser) parseDown() Command {
	return Command{
		Type: DOWN,
	}
}

// parseLeft parses a left command.
// A left command takes an optional typing speed and optional count.
//
// Left[@<time>] [count]
//
func (p *Parser) parseLeft() Command {
	return Command{
		Type: LEFT,
	}
}

// parseRight parses a right command.
// A right command takes an optional typing speed and optional count.
//
// Right[@<time>] [count]
//
func (p *Parser) parseRight() Command {
	return Command{
		Type: RIGHT,
	}
}

// parseUp parses a up command.
// A up command takes an optional typing speed and optional count.
//
// Up[@<time>] [count]
//
func (p *Parser) parseUp() Command {
	return Command{
		Type: UP,
	}
}

// Errors returns any errors that occurred during parsing.
func (p *Parser) Errors() []string {
	return p.errors
}

// nextToken gets the next token from the lexer
// and updates the parser tokens accordingly.
func (p *Parser) nextToken() {
	p.cur = p.peek
	p.peek = p.l.NextToken()
}
