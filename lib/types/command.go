//Package types implements data structures essential for parsing and rendering Rosewood tables and commands
// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package types

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

//rwArgs holds command arguments
type rwArgs []string

func (args rwArgs) String() string {
	return strings.Join(args, ",")
}

//Command is the AST for a Rosewood command.
type Command struct {
	token        RwKeyWord      //command token
	name         string         //command name
	cellSpan     *Span          //holds the complete valid description of the span that the command applies to
	spanSegments []*SpanSegment //row and/or col table spanSegments that the command applies to.
	args         rwArgs         //additional arguments passed to the command
}

//NewCommand return an empty RwCommand
func NewCommand(name string, token RwKeyWord) *Command {
	return &Command{token: token, name: name}
}

//formats command for printing; Warning: used for testing the parser, changing it might break some tests
func (c *Command) String() string {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "%s ", c.name)
	for _, s := range c.spanSegments {
		fmt.Fprintf(buf, s.String())
		fmt.Fprintf(buf, " ")
	}
	if len(c.args) > 0 {
		fmt.Fprintf(buf, "%s", c.args)
	}
	return strings.TrimSpace(buf.String())
}

//AddSpanSegment adds a SpanSegment to the command
func (c *Command) AddSpanSegment(segment *SpanSegment) error {
	//TODO: ensure only valid SpanSegments are added and no duplicates
	c.spanSegments = append(c.spanSegments, segment)
	return nil
}

//AddArg adds one or more arguments to the command
func (c *Command) AddArg(arg ...string) error {
	c.args = append(c.args, arg...)
	return nil
}

//SpanSegment returns a SpanSegment corresponding to the specified kind: row or col
func (c *Command) SpanSegment(kind string) *SpanSegment {
	for _, segment := range c.spanSegments {
		if segment.kind == kind {
			return segment
		}
	}
	return nil
}

//Args returns a list of all arguments
func (c *Command) Args() []string {
	return c.args
}

//Arg returns the index-th argument as unquoted string
func (c *Command) Arg(index int) string {
	if index < 0 || index >= len(c.args) {
		panic("Args.arg() called with invalid index")
	}
	if s, err := strconv.Unquote(c.args[index]); err == nil {
		return s
	}
	return ""
}

//Finalize creates a cell span and checks command for errors
//TODO:
func (c *Command) Finalize() error {
	checkCmd := func() error {
		c.cellSpan = NewSpanFromSpanSegments(c.spanSegments)
		if err := c.cellSpan.Validate(); err != nil {
			return err
		}
		return nil
	}
	switch c.token {
	case KwMerge:
		return checkCmd()
	case KwStyle:
		return checkCmd()
	case KwSet:
		if len(c.args) != 2 {
			return fmt.Errorf("expected 2 arguments, found %d arguments", len(c.args))
		}
	default:
		panic(fmt.Sprintf("wrong token %d in command.finalize()", c.token))
	}
	return nil
}
