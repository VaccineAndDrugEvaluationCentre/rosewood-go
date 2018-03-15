//Package types implements data structures essential for parsing and rendering Rosewood tables and commands
package types

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type rwArgs []string

func (args rwArgs) String() string {
	return strings.Join(args, ",")
}

//Command is the AST for a Rosewood command.
type Command struct {
	token    RwKeyWord
	name     string
	cellSpan *Span
	spans    []*Subspan
	args     rwArgs
}

//NewCommand return an empty RwCommand
func NewCommand(name string, token RwKeyWord) *Command {
	return &Command{token: token, name: name}
}

//formats command for printing; warning used for testing the parser
func (c *Command) String() string {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "%s ", c.name)
	for _, s := range c.spans {
		fmt.Fprintf(buf, s.String())
		fmt.Fprintf(buf, " ")
	}
	if len(c.args) > 0 {
		fmt.Fprintf(buf, "%s", c.args)
	}
	return strings.TrimSpace(buf.String())
}

func (c *Command) AddSubSpan(ss *Subspan) error {
	c.spans = append(c.spans, ss)
	return nil
}

func (c *Command) AddArg(arg ...string) error {
	c.args = append(c.args, arg...)
	return nil
}

func (c *Command) SubSpan(modifier string) *Subspan {
	for _, ss := range c.spans {
		if ss.kind == modifier {
			return ss
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
func (c *Command) Finalize() error {
	checkCmd := func() error {
		c.cellSpan = SubSpansToSpan(c.spans)
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
