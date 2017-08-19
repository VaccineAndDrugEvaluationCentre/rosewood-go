package rosewood

import (
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

//Command is the AST for a rw command.
type Command struct {
	token     rwKeyWord
	name      string
	cellSpan  span
	cellRange Range
	args      Args
	pos       scanner.Position
	//	execOrder int
}

//NewCommand return an empty RwCommand
func NewCommand(name string, token rwKeyWord, pos scanner.Position) *Command {
	return &Command{token: token, name: name, pos: pos, cellRange: newRange(), cellSpan: newSpan()}
}

//formats command for printing
func (c *Command) String() string {
	switch c.token {
	case kwSet:
		return fmt.Sprintf("%s %s", c.name, c.args)
	default:
		return fmt.Sprintf("%s %s %s", c.name, c.cellRange.testString(), c.args)
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
		return c.cellRange.validate()
	case kwStyle:
		return err
	default:
	}
	return nil
}
