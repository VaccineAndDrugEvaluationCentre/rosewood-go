package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/drgo/rosewood/types"
	"github.com/drgo/rosewood/utils"
)

type RwFile struct {
	FileName string
	sections []*types.Section //holds raw lines
	parser   *CommandParser
	settings *utils.Settings
}

func NewRwFile(fileName string, settings *utils.Settings) *RwFile {
	return &RwFile{FileName: fileName,
		parser:   NewCommandParser(settings),
		settings: settings}
}

func (rw *RwFile) Parse(r io.Reader) ([]*types.Table, error) {
	// helper functions
	isSectionSeparatorLine := func(line string) bool {
		return strings.HasPrefix(strings.TrimSpace(line), rw.settings.SectionSeparator)
	}
	var s *types.Section
	lineNum := 0
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if isSectionSeparatorLine(line) { //start a new section
			if s != nil { //there is an active section, append it to the sections array
				rw.sections = append(rw.sections, s)
			}
			s = types.NewSection(types.SectionUnknown, lineNum+1) //section if any starts on the next line
		} else {
			// if s == nil { //if text found before a SectionSeparator (at the start of a script)-> a caption section
			// 	s = types.NewSection(types.SectionCaption, lineNum)
			// }
			s.Lines = append(s.Lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, utils.NewError(utils.ErrSyntaxError, rw.parser.Pos(), err.Error()) //TODO: check parser.Pos()
	}
	tables, err := rw.createTables()
	if err != nil {
		return nil, utils.NewError(utils.ErrSyntaxError, rw.parser.Pos(), err.Error())
	}
	return tables, nil
}

func (rw *RwFile) sectionCount() int {
	return len(rw.sections)
}

func (rw *RwFile) createTables() (tables []*types.Table, err error) {
	if rw.sectionCount() == 0 || rw.sectionCount()%rw.settings.SectionsPerTable != 0 {
		return nil, fmt.Errorf("incorrect number of sections %d", rw.sectionCount())
	}
	var t *types.Table
	for i, s := range rw.sections {
		ii := i + 1 //i is zero-based, section numbers should be one-based
		kind := types.SectionDescriptor(i%rw.settings.SectionsPerTable + 1)
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
			if t.CmdList, err = rw.parser.ParseCommandLines(s); err != nil {
				return nil, err
			}
			tables = append(tables, t)
		default:
			panic(fmt.Sprintf("invalid switch case [%v] in RwFile.CreateTables()", kind))
		}
	}
	return tables, nil
}

func (rw *RwFile) Errors(err error) (eList []error) {
	//	fmt.Printf("%T\n", err)
	switch cause := err.(type) {
	case *utils.EmError:
		//		fmt.Println(cause.Type)
		switch cause.Type {
		case utils.ErrSyntaxError:
			return rw.parser.Errors()
		default:
			return nil
		}
	}
	return nil
}
