package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Parser is the structure that manages the parsing of tokens.
type Parser struct {
	l      *Lexer
	errors []ParserError
	cur    Token
	peek   Token
}

// NewParser returns a new Parser.
func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l, errors: []ParserError{}}

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
		if p.cur.Type == COMMENT {
			p.nextToken()
			continue
		}
		cmds = append(cmds, p.parseCommand())
		p.nextToken()
	}

	return cmds
}

// parseCommand parses a command.
func (p *Parser) parseCommand() Command {
	switch p.cur.Type {
	case SPACE, BACKSPACE, ENTER, ESCAPE, TAB, DOWN, LEFT, RIGHT, UP, PAGEUP, PAGEDOWN:
		return p.parseKeypress(p.cur.Type)
	case SET:
		return p.parseSet()
	case OUTPUT:
		return p.parseOutput()
	case SLEEP:
		return p.parseSleep()
	case TYPE:
		return p.parseType()
	case CTRL:
		return p.parseCtrl()
	case ALT:
		return p.parseAlt()
	case SHIFT:
		return p.parseShift()
	case HIDE:
		return p.parseHide()
	case REQUIRE:
		return p.parseRequire()
	case SHOW:
		return p.parseShow()
	case SOURCE:
		return p.parseSource()
	case SCREENSHOT:
		return p.parseScreenshot()
	case COPY:
		return p.parseCopy()
	case PASTE:
		return p.parsePaste()
	default:
		p.errors = append(p.errors, NewError(p.cur, "Invalid command: "+p.cur.Literal))
		return Command{Type: ILLEGAL}
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
	return ""
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
// <number>[ms]
func (p *Parser) parseTime() string {
	var t string

	if p.peek.Type == NUMBER {
		t = p.peek.Literal
		p.nextToken()
	} else {
		p.errors = append(p.errors, NewError(p.cur, "Expected time after "+p.cur.Literal))
	}

	// Allow TypingSpeed to have bare units (e.g. 50ms, 100ms)
	if p.peek.Type == MILLISECONDS || p.peek.Type == SECONDS {
		t += p.peek.Literal
		p.nextToken()
	} else {
		t += "s"
	}

	return t
}

// parseCtrl parses a control command.
// A control command takes one or multiples characters and/or modifiers to type while ctrl is held down.
//
// Ctrl[+Alt][+Shift]+<char>
// E.g:
// Ctrl+Shift+O
// Ctrl+Alt+Shift+P
func (p *Parser) parseCtrl() Command {
	var args []string

	for p.peek.Type == PLUS {
		p.nextToken()
		peek := p.peek

		// Get key from keywords and check if it's a valid modifier
		if k, ok := keywords[peek.Literal]; ok {
			p.nextToken()
			if IsModifier(k) {
				args = append(args, peek.Literal)
			} else {
				p.errors = append(p.errors, NewError(p.cur, "not a valid modifier"))
			}
		}

		// Add key argument
		if peek.Type == STRING && len(peek.Literal) == 1 {
			p.nextToken()
			args = append(args, peek.Literal)
		}

		// Key arguments with len > 1 are not valid
		if peek.Type == STRING && len(peek.Literal) > 1 {
			p.nextToken()
			p.errors = append(p.errors, NewError(p.cur, "Invalid control argument: "+p.cur.Literal))
		}
	}

	if len(args) == 0 {
		p.errors = append(p.errors, NewError(p.cur, "Expected control character with args, got "+p.cur.Literal))
	}

	ctrlArgs := strings.Join(args, " ")
	return Command{Type: CTRL, Args: ctrlArgs}
}

// parseAlt parses an alt command.
// An alt command takes a character to type while the modifier is held down.
//
// Alt+<character>
func (p *Parser) parseAlt() Command {
	if p.peek.Type == PLUS {
		p.nextToken()
		if p.peek.Type == STRING || p.peek.Type == ENTER || p.peek.Type == TAB {
			c := p.peek.Literal
			p.nextToken()
			return Command{Type: ALT, Args: c}
		}
	}

	p.errors = append(p.errors, NewError(p.cur, "Expected alt character, got "+p.cur.Literal))
	return Command{Type: ALT}
}

// parseShift parses a shift command.
// A shift command takes one character and types while shift is held down.
//
// Shift+<char>
// E.g.
// Shift+A
// Shift+Tab
// Shift+Enter
func (p *Parser) parseShift() Command {
	if p.peek.Type == PLUS {
		p.nextToken()
		if p.peek.Type == STRING || p.peek.Type == ENTER || p.peek.Type == TAB {
			c := p.peek.Literal
			p.nextToken()
			return Command{Type: SHIFT, Args: c}
		}
	}

	p.errors = append(p.errors, NewError(p.cur, "Expected shift character, got "+p.cur.Literal))
	return Command{Type: SHIFT}
}

// parseKeypress parses a repeatable and time adjustable keypress command.
// A keypress command takes an optional typing speed and optional count.
//
// Key[@<time>] [count]
func (p *Parser) parseKeypress(ct TokenType) Command {
	cmd := Command{Type: CommandType(ct)}
	cmd.Options = p.parseSpeed()
	cmd.Args = p.parseRepeat()
	return cmd
}

// parseOutput parses an output command.
// An output command takes a file path to which to output.
//
// Output <path>
func (p *Parser) parseOutput() Command {
	cmd := Command{Type: OUTPUT}

	if p.peek.Type != STRING {
		p.errors = append(p.errors, NewError(p.cur, "Expected file path after output"))
		return cmd
	}

	ext := filepath.Ext(p.peek.Literal)
	if ext != "" {
		cmd.Options = ext
	} else {
		cmd.Options = ".png"
		if !strings.HasSuffix(p.peek.Literal, "/") {
			p.errors = append(p.errors, NewError(p.peek, "Expected folder with trailing slash"))
		}
	}

	cmd.Args = p.peek.Literal
	p.nextToken()
	return cmd
}

// parseSet parses a set command.
// A set command takes a setting name and a value.
//
// Set <setting> <value>
func (p *Parser) parseSet() Command {
	cmd := Command{Type: SET}

	if IsSetting(p.peek.Type) {
		cmd.Options = p.peek.Literal
	} else {
		p.errors = append(p.errors, NewError(p.peek, "Unknown setting: "+p.peek.Literal))
	}
	p.nextToken()

	switch p.cur.Type {
	case LOOP_OFFSET:
		cmd.Args = p.peek.Literal
		p.nextToken()
		// Allow LoopOffset without '%'
		// Set LoopOffset 20
		cmd.Args += "%"
		if p.peek.Type == PERCENT {
			p.nextToken()
		}
	case TYPING_SPEED:
		cmd.Args = p.peek.Literal
		p.nextToken()
		// Allow TypingSpeed to have bare units (e.g. 10ms)
		// Set TypingSpeed 10ms
		if p.peek.Type == MILLISECONDS || p.peek.Type == SECONDS {
			cmd.Args += p.peek.Literal
			p.nextToken()
		} else if cmd.Options == "TypingSpeed" {
			cmd.Args += "s"
		}
	case WINDOW_BAR:
		cmd.Args = p.peek.Literal
		p.nextToken()

		windowBar := p.cur.Literal
		if !isValidWindowBar(windowBar) {
			p.errors = append(
				p.errors,
				NewError(p.cur, windowBar+" is not a valid bar style."),
			)
		}
	case MARGIN_FILL:
		cmd.Args = p.peek.Literal
		p.nextToken()

		marginFill := p.cur.Literal

		// Check if margin color is a valid hex string
		if strings.HasPrefix(marginFill, "#") {
			_, err := strconv.ParseUint(marginFill[1:], 16, 64)

			if err != nil || len(marginFill) != 7 {
				p.errors = append(
					p.errors,
					NewError(
						p.cur,
						"\""+marginFill+"\" is not a valid color.",
					),
				)
			}
		}
	case CURSOR_BLINK:
		cmd.Args = p.peek.Literal
		p.nextToken()

		if p.cur.Type != BOOLEAN {
			p.errors = append(
				p.errors,
				NewError(p.cur, "expected boolean value."),
			)
		}

	default:
		cmd.Args = p.peek.Literal
		p.nextToken()
	}

	return cmd
}

// parseSleep parses a sleep command.
// A sleep command takes a time for how long to sleep.
//
// Sleep <time>
func (p *Parser) parseSleep() Command {
	cmd := Command{Type: SLEEP}
	cmd.Args = p.parseTime()
	return cmd
}

// parseHide parses a Hide command.
//
// Hide
// ...
func (p *Parser) parseHide() Command {
	cmd := Command{Type: HIDE}
	return cmd
}

// parseRequire parses a Require command.
//
// ...
// Require
func (p *Parser) parseRequire() Command {
	cmd := Command{Type: REQUIRE}

	if p.peek.Type != STRING {
		p.errors = append(p.errors, NewError(p.peek, p.cur.Literal+" expects one string"))
	}

	cmd.Args = p.peek.Literal
	p.nextToken()

	return cmd
}

// parseShow parses a Show command.
//
// ...
// Show
func (p *Parser) parseShow() Command {
	cmd := Command{Type: SHOW}
	return cmd
}

// parseType parses a type command.
// A type command takes a string to type.
//
// Type "string"
func (p *Parser) parseType() Command {
	cmd := Command{Type: TYPE}

	cmd.Options = p.parseSpeed()

	if p.peek.Type != STRING {
		p.errors = append(p.errors, NewError(p.peek, p.cur.Literal+" expects string"))
	}

	for p.peek.Type == STRING {
		p.nextToken()
		cmd.Args += p.cur.Literal

		// If the next token is a string, add a space between them.
		// Since tokens must be separated by a whitespace, this is most likely
		// what the user intended.
		//
		// Although it is possible that there may be multiple spaces / tabs between
		// the tokens, however if the user was intending to type multiple spaces
		// they would need to use a string literal.

		if p.peek.Type == STRING {
			cmd.Args += " "
		}
	}

	return cmd
}

// parseCopy parses a copy command
// A copy command takes a string to the clipboard
//
// Copy "string"
func (p *Parser) parseCopy() Command {
	cmd := Command{Type: COPY}

	if p.peek.Type != STRING {
		p.errors = append(p.errors, NewError(p.peek, p.cur.Literal+" expects string"))
	}
	for p.peek.Type == STRING {
		p.nextToken()
		cmd.Args += p.cur.Literal

		// If the next token is a string, add a space between them.
		// Since tokens must be separated by a whitespace, this is most likely
		// what the user intended.
		//
		// Although it is possible that there may be multiple spaces / tabs between
		// the tokens, however if the user was intending to type multiple spaces
		// they would need to use a string literal.

		if p.peek.Type == STRING {
			cmd.Args += " "
		}
	}
	return cmd
}

// parsePaste parses paste command
// Paste Command the string from the clipboard buffer.
//
// Paste
func (p *Parser) parsePaste() Command {
	cmd := Command{Type: PASTE}
	return cmd
}

// parseSource parses source command.
// Source command takes a tape path to include in current tape.
//
// Source <path>
func (p *Parser) parseSource() Command {
	cmd := Command{Type: SOURCE}

	if p.peek.Type != STRING {
		p.errors = append(p.errors, NewError(p.cur, "Expected path after Source"))
		p.nextToken()
		return cmd
	}

	srcPath := p.peek.Literal

	// Check if path has .tape extension
	ext := filepath.Ext(srcPath)
	if ext != ".tape" {
		p.errors = append(p.errors, NewError(p.peek, "Expected file with .tape extension"))
		p.nextToken()
		return cmd
	}

	// Check if tape exist
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		notFoundErr := fmt.Sprintf("File %s not found", srcPath)
		p.errors = append(p.errors, NewError(p.peek, notFoundErr))
		p.nextToken()
		return cmd
	}

	// Check if source tape contains nested Source command
	d, err := os.ReadFile(srcPath)
	if err != nil {
		readErr := fmt.Sprintf("Unable to read file: %s", srcPath)
		p.errors = append(p.errors, NewError(p.peek, readErr))
		p.nextToken()
		return cmd
	}

	srcTape := string(d)
	// Check source tape is NOT empty
	if len(srcTape) == 0 {
		readErr := fmt.Sprintf("Source tape: %s is empty", srcPath)
		p.errors = append(p.errors, NewError(p.peek, readErr))
		p.nextToken()
		return cmd
	}

	srcLexer := NewLexer(srcTape)
	srcParser := NewParser(srcLexer)

	// Check not nested source
	srcCmds := srcParser.Parse()
	for _, cmd := range srcCmds {
		if cmd.Type == SOURCE {
			p.errors = append(p.errors, NewError(p.peek, "Nested Source detected"))
			p.nextToken()
			return cmd
		}
	}

	// Check src errors
	srcErrors := srcParser.Errors()
	if len(srcErrors) > 0 {
		p.errors = append(p.errors, NewError(p.peek, fmt.Sprintf("%s has %d errors", srcPath, len(srcErrors))))
		p.nextToken()
		return cmd
	}

	cmd.Args = p.peek.Literal
	p.nextToken()
	return cmd
}

// parseScreenshot parses screenshot command.
// Screenshot command takes a file path for storing screenshot.
//
// Screenshot <path>
func (p *Parser) parseScreenshot() Command {
	cmd := Command{Type: SCREENSHOT}

	if p.peek.Type != STRING {
		p.errors = append(p.errors, NewError(p.cur, "Expected path after Screenshot"))
		p.nextToken()
		return cmd
	}

	path := p.peek.Literal

	// Check if path has .png extension
	ext := filepath.Ext(path)
	if ext != ".png" {
		p.errors = append(p.errors, NewError(p.peek, "Expected file with .png extension"))
		p.nextToken()
		return cmd
	}

	cmd.Args = path
	p.nextToken()

	return cmd
}

// Errors returns any errors that occurred during parsing.
func (p *Parser) Errors() []ParserError {
	return p.errors
}

// nextToken gets the next token from the lexer
// and updates the parser tokens accordingly.
func (p *Parser) nextToken() {
	p.cur = p.peek
	p.peek = p.l.NextToken()
}
