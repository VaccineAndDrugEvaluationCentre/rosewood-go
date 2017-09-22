package parser

import (
	"fmt"
	"strconv"
	"text/scanner"

	"github.com/drgo/rosewood/types"
)

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
