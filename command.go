package carpenter

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
)

// Constant definitions
const (
	MissingUint = ^uint(0)        //flip zero to get all 1s to get max uint and use as a sentinel for missing uint
	MaxUint     = MissingUint - 1 //use this as the MaxUnit
	MinUint     = 0
)

type rwArgs []string

func (args rwArgs) String() string {
	if len(args) == 0 {
		return ""
	}
	return strings.Join(args, ",")
}

type rwCoordinate struct {
	Row, Col uint
}

func formattedUint(value uint) []byte { //return byte array for ease of concatenating with other text
	buf := []byte{}
	if value == MissingUint {
		buf = append(buf, 'N', 'A') //use NA for missing
	} else {
		buf = strconv.AppendUint(buf, uint64(value), 10)
	}
	return buf
}

func (co rwCoordinate) String() string {
	buf := formattedUint(co.Row)
	buf = append(buf, ':')
	buf = formattedUint(co.Col)
	return string(buf)
}

type rwRange struct {
	TopLeft     rwCoordinate
	BottomRight rwCoordinate
}

//NewCommand return an empty RwCommand
func newRange() rwRange {
	return rwRange{rwCoordinate{MinUint, MinUint}, rwCoordinate{MissingUint, MissingUint}} //assume topleft =(0,0)
}

func (r rwRange) String() string {
	return fmt.Sprintf("row %s col %s", r.TopLeft.String(), r.BottomRight.String())

}
func (r rwRange) AsEnteredString() string {
	return fmt.Sprintf("row %s:%s col %s:%s", formattedUint(r.TopLeft.Row), formattedUint(r.BottomRight.Row),
		formattedUint(r.TopLeft.Col), formattedUint(r.BottomRight.Col))
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
		return fmt.Sprintf("%s %s %s", c.name, c.cellRange.AsEnteredString(), c.args)
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
		err = c.cellRange.Validate()
		fallthrough
	case kwStyle:
		return err
	default:
	}
	return nil
}
