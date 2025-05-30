package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/vhs/lexer"
	"github.com/charmbracelet/vhs/token"
)

// NewError returns a new parser.Error with the given token and message.
func NewError(token token.Token, msg string) Error {
	return Error{
		Token: token,
		Msg:   msg,
	}
}

// CommandType is a type that represents a command.
type CommandType token.Type

// CommandTypes is a list of the available commands that can be executed.
var CommandTypes = []CommandType{ //nolint: deadcode
	token.BACKSPACE,
	token.DELETE,
	token.INSERT,
	token.CTRL,
	token.ALT,
	token.DOWN,
	token.ENTER,
	token.ESCAPE,
	token.ILLEGAL,
	token.LEFT,
	token.PAGE_UP,
	token.PAGE_DOWN,
	token.RIGHT,
	token.SET,
	token.OUTPUT,
	token.SLEEP,
	token.SPACE,
	token.HIDE,
	token.REQUIRE,
	token.SHOW,
	token.TAB,
	token.TYPE,
	token.UP,
	token.WAIT,
	token.SOURCE,
	token.SCREENSHOT,
	token.COPY,
	token.PASTE,
	token.ENV,
}

// String returns the string representation of the command.
func (c CommandType) String() string { return token.ToCamel(string(c)) }

// Command represents a command with options and arguments.
type Command struct {
	Type    CommandType
	Options string
	Args    string
	Source  string
}

// String returns the string representation of the command.
// This includes the options and arguments of the command.
func (c Command) String() string {
	if c.Options != "" {
		return fmt.Sprintf("%s %s %s", c.Type, c.Options, c.Args)
	}
	return fmt.Sprintf("%s %s", c.Type, c.Options)
}

// Error represents an error with parsing a tape file.
// It tracks the token causing the error and a human readable error message.
type Error struct {
	Token token.Token
	Msg   string
}

// String returns a human readable error message printing the token line number
// and message.
func (e Error) String() string {
	return fmt.Sprintf("%2d:%-2d │ %s", e.Token.Line, e.Token.Column, e.Msg)
}

// Error implements the error interface.
func (e Error) Error() string {
	return e.String()
}

// Parser is the structure that manages the parsing of tokens.
type Parser struct {
	l      *lexer.Lexer
	errors []Error
	cur    token.Token
	peek   token.Token
}

// New returns a new Parser.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []Error{}}

	// Read two tokens, so cur and peek are both set.
	p.nextToken()
	p.nextToken()

	return p
}

// Parse takes an input string provided by the lexer and parses it into a
// list of commands.
func (p *Parser) Parse() []Command {
	cmds := []Command{}

	for p.cur.Type != token.EOF {
		if p.cur.Type == token.COMMENT {
			p.nextToken()
			continue
		}
		cmds = append(cmds, p.parseCommand()...)
		p.nextToken()
	}

	return cmds
}

// parseCommand parses a command.
func (p *Parser) parseCommand() []Command {
	switch p.cur.Type {
	case token.SPACE,
		token.BACKSPACE,
		token.DELETE,
		token.INSERT,
		token.ENTER,
		token.ESCAPE,
		token.TAB,
		token.DOWN,
		token.LEFT,
		token.RIGHT,
		token.UP,
		token.PAGE_UP,
		token.PAGE_DOWN:
		return []Command{p.parseKeypress(p.cur.Type)}
	case token.SET:
		return []Command{p.parseSet()}
	case token.OUTPUT:
		return []Command{p.parseOutput()}
	case token.SLEEP:
		return []Command{p.parseSleep()}
	case token.TYPE:
		return []Command{p.parseType()}
	case token.CTRL:
		return p.parseCtrl()
	case token.ALT:
		return []Command{p.parseAlt()}
	case token.SHIFT:
		return []Command{p.parseShift()}
	case token.HIDE:
		return []Command{p.parseHide()}
	case token.REQUIRE:
		return []Command{p.parseRequire()}
	case token.SHOW:
		return []Command{p.parseShow()}
	case token.WAIT:
		return []Command{p.parseWait()}
	case token.SOURCE:
		return p.parseSource()
	case token.SCREENSHOT:
		return []Command{p.parseScreenshot()}
	case token.COPY:
		return []Command{p.parseCopy()}
	case token.PASTE:
		return []Command{p.parsePaste()}
	case token.ENV:
		return []Command{p.parseEnv()}
	default:
		p.errors = append(p.errors, NewError(p.cur, "Invalid command: "+p.cur.Literal))
		return []Command{{Type: token.ILLEGAL}}
	}
}

