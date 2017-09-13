//Package parser implements a Rosewood file and command parsers
package parser

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/drgo/rosewood/types"
	"github.com/drgo/rosewood/utils"
)

// RunMode describes Parser's running mode
// type RunMode int

// const (
// 	Interactive RunMode = iota
// 	ScriptRun
// )

// CommandParser specialized parser for format commands
type CommandParser struct {
	trace        utils.Tracer
	errors       utils.ErrorList
	lexer        *scanner.Scanner
	settings     *utils.Settings
	position     scanner.Position
	currentToken rune
	tables       []*types.TableContents //list of all loaded tables
	//	runMode      RunMode
}

//NewCommandParser initializes and returns a CommandParser
func NewCommandParser(settings *utils.Settings) *CommandParser {
	if settings == nil {
		panic("nil settings passed to NewCommandParser")
	}
	p := CommandParser{errors: utils.NewErrorList(), lexer: new(scanner.Scanner)}
	p.settings = settings
	//	fmt.Println("trace is", p.settings.Debug)
	p.trace = utils.NewTrace(false, nil) //writer to stderr
	if p.settings.Debug {
		p.trace.On()
	}
	return &p
}

//Errors returns a list of parsing errors
func (p *CommandParser) Errors() []error {
	return p.errors
}

//ErrorText returns a \n separated list of errors if index =-1, otherwise returns the indexth error;
func (p *CommandParser) ErrorText(index int) string {
	if index < 0 {
		return p.errors.Error()
	}
	return p.errors[index].Error() //intended: will panic if index out of range
}

//Pos returns the current position in the source
func (p *CommandParser) Pos() scanner.Position {
	p.position.Column = p.lexer.Pos().Column
	return p.position
}

