package rosewood

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/scanner"
)

// keyword lookup
type rwKeyWord int

const (
	kwInvalid rwKeyWord = iota
	kwMerge
	kwStyle
	kwSet
	kwUse
)

var keywords = map[string]rwKeyWord{
	"merge": kwMerge,
	"style": kwStyle,
	"set":   kwSet,
	"use":   kwUse,
}

func lookupKeyword(name string) (kw rwKeyWord, isKeyWord bool) {
	kw, isKeyWord = keywords[name]
	return
}

// RunMode describes Parser's running mode
type RunMode int

const (
	Interactive RunMode = iota
	ScriptRun
)

// CommandParser specialized parser for format commands
type CommandParser struct {
	trace    tracer
	errors   *ErrorManager
	lexer    *scanner.Scanner
	settings *Settings
	//	debug        bool
	runMode      RunMode
	currentToken rune
	tables       []*tableContents //list of all loaded tables
}

//NewCommandParser initializes and returns a CommandParser
func NewCommandParser(settings *Settings) *CommandParser {
	if settings == nil {
		panic("nil settings passed to NewCommandParser")
	}
	p := CommandParser{errors: NewErrorManager()}
	p.settings = settings
	//	fmt.Println("trace is", p.settings.Debug)
	p.trace = newTrace(off, nil) //writer to stderr
	if p.settings.Debug {
		p.trace.On()
	}
	return &p
}

//Errors returns a \n separated list of errors if index =-1, otherwise returns the indexth error;
func (p *CommandParser) Errors(index int) string {
	if index < 0 {
		return p.errors.String()
	}
	return p.errors.Errors[index].Error() //intended: will panic if index out of range
}

//nextToken advances the lexer and updates currentToken of CommandParser.
func (p *CommandParser) nextToken() {
	p.currentToken = p.lexer.Scan()
	p.trace.Printf("in nextToken: %s, current token= %q\n", p.lexer.Pos(), p.currentWord())
}

//nextNotNull like nextToken but returns an error if EOL or EOF
func (p *CommandParser) nextNotNull() error {
	p.nextToken()
	switch p.currentToken {
	case scanner.EOF:
		return fmt.Errorf("unexpected end of input")
	case '\n':
		return fmt.Errorf("unexpected end of line")
	}
	return nil
}

//currentWord returns all-lower-case current token text.
func (p *CommandParser) currentWord() string {
	return strings.ToLower(p.lexer.TokenText())
}

//currentWord returns all-lower-case current token text.
func (p *CommandParser) exactCurrentWord() string {
	return p.lexer.TokenText()
}

//wrongToken adds an error into the parser's error list
func (p *CommandParser) wrongToken(wantedText string) {
	if strings.HasPrefix(wantedText, "*") {
		wantedText = wantedText[1:] //do not print the *
	}
	p.errors.Add(NewError(ErrSyntaxError, p.lexer.Pos(), fmt.Sprintf("expected %s, found %s (%s)", wantedText, scanner.TokenString(p.currentToken), p.currentWord())))
}

//addSyntaxError adds an error into the parser's error list
func (p *CommandParser) addSyntaxError(msg string, a ...interface{}) {
	p.errors.Add(NewError(ErrSyntaxError, p.lexer.Pos(), fmt.Sprintf(msg, a...)))
}

func (p *CommandParser) isToken(wantedTok rune, wantedText string) bool {
	if p.currentToken != wantedTok {
		return false
	}
	if strings.HasPrefix(wantedText, "*") || wantedText == p.currentWord() {
		return true
	}
	return false
}

//accept
// if wantedText =="*", wantedText is not validated
func (p *CommandParser) accept(wantedTok rune, wantedText string) {
	if p.currentToken != wantedTok {
		p.wrongToken(wantedText)
	}
	if strings.HasPrefix(wantedText, "*") {
		return
	}
	if wantedText != p.currentWord() {
		p.addSyntaxError("expected %s, found %s", wantedText, p.exactCurrentWord())
	}
}

