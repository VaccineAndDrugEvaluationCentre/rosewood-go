package types

import (
	"bytes"
	"fmt"
)

//SpanSegment holds info on the span of cells that a command will apply to
//as specified in a Rosewood command. Should be translated to one or more spans
//in the physical table
type SpanSegment struct {
	kind        string // row or col
	Left, Right int    // e.g., row 1:2
	By          int    // holds the step in eg row 1:2:6
	List        []int  // holds list of row/col numbers eg row 1,2,3,4
}

//NewSpanSegment returns a segment of certain kind
func NewSpanSegment(kind string) SpanSegment {
	return SpanSegment{kind: kind, Left: RwMissing, Right: RwMissing, By: RwMissing}
}

func (ss *SpanSegment) String() string {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "%s ", ss.kind)
	if ss.Left != RwMissing {
		fmt.Fprintf(buf, "%s", formattedCellCoord(ss.Left))
		if ss.By != RwMissing {
			fmt.Fprintf(buf, ":%s", formattedCellCoord(ss.By))
		}
		fmt.Fprintf(buf, ":%s", formattedCellCoord(ss.Right))
		if len(ss.List) > 0 {
			fmt.Fprintf(buf, ", ") //add comma if we also have a comma separated list
		}
	}
	for _, item := range ss.List {
		fmt.Fprintf(buf, "%s,", formattedCellCoord(item))
	}
	//remove last comma if any
	if bytes.HasSuffix(buf.Bytes(), []byte{','}) {
		buf.Truncate(buf.Len() - 1)
	}
	return buf.String()
}
