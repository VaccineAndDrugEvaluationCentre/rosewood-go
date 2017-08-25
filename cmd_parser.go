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

//nextToken advances the lexer and updates currentToken of CommandParser. Do not
//confuse with Go's scanner.Scanner.Next()
func (p *CommandParser) nextToken() {
	p.currentToken = p.lexer.Scan()
	if p.settings.Debug {
		log.Printf("in nextToken: %s, current token= %q\n", p.lexer.Pos(), p.lexer.TokenText())
	}
}

//nextNotNull like nextToken but returns an error if EOL or EOF
func (p *CommandParser) nextNotNull() error {
	p.nextToken()
	switch p.currentToken {
	case scanner.EOF:
		return fmt.Errorf("unexpected end of input")
	case '\n': //todo: use CRLF as well
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
	p.errors.Add(NewError(ErrSyntaxError, p.lexer.Pos(), fmt.Sprintf("expected %s, found %s (%s)",
		wantedText, scanner.TokenString(p.currentToken), p.lexer.TokenText())))
}

//wrongToken adds an error into the parser's error list
func (p *CommandParser) addSyntaxError(msg string) {
	p.errors.Add(NewError(ErrSyntaxError, p.lexer.Pos(), msg))
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
	if wantedText != strings.ToLower(p.lexer.TokenText()) {
		p.addSyntaxError(fmt.Sprintf("expected %s, found %s", wantedText, p.lexer.TokenText()))
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
		p.addSyntaxError(fmt.Sprintf("unknown command %s", cmdName))
		return cmdName, kwInvalid
	}
	return cmdName, cmd
}

//acceptArgNameAndValue: reads an argument name and its value
func (p *CommandParser) acceptArg(lexeme rune) string {
	p.accept(lexeme, "*any identifier")
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
		p.addSyntaxError(fmt.Sprintf("wanted row/col number > 0; found %s", p.exactCurrentWord())) //keep parsing
		return MissingRwInt, nil
	}
	return RwInt(coordinate), nil
}

//parseCommaSepPoints reads a list of coordinates points in the form x,y,z
func (p *CommandParser) parseCommaSepPoints(ss *subspan) error {
	for {
		//read a coordinate point
		p.nextToken()
		point, err := p.parsePoint()
		if err != nil {
			return err
		}
		ss.list = append(ss.list, point)
		switch p.lexer.Peek() {
		case ',':
			p.nextToken()
			continue
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
	//p.nextToken()
	switch p.currentToken {
	case scanner.EOF:
		return nil
	case scanner.Ident:
		switch p.currentWord() {
		case "col", "row":
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
	default:
		return fmt.Errorf("expected row, col or an argument, found %s", p.exactCurrentWord())
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

//initialize lexer and its settings
func (p *CommandParser) init(r io.Reader) error {
	p.lexer = new(scanner.Scanner).Init(r)
	p.lexer.Whitespace = 1<<'\t' | 1<<'\r' | 1<<' ' //ignore spaces, tabs and CRs
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
		if !isCommandLine(line) {
			continue
		}
		p.init(strings.NewReader(line))
		cmdName, cmdToken := p.acceptCommandName()
		cmd = NewCommand(cmdName, cmdToken, scanner.Position{"", 0, i, 0}) //todo: fix line numbers
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

// //ParseCommands parses input stream and return an array of commands
// func (p *CommandParser) ParseCommands(r io.Reader) ([]*Command, error) {
// 	if err := p.init(r); err != nil {
// 		return nil, err
// 	}
// 	var cmd *Command
// 	var cmdList []*Command

// 	p.nextToken()
// 	if p.currentToken == scanner.EOF {
// 		return nil, NewError(ErrEmpty, scanner.Position{Line: -1, Column: -1}, "nothing to parse")
// 	}
// 	cmdList = make([]*Command, 0, sectionCapacity)

// 	for ; p.currentToken != scanner.EOF; p.nextToken() {
// 		if p.currentToken == '\n' { //handle lines with no text, just linefeeds
// 			p.accept('\n', "*end of line") //advance beyond it and loop back
// 			continue
// 		}
// 		cmdName, cmdToken := p.acceptCommandName()
// 		cmd = NewCommand(cmdName, cmdToken, p.lexer.Pos())
// 		switch cmdName {
// 		case "set":
// 			p.parseSetCommand(cmd)
// 		default:
// 			p.parseTableFormatCommand(cmd) //all other commands will be parsed as a formatting command
// 		}
// 		if p.currentToken == '\n' {
// 			p.accept('\n', "*end of line")
// 		}
// 		err := cmd.Validate()
// 		if err != nil {
// 			p.addSyntaxError(err.Error())
// 		}
// 		cmdList = append(cmdList, cmd)
// 	}

// 	if len(cmdList) == 0 {
// 		return nil, NewError(ErrEmpty, scanner.Position{Line: -1, Column: -1}, "found no valid commands")
// 	}
// 	if p.errors.Count() > 0 {
// 		return nil, NewError(ErrGeneric, scanner.Position{Line: -1, Column: -1}, "syntax errors")
// 	}
// 	return cmdList, nil
// }

// switch ss.kind {
// case "row":
// 	if cmd.info.sawRow {
// 		return nil, fmt.Errorf("row was already specified")
// 	}
// 	cmd.info.sawRow = true
// case "col":
// 	if cmd.info.sawCol {
// 		return nil, fmt.Errorf("col was already specified")
// 	}
// 	cmd.info.sawCol = true
// }
//now we must have a cell coordinate
//fmt.Println("just before parsePoint(),  p.currentToken=", p.currentWord())

//move to next token
//fmt.Println("just before nextNotNull(),  p.currentToken=", p.currentWord())
// if err := p.nextNotNull(); err != nil {
// 	return ss, nil //no more input, return
// }
//fmt.Println("just before switch,  p.currentToken=", p.currentWord())