//acceptCommandName: reads and validates a command name
func (p *CommandParser) acceptCommandName() (string, rwKeyWord) {
	p.nextToken()
	cmdName := p.currentWord()
	if p.currentToken != scanner.Ident {
		p.wrongToken("command")
		return cmdName, kwInvalid
	}
	cmd, found := lookupKeyword(cmdName)
	if !found {
		p.addSyntaxError("unknown command %s", cmdName)
		return cmdName, kwInvalid
	}
	return cmdName, cmd
}

//acceptArgNameAndValue: reads an argument name and its value
func (p *CommandParser) acceptArg(token rune) string {
	p.accept(token, "*any identifier")
	return p.currentWord()
}

//parsePoint: reads and validates a row/cell coordinate
func (p *CommandParser) parsePoint() (RwInt, error) {
	if err := p.nextNotNull(); err != nil {
		return MissingRwInt, err
	}
	if p.currentToken != scanner.Int {
		return MissingRwInt, fmt.Errorf("expected col or row number, found %s", p.exactCurrentWord())
	}
	coordinate, _ := strconv.Atoi(p.currentWord()) //no error check as we know it must be an int
	if coordinate < 1 {
		p.addSyntaxError("wanted row/col number > 0; found %s", p.exactCurrentWord()) //keep parsing
		return MissingRwInt, nil
	}
	return RwInt(coordinate), nil
}

//parseCommaSepPoints reads a list of coordinates points in the form x,y,z
func (p *CommandParser) parseCommaSepPoints(ss *subspan) error {
	for {
		//parser is on a comma, read a coordinate point
		point, err := p.parsePoint()
		if err != nil {
			return err
		}
		ss.list = append(ss.list, point)
		p.nextToken()
		switch p.currentToken {
		case ',':
			continue
		case p.settings.RangeOperator:
			if err := p.parseRangePoints(ss); err != nil {
				return err
			}
		default:
			return nil
		}
	}
}

//parseRangePoints read a range of coordinate either left:right or left:skipby:right
func (p *CommandParser) parseRangePoints(ss *subspan) error {
	var err error
	//p.nextToken() //skip ":"
	//read the right term of the range
	if ss.right, err = p.parsePoint(); err != nil {
		return err
	}
	p.nextToken()
	switch p.currentToken {
	case p.settings.RangeOperator: //so this a skipped range r:by:l
		ss.by = ss.right
		//read the right term of the range
		if ss.right, err = p.parsePoint(); err != nil {
			return err
		}
		p.nextToken()
	case ',':
		if err := p.parseCommaSepPoints(ss); err != nil {
			return err
		}
	default:
		return nil
	}
	return nil
}

func (p *CommandParser) parseRowOrColSegment(cmd *Command, segment string) (subspan, error) {
	var err error
	ss := newSubSpan()
	ss.kind = segment

	if ss.left, err = p.parsePoint(); err != nil {
		return ss, err
	}
	//now parse either ":", ",", either "col"/"row", an argument list or EOF
	p.nextToken()
	switch p.currentToken {
	case p.settings.RangeOperator:
		if err := p.parseRangePoints(&ss); err != nil {
			return ss, err
		}
	case ',':
		ss.list = append(ss.list, ss.left) //list has just begun
		ss.left = MissingRwInt
		if err := p.parseCommaSepPoints(&ss); err != nil {
			return ss, err
		}
	case scanner.Ident, scanner.EOF: //either "col"/"row" or an argument list or EOF
	default: //anything else is an error
		return ss, fmt.Errorf("unexpected %s", p.exactCurrentWord())
	}
	return ss, nil
}

