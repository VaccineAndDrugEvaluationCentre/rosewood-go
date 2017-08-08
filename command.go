package carpenter

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
)

//RwInt used for all table cell coordinates
type RwInt uint

//Missing, min and max values for RwInt
const (
	MissingRwInt = ^RwInt(0)        //flip bits of zero to all 1s to get max uint for use as a sentinel for missing values
	MaxRwInt     = MissingRwInt - 1 //use this as the MaxUnit
	MinRwInt     = 0
)

type rwArgs []string

func (args rwArgs) String() string {
	if len(args) == 0 {
		return ""
	}
	return strings.Join(args, ",")
}

//Arg returns the index-th argument as unquoted string
func (args rwArgs) Arg(index int) string {
	if index < 0 || index >= len(args) {
		panic("rwArgs.UnquoteString called with invalid index")
	}
	if s, err := strconv.Unquote(args[index]); err == nil {
		return s
	}
	return ""
}

type rwCoordinate struct {
	Row, Col RwInt
}

func formattedRwInt(value RwInt) []byte { //return byte array for ease of concatenating with other text
	buf := []byte{}
	if value == MissingRwInt {
		buf = append(buf, 'N', 'A') //use NA for missing
	} else {
		buf = strconv.AppendUint(buf, uint64(value), 10)
	}
	return buf
}

func (co rwCoordinate) String() string {
	buf := formattedRwInt(co.Row)
	buf = append(buf, ':')
	buf = formattedRwInt(co.Col)
	return string(buf)
}

type rwRange struct {
	TopLeft     rwCoordinate
	BottomRight rwCoordinate
}

//NewCommand return an empty RwCommand
func newRange() rwRange {
	return rwRange{rwCoordinate{MinRwInt, MinRwInt}, rwCoordinate{MissingRwInt, MissingRwInt}} //assume topleft =(0,0)
}

func (r rwRange) String() string {
	return fmt.Sprintf("row %s col %s", r.TopLeft.String(), r.BottomRight.String())

}
func (r rwRange) testString() string {
	return fmt.Sprintf("row %s:%s col %s:%s", formattedRwInt(r.TopLeft.Row), formattedRwInt(r.BottomRight.Row),
		formattedRwInt(r.TopLeft.Col), formattedRwInt(r.BottomRight.Col))
}

func (r rwRange) Validate() error {
	if r.BottomRight.Row < r.TopLeft.Row /* && r.BottomRight.Row != -1 */ { //TopLeft.Row cannot be optional (-1)
		return fmt.Errorf("top row number (%d) must be smaller than bottom row number (%d)", r.TopLeft.Row, r.BottomRight.Row)
	}
	// if r.BottomRight.Col == -1 && r.TopLeft.Col == -1 { //if both missing, nothing to validate
	// 	return nil
	// }
	if r.BottomRight.Col < r.TopLeft.Col /* && r.BottomRight.Col != -1 */ { //one or both of them are not missing
		return fmt.Errorf("left column number (%d) must be smaller than right column number (%d)",
			r.TopLeft.Col, r.BottomRight.Col)
	}
	return nil
}

//RwCommand is the AST for a rw command.
type RwCommand struct {
	token     rwKeyWord
	name      string
	cellRange rwRange
	args      rwArgs
	pos       scanner.Position
	//	execOrder int
}

//NewCommand return an empty RwCommand
func NewCommand(name string, token rwKeyWord, pos scanner.Position) *RwCommand {
	return &RwCommand{token: token, name: name, pos: pos, cellRange: newRange()}
}

//formats command for printing
func (c *RwCommand) String() string {
	switch c.token {
	case kwSet:
		return fmt.Sprintf("%s %s", c.name, c.args)
	default:
		return fmt.Sprintf("%s %s %s", c.name, c.cellRange.testString(), c.args)
	}
}

//Validate checks command for errors
func (c *RwCommand) Validate() error {
	var err error
	switch c.token {
	case kwSet:
		if len(c.args) != 2 {
			return fmt.Errorf("expected 2 arguments, found %d arguments", len(c.args))
		}
	case kwMerge:
		return c.cellRange.Validate()
	case kwStyle:
		return err
	default:
	}
	return nil
}