func (p *Parser) parseWait() Command {
	cmd := Command{Type: token.WAIT}

	if p.peek.Type == token.PLUS {
		p.nextToken()
		if p.peek.Type != token.STRING || (p.peek.Literal != "Line" && p.peek.Literal != "Screen") {
			p.errors = append(p.errors, NewError(p.peek, "Wait+ expects Line or Screen"))
			return cmd
		}
		cmd.Args = p.peek.Literal
		p.nextToken()
	} else {
		cmd.Args = "Line"
	}

	cmd.Options = p.parseSpeed()
	if cmd.Options != "" {
		dur, _ := time.ParseDuration(cmd.Options)
		if dur <= 0 {
			p.errors = append(p.errors, NewError(p.peek, "Wait expects positive duration"))
			return cmd
		}
	}

	if p.peek.Type != token.REGEX {
		// fallback to default
		return cmd
	}
	p.nextToken()
	if _, err := regexp.Compile(p.cur.Literal); err != nil {
		p.errors = append(p.errors, NewError(p.cur, fmt.Sprintf("Invalid regular expression '%s': %v", p.cur.Literal, err)))
		return cmd
	}

	cmd.Args += " " + p.cur.Literal

	return cmd
}

// parseSpeed parses a typing speed indication.
//
// i.e. @<time>
//
// This is optional (defaults to 100ms), thus skips (rather than error-ing)
// if the typing speed is not specified.
func (p *Parser) parseSpeed() string {
	if p.peek.Type == token.AT {
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
	if p.peek.Type == token.NUMBER {
		count := p.peek.Literal
		p.nextToken()
		return count
	}

	return "1"
}

// parseTime parses a time argument.
//
//	<number>[ms]
func (p *Parser) parseTime() string {
	var t string

	if p.peek.Type == token.NUMBER {
		t = p.peek.Literal
		p.nextToken()
	} else {
		p.errors = append(p.errors, NewError(p.cur, "Expected time after "+p.cur.Literal))
		return ""
	}

	// Allow TypingSpeed to have bare units (e.g. 50ms, 100ms)
	if p.peek.Type == token.MILLISECONDS || p.peek.Type == token.SECONDS || p.peek.Type == token.MINUTES {
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
//	Ctrl[+Alt][+Shift]+<char>
//	E.g:
//	Ctrl+Shift+O
//	Ctrl+Alt+Shift+P
func (p *Parser) parseCtrl() []Command {
	var args []string

	inModifierChain := true
	for p.peek.Type == token.PLUS {
		p.nextToken()
		peek := p.peek

		// Get key from keywords and check if it's a valid modifier
		if k := token.Keywords[peek.Literal]; token.IsModifier(k) {
			if !inModifierChain {
				p.errors = append(p.errors, NewError(p.cur, "Modifiers must come before other characters"))
				// Clear args so the error is returned
				args = nil
				continue
			}

			args = append(args, peek.Literal)
			p.nextToken()
			continue
		}

		inModifierChain = false

		// Add key argument.
		switch {
		case peek.Type == token.ENTER,
			peek.Type == token.SPACE,
			peek.Type == token.BACKSPACE,
			peek.Type == token.MINUS,
			peek.Type == token.AT,
			peek.Type == token.LEFT_BRACKET,
			peek.Type == token.RIGHT_BRACKET,
			peek.Type == token.CARET,
			peek.Type == token.BACKSLASH,
			peek.Type == token.STRING && len(peek.Literal) == 1:
			args = append(args, peek.Literal)
		default:
			// Key arguments with len > 1 are not valid
			p.errors = append(p.errors,
				NewError(p.cur, "Not a valid modifier"),
				NewError(p.cur, "Invalid control argument: "+p.cur.Literal))
		}

		p.nextToken()
	}

	if len(args) == 0 {
		p.errors = append(p.errors, NewError(p.cur, "Expected control character with args, got "+p.cur.Literal))
	}

	ctrlArgs := strings.Join(args, " ")
	repeat, _ := strconv.Atoi(p.parseRepeat())

	cmds := make([]Command, 0, repeat)
	for range repeat {
		cmds = append(cmds, Command{Type: token.CTRL, Args: ctrlArgs})
	}
	return cmds
}

// parseAlt parses an alt command.
// An alt command takes a character to type while the modifier is held down.
//
//	Alt+<character>
func (p *Parser) parseAlt() Command {
	if p.peek.Type == token.PLUS {
		p.nextToken()
		if p.peek.Type == token.STRING ||
			p.peek.Type == token.ENTER ||
			p.peek.Type == token.LEFT_BRACKET ||
			p.peek.Type == token.RIGHT_BRACKET ||
			p.peek.Type == token.TAB {
			c := p.peek.Literal
			p.nextToken()
			return Command{Type: token.ALT, Args: c}
		}
	}

	p.errors = append(p.errors, NewError(p.cur, "Expected alt character, got "+p.cur.Literal))
	return Command{Type: token.ALT}
}

// parseShift parses a shift command.
// A shift command takes one character and types while shift is held down.
//
//	Shift+<char>
//	E.g.
//	Shift+A
//	Shift+Tab
//	Shift+Enter
func (p *Parser) parseShift() Command {
	if p.peek.Type == token.PLUS {
		p.nextToken()
		if p.peek.Type == token.STRING ||
			p.peek.Type == token.ENTER ||
			p.peek.Type == token.LEFT_BRACKET ||
			p.peek.Type == token.RIGHT_BRACKET ||
			p.peek.Type == token.TAB {
			c := p.peek.Literal
			p.nextToken()
			return Command{Type: token.SHIFT, Args: c}
		}
	}

	p.errors = append(p.errors, NewError(p.cur, "Expected shift character, got "+p.cur.Literal))
	return Command{Type: token.SHIFT}
}

// parseKeypress parses a repeatable and time adjustable keypress command.
// A keypress command takes an optional typing speed and optional count.
//
//	Key[@<time>] [count]
func (p *Parser) parseKeypress(ct token.Type) Command {
	cmd := Command{Type: CommandType(ct)}
	cmd.Options = p.parseSpeed()
	cmd.Args = p.parseRepeat()
	return cmd
}

// parseOutput parses an output command.
// An output command takes a file path to which to output.
//
//	Output <path>
func (p *Parser) parseOutput() Command {
	cmd := Command{Type: token.OUTPUT}

	if p.peek.Type != token.STRING {
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
//	Set <setting> <value>
func (p *Parser) parseSet() Command {
	cmd := Command{Type: token.SET}

	if token.IsSetting(p.peek.Type) {
		cmd.Options = p.peek.Literal
	} else {
		p.errors = append(p.errors, NewError(p.peek, "Unknown setting: "+p.peek.Literal))
	}
	p.nextToken()

	switch p.cur.Type {
	case token.WAIT_TIMEOUT:
		cmd.Args = p.parseTime()
	case token.WAIT_PATTERN:
		cmd.Args = p.peek.Literal
		_, err := regexp.Compile(p.peek.Literal)
		if err != nil {
			p.errors = append(p.errors, NewError(p.peek, "Invalid regexp pattern: "+p.peek.Literal))
		}
		p.nextToken()
	case token.LOOP_OFFSET:
		cmd.Args = p.peek.Literal
		p.nextToken()
		// Allow LoopOffset without '%'
		// Set LoopOffset 20
		cmd.Args += "%"
		if p.peek.Type == token.PERCENT {
			p.nextToken()
		}
	case token.TYPING_SPEED:
		cmd.Args = p.peek.Literal
		p.nextToken()
		// Allow TypingSpeed to have bare units (e.g. 10ms)
		// Set TypingSpeed 10ms
		if p.peek.Type == token.MILLISECONDS ||
			p.peek.Type == token.SECONDS {
			cmd.Args += p.peek.Literal
			p.nextToken()
		} else if cmd.Options == "TypingSpeed" {
			cmd.Args += "s"
		}
	case token.WINDOW_BAR:
		cmd.Args = p.peek.Literal
		p.nextToken()

		windowBar := p.cur.Literal
		if !isValidWindowBar(windowBar) {
			p.errors = append(
				p.errors,
				NewError(p.cur, windowBar+" is not a valid bar style."),
			)
		}
	case token.MARGIN_FILL:
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
	case token.CURSOR_BLINK:
		cmd.Args = p.peek.Literal
		p.nextToken()

		if p.cur.Type != token.BOOLEAN {
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
//	Sleep <time>
func (p *Parser) parseSleep() Command {
	cmd := Command{Type: token.SLEEP}
	cmd.Args = p.parseTime()
	return cmd
}

// parseHide parses a Hide command.
//
//	Hide
func (p *Parser) parseHide() Command {
	cmd := Command{Type: token.HIDE}
	return cmd
}

// parseRequire parses a Require command.
//
//	Require
func (p *Parser) parseRequire() Command {
	cmd := Command{Type: token.REQUIRE}

	if p.peek.Type != token.STRING {
		p.errors = append(p.errors, NewError(p.peek, p.cur.Literal+" expects one string"))
	}

	cmd.Args = p.peek.Literal
	p.nextToken()

	return cmd
}

// parseShow parses a Show command.
//
//	Show
func (p *Parser) parseShow() Command {
	cmd := Command{Type: token.SHOW}
	return cmd
}

// parseType parses a type command.
// A type command takes a string to type.
//
//	Type "string"
func (p *Parser) parseType() Command {
	cmd := Command{Type: token.TYPE}

	cmd.Options = p.parseSpeed()

	if p.peek.Type != token.STRING {
		p.errors = append(p.errors, NewError(p.peek, p.cur.Literal+" expects string"))
	}

	for p.peek.Type == token.STRING {
		p.nextToken()
		cmd.Args += p.cur.Literal

		// If the next token is a string, add a space between them.
		// Since tokens must be separated by a whitespace, this is most likely
		// what the user intended.
		//
		// Although it is possible that there may be multiple spaces / tabs between
		// the tokens, however if the user was intending to type multiple spaces
		// they would need to use a string literal.

		if p.peek.Type == token.STRING {
			cmd.Args += " "
		}
	}

	return cmd
}

// parseCopy parses a copy command
// A copy command takes a string to the clipboard
//
//	Copy "string"
func (p *Parser) parseCopy() Command {
	cmd := Command{Type: token.COPY}

	if p.peek.Type != token.STRING {
		p.errors = append(p.errors, NewError(p.peek, p.cur.Literal+" expects string"))
	}
	for p.peek.Type == token.STRING {
		p.nextToken()
		cmd.Args += p.cur.Literal

		// If the next token is a string, add a space between them.
		// Since tokens must be separated by a whitespace, this is most likely
		// what the user intended.
		//
		// Although it is possible that there may be multiple spaces / tabs between
		// the tokens, however if the user was intending to type multiple spaces
		// they would need to use a string literal.

		if p.peek.Type == token.STRING {
			cmd.Args += " "
		}
	}
	return cmd
}

// parsePaste parses paste command
// Paste Command the string from the clipboard buffer.
//
//	Paste
func (p *Parser) parsePaste() Command {
	cmd := Command{Type: token.PASTE}
	return cmd
}

// parseEnv parses Env command
// Env command takes in a key-value pair which is set.
//
//	Env key "value"
func (p *Parser) parseEnv() Command {
	cmd := Command{Type: token.ENV}

	cmd.Options = p.peek.Literal
	p.nextToken()

	if p.peek.Type != token.STRING {
		p.errors = append(p.errors, NewError(p.peek, p.cur.Literal+" expects string"))
	}

	cmd.Args = p.peek.Literal
	p.nextToken()

	return cmd
}

// parseSource parses source command.
// Source command takes a tape path to include in current tape.
//
//	Source <path>
func (p *Parser) parseSource() []Command {
	cmd := Command{Type: token.SOURCE}

	if p.peek.Type != token.STRING {
		p.errors = append(p.errors, NewError(p.cur, "Expected path after Source"))
		p.nextToken()
		return []Command{cmd}
	}

	srcPath := p.peek.Literal

	// Check if path has .tape extension
	ext := filepath.Ext(srcPath)
	if ext != ".tape" {
		p.errors = append(p.errors, NewError(p.peek, "Expected file with .tape extension"))
		p.nextToken()
		return []Command{cmd}
	}

	// Check if tape exist
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		notFoundErr := fmt.Sprintf("File %s not found", srcPath)
		p.errors = append(p.errors, NewError(p.peek, notFoundErr))
		p.nextToken()
		return []Command{cmd}
	}

	// Check if source tape contains nested Source command
	d, err := os.ReadFile(srcPath)
	if err != nil {
		readErr := fmt.Sprintf("Unable to read file: %s", srcPath)
		p.errors = append(p.errors, NewError(p.peek, readErr))
		p.nextToken()
		return []Command{cmd}
	}

	srcTape := string(d)
	// Check source tape is NOT empty
	if len(srcTape) == 0 {
		readErr := fmt.Sprintf("Source tape: %s is empty", srcPath)
		p.errors = append(p.errors, NewError(p.peek, readErr))
		p.nextToken()
		return []Command{cmd}
	}

	srcLexer := lexer.New(srcTape)
	srcParser := New(srcLexer)

	// Check not nested source
	srcCmds := srcParser.Parse()
	for _, cmd := range srcCmds {
		if cmd.Type == token.SOURCE {
			p.errors = append(p.errors, NewError(p.peek, "Nested Source detected"))
			p.nextToken()
			return []Command{cmd}
		}
	}

	// Check src errors
	srcErrors := srcParser.Errors()
	if len(srcErrors) > 0 {
		p.errors = append(p.errors, NewError(p.peek, fmt.Sprintf("%s has %d errors", srcPath, len(srcErrors))))
		p.nextToken()
		return []Command{cmd}
	}

	filtered := make([]Command, 0)
	for _, srcCmd := range srcCmds {
		// Output have to be avoid in order to not overwrite output of the original tape.
		if srcCmd.Type == token.SOURCE ||
			srcCmd.Type == token.OUTPUT {
			continue
		}
		filtered = append(filtered, srcCmd)
	}

	p.nextToken()
	return filtered
}

// parseScreenshot parses screenshot command.
// Screenshot command takes a file path for storing screenshot.
//
//	Screenshot <path>
func (p *Parser) parseScreenshot() Command {
	cmd := Command{Type: token.SCREENSHOT}

	if p.peek.Type != token.STRING {
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
func (p *Parser) Errors() []Error {
	return p.errors
}

// nextToken gets the next token from the lexer
// and updates the parser tokens accordingly.
func (p *Parser) nextToken() {
	p.cur = p.peek
	p.peek = p.l.NextToken()
}

// Check if a given windowbar type is valid.
func isValidWindowBar(w string) bool {
	return w == "" ||
		w == "Colorful" || w == "ColorfulRight" ||
		w == "Rings" || w == "RingsRight"
}
