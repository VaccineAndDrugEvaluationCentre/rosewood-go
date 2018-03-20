package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"text/scanner"

	"github.com/drgo/rosewood/settings"
	"github.com/drgo/rosewood/types"
)

//Position is an alias of scanner.Position
type Position = scanner.Position

var (
	unknownPos = Position{"", -1, -1, -1}
)

const (
	//SectionsPerTable the number of section per table
	SectionsPerTable = 4
	//RwMinFileSize the size of the smallest possible rosewood file
	RwMinFileSize = SectionsPerTable * 3 // at least 3 separator chars per section
)

//File holds information on currently parsed Rosewood file
type File struct {
	FileName string
	sections []*types.Section //holds raw lines
	parser   *CommandParser
	settings *settings.Settings
	tables   []*types.Table //holds parsed tables and commands
}

//NewFile returns a Rosewood File
func NewFile(fileName string, settings *settings.Settings) *File {
	return &File{FileName: fileName,
		parser:   NewCommandParser(settings),
		settings: settings}
}

//Parse parses an io.Reader streaming a Rosewood file and returns any found tables
func (f *File) Parse(r io.Reader) error {
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
			return NewError(ErrSyntaxError, unknownPos, "file is empty or does not start by a section separator: "+f.settings.SectionSeparator)
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
		return NewError(ErrSyntaxError, unknownPos, "file is empty or does not start by a section separator: "+f.settings.SectionSeparator)
	}
	//process the rest of the file
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		//fmt.Println(line)
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
		return NewError(ErrSyntaxError, unknownPos, err.Error())
	}
	return f.createTables()
}

//SectionCount returns the number of sections found in the file
func (f *File) SectionCount() int {
	return len(f.sections)
}

func (f *File) createTables() error {
	if f.SectionCount() == 0 || f.SectionCount()%f.settings.SectionsPerTable != 0 {
		return fmt.Errorf("incorrect number of sections: %d", f.SectionCount())
	}
	var t *types.Table
	var err error
	for i, s := range f.sections {
		ii := i + 1 //i is zero-based, section numbers should be one-based
		kind := types.SectionDescriptor(i%f.settings.SectionsPerTable + 1)
		// fmt.Printf("kind: %d\n", kind)
		// fmt.Printf("section: %s\n", s)
		switch kind {
		case types.SectionCaption:
			t = types.NewTable()
			t.Caption = s
		case types.SectionBody:
			if t.Contents, err = types.NewTableContents(s.String()); err != nil {
				return NewError(ErrSyntaxError, unknownPos, fmt.Sprintf("error parsing table in section #%d: %s ", ii, err))
			}
		case types.SectionFootNotes:
			t.Footnotes = s
		case types.SectionControl:
			if t.CmdList, err = f.parser.ParseCommandLines(s); err != nil {
				return err
			}
			f.tables = append(f.tables, t)
		default:
			panic(fmt.Sprintf("invalid switch case [%v] in File.CreateTables()", kind))
		}
	}
	return nil
}

//TableCount returns the number of prased tables in the file
func (f *File) TableCount() int {
	return len(f.tables)
}

//Tables returns an array of pointers to parsed Tables
func (f *File) Tables() []*types.Table {
	return f.tables
}

//Errors returns a list of parsing errors
func (f *File) Errors() []error {
	return f.parser.Errors()
}

//Err returns a list of parsing errors
func (f *File) Err() error {
	return f.parser.errors.Err()
}
