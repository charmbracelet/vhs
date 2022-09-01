package vhs

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

// Parse takes an input string provided by the lexer and parses it into a
// list of commands.
func (p *Parser) Parse() []Command {
	cmds := []Command{}

	for p.cur.Type != EOF {
		cmds = append(cmds, p.parseCommand())
		p.nextToken()
	}

	return cmds
}

// parseCommand parses a command.
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
	case CTRL:
		return p.parseCtrl()
	default:
		p.errors = append(p.errors, fmt.Sprintf("%2d:%-2d │ Invalid command: %s", p.cur.Line, p.cur.Column, p.cur.Literal))
		return Command{Type: Unknown}
	}
}

// parseSpeed parses a typing speed indication.
//
// i.e. @<time>
//
// This is optional (defaults to 100ms), thus skips (rather than error-ing)
// if the typing speed is not specified.
func (p *Parser) parseSpeed() string {
	if p.peek.Type == AT {
		p.nextToken()
		return p.parseTime()
	}

	return "100ms"
}

// parseRepeat parses an optional repeat count for a command.
//
// i.e. Backspace [count]
//
// This is optional (defaults to 1), thus skips (rather than error-ing)
// if the repeat count is not specified.
func (p *Parser) parseRepeat() string {
	if p.peek.Type == NUMBER {
		count := p.peek.Literal
		p.nextToken()
		return count
	}

	return "1"
}

// parseTime parses a time argument.
//
// i.e. <number>(s|ms)
//
func (p *Parser) parseTime() string {
	var t string

	if p.peek.Type == NUMBER {
		t = p.peek.Literal
		p.nextToken()
	} else {
		p.errors = append(p.errors, fmt.Sprintf("%2d:%-2d │ Expected time, got %s", p.cur.Line, p.cur.Column, p.peek.Literal))
	}

	if p.peek.Type == SECONDS || p.peek.Type == MILLISECONDS {
		t += p.peek.Literal
		p.nextToken()
	} else {
		t += "ms"
	}

	return t
}

// parseCtrl parses a control command.
// A control command takes a character to type while the modifier is held down.
//
// Ctrl+<character>
//
func (p *Parser) parseCtrl() Command {
	if p.peek.Type == PLUS {
		p.nextToken()
		if p.peek.Type == STRING {
			c := p.peek.Literal
			p.nextToken()
			return Command{Type: Ctrl, Args: c}
		}
	}

	p.errors = append(p.errors, fmt.Sprintf("%2d:%-2d │ Expected character, got %s", p.cur.Line, p.cur.Column, p.peek.Literal))
	return Command{Type: Ctrl}
}

// parseBackspace parses a backspace command.
// A backspace command takes an optional typing speed and optional count.
//
// Backspace[@<time>] [count]
//
func (p *Parser) parseBackspace() Command {
	cmd := Command{Type: Backspace}
	cmd.Options = p.parseSpeed()
	cmd.Args = p.parseRepeat()
	return cmd
}

// parseEnter parses an enter command.
// An enter command takes an optional typing speed and optional count.
//
// Enter[@<time>] [count]
//
func (p *Parser) parseEnter() Command {
	cmd := Command{Type: Enter}
	cmd.Options = p.parseSpeed()
	cmd.Args = p.parseRepeat()
	return cmd
}

// parseSet parses a set command.
// A set command takes a setting name and a value.
//
// Set <setting> <value>
//
func (p *Parser) parseSet() Command {
	cmd := Command{Type: Set}

	if p.peek.Type == SETTING {
		cmd.Options = p.peek.Literal
	} else {
		p.errors = append(p.errors, fmt.Sprintf("%2d:%-2d │ Unknown setting: %s", p.peek.Line, p.peek.Column, p.peek.Literal))
	}
	p.nextToken()

	cmd.Args = p.peek.Literal
	p.nextToken()

	// Allow Padding to have bare units (e.g. 10px, 5em, 10%)
	//
	// Set Padding 5em
	//
	if p.peek.Type == EM || p.peek.Type == PX || p.peek.Type == PERCENT {
		cmd.Args += p.peek.Literal
		p.nextToken()
	}

	return cmd
}

// parseSleep parses a sleep command.
// A sleep command takes a time for how long to sleep.
//
// Sleep <time>
//
func (p *Parser) parseSleep() Command {
	cmd := Command{Type: Sleep}
	cmd.Args = p.parseTime()
	return cmd
}

// parseType parses a type command.
// A type command takes a string to type.
//
// Type "string"
//
func (p *Parser) parseType() Command {
	cmd := Command{Type: Type}

	cmd.Options = p.parseSpeed()

	if p.peek.Type == STRING {
		cmd.Args = p.peek.Literal
		p.nextToken()
	} else {
		p.errors = append(p.errors, fmt.Sprintf("%2d:%-2d │ %s expects string", p.cur.Line, p.cur.Column, p.cur.Literal))
	}

	return cmd
}

// parseDown parses a down command.
// A down command takes an optional typing speed and optional count.
//
// Down[@<time>] [count]
//
func (p *Parser) parseDown() Command {
	cmd := Command{Type: Down}
	cmd.Options = p.parseSpeed()
	cmd.Args = p.parseRepeat()
	return cmd
}

// parseLeft parses a left command.
// A left command takes an optional typing speed and optional count.
//
// Left[@<time>] [count]
//
func (p *Parser) parseLeft() Command {
	cmd := Command{Type: Left}
	cmd.Options = p.parseSpeed()
	cmd.Args = p.parseRepeat()
	return cmd
}

// parseRight parses a right command.
// A right command takes an optional typing speed and optional count.
//
// Right[@<time>] [count]
//
func (p *Parser) parseRight() Command {
	cmd := Command{Type: Right}
	cmd.Options = p.parseSpeed()
	cmd.Args = p.parseRepeat()
	return cmd
}

// parseUp parses a up command.
// A up command takes an optional typing speed and optional count.
//
// Up[@<time>] [count]
//
func (p *Parser) parseUp() Command {
	cmd := Command{Type: Up}
	cmd.Options = p.parseSpeed()
	cmd.Args = p.parseRepeat()
	return cmd
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
