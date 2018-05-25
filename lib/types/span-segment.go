package types

import (
	"bytes"
	"fmt"
)

//SpanSegment holds info on the span of cells that a command will apply to
//as specified in a Rosewood command. Should be translated to one or more spans
//in the physical table
type SpanSegment struct {
	kind        string  // row or col
	Left, Right RwInt   // e.g., row 1:2
	By          RwInt   // holds the step in eg row 1:2:6
	List        []RwInt // holds list of row/col numbers eg row 1,2,3,4
}

//NewSpanSegment returns a segment of certain kind
func NewSpanSegment(kind string) SpanSegment {
	return SpanSegment{kind: kind, Left: MissingRwInt, Right: MissingRwInt, By: MissingRwInt}
}

func (ss *SpanSegment) String() string {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "%s ", ss.kind)
	if ss.Left != MissingRwInt {
		fmt.Fprintf(buf, "%s", formattedRwInt(ss.Left))
		if ss.By != MissingRwInt {
			fmt.Fprintf(buf, ":%s", formattedRwInt(ss.By))
		}
		fmt.Fprintf(buf, ":%s", formattedRwInt(ss.Right))
		if len(ss.List) > 0 {
			fmt.Fprintf(buf, ", ") //add comma if we also have a comma separated list
		}
	}
	for _, item := range ss.List {
		fmt.Fprintf(buf, "%s,", formattedRwInt(item))
	}
	//remove last comma if any
	if bytes.HasSuffix(buf.Bytes(), []byte{','}) {
		buf.Truncate(buf.Len() - 1)
	}
	return buf.String()
}
