// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package parser

import (
	"fmt"
	"strconv"
	"text/scanner"

	"github.com/drgo/rosewood/lib/types"
)

//parseTableFormatCommand: parses a command like command row col args
func (p *CommandParser) parseTableFormatCommand(cmd *types.Command) error {
	//helper closure
	parseSegment := func(modifier string) error {
		ss, err := p.parseRowOrColSegment(cmd, modifier)
		if err == nil {
			cmd.AddSpanSegment(&ss)
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
			if cmd.SpanSegment(p.currentWord()) != nil { //already parsed
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

//TODO: convert to return a pointer to SpanSegment, so we could return nil on error
func (p *CommandParser) parseRowOrColSegment(cmd *types.Command, modifier string) (types.SpanSegment, error) {
	var err error
	ss := types.NewSpanSegment(modifier)

	if ss.Left, err = p.parsePoint(); err != nil {
		return ss, err
	}
	if ss.Left == types.MaxRwInt {
		return ss, fmt.Errorf("max is not allowed in this position")
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

//parseRangePoints read a range of coordinate either left:right or left:skip step:right
func (p *CommandParser) parseRangePoints(ss *types.SpanSegment) error {
	var err error
	//read the right term of the range
	if ss.Right, err = p.parsePoint(); err != nil {
		return err
	}
	p.nextToken()
	switch p.currentToken {
	case p.settings.RangeOperator: //another :, so this a skipped range l:step:r
		if ss.Right == types.MaxRwInt {
			return fmt.Errorf("max is not allowed in this position")
		}
		ss.By = ss.Right //what we thought was the right coordinate is actually a step
		//read the right term of the range
		if ss.Right, err = p.parsePoint(); err != nil {
			return err
		}
		p.nextToken()
		if p.currentToken == ',' {
			if ss.Right == types.MaxRwInt { //comma not allowed after a max
				return fmt.Errorf("max is not allowed in this position")
			}
			if err := p.parseCommaSepPoints(ss); err != nil {
				return err
			}
		}
	case ',':
		if ss.Right == types.MaxRwInt { //comma not allowed after a max
			return fmt.Errorf("max is not allowed in this position")
		}
		if err := p.parseCommaSepPoints(ss); err != nil {
			return err
		}
	default:
		return nil
	}
	return nil
}

//parseCommaSepPoints reads a list of coordinates points in the form x,y,z
func (p *CommandParser) parseCommaSepPoints(ss *types.SpanSegment) error {
	for {
		//parser is on a comma, read a coordinate point
		point, err := p.parsePoint()
		if err != nil {
			return err
		}
		if point == types.MaxRwInt {
			return fmt.Errorf("max is not allowed in this position")
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
func (p *CommandParser) parsePoint(signed bool) (int, error) {
	if err := p.nextNotNull(); err != nil {
		return types.MissingRwInt, err
	}
	if p.currentWord() == "max" {
		return types.MaxRwInt, nil

	}
	if p.currentToken != scanner.Int {
		return types.MissingRwInt, fmt.Errorf("expected col or row number, found %s", p.exactCurrentWord())
	}
	coordinate, _ := strconv.Atoi(p.currentWord()) //no error check as we know it must be an int
	if !signed && coordinate < 1 {
		p.addSyntaxError("wanted row/col number > 0; found %s", p.exactCurrentWord()) //keep parsing
		return types.MissingRwInt, nil
	}
	return types.RwInt(coordinate), nil
}
