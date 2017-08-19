package rosewood

import (
	"fmt"
	"io"
	"log"
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
	errors       *ErrorManager
	lexer        *scanner.Scanner
	settings     *Settings
	debug        bool
	runMode      RunMode
	currentToken rune
	tables       []*tableContents //list of all loaded tables
}

//NewCommandParser initializes and returns a CommandParser
func NewCommandParser(settings *Settings) *CommandParser {
	p := CommandParser{errors: NewErrorManager()}
	//if no custom settings use default ones
	if settings == nil {
		panic("nil settings passed to command parser")
	}
	p.settings = settings
	p.runMode = settings.RunMode
	return &p
}

//defaultSettings: returns default settings in case no settings were set.
func (p *CommandParser) defaultSettings() error {
	p.settings = NewSettings()
	p.settings.RangeOperator = ':'
	if p.settings.Debug {
		log.Printf("default settings loaded")
	}
	return nil
}

//Errors returns a \n separated list of errors
func (p *CommandParser) Errors(index int) string {
	if index < 0 {
		return p.errors.String()
	}
	return p.errors.Errors[index].Error()
}

//nextToken: advances the lexer and updates currentToken of CommandParser. Do not
//confuse with Go's scanner.Scanner.Next()
func (p *CommandParser) nextToken() {
	p.currentToken = p.lexer.Scan()
	if p.settings.Debug {
		log.Printf("in nextToken: %s, current token= %q\n", p.lexer.Pos(), p.lexer.TokenText())
	}
}

//wrongToken adds an error into the parser's error list
func (p *CommandParser) wrongToken(wantedText string) {
	if strings.HasPrefix(wantedText, "*") {
		wantedText = wantedText[1:] //do not print the *
	}
	p.errors.Add(NewError(ErrSyntaxError, p.lexer.Pos(), fmt.Sprintf("expected %s, found %s (%s)",
		wantedText, scanner.TokenString(p.currentToken), p.lexer.TokenText())))
}

//wrongToken adds an error into the parser's error list
func (p *CommandParser) addSyntaxError(msg string) {
	p.errors.Add(NewError(ErrSyntaxError, p.lexer.Pos(), msg))
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
	if wantedText != strings.ToLower(p.lexer.TokenText()) {
		p.addSyntaxError(fmt.Sprintf("expected %s, found %s", wantedText, p.lexer.TokenText()))
	}
}

//acceptCommandName: reads and validates a command name
func (p *CommandParser) acceptCommandName() (string, rwKeyWord) {
	cmdName := strings.ToLower(p.lexer.TokenText())
	if p.currentToken != scanner.Ident {
		p.wrongToken("command")
		return cmdName, kwInvalid
	}
	cmd, found := lookupKeyword(cmdName)
	if !found {
		p.addSyntaxError(fmt.Sprintf("unknown command %s", cmdName))
		return cmdName, kwInvalid
	}
	return cmdName, cmd
}

//acceptCoordinate: reads and validates a row/cell coordinate
func (p *CommandParser) acceptCoordinate() RwInt {
	if p.currentToken != scanner.Int {
		p.wrongToken("integer")
		return MissingRwInt
	}
	coordinate, _ := strconv.Atoi(p.lexer.TokenText()) //no error check as we know it must be an int
	if coordinate < 1 {
		p.addSyntaxError(fmt.Sprintf("wanted row/col number > 0; found %s", p.lexer.TokenText()))
		return MissingRwInt
	}
	return RwInt(coordinate)
}

//acceptArgNameAndValue: reads an argument name and its value
func (p *CommandParser) acceptArg(lexeme rune) string {
	p.accept(lexeme, "*any identifier")
	return strings.ToLower(p.lexer.TokenText())
}

func (p *CommandParser) init(r io.Reader) error {
	if r == nil {
		panic("ParseCommands called with nil io.Reader")
	}
	p.errors.Reset()
	//initialize lexer and its settings
	p.lexer = new(scanner.Scanner).Init(r)
	p.lexer.Whitespace = 1<<'\t' | 1<<'\r' | 1<<' ' //ignore spaces, tabs and CRs

	return nil
}

//ParseCommands parses input stream and return an array of commands
func (p *CommandParser) ParseCommands(r io.Reader) ([]*Command, error) {
	if err := p.init(r); err != nil {
		return nil, err
	}
	var cmd *Command
	var cmdList []*Command

	p.nextToken()
	if p.currentToken == scanner.EOF {
		return nil, NewError(ErrEmpty, scanner.Position{Line: -1, Column: -1}, "nothing to parse")
	}
	cmdList = make([]*Command, 0, sectionCapacity)

	for ; p.currentToken != scanner.EOF; p.nextToken() {
		if p.currentToken == '\n' { //handle lines with no text, just linefeeds
			p.accept('\n', "*end of line") //advance beyond it and loop back
			continue
		}
		cmdName, cmdToken := p.acceptCommandName()
		cmd = NewCommand(cmdName, cmdToken, p.lexer.Pos())
		switch cmdName {
		case "set":
			p.parseSetCommand(cmd)
		default:
			p.parseTableFormatCommand(cmd) //all other commands will be parsed as a formatting command
		}
		if p.currentToken == '\n' {
			p.accept('\n', "*end of line")
		}
		err := cmd.Validate()
		if err != nil {
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

//parseTableFormatCommand: parses a command like command row col args
func (p *CommandParser) parseTableFormatCommand(cmd *Command) error {
	p.nextToken()
	p.accept(scanner.Ident, "row") //read row info
	p.nextToken()
	cmd.cellRange.TopLeft.Row = p.acceptCoordinate()
	cmd.cellSpan.r1 = cmd.cellRange.TopLeft.Row
	p.nextToken()
	//	fmt.Printf("in parseTable %q: %d\n", string(p.rangeOperator), p.rangeOperator)
	colExists := false
	switch p.currentToken {
	case '\n', scanner.EOF:
		// fmt.Printf("inside end of line or EOF: mandatory: %t, colexists: %t\n", p.settings.MandatoryCol, colExists)
		if p.settings.MandatoryCol && !colExists {
			p.addSyntaxError("col is missing")
		}
		return nil
	case p.settings.RangeOperator:
		p.nextToken()
		cmd.cellRange.BottomRight.Row = p.acceptCoordinate()
		cmd.cellSpan.r2 = cmd.cellRange.BottomRight.Row
		p.nextToken()
		fallthrough
	case scanner.Ident: //either "col" or an argument list
		switch strings.ToLower(p.lexer.TokenText()) {
		case "col": //read col info
			p.accept(scanner.Ident, "col")
			p.nextToken()
			cmd.cellRange.TopLeft.Col = p.acceptCoordinate()
			cmd.cellSpan.c1 = cmd.cellRange.TopLeft.Col
			colExists = true
			p.nextToken()
			//fmt.Println(p.rangeOperator, "-->", p.currentToken)
			if p.currentToken == p.settings.RangeOperator {
				p.nextToken()
				cmd.cellRange.BottomRight.Col = p.acceptCoordinate()
				cmd.cellSpan.c2 = cmd.cellRange.BottomRight.Col
				p.nextToken()
			}
			fallthrough
		default: //read args
			for ; p.currentToken != '\n' && p.currentToken != scanner.EOF; p.nextToken() {
				arg := p.acceptArg(scanner.Ident)
				cmd.args = append(cmd.args, arg)
			}
		}
	}
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
