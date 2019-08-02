// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

//Package parser implements a Rosewood file and command parsers
package parser

import (
	"fmt"
	"io"
	"strings"
	"text/scanner"
	"unicode"

	"github.com/drgo/core/errors"
	"github.com/drgo/rosewood/table"
	"github.com/drgo/rosewood/types"
)

// CommandParser specialized parser for format commands
type CommandParser struct {
	// trace        trace.Tracer
	errors *errors.ErrorList
	lexer  *scanner.Scanner
	job    *types.Job
	// settings     *types.RosewoodSettings
	position     Position
	currentToken rune
	tables       []*table.TableContents //list of all loaded tables
}

//NewCommandParser initializes and returns a CommandParser
func NewCommandParser(job *types.Job) *CommandParser {
	if job == nil {
		panic("nil settings passed to NewCommandParser")
	}
	p := CommandParser{errors: errors.NewErrorList(), lexer: new(scanner.Scanner)}
	p.job = job
	return &p
}

//Errors returns a list of parsing errors
func (p *CommandParser) Errors() *errors.ErrorList {
	return p.errors
}

//ErrorText returns a \n separated list of errors if index =-1, otherwise returns the indexth error;
func (p *CommandParser) ErrorText(index int) string {
	if index < 0 {
		return p.errors.Error()
	}
	return p.errors.Get(index).Error() //TODO: will panic if index is out of bound
}

//Pos returns the current position in the source
func (p *CommandParser) Pos() Position {
	p.position.Column = p.lexer.Pos().Column
	return p.position
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
		p.job.UI.Logf("src[%d]:%v->", i+s.Offset, line)
		if !isCommandLine(line) {
			p.job.UI.Log("skipped")
			continue
		}
		p.init(strings.NewReader(line))
		errOffset := p.errors.Len()
		cmdName, cmdToken := p.acceptCommandName()
		cmd := types.NewCommand(cmdName, cmdToken)
		switch cmdName {
		case "set":
			err = p.parseSetCommand(cmd)
		default: //all other commands will be parsed as a formatting command
			err = p.parseTableFormatCommand(cmd)
		}
		if err != nil {
			p.addSyntaxError(err.Error())
			p.job.UI.Logf("parsing error: %s\n", err.Error())
		}
		p.job.UI.Logf(" %v\n", cmd)
		if p.errors.Len() > errOffset { //errors were detected during parsing; stop here
			continue
		}
		if err = cmd.Finalize(); err != nil {
			p.addSyntaxError(err.Error())
		}
		cmdList = append(cmdList, cmd)
	}
	p.job.UI.Log("")
	switch {
	case p.errors.Len() > 0:
		return nil, p.errors.Err()
	case len(cmdList) == 0:
		return nil, NewError(ErrSyntaxError, unknownPos, "found no valid commands")
	default:
		return cmdList, nil
	}
}

//scannerErrorHandler captures scanner errors
func (p *CommandParser) scannerErrorHandler(s *scanner.Scanner, msg string) {
	p.addSyntaxError(msg)
}

//initialize lexer and its settings
func (p *CommandParser) init(r io.Reader) error {
	p.lexer = p.lexer.Init(r)
	p.lexer.Whitespace = 1<<' ' | 1<<'\t' | 1<<'\r' //ignore spaces, tabs and CRs
	// treat '-' as part of an identifier (eg css stylenames)
	p.lexer.IsIdentRune = func(ch rune, i int) bool {
		return ch == '-' || unicode.IsLetter(ch) || unicode.IsDigit(ch) && i > 0
	}
	p.lexer.Error = p.scannerErrorHandler
	return nil
}

//parseSetCommand parses a command like "set settingname settingvalue"
func (p *CommandParser) parseSetCommand(cmd *types.Command) error {
	p.nextToken()
	settingName := p.acceptArg(scanner.Ident)
	p.nextToken()
	settingValue := p.acceptArg(scanner.String)
	p.nextToken()
	cmd.AddArg(settingName, settingValue)
	return nil
}

//acceptCommandName reads and validates a command name
func (p *CommandParser) acceptCommandName() (string, types.RwKeyWord) {
	p.nextToken()
	cmdName := p.currentWord()
	if p.currentToken != scanner.Ident {
		p.wrongToken("command")
		return cmdName, types.KwInvalid
	}
	cmd, found := types.LookupKeyword(cmdName)
	if !found {
		p.errors.Add(NewError(ErrSyntaxError, p.Pos(), fmt.Sprintf("unknown command %s", cmdName)))
		return cmdName, types.KwInvalid
	}
	return cmdName, cmd
}

//nextToken advances the lexer and updates currentToken of CommandParser.
func (p *CommandParser) nextToken() {
	p.currentToken = p.lexer.Scan()
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

//accept consumes and validates a token and optionally the token's text
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
