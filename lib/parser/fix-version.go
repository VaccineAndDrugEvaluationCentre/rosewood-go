package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/drgo/rosewood/lib/settings"
)

type RWSyntax struct {
	Version          string
	SectionSeparator string
}

var RWSyntaxV1 = RWSyntax{
	Version:          "v0.1",
	SectionSeparator: "---",
}

//ConvertV1toV2 utility to convert simple Rosewood v1 files to valid v2 files
func ConvertToCurrentVersion(settings *settings.Settings, syntax RWSyntax, in io.Reader) ([]string, error) {
	newCode := make([]string, 0, 256)
	output := func(line string) {
		newCode = append(newCode, line)
	}
	lineNum, sectionNum, headerStart, headerEnd, rulesStart := 0, 0, 0, 0, 0
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		lineNum++
		orgLine := scanner.Text()
		if strings.TrimSpace(orgLine) == "" {
			output(orgLine)
			continue
		}
		firstToken := strings.Fields(orgLine)[0] //can't panic because of the TrimSpace check
		switch firstToken {
		case syntax.SectionSeparator:
			sectionNum++
			if sectionNum == 3 && headerStart > 0 { //skip the section separator between the table header and table body
				continue
			}
			output(settings.SectionSeparator)
		case "merge", "plain", "indent":
			if rulesStart == 0 {
				rulesStart = lineNum - 1 //we deleted one +++
			}
			output(fixCommandRule(orgLine, firstToken))
		default:
			if strings.Contains(orgLine, settings.ColumnSeparator) { //found a table row
				if headerStart == 0 {
					headerStart = 1 //header section started
					headerEnd = 1
				} else {
					if sectionNum < 3 && headerStart > 0 {
						headerEnd++ //increase headerEnd until we hit the third section separator
					}
				}
			}
			output(orgLine)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, NewError(ErrSyntaxError, unknownPos, err.Error())
	}
	//TODO: make the comment more informative eg date etc
	//add comment to the generated code
	newCode = InsertToStringSlice(newCode, rulesStart-1, fmt.Sprintf(
		`//Automatically converted by Carpenter from version 0.1`))
	//now add style header command
	//TODO: change to strconv.Quote("header") to produce a quoted string for css class name
	output(fmt.Sprintf("style row %d:%d %s", headerStart, headerEnd, "header"))
	//move it before the last element which is the final section separator
	len := len(newCode)
	newCode[len-2], newCode[len-1] = newCode[len-1], newCode[len-2]
	return newCode, nil
}

func fixCommandRule(line, firstToken string) string {
	switch firstToken {
	case "plain", "indent":
		line = strings.Replace(line, firstToken, "style", 1) //replace the old command eg indent etc with style
		//TODO: change to %q to produce a quoted string for css class name
		line += fmt.Sprintf(" %s", firstToken) //add the old command as an argument to style
	}
	line = strings.Replace(line, "column", "col", 1)
	return line
}

//TODO: move to stringUtils
//from https://blog.golang.org/slices
// Insert inserts the value into the slice at the specified index,
// which must be in range.
// The slice must have room for the new element.
func InsertToStringSlice(slice []string, index int, value string) []string {
	// Grow the slice by one element.
	slice = slice[0 : len(slice)+1]
	// Use copy to move the upper part of the slice out of the way and open a hole.
	copy(slice[index+1:], slice[index:])
	// Store the new value.
	slice[index] = value
	// Return the result.
	return slice
}
