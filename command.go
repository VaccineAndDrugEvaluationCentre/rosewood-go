package rosewood

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
)

type rwArgs []string

func (args rwArgs) String() string {
	return strings.Join(args, ",")
}

//Arg returns the index-th argument as unquoted string
func (args rwArgs) Arg(index int) string {
	if index < 0 || index >= len(args) {
		panic("Args.arg() called with invalid index")
	}
	if s, err := strconv.Unquote(args[index]); err == nil {
		return s
	}
	return ""
}

//Command is the AST for a Rosewood command.
type Command struct {
	token    rwKeyWord
	name     string
	cellSpan *span
	spans    []*subspan
	args     rwArgs
	pos      scanner.Position
}

//NewCommand return an empty RwCommand
func NewCommand(name string, token rwKeyWord, pos scanner.Position) *Command {
	return &Command{token: token, name: name, pos: pos, cellSpan: newSpan()}
}

//formats command for printing; warning used for testing the parser
func (c *Command) String() string {
	switch c.token {
	case kwSet:
		return fmt.Sprintf("%s %s", c.name, c.args)
	default:
		buf := &bytes.Buffer{}
		fmt.Fprintf(buf, "%s ", c.name)
		for _, s := range c.spans {
			fmt.Fprintf(buf, "%s ", s.kind)
			if s.left != MissingRwInt {
				fmt.Fprintf(buf, "%s", formattedRwInt(s.left))
				if s.by != MissingRwInt {
					fmt.Fprintf(buf, ":%s", formattedRwInt(s.by))
				}
				fmt.Fprintf(buf, ":%s", formattedRwInt(s.right))
			}
			for _, item := range s.list {
				fmt.Fprintf(buf, "%s,", formattedRwInt(item))
			}
			//remove last comma if any
			if bytes.HasSuffix(buf.Bytes(), []byte{','}) {
				buf.Truncate(buf.Len() - 1)
			}
			fmt.Fprintf(buf, " ")
		}
		if len(c.args) > 0 {
			fmt.Fprintf(buf, "%s", c.args)
		}
		return strings.TrimSpace(buf.String())
	}
}

func (c *Command) subSpan(modifier string) *subspan {
	for _, ss := range c.spans {
		if ss.kind == modifier {
			return ss
		}
	}
	return nil
}

func (c *Command) validateMerge() error {

	return nil
}

//Validate checks command for errors
func (c *Command) Validate() error {
	var err error
	switch c.token {
	case kwSet:
		if len(c.args) != 2 {
			return fmt.Errorf("expected 2 arguments, found %d arguments", len(c.args))
		}
	case kwMerge:
		return c.validateMerge()
	case kwStyle:
		return err
	default:
	}
	return nil
}
