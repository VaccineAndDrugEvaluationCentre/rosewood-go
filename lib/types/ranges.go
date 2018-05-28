// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package types

import (
	"fmt"
	"strconv"
)

//int used for all table cell coordinates
//type int uint

//Missing, min and max values for int
const (
	RwMissing = 1<<63 - 1     //get the max int64 for use as a sentinel for missing values
	RwMax     = RwMissing - 1 //use this as the Maxint
	RwMin     = 1
)

//Coordinates holds the row, col of a table cell
type Coordinates struct {
	Row, Col int
}

//TODO: optimize
func formattedCellCoord(value int) []byte { //return byte array for ease of concatenating with other text
	var buf []byte
	switch value {
	case RwMissing:
		buf = append(buf, 'N', 'A') //use NA for missing
	case RwMax:
		buf = append(buf, 'm', 'a', 'x') //use max for missing
	default:
		buf = strconv.AppendInt(buf, int64(value), 10)
	}
	return buf
}

//TODO: optimize
func (co Coordinates) String() string {
	buf := formattedCellCoord(co.Row)
	buf = append(buf, ':')
	buf = append(buf, formattedCellCoord(co.Col)...)
	return string(buf)
}

type Range struct {
	TopLeft     Coordinates
	BottomRight Coordinates
	styleList   []string
}

//newRange return an empty a range
func newRange() Range {
	return Range{Coordinates{RwMin, RwMin}, Coordinates{RwMissing, RwMissing}, nil} //assume topleft =(1,1)
}

func makeRange(tlr, tlc, brr, brc int) Range {
	return Range{Coordinates{tlr, tlc}, Coordinates{brr, brc}, nil}
}

func (r Range) String() string {
	return fmt.Sprintf("topleft %s bottomright %s", r.TopLeft, r.BottomRight)

}
func (r Range) testString() string {
	return fmt.Sprintf("row %s:%s col %s:%s", formattedCellCoord(r.TopLeft.Row), formattedCellCoord(r.BottomRight.Row),
		formattedCellCoord(r.TopLeft.Col), formattedCellCoord(r.BottomRight.Col))
}

func (r *Range) less(s Range) bool {
	return r.TopLeft.Row < s.TopLeft.Row || (r.TopLeft.Row == s.TopLeft.Row && r.TopLeft.Col < s.TopLeft.Col)
}

func (r *Range) styles() []string {
	return r.styleList
}

//AddStyle adds one or more style names if they do not already exist in the list
//sufficiently efficient for short lists and avoids allocating a map
func (r *Range) addStyle(styles ...string) error {
outer:
	for _, s := range styles {
		for _, ss := range r.styleList { //skip s if it already exists in the list
			if ss == s {
				continue outer
			}
		}
		r.styleList = append(r.styleList, s)
	}
	return nil
}

//validate performs simple validation of the range coordinates
func (r Range) validate() error {
	if r.BottomRight.Row < r.TopLeft.Row /* && r.BottomRight.Row != -1 */ { //TopLeft.Row cannot be optional (-1)
		return fmt.Errorf("top row number (%d) must be smaller than bottom row number (%d)", r.TopLeft.Row, r.BottomRight.Row)
	}
	if r.BottomRight.Col < r.TopLeft.Col /* && r.BottomRight.Col != -1 */ { //one or both of them are not missing
		return fmt.Errorf("left column number (%d) must be smaller than right column number (%d)",
			r.TopLeft.Col, r.BottomRight.Col)
	}
	return nil
}
