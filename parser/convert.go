package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/drgo/core/str"
	"github.com/drgo/core/ui"
	"github.com/drgo/rosewood/types"
)

//RWSyntax holds info on syntax differences in each Rosewood version
type RWSyntax struct {
	Version          string
	SectionSeparator string
	ColumnSeparator  string
}

//RWSyntaxVdotzero1 syntax for Rosewood v0.1
var RWSyntaxVdotzero1 = RWSyntax{
	Version:          "v0.1",
	SectionSeparator: "---",
	ColumnSeparator:  "|",
}

//RWSyntaxVdotzero2 syntax for Rosewood v0.2
var RWSyntaxVdotzero2 = RWSyntax{
	Version:          "v0.2",
	SectionSeparator: "+++",
	ColumnSeparator:  "|",
}

//ConvertToCurrentVersion converts to current Rosewood version and saves results to io.Writer
func ConvertToCurrentVersion(settings *types.RosewoodSettings, oldSyntax RWSyntax, in io.Reader, out io.Writer) error {
	newCode, err := ConvertVersion(settings, RWSyntaxVdotzero2, oldSyntax, in)
	if err != nil {
		return err
	}
	//write all modified code to output writer using a buffer
	w := bufio.NewWriter(out)
	for _, line := range newCode {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

//ConvertVersion convert between two Rosewood version
func ConvertVersion(settings *types.RosewoodSettings, newSyntax RWSyntax, oldSyntax RWSyntax, in io.Reader) ([]string, error) {
	newCode := make([]string, 0, 256)
	lineNum, sectionNum := 0, 0
	headerStart, headerEnd, tableStart, tableEnd, footnoteStart, rulesStart := -1, -1, -1, -1, -1, -1
	// utility for debugging
	output := func(line string) {
		newCode = append(newCode, line)
		if settings.Debug == ui.DebugAll {
			fmt.Printf("%d:%s\n", lineNum, line)
		}
	}
	scanner := bufio.NewScanner(in)
	//process the first line
	if !scanner.Scan() {
		return nil, NewError(ErrSyntaxError, unknownPos, scanner.Err().Error())
	}
	lineNum++ //we found a line
	if GetFileVersion(strings.TrimSpace(scanner.Text())) != "v0.1" {
		return nil, fmt.Errorf("not a Rosewood v0.1 file")
	}
	sectionNum++ //we found the first section separator
	output(newSyntax.SectionSeparator)

	//now process the rest
	for scanner.Scan() {
		lineNum++
		orgLine := scanner.Text()
		if strings.TrimSpace(orgLine) == "" { //skip empty lines
			continue
		}
		firstToken := strings.Fields(orgLine)[0] //can't panic because of the TrimSpace check above
		switch firstToken {
		case oldSyntax.SectionSeparator:
			sectionNum++
			output(newSyntax.SectionSeparator)
		case "merge", "plain", "indent":
			if rulesStart == -1 {
				rulesStart = lineNum
			}
			output(fixCommandRule(orgLine, firstToken))
		default:
			if strings.Contains(orgLine, oldSyntax.ColumnSeparator) { //found a table row
				if tableStart == -1 {
					tableStart = lineNum //table section started
					tableEnd = lineNum
				} else {
					tableEnd++
				}
			} else { //non-table text after a table has started; must be a footnote
				if footnoteStart == -1 && tableEnd > 0 {
					footnoteStart = lineNum
				}
			}
			output(orgLine)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, NewError(ErrSyntaxError, unknownPos, err.Error())
	}
	//validate the table structure
	switch {
	case sectionNum > 6:
		return nil, fmt.Errorf("too many section separators;multiple tables? not supported")
	case len(newCode) == 0:
		return nil, fmt.Errorf("empty file")
	case tableStart == -1:
		return nil, fmt.Errorf("file does contain a valid table section")
	default:
	}

	//find a header so that we can convert it into a style header command later on
	for i := tableStart; i <= tableEnd; i++ {
		if strings.HasPrefix(newCode[i], newSyntax.SectionSeparator) {
			headerEnd = i
			headerStart = tableStart //header section always starts a table
			tableEnd++               //add one for this separator line which was not counted in the loop above
			break
		}
	}
	//if a footnote section is missing, create an empty one
	if footnoteStart == -1 {
		newCode = str.InsertToStringSlice(newCode, tableEnd+1, newSyntax.SectionSeparator)
	}

	//add comment to the beginning of the generated code
	if rulesStart == -1 { //no rulese section,
		rulesStart = len(newCode) // it should start where the last section separator is located in the output
	}
	newCode = str.InsertToStringSlice(newCode, rulesStart-1, fmt.Sprintf("//Automatically converted by Carpenter from version 0.1 on %s", time.Now().Format("2006-01-02 15:04:05")))

	//TODO: change "header" to strconv.Quote("header") to produce a quoted string for css class name
	if headerStart > -1 {
		newCode = str.InsertToStringSlice(newCode, rulesStart, fmt.Sprintf("style row %d:%d %s", 1, headerEnd-headerStart+1, "header"))
		//remove the header start section separator //TODO: refactor as RemovefromSlice
		newCode = append(newCode[:headerEnd], newCode[headerEnd+1:]...)
	}
	if settings.Debug == ui.DebugAll {
		fmt.Printf("File had %d section separators, table starts on line %d and ends on line %d, header starts on line %d and ends on line %d, footnotes section starts on line %d, rules section starts on line %d\n", sectionNum, tableStart, tableEnd, headerStart, headerEnd, footnoteStart, rulesStart)
	}
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