//parseTableFormatCommand: parses a command like command row col args
func (p *CommandParser) parseTableFormatCommand(cmd *Command) error {
	//helper closure
	parseSegment := func(segment string) error {
		ss, err := p.parseRowOrColSegment(cmd, segment)
		if err == nil {
			cmd.spans = append(cmd.spans, &ss)
		}
		return err
	}
	//next, we must parse either a row or col modifier
	p.nextToken()
	modifier := p.currentWord()
	if modifier != "row" && modifier != "col" {
		return fmt.Errorf("expected row or col, found %s", p.exactCurrentWord())
	}
	if err := parseSegment(modifier); err != nil {
		return err
	}
	//next, we either have another col/row modifier, argument or EOF
	switch p.currentToken {
	case scanner.EOF:
	default:
		return fmt.Errorf("expected row, col or an argument, found %s", p.exactCurrentWord())
	case scanner.Ident:
		switch p.currentWord() {
		case "col", "row":
			if cmd.subSpan(p.currentWord()) != nil { //already parsed
				return fmt.Errorf("duplicate %s", p.exactCurrentWord())
			}
			if err := parseSegment(p.currentWord()); err != nil {
				return err
			}
			fallthrough
		default: //read args
			for ; p.currentToken == scanner.Ident; p.nextToken() {
				arg := p.acceptArg(scanner.Ident)
				cmd.args = append(cmd.args, arg)
			}
			if p.currentToken != scanner.EOF {
				return fmt.Errorf("expected row, col or an argument, found %s", p.exactCurrentWord())
			}
		}
	}
	//success
	cmd.cellSpan = subSpansToSpan(cmd.spans)
	p.trace.Printf("parsed: %v\n", cmd)
	return nil
}

//parseSetCommand: parses a command like "set settingname settingvalue"
func (p *CommandParser) parseSetCommand(cmd *Command) error {
	p.nextToken()
	settingName := p.acceptArg(scanner.Ident)
	p.nextToken()
	settingValue := p.acceptArg(scanner.String)
	p.nextToken()
	cmd.args = append(cmd.args, settingName, settingValue)
	return nil
}

//initialize lexer and its settings
func (p *CommandParser) init(r io.Reader) error {
	p.lexer = new(scanner.Scanner).Init(r)
	p.lexer.Whitespace = 1<<' ' | 1<<'\t' | 1<<'\r' //ignore spaces, tabs and CRs
	return nil
}

//ParseCommandLines parses a list of strings into list of commands
func (p *CommandParser) ParseCommandLines(Lines []string) ([]*Command, error) {
	if len(Lines) == 0 {
		return nil, nil
	}
	isCommandLine := func(line string) bool {
		line = strings.TrimSpace(line)
		if len(line) < 2 || strings.HasPrefix(line, "//") { //must have at least 2 chars and not a line comments
			return false
		}
		return true
	}
	var cmd *Command
	cmdList := make([]*Command, 0, len(Lines))
	p.errors.Reset()
	for i, line := range Lines {
		p.trace.Printf("src: :%v->", line)
		if !isCommandLine(line) {
			p.trace.Println("skipped")
			continue
		}
		p.trace.Println("to be parsed")
		p.init(strings.NewReader(line))
		cmdName, cmdToken := p.acceptCommandName()
		cmd = NewCommand(cmdName, cmdToken, scanner.Position{"", 0, i, 0}) //TODO: fix line numbers
		switch cmdName {
		case "set":
			p.parseSetCommand(cmd)
		default: //all other commands will be parsed as a formatting command
			if err := p.parseTableFormatCommand(cmd); err != nil {
				p.addSyntaxError(err.Error())
			}
		}
		if err := cmd.Validate(); err != nil {
			p.addSyntaxError(err.Error())
		}
		cmdList = append(cmdList, cmd)
	}
	if len(cmdList) == 0 {
		return nil, NewError(ErrEmpty, scanner.Position{Line: -1, Column: -1}, "found no valid commands")
	}
	if p.errors.Count() > 0 {
		return nil, NewError(ErrGeneric, scanner.Position{Line: -1, Column: -1}, "syntax errors")
	}
	return cmdList, nil
}
