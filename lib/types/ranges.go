// Copyright 2017 Salah Mahmud and Colleagues. All rights reserved.

package types

import (
	"fmt"
	"strconv"
)

//RwInt used for all table cell coordinates
type RwInt uint

//Missing, min and max values for RwInt
const (
	MissingRwInt = ^RwInt(0)        //flip bits of zero to all 1s to get max uint for use as a sentinel for missing values
	MaxRwInt     = MissingRwInt - 1 //use this as the MaxUnit
	MinRwInt     = 1
)

//Coordinates holds the row, col of a table cell
type Coordinates struct {
	Row, Col RwInt
}

func formattedRwInt(value RwInt) []byte { //return byte array for ease of concatenating with other text
	var buf []byte
	if value == MissingRwInt {
		buf = append(buf, 'N', 'A') //use NA for missing
	} else {
		buf = strconv.AppendUint(buf, uint64(value), 10)
	}
	return buf
}

func (co Coordinates) String() string {
	buf := formattedRwInt(co.Row)
	buf = append(buf, ':')
	buf = append(buf, formattedRwInt(co.Col)...)
	return string(buf)
}

type Range struct {
	TopLeft     Coordinates
	BottomRight Coordinates
	styleList   []string
}

//newRange return an empty a range
func newRange() Range {
	return Range{Coordinates{MinRwInt, MinRwInt}, Coordinates{MissingRwInt, MissingRwInt}, nil} //assume topleft =(1,1)
}

func makeRange(tlr, tlc, brr, brc RwInt) Range {
	return Range{Coordinates{tlr, tlc}, Coordinates{brr, brc}, nil}
}

func (r Range) String() string {
	return fmt.Sprintf("topleft %s bottomright %s", r.TopLeft, r.BottomRight)

}
func (r Range) testString() string {
	return fmt.Sprintf("row %s:%s col %s:%s", formattedRwInt(r.TopLeft.Row), formattedRwInt(r.BottomRight.Row),
		formattedRwInt(r.TopLeft.Col), formattedRwInt(r.BottomRight.Col))
}

func (r *Range) less(s Range) bool {
	return r.TopLeft.Row < s.TopLeft.Row || (r.TopLeft.Row == s.TopLeft.Row && r.TopLeft.Col < s.TopLeft.Col)
}

func (r *Range) styles() []string {
	return r.styleList
}

// Normalize ... place a min or max on the range, based on the row / col count of the intended table
func (r *Range) Normalize() {

	// cap rows
	if r.TopLeft.Row == MissingRwInt {
		r.TopLeft.Row = 1
	}
	if r.BottomRight.Row == MissingRwInt {
		r.BottomRight.Row = r.TopLeft.Row
	}

	// cap columns
	if r.TopLeft.Col == MissingRwInt {
		r.TopLeft.Col = 1
	}
	if r.BottomRight.Col == MissingRwInt {
		r.BottomRight.Col = r.TopLeft.Col
	}
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
