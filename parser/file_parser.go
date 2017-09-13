package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"text/scanner"

	"github.com/drgo/rosewood/types"
	"github.com/drgo/rosewood/utils"
)

var (
	unknownPos = scanner.Position{"", -1, -1, -1}
)

//File holds information on currently parsed Rosewood file
type File struct {
	FileName string
	sections []*types.Section //holds raw lines
	parser   *CommandParser
	settings *utils.Settings
}

//NewFile returns a File
func NewFile(fileName string, settings *utils.Settings) *File {
	return &File{FileName: fileName,
		parser:   NewCommandParser(settings),
		settings: settings}
}

//Parse parses an io.Reader streaming a Rosewood file and returns any found tables
func (f *File) Parse(r io.Reader) ([]*types.Table, error) {
	// helper function
	isSectionSeparatorLine := func(line string) bool {
		return strings.HasPrefix(strings.TrimSpace(line), f.settings.SectionSeparator)
	}
	var s *types.Section
	lineNum := 0
	scanner := bufio.NewScanner(r)
	//find a line that starts with a SectionSeparator
	for {
		more := scanner.Scan()
		if !more && scanner.Err() == nil { //EOF reached first
			return nil, NewError(ErrSyntaxError, unknownPos, "file is empty or does not start by a section separator: "+f.settings.SectionSeparator)
		}
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "//") {
			continue //skip comments and empty lines
		}
		if strings.HasPrefix(line, f.settings.SectionSeparator) {
			s = types.NewSection(types.SectionUnknown, lineNum+1) //found the first section
			break                                                 //SectionSeparator was in the first valid line
		}
		//other text found
		return nil, NewError(ErrSyntaxError, unknownPos, "file is empty or does not start by a section separator: "+f.settings.SectionSeparator)
	}
	//process the rest of the file
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		fmt.Println(line)
		if isSectionSeparatorLine(line) { //start of a new section
			if s != nil { //there is an active section, append it to the sections array
				f.sections = append(f.sections, s)
			}
			s = types.NewSection(types.SectionUnknown, lineNum+1) //create a new section
		} else {
			s.Lines = append(s.Lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, NewError(ErrSyntaxError, unknownPos, err.Error())
	}
	tables, err := f.createTables()
	if err != nil {
		return nil, NewError(ErrSyntaxError, unknownPos, err.Error())
	}
	return tables, nil
}

//SectionCount returns the number of sections found in the file
func (f *File) SectionCount() int {
	return len(f.sections)
}

func (f *File) createTables() (tables []*types.Table, err error) {
	if f.SectionCount() == 0 || f.SectionCount()%f.settings.SectionsPerTable != 0 {
		return nil, fmt.Errorf("incorrect number of sections: %d", f.SectionCount())
	}
	var t *types.Table
	for i, s := range f.sections {
		ii := i + 1 //i is zero-based, section numbers should be one-based
		kind := types.SectionDescriptor(i%f.settings.SectionsPerTable + 1)
		switch kind {
		case types.SectionCaption:
			t = types.NewTable()
			t.Caption = s
		case types.SectionBody:
			if t.Contents, err = types.NewTableContents(s.String()); err != nil {
				return nil, fmt.Errorf("error parsing table in section # %d: %s ", ii, err)
			}
		case types.SectionFootNotes:
			t.Footnotes = s
		case types.SectionControl:
			if t.CmdList, err = f.parser.ParseCommandLines(s); err != nil {
				return nil, err
			}
			tables = append(tables, t)
		default:
			panic(fmt.Sprintf("invalid switch case [%v] in RwFile.CreateTables()", kind))
		}
	}
	return tables, nil
}

//Errors returns a list of parsing errors
func (f *File) Errors() []error {
	return f.parser.Errors()
}
