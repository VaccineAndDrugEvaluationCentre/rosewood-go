package rosewood

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
)

type Args []string

func (args Args) String() string {
	if len(args) == 0 {
		return ""
	}
	return strings.Join(args, ",")
}

//Arg returns the index-th argument as unquoted string
func (args Args) Arg(index int) string {
	if index < 0 || index >= len(args) {
		panic("rwArgs.UnquoteString called with invalid index")
	}
	if s, err := strconv.Unquote(args[index]); err == nil {
		return s
	}
	return ""
}

//TODO: use scanner.Position in error messages
type commandInfo struct {
	scanner.Position
	sawRow, sawCol bool
}

//Command is the AST for a rw command.
type Command struct {
	token    rwKeyWord
	name     string
	cellSpan span
	spans    []*subspan
	args     Args
	pos      scanner.Position
	info     commandInfo
}

//NewCommand return an empty RwCommand
func NewCommand(name string, token rwKeyWord, pos scanner.Position) *Command {
	return &Command{token: token, name: name, pos: pos, cellSpan: newSpan()}
}

//formats command for printing
func (c *Command) String() string {
	switch c.token {
	case kwSet:
		return fmt.Sprintf("%s %s", c.name, c.args)
	default:
		buf := &bytes.Buffer{}
		fmt.Fprintf(buf, "%s", c.name)
		for _, s := range c.spans {
			fmt.Fprintf(buf, " %s %s", s.kind, formattedRwInt(s.left))
			if s.by != MissingRwInt {
				fmt.Fprintf(buf, ":%s", formattedRwInt(s.by))
			}
			fmt.Fprintf(buf, ":%s", formattedRwInt(s.right))
		}
		if len(c.args) > 0 {
			fmt.Fprintf(buf, " %s", c.args)
		}
		return buf.String()
	}
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
		return c.cellSpan.validate()
	case kwStyle:
		return err
	default:
	}
	return nil
}