//nextToken advances the lexer and updates currentToken of CommandParser.
func (p *CommandParser) nextToken() {
	p.currentToken = p.lexer.Scan()
	//check p.lexer.ErrorCount() ==
	//	p.trace.Printf("in nextToken: %s, current token= %q\n", p.Pos(), p.currentWord())
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

//exactCurrentWord returns the current token text or a description if no text.
func (p *CommandParser) exactCurrentWord() string {
	s := strings.TrimSpace(p.lexer.TokenText())
	if s == "" {
		s = scanner.TokenString(p.currentToken)
	}
	return s
}

//wrongToken adds an error into the parser's error list
func (p *CommandParser) wrongToken(wantedText string) {
	if strings.HasPrefix(wantedText, "*") {
		wantedText = wantedText[1:] //do not print the *
	}
	p.errors.Add(NewError(ErrSyntaxError, p.Pos(), fmt.Sprintf("expected %s, found %s (%s)", wantedText, scanner.TokenString(p.currentToken), p.currentWord())))
}

//addSyntaxError adds an error into the parser's error list
func (p *CommandParser) addSyntaxError(msg string, a ...interface{}) {
	p.errors.Add(NewError(ErrSyntaxError, p.Pos(), fmt.Sprintf(msg, a...)))
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
func (p *CommandParser) acceptCommandName() (string, types.RwKeyWord) {
	p.nextToken()
	cmdName := p.currentWord()
	if p.currentToken != scanner.Ident {
		p.wrongToken("command")
		return cmdName, types.KwInvalid
	}
	cmd, found := types.LookupKeyword(cmdName)
	if !found {
		p.addSyntaxError("unknown command %s", cmdName)
		return cmdName, types.KwInvalid
	}
	return cmdName, cmd
}

//acceptArgNameAndValue: reads an argument name and its value
func (p *CommandParser) acceptArg(token rune) string {
	p.accept(token, "*any identifier")
	return p.currentWord()
}

//parsePoint: reads and validates a row/cell coordinate
func (p *CommandParser) parsePoint() (types.RwInt, error) {
	if err := p.nextNotNull(); err != nil {
		return types.MissingRwInt, err
	}
	if p.currentToken != scanner.Int {
		return types.MissingRwInt, fmt.Errorf("expected col or row number, found %s", p.exactCurrentWord())
	}
	coordinate, _ := strconv.Atoi(p.currentWord()) //no error check as we know it must be an int
	if coordinate < 1 {
		p.addSyntaxError("wanted row/col number > 0; found %s", p.exactCurrentWord()) //keep parsing
		return types.MissingRwInt, nil
	}
	return types.RwInt(coordinate), nil
}

//parseCommaSepPoints reads a list of coordinates points in the form x,y,z
func (p *CommandParser) parseCommaSepPoints(ss *types.Subspan) error {
	for {
		//parser is on a comma, read a coordinate point
		point, err := p.parsePoint()
		if err != nil {
			return err
		}
		ss.List = append(ss.List, point)
		p.nextToken()
		switch p.currentToken {
		case ',':
			continue
		case p.settings.RangeOperator: //range after a comma-list is not allowed
			return fmt.Errorf("a ranger operator [:] is not allowed following a coordinate list")
		default:
			return nil
		}
	}
}

//parseRangePoints read a range of coordinate either left:right or left:skipby:right
func (p *CommandParser) parseRangePoints(ss *types.Subspan) error {
	var err error
	//read the right term of the range
	if ss.Right, err = p.parsePoint(); err != nil {
		return err
	}
	p.nextToken()
	switch p.currentToken {
	case p.settings.RangeOperator: //so this a skipped range r:by:l
		ss.By = ss.Right
		//read the right term of the range
		if ss.Right, err = p.parsePoint(); err != nil {
			return err
		}
		p.nextToken()
		if p.currentToken == ',' {
			if err := p.parseCommaSepPoints(ss); err != nil {
				return err
			}
		}
	case ',':
		if err := p.parseCommaSepPoints(ss); err != nil {
			return err
		}
	default:
		return nil
	}
	return nil
}

func (p *CommandParser) parseRowOrColSegment(cmd *types.Command, segment string) (types.Subspan, error) {
	var err error
	ss := types.NewSubSpan(segment)

	if ss.Left, err = p.parsePoint(); err != nil {
		return ss, err
	}
	p.nextToken()
	switch p.currentToken {
	case p.settings.RangeOperator:
		if err := p.parseRangePoints(&ss); err != nil {
			return ss, err
		}
	case ',':
		ss.List = append(ss.List, ss.Left) //list has just begun
		ss.Left = types.MissingRwInt
		if err := p.parseCommaSepPoints(&ss); err != nil {
			return ss, err
		}
	case scanner.Ident, scanner.EOF: //either "col"/"row" or an argument list or EOF
	default: //anything else is an error
		return ss, fmt.Errorf("unexpected token: %s", p.exactCurrentWord())
	}

	return ss, nil
}

//parseTableFormatCommand: parses a command like command row col args
func (p *CommandParser) parseTableFormatCommand(cmd *types.Command) error {
	//helper closure
	parseSegment := func(segment string) error {
		ss, err := p.parseRowOrColSegment(cmd, segment)
		if err == nil {
			cmd.AddSubSpan(&ss)
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
			if cmd.SubSpan(p.currentWord()) != nil { //already parsed
				return fmt.Errorf("duplicate %s", p.exactCurrentWord())
			}
			if err := parseSegment(p.currentWord()); err != nil {
				return err
			}
			fallthrough
		default: //read args
			for ; p.currentToken == scanner.Ident; p.nextToken() {
				arg := p.acceptArg(scanner.Ident)
				cmd.AddArg(arg)
			}
			if p.currentToken != scanner.EOF {
				return fmt.Errorf("expected row, col or an argument, found %s", p.exactCurrentWord())
			}
		}
	}
	//success
	return nil
}

//parseSetCommand: parses a command like "set settingname settingvalue"
func (p *CommandParser) parseSetCommand(cmd *types.Command) error {
	p.nextToken()
	settingName := p.acceptArg(scanner.Ident)
	p.nextToken()
	settingValue := p.acceptArg(scanner.String)
	p.nextToken()
	cmd.AddArg(settingName, settingValue)
	return nil
}

func (p *CommandParser) scannerErrorHandler(s *scanner.Scanner, msg string) {
	p.addSyntaxError(msg)
}

//initialize lexer and its settings
func (p *CommandParser) init(r io.Reader) error {
	p.lexer = p.lexer.Init(r)
	p.lexer.Whitespace = 1<<' ' | 1<<'\t' | 1<<'\r' //ignore spaces, tabs and CRs
	p.lexer.Error = p.scannerErrorHandler
	return nil
}

//ParseCommandLines parses a list of strings into list of commands
func (p *CommandParser) ParseCommandLines(s *types.Section) ([]*types.Command, error) {
	if len(s.Lines) == 0 {
		return nil, nil
	}
	isCommandLine := func(line string) bool {
		line = strings.TrimSpace(line)
		if len(line) < 2 || strings.HasPrefix(line, "//") { //must have at least 2 chars and not a line comments
			return false
		}
		return true
	}
	cmdList := make([]*types.Command, 0, len(s.Lines))
	p.errors.Reset()
	p.position.Offset = s.Offset
	var err error
	for i, line := range s.Lines {
		p.position.Line = i
		p.trace.Printf("src[%d]: :%v->", i+s.Offset, line)
		if !isCommandLine(line) {
			p.trace.Println("skipped")
			continue
		}
		p.trace.Println("to be parsed")
		p.init(strings.NewReader(line))
		errOffset := p.errors.Count()
		cmdName, cmdToken := p.acceptCommandName()
		cmd := types.NewCommand(cmdName, cmdToken, p.Pos())
		switch cmdName {
		case "set":
			err = p.parseSetCommand(cmd)
		default: //all other commands will be parsed as a formatting command
			err = p.parseTableFormatCommand(cmd)
		}
		if err != nil {
			p.addSyntaxError(err.Error())
			p.trace.Printf("parsing error: %s\n", err.Error())
		}
		p.trace.Printf("parsed: %v\n", cmd)
		if p.errors.Count() > errOffset { //errors were detected during parsing; stop here
			continue
		}
		if err = cmd.Finalize(); err != nil {
			p.addSyntaxError(err.Error())
		}
		cmdList = append(cmdList, cmd)
	}
	//signal status to caller
	if p.errors.Count() > 0 {
		return nil, NewError(ErrSyntaxError, unknownPos, "syntax errors")
	}
	if len(cmdList) == 0 {
		return nil, NewError(ErrEmpty, unknownPos, "found no valid commands")
	}
	return cmdList, nil
}
